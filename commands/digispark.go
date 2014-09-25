package commands

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/codegangsta/cli"
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
				if runtime.GOOS == "linux" {
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
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

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
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

}
