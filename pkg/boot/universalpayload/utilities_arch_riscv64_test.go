// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build riscv64 && !tinygo

package universalpayload

import (
	"encoding/binary"
	"testing"
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
