// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build arm64 && !tinygo

package universalpayload

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

func addrOfStart() uintptr
func addrOfStackTop() uintptr
func addrOfHobAddr() uintptr

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

	trampBegin := addrOfStart()
	trampStack := addrOfStackTop()
	trampHob := addrOfHobAddr()

	tramp := ptrToSlice(trampBegin, int(trampStack-trampBegin))

	buf = append(buf, tramp...)

	padWithLength := func(slice []uint8, len uint64) []uint8 {
		tmpBytes := make([]uint8, len)
		return append(slice, tmpBytes...)
	}

	// Due to Golang Plan9 Assembly support limitation, we can only
	// fetch symbol address after relocated, and symbol address of
	// trampBegin, trampStack, trampHob should not be larger than
	// one page from PC address of trampoline entry point. If symbol
	// address is larger than one page size from PC address of
	// trampoline entry point, boot environment which is constructed
	// for UPL will be overwritten by trampoline code.
	stackOffset := trampStack & 0xFFF
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

// Get the base address and data from RDSP table
func getAcpiRsdpData() (uint64, []byte, error) {
	// Finds the RSDP in the EFI System Table.
	file, err := os.Open("/sys/firmware/efi/systab")
	if err != nil {
		return 0, nil, err
	}
	defer file.Close()

	const acpi20 = "ACPI20="

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		start := ""
		if strings.HasPrefix(line, acpi20) {
			start = strings.TrimPrefix(line, acpi20)
		}
		if start == "" {
			continue
		}
		base, err := strconv.ParseInt(start, 0, 63)
		if err != nil {
			continue
		}
		return uint64(base), nil, nil
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error while reading efi systab: %v", err)
	}

	return 0xFFFFFFFF, nil, ErrDTRsdpTableNotFound
}
