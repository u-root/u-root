// Copyright 2016-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

// TableType specifies the DMI type of the table.
// Types are defined in DMTF DSP0134.
type TableType uint8

// Supported table types.
const (
	TableTypeBIOSInformation   TableType = 0
	TableTypeSystemInformation           = 1
	TableTypeInactive                    = 126
	TableTypeEndOfTable                  = 127
)

func (t TableType) String() string {
	switch t {
	case TableTypeBIOSInformation:
		return "BIOS Information"
	case TableTypeSystemInformation:
		return "System Information"
	case TableTypeInactive:
		return "Inactive"
	case TableTypeEndOfTable:
		return "End Of Table"
	default:
		if t >= 0x80 {
			return "OEM-specific Type"
		}
		return "Unsupported"
	}
}

// ParseTypedTable parses generic Table into a typed struct.
func ParseTypedTable(t *Table) (interface{}, error) {
	switch t.Type {
	case TableTypeBIOSInformation:
		return NewBIOSInformation(t)
	case TableTypeSystemInformation:
		return NewSystemInformation(t)
	case TableTypeInactive:
		// Inactive table cannot be further parsed. Documentation suggests that it can be any table
		// that is temporarily marked inactive by tweaking the type field.
		return t, nil
	}
	return nil, ErrUnsupportedTableType
}
