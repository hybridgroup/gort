package commands

import (
	"fmt"
	"github.com/codegangsta/cli"
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
				
				fmt.Println("set-udev-rules here...")

			}
		},
	}
}
