#!/bin/bash
set -euo pipefail
set -v
# Path to a working Windows *raw* image:
WINDOWS_DISK="${WORKSPACE}"/windows.img

START_PATH=$(realpath .)

### Downloading and extracting ovmf EFI firmawre ###
mkdir -p "${EFI_WORKSPACE}"/downloads
cd "${EFI_WORKSPACE}"/downloads
FIRMWARE_URL=https://www.kraxel.org/repos/jenkins/edk2/edk2.git-ovmf-x64-0-20190704.1206.g48d8d4d80b.noarch.rpm
FIRMWARE_IMAGE=$(basename ${FIRMWARE_URL})
echo "${FIRMWARE_IMAGE}"

if [[ ! -f ${FIRMWARE_IMAGE} ]]; then
  wget "${FIRMWARE_URL}" -O "${FIRMWARE_IMAGE}"

  fakeroot alien -d "${FIRMWARE_IMAGE}"

  # Assuming there is exatcly one .deb file:
  FIRMWARE_DEB=$(find . -name '*edk2.git-ovmf*deb*')

  rm -rf ovmf
  dpkg-deb -x "${FIRMWARE_DEB}" ovmf

else
  echo "NOTE: $(realpath "${FIRMWARE_IMAGE}") exists!! " \
       "Remove file to re-download EFI firmware." 1>&2
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

