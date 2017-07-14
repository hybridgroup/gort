package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/hybridgroup/gort/commands"
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
		commands.Workshop(),
	}
	app.Run(os.Args)
}
