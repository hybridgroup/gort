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
		Usage: "Connect & disconnect bluetooth devices.",
		Action: func(c *cli.Context) {
			valid := false
			for _, s := range []string{"connect", "serial", "disconnect"} {
				if s == c.Args().First() {
					valid = true
				}
			}
			usage := func() {
				fmt.Println("Invalid/no subcommand supplied.")
				fmt.Println()
				fmt.Println("Usage:")
				fmt.Println("gort bluetooth connect <address> [hciX]")
				fmt.Println("gort bluetooth serial <address> [port] [hciX]")
				fmt.Println("gort bluetooth disconnect <address> [hciX]")
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

			switch runtime.GOOS {
			case "linux":
				switch c.Args().First() {
				case "connect":
					if len(c.Args()) >= 3 {
						hci = c.Args()[2]
					}

					cmd := exec.Command("hcitool", "-i", hci, "cc", c.Args()[1])
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					if err := cmd.Run(); err != nil {
						log.Fatal(err)
					}

				case "serial":
					port := "rfcomm0"

					if len(c.Args()) >= 4 {
						hci = c.Args()[3]
					}

					if len(c.Args()) >= 3 {
						port = c.Args()[2]
					}

					cmd := exec.Command("rfcomm", "-i", hci, "connect", port, c.Args()[1])
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					if err := cmd.Run(); err != nil {
						log.Fatal(err)
					}

				case "disconnect":
					if len(c.Args()) >= 3 {
						hci = c.Args()[2]
					}

					cmd := exec.Command("hcitool", "-i", hci, "dc", c.Args()[1])
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
