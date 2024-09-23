// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux && amd64 && tinygo

// Package trampoline sets machine to a specific state defined by multiboot v1
// spec and jumps to the intended kernel.
//
// https://www.gnu.org/software/grub/manual/multiboot/multiboot.html#Machine-state.

package trampoline

// In Go 1.17+, Go references to assembly functions resolve to an ABIInternal
// wrapper function rather than the function itself. We must reference from
// assembly to get the ABI0 (i.e., primary) address (this way of doing things
// will work for both 1.17+ and versions prior to 1.17). Note for posterity:
// runtime.funcPC (used previously) is going away in 1.18+.
//
// Each of the functions below of form 'addrOfXXX' return the starting PC
// of the assembly routine XXX.

/*
// "textflag.h" is provided by the gc compiler, tinygo does not have this
#include <stdint.h>

#define MSR_EFER    0xC0000080
#define EFER_LME    0xFFFFFEFF
#define CR0_PG      0x0FFFFFFF

#define DATA_SEGMENT    0x00CF92000000FFFF
#define CODE_SEGMENT    0x00CF9A000000FFFF

uintptr_t addrOfStart() {
	extern void start();
	return (uintptr_t)start;
}
uintptr_t addrOfEnd() {
	extern void end();
	return (uintptr_t)end;
}
uintptr_t addrOfInfo() {
	extern uint32_t info;
	return (uintptr_t)&info;
}
uintptr_t addrOfMagic() {
	extern uint32_t magic;
	return (uintptr_t)&magic;
}
uintptr_t addrOfEntry() {
	extern uint32_t entry;
	return (uintptr_t)&entry;
}
void farjump32() {
	asm volatile (
		".byte 0xEA\n"
		".long 0x0\n"
		".word 0x18\n"
	);
}
void farjump64() {
	asm volatile (
		".byte 0xFF, 0x2D\n"
		".long 0x0\n"
		".long 0x0\n"
		".long 0x8\n"
	);
}
void boot() {
	asm volatile (
		"movl %%cr0, %%eax\n"
		"andl %0, %%eax\n"
		"movl %%eax, %%cr0\n"
		"movl %1, %%ecx\n"
		"rdmsr\n"
		"andl %2, %%eax\n"
		"wrmsr\n"
		"xorl %%eax, %%eax\n"
		"movl %%eax, %%cr4\n"
		"movl $0x10, %%eax\n"
		"mov %%ax, %%ds\n"
		"mov %%ax, %%es\n"
		"mov %%ax, %%ss\n"
		"mov %%ax, %%fs\n"
		"mov %%ax, %%gs\n"
		"movl %%esi, %%eax\n"
		"jmp farjump32\n"
		:
		: "i"(CR0_PG), "i"(MSR_EFER), "i"(EFER_LME)
		: "eax", "ecx"
	);
}
void start() {
	uint64_t gdt[4] = {0x0, CODE_SEGMENT, DATA_SEGMENT, CODE_SEGMENT};
	uint64_t gdt_ptr[2];

	gdt_ptr[0] = (sizeof(gdt) - 1) | ((uint64_t)(uintptr_t)gdt << 16);
	gdt_ptr[1] = (uintptr_t)gdt >> 48;

	asm volatile (
		"lgdt (%0)\n"
		:
		: "r"(gdt_ptr)
	);

	extern uint32_t info;
	extern uint32_t magic;
	extern uint32_t entry;

	asm volatile (
		"movl %0, %%ebx\n"
		"movl %1, %%esi\n"
		"movl %2, %%eax\n"
		"movl %%eax, farjump32+1\n"
		"leaq boot, %%rcx\n"
		"movl %%ecx, farjump64+6\n"
		"jmp farjump64\n"
		:
		: "m"(info), "m"(magic), "m"(entry)
		: "eax", "ebx", "esi", "ecx"
	);
}
uint32_t info = 0x0;
uint32_t entry = 0x0;
uint32_t magic = 0x0;
void end() {}
*/
import "C"

import (
	"encoding/binary"
	"io"
	"unsafe"
)

const (
	trampolineEntry = "u-root-entry-long"
	trampolineInfo  = "u-root-info-long"
	trampolineMagic = "u-root-mb-magic"
)

// Setup scans file for trampoline code and sets
// values for multiboot info address and kernel entry point.
func Setup(path string, magic, infoAddr, entryPoint uintptr) ([]byte, error) {
	trampolineStart, d, err := extract(path)
	if err != nil {
		return nil, err
	}
	return patch(trampolineStart, d, magic, infoAddr, entryPoint)
}

// extract extracts trampoline segment from file.
// trampoline segment begins after "u-root-trampoline-begin" byte sequence + padding,
// and ends at "u-root-trampoline-end" byte sequence.
func extract(path string) (uintptr, []byte, error) {
	// TODO(https://github.com/golang/go/issues/35055): deal with
	// potentially non-contiguous trampoline. Rather than locating start
	// and end, we should locate start,boot,farjump{32,64},gdt,info,entry
	// individually and return one potentially really big trampoline slice.
	tbegin := C.addrOfStart()
	tend := C.addrOfEnd()
	if tend <= tbegin {
		return 0, nil, io.ErrUnexpectedEOF
	}
	tramp := ptrToSlice(tbegin, int(tend-tbegin))

	// tramp is read-only executable memory. So we gotta copy it to a
	// slice. Gotta modify it later.
	cp := append([]byte(nil), tramp...)
	return tbegin, cp, nil
}

func ptrToSlice(ptr uintptr, size int) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(ptr)), size)
}

// patch patches the trampoline code to store value for multiboot info address,
// entry point, and boot magic value.
//
// All 3 are determined by pretending they are functions, and finding their PC
// within our own address space.
func patch(trampolineStart uintptr, trampoline []byte, magicVal, infoAddr, entryPoint uintptr) ([]byte, error) {
	replace := func(start uintptr, d []byte, fPC uintptr, val uint32) error {
		buf := make([]byte, 4)
		binary.NativeEndian.PutUint32(buf, val)

		offset := fPC - start
		if int(offset+4) > len(d) {
			return io.ErrUnexpectedEOF
		}
		copy(d[int(offset):], buf)
		return nil
	}

	if err := replace(trampolineStart, trampoline, C.addrOfInfo(), uint32(infoAddr)); err != nil {
		return nil, err
	}
	if err := replace(trampolineStart, trampoline, C.addrOfEntry(), uint32(entryPoint)); err != nil {
		return nil, err
	}
	if err := replace(trampolineStart, trampoline, C.addrOfMagic(), uint32(magicVal)); err != nil {
		return nil, err
	}
	return trampoline, nil
}
