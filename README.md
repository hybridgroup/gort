# Gort - Command Line Interface For RobotOps

[![GitHub release](https://img.shields.io/github/release/hybridgroup/gort.svg)](https://github.com/hybridgroup/gort/releases)
[![Build Status](https://secure.travis-ci.org/hybridgroup/gort.png?branch=master)](http://travis-ci.org/hybridgroup/gort) [![Go Report Card](https://goreportcard.com/badge/github.com/hybridgroup/gort)](https://goreportcard.com/report/github.com/hybridgroup/gort) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/hybridgroup/gort/blob/master/LICENSE)

Gort (http://gort.io) is a Command Line Toolkit for RobotOps. Gort provides tools to scan for connected devices, upload firmware, and more.

Gort is written in the Go programming language (http://golang.org) for maximum speed and portability.

Want to use Golang to program your robots? Check out our open source robotics framework Gobot (http://gobot.io).

Want to use Javascript on Robots? Check out Cylon.js (http://cylonjs.com)

Want to use Ruby on robots? Check out Artoo (http://artoo.io)

## Getting Started
We now have precompiled binaries! You can also build from source.

The Gort CLI provides many useful features on many hardware platforms, and has no other dependencies. You install Gort separately from any framework, which means you can use it to program Arduinos with the Firmata firmware also compatible with Cylon.js, Gobot, Artoo, & many other libraries.

## Download

Just want to download a binary for OSX, Windows, and Linux? Go to our web site at [http://gort.io/documentation/getting_started/downloads/](http://gort.io/documentation/getting_started/downloads/) for the latest release.

Using Homebrew on OSX? You can install using:

```
brew install hybridgroup/tools/gort
```

## How To Use

```
$ ./gort
NAME:
   gort - Command Line Utility for RobotOps

USAGE:
   gort [global options] command [command options] [arguments...]

VERSION:
   0.8.0

COMMANDS:
     scan       Scan for connected devices on Serial, USB, or Bluetooth ports
     bluetooth  Connect & disconnect bluetooth devices.
     arduino    Install avrdude, and upload HEX files to your Arduino
     particle   Upload sketches to your Particle Photon or Electron
     digispark  Configure your Digispark microcontroller
     microbit   Install and upload firmware to your BBC Microbit
     crazyflie  Configure your Crazyflie
     klaatu     barada nikto
     help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

Scan for connected serial devices:

```
$ gort scan serial
[    0.000000] console [tty0] enabled
```

More help coming soon...

## Building

To install the required dependencies to build gort, run:

```
make deps
```

You build the assets into the project like this:

```
make assets
```

Then build the binary with:

```
make build
```

Compiled binaries will now be placed in `build/<VERSION>/`

## Release

When you are ready, run:

```
make release
```

For OSX users, you will also need to update the homebrew repo, located at https://github.com/hybridgroup/homebrew-tools

To obtain the needed SHA values to update homebrew recipe, run:

```
make homebrew
```

## Contributing
For our contribution guidelines, please go to [https://github.com/hybridgroup/gort/blob/master/CONTRIBUTING.md
](https://github.com/hybridgroup/gort/blob/master/CONTRIBUTING.md
).

## Licenses
Gort is copyright (c) 2014-2017 The Hybrid Group. Licensed under the Apache 2.0 license.

Firmata is copyright (c) 2006-2008 Hans-Christoph Steiner. Licensed under GNU Lesser General Public License. All rights reserved.

Rapiro is copyright (c) 2013-2014 Shota Ishiwatari. Licensed under the Creative Commons - Public Domain Dedication License.

Tinker-Servo is copyright (c) 2014-2017 Particle, Scott Beasley. Licensed under the GNU General Public License.
