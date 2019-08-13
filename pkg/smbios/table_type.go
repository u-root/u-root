// Copyright 2016-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"fmt"
)

// TableType specifies the DMI type of the table.
// Types are defined in DMTF DSP0134.
type TableType uint8

// Supported table types.
const (
	TableTypeBIOSInformation      TableType = 0
	TableTypeSystemInformation              = 1
	TableTypeBaseboardInformation           = 2
	TableTypeChassisInformation             = 3
	TableTypeProcessorInformation           = 4
	TableTypeCacheInformation               = 7
	TableTypeMemoryDevice                   = 17
	TableTypeInactive                       = 126
	TableTypeEndOfTable                     = 127
)

func (t TableType) String() string {
	switch t {
	case TableTypeBIOSInformation:
		return "BIOS Information"
	case TableTypeSystemInformation:
		return "System Information"
	case TableTypeBaseboardInformation:
		return "Base Board Information"
	case TableTypeChassisInformation:
		return "Chassis Information"
	case TableTypeProcessorInformation:
		return "Processor Information"
	case TableTypeCacheInformation:
		return "Cache Information"
	case TableTypeMemoryDevice:
		return "Memory Device"
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
func ParseTypedTable(t *Table) (fmt.Stringer, error) {
	switch t.Type {
	case TableTypeBIOSInformation: // 0
		return NewBIOSInformation(t)
	case TableTypeSystemInformation: // 1
		return NewSystemInformation(t)
	case TableTypeBaseboardInformation: // 2
		return NewBaseboardInformation(t)
	case TableTypeChassisInformation: // 3
		return NewChassisInformation(t)
	case TableTypeProcessorInformation: // 4
		return NewProcessorInformation(t)
	case TableTypeCacheInformation: // 7
		return NewCacheInformation(t)
	case TableTypeMemoryDevice: // 17
		return NewMemoryDevice(t)
	case TableTypeInactive: // 126
		// Inactive table cannot be further parsed. Documentation suggests that it can be any table
		// that is temporarily marked inactive by tweaking the type field.
		return t, nil
	case TableTypeEndOfTable: // 127
		return NewEndOfTable(t)
	}
	return nil, ErrUnsupportedTableType
}
