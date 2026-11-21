// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !tinygo

package crabefi

import (
	"debug/elf"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"slices"
	"sort"

	"github.com/u-root/u-root/pkg/align"
	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/smbios"
	"golang.org/x/sys/unix"
)

type firmwareInfo struct {
	memory                  kexec.MemoryMap
	rsdp                    uint64
	framebuffer             *framebufferInfo // nil if Linux exposes none
	smbiosEntry, smbiosData []byte
}

// Supported reports whether PrepareLoad works on this platform.
const Supported = true

// PrepareLoad stages the ELF and firmware handoff for kexec_load. It does not
// call kexec_load or reboot; callers retain the usual boot menu cleanup flow.
func PrepareLoad(payload io.ReaderAt) (uintptr, kexec.Segments, error) {
	if payload == nil {
		return 0, nil, fmt.Errorf("nil CrabEFI payload")
	}
	// Report an unsuitable payload before probing firmware interfaces.
	if _, err := parsePayload(payload); err != nil {
		return 0, nil, err
	}
	firmware, err := readFirmware()
	if err != nil {
		return 0, nil, fmt.Errorf("read firmware: %w", err)
	}
	return prepareLoad(payload, firmware)
}

func normalizeMemoryMap(memory kexec.MemoryMap) (kexec.MemoryMap, error) {
	var result kexec.MemoryMap
	for _, r := range memory {
		start, end := r.Start, r.End()
		if r.Size == 0 || end <= start || end > math.MaxUint64-(pageSize-1) {
			return nil, fmt.Errorf("invalid firmware memory range %s", r)
		}
		if r.Type == kexec.RangeRAM {
			start = align.Up(start, uintptr(pageSize))
			end = align.Down(end, uintptr(pageSize))
		} else {
			start = align.Down(start, uintptr(pageSize))
			end = align.Up(end, uintptr(pageSize))
		}
		if start >= end {
			continue
		}
		result = append(result, kexec.TypedRange{Range: kexec.RangeFromInterval(start, end), Type: r.Type})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Start < result[j].Start })
	if len(result) == 0 {
		return nil, fmt.Errorf("empty firmware memory map")
	}
	for i := 1; i < len(result); i++ {
		if result[i-1].End() > result[i].Start {
			return nil, fmt.Errorf("overlapping firmware pages")
		}
	}
	return result, nil
}

// parsePayload checks the ELF header and CrabEFI's entry ABI note.
func parsePayload(payload io.ReaderAt) (*elf.File, error) {
	file, err := elf.NewFile(payload)
	if err != nil {
		return nil, err
	}
	if file.Class != elf.ELFCLASS64 || file.Data != elf.ELFDATA2LSB || file.Type != elf.ET_EXEC || file.Machine != elf.EM_X86_64 {
		return nil, fmt.Errorf("expected executable little-endian x86-64 coreboot ELF")
	}
	if err := validateEntryABI(file); err != nil {
		return nil, err
	}
	return file, nil
}

func inRAM(memory kexec.MemoryMap, r kexec.Range) bool {
	for _, ram := range memory.RAM() {
		if ram.IsSupersetOf(r) {
			return true
		}
	}
	return false
}

