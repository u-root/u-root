// Copyright 2016-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

// Info contains the SMBIOS information.
type Info struct {
	// TODO(rojer): Add entrypoint information.
	Tables []*Table
}

// ParseInfo parses SMBIOS information from binary data.
func ParseInfo(entry, data []byte) (*Info, error) {
	var tables []*Table
	for len(data) > 0 {
		t, remainder, err := ParseTable(data)
		if err != nil && err != errEndOfTable {
			return nil, err
		}
		tables = append(tables, t)
		data = remainder
	}
	return &Info{Tables: tables}, nil
}

// GetTablesByType returns tables of specific type.
func (i *Info) GetTablesByType(tt TableType) []*Table {
	var res []*Table
	for _, t := range i.Tables {
		if t.Type == tt {
			res = append(res, t)
		}
	}
	return res
}

// GetBIOSInformation returns the Bios Information (type 0) table, if present.
func (i *Info) GetBIOSInformation() (*BIOSInformation, error) {
	bt := i.GetTablesByType(TableTypeBIOSInformation)
	if len(bt) == 0 {
		return nil, ErrTableNotFound
	}
	// There can only be one of these.
	return NewBIOSInformation(bt[0])
}

// GetSystemInformation returns the System Information (type 1) table, if present.
func (i *Info) GetSystemInformation() (*SystemInformation, error) {
	bt := i.GetTablesByType(TableTypeSystemInformation)
	if len(bt) == 0 {
		return nil, ErrTableNotFound
	}
	// There can only be one of these.
	return NewSystemInformation(bt[0])
}
