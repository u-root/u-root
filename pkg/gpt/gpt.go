// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gpt implements reading and writing of GUID Partition tables.
//
// GPTs are dumped in JSON format and written in same.  One complication is
// that we frequently only want to write a very small subset of a GPT. For
// example, we might only want to change the GUID. As it happens it is simpler
// (and more useful) just to read and write the whole thing. In for a penny, in
// for a pound.
package gpt

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"

	"github.com/rekby/gpt"
	"github.com/u-root/u-root/pkg/uio"
)

const (
	BlockSize = 512
	HeaderOff = 0x200
)

// PartitionTable defines all the partition table information.
//
// This includes the MBR and two GPTs.
//
// The GPTs are similar but not identical, as they contain "pointers" to each
// other in the HeaderCopyStartLBA in the Header. The design is defective in that if a
// given Header has an error, you are supposed to just assume that the
// HeaderCopyStartLBA is intact, which is a pretty bogus assumption. This is why you do
// standards like this in the open, not in hiding. I hope someone from Intel is
// reading this.
type PartitionTable struct {
	MasterBootRecord *MBR
	GPT              gpt.Table
}

type MBR [BlockSize]byte

func (m *MBR) String() string {
	b, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		log.Fatalf("Can't marshal %v", *m)
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
func EqualHeader(p, b gpt.Header) error {
	var err error
	if p.Signature != b.Signature {
		err = errAppend(err, "p.Signature(%#x) != b.Signature(%#x)", p.Signature, b.Signature)
	}
	if p.Revision != b.Revision {
		err = errAppend(err, "p.Revision(%v) != b.Revision(%v)", p.Revision, b.Revision)
	}
	if p.Size != b.Size {
		err = errAppend(err, "p.Size(%v) != b.Size(%v)", p.Size, b.Size)
	}
	if p.HeaderStartLBA != b.HeaderCopyStartLBA {
		err = errAppend(err, "p.HeaderStartLBA(%#x) != b.HeaderCopyStartLBA(%#x)", p.HeaderStartLBA, b.HeaderCopyStartLBA)
	}
	if p.HeaderCopyStartLBA != b.HeaderStartLBA {
		err = errAppend(err, "p.HeaderCopyStartLBA(%#x) != b.HeaderStartLBA(%#x)", p.HeaderCopyStartLBA, b.HeaderStartLBA)
	}
	if p.FirstUsableLBA != b.FirstUsableLBA {
		err = errAppend(err, "p.FirstUsableLBA(%#x) != b.FirstUsableLBA(%#x)", p.FirstUsableLBA, b.FirstUsableLBA)
	}
	if p.LastUsableLBA != b.LastUsableLBA {
		err = errAppend(err, "p.LastUsableLBA(%#x) != b.LastUsableLBA(%#x)", p.LastUsableLBA, b.LastUsableLBA)
	}
	if p.DiskGUID != b.DiskGUID {
		err = errAppend(err, "p.DiskGUID(%#x) != b.DiskGUID(%#x)", p.DiskGUID, b.DiskGUID)
	}
	if p.PartitionsArrLen != b.PartitionsArrLen {
		err = errAppend(err, "p.PartitionsArrLen(%v) != b.PartitionsArrLen(%v)", p.PartitionsArrLen, b.PartitionsArrLen)
	}
	if p.PartitionEntrySize != b.PartitionEntrySize {
		err = errAppend(err, "p.PartitionEntrySize(%v) != b.PartitionEntrySize(%v)", p.PartitionEntrySize, b.PartitionEntrySize)
	}
	return err
}

// Write writes the MBR and primary and backup GPTs to w.
func Write(w io.WriterAt, p *PartitionTable) error {
	if _, err := w.WriteAt(p.MasterBootRecord[:], 0); err != nil {
		return err
	}
	sw := uio.NewSectionWriter(w, 0, math.MaxInt64)
	if _, err := sw.Seek(HeaderOff, io.SeekStart); err != nil {
		return err
	}
	if err := p.GPT.Write(sw); err != nil {
		return err
	}
	backup := p.GPT.CreateOtherSideTable()
	if _, err := sw.Seek(int64(p.GPT.Header.HeaderCopyStartLBA*BlockSize), io.SeekStart); err != nil {
		return err
	}
	if err := backup.Write(sw); err != nil {
		return err
	}
	return nil
}

// New reads in the MBR, primary and backup GPT from a disk and returns a pointer to them.
// There are some checks it can apply. It can return with a
// one or more headers AND an error. Sorry. Experience with some real USB sticks
// is showing that we need to return data even if there are some things wrong.
func New(r io.ReaderAt) (*PartitionTable, error) {
	var mbr = &MBR{}
	n, err := r.ReadAt(mbr[:], 0)
	if n != BlockSize || err != nil {
		return nil, err
	}

	var p PartitionTable
	p.MasterBootRecord = mbr

	reader := io.NewSectionReader(r, 0, math.MaxInt64)
	if _, err := reader.Seek(HeaderOff, io.SeekStart); err != nil {
		return nil, err
	}

	primary, err := gpt.ReadTable(reader, BlockSize)
	if err != nil {
		return nil, fmt.Errorf("reading primary table failed: %v", err)
	}
	p.GPT = primary

	if _, err := reader.Seek(int64(primary.Header.HeaderCopyStartLBA*BlockSize), io.SeekStart); err != nil {
		return &p, err
	}
	backup, err := gpt.ReadTable(reader, BlockSize)
	if err != nil {
		// you go to war with the army you have
		return &p, fmt.Errorf("reading backup table failed: %v", err)
	}

	if err := EqualHeader(primary.Header, backup.Header); err != nil {
		return &p, fmt.Errorf("primary GPT and backup GPT header differ: %v", err)
	}

	if primary.Header.CRC == backup.Header.CRC {
		return &p, fmt.Errorf("primary (%v) and backup (%v) header CRC (%x) are the same and should differ", primary.Header, backup.Header, primary.Header.CRC)
	}
	if primary.Header.PartitionsCRC != backup.Header.PartitionsCRC {
		return &p, fmt.Errorf("primary (%v) and backup (%v) partitions CRC do not match", primary, backup)
	}
	return &p, nil
}

// GetBlockSize returns the block size of device.
func GetBlockSize(device string) (int, error) {
	// TODO: scan device to determine block size.
	return BlockSize, nil
}
