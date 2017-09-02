
u-root
======

[![Build Status](https://travis-ci.org/u-root/u-root.svg?branch=master)](https://travis-ci.org/u-root/u-root) [![Go Report Card](https://goreportcard.com/badge/github.com/u-root/u-root)](https://goreportcard.com/report/github.com/u-root/u-root) [![GoDoc](https://godoc.org/github.com/u-root/u-root?status.svg)](https://godoc.org/github.com/u-root/u-root) [![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://github.com/u-root/u-root/blob/master/LICENSE)

# Description

u-root is a "universal root". It's a root file system with mostly Go source with the exception of 5 binaries.

That's the interesting part. This set of utilities is all Go, and mostly source.

When you run a command that is not built, you fall through to the command that does a
`go build` of the command, and then execs the command once it is built. From that point on,
when you run the command, you get the one in tmpfs. This is fast.


# Usage

Make sure your Go version is the latest (>=1.8). Correctly set up your GOPATH like so:

    $ export GOPATH="$HOME/go"
    $ export PATH="$PATH:$GOPATH/bin"

Now, download and install u-root:

    $ go get github.com/u-root/u-root

You can now use the u-root command anywhere for building. Here are some examples:

    $ u-root --run                            # build and run in a chroot (requires sudo)
    $ u-root --format=cpio -o initramfs.cpio  # generate a cpio archive named initramfs.cpio
    $ u-root --format=cpio --run              # create a cpio in /tmp and run with qemu
    $ u-root --format=cpio --build_format=bb  # create a cpio containing a busybox
    $ u-root --format=docker --run            # build and run a docker image

It is also possible to specify packages for inclusion:

    $ go get github.com/golang/example/hello
    $ u-root --run github.com/golang/example/hello


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



# Contributions

We need help with this project, so contributions are welcome.  More information about handle dependencies you can found [here](MAINTAINERS.md)

