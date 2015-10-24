package main

import (
	"github.com/codegangsta/cli"
	"github.com/hybridgroup/gort/commands"
	"os"
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
	}
	app.Run(os.Args)
}
