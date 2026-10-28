// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !tinygo

package crabefi

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/orangecms/go-framebuffer/framebuffer"
	"github.com/u-root/u-root/pkg/acpi"
	"github.com/u-root/u-root/pkg/boot/kexec"
	"golang.org/x/sys/unix"
)

type rsdpV1 struct {
	Signature [8]byte
	Checksum  uint8
	OEMID     [6]byte
	Revision  uint8
	RSDT      uint32
}
type rsdpTable struct {
	Legacy           rsdpV1
	Length           uint32
	XSDT             uint64
	ExtendedChecksum uint8
	Reserved         [3]byte
}

func preserved(memory kexec.MemoryMap, address uint64, size uint64) bool {
	if size == 0 || address > ^uint64(0)-size {
		return false
	}
	for _, r := range memory {
		if r.Type != kexec.RangeRAM && uint64(r.Start) <= address && address+size <= uint64(r.End()) {
			return true
		}
	}
	return false
}

func validateRSDP(data []byte) error {
	var legacy rsdpV1
	if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &legacy); err != nil {
		return err
	}
	if legacy.Signature != [8]byte{'R', 'S', 'D', ' ', 'P', 'T', 'R', ' '} || byteSum(data[:binary.Size(legacy)]) != 0 {
		return fmt.Errorf("invalid RSDP signature or checksum")
	}
	if legacy.Revision >= 2 {
		var table rsdpTable
		if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &table); err != nil {
			return err
		}
		if int(table.Length) != len(data) || byteSum(data) != 0 {
			return fmt.Errorf("invalid extended RSDP")
		}
	}
	return nil
}

func rsdpFromLog(log string) (uint64, error) {
	_, rest, ok := strings.Cut(log, "ACPI: RSDP 0x")
	fields := strings.Fields(rest)
	if !ok || len(fields) == 0 {
		return 0, fmt.Errorf("RSDP missing from kernel log")
	}
	address, err := strconv.ParseUint(fields[0], 16, 64)
	if err != nil {
		return 0, err
	}
	if address == 0 {
		return 0, fmt.Errorf("null RSDP")
	}
	return address, nil
}

func readRSDP(memory kexec.MemoryMap) (uint64, []byte, error) {
	// Reuse the ACPI package's EFI and legacy discovery methods first.
	if rsdp, err := acpi.GetRSDP(); err == nil {
		return uint64(rsdp.RSDPAddr()), rsdp.AllData(), nil
	}
	// Coreboot may place the RSDP outside the legacy scan area without EFI
	// sysfs. The kernel log is a fallback, and must still contain its ACPI
	// discovery message. READ_ALL does not consume or clear that log.
	size, err := unix.Klogctl(unix.SYSLOG_ACTION_SIZE_BUFFER, nil)
	if err != nil {
		return 0, nil, fmt.Errorf("kernel log unavailable: %w", err)
	}
	if size <= 0 {
		return 0, nil, fmt.Errorf("empty kernel log")
	}
	log := make([]byte, size)
	n, err := unix.Klogctl(unix.SYSLOG_ACTION_READ_ALL, log)
	if err != nil {
		return 0, nil, err
	}
	address, err := rsdpFromLog(string(log[:n]))
	if err != nil {
		return 0, nil, err
	}
	mem, err := os.Open("/dev/mem")
	if err != nil {
		return 0, nil, err
	}
	defer mem.Close()
	data, err := readRSDPAt(mem, memory, address)
	return address, data, err
}

// readRSDPAt checks the physical extent before converting to ReaderAt's signed
// offset. In particular, a kernel-log address must not wrap into a negative
// /dev/mem offset, even if a firmware memory range happens to contain it.
func readRSDPAt(mem io.ReaderAt, memory kexec.MemoryMap, address uint64) ([]byte, error) {
	if address > math.MaxInt64 {
		return nil, fmt.Errorf("RSDP address %#x exceeds int64 offset", address)
	}
	size := binary.Size(rsdpTable{})
	if address > math.MaxInt64-uint64(size) {
		return nil, fmt.Errorf("RSDP address %#x exceeds int64 read extent", address)
	}
	if !preserved(memory, address, uint64(size)) {
		return nil, fmt.Errorf("RSDP is not in preserved firmware memory")
	}
	data := make([]byte, size)
	if _, err := mem.ReadAt(data, int64(address)); err != nil {
		return nil, fmt.Errorf("read RSDP at %#x: %w", address, err)
	}
	return data, nil
}

