// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kexec

import (
	"bytes"
	"debug/elf"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/u-root/u-root/pkg/boot/align"
)

// PurgeLoad loads an ELF (which must be one of the purgatories from the
// pkg/boot/linux package) into kmem, instructing it to jump to entry with RSI
// set to rsi.
func PurgeLoad(kmem *Memory, elfBuf []byte, start, param uintptr) (uintptr, error) {
	elfFile, err := elf.NewFile(bytes.NewReader(elfBuf))
	if err != nil {
		return 0, fmt.Errorf("parse purgatory ELF file from ELF buffer: %v", err)
	}

	log.Printf("Elf file: %#v, %d Progs", elfFile, len(elfFile.Progs))
	if len(elfFile.Progs) != 1 {
		return 0, fmt.Errorf("parse purgatory ELF file: can only handle one Prog, not %d", len(elfFile.Progs))
	}
	p := elfFile.Progs[0]

	// the package really wants things page-sized, and rather than
	// deal with all the bugs that arise from that, just keep it happy.
	p.Memsz = uint64(align.AlignUpPageSize(uint(p.Memsz)))
	b := make([]byte, p.Memsz)
	if _, err := p.ReadAt(b[:p.Filesz], 0); err != nil {
		return 0, err
	}
	entry := uintptr(elfFile.Entry)
	Debug("Start is %#x, param is %#x", start, param)

	binary.LittleEndian.PutUint64(b[8:], uint64(start))
	binary.LittleEndian.PutUint64(b[16:], uint64(param))

	// TODO: Purgatories would be easier to handle if they were position-independent.
	min := uintptr(p.Vaddr)
	max := uintptr(p.Vaddr + uint64(len(b)))

	phyRange, err := kmem.ReservePhys(uint(len(b)), RangeFromInterval(min, max))
	if err != nil {
		return 0, fmt.Errorf("purgatory: reserve phys ram of size %d between range(%d, %d): %v", len(b), min, max, err)
	}
	kmem.Segments.Insert(NewSegment(b, phyRange))
	return entry, nil
}
