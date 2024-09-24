// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package trampoline sets machine to a specific state defined by multiboot v1
// spec and jumps to the intended kernel.
//
// https://www.gnu.org/software/grub/manual/multiboot/multiboot.html#Machine-state.

#include <stdint.h>

#define MSR_EFER 0xC0000080
#define EFER_LME 0xFFFFFEFF
#define CR0_PG 0x0FFFFFFF

#define DATA_SEGMENT 0x00CF92000000FFFF
#define CODE_SEGMENT 0x00CF9A000000FFFF

// globals
uint32_t info = 0x0;
uint32_t entry = 0x0;
uint32_t magic = 0x0;
void end() {}

uintptr_t addrOfStart()
{
	extern void start();
	return (uintptr_t)start;
}

uintptr_t addrOfEnd()
{
	extern void end();
	return (uintptr_t)end;
}

uintptr_t addrOfInfo()
{
	extern uint32_t info;
	return (uintptr_t)&info;
}

uintptr_t addrOfMagic()
{
	extern uint32_t magic;
	return (uintptr_t)&magic;
}

uintptr_t addrOfEntry()
{
	extern uint32_t entry;
	return (uintptr_t)&entry;
}

void farjump32()
{
	asm volatile(
		".byte 0xEA\n"
		".long 0x0\n"
		".word 0x18\n");
}

void farjump64()
{
	asm volatile(
		".byte 0xFF, 0x2D\n"
		".long 0x0\n"
		".long 0x0\n"
		".long 0x8\n");
}

void boot()
{
	asm volatile(
		// disable paging
		"movl %%cr0, %%eax\n"
		"andl %0, %%eax\n"
		"movl %%eax, %%cr0\n"

		// disable long mode
		"movl %1, %%ecx\n"
		"rdmsr\n"
		"andl %2, %%eax\n"
		"wrmsr\n"

		// disable physical address extension (pae)
		"xorl %%eax, %%eax\n"
		"movl %%eax, %%cr4\n"

		// load data segments
		"movl $0x10, %%eax\n"
		"mov %%ax, %%ds\n"
		"mov %%ax, %%es\n"
		"mov %%ax, %%ss\n"
		"mov %%ax, %%fs\n"
		"mov %%ax, %%gs\n"

		// prepare long jump
		"movl %%esi, %%eax\n"
		"jmp farjump32\n"
		:
		: "i"(CR0_PG), "i"(MSR_EFER), "i"(EFER_LME)
		: "eax", "ecx");
}

void start()
{
	uint64_t gdt[4] = {0x0, CODE_SEGMENT, DATA_SEGMENT, CODE_SEGMENT};
	uint64_t gdt_ptr[2] = {
		[0] = (sizeof(gdt) - 1) | ((uint64_t)(uintptr_t)gdt << 16),
		[1] = (uintptr_t)gdt >> 48};

	asm volatile(
		"lgdt (%0)\n"
		:
		: "r"(gdt_ptr));

	extern uint32_t info;
	extern uint32_t magic;
	extern uint32_t entry;

	extern void farjump32();
	extern void boot();
	extern void farjump64();

	asm volatile(
		"movl %0, %%ebx\n"
		"movl %1, %%esi\n"
		"movl %2, %%eax\n"
		"movl %%eax, farjump32+1\n"
		"leaq boot, %%rcx\n"
		"movl %%ecx, farjump64+6\n"
		"jmp farjump64\n"
		:
		: "m"(info), "m"(magic), "m"(entry)
		: "eax", "ebx", "esi", "ecx");
}
