package commands

import (
	"fmt"
	"github.com/codegangsta/cli"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
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
				extractDigisparkInstaller()
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

				fmt.Println("upload here...")

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
		downloadFromUrl(dirName, "http://littlewire.cc/resources/LittleWirev13Install-Linux64.tar.gz")
	case "darwin":
		downloadFromUrl(dirName, "http://littlewire.cc/resources/LittleWirev13Install-OSX.zip")
	default:
		downloadFromUrl(dirName, "http://littlewire.cc/resources/LittleWirev13Install-Win.zip")
	}
}

func extractDigisparkInstaller() {
	fmt.Println("extract digispark installer here...")
}

func runDigisparkInstaller() {
	fmt.Println("run digispark installer here...")
}

func downloadFromUrl(dirName string, url string) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	fmt.Println("Downloading", url, "to", fileName)

	// TODO: check file existence first with io.IsExist
	output, err := os.Create(dirName + "/" + fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}

	fmt.Println(n, "bytes downloaded.")
}

func createGortDirectory() (string, error) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	dirName := usr.HomeDir + "/" + "gort"
	fileExists, err := exists(dirName)
	if fileExists {
		fmt.Println("Gort lives")
	} else {
		os.Mkdir(dirName, 0777)
	}
	return dirName, err
}
