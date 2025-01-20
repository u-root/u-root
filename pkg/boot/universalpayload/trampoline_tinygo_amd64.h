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
    // We are using AT&T inline assembly code to implement AMD64 trampoline
    // bootstrap code. In inline assembly code, we use rip-relative addressing
    // to fetch corresponding values of stack top, boot parameter and the
    // real entry point of Universal Payload FIT image.
    //
    // In this way, we need to update the actual value of stack top, boot
    // parameter and entry point after relocation.
    //
    // The layout of the trampoline code is shown as follows in Intel Syntax.
    // Then maintainer, who is familiar with either AT&T or Intel Syntax,
    // can catch up with the purpose of trampoline code quickly.
    //
    // trampoline[32 - 55] will be filled in utilities_arch_amd64_tinygo.go,
    // keep them here for easy understanding.
    //
    // trampoline[0 - 6]   : mov rax, qword ptr [rip+0x19]
    // trampoline[7 - 9]   : mov rsp, rax
    // trampoline[10 - 16] : mov rax, qword ptr [rip+0x17]
    // trampoline[17 - 19] : mov rcx, rax
    // trampoline[20 - 26] : mov rax, qword ptr [rip+0x15]
    // trampoline[27 - 28] : jmp rax
    // trampoline[29 - 31] : padding for alignment
    // trampoline[32 - 39] : Top of stack address
    // trampoline[40 - 47] : Base address of bootloader parameter
    // trampoline[48 - 55] : Entry point of FIT image

    __asm__ __volatile__ (
        "movq  0x19(%%rip), %%rax\n\t"       // Load value stack_top into RAX
        "mov   %%rax, %%rsp\n\t"             // Move RAX to RSP
        "movq  0x17(%%rip), %%rax\n\t"       // Load value hob_addr into RAX
        "mov   %%rax, %%rcx\n\t"             // Move RAX to RCX
        "movq  0x15(%%rip), %%rax\n\t"       // Load value entry_point into RAX
        "jmp   *%%rax\n\t"                   // Jump to address in RAX
        "int3\n\t"                           // Software BP for alignment
        "int3\n\t"                           // Software BP for alignment
        "int3\n\t"                           // Software BP for alignment
        ::: "rax", "rcx", "rsp"              // Clobbered registers
    );

}
