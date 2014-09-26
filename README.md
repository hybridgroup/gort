# Gort - Command Line Interface For RobotOps

Gort (http://gort.io) is a Command Line Toolkit for RobotOps. Gort provides tools to scan for connected devices, upload firmware, and more.

Gort is written in the Go programming language (http://golang.org) for maximum speed and portability.

Want to use Golang to program your robots? Check out our open source robotics framework Gobot (http://gobot.io).

Want to use Javascript on Robots? Check out Cylon.js (http://cylonjs.com)

Want to use Ruby on robots? Check out Artoo (http://artoo.io)

[![Build Status](https://secure.travis-ci.org/hybridgroup/gort.png?branch=master)](http://travis-ci.org/hybridgroup/gort)

## Getting Started
We now have precompiled binaries! You can also build from source.

The Gort CLI provides many useful features on many hardware platforms, and has no other dependencies. You install Gort separately from any framework, which means you can use it to program Arduinos with the Firmata firmware also compatible with Cylon.js, Gobot, Artoo, & Johnny-Five. 

## Downloads (version 0.2.4)

### Darwin (Apple Mac)

 * [gort\_0.2.4\_darwin\_386.zip](https://s3.amazonaws.com/gort-io/0.2.4/gort_0.2.4_darwin_386.zip)
 * [gort\_0.2.4\_darwin\_amd64.zip](https://s3.amazonaws.com/gort-io/0.2.4/gort_0.2.4_darwin_amd64.zip)

### Linux

 * [gort\_0.2.4\_amd64.deb](https://s3.amazonaws.com/gort-io/0.2.4/gort_0.2.4_amd64.deb)
 * [gort\_0.2.4\_armhf.deb](https://s3.amazonaws.com/gort-io/0.2.4/gort_0.2.4_armhf.deb)
 * [gort\_0.2.4\_i386.deb](https://s3.amazonaws.com/gort-io/0.2.4/gort_0.2.4_i386.deb)
 * [gort\_0.2.4\_linux\_386.tar.gz](https://s3.amazonaws.com/gort-io/0.2.4/gort_0.2.4_linux_386.tar.gz)
 * [gort\_0.2.4\_linux\_amd64.tar.gz](https://s3.amazonaws.com/gort-io/0.2.4/gort_0.2.4_linux_amd64.tar.gz)
 * [gort\_0.2.4\_linux\_arm.tar.gz](https://s3.amazonaws.com/gort-io/0.2.4/gort_0.2.4_linux_arm.tar.gz)

### MS Windows

 * [gort\_0.2.4\_windows\_386.zip](https://s3.amazonaws.com/gort-io/0.2.4/gort_0.2.4_windows_386.zip)
 * [gort\_0.2.4\_windows\_amd64.zip](https://s3.amazonaws.com/gort-io/0.2.4/gort_0.2.4_windows_amd64.zip)


## How To Use

```
$ ./gort
NAME:
   gort - Command Line Utility for RobotOps

USAGE:
   gort [global options] command [command options] [arguments...]

VERSION:
   0.2.4

COMMANDS:
   scan   Scan for connected devices on Serial, USB, or Bluetooth ports
   bluetooth  Pair, unpair & connect to bluetooth devices.
   arduino  Install avrdude, and upload sketches to your Arduino
   spark  Upload sketches to your Spark
   digispark  Configure your Digispark microcontroller
   crazyflie  Configure your Crazyflie
   klaatu barada nikto
   help, h  Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --version, -v  print the version
   --help, -h   show help
```

Scan for connected serial devices:

```
$ ./gort scan serial
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

Compilied binaries will now be placed in `build/<VERSION>/`


## Contributing

* All patches must be provided under the Apache 2.0 License
* Please use the -s option in git to "sign off" that the commit is your work and you are providing it under the Apache 2.0 License
* Submit a Github Pull Request to the appropriate branch and ideally discuss the changes with us in IRC.
* We will look at the patch, test it out, and give you feedback.
* Avoid doing minor whitespace changes, renamings, etc. along with merged content. These will be done by the maintainers from time to time but they can complicate merges and should be done seperately.
* Add unit tests for any new or changed functionality
* Take care to maintain the existing coding style.
* All pull requests should be "fast forward"
  * If there are commits after yours use “git rebase -i <new_head_branch>”
  * If you have local changes you may need to use “git stash”
  * For git help see [progit](http://git-scm.com/book) which is an awesome (and free) book on git

## Release History

Version 0.2.4 - Update voodoospark to 2.3.1 and bug fixes

Version 0.2.3 - Update Windows compatibility and default Spark code 

Version 0.2.2 - Correct error in avdude install for Linux

Version 0.2.1 - Update default Spark code for servo support

Version 0.2.0 - Add Windows support for Arduino, & bug fixes

Version 0.1.0 - Initial Release

## Licenses
Gort is copyright (c) 2014 The Hybrid Group. Licensed under the Apache 2.0 license.

Firmata is copyright (c) 2006-2008 Hans-Christoph Steiner. Licensed under GNU Lesser General Public License. All rights reserved.

Rapiro is copyright (c) 2013-2014 Shota Ishiwatari. Licensed under the Creative Commons - Public Domain Dedication License.

VoodooSpark is copyright (c) 2012, 2013, 2014 Rick Waldron & Chris Williams. Licensed under the MIT License.
