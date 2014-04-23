# Gort - Command Line Interface For RobotOps

Gort (http://gort.io) is a Command Line Interface (CLI) for RobotOps. Gort provides tools to scan for connected devices, upload firmware, and more.

Gort is written in the Go programming language (http://golang.org) for maximum speed and portability.

Want to use Golang to program your robots? Check out our open source robotics framework Gobot (http://gobot.io).

Want to use Javascript on Robots? Check out Cylon.js (http://cylonjs.com)

Want to use Ruby on robots? Check out Artoo (http://artoo.io)

[![Build Status](https://secure.travis-ci.org/hybridgroup/gort.png?branch=master)](http://travis-ci.org/hybridgroup/gort)

## Getting Started
We are going to have precompiled binaries available soon. You can also build from source.

The Gort CLI provides many useful features on many hardware platforms, and has no other dependencies. You install Gort separately from any framework, which means you can use it to program Arduinos with the Firmata firmware also compatible with Cylon.js, Gobot, & Artoo, as well as JavaScript frameworks such as Johnny-Five. 

## How To Use

```
$ ./gort 
NAME:
   gort - Command Line Utility for RobotOps

USAGE:
   gort [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   scan		Scan for connected devices on Serial, USB, or Bluetooth ports
   arduino	Install avrdude, and upload sketches to your Arduino
   bluetooth	Scan, pair, unpair bluetooth devices. Establishes serial to Bluetooth connection.
   help, h	Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --version, -v	print the version
   --help, -h		show help
```

Scan for connected serial devices:

```
$ ./gort scan serial
[    0.000000] console [tty0] enabled
```

More help coming soon...

## Contributing

* All patches must be provided under the Apache 2.0 License
* Please use the -s option in git to "sign off" that the commit is your work and you are providing it under the Apache 2.0 License
* Submit a Github Pull Request to the appropriate branch and ideally discuss the changes with us in IRC.
* We will look at the patch, test it out, and give you feedback.
* Avoid doing minor whitespace changes, renamings, etc. along with merged content. These will be done by the maintainers from time to time but they can complicate merges and should be done seperately.
* Take care to maintain the existing coding style.
* All pull requests should be "fast forward"
  * If there are commits after yours use “git rebase -i <new_head_branch>”
  * If you have local changes you may need to use “git stash”
  * For git help see [progit](http://git-scm.com/book) which is an awesome (and free) book on git

## Release History

Version 0.0.1 - Initial Release

## License
Copyright (c) 2014 The Hybrid Group. Licensed under the Apache 2.0 license.
