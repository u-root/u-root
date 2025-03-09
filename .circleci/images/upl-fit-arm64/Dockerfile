# Copyright 2025 the u-root Authors. All rights reserved
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

FROM ubuntu:rolling AS base

# Install dependencies
RUN apt-get update &&                          \
    apt-get install -y --no-install-recommends \
        ca-certificates                             \
        git;

RUN git clone --branch master --recursive https://github.com/tianocore/edk2 upl-fit;

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
        iasl             \
        gcc-aarch64-linux-gnu;

RUN pip3 install --break-system-packages pefile pylibfdt

SHELL ["/bin/bash", "-c"]
RUN cd upl-fit;                                                              \
    export GCC5_AARCH64_PREFIX=aarch64-linux-gnu-;                           \
    source ./edksetup.sh;                                                    \
    make -C BaseTools;                                                       \
	python3 UefiPayloadPkg/UniversalPayloadBuild.py -a AARCH64 -t GCC5       \
      -b DEBUG --Fit -DHAND_OFF_FDT_ENABLE=TRUE;                             \
	cp /upl-fit/Build/UefiPayloadPkgAARCH64/UniversalPayload.fit /UplFitARM64.fit;

FROM scratch
COPY --from=base /upl-fit/Build/UefiPayloadPkgAARCH64/UniversalPayload.fit /UplFitARM64.fit
