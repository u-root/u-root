# Copyright 2026 the u-root Authors. All rights reserved
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

FROM ubuntu:24.04 as base

# Install dependencies
RUN apt-get update &&                          \
    apt-get install -y --no-install-recommends \
        qemu-efi-riscv64

FROM scratch
COPY --from=base /usr/share/qemu-efi-riscv64/RISCV_VIRT_CODE.fd /RISCV_VIRT_CODE.fd
COPY --from=base /usr/share/qemu-efi-riscv64/RISCV_VIRT_VARS.fd /RISCV_VIRT_VARS.fd
