package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/codegangsta/cli"
)

func Arduino() cli.Command {
	return cli.Command{
		Name:  "arduino",
		Usage: "Install avrdude, and upload HEX files to your Arduino",
		Action: func(c *cli.Context) {
			valid := false
			for _, s := range []string{"scan", "install", "upload"} {
				if s == c.Args().First() {
					valid = true
				}
			}

			usage := func() {
				fmt.Println("Invalid/no subcommand supplied.\n")
				fmt.Println("Usage:")
				fmt.Println("  gort arduino install                                  # installs avrdude to allow uploading of sketches to Arduino")
				fmt.Println("  gort arduino upload firmata [port]                    # uploads Firmata sketch to Arduino")
				fmt.Println("  gort arduino upload rapiro [port]                     # uploads Rapiro sketch to Arduino")
				fmt.Println("  gort arduino upload [custom-firmware-filename] [port] # uploads a custom sketch to Arduino")
			}

			if valid == false {
				usage()
				return
			}

			switch c.Args().First() {
			case "install":
				switch runtime.GOOS {
				case "linux":
					fmt.Println("Attempting to install avrdude with apt-get.")
					cmd := exec.Command("sudo", "apt-get", "-y", "install", "avrdude")
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					if err := cmd.Run(); err != nil {
						log.Fatal(err)
					}

				case "darwin":
					fmt.Println("Attempting to install avrdude with Homebrew.")
					cmd := exec.Command("brew", "install", "avrdude")
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					if err := cmd.Run(); err != nil {
						log.Fatal(err)
					}

				case "windows":
					_, err := exec.Command("NET", "SESSION").Output()
					if err != nil {
						fmt.Println("Please run cmd.exe as administrator and try again")
						os.Exit(1)
					}

					fmt.Println("Installing winavr")
					dirName, _ := createGortDirectory()
					exeFile := "https://s3.amazonaws.com/gort-io/support/WinAVR-20100110-install.exe"
					fileName := downloadFromUrl(dirName, exeFile)
					cmd := exec.Command(gortDirName() + "\\" + fileName)
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					if err := cmd.Run(); err != nil {
						log.Fatal(err)
					}

				default:
					fmt.Println("OS not yet supported.")
				}
			case "upload":
				if len(c.Args()) < 2 {
					fmt.Println("Invalid number of arguments.")
					usage()
					return
				}

				hexfile := c.Args()[1]
				port := c.Args()[2]
				file, _ := ioutil.TempFile(os.TempDir(), "")
				defer file.Close()
				defer os.Remove(file.Name())

				if hexfile == "firmata" || hexfile == "rapiro" {
					hexfile = fmt.Sprintf("support/arduino/%v.cpp.hex", hexfile)
					data, _ := Asset(hexfile)
					file.Write(data)
					file.Sync()
					hexfile = file.Name()
				}

				switch runtime.GOOS {
				case "darwin", "linux", "windows":
					cmd := exec.Command("avrdude", "-patmega328p", "-carduino", fmt.Sprintf("-P%v", port), "-b115200", "-D", fmt.Sprintf("-Uflash:w:%v:i", hexfile))
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					if err := cmd.Run(); err != nil {
						log.Fatal(err)
					}

				default:
					fmt.Println("OS not yet supported.")
				}
			}
		},
	}
}
