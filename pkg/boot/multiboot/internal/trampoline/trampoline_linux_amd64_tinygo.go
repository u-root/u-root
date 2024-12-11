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
#include "trampoline_tinygo.h"
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

// Setup scans memory for trampoline code and sets
// values for multiboot info address and kernel entry point.
func Setup(path string, magic, infoAddr, entryPoint uintptr) ([]byte, error) {
	trampolineStart, d, err := locateTrampoline()
	if err != nil {
		return nil, err
	}
	return patch(trampolineStart, d, magic, infoAddr, entryPoint)
}

// locateTrampoline locates the trampoline segment from the loaded binary.
// The trampoline segment begins after the "u-root-trampoline-begin" byte sequence + padding,
// and ends at the "u-root-trampoline-end" byte sequence.
func locateTrampoline() (uintptr, []byte, error) {
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

	// Trampoline is read-only executable memory. So we have to copy it to a
	// slice and modify it later.
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
