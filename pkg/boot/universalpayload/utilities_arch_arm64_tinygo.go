// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build arm64 && tinygo

package universalpayload

/*
// "textflag.h" is provided by the gc compiler, tinygo does not have this
#include "trampoline_tinygo_arm64.h"
*/
import "C"

import (
	"encoding/binary"
	"reflect"
	"unsafe"
)

func getPhysicalAddressSizes() (uint8, error) {
	// Return hardcode for arm64
	// Please update to actual physical address size
	return 44, nil
}

// Construct trampoline code before jump to entry point of FIT image.
// Due to lack of support to set value of Registers in kexec,
// bootloader parameter needs to be prepared in trampoline code.
// Also stack is prepared in trampoline code snippet to ensure no data leak.
func constructTrampoline(buf []uint8, hobAddr uint64, entry uint64) []uint8 {
	ptrToSlice := func(ptr uintptr, size int) []byte {
		var data []byte

		sh := (*reflect.SliceHeader)(unsafe.Pointer(&data))
		sh.Data = ptr
		sh.Len = size
		sh.Cap = size

		return data
	}

	appendUint64 := func(slice []uint8, value uint64) []uint8 {
		tmpBytes := make([]uint8, 8)
		binary.LittleEndian.PutUint64(tmpBytes, value)
		return append(slice, tmpBytes...)
	}

	trampBegin := C.addrOfStartU()

	// Please keep 'size' parameter of 'ptrToSlice" align with implementation of
	// trampoline_startU in trampoline_tinygo_arm64.h
	tramp := ptrToSlice(trampBegin, 32)

	buf = append(buf, tramp...)

	stackTop := hobAddr + tmpStackTop
	buf = appendUint64(buf, stackTop)
	buf = appendUint64(buf, hobAddr)
	buf = appendUint64(buf, entry)

	return buf
}
