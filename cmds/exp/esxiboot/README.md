# Running ESXi with u-root and QEMU/KVM

This is a simple set of commands to get ESXi in QEMU/kvm up and running.

To run ESXi in QEMU you'll need:

-   QEMU/KVM
-   Upstream u-root
-   VMware installer image
-   Linux kernel

Linux kernel should support:

-   kexec system call
-   Sysfs
-   Vfat
-   AHCI
-   ATA
-   CDROM (iso9660, SCSI)

Sample kernel config can be found
[here](https://github.com/u-root/u-root/blob/14104d15a19773171441f5667ed5d2dce7a7da07/.circleci/images/integration/config_linux4.17_amd64.txt).

Also add the following to your Kconfig:

```
# qemu -hda support
CONFIG_BLOCK=y
CONFIG_ATA=y
CONFIG_ATA_PIIX=y

# qemu -cdrom support
CONFIG_SCSI=y
CONFIG_BLK_DEV_SD=y
CONFIG_BLK_DEV_SR=y
CONFIG_ISO9660_FS=y
```

`VMware-VMvisor-Installer-6.7.0-8169922.x86_64.iso` was used to write this
manual.

```shell
$ qemu-system-x86_64 --version
QEMU emulator version 2.11.93 (Debian 1:2.12~rc3+dfsg-1)
Copyright (c) 2003-2017 Fabrice Bellard and the QEMU Project developers

$ cat /etc/modprobe.d/kvm-intel.conf
options kvm-intel nested=y

$ uname -r
4.18.10-1rodete2-amd64
```

## Preparing the VM

```shell
$ mkdir $HOME/esxi
$ cd $HOME/esxi

$ cp <VMware installer image> ./vmware.iso
$ cp <bzImage> ./

# Create QEMU disk image.
$ qemu-img create -f qcow2 -o nocow=on esxi.qcow2 16G

# Download and install u-root.
$ go get -u github.com/u-root/u-root

# Create initramfs.
$ u-root -o ./initramfs.linux_amd64.cpio --build=bb all github.com/u-root/u-root/cmds/exp/esxiboot
```

## Installing ESXi

```shell
# Run QEMU (8G RAM is required minimum, 2 CPUs minimum).
$ qemu-system-x86_64 -cpu host -smp 2 -m 8192 -enable-kvm -kernel ./bzImage \
    -initrd ./initramfs.linux_amd64.cpio -hda esxi.qcow2 \
    -cdrom vmware.iso
```

Now you are inside the QEMU VM.

```shell
# Kexec ESXi installer.
$ esxiboot -r /dev/sr0
```

Follow GUI installer instructions. Shutdown QEMU after ESXi is installed.

## Running ESXi

Run ESXi without the `-cdrom` boot disk:

```shell
$ qemu-system-x86_64 -cpu host -smp 2 -m 8192 -enable-kvm -kernel ./bzImage \
    -initrd ./initramfs.linux_amd64.cpio -hda esxi.qcow2
```

Now you are inside the QEMU VM.

```shell
# Kexec ESXi. ESXi is normally installed to partition 5.
$ esxiboot -d /dev/sda -p 5
```

In QEMU console (Ctrl+Alt+2):

```shell
hostfwd_add tcp::4443-:443
```

In browser go to [https://localhost:4443](https://localhost:4443).

You're all set.
