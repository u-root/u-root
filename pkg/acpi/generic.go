// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package acpi

import (
	"encoding/binary"
)

// Generic is the generic ACPI table, with a Header and data
// This makes it possible for users to change certain parts
// of the Header (e.g. vendor id) without touching the data.
// When the table is Marshal'ed out checksums are regenerated.
type Generic struct {
	Header
	data []byte
}

var _ = Tabler(&Generic{})

// NewGeneric creates a new Generic table from a byte slice.
func NewGeneric(b []byte) (Tabler, error) {
	t, err := NewRaw(b)
	if err != nil {
		return nil, err
	}
	return &Generic{Header: *GetHeader(t), data: t.AllData()}, nil
}

// Marshal marshals Generic tables. The main use of this function
// is when you want to tweak the header a bit; you can convert a Raw
// table to a Generic table, do what you wish, and write it out.
func (g *Generic) Marshal() ([]byte, error) {
	// Marshal the header, as it may be changed.
	h, err := g.Header.Marshal()
	if err != nil {
		return nil, err
	}
	// Append only the table data.
	h = append(h, g.TableData()...)
	binary.LittleEndian.PutUint32(h[LengthOffset:], uint32(len(h)))
	h[CSUMOffset] = 0
	c := gencsum(h)
	Debug("CSUM is %#x", c)
	h[CSUMOffset] = c
	return h, nil
}

// Len returns the length of an entire table.
func (g *Generic) Len() uint32 {
	return uint32(len(g.data))
}

// AllData returns the entire table as a byte slice.
func (g *Generic) AllData() []byte {
	return g.data
}

// TableData returns the table, minus the common ACPI header.
func (g *Generic) TableData() []byte {
	return g.data[HeaderLength:]
}

// Sig returns the table signature.
func (g *Generic) Sig() string {
	return string(g.Header.Sig)
}

// OEMID returns the table OEMID.
func (g *Generic) OEMID() string {
	return string(g.Header.OEMID)
}

// OEMTableID returns the table OEMTableID.
func (g *Generic) OEMTableID() string {
	return string(g.Header.OEMTableID)
}

// OEMRevision returns the table OEMRevision.
func (g *Generic) OEMRevision() uint32 {
	return g.Header.OEMRevision
}

// CreatorID returns the table CreatorID.
func (g *Generic) CreatorID() uint32 {
	return g.Header.CreatorID
}

// CreatorRevision returns the table CreatorRevision.
func (g *Generic) CreatorRevision() uint32 {
	return g.Header.CreatorRevision
}

// Revision returns the table Revision.
func (g *Generic) Revision() uint8 {
	return g.Header.Revision
}

// CheckSum returns the table CheckSum.
func (g *Generic) CheckSum() uint8 {
	return g.Header.CheckSum
}
