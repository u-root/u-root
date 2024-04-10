// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package zimage contains a Parser for the arm zImage Linux format. It assumes
// little endian arm.
package zimage

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Magic values used in the zImage header and table.
const (
	Magic      = 0x016f2818
	Endianness = 0x04030201
	TableMagic = 0x45454545
)

// Tags used by TableEntry (at the time of writing, there is only one tag).
const (
	TagKernelSize Tag = 0x5a534c4b
)

// Tag is used to identify a TableEntry.
type Tag uint32

// ZImage is one of the major formats used by Linux on ARM. This struct
// is only for storing the metadata.
type ZImage struct {
	Header Header
	Table  []TableEntry
}

// Header appears near the beginning of the zImage.
//
// The layout is defined in Linux:
//
//	arch/arm/boot/compressed/head.S
type Header struct {
	Magic      uint32
	Start      uint32
	End        uint32
	Endianness uint32
	TableMagic uint32
	TableAddr  uint32
}

// TableEntry is an extension to Header. A zImage may have 0 or more entries.
//
// The layout is defined in Linux:
//
//	arch/arm/boot/compressed/vmlinux.lds.S
type TableEntry struct {
	Tag  Tag
	Data []uint32
}

// Parse a ZImage from a file.
func Parse(f io.ReadSeeker) (*ZImage, error) {
	// Parse the header.
	if _, err := f.Seek(0x24, io.SeekStart); err != nil {
		return nil, err
	}
	z := &ZImage{}
	if err := binary.Read(f, binary.LittleEndian, &z.Header); err != nil {
		return nil, err
	}
	if z.Header.Magic != Magic {
		return z, fmt.Errorf("invalid zImage magic, got %#08x, expected %#08x",
			z.Header.Magic, Magic)
	}
	if z.Header.Endianness != Endianness {
		return z, fmt.Errorf("unsupported zImage endianness, expected little")
	}
	if z.Header.End < z.Header.Start {
		return z, fmt.Errorf("invalid zImage, end is less than start, %d < %d",
			z.Header.End, z.Header.Start)
	}

	if z.Header.TableMagic != TableMagic {
		// No table.
		return z, nil
	}

	// Parse the table.
	addr := z.Header.TableAddr
	for addr != 0 {
		if _, err := f.Seek(int64(addr), io.SeekStart); err != nil {
			return nil, err
		}
		var size uint32
		if err := binary.Read(f, binary.LittleEndian, &size); err != nil {
			return nil, err
		}
		entry := TableEntry{Data: make([]uint32, size)}
		if err := binary.Read(f, binary.LittleEndian, &entry.Tag); err != nil {
			return nil, err
		}
		if err := binary.Read(f, binary.LittleEndian, &entry.Data); err != nil {
			return nil, err
		}
		z.Table = append(z.Table, entry)

		// In its current form, the Linux source code does not make it
		// super clear how multiple entries are specified in the table. Is
		// it a zero-terminated array? Is it a linked-list? Is it similar
		// to atags? See Linux commit c77256. Regardless, the kernel
		// currently only has one entry, so we exit after one iteration.
		addr = 0
	}
	return z, nil
}

// GetEntry searches through the zImage table for the given tag.
func (z *ZImage) GetEntry(t Tag) (*TableEntry, error) {
	for i := range z.Table {
		if z.Table[i].Tag == t {
			return &z.Table[i], nil
		}
	}
	return nil, fmt.Errorf("zImage table does not contain the %#08x tag", t)
}

// GetKernelSizes returns two kernel sizes relevant for kexec.
func (z *ZImage) GetKernelSizes() (piggySizeAddr uint32, kernelBSSSize uint32, err error) {
	e, err := z.GetEntry(TagKernelSize)
	if err != nil {
		return 0, 0, err
	}
	if len(e.Data) != 2 {
		return 0, 0, fmt.Errorf("zImage tag %#08x has incorrect size %d, expected 2",
			TagKernelSize, len(e.Data))
	}
	return e.Data[0], e.Data[1], nil
}
