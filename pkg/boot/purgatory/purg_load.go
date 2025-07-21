// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run genpurg.go purgatories.go

// Package purgatory provides several purgatories for use with kexec_load
// system call.
//
// The kernel's contract on x86_64 is that kexec will jump to the entry point
// in 64bit mode with identity-mapped page tables and unspecified garbage in
// registers.
//
// To drop into the right addressing mode and pass arguments to the next
// kernel, some additional code is placed before the execution of the new
// kernel called a purgatory. So it will go: old kernel jumps to purgatory
// jumps to new kernel.
//
// Our purgatory will only have these three responsibilities:
// - drop to the intended addressing mode (32bit or 64bit)
// - set the right rsi parameters (for Linux kernels)
// - jump to the entry point.
//
// Several purgatories written in Assembly are part of this package: one that
// sets up the Linux args and remains in 64bit mode, and one that sets up the
// Linux args and drops to 32bit mode.
package purgatory

import (
	"bytes"
	"debug/elf"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/u-root/u-root/pkg/align"
	"github.com/u-root/u-root/pkg/boot/kexec"
)

const defaultPurgatory = "default"

var (
	// Debug is called to print out verbose debug info.
	//
	// Set this to appropriate output stream for display
	// of useful debug info.
	Debug        = log.Printf // func(string, ...interface{}) {}
	curPurgatory = Purgatories[defaultPurgatory]
)

// Select picks a purgatory, returning an error if none is found.
func Select(name string) error {
	p, ok := Purgatories[name]
	if !ok {
		var s []string
		for i := range Purgatories {
			s = append(s, i)
		}
		return fmt.Errorf("%s: no such purgatory, try one of %v", name, s)

	}
	curPurgatory = p
	return nil
}

// Load loads the selected purgatory into kmem, instructing it to jump to entry
// with RSI set to rsi.
func Load(kmem *kexec.Memory, entry, rsi uintptr) (uintptr, error) {
	elfFile, err := elf.NewFile(bytes.NewReader(curPurgatory.Code))
	if err != nil {
		return 0, fmt.Errorf("parse purgatory ELF file from ELF buffer: %w", err)
	}

	log.Printf("Elf file: %#v, %d Progs", elfFile, len(elfFile.Progs))
	if len(elfFile.Progs) != 1 {
		return 0, fmt.Errorf("parse purgatory ELF file: can only handle one Prog, not %d", len(elfFile.Progs))
	}
	p := elfFile.Progs[0]

	// the package really wants things page-sized, and rather than
	// deal with all the bugs that arise from that, just keep it happy.
	p.Memsz = uint64(align.UpPage(uint(p.Memsz)))
	b := make([]byte, p.Memsz)
	if _, err := p.ReadAt(b[:p.Filesz], 0); err != nil {
		return 0, err
	}
	elfEntry := uintptr(elfFile.Entry)

	// Debug("Start is %#x, param is %#x", start, param)
	binary.LittleEndian.PutUint64(b[8:], uint64(entry))
	binary.LittleEndian.PutUint64(b[16:], uint64(rsi))

	// TODO: Shouldn't the purgatories be relocatable?
	minimum := uintptr(p.Vaddr)
	maximum := uintptr(p.Vaddr + uint64(len(b)))

	phyRange, err := kmem.ReservePhys(uint(len(b)), kexec.RangeFromInterval(minimum, maximum))
	if err != nil {
		return 0, fmt.Errorf("purgatory: reserve phys ram of size %d between range(%d, %d): %w", len(b), minimum, maximum, err)
	}
	kmem.Segments.Insert(kexec.NewSegment(b, phyRange))
	return elfEntry, nil
}
