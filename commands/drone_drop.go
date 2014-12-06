package commands

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/ziutek/telnet"

	"github.com/codegangsta/cli"
)

func DroneDrop() cli.Command {
	return cli.Command{
		Name:  "dronedrop",
		Usage: "Install, uninstall, update and download dronedrop firmware",
		Action: func(c *cli.Context) {
			valid := false
			for _, s := range []string{"install", "uninstall", "update", "download"} {
				if s == c.Args().First() {
					valid = true
				}
			}

			usage := func() {
				fmt.Println("Usage:")
				fmt.Println("  gort dronedrop download  # Download the latest dronedrop firmware and support scripts")
				fmt.Println("  gort dronedrop install   # Install dronedrop on your ardrone")
				fmt.Println("  gort dronedrop uninstall # Remove dronedrop from your ardrone ")
				fmt.Println("  gort dronedrop update    # Updates dronedrop on your ardrone")
			}

			if valid == false {
				fmt.Println("Invalid/no subcommand supplied.\n")
				usage()
				return
			}

			switch c.Args().First() {
			case "download":
				downloadDroneDropFirmware()
				return

			case "install":
				installDroneDrop()
				return

			case "uninstall":
				uninstallDroneDrop()
				return

			case "update":
				uninstallDroneDrop()
				installDroneDrop()
				return
			}
		},
	}
}

func droneDropSupportDir() string {
	dir, err := supportDir("dronedrop")

	if err != nil {
		log.Fatal(err)
	}

	return dir
}

func downloadDroneDropFirmware() {
	dir := droneDropSupportDir()
	downloadFromUrl(dir, "https://github.com/hybridgroup/dronedrop-ardrone/raw/master/ardrone_commander")
	downloadFromUrl(dir, "https://github.com/hybridgroup/dronedrop-ardrone/raw/master/configure_drone_drop.sh")
	downloadFromUrl(dir, "https://github.com/hybridgroup/dronedrop-ardrone/raw/master/uninstall_drone_drop.sh")
}

func uninstallDroneDrop() {
	log.Println("Uninstalling dronedrop...")
	sendCommand("sh /data/video/uninstall_drone_drop.sh\n")
	<-time.After(5 * time.Second)
}
func installDroneDrop() {
	c, err := ftp.Connect("192.168.1.1:21")
	defer c.Quit()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Pushing dronedrop...")
	for _, asset := range []string{"/ardrone_commander", "/configure_drone_drop.sh", "/uninstall_drone_drop.sh"} {
		if b, err := os.Open(droneDropSupportDir() + asset); err != nil {
			log.Fatal(err)
		} else {
			if err := c.Stor(asset, b); err != nil {
				log.Fatal("Error pushing", asset, err)
			}
		}
	}

	for _, asset := range []string{"/rcS.stock", "/rcS.dronedrop", "/usb.ids.stock", "/usb.ids.dronedrop"} {
		if b, err := Asset("support/drone_drop" + asset); err != nil {
			log.Fatal(err)
		} else {
			if err := c.Stor(asset, bytes.NewBuffer(b)); err != nil {
				log.Fatal("Error pushing", asset, err)
			}
		}
	}

	log.Println("Configuring dronedrop...")
	sendCommand("sh /data/video/configure_drone_drop.sh\n")
	log.Println("Rebooting drone...")
}

func sendCommand(command string) {
	t, err := telnet.Dial("tcp", "192.168.1.1:23")
	defer t.Close()

	if err != nil {
		log.Fatal(err)
	}
	t.SetUnixWriteMode(true)

	if _, err := t.Write([]byte(command)); err != nil {
		log.Fatal(err)
	}
}
