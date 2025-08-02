// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import "testing"

func TestTableTypeString(t *testing.T) {
	tests := []struct {
		tableType TableType
		want      string
	}{
		{
			tableType: TableTypeBIOSInfo,
			want:      "BIOS Information",
		},
		{
			tableType: TableTypeSystemInfo,
			want:      "System Information",
		},
		{
			tableType: TableTypeBaseboardInfo,
			want:      "Base Board Information",
		},
		{
			tableType: TableTypeChassisInfo,
			want:      "Chassis Information",
		},
		{
			tableType: TableTypeProcessorInfo,
			want:      "Processor Information",
		},
		{
			tableType: TableTypeCacheInfo,
			want:      "Cache Information",
		},
		{
			tableType: TableTypeSystemSlots,
			want:      "System Slots",
		},
		{
			tableType: TableTypeMemoryDevice,
			want:      "Memory Device",
		},
		{
			tableType: TableTypeIPMIDeviceInfo,
			want:      "IPMI Device Information",
		},
		{
			tableType: TableTypeTPMDevice,
			want:      "TPM Device",
		},
		{
			tableType: TableTypeInactive,
			want:      "Inactive",
		},
		{
			tableType: TableTypeEndOfTable,
			want:      "End Of Table",
		},
		{
			tableType: TableType(0x12),
			want:      "Unsupported",
		},
		{
			tableType: TableType(0x81),
			want:      "OEM-specific Type",
		},
	}

	for _, tt := range tests {
		if tt.tableType.String() != tt.want {
			t.Errorf("TableType.String(): '%s', want '%s'", tt.tableType.String(), tt.want)
		}
	}
}