func prepareLoad(payload io.ReaderAt, firmware firmwareInfo) (uintptr, kexec.Segments, error) {
	memory, err := normalizeMemoryMap(firmware.memory)
	if err != nil {
		return 0, nil, err
	}
	file, err := parsePayload(payload)
	if err != nil {
		return 0, nil, err
	}
	mem := kexec.Memory{Phys: memory}
	entryValid := false
	var loaded kexec.Ranges
	for _, p := range file.Progs {
		if p.Type != elf.PT_LOAD {
			continue
		}
		// CrabEFI replaces the caller's identity map with its own 64-GiB map.
		// kexec rounds segment extents to pages, so validate those pages too.
		if p.Memsz == 0 || p.Paddr >= mappedMemorySize || p.Memsz > mappedMemorySize-p.Paddr || p.Filesz > p.Memsz || p.Vaddr != p.Paddr || !align.IsAligned(p.Paddr, uint64(pageSize)) {
			return 0, nil, fmt.Errorf("unsupported ELF load segment")
		}
		if p.Off > math.MaxInt64 || p.Filesz > math.MaxInt64-p.Off {
			return 0, nil, fmt.Errorf("ELF segment file offset exceeds int64")
		}
		r := (kexec.Range{Start: uintptr(p.Paddr), Size: uint(p.Memsz)}).Align(pageSize)
		if !inRAM(memory, r) {
			return 0, nil, fmt.Errorf("ELF segment %s is not in firmware RAM", r)
		}
		for _, prev := range loaded {
			if prev.Overlaps(r) {
				return 0, nil, fmt.Errorf("overlapping ELF segments")
			}
		}
		loaded = append(loaded, r)
		if p.Flags&elf.PF_X != 0 && p.Paddr <= file.Entry && file.Entry < p.Paddr+p.Filesz {
			entryValid = true
		}
	}
	if !entryValid {
		return 0, nil, fmt.Errorf("ELF entry is not in executable file data")
	}
	// One additional kexec segment holds the entry code and firmware tables.
	if len(loaded) >= unix.KEXEC_SEGMENT_MAX {
		return 0, nil, fmt.Errorf("too many ELF load segments")
	}
	if _, err := mem.LoadElfSegments(payload); err != nil {
		return 0, nil, err
	}
	// The SMBIOS 2 entry point holds a 32-bit table address.
	arenaMax := uintptr(mappedMemorySize)
	if _, e64, err := smbios.ParseEntry(firmware.smbiosEntry); err != nil {
		return 0, nil, err
	} else if e64 == nil {
		arenaMax = 4 << 30
	}
	free := mem.AvailableRAM()
	arena, err := free.FindSpace(arenaSize,
		kexec.WithinRange(kexec.RangeFromInterval(arenaMin, arenaMax)),
		kexec.WithAlignment(arenaSize))
	if err != nil {
		return 0, nil, fmt.Errorf("allocate CrabEFI handoff below %#x: %w", arenaMax, err)
	}
	if !slices.ContainsFunc(free, func(r kexec.Range) bool { return r.IsSupersetOf(arena) }) {
		return 0, nil, fmt.Errorf("CrabEFI handoff %s is not in free RAM", arena)
	}
	mem.Phys.Insert(kexec.TypedRange{Range: arena, Type: kexec.RangeReserved})
	if len(mem.Phys) > maxMemoryRegions {
		return 0, nil, fmt.Errorf("too many memory regions for CrabEFI (maximum %d)", maxMemoryRegions)
	}
	handoffBase := uint64(arena.Start + handoffOffset)
	tables, err := makeHandoff(mem.Phys, handoffBase, firmware.rsdp, firmware.framebuffer, firmware.smbiosEntry, firmware.smbiosData)
	if err != nil {
		return 0, nil, err
	}
	// Linux's x86-64 kexec enters with paging enabled and identity-mapped
	// segments. Pass the full-width table pointer in RDI without changing
	// CPU mode, GDT or stack; CrabEFI owns its startup state from this point.
	stage := make([]byte, handoffOffset+len(tables))
	copy(stage, entryCode(file.Entry, handoffBase))
	copy(stage[handoffOffset:], tables)
	mem.Segments.Insert(kexec.NewSegment(stage, arena))
	return arena.Start, mem.Segments, nil
}

func entryCode(entry, handoff uint64) []byte {
	code := []byte{0xfa, 0xfc, 0x48, 0xbf} // cli; cld; movabs rdi, imm64
	code = binary.LittleEndian.AppendUint64(code, handoff)
	code = append(code, 0x48, 0xb8) // movabs rax, imm64
	code = binary.LittleEndian.AppendUint64(code, entry)
	return append(code, 0xff, 0xe0) // jmp rax
}
