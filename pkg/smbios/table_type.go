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
	TableTypeBIOSInformation       TableType = 0
	TableTypeSystemInformation     TableType = 1
	TableTypeBaseboardInformation  TableType = 2
	TableTypeChassisInformation    TableType = 3
	TableTypeProcessorInformation  TableType = 4
	TableTypeCacheInformation      TableType = 7
	TableTypeMemoryDevice          TableType = 17
	TableTypeIPMIDeviceInformation TableType = 38
	TableTypeTPMDevice             TableType = 43
	TableTypeInactive              TableType = 126
	TableTypeEndOfTable            TableType = 127
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
	case TableTypeIPMIDeviceInformation:
		return "IPMI Device Information"
	case TableTypeTPMDevice:
		return "TPM Device"
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
		return ParseBIOSInformation(t)
	case TableTypeSystemInformation: // 1
		return ParseSystemInformation(t)
	case TableTypeBaseboardInformation: // 2
		return ParseBaseboardInformation(t)
	case TableTypeChassisInformation: // 3
		return ParseChassisInformation(t)
	case TableTypeProcessorInformation: // 4
		return ParseProcessorInformation(t)
	case TableTypeCacheInformation: // 7
		return ParseCacheInformation(t)
	case TableTypeMemoryDevice: // 17
		return NewMemoryDevice(t)
	case TableTypeIPMIDeviceInformation: // 38
		return ParseIPMIDeviceInformation(t)
	case TableTypeTPMDevice: // 43
		return NewTPMDevice(t)
	case TableTypeInactive: // 126
		return NewInactiveTable(t)
	case TableTypeEndOfTable: // 127
		return NewEndOfTable(t)
	}
	return nil, ErrUnsupportedTableType
}
