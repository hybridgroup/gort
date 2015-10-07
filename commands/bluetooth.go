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
			for _, s := range []string{"pair", "connect", "unpair"} {
				if s == c.Args().First() {
					valid = true
				}
			}
			usage := func() {
				fmt.Println("Invalid/no subcommand supplied.")
				fmt.Println("Usage:")
				fmt.Println("gort bluetooth pair <address> [hciX]")
				fmt.Println("gort bluetooth connect <port> <address> [hciX]")
				fmt.Println("gort bluetooth unpair <address> [hciX]")
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
			if len(c.Args()) >= 3 {
				hci = c.Args()[2]
			}

			switch runtime.GOOS {
			case "linux":
				switch c.Args().First() {
				case "pair":
					if len(c.Args()) >= 3 {
						hci = c.Args()[2]
					}

					cmd := exec.Command("hcitool", "-i", hci, "cc", c.Args()[1])
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					if err := cmd.Run(); err != nil {
						log.Fatal(err)
					}

				case "connect":
					if len(c.Args()) >= 4 {
						hci = c.Args()[3]
					}

					cmd := exec.Command("rfcomm", "-i", hci, "connect", c.Args()[1], c.Args()[2])
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					if err := cmd.Run(); err != nil {
						log.Fatal(err)
					}

				case "unpair":
					if len(c.Args()) >= 3 {
						hci = c.Args()[2]
					}

					cmd := exec.Command("hcitool", "-i", hci, "cc", c.Args()[1])
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
