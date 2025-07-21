// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"fmt"
	"testing"
)

func TestChassisInfoString(t *testing.T) {
	tests := []struct {
		name string
		val  ChassisInfo
		want string
	}{
		{
			name: "Full Information",
			val: ChassisInfo{
				Table: Table{
					data: []byte{
						0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
						0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
					},
				},
				Manufacturer:                  "The Ancients",
				Type:                          ChassisTypeAllInOne,
				Version:                       "One",
				SerialNumber:                  "TheAncients-01",
				AssetTagNumber:                "Two",
				BootupState:                   ChassisStateSafe,
				PowerSupplyState:              ChassisStateSafe,
				ThermalState:                  ChassisStateNonrecoverable,
				SecurityStatus:                ChassisSecurityStatusUnknown,
				OEMInfo:                       0xABCD0123,
				Height:                        3,
				NumberOfPowerCords:            1,
				ContainedElementCount:         2,
				ContainedElementsRecordLength: 2,
				ContainedElements: []ChassisContainedElement{
					{
						Type: ChassisElementType(8),
						Min:  3,
						Max:  11,
					},
					{
						Type: ChassisElementType(0),
						Min:  0,
						Max:  1,
					},
				},
				SKUNumber: "Four",
			},
			want: `Handle 0x0000, DMI type 0, 0 bytes
BIOS Information
	Manufacturer: The Ancients
	Type: All In One
	Lock: Not Present
	Version: One
	Serial Number: TheAncients-01
	Asset Tag: Two
	Boot-up State: Safe
	Power Supply State: Safe
	Thermal State: Non-recoverable
	Security Status: Unknown
	OEM Information: 0xABCD0123
	Height: 3 U
	Number Of Power Cords: 1
	Contained Elements: 2
		Memory Module 3-11
		0x0 0-1
	SKU Number: Four`,
		}, {
			name: "Minimal Information",
			val: ChassisInfo{
				Table: Table{
					data: []byte{},
				},
				Manufacturer:   "The Ancients",
				Type:           ChassisTypeTower,
				Version:        "Two",
				SerialNumber:   "TheAncients-02",
				AssetTagNumber: "Three",
			},
			want: `Handle 0x0000, DMI type 0, 0 bytes
BIOS Information
	Manufacturer: The Ancients
	Type: Tower
	Lock: Not Present
	Version: Two
	Serial Number: TheAncients-02
	Asset Tag: Three`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.val.String()

			if result != tt.want {
				t.Errorf("ChassisInfo().String(): %v, want '%v'", result, tt.want)
			}
		})
	}
}

func TestChassisTypeString(t *testing.T) {
	testResults := []string{
		"Other",
		"Unknown",
		"Desktop",
		"Low Profile Desktop",
		"Pizza Box",
		"Mini Tower",
		"Tower",
		"Portable",
		"Laptop",
		"Notebook",
		"Hand Held",
		"Docking Station",
		"All In One",
		"Sub Notebook",
		"Space-saving",
		"Lunch Box",
		"Main Server Chassis",
		"Expansion Chassis",
		"Sub Chassis",
		"Bus Expansion Chassis",
		"Peripheral Chassis",
		"RAID Chassis",
		"Rack Mount Chassis",
		"Sealed-case PC",
		"Multi-system",
		"CompactPCI",
		"AdvancedTCA",
		"Blade",
		"Blade Chassis",
		"Tablet",
		"Convertible",
		"Detachable",
		"IoT Gateway",
		"Embedded PC",
		"Mini PC",
		"Stick PC",
	}

	for id := range testResults {
		val := ChassisType(id + 1)

		if val.String() != testResults[id] {
			t.Errorf("ChassisType().String(): '%s', want '%s'", val.String(), testResults[id])
		}
	}
}

func TestParseChassisInfo(t *testing.T) {
	tests := []struct {
		name  string
		val   *ChassisInfo
		table Table
		want  error
	}{
		{
			name: "Invalid Type",
			val:  &ChassisInfo{},
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
			val:  &ChassisInfo{},
			table: Table{
				Header: Header{
					Type: TableTypeChassisInfo,
				},
				data: []byte{},
			},
			want: fmt.Errorf("required fields missing"),
		},
		{
			name: "Error parsing structure",
			val:  &ChassisInfo{},
			table: Table{
				Header: Header{
					Type: TableTypeChassisInfo,
				},
				data: []byte{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
					0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
					0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
					0x1a,
				},
			},
			want: fmt.Errorf("error parsing structure"),
		},
		{
			name: "Parse valid SystemInfo",
			val: &ChassisInfo{
				Table: Table{
					Header: Header{
						Type: TableTypeChassisInfo,
					},
					data: []byte{
						0x7, 0x01, 0x02, 0x07, 0x04, 0x05, 0x06, 0x07,
						0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
						0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
						0x1a, 0x7, 0x01, 0x02, 0x07, 0x04, 0x05, 0x06, 0x07,
						0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
						0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
						0x1a,
					},
				},
				ContainedElementCount:         7,
				ContainedElementsRecordLength: 0x10,
			},
			table: Table{
				Header: Header{
					Type: TableTypeChassisInfo,
				},
				data: []byte{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
					0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
					0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
					0x1a,
				},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parseStruct := func(t *Table, off int, complete bool, sp interface{}) (int, error) {
				return len(tt.val.data), tt.want
			}
			_, err := parseChassisInfo(parseStruct, &tt.table)

			if !checkError(err, tt.want) {
				t.Errorf("parseChassisInfo(): '%v', want '%v'", err, tt.want)
			}
		})
	}
}
