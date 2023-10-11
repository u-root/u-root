# VMBoot - Experimental
VMBoot starts a virtual machine based on [gokvm](github.com/bobuhiro11/gokvm) and executes [EDK2/CloudHV](https://github.com/cloud-hypervisor/edk2/tree/ch) via PVH boot protocol.
The EDK2/CloudHV firmware image must be placed at a /dev/sda1 with a xfs file system.

## Introduction
Ever wanted to start a system with Open Source, non-UEFI firmware (coreboo+Linuxboot/u-root) and still be able to boot UEFI/EDK2? Look no further.
VMBoot allows to execute EDK2 in a VM started from Linuxboot/u-root.
Why you ask? Because noone wants to implement UEFI-compliance in u-root and lose their sanity.

Booting EDK2 in the VM technically allows booting into UEFI-compliant or UEFI-required operating systems without relying on UEFI as host system firmware.
(Though gokvm is not able to do so yet!)

gokvm allows to execute [EDK2/CloudHV](https://github.com/cloud-hypervisor/edk2/tree/ch) until EFI-Shell.
CloudHV is an adaption of [OVMF](https://github.com/tianocore/tianocore.github.io/wiki/OVMF) specifically for [Cloud-Hypervisor](https://github.com/cloud-hypervisor/cloud-hypervisor).
Focusing on paravirtualisation with EDK2/CloudHV keeps emulation minimal on the hosts side.

## Hurdles
Flashchip size may be an issue. It needs to hold an linux kernel with more drivers to support the hardware and kvm.

## Best practices
- Build linux kernel with either AMD-V or Intel-VT support, not both.
- Reduce hardware support only to required by actual devices in the system.
- Reduce filesystem, compression, and other support to bare minimum.

## Example Linux config files
[Linux config with Intel-VT support](./linux_intel.config)

## Links
[VMBoot project repo](https://www.github.com/9elements/vmboot)