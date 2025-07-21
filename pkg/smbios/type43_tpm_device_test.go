// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"fmt"
	"testing"
)

func TestTPMDeviceString(t *testing.T) {
	tests := []struct {
		name string
		val  TPMDevice
		want string
	}{
		{
			name: "Infineon TPM",
			val: TPMDevice{
				VendorID:         TPMDeviceVendorID{0x0, 'X', 'F', 'I'},
				MajorSpecVersion: 1,
				MinorSpecVersion: 7,
				FirmwareVersion1: 2,
				FirmwareVersion2: 3,
				Description:      "Test TPM",
				Characteristics:  TPMDeviceCharacteristics(8),
				OEMDefined:       2,
			},
			want: `Handle 0x0000, DMI type 0, 0 bytes
BIOS Information
	Vendor ID: IFX
	Specification Version: 1.7
	Firmware Revision: 0.0
	Description: Test TPM
	Characteristics:
		Family configurable via firmware update
	OEM-specific Info: 0x00000002`,
		},
		{
			name: "Random TPM",
			val: TPMDevice{
				VendorID:         TPMDeviceVendorID{'A', 'B', 'C', 'D'},
				MajorSpecVersion: 2,
				MinorSpecVersion: 9,
				FirmwareVersion1: 1,
				FirmwareVersion2: 9,
				Description:      "Test TPM",
				Characteristics:  TPMDeviceCharacteristics(16),
				OEMDefined:       2,
			},
			want: `Handle 0x0000, DMI type 0, 0 bytes
BIOS Information
	Vendor ID: ABCD
	Specification Version: 2.9
	Firmware Revision: 0.1
	Description: Test TPM
	Characteristics:
		Family configurable via platform software support
	OEM-specific Info: 0x00000002`,
		},
		{
			name: "Random TPM #2",
			val: TPMDevice{
				VendorID:         TPMDeviceVendorID{'A', 'B', 'C', 'D'},
				MajorSpecVersion: 2,
				MinorSpecVersion: 9,
				FirmwareVersion1: 1,
				FirmwareVersion2: 9,
				Description:      "Test TPM",
				Characteristics:  TPMDeviceCharacteristics(32),
				OEMDefined:       2,
			},
			want: `Handle 0x0000, DMI type 0, 0 bytes
BIOS Information
	Vendor ID: ABCD
	Specification Version: 2.9
	Firmware Revision: 0.1
	Description: Test TPM
	Characteristics:
		Family configurable via OEM proprietary mechanism
	OEM-specific Info: 0x00000002`,
		},
		{
			name: "Random TPM #2",
			val: TPMDevice{
				VendorID:         TPMDeviceVendorID{'A', 'B', 'C', 'D'},
				MajorSpecVersion: 2,
				MinorSpecVersion: 9,
				FirmwareVersion1: 1,
				FirmwareVersion2: 9,
				Description:      "Test TPM",
				Characteristics:  TPMDeviceCharacteristics(4),
				OEMDefined:       2,
			},
			want: `Handle 0x0000, DMI type 0, 0 bytes
BIOS Information
	Vendor ID: ABCD
	Specification Version: 2.9
	Firmware Revision: 0.1
	Description: Test TPM
	Characteristics:
		TPM Device characteristics not supported
	OEM-specific Info: 0x00000002`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.val.String()

			if result != tt.want {
				t.Errorf("%q failed. Got: %q, Want: %q", tt.name, result, tt.want)
			}
		})
	}
}

func TestNewTPMDevice(t *testing.T) {
	tests := []struct {
		name  string
		val   TPMDevice
		table Table
		want  error
	}{
		{
			name: "Invalid Type",
			val:  TPMDevice{},
			table: Table{
				Header: Header{
					Type: TableTypeBIOSInfo,
				},
				data: []byte{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
					0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
					0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
					0x1a,
				},
			},
			want: fmt.Errorf("invalid table type 0"),
		},
		{
			name: "Required fields are missing",
			val:  TPMDevice{},
			table: Table{
				Header: Header{
					Type: TableTypeTPMDevice,
				},
				data: []byte{},
			},
			want: fmt.Errorf("required fields missing"),
		},
		{
			name: "Error parsing structure",
			val:  TPMDevice{},
			table: Table{
				Header: Header{
					Type: TableTypeTPMDevice,
				},
				data: []byte{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
					0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
					0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
					0x1a, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
					0x1a,
				},
			},
			want: fmt.Errorf("error parsing structure"),
		},
		{
			name: "Parse valid TPMDevice",
			val:  TPMDevice{},
			table: Table{
				Header: Header{
					Type: TableTypeTPMDevice,
				},
				data: []byte{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
					0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
					0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
					0x1a, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
					0x1a,
				},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parseStruct := func(t *Table, off int, complete bool, sp interface{}) (int, error) {
				return 0, tt.want
			}
			_, err := newTPMDevice(parseStruct, &tt.table)

			if !checkError(err, tt.want) {
				t.Errorf("%q failed. Got: %q, Want: %q", tt.name, err, tt.want)
			}
		})
	}
}
