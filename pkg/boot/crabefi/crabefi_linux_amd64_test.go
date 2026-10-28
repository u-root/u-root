// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !tinygo

package crabefi

import (
	"bytes"
	"debug/elf"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/orangecms/go-framebuffer/framebuffer"
	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/smbios"
)

func testPayload(t *testing.T, change func(*elf.Header64, *elf.Prog64)) []byte {
	t.Helper()
	h := elf.Header64{Type: uint16(elf.ET_EXEC), Machine: uint16(elf.EM_X86_64), Version: 1, Entry: 0x100000, Phoff: 64, Ehsize: 64, Phentsize: 56, Phnum: 1}
	copy(h.Ident[:], []byte{0x7f, 'E', 'L', 'F', 2, 1, 1})
	p := elf.Prog64{Type: uint32(elf.PT_LOAD), Flags: uint32(elf.PF_R | elf.PF_X), Off: 4096, Vaddr: 0x100000, Paddr: 0x100000, Filesz: 4, Memsz: 8192, Align: 4096}
	if change != nil {
		change(&h, &p)
	}
	return testPayloadWithProgs(t, h, []elf.Prog64{p})
}

func testPayloadWithProgs(t *testing.T, h elf.Header64, progs []elf.Prog64) []byte {
	t.Helper()
	h.Phnum = uint16(len(progs))
	var out bytes.Buffer
	if err := binary.Write(&out, binary.LittleEndian, h); err != nil {
		t.Fatal(err)
	}
	if err := binary.Write(&out, binary.LittleEndian, progs); err != nil {
		t.Fatal(err)
	}
	data := make([]byte, 4100)
	copy(data, out.Bytes())
	copy(data[4096:], []byte{0xfa, 0xf4, 0xeb, 0xfd})
	return data
}

func testFirmware(t *testing.T, version int) firmwareInfo {
	t.Helper()
	data := []byte{127, 4, 0, 0, 0, 0}
	var entry []byte
	var err error
	if version == 3 {
		ep := smbios.Entry64{Length: 24, SMBIOSMajorVersion: 3, StructMaxSize: uint32(len(data))}
		copy(ep.Anchor[:], "_SM3_")
		entry, err = ep.MarshalBinary()
	} else {
		ep := smbios.Entry32{Length: 31, SMBIOSMajorVersion: 2, StructTableLength: uint16(len(data))}
		copy(ep.Anchor[:], "_SM_")
		copy(ep.IntAnchor[:], "_DMI_")
		entry, err = ep.MarshalBinary()
	}
	if err != nil {
		t.Fatal(err)
	}
	return firmwareInfo{
		memory: kexec.MemoryMap{
			{Range: kexec.Range{Start: 0xf0000, Size: 4096}, Type: kexec.RangeReserved},
			{Range: kexec.Range{Start: 0x100000, Size: 0x20000000}, Type: kexec.RangeRAM},
		},
		rsdp: 0xf0000, smbiosEntry: entry, smbiosData: data,
		framebuffer: framebufferInfo{Address: 0xe0000000, Width: 1024, Height: 768, Stride: 4096, BitsPerPixel: 32},
	}
}

