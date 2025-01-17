// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include <stdint.h>

// globals; initial values taken from pkg/boot/universalpayload/trampoline_amd64.go
long stack_top = 0xdeadbeef;
long hob_addr = 0xdeadbeef;
long entry_point = 0xdeadbeef;

void *addrOfStartU()
{
    extern void trampoline_startU();
    return trampoline_startU;
}

void *addrOfStackTopU()
{
    extern long stack_top;
    return &stack_top;
}

void *addrOfHobAddrU()
{
    extern long hob_addr;
    return &hob_addr;
}

void trampoline_startU()
{
    extern long stack_top;
    extern long hob_addr;
    extern long entry_point;

    asm volatile(
        // Load stack_top address into rax
        "leaq stack_top(%%rip), %%rax\n"
        "movq (%%rax), %%rax\n"
        "movq %%rax, %%rsp\n"

        // Load hob_addr into rax and move to rcx
        "leaq hob_addr(%%rip), %%rax\n"
        "movq (%%rax), %%rax\n"
        "movq %%rax, %%rcx\n"

        // Load entry_point into AX and jump
        "leaq entry_point(%%rip), %%rax\n"
        "movq (%%rax), %%rax\n"
        "jmp *%%rax\n"
        :
        : "m"(stack_top), "m"(hob_addr), "m"(entry_point)
        : "rax", "rcx" // Clobbered registers
    );
}