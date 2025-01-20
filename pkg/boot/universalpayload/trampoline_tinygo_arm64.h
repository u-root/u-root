// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include <stdint.h>

uintptr_t addrOfStartU()
{
    extern void trampoline_startU();
    return (uintptr_t)trampoline_startU;
}

void trampoline_startU()
{
    // We are using inline assembly code to implement ARM64 trampoline
    // bootstrap code. In inline assembly code, we use pc-relative addressing
    // to fetch corresponding value of stack top, boot parameter and the
    // real entry point of Universal Payload FIT image.
    //
    // In this way, we need to update the actual value of stack, parameter
    // and entry point after relocation.
    //
    // trampoline[32 - 55] will be filled in utilities_arch_arm64_tinygo.go,
    // keep them here for easy understanding.
    //
    // trampoline[0 - 3]   : ldr x4, #0x30 (PC relative: buf[48 - 55], entry_point)
    // trampoline[4 - 7]   : ldr x0, #0x24 (PC relative: buf[40 - 47], hob_addr)
    // trampoline[8 - 11]  : clear x1
    // trampoline[12 - 15] : ldr x2, #0x14 (PC relative: buf[32 - 39], stack_top)
    // trampoline[16 - 19] : mov sp, x2
    // trampoline[20 - 23] : clear x2
    // trampoline[24 - 27] : clear x3
    // trampoline[28 - 31] : br x4
    // trampoline[32 - 39] : Top of stack address
    // trampoline[40 - 47] : Base address of bootloader parameter
    // trampoline[48 - 55] : Entry point of FIT image

	__asm__ __volatile__ (
        "ldr x4, [pc, #0x30]\n\t"
        "ldr x0, [pc, #0x24]\n\t"
        "mov x1, xzr\n\t"
        "ldr x2, [pc, #0x14]\n\t"
        "mov sp, x2\n\t"
        "mov x2, xzr\n\t"
        "mov x3, xzr\n\t"
        "br  x4\n\t"
        ::: "x0", "x1", "x2", "x3", "x4", "sp"    // Clobbered registers
    );

}
