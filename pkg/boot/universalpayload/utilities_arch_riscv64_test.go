// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build riscv64

package universalpayload

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/u-root/u-root/pkg/acpi"
	"github.com/u-root/u-root/pkg/dt"
)

func TestConstructTrampolineRISCV(t *testing.T) {
	u := New()
	u.trampolineOffset = 0x12000
	u.fdtDtbOffset = 0x11000

	const load = uint64(0x80000000)
	const entry = uint64(0x80201000)
	buf := u.constructTrampoline(nil, load, entry)
	if buf == nil {
		t.Fatal("constructTrampoline() returned nil")
	}

	trampBegin := addrOfStart()
	trampStack := addrOfStackTop()
	trampHob := addrOfHobAddr()
	trampEntry := addrOfEntry()
	if trampHob < trampStack+8 {
		t.Fatalf("hob slot %#x is not 8 bytes after stack slot %#x", trampHob, trampStack)
	}
	if trampEntry-trampHob != trampHob-trampStack {
		t.Fatalf("slot stride mismatch: stack %#x hob %#x entry %#x", trampStack, trampHob, trampEntry)
	}
	padLen := int(trampHob - trampStack - 8)
	codeLen := int(trampStack - trampBegin)
	stackAt := codeLen
	fdtAt := stackAt + 8 + padLen
	entryAt := fdtAt + 8 + padLen

	if len(buf) != entryAt+8 {
		t.Fatalf("trampoline length %d, want %d", len(buf), entryAt+8)
	}
	if got := binary.LittleEndian.Uint64(buf[stackAt:]); got != load+u.trampolineOffset {
		t.Fatalf("stack slot %#x, want %#x", got, load+u.trampolineOffset)
	}
	if got := binary.LittleEndian.Uint64(buf[fdtAt:]); got != load+u.fdtDtbOffset {
		t.Fatalf("fdt slot %#x, want %#x", got, load+u.fdtDtbOffset)
	}
	if got := binary.LittleEndian.Uint64(buf[entryAt:]); got != entry {
		t.Fatalf("entry slot %#x, want %#x", got, entry)
	}

	// The stub addresses stack_top, hob_addr, and entry_point with
	// PC-relative AUIPC sequences. Those targets have to land on the
	// slots written above, at whatever address kexec places the copy.
	targets := pcRelTargets(t, buf[:codeLen])
	want := []int{stackAt, fdtAt, entryAt}
	if len(targets) != len(want) {
		t.Fatalf("pc-relative targets %v, want %v", targets, want)
	}
	for i := range want {
		if targets[i] != want[i] {
			t.Fatalf("pc-relative target %d = %d, want %d", i, targets[i], want[i])
		}
	}
}

// pcRelTargets returns the byte offsets addressed by AUIPC pairs in code.
func pcRelTargets(t *testing.T, code []byte) []int {
	t.Helper()
	var targets []int
	for pc := 0; pc+2 <= len(code); {
		if code[pc]&0x3 != 0x3 {
			pc += 2
			continue
		}
		if pc+4 > len(code) {
			break
		}
		word := binary.LittleEndian.Uint32(code[pc:])
		if word&0x7f != 0x17 {
			pc += 4
			continue
		}
		if pc+8 > len(code) {
			t.Fatalf("truncated AUIPC pair at offset %d", pc)
		}
		second := binary.LittleEndian.Uint32(code[pc+4:])
		hi := int(int32(word &^ 0xfff))
		lo := int(int32(second) >> 20)
		targets = append(targets, pc+hi+lo)
		pc += 8
	}
	return targets
}

func TestSerialPortFromFDT(t *testing.T) {
	fdt := &dt.FDT{
		RootNode: dt.NewNode("/", dt.WithChildren(
			dt.NewNode("serial@10000000", dt.WithProperty(
				dt.PropertyString("compatible", "ns16550a"),
				dt.PropertyU32Array("reg", []uint32{0, 0x10000000, 0, 0x100}),
				dt.PropertyU32("current-speed", 57600),
				dt.PropertyU32("reg-shift", 2),
			)),
			dt.NewNode("serial@10001000", dt.WithProperty(
				dt.PropertyString("compatible", "ns16550a"),
				dt.PropertyString("status", "disabled"),
				dt.PropertyU32Array("reg", []uint32{0, 0x10001000, 0, 0x100}),
			)),
		)),
	}

	useMmio, stride, baud, base, ok := serialPortFromFDT(fdt)
	if !ok {
		t.Fatal("serialPortFromFDT() ok = false")
	}
	if useMmio != 1 || stride != 4 || baud != 57600 || base != 0x10000000 {
		t.Fatalf("serial = mmio %d stride %d baud %d base %#x", useMmio, stride, baud, base)
	}
}

