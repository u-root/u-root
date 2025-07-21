// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"fmt"
	"reflect"
	"testing"
)

var validType1SystemInfoData = []byte{1, 27, 1, 0, 1, 2, 3, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6, 5, 6}

func validType1SystemInfoRaw(t *testing.T) []byte {
	return joinBytesT(
		t,
		validType1SystemInfoData,
		"MockManufacturer", 0,
		"MockProductName", 0,
		"MockVersion", 0,
		"MockSerialNumber", 0,
		"MockSKUNumber", 0,
		"MockFamily", 0,
		0, // Table terminator
	)
}

func TestSystemInfoString(t *testing.T) {
	tests := []struct {
		name string
		val  SystemInfo
		want string
	}{
		{
			name: "All Infos provided",
			val: SystemInfo{
				Header: Header{
					Type:   TableTypeSystemInfo,
					Length: 27,
					Handle: 0,
				},
				Manufacturer: "u-root testing",
				ProductName:  "Illusion",
				Version:      "1.0",
				SerialNumber: "UR00T1234",
				UUID:         UUID{0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf},
				SKUNumber:    "3a",
				Family:       "UR00T1234",
			},
			want: `Handle 0x0000, DMI type 1, 27 bytes
System Information
	Manufacturer: u-root testing
	Product Name: Illusion
	Version: 1.0
	Serial Number: UR00T1234
	UUID: 03020100-0504-0706-0809-0a0b0c0d0e0f
	Wake-up Type: Reserved
	SKU Number: 3a
	Family: UR00T1234`,
		},
		{
			name: "UUID not present",
			val: SystemInfo{
				Header: Header{
					Type:   TableTypeSystemInfo,
					Length: 14,
					Handle: 0,
				},
				Manufacturer: "u-root testing",
				ProductName:  "Illusion",
				Version:      "1.0",
				SerialNumber: "UR00T1234",
				UUID:         UUID{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
				SKUNumber:    "3a",
				Family:       "UR00T1234",
			},
			want: `Handle 0x0000, DMI type 1, 14 bytes
System Information
	Manufacturer: u-root testing
	Product Name: Illusion
	Version: 1.0
	Serial Number: UR00T1234
	UUID: Not Present
	Wake-up Type: Reserved`,
		},
		{
			name: "UUID not present",
			val: SystemInfo{
				Header: Header{
					Type:   TableTypeSystemInfo,
					Length: 14,
					Handle: 0,
				},
				Manufacturer: "u-root testing",
				ProductName:  "Illusion",
				Version:      "1.0",
				SerialNumber: "UR00T1234",
				UUID:         UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
				SKUNumber:    "3a",
				Family:       "UR00T1234",
			},
			want: `Handle 0x0000, DMI type 1, 14 bytes
System Information
	Manufacturer: u-root testing
	Product Name: Illusion
	Version: 1.0
	Serial Number: UR00T1234
	UUID: Not Settable
	Wake-up Type: Reserved`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.val.String()

			if result != tt.want {
				t.Errorf("SystemInfo().String(): '%s', want '%s'", result, tt.want)
			}
		})
	}
}

func TestUUIDParseField(t *testing.T) {
	tests := []struct {
		name string
		val  Table
		want string
	}{
		{
			name: "Valid UUID",
			val: Table{
				data: []byte{
					0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03,
					0x00, 0x01, 0x02, 0x03, 0x00, 0x01, 0x02, 0x03,
				},
			},
			want: "03020100-0100-0302-0001-020300010203",
		},
	}

	for _, tt := range tests {
		uuid := UUID([16]byte{})
		_, err := uuid.ParseField(&tt.val, 0)
		if err != nil {
			t.Errorf("ParseField(): '%v', want nil", err)
		}
		if uuid.String() != tt.want {
			t.Errorf("ParseField(): '%s', want '%s'", uuid.String(), tt.want)
		}
	}
}

func TestParseSystemInfo(t *testing.T) {
	tests := []struct {
		name  string
		val   SystemInfo
		table Table
		want  error
	}{
		{
			name: "Invalid Type",
			val:  SystemInfo{},
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
			name: "Invalid Type",
			val:  SystemInfo{},
			table: Table{
				Header: Header{
					Type: TableTypeSystemInfo,
				},
				data: []byte{},
			},
			want: fmt.Errorf("required fields missing"),
		},
		{
			name: "Parse valid SystemInfo",
			val:  SystemInfo{},
			table: Table{
				Header: Header{
					Type: TableTypeSystemInfo,
				},
				data: []byte{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
					0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
					0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
					0x1a,
				},
			},
		},
		{
			name: "Parse valid SystemInfo",
			val:  SystemInfo{},
			table: Table{
				Header: Header{
					Type: TableTypeSystemInfo,
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parseStruct := func(t *Table, off int, complete bool, sp interface{}) (int, error) {
				return 0, tt.want
			}
			_, err := parseSystemInfo(parseStruct, &tt.table)

			if !checkError(err, tt.want) {
				t.Errorf("parseSystemInfo(): '%v', want '%v'", err, tt.want)
			}
		})
	}
}

func TestSystemInfoToTable(t *testing.T) {
	tests := []struct {
		name string
		si   *SystemInfo
		want *Table
	}{
		{
			name: "Full of strings",
			si: &SystemInfo{
				Header: Header{
					Type:   TableTypeSystemInfo,
					Length: 27,
					Handle: 0,
				},
				Manufacturer: "Manufacturer",
				ProductName:  "ProductName",
				Version:      "Version",
				SerialNumber: "SerialNumber",
				UUID:         [16]byte{0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3},
				WakeupType:   WakeupTypeReserved,
				SKUNumber:    "SKUNumber",
				Family:       "Family",
			},
			want: &Table{
				Header: Header{
					Type:   TableTypeSystemInfo,
					Length: 27,
					Handle: 0,
				},
				data: []byte{
					1, 27, 0, 0, // Header
					1, 2, 3, 4, // string number
					0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, // UUID
					0,    // WakeupType
					5, 6, // string number
				},
				strings: []string{"Manufacturer", "ProductName", "Version", "SerialNumber", "SKUNumber", "Family"},
			},
		},
		{
			name: "Have empty strings",
			si: &SystemInfo{
				Header: Header{
					Type:   TableTypeSystemInfo,
					Length: 27,
					Handle: 0,
				},
				Manufacturer: "Manufacturer",
				ProductName:  "",
				Version:      "Version",
				SerialNumber: "",
				UUID:         [16]byte{0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3},
				WakeupType:   WakeupTypeReserved,
				SKUNumber:    "SKUNumber",
				Family:       "",
			},
			want: &Table{
				Header: Header{
					Type:   TableTypeSystemInfo,
					Length: 27,
					Handle: 0,
				},
				data: []byte{
					1, 27, 0, 0, // Header
					1, 0, 2, 0, // string number
					0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, // UUID
					0,    // WakeupType
					3, 0, // string number
				},
				strings: []string{"Manufacturer", "Version", "SKUNumber"},
			},
		},
	}

	for _, tt := range tests {
		got, err := tt.si.toTable()
		if err != nil {
			t.Errorf("toTable() should pass but return error: %v", err)
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("toTable(): '%v', want '%v'", got, tt.want)
		}
	}
}
