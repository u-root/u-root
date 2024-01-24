# Copyright 2018-2021 the u-root Authors. All rights reserved
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

FROM ubuntu:rolling AS base

# Install dependencies
RUN apt-get update &&                          \
    apt-get install -y --no-install-recommends \
        ca-certificates                             \
        git;

RUN git clone --branch uefipayload-2024 --recursive https://github.com/linuxboot/edk2 uefipayload;

RUN apt-get install -y --no-install-recommends \
        make             \
        python3          \
        gcc              \
        g++              \
        uuid-dev                                    \
        nasm                                        \
        bash                                        \
        iasl &&                                     \
    rm -rf /var/lib/apt/lists/*

SHELL ["/bin/bash", "-c"]
RUN cd uefipayload;                                                    \
    source ./edksetup.sh;                                              \
    make -C BaseTools;                                                 \
    build -a X64 -p UefiPayloadPkg/UefiPayloadPkg.dsc -b DEBUG         \
      -t GCC5 -D BOOTLOADER=LINUXBOOT -D DISABLE_MMX_SSE=true;

FROM scratch
COPY --from=base /uefipayload/Build/UefiPayloadPkgX64/DEBUG_GCC5/FV/UEFIPAYLOAD.fd /UEFIPAYLOAD.fd
