// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package acpi is designed to support copying individual tables or
// a blob containing many tables from one spot to another, supporting
// filtering. For example, one might read tables from /dev/mem, using
// the RSDP, so as to create an ACPI table blob for use in coreboot.
// In this case, we only care about checking the signature.
package acpi

import (
	"fmt"
	"io"
)

const (
	// lengthOffset is the offset of the table length
	lengthOffset = 4
	// checksum1 offset in RSDP struct.
	cSUM1Off = 8
	// cSUMOffset is the offset of the single byte checksum in *most* ACPI tables
	cSUMOffset = 9
	// minTableLength is the minimum length: 4 byte tag, 4 byte length, 1 byte revision, 1 byte checksum,
	minTableLength = 10

	// checksum2 offset in RSDP struct.
	cSUM2Off    = 32
	xSDTLenOff  = 20
	xSDTAddrOff = 24

	// headerLength is a common header length for (almost)
	// all ACPI tables.
	headerLength = 36
)

type (
	// Table is an individual ACPI table.
	Table interface {
		Sig() string
		Len() uint32
		Revision() uint8
		CheckSum() uint8
		OEMID() string
		OEMTableID() string
		OEMRevision() uint32
		CreatorID() uint32
		CreatorRevision() uint32
		Data() []byte
		TableData() []byte
	}
	// TableMethod defines the type of functions used to read a table.
	TableMethod func() ([]Table, error)
)

// gencsum generates a uint8 checksum of a []uint8
func gencsum(b []uint8) uint8 {
	var csum uint8
	for _, bb := range b {
		csum += bb
	}
	return ^csum + 1
}

// Method accepts a method name and returns a TableMethod if one exists, or error othewise.
func Method(n string) (TableMethod, error) {
	f, ok := Methods[n]
	if !ok {
		return nil, fmt.Errorf("only method[s] %q are available, not %q", MethodNames(), n)
	}
	return f, nil
}

// String pretty-prints a Table
func String(t Table) string {
	return fmt.Sprintf("%s %d %d %#02x %s %s %#08x %#08x %#08x",
		t.Sig(),
		t.Len(),
		t.Revision(),
		t.CheckSum(),
		t.OEMID(),
		t.OEMTableID(),
		t.OEMRevision(),
		t.CreatorID(),
		t.CreatorRevision())

}

// WriteTables writes one or more tables to an io.Writer.
func WriteTables(w io.Writer, tab Table, tabs ...Table) error {
	for _, tt := range append([]Table{tab}, tabs...) {
		if _, err := w.Write(tt.Data()); err != nil {
			return fmt.Errorf("Writing %s: %v", tt.Sig(), err)
		}
	}
	return nil
}

// ReadTables reads tables, given a method name.
func ReadTables(n string) ([]Table, error) {
	f, err := Method(n)
	if err != nil {
		return nil, err
	}
	return f()
}
