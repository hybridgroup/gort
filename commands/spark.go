package commands

import (
	"fmt"
	"github.com/codegangsta/cli"
)

func Spark() cli.Command {
	return cli.Command{
		Name:  "spark",
		Usage: "Upload sketches to your Spark",
		Action: func(c *cli.Context) {
			valid := false
			for _, s := range []string{"upload"} {
				if s == c.Args().First() {
					valid = true
				}
			}

			usage := func() {
				fmt.Println("Invalid/no subcommand supplied.\n")
				fmt.Println("Usage:")
				fmt.Println("  gort spark upload [accessToken] [deviceId] [default|path name] # uploads sketch to Spark")
			}

			if valid == false {
				usage()
				return
			}

			switch c.Args().First() {
			case "upload":

				if len(c.Args()) < 5 {
					fmt.Println("Invalid number of arguments.")
					usage()
					return
				}

			}
		},
	}
}
