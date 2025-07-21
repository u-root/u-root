// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import "testing"

func TestBMCInterfaceTypeString(t *testing.T) {
	tests := []struct {
		name       string
		BMCIntType BMCInterfaceType
		want       string
	}{
		{
			name:       "BMCInterfaceTypeUnknown",
			BMCIntType: BMCInterfaceTypeUnknown,
			want:       "Unknown",
		},
		{
			name:       "BMCInterfaceTypeKCSKeyboardControllerStyle",
			BMCIntType: BMCInterfaceTypeKCSKeyboardControllerStyle,
			want:       "KCS (Keyboard Control Style)",
		},
		{
			name:       "BMCInterfaceTypeSMICServerManagementInterfaceChip",
			BMCIntType: BMCInterfaceTypeSMICServerManagementInterfaceChip,
			want:       "SMIC (Server Management Interface Chip)",
		},
		{
			name:       "BMCInterfaceTypeBTBlockTransfer",
			BMCIntType: BMCInterfaceTypeBTBlockTransfer,
			want:       "BT (Block Transfer)",
		},
		{
			name:       "BMCInterfaceTypeSSIFSMBusSystemInterface",
			BMCIntType: BMCInterfaceTypeSSIFSMBusSystemInterface,
			want:       "SSIF (SMBus System Interface)",
		},
		{
			name:       "Unknown Type 0x05",
			BMCIntType: 0x05,
			want:       "0x5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.BMCIntType.String()
			if got != tt.want {
				t.Errorf("BMCInterfaceType.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIPMIDeviceInfoString(t *testing.T) {
	tests := []struct {
		name    string
		IPMIDev IPMIDeviceInfo
		want    string
	}{
		{
			name:    "Empty Struct",
			IPMIDev: IPMIDeviceInfo{},
			want: `Handle 0x0000, DMI type 0, 0 bytes
BIOS Information
	Interface Type: Unknown
	Specification Version: 0.0
	I2C Slave Address: 0x00
	NV Storage Device: 0
	Base Address: 0x0000000000000000 (Memory-mapped)
	Register Spacing: Successive Byte Boundaries`,
		},
		{
			name: "I/O mapped,32bit Boundaries, Interrupt Active High Edge",
			IPMIDev: IPMIDeviceInfo{
				BaseAddress:                      0x1,
				BaseAddressModifierInterruptInfo: 0x4a,
			},
			want: `Handle 0x0000, DMI type 0, 0 bytes
BIOS Information
	Interface Type: Unknown
	Specification Version: 0.0
	I2C Slave Address: 0x00
	NV Storage Device: 0
	Base Address: 0x0000000000000000 (I/O)
	Register Spacing: 32-bit Boundaries
	Interrupt Polarity: Active High
	Interrupt Trigger Mode: Edge`,
		},
		{
			name: "I/O mapped,32bit Boundaries, Interrupt Active High Level Interrupt 1",
			IPMIDev: IPMIDeviceInfo{
				BaseAddress:                      0x1,
				BaseAddressModifierInterruptInfo: 0x89,
				InterruptNumber:                  1,
			},
			want: `Handle 0x0000, DMI type 0, 0 bytes
BIOS Information
	Interface Type: Unknown
	Specification Version: 0.0
	I2C Slave Address: 0x00
	NV Storage Device: 0
	Base Address: 0x0000000000000000 (I/O)
	Register Spacing: 16-bit Boundaries
	Interrupt Polarity: Active Low
	Interrupt Trigger Mode: Level
	Interrupt Number: 1`,
		},
		{
			name: "Out of Spec",
			IPMIDev: IPMIDeviceInfo{
				BaseAddressModifierInterruptInfo: 0xc0,
			},
			want: `Handle 0x0000, DMI type 0, 0 bytes
BIOS Information
	Interface Type: Unknown
	Specification Version: 0.0
	I2C Slave Address: 0x00
	NV Storage Device: 0
	Base Address: 0x0000000000000000 (Memory-mapped)
	Register Spacing: <OUT OF SPEC>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.IPMIDev.String()
			if got != tt.want {
				t.Errorf("IPMIDeviceInfo.String() = %q, want %q", got, tt.want)
			}
		})
	}
}
