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

	"github.com/u-root/u-root/pkg/align"
)

func addrOfStart() uintptr
func addrOfStackTop() uintptr
func addrOfHobAddr() uintptr

func getPhysicalAddressSizes() (uint8, error) {
	// Return hardcode for arm64
	// Please update to actual physical address size
	physicalAddrSize := os.Getenv("UROOT_PHYS_ADDR_SIZE")
	if physicalAddrSize != "" {
		if num, err := strconv.ParseUint(physicalAddrSize, 10, 8); err == nil {
			return uint8(num), nil
		} else {
			return 0, fmt.Errorf("Malformed UROOT_PHYS_ADDR_SIZE value \"%s\": %v\n", physicalAddrSize, err)
		}
	}
	return 48, nil
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

	appendUint64 := func(slice []uint8, value uint64) []uint8 {
		tmpBytes := make([]uint8, 8)
		binary.LittleEndian.PutUint64(tmpBytes, value)
		return append(slice, tmpBytes...)
	}

	padLen := uint64(trampHob - trampStack - 8)

	// Update temporary stack top
	buf = appendUint64(buf, addr+trampolineOffset)
	buf = padWithLength(buf, padLen)
	// Update FDT DTB info address
	buf = appendUint64(buf, addr+fdtDtbOffset)
	buf = padWithLength(buf, padLen)
	buf = appendUint64(buf, entry)

	return buf
}

// Get the base address and data from RDSP table
func archGetAcpiRsdpData() (uint64, []byte, error) {
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
