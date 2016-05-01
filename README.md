# Gort - Command Line Interface For RobotOps

Gort (http://gort.io) is a Command Line Toolkit for RobotOps. Gort provides tools to scan for connected devices, upload firmware, and more.

Gort is written in the Go programming language (http://golang.org) for maximum speed and portability.

Want to use Golang to program your robots? Check out our open source robotics framework Gobot (http://gobot.io).

Want to use Javascript on Robots? Check out Cylon.js (http://cylonjs.com)

Want to use Ruby on robots? Check out Artoo (http://artoo.io)

[![Build Status](https://secure.travis-ci.org/hybridgroup/gort.png?branch=master)](http://travis-ci.org/hybridgroup/gort) [![Go Report Card](https://goreportcard.com/badge/github.com/hybridgroup/gort)](https://goreportcard.com/report/github.com/hybridgroup/gort)

## Getting Started
We now have precompiled binaries! You can also build from source.

The Gort CLI provides many useful features on many hardware platforms, and has no other dependencies. You install Gort separately from any framework, which means you can use it to program Arduinos with the Firmata firmware also compatible with Cylon.js, Gobot, Artoo, & Johnny-Five.

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
   0.6.1

COMMANDS:
   scan         Scan for connected devices on Serial, USB, or Bluetooth ports
   bluetooth    Connect & disconnect bluetooth devices.
   arduino      Install avrdude, and upload HEX files to your Arduino
   particle     Upload sketches to your Particle Photon
   digispark    Configure your Digispark microcontroller
   crazyflie    Configure your Crazyflie
   klaatu       barada nikto
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h           show help
   --version, -v        print the version
```

Scan for connected serial devices:

```
$ gort scan serial
[    0.000000] console [tty0] enabled
```

More help coming soon...

## Building

You need to have installed go-bindata to build the file assets into Gort for a single standalone executable:

```
go get github.com/jteeuwen/go-bindata/...
```

Once installed, you build the assets into the project like this:
```
cd commands && go-bindata -pkg="commands" support/... && cd ..
```

## Release

You need to have goxc installed in order to cross compile Gort:

```
go get github.com/laher/goxc
```

Once installed, you can build the binary with
```
make release
```

Compiled binaries will now be placed in `build/<VERSION>/`

You will probably need to set

`export GOBIN=$GOPATH/bin` in order to run the build.


## Contributing
For our contribution guidelines, please go to [https://github.com/hybridgroup/gort/blob/master/CONTRIBUTING.md
](https://github.com/hybridgroup/gort/blob/master/CONTRIBUTING.md
).

## Release History
Version 0.6.1 - Add Debian control file

Version 0.6.0 - simplify and cleanup bluetooth commands

Version 0.5.0 - Adds board-type flag to Arduino upload command to support many more kinds of boards, refactor Bluetooth command params to be in a more logical order

Version 0.4.1 - Corrections for bad merge

Version 0.4.0 - Bluetooth commands use pure hcitool/rfcomm to avoid dependencies on Linux, Spark is now Particle, remove dronedrop commands

Version 0.3.0 - Add dronedrop commands

Version 0.2.4 - Update voodoospark to 2.3.1 and bug fixes

Version 0.2.3 - Update Windows compatibility and default Spark code

Version 0.2.2 - Correct error in avdude install for Linux

Version 0.2.1 - Update default Spark code for servo support

Version 0.2.0 - Add Windows support for Arduino, & bug fixes

Version 0.1.0 - Initial Release

## Licenses
Gort is copyright (c) 2014-2016 The Hybrid Group. Licensed under the Apache 2.0 license.

Firmata is copyright (c) 2006-2008 Hans-Christoph Steiner. Licensed under GNU Lesser General Public License. All rights reserved.

Rapiro is copyright (c) 2013-2014 Shota Ishiwatari. Licensed under the Creative Commons - Public Domain Dedication License.

VoodooSpark is copyright (c) 2012, 2013, 2014 Rick Waldron & Chris Williams. Licensed under the MIT License.
