package commands

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"

	"os"

	"github.com/codegangsta/cli"
)

// Workshop uploads workshop specific sketch for firmta.
func Workshop() cli.Command {
	return cli.Command{
		Name:  "workshop",
		Usage: "Performs workshop realated Arduino commands.",
		Subcommands: []cli.Command{
			{
				Name:  "upload",
				Usage: "upload a sketch via Arduino cli",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "board",
						Usage: "board to flash",
						Value: "Arduino101",
					},
					cli.StringFlag{
						Name:  "port",
						Usage: "port to use",
					},
					cli.StringFlag{
						Name:  "sketch",
						Usage: "sketch to upload",
					},
				},
				Action: func(c *cli.Context) {
					// usage := func() {
					// 	fmt.Println("")
					// }

					//valid := true

					var arduino string

					var board string
					port := c.String("port")
					sketch := c.String("sketch")

					switch runtime.GOOS {
					case "windows":
						arduino = "arduino_debug"
					case "darwin":
						arduino = "Arduino.app/Contents/MacOS/Arduino"
					case "linux":
						arduino = "arduino"
					}

					switch c.String("board") {
					case "Arduino101", "TinyTile":
						board = "Intel:arc32:arduino_101"
					default:
						//valid = false
					}

					fmt.Printf("Arduino upload parameters: %v, %v, %v", board, port, sketch)

					cmd := exec.Command(arduino, "--upload", "--board", board, "--port", port, sketch)
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					if err := cmd.Run(); err != nil {
						log.Fatal(err)
					}
				},
			},
		},
	}
}
