// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package trampoline sets machine to a specific state defined by multiboot v1
// spec and jumps to the intended kernel.
//
// https://www.gnu.org/software/grub/manual/multiboot/multiboot.html#Machine-state.
package trampoline

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"unsafe"

	"github.com/u-root/u-root/pkg/ubinary"
)

const (
	trampolineEntry = "u-root-entry-long"
	trampolineInfo  = "u-root-info-long"
)

func start()
func end()

// funcPC gives the program counter of the given function.
//
//go:linkname funcPC runtime.funcPC
func funcPC(f interface{}) uintptr

// alignUp aligns x to a 0x10 bytes boundary.
// go compiler aligns TEXT parts at 0x10 bytes boundary.
func alignUp(x int) int {
	const mask = 0x10 - 1
	return (x + mask) & ^mask
}

// Setup scans file for trampoline code and sets
// values for multiboot info address and kernel entry point.
func Setup(path string, infoAddr, entryPoint uintptr) ([]byte, error) {
	d, err := extract(path)
	if err != nil {
		return nil, err
	}
	return patch(d, infoAddr, entryPoint)
}

// extract extracts trampoline segment from file.
// trampoline segment begins after "u-root-trampoline-begin" byte sequence + padding,
// and ends at "u-root-trampoline-end" byte sequence.
func extract(path string) ([]byte, error) {
	// TODO(https://github.com/golang/go/issues/35055): deal with
	// potentially non-contiguous trampoline. Rather than locating start
	// and end, we should locate start,boot,farjump{32,64},gdt,info,entry
	// individually and return one potentially really big trampoline slice.
	tbegin := funcPC(start)
	tend := funcPC(end)
	if tend <= tbegin {
		return nil, io.ErrUnexpectedEOF
	}
	tramp := ptrToSlice(tbegin, int(tend-tbegin))

	// tramp is read-only executable memory. So we gotta copy it to a
	// slice. Gotta modify it later.
	cp := append([]byte(nil), tramp...)
	return cp, nil
}

func ptrToSlice(ptr uintptr, size int) []byte {
	var data []byte

	sh := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	sh.Data = ptr
	sh.Len = size
	sh.Cap = size

	return data
}

// patch patches the trampoline code to store value for multiboot info address
// after "u-root-header-long" byte sequence + padding and value
// for kernel entry point, after "u-root-entry-long" byte sequence + padding.
func patch(trampoline []byte, infoAddr, entryPoint uintptr) ([]byte, error) {
	replace := func(d, label []byte, val uint32) error {
		buf := make([]byte, 4)
		ubinary.NativeEndian.PutUint32(buf, val)

		ind := bytes.Index(d, label)
		if ind == -1 {
			return fmt.Errorf("%q label not found in file", label)
		}
		ind = alignUp(ind + len(label))
		if len(d) < ind+len(buf) {
			return io.ErrUnexpectedEOF
		}
		copy(d[ind:], buf)
		return nil
	}

	if err := replace(trampoline, []byte(trampolineInfo), uint32(infoAddr)); err != nil {
		return nil, err
	}
	if err := replace(trampoline, []byte(trampolineEntry), uint32(entryPoint)); err != nil {
		return nil, err
	}
	return trampoline, nil
}
