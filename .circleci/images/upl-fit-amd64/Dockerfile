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
        iasl;

RUN pip3 install --break-system-packages pefile pylibfdt

SHELL ["/bin/bash", "-c"]
RUN cd upl-fit;                                                              \
    source ./edksetup.sh;                                                    \
    make -C BaseTools;                                                       \
	python3 UefiPayloadPkg/UniversalPayloadBuild.py -t GCC5 -b DEBUG         \
      --Fit -DCPU_TIMER_LIB_ENABLE=FALSE                                     \
      -DTIMER_SUPPORT=LAPIC -DHAND_OFF_FDT_ENABLE=TRUE                       \
      --add_cc_flag=-mno-sse --add_cc_flag=-mno-mmx;                         \
	cp /upl-fit/Build/UefiPayloadPkgX64/UniversalPayload.fit /UplFitX64.fit;

FROM scratch
COPY --from=base /UplFitX64.fit /UplFitX64.fit
