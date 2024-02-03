# Copyright 2018-2021 the u-root Authors. All rights reserved
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

FROM ubuntu:rolling as base

# Install dependencies
RUN apt-get update &&                          \
    apt-get install -y --no-install-recommends \
        ca-certificates                        \
        git                                    \
        make                                   \
        gcc-multilib                           \
        gzip                                   \
        gcc;

# Create working directory
WORKDIR /root

# Build Multiboot kernel
RUN git clone --depth=1 https://github.com/u-root/multiboot-test-kernel;
RUN cd multiboot-test-kernel;                                          \
    git checkout f5d529c940603ea9398a488f910142529387e4a6;             \
    make;

FROM scratch
COPY --from=base /root/multiboot-test-kernel/kernel ./mb/kernel
COPY --from=base /root/multiboot-test-kernel/kernel.gz ./mb/kernel.gz
