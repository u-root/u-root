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
		"mov %%cr0, %%rax\n"
		"and %0, %%rax\n"
		"mov %%rax, %%cr0\n"

		// disable long mode
		"mov %1, %%rcx\n"
		"rdmsr\n"
		"and %2, %%rax\n"
		"wrmsr\n"

		// disable physical address extension (pae)
		"xor %%rax, %%rax\n"
		"mov %%rax, %%cr4\n"

		// load data segments
		"mov $0x10, %%rax\n"
		"mov %%ax, %%ds\n"
		"mov %%ax, %%es\n"
		"mov %%ax, %%ss\n"
		"mov %%ax, %%fs\n"
		"mov %%ax, %%gs\n"

		// prepare long jump
		"mov %%rsi, %%rax\n"
		"jmp farjump32\n"
		:
		: "i"(CR0_PG), "i"(MSR_EFER), "i"(EFER_LME)
		: "rax", "rcx");
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

	uintptr_t farjump32_addr = (uintptr_t)farjump32;
	uintptr_t boot_addr = (uintptr_t)boot;
	uintptr_t farjump64_addr = (uintptr_t)farjump64;

	asm volatile(
		"mov %0, %%rbx\n"
		"mov %1, %%rsi\n"
		"mov %2, %%rax\n"
		"mov %%rax, (%3)\n"
		"lea (%4), %%rcx\n"
		"mov %%rcx, 6(%5)\n"
		"jmp *%5\n"
		:
		: "m"(info), "m"(magic), "m"(entry), "r"(farjump32_addr + 1), "r"(boot_addr), "r"(farjump64_addr)
		: "rax", "rbx", "rsi", "rcx");
}

void end(void) {}