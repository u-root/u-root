#!/bin/bash
set -v
set -euo pipefail

# Get the first two arguments. Default is "none". We need this to avoid
# "unbound variable" error on $1 or $2
PARAM1="${1-none}"
PARAM2="${2-none}"

if [[ "${PARAM1}" == "help" ]]; then
cat <<USAGETEXT
  Usage:
  'run_vm.sh [rebuild_kernel] [rebuild_uroot]'
     - Run NERF (uroot). Rebuild the kernel or uroot (or both).
USAGETEXT
  exit
fi

MEM=4096 # 4G of DRAM for the virtual machine

pushd "${EFI_WORKSPACE}"

BIOS_PATH=${EFI_WORKSPACE}/downloads/ovmf/usr/share/edk2.git/ovmf-x64/OVMF-pure-efi.fd

if [[ "${PARAM1}" == "rebuild_kernel" ]] || \
   [[ "${PARAM2}" == "rebuild_kernel" ]]; then
  pushd linux
  make -j$(($(nproc)+2))
  popd
fi

if [[ "${PARAM1}" == "rebuild_uroot" ]] || \
   [[ "${PARAM2}" == "rebuild_uroot" ]]; then
  pushd efi_fs
  u-root -build=bb -files=bootmgfw.efi
  popd
fi


# We want to mask all writes to the original Windows image. We therefore
# create a new snampshot file, which is a delta file. We remove it before
# every run to make sure we start fresh with the original image.
WINDOWS_IMAGE="${WORKSPACE}"/windows.img
WINDOWS_DISK=windows_write_masking.img

if [[ ! -f "${WINDOWS_DISK}" ]]; then
  qemu-img create -f qcow2 -b "${WINDOWS_IMAGE}" "${WINDOWS_DISK}"
fi

KERNEL_PATH=linux/arch/x86_64/boot/bzImage
BOOT_FLAGS="earlyprintk=ttyS0 printk=ttyS0 console=ttyS0 root=/dev/sda1 noefi"

# Make sure we get the graphical terminal
export DISPLAY=:0

qemu-system-x86_64 -L .                   \
  -bios "${BIOS_PATH}"                    \
  -kernel "${KERNEL_PATH}"                \
  -initrd /tmp/initramfs.linux_amd64.cpio \
  -hda ${WINDOWS_DISK}                    \
  -append "${BOOT_FLAGS}"                 \
  -m ${MEM}                               \
  -smp "$(nproc)"                         \
  -serial stdio                           \
  -s

stty sane

