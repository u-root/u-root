// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package smbios

import "testing"

func defaultEntry32() *Entry32 {
	return &Entry32{
		Anchor:             [4]byte{95, 83, 77, 95},
		Checksum:           0x73,
		Length:             0x1F,
		SMBIOSMajorVersion: 0x01,
		SMBIOSMinorVersion: 0x01,
		StructMaxSize:      0x000E,
		Revision:           0x00,
		Reserved:           [5]byte{0x00, 0x00, 0x00, 0x00, 0x00},
		IntAnchor:          [5]byte{95, 68, 77, 73, 95},
		IntChecksum:        0x68,
		StructTableLength:  0x0000,
		StructTableAddr:    0x00000000,
		NumberOfStructs:    0x0000,
		BCDRevision:        0x00,
	}
}

func TestEntry32Marshall(t *testing.T) {
	for _, tt := range []struct {
		name  string
		entry *Entry32
	}{
		{
			name:  "Test valid Entry32",
			entry: defaultEntry32(),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.entry.MarshalBinary()
			if err != nil {
				t.Errorf("MarshalBinary(): %v", err)
			}
		})
	}
}

var validEntry32Bytes = []byte{95, 83, 77, 95, 115, 31, 1, 1, 14, 0, 0, 0, 0, 0, 0, 0, 95, 68, 77, 73, 95, 104, 0, 0, 0, 0, 0, 0, 0, 0, 0}

func TestEntry32Unmarshall(t *testing.T) {
	for _, tt := range []struct {
		name string
		data []byte
	}{
		{
			name: "Test valid data unmarshalling",
			data: validEntry32Bytes,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var e Entry32
			if err := e.UnmarshalBinary(tt.data); err != nil {
				t.Errorf("UnmarshalBinary(): %v", err)
			}
		})
	}
}

func defaultEntry64() *Entry64 {
	return &Entry64{
		Anchor:             [5]byte{95, 83, 77, 51, 95},
		Checksum:           0x5F,
		Length:             0x18,
		SMBIOSMajorVersion: 0x02,
		SMBIOSMinorVersion: 0x01,
		SMBIOSDocRev:       0x01,
		Revision:           0x00,
		Reserved:           0x00,
		StructMaxSize:      0xFFFFFFFF,
		StructTableAddr:    0xFFFFFFFFFFFFFFFF,
	}
}

func TestEntry64Marshall(t *testing.T) {
	for _, tt := range []struct {
		name  string
		entry *Entry64
	}{
		{
			name:  "Test valid Entry64",
			entry: defaultEntry64(),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.entry)
			data, err := tt.entry.MarshalBinary()
			if err != nil {
				t.Errorf("MarshalBinary(): %v", err)
			}
			t.Log(data)
		})
	}
}

var validEntry64Bytes = []byte{95, 83, 77, 51, 95, 95, 24, 2, 1, 1, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}

func TestEntry64Unmarshall(t *testing.T) {
	for _, tt := range []struct {
		name string
		data []byte
	}{
		{
			name: "Test valid data unmarshalling",
			data: validEntry64Bytes,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var e Entry64
			if err := e.UnmarshalBinary(tt.data); err != nil {
				t.Errorf("UnmarshalBinary(): %v", err)
			}
		})
	}
}

func TestParseEntry(t *testing.T) {
	for _, tt := range []struct {
		name string
		data []byte
	}{
		{
			name: "Test valid data32 unmarshalling",
			data: validEntry32Bytes,
		},
		{
			name: "Test valid data64 unmarshalling",
			data: validEntry64Bytes,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := ParseEntry(tt.data)
			if err != nil {
				t.Errorf("ParseEntry(): %v", err)
			}
		})
	}
}
