# Copyright 2018-2019 the u-root Authors. All rights reserved
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

FROM circleci/golang:1.13

# Install dependencies
RUN sudo apt-get update &&                          \
    sudo apt-get install -y --no-install-recommends \
        `# Linux dependencies`                      \
        git                                         \
        bc                                          \
        bison                                       \
        flex                                        \
        gcc                                         \
        make                                        \
        `# QEMU dependencies`                       \
        libglib2.0-dev                              \
        libfdt-dev                                  \
        libpixman-1-dev                             \
        zlib1g-dev                                  \
        libcap-dev                                  \
        libattr1-dev                                \
        `# Linux kernel build deps`                 \
        libelf-dev                                  \
        `# Multiboot kernel build deps`             \
        gcc-multilib                                \
        gzip                                        \
        `# Edk2 build deps`                         \
        uuid-dev                                    \
        nasm                                        \
        bash                                        \
        iasl &&                                     \
    sudo rm -rf /var/lib/apt/lists/*

# Create working directory
WORKDIR /home/circleci
COPY config_linux4.17_x86_64.txt .config

# Build linux
RUN set -eux;                                                             \
    git clone --depth=1 --branch=v4.17 https://github.com/torvalds/linux; \
    sudo chmod 0444 .config;                                              \
    mv .config linux/;                                                    \
    cd linux;                                                             \
    make olddefconfig;                                                    \
    make -j$(($(nproc) * 2 + 1));                                         \
    cd ~;                                                                 \
    cp linux/arch/x86_64/boot/bzImage bzImage;                            \
    rm -rf linux/

# Build QEMU
RUN set -eux;                                                          \
    git clone --depth=1 --branch=v4.2.0 https://github.com/qemu/qemu;  \
    cd qemu;                                                           \
    mkdir build;                                                       \
    cd build;                                                          \
    ../configure                                                       \
        --target-list=x86_64-softmmu                                   \
        --enable-virtfs                                                \
        --disable-docs                                                 \
        --disable-sdl                                                  \
        --disable-kvm;                                                 \
    make -j$(($(nproc) * 2 + 1));                                      \
    cd ~;                                                              \
    cp -rL qemu/build/pc-bios/ ~/pc-bios;                              \
    cp qemu/build/x86_64-softmmu/qemu-system-x86_64 .;                 \
    rm -rf qemu/

# Build Multiboot kernel
RUN set -eux;                                                          \
    git clone --depth=1 --branch=v1.01                                 \
           https://github.com/u-root/multiboot-test-kernel;            \
    cd multiboot-test-kernel;                                          \
    make;                                                              \
    cd ~;                                                              \
    cp multiboot-test-kernel/kernel ./;                                \
    gzip -k kernel;                                                    \
    rm -rf multiboot-test-kernel/

SHELL ["/bin/bash", "-c"]
RUN set -ex;                                                           \
    git clone --branch uefipayload --recursive                         \
    https://github.com/chengchiehhuang/edk2 uefipayload;               \
    cd uefipayload;                                                    \
    source ./edksetup.sh;                                              \
    make -C BaseTools;                                                 \
    build -a X64 -p UefiPayloadPkg/UefiPayloadPkg.dsc -b DEBUG         \
    -t GCC5 -D BOOTLOADER=LINUXBOOT;                                   \
    cp Build/UefiPayloadPkgX64/DEBUG_GCC5/FV/UEFIPAYLOAD.fd ~/;        \
    cd ~;                                                              \
    rm -rf uefipayload/

# Export paths to binaries.
ENV UROOT_KERNEL /home/circleci/bzImage
ENV UROOT_QEMU "/home/circleci/qemu-system-x86_64 -L /home/circleci/pc-bios -m 1G"
ENV UROOT_TESTARCH amd64
