// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"fmt"
	"reflect"
	"testing"
)

func TestBoardTypeString(t *testing.T) {
	tests := []struct {
		name string
		val  BoardType
		want string
	}{
		{
			name: "Unknown board type",
			val:  BoardTypeUnknown,
			want: "Unknown",
		},
		{
			name: "Other board type",
			val:  BoardTypeOther,
			want: "Other",
		},
		{
			name: "Server Blade board type",
			val:  BoardTypeServerBlade,
			want: "Server Blade",
		},
		{
			name: "System Management Module board type",
			val:  BoardTypeSystemManagementModule,
			want: "System Management Module",
		},
		{
			name: "Connectivity Switch board type",
			val:  BoardTypeConnectivitySwitch,
			want: "Connectivity Switch",
		},
		{
			name: "Processor Module board type",
			val:  BoardTypeProcessorModule,
			want: "Processor Module",
		},
		{
			name: "I/O Module board type",
			val:  BoardTypeIOModule,
			want: "I/O Module",
		},
		{
			name: "Memory Module board type",
			val:  BoardTypeMemoryModule,
			want: "Memory Module",
		},
		{
			name: "Daughter board board type",
			val:  BoardTypeDaughterBoard,
			want: "Daughter board",
		},
		{
			name: "Motherboard board type",
			val:  BoardTypeMotherboardIncludesProcessorMemoryAndIO,
			want: "Motherboard",
		},
		{
			name: "Processor/Memory Module board type",
			val:  BoardTypeProcessorMemoryModule,
			want: "Processor/Memory Module",
		},
		{
			name: "Processor/IO Module board type",
			val:  BoardTypeProcessorIOModule,
			want: "Processor/IO Module",
		},
		{
			name: "Interconnect board board type",
			val:  BoardTypeInterconnectBoard,
			want: "Interconnect board",
		},
		{
			name: "Not known board type",
			val:  BoardType(0x10),
			want: "0x10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.val.String() != tt.want {
				t.Errorf("BoardType().String(): '%s', want '%s'", tt.val.String(), tt.want)
			}
		})
	}
}

func TestBoardFeaturesString(t *testing.T) {
	tests := []struct {
		name string
		val  BoardFeatures
		want string
	}{
		{
			name: "All options",
			val:  BoardFeatures(0x1F),
			want: `		Board is a hosting board
		Board requires at least one daughter board
		Board is removable
		Board is replaceable
		Board is hot swappable`,
		},
		{
			name: "Unknown Option",
			val:  BoardFeatures(0x80),
			want: "		",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.val.String()

			if result != tt.want {
				t.Errorf("BoardFeatures().String(): '%s', want '%s'", result, tt.want)
			}
		})
	}
}

func TestBaseBoardInfoString(t *testing.T) {
	tests := []struct {
		name string
		val  BaseboardInfo
		want string
	}{
		{
			name: "Fully populated",
			val: BaseboardInfo{
				Manufacturer:                   "Astria Porta",
				Product:                        "Stargate",
				Version:                        "1",
				SerialNumber:                   "0a 0b 0c 0d 0e 0f 01 04",
				AssetTag:                       "0a",
				BoardFeatures:                  BoardFeaturesRequiresAtLeastOneDaughterBoard,
				LocationInChassis:              "Free floating",
				ChassisHandle:                  2, // for easy carrying around
				BoardType:                      BoardTypeUnknown,
				NumberOfContainedObjectHandles: 2,
				ContainedObjectHandles:         []uint16{0xABCD, 0xE7E7},
			},
			want: `Handle 0x0000, DMI type 0, 0 bytes
BIOS Information
	Manufacturer: Astria Porta
	Product Name: Stargate
	Version: 1
	Serial Number: 0a 0b 0c 0d 0e 0f 01 04
	Asset Tag: 0a
	Features:
		Board requires at least one daughter board
	Location In Chassis: Free floating
	Chassis Handle: 0x0002
	Type: Unknown
	Contained Object Handles: 2
		0xABCD
		0xE7E7`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.val.String()
			if result != tt.want {
				t.Errorf("BaseboardInfo().String(): '%s', want '%s'", result, tt.want)
			}
		})
	}
}

