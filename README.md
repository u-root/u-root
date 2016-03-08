u-root
======

A universal root. You mount it, and it's mostly Go source with the exception of 5 binaries. 

And that's the interesting part. This set of utilities is all Go, and mostly source.

The /bin should be mounted in a tmpfs. The directory with the source should be in your path.
The bin in ram comes in your path before the directory with the source code.

When you run a command that is not built, you fall through to the command that does a
'go build' of the command, and then execs the command once it is built. From that point on,
when you run the command, you get the one in tmpfs. This is fast.

To try the chroot, just run 
./README.

In the kernel and coreboot case, you need to configure ethernet. We have a primitive
ip command for that case. Since it's qemu:
ip addr add 10.0.2.15/8 dev eth0
ip link set dev eth0 up

Note that in the kernel and coreboot case, you need to build and run
src/srvfiles/srvfiles -p 80 -d path-to-local-tinycore-packages
UNLESS you set up the qemu to NAT for you.

You can get tinycore packages by running
tcz bash

In the EXAMPLES directory you can see examples of running in a chroot, kernel, and coreboot.

We need help with this project, so contributions are welcome.

[![Build Status](https://travis-ci.org/u-root/uroot.svg?branch=master)](https://travis-ci.org/u-root/u-root)
