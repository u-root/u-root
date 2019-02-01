// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package dt contains utilities for device tree.
package dt

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"unsafe"

	"github.com/u-root/u-root/pkg/uio"
)

const (
	// Magic value seen in the FDT Header.
	Magic = 0xd00dfeed

	// MaxTotalSize is a limitation imposed by this implementation. This
	// prevents the integers from wrapping around. Typically, the total size is
	// a few megabytes, so this is not restrictive.
	MaxTotalSize = 1024 * 1024 * 1024
)

type token uint32

const (
	tokenBeginNode token = 0x1
	tokenEndNode         = 0x2
	tokenProp            = 0x3
	tokenNop             = 0x4
	tokenEnd             = 0x9
)

// FDT contains the parsed contents of a Flattend Device Tree (.dtb).
//
// The format is relatively simple and defined in chapter 5 of the Devicetree
// Specification Release 0.2.
//
// See: https://github.com/devicetree-org/devicetree-specification/releases/tag/v0.2
//
// This package is compatible with version 16 and 17 of DTSpec.
type FDT struct {
	Header         Header
	ReserveEntries []ReserveEntry
	RootNode       *Node
}

// Header appears at offset 0.
type Header struct {
	Magic           uint32
	TotalSize       uint32
	OffDtStruct     uint32
	OffDtStrings    uint32
	OffMemRsvmap    uint32
	Version         uint32
	LastCompVersion uint32
	BootCpuidPhys   uint32
	SizeDtStrings   uint32
	SizeDtStruct    uint32
}

// ReserveEntry defines a memory region which is reserved.
type ReserveEntry struct {
	Address uint64
	Size    uint64
}

// ReadFDT reads FDT from an io.ReaderSeeker.
func ReadFDT(f io.ReadSeeker) (*FDT, error) {
	fdt := &FDT{}
	if err := fdt.readHeader(f); err != nil {
		return nil, err
	}
	if err := fdt.readMemoryReservationBlock(f); err != nil {
		return nil, err
	}
	if err := fdt.checkLayout(); err != nil {
		return nil, err
	}
	strs, err := fdt.readStringsBlock(f)
	if err != nil {
		return nil, err
	}
	if err := fdt.readStructBlock(f, strs); err != nil {
		return nil, err
	}
	return fdt, nil
}

func (fdt *FDT) readHeader(f io.ReadSeeker) error {
	h := &fdt.Header
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return err
	}
	if err := binary.Read(f, binary.BigEndian, h); err != nil {
		return err
	}
	if h.Magic != Magic {
		return fmt.Errorf("invalid FDT magic, got %#08x, expected %#08x",
			h.Magic, Magic)
	}
	if !(h.Version == 16 || h.Version == 17 ||
		(h.LastCompVersion <= 17 && h.Version > 17)) {
		return fmt.Errorf(
			"incompatible FDT version, must be compatible with 16/17,"+
				"but version is %d last compatible with %d",
			h.Version, h.LastCompVersion)
	}
	if h.TotalSize > MaxTotalSize {
		return fmt.Errorf("FDT too large, %d > %d", h.TotalSize, MaxTotalSize)
	}
	return nil
}

func (fdt *FDT) readMemoryReservationBlock(f io.ReadSeeker) error {
	if fdt.Header.OffMemRsvmap < uint32(unsafe.Sizeof(fdt.Header)) {
		return fmt.Errorf("memory reservation block may not overlap header, %#x < %#x",
			fdt.Header.OffMemRsvmap, unsafe.Sizeof(fdt.Header))
	}
	if fdt.Header.OffMemRsvmap%8 != 0 {
		return fmt.Errorf(
			"memory reservation offset must be aligned to 8 bytes, but is %#x",
			fdt.Header.OffMemRsvmap)
	}
	if _, err := f.Seek(int64(fdt.Header.OffMemRsvmap), io.SeekStart); err != nil {
		return err
	}
	fdt.ReserveEntries = []ReserveEntry{}
	for {
		entry := ReserveEntry{}
		if err := binary.Read(f, binary.BigEndian, &entry); err != nil {
			return err
		}
		if entry.Address == 0 && entry.Size == 0 {
			break
		}
		fdt.ReserveEntries = append(fdt.ReserveEntries, entry)
	}
	return nil
}

// checkLayout returns any errors if any of the blocks overlap, are in the
// wrong order or stray past the end of the file. This function must be called
// after readHeader and readMemoryReservationBlock.
func (fdt *FDT) checkLayout() error {
	memRscEnd := fdt.Header.OffMemRsvmap +
		uint32(len(fdt.ReserveEntries)+1)*uint32(unsafe.Sizeof(ReserveEntry{}))
	if fdt.Header.OffDtStruct < memRscEnd {
		return fmt.Errorf(
			"struct block must not overlap memory reservation block, %#x < %#x",
			fdt.Header.OffDtStruct, memRscEnd)
	}
	// TODO: there are more checks which should be done
	return nil
}

