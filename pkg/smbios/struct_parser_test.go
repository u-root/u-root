// Copyright 2016-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"encoding/binary"
	"fmt"
	"strings"
	"testing"
)

type UnknownTypes struct {
	Table
	SupportedField   uint64
	UnsupportedField float32
}

func TestParseStructUnsupported(t *testing.T) {
	buffer := []byte{
		0x77,
		0xFF,
		0x00, 0x11,
		0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
		0x00, 0x01, 0x02, 0x03,
	}

	want := "unsupported type float32"

	table := Table{
		data: buffer,
	}

	UnknownType := &UnknownTypes{
		Table: table,
	}

	off, err := parseStruct(&table, 0, false, UnknownType)
	if err == nil {
		t.Errorf("TestParseStructUnsupported : parseStruct() = %d, '%v' want: %q", off, err, want)
	} else {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("TestParseStructUnsupported : parseStruct() = %d, '%v' want: %q", off, err, want)
		}
	}
}

func TestParseStructSupported(t *testing.T) {
	buffer := []byte{
		0x77,
		0xFF,
		0x00, 0x11,
		0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
	}

	table := Table{
		data: buffer,
	}

	UnknownType := &UnknownTypes{
		Table: table,
	}

	off, err := parseStruct(&table, 0, false, UnknownType)
	if err != nil {
		t.Errorf("TestParseStructUnsupported : parseStruct() = %d, '%v' want: 'nil'", off, err)
	}
}

func TestParseStructWithTPMDevice(t *testing.T) {
	tests := []struct {
		name     string
		buffer   []byte
		strings  []string
		complete bool
		want     TPMDevice
		wantErr  error
	}{
		{
			name: "Type43TPMDevice",
			buffer: []byte{
				0x2B,       // Type
				0xFF,       // Length
				0x00, 0x11, // Handle
				0x00, 0x00, 0x00, 0x00, // VendorID
				0x02,       // Major
				0x03,       // Minor
				0x01, 0x00, // FirmwareVersion1
				0x02, 0x00, // FirmwareVersion1
				0x00, 0x00, 0x00, 0x00, // FirmwareVersion2
				0x01,                   // String Index
				1 << 3,                 // Characteristics
				0x78, 0x56, 0x34, 0x12, // OEMDefined
			},
			strings:  []string{"Test TPM"},
			complete: false,
			want: TPMDevice{
				VendorID:         [4]byte{0x00, 0x00, 0x00, 0x00},
				MajorSpecVersion: 2,
				MinorSpecVersion: 3,
				FirmwareVersion1: 0x00020001,
				FirmwareVersion2: 0x00000000,
				Description:      "Test TPM",
				Characteristics:  TPMDeviceCharacteristicsFamilyConfigurableViaFirmwareUpdate,
				OEMDefined:       0x12345678,
			},
			wantErr: nil,
		},
		{
			name: "Type43TPMDevice Incomplete",
			buffer: []byte{
				0x2B,       // Type
				0xFF,       // Length
				0x00, 0x11, // Handle
				0x00, 0x00, 0x00, 0x00, // VendorID
				0x02,       // Major
				0x03,       // Minor
				0x01, 0x00, // FirmwareVersion1
				0x02, 0x00, // FirmwareVersion1
				0x00, 0x00, 0x00, 0x00, // FirmwareVersion2
				0x01,   // String Index
				1 << 3, // Characteristics
			},
			strings:  []string{"Test TPM"},
			complete: true,
			want: TPMDevice{
				VendorID:         [4]byte{0x00, 0x00, 0x00, 0x00},
				MajorSpecVersion: 2,
				MinorSpecVersion: 3,
				FirmwareVersion1: 0x00020001,
				FirmwareVersion2: 0x00000000,
				Description:      "Test TPM",
				Characteristics:  TPMDeviceCharacteristicsFamilyConfigurableViaFirmwareUpdate,
				OEMDefined:       0x12345678,
			},
			wantErr: fmt.Errorf("TPMDevice incomplete, got 8 of 9 fields"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := Table{
				data:    tt.buffer,
				strings: tt.strings,
			}
			TPMDev := &TPMDevice{
				Table: table,
			}

			// We need to modify tt.want with runtime data
			tt.want.Table = Table{
				Header: Header{
					Type:   TableType(tt.buffer[0]),
					Length: tt.buffer[1],
					Handle: binary.BigEndian.Uint16([]byte{tt.buffer[3], tt.buffer[2]}),
				},
				data:    tt.buffer,
				strings: tt.strings,
			}

			off, err := parseStruct(&table, 0, tt.complete, TPMDev)
			if err != tt.wantErr {
				if !strings.Contains(err.Error(), tt.wantErr.Error()) {
					t.Errorf("parseStruct() = %d, '%v' want '%v'", off, err, tt.wantErr)
				}
			}
			if tt.wantErr == nil {
				if TPMDev.VendorID != tt.want.VendorID {
					t.Errorf("parseStruct().VendorID = %q, want %q", TPMDev.VendorID, tt.want.VendorID)
				}

				if TPMDev.MajorSpecVersion != tt.want.MajorSpecVersion {
					t.Errorf("parseStruct().MajorSpecVersion = %q, want %q", TPMDev.MajorSpecVersion, tt.want.MajorSpecVersion)
				}

				if TPMDev.MinorSpecVersion != tt.want.MinorSpecVersion {
					t.Errorf("parseStruct().MinorSpecVersion = %q, want %q", TPMDev.MinorSpecVersion, tt.want.MinorSpecVersion)
				}

				if TPMDev.FirmwareVersion1 != tt.want.FirmwareVersion1 {
					t.Errorf("parseStruct().FirmwareVersion1 = %q, want %q", TPMDev.FirmwareVersion1, tt.want.FirmwareVersion1)
				}

				if TPMDev.FirmwareVersion2 != tt.want.FirmwareVersion2 {
					t.Errorf("parseStruct().FirmwareVersion2 = %q, want %q", TPMDev.FirmwareVersion2, tt.want.FirmwareVersion2)
				}

				if TPMDev.Description != tt.want.Description {
					t.Errorf("parseStruct().Description = %q, want %q", TPMDev.Description, tt.want.Description)
				}

				if TPMDev.Characteristics != tt.want.Characteristics {
					t.Errorf("parseStruct().Characteristics = %q, want %q", TPMDev.Characteristics, tt.want.Characteristics)
				}

				if TPMDev.OEMDefined != tt.want.OEMDefined {
					t.Errorf("parseStruct().OEMDefined = %q, want %q", TPMDev.OEMDefined, tt.want.OEMDefined)
				}
			}
		})
	}
}
