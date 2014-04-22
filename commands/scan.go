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
		Usage: "Scans serial, Bluetooth, or USB for connected devices",
		Action: func(c *cli.Context) {
			valid := false
			for _, s := range []string{"serial", "bluetooth", "usb"} {
				if s == c.Args().First() {
					valid = true
				}
			}

			if valid == false {
				fmt.Println("Invalid/no type supplied.")
				fmt.Println("Usage: gort scan [serial|bluetooth|usb]")
				return
			}

			switch runtime.GOOS {
			case "darwin":
				cmd := exec.Command("/bin/sh", "-c", "ls /dev/tty*")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()

				cmd = exec.Command("/bin/sh", "-c", "ls /dev/cu*")
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
			default:
				fmt.Println("OS not yet supported.")
			}
		},
	}
}
