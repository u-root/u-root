#!/bin/bash
set -euo pipefail

# Path to a working Windows *raw* image:
WINDOWS_DISK="${WORKSPACE}"/windows.img

START_PATH=$(realpath .)

### Extracting ovmf EFI firmawre ###
cd "${EFI_WORKSPACE}"
FIRMWARE_TAR="${START_PATH}"/ovmf_uefi.tar.gz
FIRMWARE_IMAGE=OVMF-pure-efi.fd

# The OVMF EFI can be downloaded from:
# https://www.kraxel.org/repos/jenkins/edk2/edk2.git-ovmf-x64-0-20190704.1206.g48d8d4d80b.noarch.rpm
# Exact filename might be different

echo "Will try to unpack ${FIRMWARE_TAR}"

if [[ ! -f ${FIRMWARE_IMAGE} ]]; then
  echo "Unpacking OVMF EFI"
  tar -xvf "${FIRMWARE_TAR}"

else
  echo "NOTE: $(realpath "${FIRMWARE_IMAGE}") exists!! "
fi

### Creating EFI filesystem dir ###
cd "${EFI_WORKSPACE}"
mkdir -p efi_fs

echo "Extracting Windows Loader from Windows Image, requires sudo:"

WINDOWS_LOADER=${EFI_WORKSPACE}/efi_fs/bootmgfw.efi

# Loop device for mounting Windows image. You may want to run
# `losetup --list` to see which device is available
LOOP_DEVICE=loop4

if [[ ! -f "${WINDOWS_LOADER}" ]]; then
  sudo losetup "${LOOP_DEVICE}" "${WINDOWS_DISK}"  # Attach raw disk to loop1
  sudo kpartx -a /dev/"${LOOP_DEVICE}"         # Create /dev/mapper/loop1* partitions
  sudo mkdir -p /mnt/win_disk
  sudo mount /dev/mapper/"${LOOP_DEVICE}"p2 /mnt/win_disk

  cp /mnt/win_disk/EFI/Microsoft/Boot/bootmgfw.efi "${WINDOWS_LOADER}"

  sudo umount /mnt/win_disk
  sudo kpartx -d /dev/"${LOOP_DEVICE}"         # Remove /dev/mapper paritions
  sudo losetup -d /dev/"${LOOP_DEVICE}"        # Dettach WINDOWS_DISK

else
  echo "NOTE: Windows loader already exists, " 1>&2
  echo "      no need to extract it again" 1>&2
fi

# Clone the forked Linux kernel and build it. There may be some pre-requisties
# missing.

if [[ ! -d "${EFI_WORKSPACE}/linux" ]]; then
  git clone https://github.com/oweisse/linux.git
fi

pushd linux
cp "$START_PATH/linux_config/dot_config" .config
make olddefconfig # populate config with default values which may be missing

echo "Installing libelf-dev, may prompt for sudo password:"
sudo apt-get install libelf-dev

