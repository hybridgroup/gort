package commands

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"runtime"
)

func Crazyflie() cli.Command {
	return cli.Command{
		Name:  "crazyflie",
		Usage: "Configure your Crazyflie",
		Action: func(c *cli.Context) {
			valid := false
			for _, s := range []string{"set-udev-rules"} {
				if s == c.Args().First() {
					valid = true
				}
			}

			usage := func() {
				fmt.Println("Invalid/no subcommand supplied.\n")
				fmt.Println("Usage:")
				fmt.Println("  gort crazyflie set-udev-rules # set udev rules needed to connect to Crazyflie")
			}

			if valid == false {
				usage()
				return
			}

			switch c.Args().First() {
			case "set-udev-rules":
				if runtime.GOOS == "linux" {
					crazyflieSetUdevRules()
				} else {
					fmt.Println("No need to set-udev-rules on your OS")
				}

			}
		},
	}
}

func crazyflieSetUdevRules() {
	fileExists, _ := exists("/etc/udev/rules.d/99-crazyradio.rules")
	if !fileExists {
		file, err := os.Create("/etc/udev/rules.d/99-crazyradio.rules")
		if err != nil {
			if os.IsPermission(err) {
				fmt.Println("You do not have the required permissions. Try running this command using 'sudo'")
			} else {
				fmt.Println(err)
			}
			return
		}
		defer file.Close()

		data, _ := Asset("support/crazyflie/crazyradio.rules")
		file.Write(data)
		file.Sync()

		if err != nil {
			fmt.Println("Error while copying Crazyradio udev rules")
			return
		}
	}
}