func TestPrepareLoad(t *testing.T) {
	for _, version := range []int{2, 3} {
		f := testFirmware(t, version)
		payload := testPayload(t, nil)
		entry, segments, err := prepareLoad(bytes.NewReader(payload), f)
		if err != nil {
			t.Fatalf("prepareLoad(SMBIOS %d) = %v", version, err)
		}
		if entry != arenaMin+trampolineOffset || len(segments) != 2 {
			t.Fatalf("prepareLoad() = entry %#x, %d segments; want trampoline in arena and two segments", entry, len(segments))
		}
		var stage []byte
		for _, s := range segments {
			if s.Phys.Start == 0x100000 && (s.Phys.Size != 8192 || !bytes.Equal(s.Buf, payload[4096:])) {
				t.Errorf("ELF segment = %v; want original data and 8192 bytes including BSS", s)
			}
			if s.Phys.Start == 0x8000000 {
				if s.Phys.Size != arenaSize {
					t.Fatalf("arena size = %#x; want %#x", s.Phys.Size, arenaSize)
				}
				stage = s.Buf[handoffOffset:]
				if binary.LittleEndian.Uint32(s.Buf[8:]) != 0x100000 {
					t.Errorf("adapter target is not the ELF entry")
				}
			}
		}
		if stage == nil {
			t.Fatal("handoff arena was not loaded")
		}
		var header tableHeader
		if err := binary.Read(bytes.NewReader(stage), binary.LittleEndian, &header); err != nil {
			t.Fatal(err)
		}
		if string(header.Signature[:]) != "LBIO" || header.HeaderBytes != 24 || header.TableEntries != 4 || checksum(stage[:24]) != 0 || checksum(stage[24:24+header.TableBytes]) != uint16(header.TableChecksum) {
			t.Fatalf("invalid coreboot header: %+v", header)
		}
		r := bytes.NewReader(stage[24 : 24+header.TableBytes])
		reserved := false
		for r.Len() != 0 {
			var record recordHeader
			if err := binary.Read(r, binary.LittleEndian, &record); err != nil {
				t.Fatal(err)
			}
			if record.Size < 8 || record.Size%4 != 0 || uint64(record.Size-8) > uint64(r.Len()) {
				t.Fatalf("record size = %d; want an aligned record within the table", record.Size)
			}
			body := make([]byte, record.Size-8)
			if _, err := r.Read(body); err != nil {
				t.Fatal(err)
			}
			if record.Tag == 1 {
				m := bytes.NewReader(body)
				for m.Len() > 0 {
					var region memoryRange
					if err := binary.Read(m, binary.LittleEndian, &region); err != nil {
						t.Fatal(err)
					}
					if region.Base == 0x8000000 && region.Size == arenaSize && region.Kind == 2 {
						reserved = true
					}
				}
			}
		}
		if !reserved {
			t.Error("handoff arena is not reserved in firmware memory map")
		}
		e32, e64, err := smbios.ParseEntry(stage[smbiosEntryOffset:smbiosDataOffset])
		if err != nil {
			t.Fatal(err)
		}
		address := uint64(0)
		if e64 != nil {
			address = e64.StructTableAddr
		} else {
			address = uint64(e32.StructTableAddr)
		}
		if address != 0x8004000+smbiosDataOffset || !bytes.Equal(stage[smbiosDataOffset:], f.smbiosData) {
			t.Errorf("relocated SMBIOS address/data = %#x/%x", address, stage[smbiosDataOffset:])
		}
		if _, err := kexec.AlignAndMerge(segments); err != nil {
			t.Fatalf("kexec segments are invalid: %v", err)
		}
	}
}

func TestRejectInvalidPayload(t *testing.T) {
	for name, change := range map[string]func(*elf.Header64, *elf.Prog64){
		"wrong architecture":  func(h *elf.Header64, p *elf.Prog64) { h.Machine = uint16(elf.EM_AARCH64) },
		"PIE":                 func(h *elf.Header64, p *elf.Prog64) { h.Type = uint16(elf.ET_DYN) },
		"entry in BSS":        func(h *elf.Header64, p *elf.Prog64) { h.Entry += 4 },
		"file exceeds memory": func(h *elf.Header64, p *elf.Prog64) { p.Filesz = 8193 },
		"crosses 4GiB":        func(h *elf.Header64, p *elf.Prog64) { p.Paddr = 0xfffff000; p.Vaddr = p.Paddr },
		"not identity mapped": func(h *elf.Header64, p *elf.Prog64) { p.Vaddr += 4096 },
		"empty segment":       func(h *elf.Header64, p *elf.Prog64) { p.Memsz = 0 },
		"above 4GiB":          func(h *elf.Header64, p *elf.Prog64) { p.Paddr = 1 << 32; p.Vaddr = p.Paddr },
		"unaligned segment":   func(h *elf.Header64, p *elf.Prog64) { p.Paddr++; p.Vaddr = p.Paddr },
		"non-executable":      func(h *elf.Header64, p *elf.Prog64) { p.Flags = uint32(elf.PF_R) },
		"no load segments":    func(h *elf.Header64, p *elf.Prog64) { p.Type = uint32(elf.PT_NOTE) },
		"not in RAM":          func(h *elf.Header64, p *elf.Prog64) { p.Paddr = 0xf0000; p.Vaddr = p.Paddr },
		"offset above int64":  func(h *elf.Header64, p *elf.Prog64) { p.Off = 1 << 63 },
		"read crosses int64":  func(h *elf.Header64, p *elf.Prog64) { p.Off = math.MaxInt64 - 1 },
	} {
		t.Run(name, func(t *testing.T) {
			if _, _, err := prepareLoad(bytes.NewReader(testPayload(t, change)), testFirmware(t, 3)); err == nil {
				t.Fatal("prepareLoad(invalid ELF) succeeded")
			}
		})
	}
	for _, n := range []int{0, 63, 119, 4099} {
		if _, _, err := prepareLoad(bytes.NewReader(testPayload(t, nil)[:n]), testFirmware(t, 3)); err == nil {
			t.Errorf("prepareLoad(truncated at %d) succeeded", n)
		}
	}
}

