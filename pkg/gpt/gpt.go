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
)

const (
	BlockSize         = 512
	HeaderOff         = 0x200
	HeaderSize        = 0x5c               // They claim it can vary. Give me a break.
	Signature  uint64 = 0x5452415020494645 // ("EFI PART", 45h 46h 49h 20h 50h 41h 52h 54h on little-endian machines)
	Revision          = 0x10000
	MaxNPart          = 0x80
)

type GUID struct {
	L  uint32
	W1 uint16
	W2 uint16
	B  [8]byte
}

type MBR [BlockSize]byte
type Header struct {
	Signature  uint64
	Revision   uint32 // (for GPT version 1.0 (through at least UEFI version 2.7 (May 2017)), the value is 00h 00h 01h 00h)
	HeaderSize uint32 // size in little endian (in bytes, usually 5Ch 00h 00h 00h or 92 bytes)
	CRC        uint32 // CRC32/zlib of header (offset +0 up to header size) in little endian, with this field zeroed during calculation
	Reserved   uint32 // ; must be zero
	CurrentLBA uint64 // (location of this header copy)
	BackupLBA  uint64 // (location of the other header copy)
	FirstLBA   uint64 // usable LBA for partitions (primary partition table last LBA + 1)
	LastLBA    uint64 // usable LBA (secondary partition table first LBA - 1)
	DiskGUID   GUID   // (also referred as UUID on UNIXes)
	PartStart  uint64 // LBA of array of partition entries (always 2 in primary copy)
	NPart      uint32 // Number of partition entries in array
	PartSize   uint32 // Size of a single partition entry (usually 80h or 128)
	PartCRC    uint32 // CRC32/zlib of partition array in little endian
}

type PartAttr uint64
type PartName [72]byte
type Part struct {
	PartGUID   GUID     // Partition type GUID
	UniqueGUID GUID     // Unique partition GUID
	FirstLBA   uint64   // LBA (little endian)
	LastLBA    uint64   // LBA (inclusive, usually odd)
	Attribute  PartAttr // flags (e.g. bit 60 denotes read-only)
	Name       PartName // Partition name (36 UTF-16LE code units)
}

type GPT struct {
	Header
	Parts []Part
}

func (g *GUID) String() string {
	return fmt.Sprintf("%08x-%04x-%04x-%02x-%02x", g.L, g.W1, g.W2, g.B[0:2], g.B[2:])
}

// PartitionTable defines all the partition table information.
// This includes the MBR and two GPTs. The GPTs are
// similar but not identical, as they contain "pointers"
// to each other in the BackupLBA in the Header.
// The design is defective in that if a given Header has
// an error, you are supposed to just assume that the BackupLBA
// is intact, which is a pretty bogus assumption. This is why
// you do standards like this in the open, not in hiding.
// I hope someone from Intel is reading this.
type PartitionTable struct {
	MasterBootRecord *MBR
	Primary          *GPT
	Backup           *GPT
}

func (m *MBR) String() string {
	b, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		log.Fatalf("Can't marshal %v", *m)
	}
	return string(b)
}

func (g *GPT) String() string {
	b, err := json.MarshalIndent(g, "", "\t")
	if err != nil {
		log.Fatalf("Can't marshal %v", *g)
	}
	return string(b)
}

func (p *PartitionTable) String() string {
	b, err := json.MarshalIndent(p, "", "\t")
	if err != nil {
		log.Fatalf("Can't marshal %v", *p)
	}
	return string(b)

}

func errAppend(err error, s string, a ...interface{}) error {
	var p string
	if err != nil {
		p = err.Error() + "; "
	}
	return fmt.Errorf(p+s, a...)
}

// EqualHeader compares two headers and returns true if they match.
// Those fields which by definition must differ are ignored.
func EqualHeader(p, b Header) error {
	var err error
	if p.Signature != b.Signature {
		err = errAppend(err, "p.Signature(%#x) != b.Signature(%#x)", p.Signature, b.Signature)
	}
	if p.Revision != b.Revision {
		err = errAppend(err, "p.Revision(%v) != b.Revision(%v)", p.Revision, b.Revision)
	}
	if p.HeaderSize != b.HeaderSize {
		err = errAppend(err, "p.HeaderSize(%v) != b.HeaderSize(%v)", p.HeaderSize, b.HeaderSize)
	}
	if p.CurrentLBA != b.BackupLBA {
		err = errAppend(err, "p.CurrentLBA(%#x) != b.BackupLBA(%#x)", p.CurrentLBA, b.BackupLBA)
	}
	if p.BackupLBA != b.CurrentLBA {
		err = errAppend(err, "p.BackupLBA(%#x) != b.CurrentLBA(%#x)", p.BackupLBA, b.CurrentLBA)
	}
	if p.FirstLBA != b.FirstLBA {
		err = errAppend(err, "p.FirstLBA(%#x) != b.FirstLBA(%#x)", p.FirstLBA, b.FirstLBA)
	}
	if p.LastLBA != b.LastLBA {
		err = errAppend(err, "p.LastLBA(%#x) != b.LastLBA(%#x)", p.LastLBA, b.LastLBA)
	}
	if p.DiskGUID != b.DiskGUID {
		err = errAppend(err, "p.DiskGUID(%#x) != b.DiskGUID(%#x)", p.DiskGUID, b.DiskGUID)
	}
	if p.NPart != b.NPart {
		err = errAppend(err, "p.NPart(%v) != b.NPart(%v)", p.NPart, b.NPart)
	}
	if p.PartSize != b.PartSize {
		err = errAppend(err, "p.PartSize(%v) != b.PartSize(%v)", p.PartSize, b.PartSize)
	}
	return err
}

