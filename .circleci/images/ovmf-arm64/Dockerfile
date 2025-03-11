# Copyright 2025 the u-root Authors. All rights reserved
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

FROM ubuntu:rolling as base

# Install dependencies
RUN apt-get update &&                          \
    apt-get install -y --no-install-recommends \
        qemu-system-arm

FROM scratch
COPY --from=base /usr/share/qemu-efi-aarch64/QEMU_EFI.fd /QEMU_EFI.fd
