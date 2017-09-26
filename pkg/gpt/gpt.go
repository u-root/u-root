// Copyright 2009-2017 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// gpt implements reading and writing of GUID Partition tables.
// GPTs are dumped in JSON format and written in same.
// One complication is that we frequently only want to
// write a very small subset of a GPT. For example,
// we might only want to change the GUID. As it happens
// it is simpler (and more useful) just to read and write
// the whole thing. In for a penny, in for a pound.
package gpt

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"reflect"

	"github.com/google/uuid"
)

const (
	BlockSize  = 512
	HeaderOff  = 0x200
	HeaderSize = 0x5c               // They claim it can vary. Give me a break.
	Signature  = 0x5452415020494645 // ("EFI PART", 45h 46h 49h 20h 50h 41h 52h 54h on little-endian machines)
	Revision   = 0x10000
	MaxNPart   = 0x80
)

type Header struct {
	Signature  uint64
	Revision   uint32    // (for GPT version 1.0 (through at least UEFI version 2.7 (May 2017)), the value is 00h 00h 01h 00h)
	HeaderSize uint32    // size in little endian (in bytes, usually 5Ch 00h 00h 00h or 92 bytes)
	CRC        uint32    // CRC32/zlib of header (offset +0 up to header size) in little endian, with this field zeroed during calculation
	Reserved   uint32    // ; must be zero
	CurrentLBA uint64    // (location of this header copy)
	BackupLBA  uint64    // (location of the other header copy)
	FirstLBA   uint64    // usable LBA for partitions (primary partition table last LBA + 1)
	LastLBA    uint64    // usable LBA (secondary partition table first LBA - 1)
	DiskGUID   uuid.UUID // (also referred as UUID on UNIXes)
	PartStart  uint64    // LBA of array of partition entries (always 2 in primary copy)
	NPart      uint32    // Number of partition entries in array
	PartSize   uint32    // Size of a single partition entry (usually 80h or 128)
	PartCRC    uint32    // CRC32/zlib of partition array in little endian
}

type PartAttr uint64
type PartName [72]byte
type Part struct {
	PartGUID   uuid.UUID // Partition type GUID
	UniqueGUID uuid.UUID // Unique partition GUID
	FirstLBA   uint64    // LBA (little endian)
	LastLBA    uint64    // LBA (inclusive, usually odd)
	Attribute  PartAttr  // flags (e.g. bit 60 denotes read-only)
	Name       PartName  // Partition name (36 UTF-16LE code units)
}

type GPT struct {
	Header
	Parts []Part
}

func (g *GPT) String() string {
	b, err := json.MarshalIndent(g, "", "\t")
	if err != nil {
		log.Fatalf("Can't marshal %v", *g)
	}
	return string(b)
}

// EqualHeader compares two headers and returns true if they match.
// Those fields which by definition must differ are ignored.
func EqualHeader(p, b Header) bool {
	return p.Signature == b.Signature &&
		p.Revision == b.Revision &&
		p.HeaderSize == b.HeaderSize &&
		p.CurrentLBA == b.BackupLBA &&
		p.BackupLBA == b.CurrentLBA &&
		p.FirstLBA == b.FirstLBA &&
		p.LastLBA == b.LastLBA &&
		p.DiskGUID == b.DiskGUID &&
		p.NPart == b.NPart &&
		p.PartSize == b.PartSize &&
		p.PartCRC == b.PartCRC
}

