package commands

import (
	"archive/zip"
	"fmt"
	"github.com/codegangsta/cli"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

func Digispark() cli.Command {
	return cli.Command{
		Name:  "digispark",
		Usage: "Configure your Digispark microcontroller",
		Action: func(c *cli.Context) {
			valid := false
			for _, s := range []string{"install", "upload", "set-udev-rules"} {
				if s == c.Args().First() {
					valid = true
				}
			}

			usage := func() {
				fmt.Println("Usage:")
				fmt.Println("  gort digispark install # installs software to upload firmware to Digispark")
				fmt.Println("  gort digispark upload [littlewire] # uploads firmware to Digispark")
				fmt.Println("  gort digispark set-udev-rules # set udev rules needed to connect to Digispark")
			}

			if valid == false {
				fmt.Println("Invalid/no subcommand supplied.\n")
				usage()
				return
			}

			switch c.Args().First() {
			case "install":
				downloadDigisparkInstaller()
				runDigisparkInstaller()
				return

			case "upload":
				if len(c.Args()) != 2 {
					fmt.Println("Invalid number of arguments.")
					usage()
					return
				}

				firmware := c.Args()[1]
				if firmware != "littlewire" {
					fmt.Println("Unknown firmware.")
					usage()
					return
				}

				runDigisparkInstaller()
				return

			case "set-udev-rules":
				if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
					digisparkSetUdevRules()
				} else {
					fmt.Println("No need to set-udev-rules on your OS")
				}
			}
		},
	}
}

func digisparkSetUdevRules() {
	fileExists, _ := exists("/etc/udev/rules.d/49-micronucleus.rules")
	if !fileExists {
		file, err := os.Create("/etc/udev/rules.d/49-micronucleus.rules")
		if err != nil {
			if os.IsPermission(err) {
				fmt.Println("You do not have the required permissions. Try running this command using 'sudo'")
			} else {
				fmt.Println(err)
			}
			return
		}
		defer file.Close()

		data, _ := Asset("support/digispark/micronucleus.rules")
		file.Write(data)
		file.Sync()

		if err != nil {
			fmt.Println("Error while copying Digispark udev rules")
			return
		}
	}
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func downloadDigisparkInstaller() {
	dirName, _ := createGortDirectory()
	switch runtime.GOOS {
	case "linux":
		zipFile := "http://littlewire.cc/resources/LittleWirev13Install-Linux64.tar.gz"
		fileName := downloadFromUrl(dirName, zipFile)
		extractDigisparkInstaller(dirName, dirName+"/"+fileName)
	case "darwin":
		zipFile := "http://littlewire.cc/resources/LittleWirev13Install-OSX.zip"
		fileName := downloadFromUrl(dirName, zipFile)
		unzipDigisparkInstaller(dirName, dirName+"/"+fileName)
	default:
		zipFile := "http://littlewire.cc/resources/LittleWirev13Install-Win.zip"
		fileName := downloadFromUrl(dirName, zipFile)
		unzipDigisparkInstaller(dirName, dirName+"/"+fileName)
	}
}

func extractDigisparkInstaller(dirName string, zipFile string) {
	cmd := exec.Command("tar", "-C", dirName, "-zxvf", zipFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func unzipDigisparkInstaller(dirName string, zipFile string) {
	err := Unzip(zipFile, dirName)
	if err != nil {
		log.Fatal(err)
	}
}

func runDigisparkInstaller() {
	cmd := new(exec.Cmd) 
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command(gortDirName() + "/littlewireLoader_v13")
	case "darwin":
		cmd = exec.Command(gortDirName() + "/LittleWirev13Install")
	default:
		cmd = exec.Command(gortDirName() + "/LittleWirev13Install.exe")
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func downloadFromUrl(dirName string, url string) string {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	fmt.Println("Downloading", url, "to", fileName)

	// TODO: check file existence first with io.IsExist
	output, err := os.Create(dirName + "/" + fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return ""
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return ""
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return ""
	}

	fmt.Println(n, "bytes downloaded.")
	return fileName
}

func gortDirName() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir + "/" + "gort"
}

func createGortDirectory() (string, error) {
	dirName := gortDirName()
	fileExists, err := exists(dirName)
	if fileExists {
		fmt.Println("Gort lives")
	} else {
		os.Mkdir(dirName, 0777)
	}
	return dirName, err
}

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			f, err := os.OpenFile(
				path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
