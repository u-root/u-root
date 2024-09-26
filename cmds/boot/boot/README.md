# uses of net in boot cmd

```
grep -rIi "net" ../../../pkg/{boot,cmdline,ulog,mount} | cut -d : -f 1 | uniq
```
../../../pkg/boot/bzimage/testdata/CONFIG
../../../pkg/boot/grub/echo_test.go
../../../pkg/boot/grub/testdata_new/rhel_7_8_another.json
../../../pkg/boot/grub/testdata_new/rhel_7_8_another/boot/grub2/grub.cfg
../../../pkg/boot/grub/grub.go
../../../pkg/boot/ibft/ibft.go
../../../pkg/boot/ibft/ibft_test.go
../../../pkg/boot/linux_test.go
../../../pkg/boot/netboot/ipxe/ipxe.go
../../../pkg/boot/netboot/ipxe/ipxe_test.go
../../../pkg/boot/netboot/netboot.go
../../../pkg/boot/netboot/pxe/pxe.go
../../../pkg/boot/netboot/pxe/pxe_test.go
../../../pkg/boot/netboot/simple/simple.go
../../../pkg/boot/purgatory/purgatory.go
../../../pkg/boot/syslinux/syslinux_test.go
../../../pkg/boot/syslinux/syslinux.go
../../../pkg/boot/systembooter/README.md
../../../pkg/boot/systembooter/bootentry.go
../../../pkg/boot/systembooter/bootentry_test.go
../../../pkg/boot/systembooter/booter.go
../../../pkg/boot/systembooter/netbooter.go
../../../pkg/cmdline/cmdline_test.go
../../../pkg/cmdline/filters_test.go
../../../pkg/mount/block/testdata/mounts
