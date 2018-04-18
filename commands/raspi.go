package commands

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/codegangsta/cli"
)

// Raspi command implementation
func Raspi() cli.Command {

	return cli.Command{
		Name:  "raspi",
		Usage: "Download and install components for your Raspberry Pi",
		Action: func(c *cli.Context) {
			if runtime.GOOS != "linux" {
				fmt.Println("Raspberry Pis commands must be run on the device, current OS not supported")
				return
			}

			valid := true

			if c.Args().First() != "install" {
				valid = false
			} else if len(c.Args()) < 2 {
				valid = false
			}

			usage := func() {
				fmt.Println("Invalid/no subcommand supplied.")
				fmt.Println()
				fmt.Println("Usage:")
				fmt.Println("  # downloads and installs piblaster")
				fmt.Println("  gort raspi install piblaster")
			}

			if valid == false {
				usage()
				return
			}

			subCommand := c.Args()[1]

			switch subCommand {
			case "piblaster":
				raspiInstallPiBlaster()
			default:
				usage()
			}

		},
	}
}

func raspiInstallPiBlaster() (err error) {
	dir, _ := createGortDirectory()

	var cmd *exec.Cmd

	if !linuxCommandExists("autoconf") {
		fmt.Println("Attempting to install dev tools")
		switch getLinuxDist() {
		case "arch":
			cmd = exec.Command("sudo", "pacman", "-S", "--noconfirm", "autoconf")
		case "fedora":
			cmd = exec.Command("sudo", "yum", "-y", "install", "autoconf")
		case "ubuntu":
			cmd = exec.Command("sudo", "apt-get", "-y", "install", "autoconf")
		default:
			return fmt.Errorf("Unable to detect Linux distribution. Try installing autoconf manually first")
		}
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Attempting to fetch 'pi-blaster' from github.")
	cmd = exec.Command("wget", "https://github.com/sarfata/pi-blaster/archive/master.zip")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Attempting to unzip 'pi-blaster'.")
	cmd = exec.Command("unzip", "master.zip")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Preparing to build 'pi-blaster'.")
	cmd = exec.Command("./autogen.sh")
	cmd.Dir = fmt.Sprintf("%v/%v", dir, "pi-blaster-master")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Attempting to build 'pi-blaster'.")
	cmd = exec.Command("./configure && make")
	cmd.Dir = fmt.Sprintf("%v/%v", dir, "pi-blaster-master")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Attempting to install 'pi-blaster'.")
	cmd = exec.Command("sudo make install")
	cmd.Dir = fmt.Sprintf("%v/%v", dir, "pi-blaster-master")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	return
}
