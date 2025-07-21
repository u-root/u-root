// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package smbios

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"
)

var testbinary = "testdata/satellite_pro_l70_testdata.bin"

func checkError(got error, want error) bool {
	if got != nil && want != nil {
		if got.Error() == want.Error() {
			return true
		}
	}

	return errors.Is(got, want)
}

func TestParseSMBIOS(t *testing.T) {
	data, err := os.ReadFile(testbinary)
	if err != nil {
		t.Error(err)
	}
	datalen := len(data)
	readlen := 0
	for i := 0; datalen > i; i += readlen {
		_, rest, err := ParseTable(data)
		if err != nil {
			t.Error(err)
		}
		readlen = datalen - len(rest)
	}
}

func Test64Len(t *testing.T) {
	info, err := setupMockData()
	if err != nil {
		t.Errorf("error parsing info data: %v", err)
	}

	if info.Tables != nil {
		if info.Tables[0].Len() != 14 {
			t.Errorf("Wrong length: Got %d, want %d", info.Tables[0].Len(), 14)
		}
	}
}

func Test64String(t *testing.T) {
	tableString := `Handle 0x0000, DMI type 222, 14 bytes
OEM-specific Type
	Header and Data:
		DE 0E 00 00 01 99 00 03 10 01 20 02 30 03
	Strings:
		Memory Init Complete
		End of DXE Phase
		BIOS Boot Complete`

	info, err := setupMockData()
	if err != nil {
		t.Errorf("error parsing info data: %v", err)
	}

	if info.Tables != nil {
		if info.Tables[0].String() != tableString {
			t.Errorf("Wrong length: Got %s, want %s", info.Tables[0].String(), tableString)
		}
	}
}

func Test64MarshalBinary(t *testing.T) {
	tests := []struct {
		name  string
		table Table
		want  []byte
	}{
		{
			name: "GetRaw",
			table: Table{
				Header: Header{
					Type:   224,
					Length: 14,
					Handle: 0,
				},
				data:    []byte{222, 14, 0, 0, 1, 153, 0, 3, 16, 1, 32, 2, 48, 3},
				strings: []string{"Memory Init Complete", "End of DXE Phase", "BIOS Boot Complete"},
			},
			want: []byte{
				// Header and Data
				222, 14, 0, 0, 1, 153, 0, 3, 16, 1, 32, 2, 48, 3,
				// Strings
				77, 101, 109, 111, 114, 121, 32, 73, 110, 105, 116, 32, 67, 111, 109, 112, 108, 101, 116, 101, 0, // Memory Init Complete
				69, 110, 100, 32, 111, 102, 32, 68, 88, 69, 32, 80, 104, 97, 115, 101, 0, // End of DXE Phase
				66, 73, 79, 83, 32, 66, 111, 111, 116, 32, 67, 111, 109, 112, 108, 101, 116, 101, 0, //  BIOS Boot Complete
				0, // Table terminator
			},
		},
		{
			name: "GetRawNoString",
			table: Table{
				Header: Header{
					Type:   224,
					Length: 14,
					Handle: 0,
				},
				data: []byte{222, 14, 0, 0, 1, 153, 0, 3, 16, 1, 32, 2, 48, 3},
			},
			want: []byte{
				// Header and Data
				222, 14, 0, 0, 1, 153, 0, 3, 16, 1, 32, 2, 48, 3,
				// Strings
				0, // String terminator
				0, // Table terminator
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.table.MarshalBinary()
			if err != nil {
				t.Errorf("MarshalBinary returned error: %v", err)
			}
			if !bytes.Equal(got, tt.want) {
				t.Errorf("Wrong raw data: Got %v, want %v", got, tt.want)
			}
		})
	}
}

