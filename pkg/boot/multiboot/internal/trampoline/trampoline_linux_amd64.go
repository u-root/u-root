// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package trampoline sets machine to a specific state defined by multiboot v1
// spec and jumps to the intended kernel.
//
// https://www.gnu.org/software/grub/manual/multiboot/multiboot.html#Machine-state.
package trampoline

import (
	"io"
	"reflect"
	"unsafe"

	"github.com/u-root/u-root/pkg/ubinary"
)

const (
	trampolineEntry = "u-root-entry-long"
	trampolineInfo  = "u-root-info-long"
	trampolineMagic = "u-root-mb-magic"
)

func start()
func end()
func info()
func magic()
func entry()

// funcPC gives the program counter of the given function.
//
//go:linkname funcPC runtime.funcPC
func funcPC(f interface{}) uintptr

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
	tbegin := funcPC(start)
	tend := funcPC(end)
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
	var data []byte

	sh := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	sh.Data = ptr
	sh.Len = size
	sh.Cap = size

	return data
}

// patch patches the trampoline code to store value for multiboot info address,
// entry point, and boot magic value.
//
// All 3 are determined by pretending they are functions, and finding their PC
// within our own address space.
func patch(trampolineStart uintptr, trampoline []byte, magicVal, infoAddr, entryPoint uintptr) ([]byte, error) {
	replace := func(start uintptr, d []byte, f func(), val uint32) error {
		buf := make([]byte, 4)
		ubinary.NativeEndian.PutUint32(buf, val)

		offset := funcPC(f) - start
		if int(offset+4) > len(d) {
			return io.ErrUnexpectedEOF
		}
		copy(d[int(offset):], buf)
		return nil
	}

	if err := replace(trampolineStart, trampoline, info, uint32(infoAddr)); err != nil {
		return nil, err
	}
	if err := replace(trampolineStart, trampoline, entry, uint32(entryPoint)); err != nil {
		return nil, err
	}
	if err := replace(trampolineStart, trampoline, magic, uint32(magicVal)); err != nil {
		return nil, err
	}
	return trampoline, nil
}
