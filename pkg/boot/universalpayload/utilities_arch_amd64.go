// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package universalpayload

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"unsafe"
)

func addrOfStart() uintptr
func addrOfStackTop() uintptr
func addrOfHobAddr() uintptr

// Get Physical Address size from sysfs node /proc/cpuinfo.
// Both Physical and Virtual Address size will be prompted as format:
// "address sizes	: 39 bits physical, 48 bits virtual"
// Use regular expression to fetch the integer of Physical Address
// size before "bits physical" keyword
func getPhysicalAddressSizes() (uint8, error) {
	file, err := os.Open(sysfsCPUInfoPath)
	if err != nil {
		return 0, fmt.Errorf("failed to open %s: %w", sysfsCPUInfoPath, err)
	}
	defer file.Close()

	// Regular expression to match the address size line
	re := regexp.MustCompile(`address sizes\s*:\s*(\d+)\s+bits physical,\s*(\d+)\s+bits virtual`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if match := re.FindStringSubmatch(line); match != nil {
			// Convert the physical bits size to integer
			physicalBits, err := strconv.ParseUint(match[1], 10, 8)
			if err != nil {
				return 0, errors.Join(ErrCPUAddressConvert, err)
			}
			return uint8(physicalBits), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("%w: file: %s, err: %w", ErrCPUAddressRead, sysfsCPUInfoPath, err)
	}

	return 0, ErrCPUAddressNotFound
}

// Construct trampoline code before jump to entry point of FIT image.
// Due to lack of support to set value of General Purpose Registers in kexec,
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

	padLen := uint64(trampHob - trampStack - 8)

	tramp := ptrToSlice(trampBegin, int(trampStack-trampBegin))

	buf = append(buf, tramp...)

	stackTop := hobAddr + tmpStackTop
	appendUint64 := func(slice []uint8, value uint64) []uint8 {
		tmpBytes := make([]uint8, 8)
		binary.LittleEndian.PutUint64(tmpBytes, value)
		return append(slice, tmpBytes...)
	}

	padWithLength := func(slice []uint8, len uint64) []uint8 {
		tmpBytes := make([]uint8, len)
		return append(slice, tmpBytes...)
	}

	buf = appendUint64(buf, stackTop)
	buf = padWithLength(buf, padLen)
	buf = appendUint64(buf, hobAddr)
	buf = padWithLength(buf, padLen)
	buf = appendUint64(buf, entry)

	return buf
}
