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

	"github.com/u-root/u-root/pkg/align"
)

// validateEntryABI requires CrabEFI's version-1 64-bit entry contract. The old
// 32-bit coreboot payload is also ELF64, so the ELF class alone is insufficient.
// Use program headers: the note must survive stripping section/symbol tables.
func validateEntryABI(file *elf.File) error {
	for _, p := range file.Progs {
		if p.Type != elf.PT_NOTE {
			continue
		}
		if p.Off > math.MaxInt64 || p.Filesz > math.MaxInt64-p.Off {
			return fmt.Errorf("ELF note file offset exceeds int64")
		}
		r := p.Open()
		for {
			var header struct{ NameSize, DescSize, Type uint32 }
			if err := binary.Read(r, binary.LittleEndian, &header); err != nil {
				if err == io.EOF {
					break
				}
				return fmt.Errorf("read CrabEFI ELF note: %w", err)
			}
			nameSize := align.Up(uint64(header.NameSize), 4)
			descSize := align.Up(uint64(header.DescSize), 4)
			if header.NameSize != 8 || header.Type != 1 {
				if _, err := io.CopyN(io.Discard, r, int64(nameSize+descSize)); err != nil {
					return fmt.Errorf("skip ELF note: %w", err)
				}
				continue
			}
			var name [8]byte
			if _, err := io.ReadFull(r, name[:]); err != nil {
				return fmt.Errorf("read ELF note owner: %w", err)
			}
			if string(name[:]) != "CrabEFI\x00" {
				if _, err := io.CopyN(io.Discard, r, int64(descSize)); err != nil {
					return fmt.Errorf("skip ELF note description: %w", err)
				}
				continue
			}
			if header.DescSize != 4 {
				return fmt.Errorf("invalid CrabEFI entry ABI note size %d", header.DescSize)
			}
			var version uint32
			if err := binary.Read(r, binary.LittleEndian, &version); err != nil {
				return fmt.Errorf("read CrabEFI entry ABI version: %w", err)
			}
			if version != 1 {
				return fmt.Errorf("unsupported CrabEFI entry ABI version %d (want 1)", version)
			}
			return nil
		}
	}
	return fmt.Errorf("missing CrabEFI 64-bit entry ABI note; rebuild the payload")
}
