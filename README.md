# u-root

[![Build Status](https://circleci.com/gh/u-root/u-root/tree/master.png?style=shield&circle-token=8d9396e32f76f82bf4257b60b414743e57734244)](https://circleci.com/gh/u-root/u-root/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/u-root/u-root)](https://goreportcard.com/report/github.com/u-root/u-root)
[![GoDoc](https://godoc.org/github.com/u-root/u-root?status.svg)](https://godoc.org/github.com/u-root/u-root)
[![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://github.com/u-root/u-root/blob/master/LICENSE)

# Description

u-root embodies four different projects.

*   Go versions of many standard Linux tools, such as [ls](cmds/core/ls/ls.go),
    [cp](cmds/core/cp/cp.go), or [shutdown](cmds/core/shutdown/shutdown.go). See
    [cmds/core](cmds/core) for most of these.

*   Go bootloaders that use `kexec` to boot Linux or multiboot kernels such as
    ESXi, Xen, or tboot. They are meant to be used with
    [LinuxBoot](https://www.linuxboot.org). With that, parsers for
    [GRUB config files](pkg/diskboot) or [syslinux config files](pkg/syslinux)
    are to make transition to LinuxBoot easier.

*   A way to create very small Go programs using
    [busybox mode](pkg/bb/README.md) or source mode (see below).

*   A way to create initramfs (an archive of files) to use with Linux kernels.

# Creating Initramfs Archives

u-root can create an initramfs in two different modes:

*   source mode includes Go toolchain binaries + simple shell + Go source files
    in the initramfs archive. Tools are compiled from source on the fly by the
    shell.

    When you try to run a command that is not built, it is compiled first and
    stored in tmpfs. From that point on, when you run the command, you get the
    one in tmpfs. Don't worry: the Go compiler is pretty fast.

*   bb mode: One busybox-like binary comprising all the Go tools you ask to
    include. See [here for how it works](pkg/bb/README.md).

    In this mode, u-root copies and rewrites the source of the tools you asked
    to include to be able to compile everything into one busybox-like binary.

# SystemBoot

SystemBoot is a set of bootloaders written in Go. It is meant to be a
distribution for LinuxBoot to create a system firmware + bootloader. All of
these use `kexec` to boot. The commands are in [cmds/boot](cmds/boot).

*   `pxeboot`: a network boot client that uses DHCP and HTTP or TFTP to get a
    boot configuration which can be parsed as PXELinux or iPXE configuration
    files to get a boot program.

*   `fbnetboot`: a network boot client that uses DHCP and HTTP to get a boot
    program based on Linux, and boots it. To be merged with `pxeboot`.

*   `localboot`: a tool that finds bootable kernel configurations on the local
    disks and boots them.

*   `boot2`: similar to `localboot`, finds a bootable kernel configuration on
    disk (GRUB or syslinux) and boots it. To be merged into `localboot`.

*   `uinit`: a wrapper around `netboot` and `localboot` that just mimicks a
    BIOS/UEFI BDS behaviour, by looping between network booting and local
    booting. The name `uinit` is necessary to be picked up as boot program by
    u-root.

This project started as a loose collection of programs in u-root by various
LinuxBoot contributors, as well as a personal experiment by
[Andrea Barberio](https://github.com/insomniacslk) that has since been merged
in. It is now an effort of a broader community and graduated to a real project
for system firmwares.

More detailed information about the build process for a full LinuxBoot firmware
image using u-root/systemboot and coreboot can be found in the
[LinuxBoot book](https://github.com/linuxboot/book) chapter 11,
[LinuxBoot using coreboot, u-root and systemboot](https://github.com/linuxboot/book/blob/master/11.coreboot.u-root.systemboot/README.md).

# Usage

Make sure your Go version is 1.12. Make sure your `GOPATH` is set up correctly.

Download and install u-root:

```shell
go get github.com/u-root/u-root
```

You can now use the u-root command to build an initramfs. Here are some
examples:

```shell
# Build a bb-mode cpio initramfs of all the Go cmds in ./cmds/core/...
u-root -build=bb

# Generate a bb-mode archive with bootloaders
u-root -build=bb core boot

# Generate a cpio archive named initramfs.cpio.
u-root -format=cpio -build=source -o initramfs.cpio

# Generate a bb-mode archive with only these given commands.
u-root -format=cpio -build=bb ./cmds/core/{init,ls,ip,dhclient,wget,cat}
```

`-format=cpio` and `-build=source` are the default flag values. The default set
of packages included is all packages in
`github.com/u-root/u-root/cmds/core/...`.

In addition to using paths to specify Go source packages to include, you may
also use Go package import paths (e.g. `golang.org/x/tools/imports`) to include
commands. Only the `main` package and its dependencies in those source
directories will be included. For example:

You can build the initramfs built by u-root into the kernel via the
`CONFIG_INITRAMFS_SOURCE` config variable or you can load it separately via an
option in for example Grub or the QEMU command line or coreboot config variable.

A good way to test the initramfs generated by u-root is with qemu:

```shell
qemu-system-x86_64 -nographic -kernel path/to/kernel -initrd /tmp/initramfs.linux_amd64.cpio
```

Note that you do not have to build a special kernel on your own, it is
sufficient to use an existing one. Usually you can find one in `/boot`.

> NOTE: you can compress the initramfs but for xz compression, the kernel has
> some restrictions on the compression options and it is suggested to align the
> file to 512 byte boundaries `shell xz --check=crc32 -9 --lzma2=dict=1MiB
> --stdout /tmp/initramfs.linux_amd64.cpio \ | dd conv=sync bs=512
> of=/tmp/initramfs.linux_amd64.cpio.xz`

You may also include additional files in the initramfs using the `-files` flag.
If you add binaries with `-files` are listed, their ldd dependencies will be
included as well. As example for Debian, you want to add two kernel modules for
testing, executing your currently booted kernel:

> NOTE: these files will be placed in the `$HOME` dir in the initramfs.

```shell
u-root -files "$HOME/hello.ko $HOME/hello2.ko"
qemu-system-x86_64 -kernel /boot/vmlinuz-$(uname -r) -initrd /tmp/initramfs.linux_amd64.cpio
```

To specify the location in the initramfs, use `<sourcefile>:<destinationfile>`.
For example:

```shell
u-root -files "root-fs/usr/bin/runc:usr/bin/run"
```

## Getting Packages of TinyCore

Using the `tcz` command included in u-root, you can install tinycore linux
packages for things you want.

You can use QEMU NAT to allow you to fetch packages. Let's suppose, for example,
you want bash. Once u-root is running, you can do this:

```shell
% tcz bash
```

The tcz command computes and fetches all dependencies. If you can't get to
tinycorelinux.net, or you want package fetching to be faster, you can run your
own server for tinycore packages.

You can do this to get a local server using the u-root srvfiles command:

```shell
% srvfiles -p 80 -d path-to-local-tinycore-packages
```

Of course you have to fetch all those packages first somehow :-)

## Build an Embeddable U-root

You can build this environment into a kernel as an initramfs, and further embed
that into firmware as a coreboot payload.

In the kernel and coreboot case, you need to configure ethernet. We have a
`dhclient` command that works for both ipv4 and ipv6. Since v6 does not yet work
that well for most people, a typical invocation looks like this:

```shell
% dhclient -ipv4 -ipv6=false
```

Or, on newer linux kernels (> 4.x) boot with ip=dhcp in the command line,
assuming your kernel is configured to work that way.

## Updating Dependencies

```shell
# The latest released version of dep is required:
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
dep ensure
```

# Hardware

If you want to see u-root on real hardware, this
[board](https://www.pcengines.ch/apu2.htm) is a good start.

# Contributions

For information about contributing, including how we sign off commits, please
see [CONTRIBUTING.md](CONTRIBUTING.md).

Improving existing commands (e.g., additional currently unsupported flags) is
very welcome. In this case it is not even required to build an initramfs, just
enter the `cmds/` directory and start coding. A list of commands that are on the
roadmap can be found [here](roadmap.md).
