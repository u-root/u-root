// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package acpi

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/u-root/u-root/pkg/memio"
)

// Raw ACPI table support. Raw tables are those tables
// one needs to read in, write out, but not change in any way.
// This is needed when, e.g., there is need to create files for
// coreboot cbfs.

// Raw is just a table embedded in a []byte.  Operations on Raw are
// for figuring out how to skip a table you don't care about or, possibly,
// truncating a table and regenerating a checksum.
type Raw struct {
	addr int64
	data []byte
}

var _ = Table(&Raw{})

// NewRaw returns a new Raw []Table fron a given byte slice.
func NewRaw(b []byte) ([]Table, error) {
	var tab []Table
	for len(b) != 0 {
		if len(b) < headerLength {
			return nil, fmt.Errorf("NewRaw: byte slice is only %d bytes and must be at least %d bytes", len(b), headerLength)
		}
		u := binary.LittleEndian.Uint32(b[lengthOffset : lengthOffset+4])
		if int(u) > len(b) {
			return nil, fmt.Errorf("Table length %d is too large for %d byte table (Signature %q", u, len(b), string(b[0:4]))
		}
		tab = append(tab, &Raw{data: b[:u]})
		b = b[u:]
	}
	return tab, nil
}

// RawFromFile reads from an io.Reader and returns a []Table and error if any.
func RawFromFile(r io.Reader) ([]Table, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return NewRaw(b)
}

// RawFromName reads a raw []Table in from a named file.
func RawFromName(n string) ([]Table, error) {
	f, err := os.Open(n)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return RawFromFile(f)
}

// ReadRawTable reads a full table in, given an address.
//
// ReadRawTable uses the io package. This may not always work
// if the kernel has restrictions on reading memory above
// the 1M boundary, and the tables are above boundary.
func ReadRawTable(physAddr int64) (Table, error) {
	var u memio.Uint32
	// Read the table size at a+4
	if err := memio.Read(physAddr+4, &u); err != nil {
		return nil, err
	}
	dat := memio.ByteSlice(make([]byte, u))
	if err := memio.Read(physAddr, &dat); err != nil {
		return nil, err
	}
	return &Raw{addr: physAddr, data: []byte(dat)}, nil
}

// Address returns the table's base address
func (r *Raw) Address() int64 {
	return r.addr
}

// Data returns all the data in a Raw table.
func (r *Raw) Data() []byte {
	return r.data
}

// TableData returns the Raw table, minus the standard ACPI header.
func (r *Raw) TableData() []byte {
	return r.data[headerLength:]
}

// Sig returns the table signature.
func (r *Raw) Sig() string {
	return string(r.data[:4])
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
	return fmt.Sprintf("%q", r.data[10:16])
}

// OEMTableID returns the table OEMTableID.
func (r *Raw) OEMTableID() string {
	return fmt.Sprintf("%q", r.data[16:24])
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