func TestPhysicalAddressSizeRISCV(t *testing.T) {
	u := New()
	t.Setenv("UROOT_PHYS_ADDR_SIZE", "")
	got, err := u.getPhysicalAddressSizes()
	if err != nil {
		t.Fatal(err)
	}
	if got != 56 {
		t.Fatalf("physical address size %d, want 56", got)
	}

	t.Setenv("UROOT_PHYS_ADDR_SIZE", "48")
	got, err = u.getPhysicalAddressSizes()
	if err != nil {
		t.Fatal(err)
	}
	if got != 48 {
		t.Fatalf("physical address size %d, want 48", got)
	}
}

func TestPhysicalAddressSizeRISCVInvalid(t *testing.T) {
	u := New()
	for _, v := range []string{"abc", "-1", "256"} {
		t.Setenv("UROOT_PHYS_ADDR_SIZE", v)
		if _, err := u.getPhysicalAddressSizes(); err == nil {
			t.Fatalf("UROOT_PHYS_ADDR_SIZE=%s: getPhysicalAddressSizes() = nil, want error", v)
		}
	}

	t.Setenv("UROOT_PHYS_ADDR_SIZE", "0")
	got, err := u.getPhysicalAddressSizes()
	if err != nil {
		t.Fatal(err)
	}
	if got != 0 {
		t.Fatalf("physical address size %d, want 0", got)
	}
}

func TestSerialDisabled(t *testing.T) {
	if serialDisabled(dt.NewNode("serial")) {
		t.Fatal("missing status was treated as disabled")
	}
	if serialDisabled(dt.NewNode("serial", dt.WithProperty(dt.PropertyString("status", "okay")))) {
		t.Fatal("status okay was treated as disabled")
	}
	nonString := dt.NewNode("serial", dt.WithProperty(dt.Property{Name: "status", Value: []byte{1}}))
	if serialDisabled(nonString) {
		t.Fatal("non-string status was treated as disabled")
	}
	disabled := dt.NewNode("serial", dt.WithProperty(dt.PropertyString("status", "disabled")))
	if !serialDisabled(disabled) {
		t.Fatal("status disabled was not treated as disabled")
	}
}

func TestSerialCompatible(t *testing.T) {
	if serialCompatible(dt.NewNode("serial")) {
		t.Fatal("missing compatible matched")
	}
	bad := dt.NewNode("serial", dt.WithProperty(dt.Property{Name: "compatible", Value: []byte{0xff}}))
	if serialCompatible(bad) {
		t.Fatal("non-string compatible matched")
	}
	unknown := dt.NewNode("serial", dt.WithProperty(dt.PropertyString("compatible", "uart")))
	if serialCompatible(unknown) {
		t.Fatal("unknown compatible matched")
	}
	for _, compatible := range []string{"ns16550a", "ns16550", "snps,dw-apb-uart", "sifive,uart0"} {
		node := dt.NewNode("serial", dt.WithProperty(dt.PropertyString("compatible", compatible)))
		if !serialCompatible(node) {
			t.Fatalf("compatible %q did not match", compatible)
		}
	}
}

