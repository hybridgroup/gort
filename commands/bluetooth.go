package commands

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/codegangsta/cli"
)

func Bluetooth() cli.Command {
	return cli.Command{
		Name:  "bluetooth",
		Usage: "Pair, unpair & connect to bluetooth devices.",
		Action: func(c *cli.Context) {
			valid := false
			for _, s := range []string{"pair", "unpair", "connect"} {
				if s == c.Args().First() {
					valid = true
				}
			}
			usage := func() {
				fmt.Println("Invalid/no subcommand supplied.\n")
				fmt.Println("Usage:")
				fmt.Println("gort bluetooth pair <address> [hciX]")
				fmt.Println("gort bluetooth unpair <address> [hciX]")
				fmt.Println("gort bluetooth connect <dev> [hciX]\n")
			}

			if runtime.GOOS == "darwin" {
				fmt.Println("OS X manages Bluetooth pairing/unpairing/binding itself.")
				return
			}

			if len(c.Args()) < 2 {
				valid = false
			}

			if valid == false {
				usage()
				return
			}

			hci := "hci0"
			if len(c.Args()) <= 3 {
				hci = c.Args()[2]
			}

			switch runtime.GOOS {
			case "linux":
				switch c.Args().First() {
				case "pair":
					cmd := exec.Command("bluez-simple-agent", hci, c.Args()[1])
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					if err := cmd.Run(); err != nil {
						log.Fatal(err)
					}

				case "unpair":
					cmd := exec.Command("bluez-simple-agent", hci, c.Args()[1], "remove")
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					if err := cmd.Run(); err != nil {
						log.Fatal(err)
					}

				case "connect":
					cmd := exec.Command("bluez-test-serial", "-i", hci, c.Args()[1])
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					if err := cmd.Run(); err != nil {
						log.Fatal(err)
					}

				default:
					usage()
				}
			default:
				fmt.Println("OS not yet supported.")
			}
		},
	}
}
