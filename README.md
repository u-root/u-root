# u-root

[![Build Status](https://circleci.com/gh/u-root/u-root/tree/main.png?style=shield)](https://circleci.com/gh/u-root/u-root/tree/main)
[![codecov](https://codecov.io/gh/u-root/u-root/branch/main/graph/badge.svg?token=1qjHT02oCB)](https://codecov.io/gh/u-root/u-root)
[![Go Report Card](https://goreportcard.com/badge/github.com/u-root/u-root)](https://goreportcard.com/report/github.com/u-root/u-root)
[![CodeQL](https://github.com/u-root/u-root/workflows/CodeQL/badge.svg)](https://github.com/u-root/u-root/actions?query=workflow%3ACodeQL)
[![GoDoc](https://godoc.org/github.com/u-root/u-root?status.svg)](https://godoc.org/github.com/u-root/u-root)
[![Slack](https://slack.osfw.dev/badge.svg)](https://slack.osfw.dev)
[![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://github.com/u-root/u-root/blob/main/LICENSE)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/8160/badge)](https://www.bestpractices.dev/projects/8160)


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

Make sure your Go version is >= 1.21.

Download and install u-root either via git:

```shell
git clone https://github.com/u-root/u-root
cd u-root
go install
```

Or install directly with go:

```shell
go install github.com/u-root/u-root@latest
```

> [!NOTE]
> The `u-root` command will end up in `$GOPATH/bin/u-root`, so you may
> need to add `$GOPATH/bin` to your `$PATH`.

# Examples

Here are some examples of using the `u-root` command to build an initramfs.

```shell
git clone https://github.com/u-root/u-root
cd u-root

# Build an initramfs of all the Go cmds in ./cmds/core/... (default)
u-root

# Generate an archive with bootloaders
#
# core and boot are templates that expand to sets of commands
u-root core boot

# Generate an archive with only these given commands
u-root ./cmds/core/{init,ls,ip,dhclient,wget,cat,gosh}

# Generate an archive with all of the core tools with some exceptions
u-root core -cmds/core/{ls,losetup}
```

> [!IMPORTANT]
>
> `u-root` works exactly when `go build` and `go list` work as well.

> [!NOTE]
>
> The `u-root` tool is the same as the
> [mkuimage](https://github.com/u-root/mkuimage) tool with some defaults
> applied.
>
> In the near future, `uimage` will replace `u-root`.

> [!TIP]
>
> To just build Go busybox binaries, try out
> [gobusybox](https://github.com/u-root/gobusybox)'s `makebb` tool.

## Multi-module workspace builds

There are several ways to build multi-module command images using standard Go
tooling.

```shell
$ mkdir workspace
$ cd workspace
$ git clone https://github.com/u-root/u-root
$ git clone https://github.com/u-root/cpu

$ go work init ./u-root
$ go work use ./cpu

$ u-root ./u-root/cmds/core/{init,gosh} ./cpu/cmds/cpud

$ cpio -ivt < /tmp/initramfs.linux_amd64.cpio
...
-rwxr-x---   0 root     root      6365184 Jan  1  1970 bbin/bb
lrwxrwxrwx   0 root     root            2 Jan  1  1970 bbin/cpud -> bb
lrwxrwxrwx   0 root     root            2 Jan  1  1970 bbin/gosh -> bb
lrwxrwxrwx   0 root     root            2 Jan  1  1970 bbin/init -> bb
...

# Works for offline vendored builds as well.
$ go work vendor # Go 1.22 feature.

$ u-root ./u-root/cmds/core/{init,gosh} ./cpu/cmds/cpud
```

When creating a new Go workspace is too much work, the `goanywhere` tool can
create one on the fly. This works **only with local file system paths**:

```shell
$ go install github.com/u-root/gobusybox/src/cmd/goanywhere@latest

$ goanywhere ./u-root/cmds/core/{init,gosh} ./cpu/cmds/cpud -- u-root
```

`goanywhere` creates a workspace in a temporary directory with the given
modules, and then execs `u-root` in the workspace passing along the command
names.

> [!TIP]
>
> While workspaces are good for local compilation, they are not meant to be
> checked in to version control systems.
>
> For a non-workspace way of building multi-module initramfs images, read more
> in the [mkuimage README](https://github.com/u-root/mkuimage). (The `u-root`
> tool is `mkuimage` with more defaults applied.)

## Extra Files

You may also include additional files in the initramfs using the `-files` flag.

If you add binaries with `-files` are listed, their ldd dependencies will be
included as well.

```shell
$ u-root -files /bin/bash

$ cpio -ivt < /tmp/initramfs.linux_amd64.cpio
...
-rwxr-xr-x   0 root     root      1277936 Jan  1  1970 bin/bash
...
drwxr-xr-x   0 root     root            0 Jan  1  1970 lib/x86_64-linux-gnu
-rwxr-xr-x   0 root     root       210792 Jan  1  1970 lib/x86_64-linux-gnu/ld-linux-x86-64.so.2
-rwxr-xr-x   0 root     root      1926256 Jan  1  1970 lib/x86_64-linux-gnu/libc.so.6
lrwxrwxrwx   0 root     root           15 Jan  1  1970 lib/x86_64-linux-gnu/libtinfo.so.6 -> libtinfo.so.6.4
-rw-r--r--   0 root     root       216368 Jan  1  1970 lib/x86_64-linux-gnu/libtinfo.so.6.4
drwxr-xr-x   0 root     root            0 Jan  1  1970 lib64
lrwxrwxrwx   0 root     root           42 Jan  1  1970 lib64/ld-linux-x86-64.so.2 -> /lib/x86_64-linux-gnu/ld-linux-x86-64.so.2
...
```

You can determine placement with colons:

```shell
$ u-root -files "/bin/bash:sbin/sh"

$ cpio -ivt < /tmp/initramfs.linux_amd64.cpio
...
-rwxr-xr-x   0 root     root      1277936 Jan  1  1970 sbin/sh
...
```

For example on Debian, if you want to add two kernel modules for testing,
executing your currently booted kernel:

```shell
$ u-root -files "$HOME/hello.ko:etc/hello.ko" -files "$HOME/hello2.ko:etc/hello2.ko"
$ qemu-system-x86_64 -kernel /boot/vmlinuz-$(uname -r) -initrd /tmp/initramfs.linux_amd64.cpio
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
u-root -uinitcmd="echo Go Gopher" ./cmds/core/{init,echo,gosh}

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

Passing command line arguments like above is equivalent to passing the arguments
to uinit via a flags file in `/etc/uinit.flags`, see [Extra Files](#extra-files).

Additionally, you can pass arguments to uinit via the `uroot.uinitargs` kernel
parameters, for example:

```bash
u-root -uinitcmd="echo Gopher" ./cmds/core/{init,echo,gosh}

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
u-root -uinitcmd="echo Go Gopher" ./cmds/core/{init,gosh}
# 21:05:57 could not create symlink from "bin/uinit" to "echo": command or path "echo" not included in u-root build: specify -uinitcmd="" to ignore this error and build without a uinit
```

You can also refer to non-u-root-commands; they will be added as symlinks. We
don't presume to know whether your symlink target is correct or not.

This will build, but not work unless you add a /bin/foobar to the initramfs.

```bash
u-root -uinitcmd="/bin/foobar Go Gopher" ./cmds/core/{init,gosh}
```

This will boot the same as the above.

```bash
u-root -uinitcmd="/bin/foobar Go Gopher" -files /bin/echo:bin/foobar -files your-hosts-file:/etc/hosts ./cmds/core/{init,gosh}
```

The effect of the above command:
*   Sets up the uinit command to be /bin/foobar, with 2 arguments: Go Gopher
*   Adds /bin/echo as bin/foobar
*   Adds your-hosts-file as etc/hosts
*   builds in the cmds/core/init, and cmds/core/gosh commands.

This will bypass the regular u-root init and just launch a shell:

```bash
u-root -initcmd=gosh ./cmds/core/{gosh,ls}

cpio -ivt < /tmp/initramfs.linux_amd64.cpio
# ...
# lrwxrwxrwx   0 root     root            9 Dec 31  1969 init -> bbin/gosh

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

Further, we run integration tests on linux/amd64, and linux/arm64, using several
CI systems. If you need to add another CI system, processor or OS, please let us
know.

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

If you don't have a kernel handy, you can also get the one we use for VM testing:

```shell
go install github.com/hugelgupf/vmtest/tools/runvmtest@latest

runvmtest -- bash -c "cp \$VMTEST_KERNEL ./kernel"
```

It may not have all features you require, however.

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

More detailed information about the build process for a full LinuxBoot firmware
image using u-root/systemboot and coreboot can be found in the
[LinuxBoot book](https://github.com/linuxboot/book) chapter about
[LinuxBoot using coreboot, u-root and systemboot](https://github.com/linuxboot/book/blob/master/coreboot.u-root.systemboot/README.md).

This project started as a loose collection of programs in u-root by various
LinuxBoot contributors, as well as a personal experiment by
[Andrea Barberio](https://github.com/insomniacslk) that has since been merged
in. It is now an effort of a broader community and graduated to a real project
for system firmwares.

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

The u-root command supports building with workspace vendoring and module
vendoring. In both of those cases, if all dependencies are found in the vendored
directories, the build happens completely offline.

Read more in the [mkuimage README](https://github.com/u-root/mkuimage).

u-root also still supports `GO111MODULE=off` builds.

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

## Website

The sources of [u-root.org](https://u-root.org) are inside the `docs/` directory and
are deployed to the gh-pages branch. The CNAME file is currently not part of the CI
which deploys to the branch which shall be evaluated if this makes futures deployments easier.
