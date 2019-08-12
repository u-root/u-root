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

// GetBaseboardInformation returns all the Baseboard Information (type 2) tables present.
func (i *Info) GetBaseboardInformation() ([]*BaseboardInformation, error) {
	var res []*BaseboardInformation
	for _, t := range i.GetTablesByType(TableTypeBaseboardInformation) {
		bi, err := NewBaseboardInformation(t)
		if err != nil {
			return nil, err
		}
		res = append(res, bi)
	}
	return res, nil
}

// GetChassisInformation returns all the Chassis Information (type 3) tables present.
func (i *Info) GetChassisInformation() ([]*ChassisInformation, error) {
	var res []*ChassisInformation
	for _, t := range i.GetTablesByType(TableTypeChassisInformation) {
		bi, err := NewChassisInformation(t)
		if err != nil {
			return nil, err
		}
		res = append(res, bi)
	}
	return res, nil
}

// GetProcessorInformation returns all the Processor Information (type 4) tables present.
func (i *Info) GetProcessorInformation() ([]*ProcessorInformation, error) {
	var res []*ProcessorInformation
	for _, t := range i.GetTablesByType(TableTypeProcessorInformation) {
		bi, err := NewProcessorInformation(t)
		if err != nil {
			return nil, err
		}
		res = append(res, bi)
	}
	return res, nil
}
