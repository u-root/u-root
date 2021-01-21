# Copyright 2020-2021 the u-root Authors. All rights reserved
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

FROM cimg/go:1.13

# Install dependencies
RUN sudo apt-get update &&                          \
    sudo apt-get install -y --no-install-recommends \
        `# Linux dependencies`                      \
        git                                         \
        bc                                          \
        bison                                       \
        flex                                        \
        gcc-aarch64-linux-gnu                       \
        libssl-dev                                  \
        make                                        \
        `# QEMU dependencies`                       \
        python                                      \
        libglib2.0-dev                              \
        libfdt-dev                                  \
        libpixman-1-dev                             \
        zlib1g-dev                                  \
        libcap-dev                                  \
        libattr1-dev                                \
        qemu-efi-aarch64                            \
        `# Linux kernel build deps`                 \
        libelf-dev &&                               \
    sudo rm -rf /var/lib/apt/lists/*

# Create working directory
WORKDIR /home/circleci
COPY config_linux5.2.0_arm64.txt .config

# Build linux
RUN set -eux;                                                             \
    git clone --depth=1 --branch=v5.2  https://github.com/torvalds/linux; \
    sudo chmod 0444 .config;                                              \
    mv .config linux/;                                                    \
    cd linux;                                                             \
    export ARCH=arm64;                                                    \
    export CROSS_COMPILE=/usr/bin/aarch64-linux-gnu-;                     \
    make olddefconfig;                                                    \
    make -j$(($(nproc) * 2 + 1));                                         \
    cd ~;                                                                 \
    cp linux/arch/arm64/boot/Image Image;                                 \
    rm -rf linux/

# Build QEMU
RUN set -eux;                                                          \
    git clone --depth=1 --branch=v3.1.0 https://github.com/qemu/qemu;  \
    cd qemu;                                                           \
    mkdir build;                                                       \
    cd build;                                                          \
    ../configure                                                       \
        --target-list=aarch64-softmmu                                  \
        --enable-virtfs                                                \
        --disable-docs                                                 \
        --disable-sdl                                                  \
        --disable-kvm;                                                 \
    make -j$(($(nproc) * 2 + 1));                                      \
    cd ~;                                                              \
    cp -rL qemu/build/pc-bios/ ~/pc-bios;                              \
    cp qemu/build/aarch64-softmmu/qemu-system-aarch64 .;               \
    rm -rf qemu/

# Generate UEFI file images
RUN dd if=/dev/zero bs=1M count=64 of=flash0.img;                                   \
    dd if=/usr/share/qemu-efi-aarch64/QEMU_EFI.fd bs=1M of=flash0.img conv=notrunc;

# Export paths to binaries.
ENV UROOT_KERNEL /home/circleci/Image
ENV UROOT_QEMU "/home/circleci/qemu-system-aarch64 -machine virt -cpu cortex-a57 -L /home/circleci/pc-bios -m 1G -pflash /home/circleci/flash0.img"
ENV UROOT_TESTARCH arm64
