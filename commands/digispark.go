package commands

import (
	"fmt"
	"github.com/codegangsta/cli"
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

				fmt.Println("set-udev-rules here...")

			}
		},
	}
}
