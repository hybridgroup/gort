package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/codegangsta/cli"
)

// Chip command implementation
func Chip() cli.Command {

	return cli.Command{
		Name:  "chip",
		Usage: "Download and install device tree components for your C.H.I.P",
		Action: func(c *cli.Context) {
			if runtime.GOOS != "linux" {
				fmt.Println("C.H.I.P commands must be run on the device, current OS not supported")
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
				fmt.Println("  # downloads and installs the required 'dtc' tool")
				fmt.Println("  gort chip install dtc")
				fmt.Println()
				fmt.Println("  # installs the device tree overlay for PWM0")
				fmt.Println("  # requires an installed 'dtc' tool, and most likely 'sudo'")
				fmt.Println("  gort chip install pwm")
				fmt.Println()
				fmt.Println("  # installs the device tree overlay for SPI2")
				fmt.Println("  # requires an installed 'dtc' tool, and most likely 'sudo'")
				fmt.Println("  gort chip install spi")
			}

			if valid == false {
				usage()
				return
			}

			subCommand := c.Args()[1]

			switch subCommand {
			case "pwm":
				info := chipOverlayInfo{"pwm", "chip-pwm0.dtbo", "chip-pwm", "support/chip/chip-pwm0.dts"}
				err := chipBuildAndInstallOverlay(info)
				if err != nil {
					log.Fatal(err)
				}
				err = chipUpdateInitscript(info)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("PWM0 overlay installation complete.")
				fmt.Println("Now reboot to make sure init scripts load the overlay as expected.")
			case "spi":
				info := chipOverlayInfo{"spi", "chip-spi2.dtbo", "chip-spi", "support/chip/chip-spi2.dts"}
				err := chipBuildAndInstallOverlay(info)
				if err != nil {
					log.Fatal(err)
				}
				err = chipUpdateInitscript(info)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("SPI2 overlay installation complete.")
				fmt.Println("Now reboot to make sure init scripts load the overlay as expected.")
			case "dtc":
				chipInstallDTC()
			default:
				usage()
			}

		},
	}
}

const chipOverlayInstallPath = "/lib/firmware/gort.io"
const chipOverlayConfigPath = "/sys/kernel/config/device-tree/overlays"

type chipOverlayInfo struct {
	id     string
	dtbo   string
	folder string
	source string
}

func chipBuildOverlay(tmpDir string, source string) (err error) {
	path, err := exec.LookPath("dtc")
	if err != nil {
		return fmt.Errorf("Failed to find 'dtc' command on path")
	}

	sourcePath := tmpDir + "/overlay.dts"
	blobPath := tmpDir + "/overlay.dtbo"

	if err = ioutil.WriteFile(sourcePath, []byte(source), 0666); err != nil {
		return err
	}

	dtc := exec.Command(path, "-O", "dtb", "-o", blobPath, "-b", "o", "-@", sourcePath)
	if err = dtc.Run(); err != nil {
		return err
	}

	return nil
}

func chipInstallOverlay(tmpDir string, blobFile string) (err error) {
	err = os.MkdirAll(chipOverlayInstallPath, 0777)
	if err == nil {
		blobSource := tmpDir + "/overlay.dtbo"
		blobTarget := chipOverlayInstallPath + "/" + blobFile
		err = copyFile(blobSource, blobTarget)
	}

	if os.IsPermission(err) {
		return fmt.Errorf("No permission to install overlay at %q, retry using 'sudo'", chipOverlayInstallPath)
	}

	return err
}

func chipBuildAndInstallOverlay(info chipOverlayInfo) (err error) {
	blobPath := chipOverlayInstallPath + "/" + info.dtbo
	if _, err = os.Stat(blobPath); err == nil {
		// this blob is in place, return
		return
	}
	if os.IsNotExist(err) {
		var tmpDir string
		tmpDir, err = ioutil.TempDir("", "dtbo")
		defer os.RemoveAll(tmpDir)

		code, _ := Asset(info.source)
		err = chipBuildOverlay(tmpDir, string(code))
		if err != nil {
			return fmt.Errorf("Failed to build overlay: %v", err)
		}
		err = chipInstallOverlay(tmpDir, info.dtbo)
		if err != nil {
			return fmt.Errorf("Failed to install overlay: %v", err)
		}
	} else {
		return fmt.Errorf("Failed to check for installed overlay: %v", err)
	}

	return
}

func chipInstallDTC() (err error) {
	dir, _ := createGortDirectory()

	fmt.Println("Attempting to install dev tools with apt-get.")
	cmd := exec.Command("sudo", "apt-get", "-y", "install", "flex", "bison", "git")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Attempting to fetch 'dtc' from github.")
	cmd = exec.Command("git", "clone", "https://github.com/atenart/dtc")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Attempting to build 'dtc'.")
	cmd = exec.Command("make")
	cmd.Dir = fmt.Sprintf("%v/%v", dir, "dtc")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Attempting to install 'dtc'.")
	cmd = exec.Command("sudo", "make", "install", "PREFIX=/usr")
	cmd.Dir = fmt.Sprintf("%v/%v", dir, "dtc")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	return
}

func chipUpdateInitscript(info chipOverlayInfo) (err error) {
	const initScript = "/etc/rc.local"

	configPath := chipOverlayConfigPath + "/" + info.folder
	overlayPath := chipOverlayInstallPath + "/" + info.dtbo

	prefix := fmt.Sprintf("# GORT BEGIN %s\n", info.id)
	action := fmt.Sprintf("mkdir -p %s\ncp %s %s/dtbo\n", configPath, overlayPath, configPath)
	suffix := fmt.Sprintf("# GORT END %s\n", info.id)

	initSegment := prefix + action + suffix

	scriptBytes, err := ioutil.ReadFile(initScript)
	if err != nil {
		return
	}

	script := string(scriptBytes)
	if strings.Contains(script, prefix) {
		// assume previously edited script
		headAndRest := strings.Split(script, prefix)
		head, rest := headAndRest[0], headAndRest[1]
		if strings.Contains(rest, suffix) {
			midTail := strings.Split(rest, suffix)
			_, tail := midTail[0], midTail[1]

			script = head + initSegment + tail
		} else {
			return fmt.Errorf("%q init script GORT section malformed. Manual editing required", initScript)
		}
	} else {
		// assume unedited script, should end with "exit 0"
		const exitLine = "exit 0\n"
		if strings.HasSuffix(script, exitLine) {
			script = strings.TrimSuffix(script, exitLine)
			script = script + initSegment + exitLine
		} else {
			return fmt.Errorf("%q init script editing failed. Manual editing required", initScript)
		}
	}

	err = ioutil.WriteFile(initScript, []byte(script), 0755)

	if err != nil {
		return fmt.Errorf("Failed to write edited init script: %v", err)
	}

	return
}
