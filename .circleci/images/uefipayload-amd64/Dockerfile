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
        python3-dev      \
        python3-pip      \
        gcc              \
        g++              \
        uuid-dev         \
        nasm             \
        bash             \
        libfdt-dev       \
        swig             \
        iasl;

RUN pip3 install --break-system-packages pefile pylibfdt

SHELL ["/bin/bash", "-c"]
RUN cd uefipayload;                                                    \
    source ./edksetup.sh;                                              \
    make -C BaseTools;                                                 \
    build -a X64 -p UefiPayloadPkg/UefiPayloadPkg.dsc -b DEBUG         \
      -t GCC5 -D BOOTLOADER=LINUXBOOT -D DISABLE_MMX_SSE=true;         \
    cp /uefipayload/Build/UefiPayloadPkgX64/DEBUG_GCC5/FV/UEFIPAYLOAD.fd /UEFIPAYLOAD.fd; \
    python3 UefiPayloadPkg/UniversalPayloadBuild.py -t GCC5 --Fit;

FROM scratch
COPY --from=base /UEFIPAYLOAD.fd /UEFIPAYLOAD.fd
COPY --from=base /uefipayload/Build/UefiPayloadPkgX64/UniversalPayload.fit /UEFIPayload.fit
