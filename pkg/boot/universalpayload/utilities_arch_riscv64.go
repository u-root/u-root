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

	"github.com/u-root/u-root/pkg/dt"
)

// QEMU's virt machine places a 16550 at this MMIO address. Used when the
// running system's device tree does not name a UART.
const qemuVirtSerialBase = 0x10000000

// firmwareFDTPath is the firmware device tree. Tests point it at a fixture.
var firmwareFDTPath = "/sys/firmware/fdt"

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

// discoverPlatformSerial reads the UART from the firmware device tree.
// The generic handoff describes an ISA port at 0x3f8. RISC-V has no port
// I/O, and EDK2 polls that address until the transmitter looks empty, so
// the payload never leaves its serial loop.
func discoverPlatformSerial() (useMmio uint8, stride uint8, baud uint32, base uint64, ok bool) {
	if fdt, err := dt.ReadFile(firmwareFDTPath); err == nil {
		if mmio, stride, baud, base, ok := serialPortFromFDT(fdt); ok {
			return mmio, stride, baud, base, true
		}
	}
	return 1, UniversalPayloadSerialPortRegisterStride, UniversalPayloadSerialPortBaudRate, qemuVirtSerialBase, true
}

func serialPortFromFDT(fdt *dt.FDT) (useMmio uint8, stride uint8, baud uint32, base uint64, ok bool) {
	if fdt == nil || fdt.RootNode == nil {
		return 0, 0, 0, 0, false
	}
	node, found := fdt.RootNode.Find(func(n *dt.Node) bool {
		if serialDisabled(n) || !serialCompatible(n) {
			return false
		}
		_, ok := uartBase(n)
		return ok
	})
	if !found {
		return 0, 0, 0, 0, false
	}
	base, _ = uartBase(node)
	return 1, serialStride(node), serialBaud(node), base, true
}

func serialDisabled(n *dt.Node) bool {
	prop, ok := n.LookProperty("status")
	if !ok {
		return false
	}
	status, err := prop.AsString()
	return err == nil && status == "disabled"
}

func serialCompatible(n *dt.Node) bool {
	prop, ok := n.LookProperty("compatible")
	if !ok {
		return false
	}
	list, err := prop.AsStringList()
	if err != nil {
		return false
	}
	for _, compatible := range list {
		switch compatible {
		case "ns16550a", "ns16550", "snps,dw-apb-uart", "sifive,uart0":
			return true
		}
	}
	return false
}

func uartBase(n *dt.Node) (uint64, bool) {
	prop, ok := n.LookProperty("reg")
	if !ok || len(prop.Value) < 4 || len(prop.Value)%4 != 0 {
		return 0, false
	}
	cells := make([]uint32, len(prop.Value)/4)
	for i := range cells {
		cells[i] = binary.BigEndian.Uint32(prop.Value[i*4 : (i+1)*4])
	}
	switch len(cells) {
	case 4:
		// #address-cells = 2, #size-cells = 2.
		return uint64(cells[0])<<32 | uint64(cells[1]), true
	case 2:
		// One address cell and one size cell.
		return uint64(cells[0]), cells[0] != 0
	default:
		return uint64(cells[0]), cells[0] != 0
	}
}

func serialStride(n *dt.Node) uint8 {
	prop, ok := n.LookProperty("reg-shift")
	if !ok {
		return UniversalPayloadSerialPortRegisterStride
	}
	shift, err := prop.AsU32()
	if err != nil || shift >= 8 {
		return UniversalPayloadSerialPortRegisterStride
	}
	return 1 << shift
}

func serialBaud(n *dt.Node) uint32 {
	prop, ok := n.LookProperty("current-speed")
	if !ok {
		return UniversalPayloadSerialPortBaudRate
	}
	baud, err := prop.AsU32()
	if err != nil || baud == 0 {
		return UniversalPayloadSerialPortBaudRate
	}
	return baud
}

func (u *UPL) appendAddonMemMap(_ *EFIMemoryMapHOB) uint64 {
	return 0
}

func (u *UPL) isMemReserved(memType string) bool {
	return false
}
