package commands

import (
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
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
					ports, _ := ioutil.ReadDir("/dev/serial/by-id/")

					numOfPorts := len(ports)

					if numOfPorts == 0 {
						fmt.Println()
						fmt.Println("No serial ports found.")
						fmt.Println()
						return
					}

					fmt.Println()
					fmt.Println(numOfPorts, "serial port(s) found.")
					fmt.Println()

					for i, port := range ports {
						portPath, _ := filepath.EvalSymlinks("/dev/serial/by-id/" + port.Name())

						deviceInfoPath := "/sys/class/tty/" + filepath.Base(portPath) + "/device/../"

						busNumber, busErr := ioutil.ReadFile(deviceInfoPath + "busnum")
						deviceNumber, devErr := ioutil.ReadFile(deviceInfoPath + "devnum")

						usbDevice := []byte(nil)
						if busErr == nil && devErr == nil {
							usb, usbErr := exec.Command("lsusb", "-s", "00"+string(busNumber)+":00"+string(deviceNumber)).Output()
							if usbErr == nil {
								usbDevice = usb
							}
						}

						fmt.Printf("%d. [%s] - [%s]\n", i+1, portPath, port.Name())

						if usbDevice != nil {
							fmt.Println("  USB device: ", string(usbDevice))
						}

					}
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
					fmt.Println("Connected serialport devices: ")
					fmt.Println(string(out))
				default:
					fmt.Println("Command not available on this OS.")
				}
			default:
				fmt.Println("OS not yet supported.")
			}
		},
	}
}
