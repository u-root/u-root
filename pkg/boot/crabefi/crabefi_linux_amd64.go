// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !tinygo

// Package crabefi loads the x86-64 CrabEFI coreboot payload from Linux.
package crabefi

import (
	"debug/elf"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"sort"

	"github.com/u-root/u-root/pkg/boot/internal/trampoline"
	"github.com/u-root/u-root/pkg/boot/kexec"
)

type firmwareInfo struct {
	memory                  kexec.MemoryMap
	rsdp                    uint64
	framebuffer             framebufferInfo
	smbiosEntry, smbiosData []byte
}

// PrepareLoad stages the ELF and firmware handoff for kexec_load. It does not
// call kexec_load or reboot; callers retain the usual boot menu cleanup flow.
func PrepareLoad(payload io.ReaderAt) (uintptr, kexec.Segments, error) {
	firmware, err := readFirmware()
	if err != nil {
		return 0, nil, err
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
			start = (start + pageSize - 1) &^ (pageSize - 1)
			end &^= pageSize - 1
		} else {
			start &^= pageSize - 1
			end = (end + pageSize - 1) &^ (pageSize - 1)
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
	file, err := elf.NewFile(payload)
	if err != nil {
		return 0, nil, err
	}
	if file.Class != elf.ELFCLASS64 || file.Data != elf.ELFDATA2LSB || file.Type != elf.ET_EXEC || file.Machine != elf.EM_X86_64 {
		return 0, nil, fmt.Errorf("expected executable little-endian x86-64 coreboot ELF")
	}
	mem := kexec.Memory{Phys: memory}
	entryValid := false
	var loaded kexec.Ranges
	for _, p := range file.Progs {
		if p.Type != elf.PT_LOAD {
			continue
		}
		if p.Memsz == 0 || p.Paddr >= 1<<32 || p.Memsz > (1<<32)-p.Paddr || p.Filesz > p.Memsz || p.Vaddr != p.Paddr || p.Paddr%pageSize != 0 {
			return 0, nil, fmt.Errorf("unsupported ELF load segment")
		}
		r := kexec.Range{Start: uintptr(p.Paddr), Size: uint(p.Memsz)}
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
	if !entryValid || len(loaded) > 12 {
		return 0, nil, fmt.Errorf("invalid ELF entry or too many load segments")
	}
	if _, err := mem.LoadElfSegments(payload); err != nil {
		return 0, nil, err
	}
	var arena kexec.Range
	for base := uintptr(0x08000000); base <= 0xffe00000; base += arenaSize {
		r := kexec.Range{Start: base, Size: arenaSize}
		if !inRAM(memory, r) {
			continue
		}
		free := true
		for _, segment := range mem.Segments {
			if segment.Phys.Overlaps(r) {
				free = false
				break
			}
		}
		if free {
			arena = r
			break
		}
	}
	if arena.Size == 0 {
		return 0, nil, fmt.Errorf("no free handoff arena below 4 GiB")
	}
	mem.Phys.Insert(kexec.TypedRange{Range: arena, Type: kexec.RangeReserved})
	if len(mem.Phys) > 64 {
		return 0, nil, fmt.Errorf("too many memory regions for CrabEFI")
	}
	handoffBase := uint64(arena.Start + handoffOffset)
	tables, err := makeHandoff(mem.Phys, handoffBase, firmware.rsdp, firmware.framebuffer, firmware.smbiosEntry, firmware.smbiosData)
	if err != nil {
		return 0, nil, err
	}
	// The shared trampoline enters 32-bit mode with EBX holding the handoff.
	// Establish coreboot's cdecl stack argument before entering the payload.
	adapter := []byte{
		0xfa, 0xfc, // cli; cld
		0x89, 0xdc, // mov esp, ebx
		0x53,       // push ebx
		0x6a, 0x00, // push 0 (return address)
		0xb8, // mov eax, imm32
	}
	adapter = binary.LittleEndian.AppendUint32(adapter, uint32(file.Entry))
	adapter = append(adapter, 0xff, 0xe0) // jmp eax
	stage := make([]byte, handoffOffset+len(tables))
	copy(stage, adapter)
	copy(stage[handoffOffset:], tables)
	mem.Segments.Insert(kexec.NewSegment(stage, arena))
	code, err := trampoline.Setup("", 0, uintptr(handoffBase), arena.Start)
	if err != nil {
		return 0, nil, err
	}
	jump, err := mem.AddKexecSegment(code)
	if err != nil {
		return 0, nil, err
	}
	return jump.Start, mem.Segments, nil
}
