package commands

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"runtime"
)

func Digispark() cli.Command {
	return cli.Command{
		Name:  "digispark",
		Usage: "Upload sketches to your Digispark",
		Action: func(c *cli.Context) {
			valid := false
			for _, s := range []string{"upload", "set-udev-rules"} {
				if s == c.Args().First() {
					valid = true
				}
			}

			usage := func() {
				fmt.Println("Invalid/no subcommand supplied.\n")
				fmt.Println("Usage:")
				fmt.Println("  gort digispark upload [default|path name] # uploads sketch to Digispark")
				fmt.Println("  gort digispark set-udev-rules # set udev rules needed to connect to Digispark")
			}

			if valid == false {
				usage()
				return
			}

			switch c.Args().First() {
			case "upload":

				if len(c.Args()) < 3 {
					fmt.Println("Invalid number of arguments.")
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