func TestUARTBase(t *testing.T) {
	if _, ok := uartBase(dt.NewNode("serial")); ok {
		t.Fatal("missing reg")
	}
	short := dt.NewNode("serial", dt.WithProperty(dt.Property{Name: "reg", Value: []byte{1, 2}}))
	if _, ok := uartBase(short); ok {
		t.Fatal("short reg")
	}
	odd := dt.NewNode("serial", dt.WithProperty(dt.Property{Name: "reg", Value: []byte{1, 2, 3, 4, 5}}))
	if _, ok := uartBase(odd); ok {
		t.Fatal("reg length is not a multiple of 4")
	}

	two := dt.NewNode("serial", dt.WithProperty(dt.PropertyU32Array("reg", []uint32{0x10000000, 0x100})))
	base, ok := uartBase(two)
	if !ok || base != 0x10000000 {
		t.Fatalf("two-cell base %#x ok %v", base, ok)
	}
	zero := dt.NewNode("serial", dt.WithProperty(dt.PropertyU32Array("reg", []uint32{0, 0x100})))
	if _, ok := uartBase(zero); ok {
		t.Fatal("zero address")
	}

	four := dt.NewNode("serial", dt.WithProperty(dt.PropertyU32Array("reg", []uint32{0x1, 0, 0, 0x100})))
	base, ok = uartBase(four)
	if !ok || base != 0x100000000 {
		t.Fatalf("four-cell base %#x ok %v", base, ok)
	}
	one := dt.NewNode("serial", dt.WithProperty(dt.PropertyU32Array("reg", []uint32{0x3f8})))
	base, ok = uartBase(one)
	if !ok || base != 0x3f8 {
		t.Fatalf("one-cell base %#x ok %v", base, ok)
	}
}

func TestSerialStrideAndBaud(t *testing.T) {
	bare := dt.NewNode("serial")
	if stride := serialStride(bare); stride != UniversalPayloadSerialPortRegisterStride {
		t.Fatalf("missing reg-shift stride %d", stride)
	}
	if baud := serialBaud(bare); baud != UniversalPayloadSerialPortBaudRate {
		t.Fatalf("missing current-speed baud %d", baud)
	}

	bad := dt.NewNode("serial", dt.WithProperty(
		dt.PropertyString("reg-shift", "nope"),
		dt.PropertyString("current-speed", "nope"),
	))
	if stride := serialStride(bad); stride != UniversalPayloadSerialPortRegisterStride {
		t.Fatalf("bad reg-shift stride %d", stride)
	}
	if baud := serialBaud(bad); baud != UniversalPayloadSerialPortBaudRate {
		t.Fatalf("bad current-speed baud %d", baud)
	}

	wide := dt.NewNode("serial", dt.WithProperty(dt.PropertyU32("reg-shift", 8)))
	if stride := serialStride(wide); stride != UniversalPayloadSerialPortRegisterStride {
		t.Fatalf("reg-shift 8 stride %d", stride)
	}
	zeroBaud := dt.NewNode("serial", dt.WithProperty(dt.PropertyU32("current-speed", 0)))
	if baud := serialBaud(zeroBaud); baud != UniversalPayloadSerialPortBaudRate {
		t.Fatalf("zero baud %d", baud)
	}

	shifted := dt.NewNode("serial", dt.WithProperty(dt.PropertyU32("reg-shift", 0)))
	if stride := serialStride(shifted); stride != 1 {
		t.Fatalf("reg-shift 0 stride %d, want 1", stride)
	}
}

func TestSerialPortFromFDTRejects(t *testing.T) {
	if _, _, _, _, ok := serialPortFromFDT(nil); ok {
		t.Fatal("nil FDT")
	}
	if _, _, _, _, ok := serialPortFromFDT(&dt.FDT{}); ok {
		t.Fatal("FDT without a root node")
	}

	disabled := &dt.FDT{RootNode: dt.NewNode("/", dt.WithChildren(
		dt.NewNode("serial", dt.WithProperty(
			dt.PropertyString("compatible", "ns16550a"),
			dt.PropertyString("status", "disabled"),
			dt.PropertyU32Array("reg", []uint32{0x10000000, 0x100}),
		)),
	))}
	if _, _, _, _, ok := serialPortFromFDT(disabled); ok {
		t.Fatal("disabled UART was selected")
	}

	defaults := &dt.FDT{RootNode: dt.NewNode("/", dt.WithChildren(
		dt.NewNode("serial", dt.WithProperty(
			dt.PropertyString("compatible", "ns16550"),
			dt.PropertyU32Array("reg", []uint32{0x10000000, 0x100}),
		)),
	))}
	useMmio, stride, baud, base, ok := serialPortFromFDT(defaults)
	if !ok || useMmio != 1 || stride != UniversalPayloadSerialPortRegisterStride || baud != UniversalPayloadSerialPortBaudRate || base != 0x10000000 {
		t.Fatalf("defaults = mmio %d stride %d baud %d base %#x ok %v", useMmio, stride, baud, base, ok)
	}
}

