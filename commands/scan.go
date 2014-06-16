package commands

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"os/exec"
	"runtime"
)

func Scan() cli.Command {
	return cli.Command{
		Name:  "scan",
		Usage: "Scan for connected devices on Serial, USB, or Bluetooth ports",
		Action: func(c *cli.Context) {
			valid := false
			for _, s := range []string{"serial", "usb", "bluetooth"} {
				if s == c.Args().First() {
					valid = true
				}
			}

			usage := func() {
				fmt.Println("Usage: gort scan [serial|usb|bluetooth]")
			}

			if valid == false {
				fmt.Println("Invalid/no subcommand supplied.\n")
				usage()
				return
			}

			switch runtime.GOOS {
			case "darwin":
				cmd := exec.Command("/bin/sh", "-c", "ls /dev/{tty,cu}.*")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
			case "linux":
				switch c.Args().First() {
				case "serial":
					dmesg := exec.Command("dmesg")
					grep := exec.Command("grep", "tty")
					grep.Stdin, _ = dmesg.StdoutPipe()
					grep.Stdout = os.Stdout
					grep.Start()
					dmesg.Run()
					grep.Wait()
				case "bluetooth":
					cmd := exec.Command("hcitool", "scan")
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					cmd.Run()
				case "usb":
					cmd := exec.Command("lsusb")
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					cmd.Run()
				default:
					fmt.Println("Device type not yet supported.")
				}
			case "windows":
				switch c.Args().First() {
				case "serial":
					out, _ := exec.Command("powershell", "Get-WmiObject Win32_SerialPort", "|", "Select-Object Name, DeviceID, Description").Output()
					fmt.Println(string(out))
				default :
					fmt.Println("Command not available on this OS.")
				}
			default:
				fmt.Println("OS not yet supported.")
			}
		},
	}
}
