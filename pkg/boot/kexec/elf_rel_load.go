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

// RelocateAndLoad does ...
func RelocateAndLoad(kmem *Memory, elfBuf []byte, min uint, max uint, end int, flags uint32) (entry uintptr, err error) {
	elfFile, err := elf.NewFile(bytes.NewReader(elfBuf))
	if err != nil {
		return uintptr(0), fmt.Errorf("parse elf file from elf buffer: %v", err)
	}

	if len(elfFile.Sections) == 0 {
		return uintptr(0), fmt.Errorf("empty sections")
	}

	if elfFile.Type != elf.ET_REL {
		return entry, fmt.Errorf("the elf is not a relocatable")
	}

	/* Find which section the entry is in */
	var entrySection *elf.Section
	var entryVal uint64
	for _, section := range elfFile.Sections {
		if (section.Flags&elf.SHF_ALLOC) == 0 || (section.Flags&elf.SHF_EXECINSTR) == 0 {
			// Does not occupy mem or contains no instructions.
			continue
		}
		if section.Addr <= elfFile.Entry && (section.Addr+section.Size) > elfFile.Entry {
			entrySection = section
			/* Make entry section relative */
			entryVal -= section.Addr
		}
	}
	log.Printf("entry section identified as: %v", entrySection)

	/* Find memory footprint of relocatable objects  */
	var bufAlign, bssAlign, bufsz, bsssz uint64 = 1, 1, 0, 0
	for _, section := range elfFile.Sections {
		if (section.Flags & elf.SHF_ALLOC) == 0 {
			continue
		}
		if section.Type != elf.SHT_NOBITS {
			align := section.Addralign
			if bufAlign < align {
				bufAlign = align
			}
			bufsz = alignUpUint64(bufsz, align)
			bufsz += section.Size
		} else { // bss: block start symbol sections.
			align := section.Addralign
			if bssAlign < align {
				bssAlign = align
			}
			bsssz = alignUpUint64(bsssz, align)
			bsssz += section.Size
		}
	}
	if bufAlign < bssAlign {
		bufAlign = bssAlign
	}
	var bssPad uint64 = 0
	if (bufsz & (bsssz - 1)) != 0 {
		bssPad = bssAlign - (bufsz & (bssAlign - 1))
	}
	log.Printf("bssPad: %d", bssPad)

	/* Allocate for relocatable objects. */
	buf := make([]byte, bufsz)
	memsz := bufsz + bssPad + bsssz // Allocate additional ram for bss.
	phyRange, err := kmem.ReservePhys(
		uint(memsz),
		RangeFromInterval(
			uintptr(min),
			uintptr(max),
		),
	)
	if err != nil {
		return entry, fmt.Errorf("reserve phys ram of size %d between range(%d, %d): %v", memsz, min, max, err)
	}
	kmem.Segments.Insert(NewSegment(buf, phyRange))

	log.Printf("Added segment for relocatable objects at: %v", phyRange)

	/* Update addresses for SHF_ALLOC sections. */
	dataAddr := uint64(phyRange.Start)
	bssAddr := dataAddr + bufsz + bssPad
	for _, section := range elfFile.Sections {
		if (section.Flags & elf.SHF_ALLOC) == 0 {
			continue
		}
		align := section.Addralign
		if section.Type != elf.SHT_NOBITS {
			dataAddr = alignUpUint64(dataAddr, align)
			off := dataAddr - uint64(phyRange.Start)
			secData, err := section.Data()
			if err != nil {
				return entry, fmt.Errorf("read data from section %v: %v", section, err)
			}
			copy(buf[off:], secData)
			section.Addr = dataAddr
			// Advance to addr for next section.
			dataAddr += section.Size
		} else {
			// TODO(10000TB): update elf section hdrs once we can edit it.
			bssAddr = alignUpUint64(bssAddr, align)
			section.Addr = bssAddr
			// Advance to addr for ext section.
			bssAddr += section.Size
		}
	}

	if entrySection != nil {
		entry += uintptr(entrySection.Addr)
		elfFile.Entry = uint64(entry)
	}

	entry = uintptr(elfFile.Entry)

	return uintptr(entry), nil
}

// ElfRelFindSymbol finds and return symbol by name from the given ELF file.
func ElfRelFindSymbol(e *elf.File, name string) (*elf.Symbol, error) {
	symbols, err := e.Symbols()
	if err != nil {
		return nil, fmt.Errorf("retrieve symbol tables from elf: %v", err)
	}
	for _, sym := range symbols {
		if elf.ST_BIND(sym.Info) != elf.STB_GLOBAL {
			continue
		}
		if sym.Name != name {
			continue
		}
		/* found the named symbol */
		if sym.Section == elf.SHN_UNDEF {
			return nil, fmt.Errorf("symbol: %v has bad section index: %v", sym, sym.Section)
		}
		return &sym, nil
	}
	return nil, fmt.Errorf("did not find a symbol named after: %s", name)
}

func ElfRelSetSymbol(e *elf.File, name string) error {
	//var shdr *MemSymbolHdr
	//var memSym MemSymbol

	sym, err := ElfRelFindSymbol(e, name)
	if err != nil {
		return fmt.Errorf("ElfRelFindSymbol(%v, %s): %v", e, name, err)
	}
	log.Printf("found symbol %s: %v", name, *sym)

	return nil
}
