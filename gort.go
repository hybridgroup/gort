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
		commands.Spark(),
		commands.Digispark(),
		commands.Crazyflie(),
		commands.Klaatu(),
		commands.DroneDrop(),
	}
	app.Run(os.Args)
}
