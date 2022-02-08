package kexec

import (
	"bytes"
	"debug/elf"
	"fmt"
	"log"
)

func alignUpUint64(v, align uint64) uint64 {
	return (v + align) &^ align
}

// ELFLoad loads an elf file at the desired address.
func ELFLoad(kmem *Memory, elfBuf []byte, min uint, max uint, end int, flags uint32) (uintptr, error) {
	elfFile, err := elf.NewFile(bytes.NewReader(elfBuf))
	if err != nil {
		return 0, fmt.Errorf("parse elf file from elf buffer: %v", err)
	}

	log.Printf("Elf file: %#v, %d Progs", elfFile, len(elfFile.Progs))
	if len(elfFile.Progs) != 1 {
		return 0, fmt.Errorf("parse elf file: can only handle one Prog, not %d", len(elfFile.Progs))
	}
	p := elfFile.Progs[0]
	// the package really wants things page-sized, and rather than
	// deal with all the bugs that arise from that, just keep it happy.
	p.Memsz = uint64(alignUp(uint(p.Memsz)))
	b := make([]byte, p.Memsz)
	if _, err := p.ReadAt(b[:p.Filesz], 0); err != nil {
		return 0, err
	}
	entry := elfFile.Entry
	phyRange, err := kmem.ReservePhys(uint(len(b)), RangeFromInterval(uintptr(p.Vaddr), uintptr(p.Vaddr+uint64(len(b)))))
	if err != nil {
		return uintptr(entry), fmt.Errorf("reserve phys ram of size %d between range(%d, %d): %v", len(b), min, max, err)
	}
	kmem.Segments.Insert(NewSegment(b, phyRange))
	return uintptr(entry), nil
}
