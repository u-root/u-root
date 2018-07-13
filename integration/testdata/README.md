# Minimal Go-Capable Linux Kernel

Build Setup:

- Build machine: Ubuntu 16.04 LTS
- Linux repo: github.com/torvalds/linux, v4.17 tag
- Go version: go version go1.10.3 linux/amd64

Minimal kernel config needed for Go:

    CONFIG_64BIT=y
    CONFIG_BINFMT_ELF=y
    CONFIG_BLK_DEV_INITRD=y
    CONFIG_DEVTMPFS=y
    CONFIG_EARLY_PRINTK=y
    CONFIG_EPOLL=y
    CONFIG_FUTEX=y
    CONFIG_PRINTK=y
    CONFIG_PROC_FS=y
    CONFIG_SERIAL_8250=y
    CONFIG_SERIAL_8250_CONSOLE=y
    CONFIG_TTY=y

Build Linux:

1. Run `make mrproper`.
2. Run `make tinyconfig`.
3. Append above lists to `.config`.
4. Run `make menuconfig`. Exit and save.
5. make -j$(($(nproc) * 2 + 1))

Build u-root:

1. `go get github.com/u-root/u-root`
2. `u-root -format=cpio -build=bb`

Test:

1. `qemu-system-x86_64 -kernel arch/x86_64/boot/bzImage -initrd /tmp/initramfs.linux_amd64.cpio -nographic -append 'earlyprintk=ttyS0 console=ttyS0'`
2. Exit with CTRL-A + X
