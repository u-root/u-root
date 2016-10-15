
u-root
======

[![Build Status](https://travis-ci.org/u-root/u-root.svg?branch=master)](https://travis-ci.org/u-root/u-root)

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



# Contributions

We need help with this project, so contributions are welcome.  More information about handle dependencies you can found [here](MAINTENERS.MD)