func framebufferRecord(fb *framebuffer.Framebuffer) (framebufferInfo, error) {
	fixed, v := fb.Finfo, fb.Vinfo
	if fixed.Type != 0 || fixed.Visual != 2 || v.Bits_per_pixel != 32 ||
		v.Xoffset != 0 || v.Yoffset != 0 || v.Grayscale != 0 || v.Nonstd != 0 ||
		v.Xres == 0 || v.Yres == 0 || fixed.Smem_start == 0 || fixed.Line_length%4 != 0 ||
		// CrabEFI's x86 entry identity-maps only the first 64 GiB.
		fixed.Smem_start >= 64<<30 || uint64(fixed.Smem_len) > (64<<30)-fixed.Smem_start ||
		uint64(v.Xres)*4 > uint64(fixed.Line_length) ||
		uint64(fixed.Line_length)*uint64(v.Yres) > uint64(fixed.Smem_len) {
		return framebufferInfo{}, fmt.Errorf("requires a physical 32-bit linear true-color framebuffer")
	}
	var channels [4]colorChannel
	var used uint32
	for i, c := range []struct{ offset, length, right uint32 }{
		{v.Red.Offset, v.Red.Length, v.Red.Right}, {v.Green.Offset, v.Green.Length, v.Green.Right},
		{v.Blue.Offset, v.Blue.Length, v.Blue.Right}, {v.Transp.Offset, v.Transp.Length, v.Transp.Right},
	} {
		if c.right != 0 || c.length > 8 || c.offset > 32-c.length || i < 3 && c.length != 8 {
			return framebufferInfo{}, fmt.Errorf("unsupported framebuffer channel")
		}
		mask := uint32((uint64(1)<<c.length)-1) << c.offset
		if used&mask != 0 {
			return framebufferInfo{}, fmt.Errorf("overlapping framebuffer channels")
		}
		used |= mask
		channels[i] = colorChannel{uint8(c.offset), uint8(c.length)}
	}
	return framebufferInfo{
		Address: fixed.Smem_start, Width: v.Xres, Height: v.Yres,
		Stride: fixed.Line_length, BitsPerPixel: 32,
		Red: channels[0], Green: channels[1], Blue: channels[2], Reserved: channels[3],
	}, nil
}

func readFirmware() (firmwareInfo, error) {
	var f firmwareInfo
	var err error
	f.memory, err = kexec.MemoryMapFromSysfsMemmap()
	if err != nil {
		return f, fmt.Errorf("firmware memory map: %w", err)
	}
	var data []byte
	f.rsdp, data, err = readRSDP(f.memory)
	if err != nil {
		return f, err
	}
	if !preserved(f.memory, f.rsdp, uint64(len(data))) {
		return f, fmt.Errorf("RSDP is not preserved")
	}
	if err := validateRSDP(data); err != nil {
		return f, err
	}
	fb, err := framebuffer.Init("/dev/fb0")
	if err != nil {
		return f, fmt.Errorf("open framebuffer: %w", err)
	}
	defer fb.Close()
	f.framebuffer, err = framebufferRecord(fb)
	if err != nil {
		return f, err
	}
	for _, r := range f.memory.RAM() {
		if r.Overlaps(kexec.Range{Start: uintptr(fb.Finfo.Smem_start), Size: uint(fb.Finfo.Smem_len)}) {
			return f, fmt.Errorf("framebuffer overlaps reusable RAM")
		}
	}
	f.smbiosEntry, err = os.ReadFile("/sys/firmware/dmi/tables/smbios_entry_point")
	if err != nil {
		return f, fmt.Errorf("read SMBIOS entry: %w", err)
	}
	f.smbiosData, err = os.ReadFile("/sys/firmware/dmi/tables/DMI")
	if err != nil {
		return f, fmt.Errorf("read SMBIOS data: %w", err)
	}
	return f, nil
}