func TestNormalizeMemoryMap(t *testing.T) {
	for name, memory := range map[string]kexec.MemoryMap{
		"empty":             nil,
		"zero size":         {{Range: kexec.Range{Start: 4096}, Type: kexec.RangeRAM}},
		"wraparound":        {{Range: kexec.Range{Start: math.MaxUint64 - 4095, Size: 8192}, Type: kexec.RangeRAM}},
		"rounding overflow": {{Range: kexec.Range{Start: math.MaxUint64 - 4095, Size: 4094}, Type: kexec.RangeReserved}},
		"no full RAM page":  {{Range: kexec.Range{Start: 1, Size: 4094}, Type: kexec.RangeRAM}},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := normalizeMemoryMap(memory); err == nil {
				t.Fatalf("normalizeMemoryMap(%v) = nil; want error", memory)
			}
		})
	}
	memory := kexec.MemoryMap{
		{Range: kexec.Range{Start: 0xf001, Size: 30}, Type: kexec.RangeReserved},
		{Range: kexec.Range{Start: 0x1001, Size: 0x3ffe}, Type: kexec.RangeRAM},
	}
	want := kexec.MemoryMap{
		{Range: kexec.Range{Start: 0x2000, Size: 0x2000}, Type: kexec.RangeRAM},
		{Range: kexec.Range{Start: 0xf000, Size: 0x1000}, Type: kexec.RangeReserved},
	}
	got, err := normalizeMemoryMap(memory)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("normalizeMemoryMap() mismatch (-want +got):\n%s", diff)
	}
	if memory[0].Start != 0xf001 {
		t.Error("normalizeMemoryMap() modified its input")
	}
}

func TestHandoffAllocation(t *testing.T) {
	for _, base := range []uintptr{arenaMin, 0xfff00000} {
		t.Run(fmt.Sprintf("%#x", base), func(t *testing.T) {
			f := testFirmware(t, 3)
			f.memory = append(f.memory[:1],
				kexec.TypedRange{Range: kexec.Range{Start: 0x100000, Size: 8192}, Type: kexec.RangeRAM},
				kexec.TypedRange{Range: kexec.Range{Start: base, Size: arenaSize}, Type: kexec.RangeRAM})
			entry, segments, err := prepareLoad(bytes.NewReader(testPayload(t, nil)), f)
			if err != nil || entry != base+trampolineOffset {
				t.Fatalf("prepareLoad() = %#x, %v; want %#x, nil", entry, err, base+trampolineOffset)
			}
			if _, err := kexec.AlignAndMerge(segments); err != nil {
				t.Fatal(err)
			}
		})
	}
	f := testFirmware(t, 3)
	f.memory[1].Size = 8192
	if _, _, err := prepareLoad(bytes.NewReader(testPayload(t, nil)), f); !errors.Is(err, kexec.ErrNotEnoughSpace) {
		t.Errorf("prepareLoad(no arena) = %v; want ErrNotEnoughSpace", err)
	}

	f = testFirmware(t, 3)
	payload := testPayload(t, func(h *elf.Header64, p *elf.Prog64) {
		h.Entry, p.Paddr, p.Vaddr = arenaMin, arenaMin, arenaMin
	})
	entry, _, err := prepareLoad(bytes.NewReader(payload), f)
	if err != nil || entry != arenaMin+arenaSize+trampolineOffset {
		t.Errorf("prepareLoad(occupied arena) = %#x, %v; want next arena", entry, err)
	}

	for i := range maxMemoryRegions {
		f.memory = append(f.memory, kexec.TypedRange{
			Range: kexec.Range{Start: 0x30000000 + uintptr(i)*8192, Size: 4096}, Type: kexec.RangeReserved,
		})
	}
	if _, _, err := prepareLoad(bytes.NewReader(testPayload(t, nil)), f); err == nil {
		t.Error("prepareLoad(too many memory regions) = nil; want error")
	}
}

func TestELFSegmentValidation(t *testing.T) {
	var h elf.Header64
	var p elf.Prog64
	r := bytes.NewReader(testPayload(t, nil))
	if err := binary.Read(r, binary.LittleEndian, &h); err != nil {
		t.Fatal(err)
	}
	if err := binary.Read(r, binary.LittleEndian, &p); err != nil {
		t.Fatal(err)
	}
	if _, _, err := prepareLoad(bytes.NewReader(testPayloadWithProgs(t, h, []elf.Prog64{p, p})), testFirmware(t, 3)); err == nil {
		t.Error("prepareLoad(overlapping segments) = nil; want error")
	}
	var progs []elf.Prog64
	for i := range 16 {
		q := p
		q.Paddr += uint64(i) * 8192
		q.Vaddr = q.Paddr
		progs = append(progs, q)
	}
	if _, _, err := prepareLoad(bytes.NewReader(testPayloadWithProgs(t, h, progs)), testFirmware(t, 3)); err == nil {
		t.Error("prepareLoad(too many segments) = nil; want error")
	}
	f := testFirmware(t, 3)
	f.memory[1].Size = 8191
	if _, _, err := prepareLoad(bytes.NewReader(testPayload(t, nil)), f); err == nil {
		t.Error("prepareLoad(segment crosses RAM boundary) = nil; want error")
	}
}

