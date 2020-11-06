# My u-root cmd

`#!/bin/bash
 u-root -build=bb -uinitcmd="netdiskboot" -o payload.cpio \
 -files="config.json" \
 core \
 github.com/u-root/u-root/cmds/boot/netdiskboot
´
# Test QEMU cmd with args

`qemu-system-x86_64 \
-drive file=boot_S0.img,format=raw,if=virtio,media=disk \
-kernel bzImage \
-initrd payload.cpio \
-append "console=ttyS0 uroot.nohwrng uroot.uinitargs='-dryRun'" \
-nographic \
-m 32G \
-nic user,model=virtio-net-pci \
-device virtio-rng-pci \
-smp 8 ´

#Example Config
`
{
	"ImgURL":"https://blobs.9esec.io/os/nightly/Fedora-HWT-disk-31-buildserver-MBR.img.xz",
	"KernelPrefix":"vmlinuz-5.",
	"InitramPrefix":"initramfs-5",
	"Args":"ro earlyprintk=ttyS0,io,0x2f8,57600n1 console=ttyS0,io,0x2f8,57600n1 loglevel=7 cpuidle.off=1 iomem=relaxed no_timer_check verbose nokaslr audit=0 systemd.log_color=0 systemd.log_level=debug rootovl systemd.unified_cgroup_hierarchy=0 selinux=0 biosdevname=0 net.ifnames=0 root=/dev/vda1",
	"Device":"/dev/vda1"
}
`