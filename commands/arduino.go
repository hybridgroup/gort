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

func uploadParams(board string) (string, string) {
	switch board {
	case "m328", "uno", "nano", "mini", "ethernet", "fio":
		return "arduino", "m328p"

	case "m168", "diecimila", "stamp":
		return "arduino", "m168"

	case "mega", "mega1280":
		return "arduino", "m1280"

	case "mega2560", "megaADK":
		return "stk500v2", "m2560"

	case "leonardo", "robot", "micro", "esplora":
		return "avr109", "atmega32u4"

	default:
		return "arduino", "m328p"
	}
}

func Arduino() cli.Command {
	return cli.Command{
		Name:  "arduino",
		Usage: "Install avrdude, and upload HEX files to your Arduino",
		Flags: []cli.Flag {
		  cli.StringFlag{
		    Name: "board, b",
		    Value: "uno",
		    Usage: "board type of arduino",
		  },
		},
		Action: func(c *cli.Context) {
			valid := false
			for _, s := range []string{"scan", "install", "upload"} {
				if s == c.Args().First() {
					valid = true
				}
			}

			usage := func() {
				fmt.Println("Invalid/no subcommand supplied.")
				fmt.Println()
				fmt.Println("Usage:")
				fmt.Println("  # installs avrdude to allow uploading of sketches to Arduino")
				fmt.Println("  gort arduino install")
				fmt.Println()
				fmt.Println("  # uploads Firmata sketch to Arduino")
				fmt.Println("  gort arduino upload firmata <port> [flags]")
				fmt.Println()
				fmt.Println("  # uploads Rapiro sketch to Arduino")
				fmt.Println("  gort arduino upload rapiro <port> [flags]")
				fmt.Println()
				fmt.Println("  # uploads a custom sketch to Arduino")
				fmt.Println("  gort arduino upload <custom-firmware-filename> <port> [flags]")
				fmt.Println()
				fmt.Println("    upload flags:")
				fmt.Println("      -b < m328 | uno | nano | mini | ethernet | fio | m168 |")
				fmt.Println("           diecimila | stamp | mega | mega1280 | mega2560 | megaADK |")
				fmt.Println("           leonardo | robot | micro | esplora >")
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
				programmer, part := uploadParams(c.String("board"))
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
					cmd := exec.Command("avrdude", fmt.Sprintf("-p%v", part), fmt.Sprintf("-c%v", programmer), fmt.Sprintf("-P%v", port), "-D", fmt.Sprintf("-Uflash:w:%v:i", hexfile))
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