// Table reads a GPT table at the given offset.  It checks that
// the Signature, Revision, HeaderSize, and MaxPart are reasonable. It
// also verifies the header and partition table CRC32 values.
func Table(r io.ReaderAt, off int64) (*GPT, error) {
	which := "Primary"
	if off != BlockSize {
		which = "Backup"
	}
	var g = &GPT{}
	if err := binary.Read(io.NewSectionReader(r, off, HeaderSize), binary.LittleEndian, &g.Header); err != nil {
		return nil, err
	}

	if g.Signature != Signature {
		return nil, fmt.Errorf("%s GPT signature invalid (%x), needs to be %x", which, g.Signature, Signature)
	}
	if g.Revision != Revision {
		return nil, fmt.Errorf("%s GPT revision (%x) is not supported value (%x)", which, g.Revision, Revision)
	}
	if g.HeaderSize != HeaderSize {
		return nil, fmt.Errorf("%s GPT HeaderSize (%x) is not supported value (%x)", which, g.HeaderSize, HeaderSize)
	}
	if g.NPart > MaxNPart {
		return nil, fmt.Errorf("%s GPT MaxNPart (%x) is above maximum of %x", which, g.NPart, MaxNPart)
	}

	// Read in all the partition data and check the hash.
	// Since the partition hash is included in the header CRC,
	// it's sensible to check it first.
	s := int64(g.PartSize)
	partBlocks := make([]byte, int64(g.NPart)*s)
	n, err := r.ReadAt(partBlocks, int64(g.PartStart)*BlockSize)
	if n != len(partBlocks) || err != nil {
		return nil, fmt.Errorf("%s Reading partitions: Wanted %d bytes, got %d: %v", which, n, len(partBlocks), err)
	}
	if h := crc32.ChecksumIEEE(partBlocks); h != g.PartCRC {
		return g, fmt.Errorf("%s Partition CRC: Header %v, computed checksum is %08x, header has %08x", which, g, h, g.PartCRC)
	}

	hdr := make([]byte, g.HeaderSize)
	n, err = r.ReadAt(hdr, off)
	if n != len(hdr) || err != nil {
		return nil, fmt.Errorf("%s Reading Header: Wanted %d bytes, got %d: %v", which, n, len(hdr), err)
	}
	// Zap the checksum in the header to 0.
	copy(hdr[16:], []byte{0, 0, 0, 0})
	if h := crc32.ChecksumIEEE(hdr); h != g.CRC {
		return g, fmt.Errorf("%s Header CRC: computed checksum is %08x, header has %08x", which, h, g.CRC)
	}

	// Now read in the partition table entries.
	g.Parts = make([]Part, g.NPart)

	for i := range g.Parts {
		if err := binary.Read(io.NewSectionReader(r, int64(g.PartStart*BlockSize)+int64(i)*s, s), binary.LittleEndian, &g.Parts[i]); err != nil {
			return nil, fmt.Errorf("%s GPT partition %d failed: %v", which, i, err)
		}
	}

	return g, nil

}

// Write writes the GPT to w, both primary and backup. It generates the CRCs before writing.
// It takes an io.Writer and assumes that you are correctly positioned in the output stream.
// This means we must adjust partition numbers by subtracting one from them.
func Write(w io.WriterAt, g *GPT) error {
	// The maximum extent is NPart * PartSize
	var h = make([]byte, uint64(g.NPart*g.PartSize))
	s := int64(g.PartSize)
	for i := int64(0); i < int64(g.NPart); i++ {
		var b bytes.Buffer
		if err := binary.Write(&b, binary.LittleEndian, &g.Parts[i]); err != nil {
			return err
		}
		copy(h[i*s:], b.Bytes())
	}

	ps := int64(g.PartStart * BlockSize)
	if _, err := w.WriteAt(h, ps); err != nil {
		return fmt.Errorf("Writing %d bytes of partition table at %v: %v", len(h), ps, err)
	}

	g.PartCRC = crc32.ChecksumIEEE(h[:])
	g.CRC = 0
	var b bytes.Buffer
	if err := binary.Write(&b, binary.LittleEndian, &g.Header); err != nil {
		return err
	}
	h = make([]byte, g.HeaderSize)
	copy(h, b.Bytes())
	g.CRC = crc32.ChecksumIEEE(h[:])
	b.Reset()
	if err := binary.Write(&b, binary.LittleEndian, g.CRC); err != nil {
		return err
	}
	copy(h[16:], b.Bytes())

	_, err := w.WriteAt(h, int64(g.CurrentLBA*BlockSize))
	return err
}

// New reads in the primary and backup GPT from a disk and returns a pointer to them.
// There are some checks it can apply. It can return with a
// one or more headers AND an error. Sorry. Experience with some real USB sticks
// is showing that we need to return data even if there are some things wrong.
func New(r io.ReaderAt) (*GPT, *GPT, error) {
	g, err := Table(r, HeaderOff)
	// If we can't read the table it's kinda hard to find the backup.
	// Bit of a flaw in the design, eh?
	if err != nil {
		return nil, nil, err
	}

	b, err := Table(r, int64(g.BackupLBA*BlockSize))
	if err != nil {
		// you go to war with the army you have
		return g, nil, err
	}

	if !EqualHeader(g.Header, b.Header) {
		return g, b, fmt.Errorf("Primary GPT(%s) and backup GPT(%s) Header differ", g, b)
	}

	if g.CRC == b.CRC {
		return g, b, fmt.Errorf("Primary (%v) and Backup (%v) Header CRC (%x) are the same and should differ", g.Header, b.Header, g.CRC)
	}

	if !reflect.DeepEqual(g.Parts, b.Parts) {
		return b, g, fmt.Errorf("Primary GPT(%s) and backup GPT(%s) Parts differ", g, b)
	}
	return g, b, nil
}
