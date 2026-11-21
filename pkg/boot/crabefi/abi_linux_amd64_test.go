// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !tinygo

package crabefi

import (
	"bytes"
	"debug/elf"
	"encoding/binary"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/boot/kexec"
)

// Same wire format as CrabEFI's .note.crabefi, independent of section headers.
const testABINote = "\x08\x00\x00\x00\x04\x00\x00\x00\x01\x00\x00\x00CrabEFI\x00\x01\x00\x00\x00"

func TestEntryABI(t *testing.T) {
	for name, change := range map[string]func([]byte){
		"version 1":               func([]byte) {},
		"default coreboot build":  func(b []byte) { clear(b[2048 : 2048+len(testABINote)]) },
		"unknown version":         func(b []byte) { b[2048+20] = 2 },
		"invalid descriptor size": func(b []byte) { b[2048+4] = 8 },
		"wrong owner":             func(b []byte) { b[2048+12] = 'X' },
		"unknown note type":       func(b []byte) { b[2048+8] = 2 },
		"oversized owner":         func(b []byte) { binary.LittleEndian.PutUint32(b[2048:], 0xffffffff) },
	} {
		t.Run(name, func(t *testing.T) {
			payload := testPayload(t, nil)
			change(payload)
			file, err := elf.NewFile(bytes.NewReader(payload))
			if err != nil {
				t.Fatal(err)
			}
			err = validateEntryABI(file)
			if got, want := err == nil, name == "version 1"; got != want {
				t.Errorf("validateEntryABI() = %v; want success %v", err, want)
			}
		})
	}
}

func TestRejectCorebootEntryABI(t *testing.T) {
	payload := testPayload(t, nil)
	// The default coreboot ELF64/.code32 payload has no entry64 note.
	binary.LittleEndian.PutUint16(payload[56:], 1)
	_, _, err := prepareLoad(bytes.NewReader(payload), testFirmware(t, 3))
	if err == nil || !strings.Contains(err.Error(), "--features entry64") {
		t.Errorf("prepareLoad(default coreboot entry ABI) = %v; want entry64 build hint", err)
	}
}

func TestEntryCode(t *testing.T) {
	got := entryCode(0x123456789, 0x23456789a)
	want := []byte{
		0xfa, 0xfc, // cli; cld
		0x48, 0xbf, // movabs rdi, handoff
		0x9a, 0x78, 0x56, 0x34, 0x02, 0, 0, 0,
		0x48, 0xb8, // movabs rax, entry
		0x89, 0x67, 0x45, 0x23, 0x01, 0, 0, 0,
		0xff, 0xe0, // jmp rax
	}
	if !bytes.Equal(got, want) {
		t.Errorf("entryCode() = %x; want %x", got, want)
	}
}

func TestHigh64BitAddresses(t *testing.T) {
	const payloadBase, arenaBase = 0x100100000, 0x200000000
	payload := testPayload(t, func(h *elf.Header64, p *elf.Prog64) {
		h.Entry, p.Paddr, p.Vaddr = payloadBase, payloadBase, payloadBase
	})
	f := testFirmware(t, 3)
	f.memory = append(f.memory[:1],
		kexec.TypedRange{Range: kexec.Range{Start: payloadBase, Size: 8192}, Type: kexec.RangeRAM},
		kexec.TypedRange{Range: kexec.Range{Start: arenaBase, Size: arenaSize}, Type: kexec.RangeRAM})
	entry, segments, err := prepareLoad(bytes.NewReader(payload), f)
	if err != nil || entry != arenaBase {
		t.Fatalf("prepareLoad(above 4 GiB) = %#x, %v; want %#x, nil", entry, err, uint64(arenaBase))
	}
	for _, segment := range segments {
		if segment.Phys.Start == arenaBase {
			if got := binary.LittleEndian.Uint64(segment.Buf[4:12]); got != arenaBase+handoffOffset {
				t.Errorf("handoff pointer = %#x; want %#x", got, uint64(arenaBase+handoffOffset))
			}
			if got := binary.LittleEndian.Uint64(segment.Buf[14:22]); got != payloadBase {
				t.Errorf("ELF entry = %#x; want %#x", got, uint64(payloadBase))
			}
		}
	}
}
