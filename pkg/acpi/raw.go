// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package acpi

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"

	"github.com/u-root/u-root/pkg/memio"
)

// Raw ACPI table support. Raw tables are those tables
// one needs to read in, write out, but not change in any way.
// This is needed when, e.g., a program has to reassemble all the
// tables in /sys for kexec.

// Raw is just a table embedded in a []byte.  Operations on Raw are
// useful for unpacking into a more refined table or just figuring out
// how to skip a table you don't care about.
type Raw struct {
	data []byte
}

var _ = Tabler(&Raw{})

// NewRaw returns a new Raw table given a byte slice.
func NewRaw(b []byte) (Tabler, error) {
	if len(b) < HeaderLength {
		return nil, fmt.Errorf("NewRaw: byte slice is only %d bytes and must be at least %d bytes", len(b), HeaderLength)
	}
	u := binary.LittleEndian.Uint32(b[LengthOffset : LengthOffset+4])
	return &Raw{data: b[:u]}, nil
}

// RawFromFile reads a raw table in from a file.
func RawFromFile(n string) (Tabler, error) {
	b, err := ioutil.ReadFile(n)
	if err != nil {
		return nil, err
	}
	return NewRaw(b)
}

// ReadRaw reads a full table in, given an address.
// ReadRaw uses the io package. This may not always work
// if the kernel has restrictions on reading memory above
// the 1M boundary, and the tables are above boundary.
func ReadRaw(a int64) (Tabler, error) {
	var u memio.Uint32
	// Read the table size at a+4
	if err := memio.Read(a+4, &u); err != nil {
		return nil, err
	}
	Debug("ReadRaw: Size is %d", u)
	dat := memio.ByteSlice(make([]byte, u))
	if err := memio.Read(a, &dat); err != nil {
		return nil, err
	}
	return &Raw{data: []byte(dat)}, nil
}

// Marshal marshals Raw tables to a byte slice.
func (r *Raw) Marshal() ([]byte, error) {
	return r.data, nil
}

// AllData returns all the data in a Raw table.
func (r *Raw) AllData() []byte {
	return r.data
}

// TableData returns the Raw table, minus the standard ACPI header.
func (r *Raw) TableData() []byte {
	return r.data[HeaderLength:]
}

// Sig returns the table signature.
func (r *Raw) Sig() string {
	return fmt.Sprintf("%s", r.data[:4])
}

// Len returns the total table length.
func (r *Raw) Len() uint32 {
	return uint32(len(r.data))
}

// Revision returns the table Revision.
func (r *Raw) Revision() uint8 {
	return uint8(r.data[8])
}

// CheckSum returns the table CheckSum.
func (r *Raw) CheckSum() uint8 {
	return uint8(r.data[9])
}

// OEMID returns the table OEMID.
func (r *Raw) OEMID() string {
	return fmt.Sprintf("%s", r.data[10:16])
}

// OEMTableID returns the table OEMTableID.
func (r *Raw) OEMTableID() string {
	return fmt.Sprintf("%s", r.data[16:24])
}

// OEMRevision returns the table OEMRevision.
func (r *Raw) OEMRevision() uint32 {
	return binary.LittleEndian.Uint32(r.data[24 : 24+4])
}

// CreatorID returns the table CreatorID.
func (r *Raw) CreatorID() uint32 {
	return binary.LittleEndian.Uint32(r.data[28 : 28+4])
}

// CreatorRevision returns the table CreatorRevision.
func (r *Raw) CreatorRevision() uint32 {
	return binary.LittleEndian.Uint32(r.data[32 : 32+4])
}
