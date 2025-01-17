// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include <stdint.h>

// globals; initial values taken from pkg/boot/universalpayload/trampoline_amd64.go
long stack_top;
long hob_addr;
long entry_point;

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
    asm volatile(
        // Load the address of entry_point into x4 and dereference it
        "adrp    x4, entry_point\n"
        "add     x4, x4, #:lo12:entry_point\n" // use lower 12 bits of entry_point for page alignment
        "ldr     x4, [x4]\n"

        // Load the address of hob_addr into x0 and dereference it
        "adrp    x0, hob_addr\n"
        "add     x0, x0, #:lo12:hob_addr\n"
        "ldr     x0, [x0]\n"

        // Load the address of stack_top into x2, dereference it, and set the stack pointer
        "adrp    x2, stack_top\n"
        "add     x2, x2, #:lo12:stack_top\n"
        "ldr     x2, [x2]\n"
        "mov     sp, x2\n"

        // Jump to the address stored in x4
        "br      x4\n"
        :
        : "m"(stack_top), "m"(hob_addr), "m"(entry_point)
        : "x0", "x2", "x4" // Clobbered registers
    );
}