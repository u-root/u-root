// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package acpi

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

var sdtLen = 36

// SDT contains information about tables. It does not differentiate 32- vs 64-bit
// tables.
type SDT struct {
	// Table is the SDT itself.
	Table

	// Addrs is the array of physical addresses in the SDT.
	Addrs []int64

	// Base is the SDT base, used to generate a new SDT for, e.g, coreboot.
	Base int64
}

var _ = Table(&SDT{})

// NewSDTAddr returns an SDT, given an address.
func NewSDTAddr(addr int64) (*SDT, error) {
	t, err := ReadRawTable(addr)
	if err != nil {
		return nil, fmt.Errorf("can not read SDT at %#x", addr)
	}
	Debug("NewSDTAddr: %s %#x", String(t), t.TableData())
	return NewSDT(t, addr)
}

// NewSDT returns an SDT, given a Table
func NewSDT(t Table, addr int64) (*SDT, error) {
	s := &SDT{
		Table: t,
		Base:  addr,
	}
	r := bytes.NewReader(t.TableData())
	x := true
	if t.Sig() == "RSDT" {
		x = false
	}
	for r.Len() > 0 {
		var a int64
		if x {
			if err := binary.Read(r, binary.LittleEndian, &a); err != nil {
				return nil, err
			}
		} else {
			var a32 uint32
			if err := binary.Read(r, binary.LittleEndian, &a32); err != nil {
				return nil, err
			}
			a = int64(a32)
		}
		Debug("Add addr %#x", a)
		s.Addrs = append(s.Addrs, a)
	}
	Debug("NewSDT: %v", s.String())
	return s, nil
}

// String implements string for an SDT.
func (sdt *SDT) String() string {
	return fmt.Sprintf("%s at %#x with %d tables: %#x", String(sdt), sdt.Base, len(sdt.Addrs), sdt.Addrs)
}
