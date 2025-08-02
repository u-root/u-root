// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"fmt"
	"testing"
)

func TestCacheSizeBytes2Or1(t *testing.T) {
	tests := []struct {
		name  string
		size1 uint16
		size2 uint32
		want  uint64
	}{
		{
			name:  "No high bit set",
			size1: 0x1234,
			size2: 0x12345678,
			want:  0x48D159E000,
		},
		{
			name:  "High bit set",
			size1: 0x1234,
			size2: 0x80023456,
			want:  0x234560000,
		},
		{
			name:  "size2 zero",
			size1: 0x1234,
			size2: 0x80000000,
			want:  0x48D000,
		},
		{
			name:  "Zero",
			size1: 0x8000,
			size2: 0x80000000,
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			size := cacheSizeBytes2Or1(tt.size1, tt.size2)
			if size != tt.want {
				t.Errorf("%q failed. Got: %q, Want: %q", tt.name, size, tt.want)
			}
		})
	}
}

func TestCacheInfoString(t *testing.T) {
	tests := []struct {
		name string
		val  CacheInfo
		want string
	}{
		{
			name: "Full details",
			val: CacheInfo{
				Table: Table{
					data: []byte{
						0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
						0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
					},
				},
				SocketDesignation:   "",
				Configuration:       0x03,
				MaximumSize:         0x100,
				InstalledSize:       0x3F,
				SupportedSRAMType:   CacheSRAMTypePipelineBurst,
				CurrentSRAMType:     CacheSRAMTypeOther,
				Speed:               0x4,
				ErrorCorrectionType: CacheErrorCorrectionTypeParity,
				SystemType:          CacheSystemTypeUnified,
				Associativity:       CacheAssociativity16waySetAssociative,
				MaximumSize2:        0x200,
				InstalledSize2:      0x00,
			},
			want: `Handle 0x0000, DMI type 0, 0 bytes
BIOS Information
	Socket Designation: 
	Configuration: Disabled, Not Socketed, Level 4
	Operational Mode: Write Through
	Location: Internal
	Installed Size: 63 kB
	Maximum Size: 512 kB
	Supported SRAM Types:
		Pipeline Burst
	Installed SRAM Type: Other
	Speed: 4 ns
	Error Correction Type: Parity
	System Type: Unified
	Associativity: 16-way Set-associative`,
		},
		{
			name: "More details",
			val: CacheInfo{
				Table: Table{
					data: []byte{
						0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
						0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
					},
				},
				SocketDesignation:   "",
				Configuration:       0x3A8,
				MaximumSize:         0x100,
				InstalledSize:       0x3F,
				SupportedSRAMType:   CacheSRAMTypePipelineBurst,
				CurrentSRAMType:     CacheSRAMTypeOther,
				Speed:               0x4,
				ErrorCorrectionType: CacheErrorCorrectionTypeParity,
				SystemType:          CacheSystemTypeUnified,
				Associativity:       CacheAssociativity16waySetAssociative,
				MaximumSize2:        0x200,
				InstalledSize2:      0x00,
			},
			want: `Handle 0x0000, DMI type 0, 0 bytes
BIOS Information
	Socket Designation: 
	Configuration: Enabled, Socketed, Level 1
	Operational Mode: Unknown
	Location: External
	Installed Size: 63 kB
	Maximum Size: 512 kB
	Supported SRAM Types:
		Pipeline Burst
	Installed SRAM Type: Other
	Speed: 4 ns
	Error Correction Type: Parity
	System Type: Unified
	Associativity: 16-way Set-associative`,
		},
		{
			name: "More details",
			val: CacheInfo{
				Table: Table{
					data: []byte{
						0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
						0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
					},
				},
				SocketDesignation:   "",
				Configuration:       0x2CA,
				MaximumSize:         0x100,
				InstalledSize:       0x3F,
				SupportedSRAMType:   CacheSRAMTypePipelineBurst,
				CurrentSRAMType:     CacheSRAMTypeOther,
				Speed:               0x4,
				ErrorCorrectionType: CacheErrorCorrectionTypeParity,
				SystemType:          CacheSystemTypeUnified,
				Associativity:       CacheAssociativity16waySetAssociative,
				MaximumSize2:        0x200,
				InstalledSize2:      0x00,
			},
			want: `Handle 0x0000, DMI type 0, 0 bytes
BIOS Information
	Socket Designation: 
	Configuration: Enabled, Socketed, Level 3
	Operational Mode: Varies With Memory Address
	Location: Reserved
	Installed Size: 63 kB
	Maximum Size: 512 kB
	Supported SRAM Types:
		Pipeline Burst
	Installed SRAM Type: Other
	Speed: 4 ns
	Error Correction Type: Parity
	System Type: Unified
	Associativity: 16-way Set-associative`,
		},
		{
			name: "More details",
			val: CacheInfo{
				Table: Table{
					data: []byte{
						0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
						0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
					},
				},
				SocketDesignation:   "",
				Configuration:       0x1EA,
				MaximumSize:         0x100,
				InstalledSize:       0x3F,
				SupportedSRAMType:   CacheSRAMTypePipelineBurst,
				CurrentSRAMType:     CacheSRAMTypeOther,
				Speed:               0x4,
				ErrorCorrectionType: CacheErrorCorrectionTypeParity,
				SystemType:          CacheSystemTypeUnified,
				Associativity:       CacheAssociativity16waySetAssociative,
				MaximumSize2:        0x200,
				InstalledSize2:      0x00,
			},
			want: `Handle 0x0000, DMI type 0, 0 bytes
BIOS Information
	Socket Designation: 
	Configuration: Enabled, Socketed, Level 3
	Operational Mode: Write Back
	Location: Unknown
	Installed Size: 63 kB
	Maximum Size: 512 kB
	Supported SRAM Types:
		Pipeline Burst
	Installed SRAM Type: Other
	Speed: 4 ns
	Error Correction Type: Parity
	System Type: Unified
	Associativity: 16-way Set-associative`,
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

func TestParseInfoCache(t *testing.T) {
	tests := []struct {
		name  string
		val   CacheInfo
		table Table
		want  error
	}{
		{
			name: "Invalid Type",
			val:  CacheInfo{},
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
			val:  CacheInfo{},
			table: Table{
				Header: Header{
					Type: TableTypeCacheInfo,
				},
				data: []byte{},
			},
			want: fmt.Errorf("required fields missing"),
		},
		{
			name: "Error parsing structure",
			val:  CacheInfo{},
			table: Table{
				Header: Header{
					Type: TableTypeCacheInfo,
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
			name: "Parse valid CacheInfo",
			val:  CacheInfo{},
			table: Table{
				Header: Header{
					Type: TableTypeCacheInfo,
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
				return 0, tt.want
			}
			_, err := parseCacheInfo(parseStruct, &tt.table)

			if !checkError(err, tt.want) {
				t.Errorf("%q failed. Got: %q, Want: %q", tt.name, err, tt.want)
			}
		})
	}
}
