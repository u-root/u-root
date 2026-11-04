// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build riscv64 && !tinygo

package universalpayload

import (
	"encoding/binary"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"unsafe"
)

func addrOfStart() uintptr
func addrOfStackTop() uintptr
func addrOfHobAddr() uintptr
func addrOfEntry() uintptr

func (u *UPL) getPhysicalAddressSizes() (uint8, error) {
	// QEMU's virt machine implements a 56-bit physical address space.
	// UROOT_PHYS_ADDR_SIZE overrides that for other platforms.
	physicalAddrSize := os.Getenv("UROOT_PHYS_ADDR_SIZE")
	if physicalAddrSize == "" {
		return 56, nil
	}
	num, err := strconv.ParseUint(physicalAddrSize, 10, 8)
	if err != nil {
		return 0, fmt.Errorf("malformed UROOT_PHYS_ADDR_SIZE value %q: %w", physicalAddrSize, err)
	}
	return uint8(num), nil
}

// constructTrampoline builds the RISC-V stub in trampoline_riscv64.s.
// addr is the base of the load area. entry is the FIT entry point.
func (u *UPL) constructTrampoline(buf []uint8, addr uint64, entry uint64) []uint8 {
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
	trampEntry := addrOfEntry()

	if trampHob < trampStack+8 || trampEntry < trampHob+8 {
		return nil
	}

	buf = append(buf, ptrToSlice(trampBegin, int(trampStack-trampBegin))...)

	appendUint64 := func(slice []uint8, value uint64) []uint8 {
		tmpBytes := make([]uint8, 8)
		binary.LittleEndian.PutUint64(tmpBytes, value)
		return append(slice, tmpBytes...)
	}

	padWithLength := func(slice []uint8, length uint64) []uint8 {
		tmpBytes := make([]uint8, length)
		return append(slice, tmpBytes...)
	}

	// Temporary stack top, handoff FDT, then FIT entry point.
	// Padding follows the assembly symbol distances, so each AUIPC
	// still lands on its slot after the stub is copied.
	buf = appendUint64(buf, addr+u.trampolineOffset)
	buf = padWithLength(buf, uint64(trampHob-trampStack-8))
	buf = appendUint64(buf, addr+u.fdtDtbOffset)
	buf = padWithLength(buf, uint64(trampEntry-trampHob-8))
	buf = appendUint64(buf, entry)

	return buf
}

// archGetAcpiRsdpData returns RSDP bytes to copy into the load area.
// RISC-V firmware often boots with ACPI disabled. Missing ACPI is not a
// load failure; the caller omits the ACPI handoff when both results are zero.
func (u *UPL) archGetAcpiRsdpData() (uint64, []byte, error) {
	if u.getAcpiRsdp != nil {
		if rsdp, err := u.getAcpiRsdp(); err == nil && rsdp != nil && len(rsdp.AllData()) > 0 {
			return 0, rsdp.AllData(), nil
		}
	}
	return 0, nil, nil
}

func (u *UPL) appendAddonMemMap(_ *EFIMemoryMapHOB) uint64 {
	return 0
}

func (u *UPL) isMemReserved(memType string) bool {
	return false
}
