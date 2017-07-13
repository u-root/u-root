
u-root
======

[![Build Status](https://travis-ci.org/u-root/u-root.svg?branch=master)](https://travis-ci.org/u-root/u-root) [![Go Report Card](https://goreportcard.com/badge/github.com/u-root/u-root)](https://goreportcard.com/report/github.com/u-root/u-root) [![GoDoc](https://godoc.org/github.com/u-root/u-root?status.svg)](https://godoc.org/github.com/u-root/u-root) [![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://github.com/u-root/u-root/blob/master/LICENSE)

# Description

u-root is a "universal root". It's a root file system with mostly Go source with the exception of 5 binaries. 

That's the interesting part. This set of utilities is all Go, and mostly source.

When you run a command that is not built, you fall through to the command that does a
`go build` of the command, and then execs the command once it is built. From that point on,
when you run the command, you get the one in tmpfs. This is fast.

# Setup

You'll need a GOPATH. Be sure to set it to something, e.g.

`export GOPATH=/usr/local/src/go`

On my machine, my gopath is
`export GOPATH=/home/$USER/go`

Then
`go get github.com/u-root/u-root`

`cd $GOPATH/src/github.com/u-root/u-root`

You may hit a problem where it can't find some standard Go packages, if so, you'll need
to set GOROOT, e.g.
`export GOROOT=/path/to/some_go_>=1.6`

# Using

To try the chroot, just run the README:
`bash RUN`

## Build an Embeddable U-root

You can build this environment into a kernel as an initramfs, and further
embed that into firmware as a coreboot payload. 

In the kernel and coreboot case, you need to configure ethernet. We have a primitive
ip command for that case. In qemu:
`ip addr add 10.0.2.15/8 dev eth0`
`ip link set dev eth0 up`

Or, on newer linux kernels (> 4.x) boot with ip=dhcp in the command line. There's also a dhcp command.



## Getting Packages of TinyCore

You can install tinycore linux packages for things you want.
You can use QEMU NAT to allow you to fetch packages.
Let's suppose, for example, you want bash. Once u-root is
running, you can do this:
`% tcz bash`

The tcz command computes and fetches all dependencies.

If you can't get to tinycorelinux.net, or you want package fetching to be faster,
you can run your own server for tinycore packages. 

You can do this to get a local server using the u-root srvfiles command:
`% src/srvfiles/srvfiles -p 80 -d path-to-local-tinycore-packages`

Of course you have to fetch all those packages first somehow :-)

In the EXAMPLES directory you can see examples of running in a chroot, kernel, and coreboot.



## Using elvish: a more handy shell

In default, rush is the shell in u-root. Now, thanks to Qi Xiao(\<xiaqqaix@gmail.com\>), u-root users are also able to use a friendly and expressive unix-like shell: __elvish__. Users are free to choose whether to include elvish in u-root or not. Basically, elvish has handy functionalities such as auto completion, command-line existence checks, etc. More info of elvish can be found at: [http://github/elves/elvish](http://github.com/elves/elvish).

If you prefer to use elvish as shell in u-root, here are the instructions:

1. Get project __elvish__:
  `go get github.com/elves/elvish`

2. Temporarily, since package `sqlite3` used in elvish has been updated, and its latest
   version includes codes in C (which u-root does not support), users have to
   roll back to last good commit of elvish:
   `cd $GOPATH/src/elves/elvish`
   `git checkout bc5543aef2c493b658d6bd1bb81e3de298de8d2f`

3. Go to u-root repo. If you did `go get github.com/u-root/u-root` before, do:
  `cd $GOPATH/src/u-root/u-root`

4. If you prefer to build under bb mode, please do the following command line
   in u-root/u-root/:
   `cd ./bb/`
   `go build .`
   `CGO_ENABLED=0 ./bb 'src/github.com/u-root/u-root/cmds/[a-z]*' src/github.com/elves/elvish`
   which generates a cpio file, /tmp/initramfs.linux\_amd64.cpio for you to
   start up u-root in qemu.

   If you prefer dynamic buildup mode, do the following command line in u-root/u-root:
   `CGO_ENABLED=0 go run scripts/ramfs.go 'src/github.com/u-root/u-root/cmds/[a-z]*' src/github.com/elves/elvish`
   which also generates /tmp/initramfs.linux\_amd64.cpio.

5. Afterwards, users can type command line `elvish` in u-root and start to use elvish as shell.



# Contributions

We need help with this project, so contributions are welcome.  More information about handle dependencies you can found [here](MAINTAINERS.md)

