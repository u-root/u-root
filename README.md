# u-root

[![Build Status](https://circleci.com/gh/u-root/u-root/tree/main.png?style=shield&circle-token=8d9396e32f76f82bf4257b60b414743e57734244)](https://circleci.com/gh/u-root/u-root/tree/main)
[![codecov](https://codecov.io/gh/u-root/u-root/branch/main/graph/badge.svg?token=1qjHT02oCB)](https://codecov.io/gh/u-root/u-root)
[![Go Report Card](https://goreportcard.com/badge/github.com/u-root/u-root)](https://goreportcard.com/report/github.com/u-root/u-root)
[![CodeQL](https://github.com/u-root/u-root/workflows/CodeQL/badge.svg)](https://github.com/u-root/u-root/actions?query=workflow%3ACodeQL)
[![GoDoc](https://godoc.org/github.com/u-root/u-root?status.svg)](https://godoc.org/github.com/u-root/u-root)
[![Slack](https://slack.osfw.dev/badge.svg)](https://slack.osfw.dev)
[![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://github.com/u-root/u-root/blob/main/LICENSE)


# Description

u-root embodies four different projects.

*   Go versions of many standard Linux tools, such as [ls](cmds/core/ls/ls.go),
    [cp](cmds/core/cp/cp.go), or [shutdown](cmds/core/shutdown/shutdown_linux.go). See
    [cmds/core](cmds/core) for most of these.

*   A way to compile many Go programs into a single binary with
    [busybox mode](#build-modes).

*   A way to create initramfs (an archive of files) to use with Linux kernels,
    [embeddable into firmware](#build-an-embeddable-u-root).

*   Go [bootloaders](#systemboot) that use `kexec` to boot Linux or multiboot
    kernels such as ESXi, Xen, or tboot. They are meant to be used with
    [LinuxBoot](https://www.linuxboot.org).

# Usage

Make sure your Go version is >=1.19.

Download and install u-root either via git:

```shell
git clone https://github.com/u-root/u-root
cd u-root
go build
```

The resulting binary will the be placed where `go build` was invoked

Or install directly with go:

```shell
go install github.com/u-root/u-root
```

**Note: The `u-root` command will end up in `$GOPATH/bin/u-root`, so you may
need to add `$GOPATH/bin` to your `$PATH`.**

## Templates

To quickly specify a set of commands from u-root, you can use any of the
templates as defined in [templates.go](templates.go).

## Examples

Here are some examples of using the `u-root` command to build an initramfs,
with `$UROOT_PATH` being the path to where the u-root sources are on the disk
(explicitly specifiying this is only necessary if not running u-root inside the
root of the repository):

```shell
# Build an initramfs of all the Go cmds in ./cmds/core/... (default)
u-root

# Build an initramfs of all the Go cmds in ./cmds/core/...
# But running the command outside of the repository root
(cd /tmp && GBB_PATH=$UROOT_PATH u-root)

# Generate an archive with bootloaders
#
# core and boot are templates that expand to sets of commands
u-root core boot

# Generate an archive with only these given commands
u-root ./cmds/core/{init,ls,ip,dhclient,wget,cat,elvish}

# Generate an archive with all of the core tools with some exceptions
u-root core -cmds/core/{ls,losetup}

# Generate an archive with a tool outside of u-root
git clone https://github.com/u-root/cpu
u-root ./cmds/core/{init,ls,elvish} ./cpu/cmds/cpud

# Generate an archive with a tool outside of u-root, in any PWD
(cd /tmp && GBB_PATH=$UROOT_PATH:$CPU_PATH u-root ./cmds/core/{init,ls,elvish} ./cmds/cpud)
```

The default set of packages included is all packages in
`github.com/u-root/u-root/cmds/core/...`.

`GBB_PATH` is a place that u-root will look for commands. Each colon-separated
`GBB_PATH` element is concatenated with patterns from the command-line and
checked for existence. For example:

```shell
GBB_PATH=$HOME/u-root:$HOME/u-bmc u-root \
    cmds/core/init \
    cmds/core/elvish \
    cmd/socreset

# matches:
#   $HOME/u-root/cmds/core/init
#   $HOME/u-root/cmds/core/elvish
#   $HOME/u-bmc/cmd/socreset
```

## Extra Files

You may also include additional files in the initramfs using the `-files` flag.
If you add binaries with `-files` are listed, their ldd dependencies will be
included as well. As example for Debian, you want to add two kernel modules for
testing, executing your currently booted kernel:

> NOTE: these files will be placed in the `$HOME` dir in the initramfs.

```shell
u-root -files $HOME/hello.ko -files $HOME/hello2.ko
qemu-system-x86_64 -kernel /boot/vmlinuz-$(uname -r) -initrd /tmp/initramfs.linux_amd64.cpio
```

To specify the location in the initramfs, use `<sourcefile>:<destinationfile>`.
For example:

```shell
u-root -files "root-fs/usr/bin/runc:usr/bin/run"
```

## Init and Uinit

u-root has a very simple (exchangable) init system controlled by the `-initcmd`
and `-uinitcmd` command-line flags.

*   `-initcmd` determines what `/init` is symlinked to. `-initcmd` may be a
    u-root command name or a symlink target.
*   `-uinitcmd` is run by the default u-root [init](cmds/core/init) after some
    basic file system setup. There is no default, users should optionally supply
    their own. `-uinitcmd` may be a u-root command name with arguments or a
    symlink target with arguments.
*   After running a uinit (if there is one), [init](cmds/core/init) will start a
    shell determined by the `-defaultsh` argument.

We expect most users to keep their `-initcmd` as [init](cmds/core/init), but to
supply their own uinit for additional initialization or to immediately load
another operating system.

All three command-line args accept both a u-root command name or a target
symlink path. **Only `-uinitcmd` accepts command-line arguments, however.** For
example,

```bash
u-root -uinitcmd="echo Go Gopher" ./cmds/core/{init,echo,elvish}

cpio -ivt < /tmp/initramfs.linux_amd64.cpio
# ...
# lrwxrwxrwx   0 root     root           12 Dec 31  1969 bin/uinit -> ../bbin/echo
# lrwxrwxrwx   0 root     root            9 Dec 31  1969 init -> bbin/init

qemu-system-x86_64 -kernel $KERNEL -initrd /tmp/initramfs.linux_amd64.cpio -nographic -append "console=ttyS0"
# ...
# [    0.848021] Freeing unused kernel memory: 896K
# 2020/05/01 04:04:39 Welcome to u-root!
#                              _
#   _   _      _ __ ___   ___ | |_
#  | | | |____| '__/ _ \ / _ \| __|
#  | |_| |____| | | (_) | (_) | |_
#   \__,_|    |_|  \___/ \___/ \__|
#
# Go Gopher
# ~/>
```
Passing command line arguments like above is equivalent to passing the arguments to uinit via a flags file in `/etc/uinit.flags`, see [Extra Files](#extra-files).

Additionally, you can pass arguments to uinit via the `uroot.uinitargs` kernel parameters, for example:

```bash
u-root -uinitcmd="echo Gopher" ./cmds/core/{init,echo,elvish}

cpio -ivt < /tmp/initramfs.linux_amd64.cpio
# ...
# lrwxrwxrwx   0 root     root           12 Dec 31  1969 bin/uinit -> ../bbin/echo
# lrwxrwxrwx   0 root     root            9 Dec 31  1969 init -> bbin/init

qemu-system-x86_64 -kernel $KERNEL -initrd /tmp/initramfs.linux_amd64.cpio -nographic -append "console=ttyS0 uroot.uinitargs=Go"
# ...
# [    0.848021] Freeing unused kernel memory: 896K
# 2020/05/01 04:04:39 Welcome to u-root!
#                              _
#   _   _      _ __ ___   ___ | |_
#  | | | |____| '__/ _ \ / _ \| __|
#  | |_| |____| | | (_) | (_) | |_
#   \__,_|    |_|  \___/ \___/ \__|
#
# Go Gopher
# ~/>
```
Note the order of the passed arguments in the above example.


The command you name must be present in the command set. The following will *not
work*:

```bash
u-root -uinitcmd="echo Go Gopher" ./cmds/core/{init,elvish}
# 2020/04/30 21:05:57 could not create symlink from "bin/uinit" to "echo": command or path "echo" not included in u-root build: specify -uinitcmd="" to ignore this error and build without a uinit
```

You can also refer to non-u-root-commands; they will be added as symlinks. We
don't presume to know whether your symlink target is correct or not.

This will build, but not work unless you add a /bin/foobar to the initramfs.

```bash
u-root -uinitcmd="/bin/foobar Go Gopher" ./cmds/core/{init,elvish}
```

This will boot the same as the above.

```bash
u-root -uinitcmd="/bin/foobar Go Gopher" -files /bin/echo:bin/foobar -files your-hosts-file:/etc/hosts ./cmds/core/{init,elvish}
```

The effect of the above command:
*   Sets up the uinit command to be /bin/foobar, with 2 arguments: Go Gopher
*   Adds /bin/echo as bin/foobar
*   Adds your-hosts-file as etc/hosts
*   builds in the cmds/core/init, and cmds/core/elvish commands.
    The {} are expanded by the shell 

This will bypass the regular u-root init and just launch a shell:

```bash
u-root -initcmd=elvish ./cmds/core/{elvish,ls}

cpio -ivt < /tmp/initramfs.linux_amd64.cpio
# ...
# lrwxrwxrwx   0 root     root            9 Dec 31  1969 init -> bbin/elvish

qemu-system-x86_64 -kernel $KERNEL -initrd /tmp/initramfs.linux_amd64.cpio -nographic -append "console=ttyS0"
# ...
# [    0.848021] Freeing unused kernel memory: 896K
# failed to put myself in foreground: ioctl: inappropriate ioctl for device
# ~/>
```

(It fails to do that because some initialization is missing when the shell is
started without a proper init.)

## Cross Compilation (targeting different architectures and OSes)

Cross-OS and -architecture compilation comes for free with Go. In fact, every PR
to the u-root repo is built against the following architectures: amd64, x86
(i.e. 32bit), mipsle, armv7, arm64, and ppc64le.

Further, we run integration tests on linux/amd64, freebsd/amd64 and linux/arm64,
using several CI systems. If you need to add another CI system, processor or OS,
please let us know.

To cross compile for an ARM, on Linux:

```shell
GOARCH=arm u-root
```

If you are on OSX, and wish to build for Linux on AMD64:

```shell
GOOS=linux GOARCH=amd64 u-root
```

## Testing in QEMU

A good way to test the initramfs generated by u-root is with qemu:

```shell
qemu-system-x86_64 -nographic -kernel path/to/kernel -initrd /tmp/initramfs.linux_amd64.cpio
```

Note that you do not have to build a special kernel on your own, it is
sufficient to use an existing one. Usually you can find one in `/boot`.

If you quickly need to obtain a kernel, for example, when you are on a non-Linux
system, you can assemble a URL to download one through Arch Linux's
[iPXE menu file](https://www.archlinux.org/releng/netboot/archlinux.ipxe). It
would download from `${mirrorurl}iso/${release}/arch/boot/x86_64/vmlinuz-linux`, so
just search for a mirror URL you prefer and a release version, for example,
`http://mirror.rackspace.com/archlinux/iso/2022.05.01/arch/boot/x86_64/vmlinuz-linux`.

### Framebuffer

For framebuffer support, append a VESA mode via the `vga` kernel parameter:

```shell
qemu-system-x86_64 \
  -kernel path/to/kernel \
  -initrd /tmp/initramfs.linux_amd64.cpio \
  -append "vga=786"
```

For a list of modes, refer to the
[Linux kernel documentation](https://github.com/torvalds/linux/blob/master/Documentation/fb/vesafb.rst#how-to-use-it).

### Entropy / Random Number Generator

Some utilities, e.g., `dhclient`, require entropy to be present. For a speedy
virtualized random number generator, the kernel should have the following:

```shell
CONFIG_VIRTIO_PCI=y
CONFIG_HW_RANDOM_VIRTIO=y
CONFIG_CRYPTO_DEV_VIRTIO=y
```

Then you can run your kernel in QEMU with a `virtio-rng-pci` device:

```shell
qemu-system-x86_64 \
    -device virtio-rng-pci \
    -kernel vmlinuz \
    -initrd /tmp/initramfs.linux_amd64.cpio
```

In addition, you can pass your host's RNG:

```shell
qemu-system-x86_64 \
    -object rng-random,filename=/dev/urandom,id=rng0 \
    -device virtio-rng-pci,rng=rng0 \
    -kernel vmlinuz \
    -initrd /tmp/initramfs.linux_amd64.cpio
```

## u-root with Go package paths

For Go package paths to be usable, the path passed to `u-root` must be in the
go.mod of the working directory or one of its parents. This is mostly useful for
repositories making programmatic use of u-root's APIs.

```sh
cd ./u-root

# In u-root's directory itself, github.com/u-root/u-root is resolvable. There is
# a go.mod here that can refer to u-root.
u-root github.com/u-root/u-root/cmds/core/...
u-root github.com/u-root/u-root/cmds/core/*
u-root github.com/u-root/u-root/cmds/core/i*
```

To depend on commands outside of ones own repository, the easiest way to depend
on Go commands is the following:

```sh
TMPDIR=$(mktemp -d)
cd $TMPDIR
go mod init foobar
```

Create a file with some unused build tag like this to create dependencies on
commands:

```go
//go:build tools

package something

import (
        "github.com/u-root/u-root/cmds/core/ip"
        "github.com/u-root/u-root/cmds/core/init"
        "github.com/hugelgupf/p9/cmd/p9ufs"
)
```

The unused build tag keeps it from being compiled, but its existence forces `go
mod tidy` to add these dependencies to `go.mod`:

```sh
go mod tidy

u-root \
  github.com/u-root/u-root/cmds/core/ip \
  github.com/u-root/u-root/cmds/core/init \
  github.com/hugelgupf/p9/cmd/p9ufs
```

## SystemBoot

SystemBoot is a set of bootloaders written in Go. It is meant to be a
distribution for LinuxBoot to create a system firmware + bootloader. All of
these use `kexec` to boot. The commands are in [cmds/boot](cmds/boot).
Parsers are available for [GRUB](pkg/boot/grub), [syslinux](pkg/boot/syslinux),
and other config files to make the transition to LinuxBoot easier.

*   `pxeboot`: a network boot client that uses DHCP and HTTP or TFTP to get a
    boot configuration which can be parsed as PXELinux or iPXE configuration
    files to get a boot program.

*   `boot`: finds all bootable kernels on local disk, shows a menu, and boots
    them. Supports (basic) GRUB, (basic) syslinux, (non-EFI) BootLoaderSpec, and
    ESXi configurations.

*   `fbnetboot`: a network boot client that uses DHCP and HTTP to get a boot
    program based on Linux, and boots it. To be merged with `pxeboot`.

*   `localboot`: a tool that finds bootable kernel configurations on the local
    disks and boots them.

*   `systemboot`: a wrapper around `pxeboot` and `boot` that just mimicks
    a BIOS/UEFI BDS behaviour, by looping between network booting and local
    booting. Use `-uinitcmd` argument to the u-root build tool to make it the
    boot program.

This project started as a loose collection of programs in u-root by various
LinuxBoot contributors, as well as a personal experiment by
[Andrea Barberio](https://github.com/insomniacslk) that has since been merged
in. It is now an effort of a broader community and graduated to a real project
for system firmwares.

More detailed information about the build process for a full LinuxBoot firmware
image using u-root/systemboot and coreboot can be found in the
[LinuxBoot book](https://github.com/linuxboot/book) chapter about
[LinuxBoot using coreboot, u-root and systemboot](https://github.com/linuxboot/book/blob/master/coreboot.u-root.systemboot/README.md).

You can build systemboot like this:

```shell
u-root -uinitcmd=systemboot core ./cmds/boot/{systemboot,localboot,fbnetboot}
```

## Compression

You can compress the initramfs. However, for xz compression, the kernel has some
restrictions on the compression options and it is suggested to align the file to
512 byte boundaries:

```shell
xz --check=crc32 -9 --lzma2=dict=1MiB \
   --stdout /tmp/initramfs.linux_amd64.cpio \
   | dd conv=sync bs=512 \
   of=/tmp/initramfs.linux_amd64.cpio.xz
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

## Build an Embeddable u-root

You can build the cpio image created by u-root into a Linux kernel via the
`CONFIG_INITRAMFS_SOURCE` config variable or coreboot config variable, and
further embed the kernel image into firmware as a coreboot payload.

In the kernel and coreboot case, you may need to configure ethernet. We have a
`dhclient` command that works for both ipv4 and ipv6. Since v6 does not yet work
that well for most people, a typical invocation looks like this:

```shell
% dhclient -ipv4 -ipv6=false
```

Or, on newer linux kernels (> 4.x) boot with ip=dhcp in the command line,
assuming your kernel is configured to work that way.

## Build Modes

u-root can create an initramfs in two different modes, specified by `-build`:

*   `gbb` mode: One busybox-like binary comprising all the Go tools you ask to
    include.
    See [the gobusybox README for how it works](https://github.com/u-root/gobusybox).
    In this mode, u-root copies and rewrites the source of the tools you asked
    to include to be able to compile everything into one busybox-like binary.

*   `binary` mode: each specified binary is compiled separately and all binaries
    are added to the initramfs.

## Updating Dependencies

```shell
go get -u
go mod tidy
go mod vendor
```

## Building without network access

Go modules require network access. If you need to make a repeatable build with
no network access, make sure that your code is under `$GOPATH` and the
environment variable `GO111MODULE` is set to `off`. This is:

1. Pick a location for your off-network build, it can be anywhere and
the directory does not need to exist ahead of time:

```shell
export GOPATH=$(mktemp -d)

```

2. Fetch the code, you can use `git`, `go get` or even a release file, just
make sure that the code ends in: `${GOPATH}/src/github.com/u-root/u-root` E.g:

```shell
mkdir -p ${GOPATH}/src/github.com/u-root/
cd ${GOPATH}/src/github.com/u-root/
git clone https://github.com/u-root/u-root.git
cd u-root
```

Or simply:

```shell
GO111MODULE=off go get github.com/u-root/u-root
cd $GOPATH/src/github.com/u-root/u-root
```

3. Build u-root and use it normally:

```shell
GO111MODULE=off GOPROXY=off go build
GO111MODULE=off GOPROXY=off ./u-root
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