func EqualPart(p, b Part) (err error) {
	if p.PartGUID != b.PartGUID {
		err = errAppend(err, "p.PartGUID(%#x) != b.PartGUID(%#x)", p.PartGUID, b.PartGUID)
	}
	if p.UniqueGUID != b.UniqueGUID {
		err = errAppend(err, "p.UniqueGUID(%#x) != b.UniqueGUID(%#x)", p.UniqueGUID, b.UniqueGUID)
	}
	if p.FirstLBA != b.FirstLBA {
		err = errAppend(err, "p.FirstLBA(%#x) != b.FirstLBA(%#x)", p.FirstLBA, b.FirstLBA)
	}
	if p.LastLBA != b.LastLBA {
		err = errAppend(err, "p.LastLBA(%#x) != b.LastLBA(%#x)", p.LastLBA, b.LastLBA)
	}
	// TODO: figure out what Attributes actually matter. We're not able to tell what differences
	// matter and what differences don't.
	if false && p.Attribute != b.Attribute {
		err = errAppend(err, "p.Attribute(%#x) != b.Attribute(%#x)", p.Attribute, b.Attribute)
	}
	if p.Name != b.Name {
		err = errAppend(err, "p.Name(%#x) != b.Name(%#x)", p.Name, b.Name)
	}
	return err
}

// EqualParts compares the Parts arrays from two GPTs
// and returns an error if they differ.
// If they length differs we just give up, since there's no way
// to know which should have matched.
// Otherwise, we do a 1:1 comparison.
func EqualParts(p, b *GPT) (err error) {
	if len(p.Parts) != len(b.Parts) {
		return fmt.Errorf("Primary Number of partitions (%d) differs from Backup (%d)", len(p.Parts), len(b.Parts))
	}
	for i := range p.Parts {
		if e := EqualPart(p.Parts[i], b.Parts[i]); e != nil {
			err = errAppend(err, "Partition %d: %v", i, e)
		}
	}
	return err
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

// Write writes the MBR and primary and backup GPTs to w.
func Write(w io.WriterAt, p *PartitionTable) error {
	if _, err := w.WriteAt(p.MasterBootRecord[:], 0); err != nil {
		return err
	}
	if err := writeGPT(w, p.Primary); err != nil {
		return err
	}

	if err := writeGPT(w, p.Backup); err != nil {
		return err
	}
	return nil

}

// Write writes the GPT to w. It generates the partition and header CRC before writing.
func writeGPT(w io.WriterAt, g *GPT) error {
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

// New reads in the MBR, primary and backup GPT from a disk and returns a pointer to them.
// There are some checks it can apply. It can return with a
// one or more headers AND an error. Sorry. Experience with some real USB sticks
// is showing that we need to return data even if there are some things wrong.
func New(r io.ReaderAt) (*PartitionTable, error) {
	var p = &PartitionTable{}
	var mbr = &MBR{}
	n, err := r.ReadAt(mbr[:], 0)
	if n != BlockSize || err != nil {
		return p, err
	}
	p.MasterBootRecord = mbr
	g, err := Table(r, HeaderOff)
	// If we can't read the table it's kinda hard to find the backup.
	// Bit of a flaw in the design, eh?
	// "We can't recover the backup from the error with the primary because there
	// was an error in the primary"
	// uh, what?
	if err != nil {
		return p, err
	}
	p.Primary = g

	b, err := Table(r, int64(g.BackupLBA*BlockSize))
	if err != nil {
		// you go to war with the army you have
		return p, err
	}

	if err := EqualHeader(g.Header, b.Header); err != nil {
		return p, fmt.Errorf("Primary GPT and backup GPT Header differ: %v", err)
	}

	if g.CRC == b.CRC {
		return p, fmt.Errorf("Primary (%v) and Backup (%v) Header CRC (%x) are the same and should differ", g.Header, b.Header, g.CRC)
	}

	p.Backup = b
	return p, EqualParts(g, b)
}
