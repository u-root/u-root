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
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/u-root/u-root/pkg/align"
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
func constructTrampoline(buf []uint8, addr uint64, entry uint64) []uint8 {
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

	buf = appendUint64(buf, addr+trampolineOffse)
	buf = appendUint64(buf, addr+fdtDtbOffset)
	buf = appendUint64(buf, entry)

	return buf
}

// According to Arm Server Base System Architecture 7.2 (DEN00291) chapter 1.2.7
// "Peripheral subsystems" for Level 3:
// " The base server system must implement a UART as specified by B_PER_O5 in
// Peripheral subsystems section from Arm BSA [4]. "
// Due to limitation of memoryMapFromIOMem, memory region of UART device cannot
// be parsed, we append memory region of UART device here.
func appendUARTMemMap(memMapHOB *EFIMemoryMapHOB) uint64 {
	f, err := os.Open("/proc/iomem")
	if err != nil {
		return 0
	}
	defer f.Close()

	b := bufio.NewScanner(f)
	for b.Scan() {
		var start uint64
		var end uint64
		content := b.Text()

		if strings.Contains(content, "ARMH0011") {
			els := strings.Split(content, ":")
			addrs := strings.Split(strings.TrimSpace(els[0]), "-")
			if len(addrs) != 2 {
				fmt.Printf("Address format incorrect for device 'ARMH0011'\n")
				continue
			}

			start, err = strconv.ParseUint(addrs[0], 16, 64)
			if err != nil {
				fmt.Printf("Failed to parse start address for device 'ARMH0011'\n")
				continue
			}

			end, err = strconv.ParseUint(addrs[1], 16, 64)
			if err != nil {
				fmt.Printf("Failed to parse end address for device 'ARMH0011'\n")
				continue
			}

			*memMapHOB = append(*memMapHOB, EFIHOBResourceDescriptor{
				Header: EFIHOBGenericHeader{
					HOBType:   EFIHOBTypeResourceDescriptor,
					HOBLength: EFIHOBLength(unsafe.Sizeof(EFIHOBResourceDescriptor{})),
				},
				ResourceType: EFIResourceMemoryMappedIO,
				ResourceAttribute: EFIResourceAttributePresent |
					EFIResourceAttributeInitialized |
					EFIResourceAttributeTested |
					EFIResourceAttributeUncacheable |
					EFIResourceAttributeWriteCombineable |
					EFIResourceAttributeWriteThroughCacheable |
					EFIResourceAttributeWriteBackCacheable,
				PhysicalStart:  EFIPhysicalAddress(start),
				ResourceLength: uint64(align.UpPage(end - start)),
			})

			return uint64(unsafe.Sizeof(EFIHOBResourceDescriptor{}))
		}
	}

	return 0
}

func appendAddonMemMap(memMapHOB *EFIMemoryMapHOB) uint64 {
	return appendUARTMemMap(memMapHOB)
}
