# Copyright 2018-2021 the u-root Authors. All rights reserved
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

FROM cimg/go:1.21

# Install dependencies
RUN sudo apt-get update &&                          \
    sudo apt-get install -y --no-install-recommends \
        `# Linux dependencies`                      \
        bc                                          \
        bison                                       \
        flex                                        \
        gcc-arm-linux-gnueabi                       \
        git                                         \
        make                                        \
        `# QEMU dependencies`                       \
        libattr1-dev                                \
        libcap-dev                                  \
        libcap-ng-dev                               \
        libfdt-dev                                  \
        libglib2.0-dev                              \
        libpixman-1-dev                             \
        meson                                       \
        ninja-build                                 \
        python3                                     \
        zlib1g-dev                                  \
        `# Linux kernel build deps`                 \
        libelf-dev &&                               \
    sudo rm -rf /var/lib/apt/lists/*

# Create working directory
WORKDIR /home/circleci
COPY config_linux5.10.0_arm.txt .config

# Build linux
RUN set -eux;                                                             \
    git clone --depth=1 --branch=v5.10 https://github.com/torvalds/linux; \
    sudo chmod 0444 .config;                                              \
    mv .config linux/;                                                    \
    cd linux;                                                             \
    export ARCH=arm;                                                      \
    export CROSS_COMPILE=/usr/bin/arm-linux-gnueabi-;                     \
    make -j$(($(nproc) * 2 + 1));                                         \
    cd ~;                                                                 \
    cp linux/arch/arm/boot/zImage zImage;                                 \
    rm -rf linux/

# Build QEMU
RUN set -eux;                                                          \
    git clone --depth=1 --branch=v7.0.0 https://github.com/qemu/qemu; \
    cd qemu;                                                           \
    mkdir build;                                                       \
    cd build;                                                          \
    ../configure                                                       \
        --target-list=arm-softmmu                                      \
        --enable-virtfs                                                \
        --disable-docs                                                 \
        --disable-sdl                                                  \
        --disable-kvm;                                                 \
    make -j$(($(nproc) * 2 + 1));                                      \
    cd ~;                                                              \
    cp -rL qemu/build/pc-bios/ ~/pc-bios;                              \
    cp qemu/build/arm-softmmu/qemu-system-arm .;                       \
    rm -rf qemu/

# Export paths to binaries.
ENV UROOT_KERNEL /home/circleci/zImage
ENV UROOT_QEMU "/home/circleci/qemu-system-arm -M virt -L /home/circleci/pc-bios"
ENV UROOT_TESTARCH arm
