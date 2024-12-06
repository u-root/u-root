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

	trampBegin := C.addrOfStartU()
	trampStack := C.addrOfStackTopU()
	trampHob := C.addrOfHobAddrU()

	tramp := ptrToSlice(trampBegin, int(trampStack-trampBegin))

	buf = append(buf, tramp...)

	padWithLength := func(slice []uint8, len uint64) []uint8 {
		tmpBytes := make([]uint8, len)
		return append(slice, tmpBytes...)
	}

	stackOffset := trampStack & 0xFFFF
	gapLen := stackOffset - (trampStack - trampBegin)
	buf = padWithLength(buf, uint64(gapLen))

	stackTop := hobAddr + tmpStackTop
	appendUint64 := func(slice []uint8, value uint64) []uint8 {
		tmpBytes := make([]uint8, 8)
		binary.LittleEndian.PutUint64(tmpBytes, value)
		return append(slice, tmpBytes...)
	}

	padLen := uint64(trampHob - trampStack - 8)

	buf = appendUint64(buf, stackTop)
	buf = padWithLength(buf, padLen)
	buf = appendUint64(buf, hobAddr)
	buf = padWithLength(buf, padLen)
	buf = appendUint64(buf, entry)

	return buf
}
