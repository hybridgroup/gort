package commands

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/codegangsta/cli"
)

// Microbit function returns the CLI commands for gort microbit
func Microbit() cli.Command {
	return cli.Command{
		Name:  "microbit",
		Usage: "Install and upload firmware to your BBC Microbit",
		Action: func(c *cli.Context) {
			valid := false
			for _, s := range []string{"install", "download"} {
				if s == c.Args().First() {
					valid = true
				}
			}

			usage := func() {
				fmt.Println("Invalid/no subcommand supplied.")
				fmt.Println()
				fmt.Println("Usage:")
				fmt.Println("  # downloads latest microbit firmware file")
				fmt.Println("  gort microbit download")
				fmt.Println()
				fmt.Println("  # installs firmware on microbit")
				fmt.Println("  gort microbit install <path>")
			}

			if valid == false {
				usage()
				return
			}

			switch c.Args().First() {
			case "download":
				fmt.Println("Downloading microbit firmware...")
				dirName, _ := createGortDirectory()
				hexFile := "https://github.com/sandeepmistry/node-bbc-microbit/raw/master/firmware/node-bbc-microbit-v0.1.0.hex"
				fileName := downloadFromUrl(dirName, hexFile)
				fmt.Println(dirName, fileName)

			case "install":
				if len(c.Args()) < 2 {
					fmt.Println("Invalid number of arguments.")
					usage()
					return
				}

				fmt.Println("Installing microbit firmware...")
				targetDir := c.Args()[1]
				fileName := gortDirName() + "/node-bbc-microbit-v0.1.0.hex"
				data, _ := ioutil.ReadFile(fileName)

				hexFile, _ := os.Create(targetDir + "/node-bbc-microbit-v0.1.0.hex")
				defer hexFile.Close()

				hexFile.Write(data)
				hexFile.Sync()
			}
		},
	}
}