func TestParseBaseboardInfo(t *testing.T) {
	tests := []struct {
		name  string
		val   BaseboardInfo
		table Table
		want  error
	}{
		{
			name: "Invalid Type",
			val:  BaseboardInfo{},
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
			val:  BaseboardInfo{},
			table: Table{
				Header: Header{
					Type: TableTypeBaseboardInfo,
				},
				data: []byte{},
			},
			want: fmt.Errorf("required fields missing"),
		},
		{
			name: "Error parsing structure",
			val:  BaseboardInfo{},
			table: Table{
				Header: Header{
					Type: TableTypeBaseboardInfo,
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
			name: "Parse valid BaseboardInfo",
			val: BaseboardInfo{
				NumberOfContainedObjectHandles: 2,
			},
			table: Table{
				Header: Header{
					Type: TableTypeBaseboardInfo,
				},
				data: []byte{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
					0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
					0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
					0x1a,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parseStruct := func(t *Table, off int, complete bool, sp interface{}) (int, error) {
				return 0, tt.want
			}
			_, err := parseBaseboardInfo(parseStruct, &tt.table)

			if !checkError(err, tt.want) {
				t.Errorf("parseBaseboardInfo(): '%v', want '%v'", err, tt.want)
			}
		})
	}
}

func TestBaseboardInfoToTablePass(t *testing.T) {
	tests := []struct {
		name string
		bi   *BaseboardInfo
		want *Table
	}{
		{
			name: "Full of strings",
			bi: &BaseboardInfo{
				Header: Header{
					Type:   TableTypeBaseboardInfo,
					Length: 17,
					Handle: 0,
				},
				Manufacturer:                   "Manufacturer",
				Product:                        "Product",
				Version:                        "Version",
				SerialNumber:                   "1234-5678",
				AssetTag:                       "8765-4321",
				BoardFeatures:                  0,
				LocationInChassis:              "Location",
				ChassisHandle:                  0,
				BoardType:                      0,
				NumberOfContainedObjectHandles: 1,
				ContainedObjectHandles:         []uint16{10},
			},
			want: &Table{
				Header: Header{
					Type:   TableTypeBaseboardInfo,
					Length: 17,
					Handle: 0,
				},
				data: []byte{
					2, 17, 0, 0, // Header
					1, 2, 3, 4, 5, // string number
					0,    // BoardFeatures
					6,    // string number
					0, 0, // ChassisHandle
					0,     // BoardType
					1,     // NumberOfContainedObjectHandles
					10, 0, // ContainedObjectHandles
				},
				strings: []string{"Manufacturer", "Product", "Version", "1234-5678", "8765-4321", "Location"},
			},
		},
		{
			name: "Have empty strings",
			bi: &BaseboardInfo{
				Header: Header{
					Type:   TableTypeBaseboardInfo,
					Length: 17,
					Handle: 0,
				},
				Manufacturer:                   "Manufacturer",
				Product:                        "",
				Version:                        "Version",
				SerialNumber:                   "",
				AssetTag:                       "8765-4321",
				BoardFeatures:                  0,
				LocationInChassis:              "Location",
				ChassisHandle:                  0,
				BoardType:                      0,
				NumberOfContainedObjectHandles: 1,
				ContainedObjectHandles:         []uint16{10},
			},
			want: &Table{
				Header: Header{
					Type:   TableTypeBaseboardInfo,
					Length: 17,
					Handle: 0,
				},
				data: []byte{
					2, 17, 0, 0, // Header
					1, 0, 2, 0, 3, // string number
					0,    // BoardFeatures
					4,    // string number
					0, 0, // ChassisHandle
					0,     // BoardType
					1,     // NumberOfContainedObjectHandles
					10, 0, // ContainedObjectHandles
				},
				strings: []string{"Manufacturer", "Version", "8765-4321", "Location"},
			},
		},
	}

	for _, tt := range tests {
		got, err := tt.bi.toTable()
		if err != nil {
			t.Errorf("toTable() should pass but return error: %v", err)
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("toTable(): '%v', want '%v'", got, tt.want)
		}
	}
}

func TestBaseboardInfoToTableFail(t *testing.T) {
	bi := &BaseboardInfo{
		Header: Header{
			Type:   TableTypeBaseboardInfo,
			Length: 17,
			Handle: 0,
		},
		NumberOfContainedObjectHandles: 2, // Wrong NumberOfContainedObjectHandles
		ContainedObjectHandles:         []uint16{10},
	}

	_, err := bi.toTable()

	if err == nil {
		t.Fatalf("toTable() should fail but pass")
	}
}