func Test64GetByteAt(t *testing.T) {
	testStruct := Table{
		Header: Header{
			Type:   TableTypeBIOSInfo,
			Length: 16,
			Handle: 0,
		},
		data:    []byte{1, 0, 0, 0, 213, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		strings: []string{"BIOS Boot Complete", "TestString #1"},
	}

	tests := []struct {
		name         string
		offset       int
		expectedByte uint8
		want         error
	}{
		{
			name:         "GetByteAt",
			offset:       0,
			expectedByte: 1,
			want:         nil,
		},
		{
			name:         "GetByteAt Wrong Offset",
			offset:       213,
			expectedByte: 0,
			want:         fmt.Errorf("invalid offset %d", 213),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultByte, err := testStruct.GetByteAt(tt.offset)
			if !checkError(err, tt.want) {
				t.Errorf("GetByteAt(): '%v', want '%v'", err, tt.want)
			}
			if resultByte != tt.expectedByte {
				t.Errorf("GetByteAt() = %x, want %x", resultByte, tt.expectedByte)
			}
		})
	}
}

func Test64GetBytesAt(t *testing.T) {
	testStruct := Table{
		Header: Header{
			Type:   TableTypeBIOSInfo,
			Length: 16,
			Handle: 0,
		},
		data:    []byte{1, 0, 0, 0, 213, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		strings: []string{"BIOS Boot Complete", "TestString #1"},
	}

	tests := []struct {
		name          string
		offset        int
		length        int
		expectedBytes []byte
		want          error
	}{
		{
			name:          "Get two bytes",
			offset:        0,
			length:        2,
			expectedBytes: []byte{1, 0},
			want:          nil,
		},
		{
			name:          "Wrong Offset",
			offset:        213,
			expectedBytes: []byte{},
			want:          fmt.Errorf("invalid offset 213"),
		},
		{
			name:          "Read out-of-bounds",
			offset:        7,
			length:        16,
			expectedBytes: []byte{},
			want:          fmt.Errorf("invalid offset 7"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultBytes, err := testStruct.GetBytesAt(tt.offset, tt.length)

			if !checkError(err, tt.want) {
				t.Errorf("GetBytesAt(): '%v', want '%v'", err, tt.want)
			}
			if !bytes.Equal(resultBytes, tt.expectedBytes) && err == nil {
				t.Errorf("GetBytesAt(): Wrong byte size, %x, want %x", resultBytes, tt.expectedBytes)
			}
		})
	}
}

func Test64GetWordAt(t *testing.T) {
	testStruct := Table{
		Header: Header{
			Type:   TableTypeBIOSInfo,
			Length: 16,
			Handle: 0,
		},
		data:    []byte{1, 0, 0, 0, 213, 0, 0, 11, 12, 0, 0, 0, 0, 0, 0},
		strings: []string{"BIOS Boot Complete", "TestString #1"},
	}

	tests := []struct {
		name          string
		offset        int
		expectedBytes uint16
		want          error
	}{
		{
			name:          "Get two bytes",
			offset:        0,
			expectedBytes: 1,
			want:          nil,
		},
		{
			name:          "Wrong Offset",
			offset:        213,
			expectedBytes: 0,
			want:          fmt.Errorf("invalid offset 213"),
		},
		{
			name:          "Read position 7",
			offset:        7,
			expectedBytes: 0xc0b,
			want:          nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultBytes, err := testStruct.GetWordAt(tt.offset)
			if !checkError(err, tt.want) {
				t.Errorf("GetBytesAt(): '%v', want '%v'", err, tt.want)
			}
			if resultBytes != tt.expectedBytes && err == nil {
				t.Errorf("GetBytesAt(): Wrong byte size, %x, want %x", resultBytes, tt.expectedBytes)
			}
		})
	}
}

