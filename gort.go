package main

import (
	"os"

	"github.com/hybridgroup/gort/commands"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gort"
	app.Version = VERSION
	app.Usage = "Command Line Utility for RobotOps"
	app.Commands = []cli.Command{
		commands.Scan(),
		commands.Bluetooth(),
		commands.Arduino(),
		commands.Raspi(),
		commands.Chip(),
		commands.Particle(),
		commands.Digispark(),
		commands.Microbit(),
		commands.Crazyflie(),
		commands.Klaatu(),
	}
	app.Run(os.Args)
}
