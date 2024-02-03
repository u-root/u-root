// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package acpi

import "fmt"

const (
	llDSDTAddr = 140
	lDSDTAddr  = 40
)

// BiosTable contains the information needed to create table images for
// firmware such as coreboot or oreboot. It hence includes the RSDP,
// [XR]SDT, and the raw table data. The *SDT is always Tables[0]
// as in the real tables.
type BiosTable struct {
	RSDP   *RSDP
	Tables []Table
}

// ReadBiosTables reads tables that are not interpreted by the OS,
// i.e. it goes straight to memory and gets them there. We optimistically
// hope the bios has not stomped around in low memory messing around.
func ReadBiosTables() (*BiosTable, error) {
	r, err := GetRSDPEBDA()
	if err != nil {
		r, err = GetRSDPMem()
		if err != nil {
			return nil, err
		}
	}
	Debug("Found an RSDP at %#x", r.base)
	x, err := NewSDTAddr(r.SDTAddr())
	if err != nil {
		return nil, err
	}
	Debug("Found an SDT: %s", String(x))
	bios := &BiosTable{
		RSDP:   r,
		Tables: []Table{x},
	}
	for _, a := range x.Addrs {
		Debug("Check Table at %#x", a)
		t, err := ReadRawTable(a)
		if err != nil {
			return nil, fmt.Errorf("%#x:%w", a, err)
		}
		Debug("Add table %s", String(t))
		bios.Tables = append(bios.Tables, t)
		// What I love about ACPI is its unchanging
		// consistency. One table, the FADT, points
		// to another table, the DSDT. There are
		// very good reasons for this:
		// (1) ACPI is a bad design
		// (2) see (1)
		// The signature of the FADT is "FACP".
		// Most appropriate that the names
		// start with F. So does Failure Of Vision.
		if t.Sig() != "FACP" {
			continue
		}
		// 64-bit CPUs had been around for 30 years when ACPI
		// was defined. Nevertheless, they filled it chock full
		// of 32-bit pointers, and then had to go back and paste
		// in 64-bit pointers. The mind reels.
		dsdt, err := getaddr(t.Data(), llDSDTAddr, lDSDTAddr)
		if err != nil {
			return nil, err
		}
		// This is sometimes a kernel virtual address.
		// Fix that.
		t, err = ReadRawTable(int64(uint32(dsdt)))
		if err != nil {
			return nil, fmt.Errorf("%#x:%w", uint64(dsdt), err)
		}
		Debug("Add table %s", String(t))
		bios.Tables = append(bios.Tables, t)

	}
	return bios, nil
}

// RawTablesFromMem reads all the tables from Mem, using the SDT.
func RawTablesFromMem() ([]Table, error) {
	x, err := ReadBiosTables()
	if err != nil {
		return nil, err
	}
	return x.Tables, err
}