func Test64GetDWordAt(t *testing.T) {
	testStruct := Table{
		Header: Header{
			Type:   TableTypeBIOSInfo,
			Length: 16,
			Handle: 0,
		},
		data:    []byte{1, 0, 0, 0, 213, 0, 0, 11, 12, 13, 14, 0, 0, 0, 0},
		strings: []string{"BIOS Boot Complete", "TestString #1"},
	}

	tests := []struct {
		name          string
		offset        int
		expectedBytes uint32
		want          error
	}{
		{
			name:          "Get two bytes",
			offset:        0,
			expectedBytes: 1,
			want:          nil,
		},
		{
			name:          "Wrong Offset",
			offset:        213,
			expectedBytes: 0,
			want:          fmt.Errorf("invalid offset 213"),
		},
		{
			name:          "Read position 7",
			offset:        7,
			expectedBytes: 0xe0d0c0b,
			want:          nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultBytes, err := testStruct.GetDWordAt(tt.offset)
			if !checkError(err, tt.want) {
				t.Errorf("GetBytesAt(): '%v', want '%v'", err, tt.want)
			}
			if resultBytes != tt.expectedBytes && err == nil {
				t.Errorf("GetBytesAt(): Wrong byte size, %x, want %x", resultBytes, tt.expectedBytes)
			}
		})
	}
}

func Test64GetQWordAt(t *testing.T) {
	testStruct := Table{
		Header: Header{
			Type:   TableTypeBIOSInfo,
			Length: 16,
			Handle: 0,
		},
		data:    []byte{1, 0, 0, 0, 213, 0, 0, 11, 12, 13, 14, 15, 16, 17, 18},
		strings: []string{"BIOS Boot Complete", "TestString #1"},
	}

	tests := []struct {
		name          string
		offset        int
		expectedBytes uint64
		want          error
	}{
		{
			name:          "Get two bytes",
			offset:        0,
			expectedBytes: 0xb0000d500000001,
			want:          nil,
		},
		{
			name:          "Wrong Offset",
			offset:        213,
			expectedBytes: 0,
			want:          fmt.Errorf("invalid offset 213"),
		},
		{
			name:          "Read position 7",
			offset:        7,
			expectedBytes: 0x1211100f0e0d0c0b,
			want:          nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultBytes, err := testStruct.GetQWordAt(tt.offset)
			if !checkError(err, tt.want) {
				t.Errorf("GetBytesAt(): '%v', want '%v'", err, tt.want)
			}
			if resultBytes != tt.expectedBytes && err == nil {
				t.Errorf("GetBytesAt(): Wrong byte size, %x, want %x", resultBytes, tt.expectedBytes)
			}
		})
	}
}

func TestKmgt(t *testing.T) {
	tests := []struct {
		name   string
		value  uint64
		expect string
	}{
		{
			name:   "Just bytes",
			value:  512,
			expect: "512 bytes",
		},
		{
			name:   "Two Kb",
			value:  2 * 1024,
			expect: "2 kB",
		},
		{
			name:   "512 MB",
			value:  512 * 1024 * 1024,
			expect: "512 MB",
		},
		{
			name:   "8 GB",
			value:  8 * 1024 * 1024 * 1024,
			expect: "8 GB",
		},
		{
			name:   "3 TB",
			value:  3 * 1024 * 1024 * 1024 * 1024,
			expect: "3 TB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if kmgt(tt.value) != tt.expect {
				t.Errorf("kgmt(): %v - want '%v'", kmgt(tt.value), tt.expect)
			}
		})
	}
}

func Test64GetStringAt(t *testing.T) {
	testStruct := Table{
		Header: Header{
			Type:   TableTypeBIOSInfo,
			Length: 16,
			Handle: 0,
		},
		data:    []byte{1, 0, 0, 0, 213, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		strings: []string{"BIOS Boot Complete", "TestString #1"},
	}

	tests := []struct {
		name           string
		offset         int
		expectedString string
	}{
		{
			name:           "Valid offset",
			offset:         0,
			expectedString: "BIOS Boot Complete",
		},
		{
			name:           "Not Specified",
			offset:         2,
			expectedString: "Not Specified",
		},
		{
			name:           "Bad Index",
			offset:         4,
			expectedString: "<BAD INDEX>",
		},
	}

	for _, tt := range tests {
		resultString, _ := testStruct.GetStringAt(tt.offset)
		if resultString != tt.expectedString {
			t.Errorf("GetStringAt(): %s, want '%s'", resultString, tt.expectedString)
		}
	}
}