func TestDiscoverPlatformSerialRISCV(t *testing.T) {
	orig := firmwareFDTPath
	t.Cleanup(func() { firmwareFDTPath = orig })

	firmwareFDTPath = filepath.Join(t.TempDir(), "missing")
	mmio, stride, baud, base, ok := discoverPlatformSerial()
	if !ok || mmio != 1 || stride != UniversalPayloadSerialPortRegisterStride || baud != UniversalPayloadSerialPortBaudRate || base != qemuVirtSerialBase {
		t.Fatalf("fallback = mmio %d stride %d baud %d base %#x ok %v", mmio, stride, baud, base, ok)
	}

	empty := &dt.FDT{RootNode: dt.NewNode("/")}
	firmwareFDTPath = writeFirmwareFDT(t, empty)
	_, _, _, base, ok = discoverPlatformSerial()
	if !ok || base != qemuVirtSerialBase {
		t.Fatalf("empty FDT base %#x ok %v, want qemu fallback", base, ok)
	}

	uart := &dt.FDT{RootNode: dt.NewNode("/", dt.WithChildren(
		dt.NewNode("serial@10000000", dt.WithProperty(
			dt.PropertyString("compatible", "ns16550a"),
			dt.PropertyU32Array("reg", []uint32{0, 0x10000000, 0, 0x100}),
			dt.PropertyU32("current-speed", 57600),
			dt.PropertyU32("reg-shift", 2),
		)),
	))}
	firmwareFDTPath = writeFirmwareFDT(t, uart)
	mmio, stride, baud, base, ok = discoverPlatformSerial()
	if !ok || mmio != 1 || stride != 4 || baud != 57600 || base != 0x10000000 {
		t.Fatalf("fdt serial = mmio %d stride %d baud %d base %#x ok %v", mmio, stride, baud, base, ok)
	}
}

func TestArchGetAcpiRsdpDataRISCV(t *testing.T) {
	u := New()
	u.getAcpiRsdp = nil
	base, data, err := u.archGetAcpiRsdpData()
	if err != nil || base != 0 || data != nil {
		t.Fatalf("nil getter = base %#x data %d err %v", base, len(data), err)
	}

	u.getAcpiRsdp = func() (*acpi.RSDP, error) {
		return nil, errors.New("no acpi")
	}
	base, data, err = u.archGetAcpiRsdpData()
	if err != nil || base != 0 || data != nil {
		t.Fatalf("getter error = base %#x data %d err %v", base, len(data), err)
	}

	u.getAcpiRsdp = func() (*acpi.RSDP, error) {
		return nil, nil
	}
	base, data, err = u.archGetAcpiRsdpData()
	if err != nil || base != 0 || data != nil {
		t.Fatalf("nil rsdp = base %#x data %d err %v", base, len(data), err)
	}

	u.getAcpiRsdp = func() (*acpi.RSDP, error) {
		return riscvTestRSDP(0x4000), nil
	}
	base, data, err = u.archGetAcpiRsdpData()
	if err != nil || base != 0 || len(data) == 0 {
		t.Fatalf("rsdp = base %#x data %d err %v", base, len(data), err)
	}
}

func TestRISCVMemMapHelpers(t *testing.T) {
	u := New()
	if n := u.appendAddonMemMap(nil); n != 0 {
		t.Fatalf("appendAddonMemMap() = %d, want 0", n)
	}
	if u.isMemReserved("Reserved") {
		t.Fatal("isMemReserved() = true, want false")
	}
}

func writeFirmwareFDT(t *testing.T, fdt *dt.FDT) string {
	t.Helper()
	fdt.Header.Magic = dt.Magic
	fdt.Header.Version = 17
	var buf bytes.Buffer
	if _, err := fdt.Write(&buf); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "fdt")
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// riscvTestRSDP builds a non-empty RSDP. The acpi package does not export a
// constructor that returns *RSDP.
func riscvTestRSDP(base int64) *acpi.RSDP {
	type mirror struct {
		base int64
		data [36]byte
	}
	m := &mirror{base: base}
	copy(m.data[:], "RSD PTR ")
	return (*acpi.RSDP)(unsafe.Pointer(m))
}