func (fdt *FDT) readStringsBlock(f io.ReadSeeker) (strs []byte, err error) {
	if _, err = f.Seek(int64(fdt.Header.OffDtStrings), io.SeekStart); err != nil {
		return
	}
	strs = make([]byte, fdt.Header.SizeDtStrings)
	_, err = f.Read(strs)
	return
}

// readStructBlock reads the nodes and properties of the device and creates the
// tree structure. strs contains the strings block.
func (fdt *FDT) readStructBlock(f io.ReadSeeker, strs []byte) error {
	if fdt.Header.OffDtStruct%4 != 0 {
		return fmt.Errorf(
			"struct offset must be aligned to 4 bytes, but is %#v",
			fdt.Header.OffDtStruct)
	}
	if _, err := f.Seek(int64(fdt.Header.OffDtStruct), io.SeekStart); err != nil {
		return err
	}

	// Buffer file so we don't perform a bajillion syscalls when looking for
	// null-terminating characters.
	r := &uio.AlignReader{
		R: bufio.NewReader(
			&io.LimitedReader{
				R: f,
				N: int64(fdt.Header.SizeDtStruct),
			},
		),
	}

	stack := []*Node{}
	for {
		var t token
		if err := binary.Read(r, binary.BigEndian, &t); err != nil {
			return err
		}
		switch t {
		case tokenBeginNode:
			// Push new node onto the stack.
			child := &Node{}
			stack = append(stack, child)
			if len(stack) == 1 {
				// Root node
				if fdt.RootNode != nil {
					return errors.New("invalid multiple root nodes")
				}
				fdt.RootNode = child
			} else if len(stack) > 1 {
				// Non-root node
				parent := stack[len(stack)-2]
				parent.Children = append(parent.Children, child)
			}

			// The name is a null-terminating string.
			for {
				c := make([]byte, 1)
				if _, err := io.ReadFull(r, c); err != nil {
					return err
				} else if c[0] == 0 {
					break
				} else {
					child.Name += string(c[0])
				}
			}

		case tokenEndNode:
			if len(stack) == 0 {
				return errors.New(
					"unbalanced FDT_BEGIN_NODE and FDT_END_NODE tokens")
			}
			stack = stack[:len(stack)-1]

		case tokenProp:
			pHeader := struct {
				Len, Nameoff uint32
			}{}
			if err := binary.Read(r, binary.BigEndian, &pHeader); err != nil {
				return err
			}
			if pHeader.Nameoff >= uint32(len(strs)) {
				return fmt.Errorf(
					"name offset is larger than strings block: %#x >= %#x",
					pHeader.Nameoff, len(strs))
			}
			null := bytes.IndexByte(strs[pHeader.Nameoff:], 0)
			if null == -1 {
				return fmt.Errorf(
					"property name does not having terminating null at %#x",
					pHeader.Nameoff)
			}
			p := Property{
				Name:  string(strs[pHeader.Nameoff : pHeader.Nameoff+uint32(null)]),
				Value: make([]byte, pHeader.Len),
			}
			_, err := io.ReadFull(r, p.Value)
			if err != nil {
				return err
			}
			if len(stack) == 0 {
				return fmt.Errorf("property %q appears outside a node", p.Name)
			}
			curNode := stack[len(stack)-1]
			curNode.Properties = append(curNode.Properties, p)

		case tokenNop:

		case tokenEnd:
			if uint32(r.N) < fdt.Header.SizeDtStruct {
				return fmt.Errorf(
					"extra data at end of structure block, %#x < %#x",
					uint32(r.N), fdt.Header.SizeDtStruct)
			}
			if fdt.RootNode == nil {
				return errors.New("no root node")
			}
			return nil

		default:
			return fmt.Errorf("undefined token %d", t)
		}

		// Align to four bytes.
		pad, err := r.Align(4)
		if err != nil {
			return err
		}
		for _, v := range pad {
			if v != 0 {
				// TODO: Some of the padding is not zero. Is this a mistake?
				//return fmt.Errorf("padding is non-zero: %d", v)
			}
		}
	}
}

// Write marshals the FDT to an io.Writer and returns the size.
// TODO: implement
func (fdt *FDT) Write(f io.Writer) (int, error) {
	// TODO: implement
	return -1, errors.New("not yet implemented")
}
