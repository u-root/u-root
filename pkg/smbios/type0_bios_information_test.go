// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"fmt"
	"testing"
)

var validType0BIOSInfoData = []byte{0, 26, 0, 0, 1, 2, 0, 0, 3, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

func validType0BIOSInfoRaw(t *testing.T) []byte {
	return joinBytesT(
		t,
		validType0BIOSInfoData,
		"MockVendor", 0,
		"1.0", 0,
		"01/01/2024", 0,
		0, // Table terminator
	)
}

func TestBIOSCharacteristicsString(t *testing.T) {
	tests := []struct {
		name string
		val  uint64
		want string
	}{
		{
			name: "Reserved",
			val:  0x1,
			want: "\t\tReserved",
		},
		{
			name: "Every Option",
			val:  0xFFFFFFFFFFFF,
			want: `		Reserved
		Reserved
		Unknown
		BIOS characteristics not supported
		ISA is supported
		MCA is supported
		EISA is supported
		PCI is supported
		PC Card (PCMCIA) is supported
		PNP is supported
		APM is supported
		BIOS is upgradeable
		BIOS shadowing is allowed
		VLB is supported
		ESCD support is available
		Boot from CD is supported
		Selectable boot is supported
		BIOS ROM is socketed
		Boot from PC Card (PCMCIA) is supported
		EDD is supported
		Japanese floppy for NEC 9800 1.2 MB is supported (int 13h)
		Japanese floppy for Toshiba 1.2 MB is supported (int 13h)
		5.25"/360 kB floppy services are supported (int 13h)
		5.25"/1.2 MB floppy services are supported (int 13h)
		3.5"/720 kB floppy services are supported (int 13h)
		3.5"/2.88 MB floppy services are supported (int 13h)
		Print screen service is supported (int 5h)
		8042 keyboard services are supported (int 9h)
		Serial services are supported (int 14h)
		Printer services are supported (int 17h)
		CGA/mono video services are supported (int 10h)
		NEC PC-98`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultString := BIOSCharacteristics(tt.val).String()
			if resultString != tt.want {
				t.Errorf("BIOSCharacteristics().String(): '%s', want '%s'", resultString, tt.want)
			}
		})
	}
}

func TestBIOSCharacteristicsExt1String(t *testing.T) {
	tests := []struct {
		name string
		val  uint8
		want string
	}{
		{
			name: "All options Ext1",
			val:  0xFF,
			want: `		ACPI is supported
		USB legacy is supported
		AGP is supported
		I2O boot is supported
		LS-120 boot is supported
		ATAPI Zip drive boot is supported
		IEEE 1394 boot is supported
		Smart battery is supported`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testVal := BIOSCharacteristicsExt1(tt.val)

			resultString := testVal.String()

			if resultString != tt.want {
				t.Errorf("BIOSCharacteristicsExt1().String(): '%s', want '%s'", resultString, tt.want)
			}
		})
	}
}

func TestBIOSCharacteristicsExt2String(t *testing.T) {
	tests := []struct {
		name string
		val  uint8
		want string
	}{
		{
			name: "All options Ex2",
			val:  0xFF,
			want: `		BIOS boot specification is supported
		Function key-initiated network boot is supported
		Targeted content distribution is supported
		UEFI is supported
		System is a virtual machine`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testVal := BIOSCharacteristicsExt2(tt.val)

			resultString := testVal.String()

			if resultString != tt.want {
				t.Errorf("BIOSCharacteristicsExt2().String(): '%s', want '%s'", resultString, tt.want)
			}
		})
	}
}

func TestBIOSInfoString(t *testing.T) {
	tests := []struct {
		name string
		val  BIOSInfo
		want string
	}{
		{
			name: "Valid BIOSInfo",
			val: BIOSInfo{
				Vendor:                                 "u-root",
				Version:                                "1.0",
				StartingAddressSegment:                 0x4,
				ReleaseDate:                            "2021/11/23",
				ROMSize:                                8,
				Characteristics:                        BIOSCharacteristics(0x8),
				CharacteristicsExt1:                    BIOSCharacteristicsExt1(0x4),
				CharacteristicsExt2:                    BIOSCharacteristicsExt2(0x2),
				SystemBIOSMajorRelease:                 1,
				SystemBIOSMinorRelease:                 2,
				EmbeddedControllerFirmwareMajorRelease: 3,
				EmbeddedControllerFirmwareMinorRelease: 1,
				ExtendedROMSize:                        0x10,
			},
			want: `Handle 0x0000, DMI type 0, 0 bytes
BIOS Information
	Vendor: u-root
	Version: 1.0
	Release Date: 2021/11/23
	Address: 0x00040
	Runtime Size: 1048512 bytes
	ROM Size: 576 kB
	Characteristics:
		BIOS characteristics not supported
		AGP is supported
		Function key-initiated network boot is supported`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.val.String() != tt.want {
				t.Errorf("BiosInfo().String(): '%s', want '%s'", tt.val.String(), tt.want)
			}
		})
	}
}

func TestGetROMSizeBytes(t *testing.T) {
	tests := []struct {
		name string
		val  BIOSInfo
		want uint64
	}{
		{
			name: "Rom Size 0xFF",
			val: BIOSInfo{
				ROMSize: 0xFF,
			},
			want: 0x1000000,
		},
		{
			name: "Rom Size 0xAB",
			val: BIOSInfo{
				ROMSize: 0xAB,
			},
			want: 0xAC0000,
		},
		{
			name: "Big Ext Size",
			val: BIOSInfo{
				ROMSize: 0xFF,
				Table: Table{
					data: []byte{
						0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
						0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
						0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
						0x1a,
					},
				},
				ExtendedROMSize: 0xFFFF,
			},
			want: 0x3FFF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			romSize := tt.val.GetROMSizeBytes()

			if romSize != tt.want {
				t.Errorf("BiosInfo().GetROMSizeBytes(): '%v', want '%v'", romSize, tt.want)
			}
		})
	}
}

func TestParseBIOSInfo(t *testing.T) {
	tests := []struct {
		name  string
		val   BIOSInfo
		table Table
		want  error
	}{
		{
			name: "Parse BIOS Info",
			val:  BIOSInfo{},
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
		},
		{
			name: "Length too short",
			val:  BIOSInfo{},
			table: Table{
				Header: Header{
					Type: TableTypeBIOSInfo,
				},
				data: []byte{},
			},
			want: fmt.Errorf("required fields missing"),
		},
		{
			name: "Error parsing data",
			val:  BIOSInfo{},
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
			want: fmt.Errorf("error parsing data"),
		},
		{
			name: "Length too short",
			val:  BIOSInfo{},
			table: Table{
				Header: Header{
					Type: TableTypeCacheInfo,
				},
				data: []byte{},
			},
			want: fmt.Errorf("invalid table type 7"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parseStruct := func(t *Table, off int, complete bool, sp interface{}) (int, error) {
				return 0, tt.want
			}
			_, err := parseBIOSInfo(parseStruct, &tt.table)

			if !checkError(err, tt.want) {
				t.Errorf("parseBIOSInfo(): '%v', want '%v'", err, tt.want)
			}
		})
	}
}
