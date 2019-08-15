#!/bin/bash
set -v
set -euo pipefail
MEM=4096 # 4G of DRAM for the virtual machine

pushd "${EFI_WORKSPACE}"

BIOS_PATH=${EFI_WORKSPACE}/downloads/ovmf/usr/share/edk2.git/ovmf-x64/OVMF-pure-efi.fd
WINDOWS_INSTALLER_ISO=~/Downloads/windows_installer.iso

# An empty image can be created by
# `qemu-img create -f raw "${WORKSPACE}"/windows.img 20G`
WINDOWS_DISK="${WORKSPACE}"/windows.img

qemu-system-x86_64 -L .                   \
  --bios "${BIOS_PATH}"                   \
  -m "${MEM}"                             \
  -cdrom "${WINDOWS_INSTALLER_ISO}"       \
  -boot d                                 \
  -hda "${WINDOWS_DISK}"                  \
  -smp "$(nproc)"                         \
  -serial stdio                           \
  -net none                               \
  -enable-kvm                             \
  -s

stty sane