func TestPrepareLoadNil(t *testing.T) {
	if _, _, err := PrepareLoad(nil); err == nil {
		t.Error("PrepareLoad(nil) = nil; want error")
	}
}

func TestMemoryAndSMBIOSValidation(t *testing.T) {
	memory := testFirmware(t, 3).memory
	memory = append(memory, kexec.TypedRange{Range: kexec.Range{Start: 0x1fffff, Size: 4096}, Type: kexec.RangeACPI})
	if _, err := normalizeMemoryMap(memory); err == nil {
		t.Error("normalizeMemoryMap(overlapping pages) succeeded")
	}
	if _, err := normalizeMemoryMap(nil); err == nil {
		t.Error("normalizeMemoryMap(empty map) succeeded")
	}
	for _, version := range []int{2, 3} {
		f := testFirmware(t, version)
		f.smbiosEntry[7] ^= 1
		if _, err := relocateSMBIOS(f.smbiosEntry, f.smbiosData, 0x8000000); err == nil {
			t.Errorf("relocateSMBIOS(corrupt version %d) succeeded", version)
		}
		f = testFirmware(t, version)
		for _, entry := range [][]byte{nil, f.smbiosEntry[:len(f.smbiosEntry)-1], append(append([]byte(nil), f.smbiosEntry...), 0)} {
			if _, err := relocateSMBIOS(entry, f.smbiosData, 0x8000000); err == nil {
				t.Errorf("relocateSMBIOS(version %d, entry length %d) = nil; want error", version, len(entry))
			}
		}
		if _, err := makeHandoff(f.memory, 0x8004000, f.rsdp, f.framebuffer, f.smbiosEntry, make([]byte, arenaSize)); err == nil {
			t.Error("makeHandoff(oversized SMBIOS) succeeded")
		}
	}
}

func TestRSDPValidation(t *testing.T) {
	data := make([]byte, 36)
	copy(data, "RSD PTR ")
	data[15] = 2
	binary.LittleEndian.PutUint32(data[20:], 36)
	data[8] = -byteSum(data[:20])
	data[32] = -byteSum(data)
	if err := validateRSDP(data); err != nil {
		t.Fatal(err)
	}
	data[32] ^= 1
	if err := validateRSDP(data); err == nil {
		t.Error("validateRSDP(corrupt extended checksum) succeeded")
	}
	address, err := rsdpFromLog("ACPI: RSDP 0x000000007FE80040 000024 (v02 CORE)")
	if err != nil || address != 0x7fe80040 {
		t.Fatalf("rsdpFromLog = %#x, %v", address, err)
	}
	for _, log := range []string{"", "ACPI: RSDP 0x", "ACPI: RSDP 0x0", "ACPI: RSDP 0xBADADDRESS"} {
		if _, err := rsdpFromLog(log); err == nil {
			t.Errorf("rsdpFromLog(%q) succeeded", log)
		}
	}
	if preserved(testFirmware(t, 3).memory, 0x100000, 36) {
		t.Error("RSDP in reusable RAM was accepted")
	}
}

func TestFramebufferValidation(t *testing.T) {
	fb := &framebuffer.Framebuffer{}
	fb.Finfo.Smem_start = 0xe0000000
	fb.Finfo.Smem_len = 4096 * 768
	fb.Finfo.Visual = 2
	fb.Finfo.Line_length = 4096
	fb.Vinfo.Xres, fb.Vinfo.Yres, fb.Vinfo.Bits_per_pixel = 1024, 768, 32
	fb.Vinfo.Red.Offset, fb.Vinfo.Red.Length = 16, 8
	fb.Vinfo.Green.Offset, fb.Vinfo.Green.Length = 8, 8
	fb.Vinfo.Blue.Length = 8
	fb.Vinfo.Transp.Offset, fb.Vinfo.Transp.Length = 24, 8
	r, err := framebufferRecord(fb)
	if err != nil || r.Address != 0xe0000000 || binary.Size(r) != 29 {
		t.Fatalf("framebufferRecord = %+v, %v", r, err)
	}
	body, err := binary.Append(nil, binary.LittleEndian, r)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(body[20:], []byte{32, 16, 8, 8, 8, 0, 8, 24, 8}) {
		t.Errorf("framebuffer pixel format = %x; want coreboot's 32-bit BGR format", body[20:])
	}
	fb.Vinfo.Green.Offset = 16
	if _, err := framebufferRecord(fb); err == nil {
		t.Error("framebufferRecord(overlapping channels) succeeded")
	}
	fb.Vinfo.Green.Offset = 8
	fb.Finfo.Smem_start = 0
	if _, err := framebufferRecord(fb); err == nil {
		t.Error("framebufferRecord(shadow buffer) succeeded")
	}
}
