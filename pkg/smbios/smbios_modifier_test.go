// Copyright 2016-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	fakeMemFile = "/tmp/fakeMemFile"
)

func joinBytesT(t *testing.T, args ...any) []byte {
	b, err := joinBytes(args...)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func defaultMock64Memory(t *testing.T) []byte {
	return joinBytesT(
		t,
		mockEntry64Raw(t, uint32(len(tablesRaw(t))), 24),
		tablesRaw(t),
	)
}

// tablesRaw returns the actual raw bytes of tables struct
func tablesRaw(t *testing.T) []byte {
	return joinBytesT(
		t,
		validType0BIOSInfoRaw(t),
		validType1SystemInfoRaw(t),
		validEndOfTableRaw(t),
	)
}

func makeMemFile(t *testing.T, data []byte) *os.File {
	t.Helper()
	file, err := os.Create(fakeMemFile)
	if err != nil {
		t.Fatalf("Failed to create fake mem file: %v", err)
	}

	if _, err = file.Write(data); err != nil {
		t.Fatalf("Failed to write fake mem file: %v", err)
	}
	return file
}

func mockEntry64Raw(t *testing.T, structMaxSize uint32, tableAddr uint64) []byte {
	data, err := mockEntry64(t, structMaxSize, tableAddr).MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func mockEntry32Raw(t *testing.T, maxSize uint16, tableLength uint16, tableAddr uint32, numberOfStructs uint16) []byte {
	data, err := mockEntry32(t, maxSize, tableLength, tableAddr, numberOfStructs).MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// mockEntry64 creates a valid Entry64 with custom attributes
func mockEntry64(t *testing.T, structMaxSize uint32, tableAddr uint64) *Entry64 {
	t.Helper()
	e64 := defaultEntry64()
	e64.StructMaxSize = structMaxSize
	e64.StructTableAddr = tableAddr

	// Recalculate checksum
	data, err := e64.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	e64.UnmarshalBinary(data)
	return e64
}

// mockEntry32 creates a valid Entry32 with custom attributes
func mockEntry32(t *testing.T, maxSize uint16, tableLength uint16, tableAddr uint32, numberOfStructs uint16) *Entry32 {
	t.Helper()
	e32 := defaultEntry32()
	e32.StructMaxSize = maxSize
	e32.StructTableLength = tableLength
	e32.StructTableAddr = tableAddr
	e32.NumberOfStructs = numberOfStructs

	// Recalculate checksum
	data, err := e32.MarshalBinary()
	if err != nil {
		panic(err)
	}
	e32.UnmarshalBinary(data)
	return e32
}

func newMock64Modifier(t *testing.T) *Modifier {
	t.Helper()
	memFile := makeMemFile(t, defaultMock64Memory(t))

	getMemFileMock := func() (*os.File, error) { return memFile, nil }
	smbiosBaseMock := func() (int64, int64, error) { return 0, 24, nil }

	m, err := newModifier(getMemFileMock, smbiosBaseMock)
	if err != nil {
		t.Fatalf("Failed to initialize modifier: %v", err)
	}

	return m
}

func readFileFromStart(f *os.File, size int64) ([]byte, error) {
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	buf := make([]byte, size)
	if _, err := io.ReadFull(f, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

func TestModifySystemInfo(t *testing.T) {
	tests := []struct {
		name                                             string
		manufacturer, productName, version, serialNumber string
		want                                             []byte
	}{
		{
			name:         "Default",
			manufacturer: "NewManufacturer",
			productName:  "NewProductName",
			version:      "NewVersion",
			serialNumber: "NewSerialNumber",
			want: joinBytesT(t,
				mockEntry64Raw(t, 170, 24),
				validType0BIOSInfoRaw(t),
				validType1SystemInfoData,
				"NewManufacturer", 0,
				"NewProductName", 0,
				"NewVersion", 0,
				"NewSerialNumber", 0,
				"MockSKUNumber", 0,
				"MockFamily", 0,
				0, // Table terminator
				validEndOfTableRaw(t),
			),
		},
		{
			name:         "SetOneStringToEmpty",
			manufacturer: "NewManufacturer",
			productName:  "NewProductName",
			version:      "",
			serialNumber: "NewSerialNumber",
			want: joinBytesT(t,
				mockEntry64Raw(t, 159, 24),
				validType0BIOSInfoRaw(t),
				[]byte{1, 27, 1, 0, 1, 2, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6, 4, 5}, // Structure data, the number of string skip the one for "version"
				"NewManufacturer", 0,
				"NewProductName", 0,
				"NewSerialNumber", 0,
				"MockSKUNumber", 0,
				"MockFamily", 0,
				0, // Table terminator
				validEndOfTableRaw(t),
			),
		},
		{
			name:         "NewStringsHaveDifferentLength",
			manufacturer: "VeryLooooooooooooooooooooooooongNewManufacturer",
			productName:  "NewProductName",
			version:      "NewVersion",
			serialNumber: "NewSerialNumber",
			want: joinBytesT(t,
				mockEntry64Raw(t, 202, 24),
				validType0BIOSInfoRaw(t),
				validType1SystemInfoData,
				"VeryLooooooooooooooooooooooooongNewManufacturer", 0,
				"NewProductName", 0,
				"NewVersion", 0,
				"NewSerialNumber", 0,
				"MockSKUNumber", 0,
				"MockFamily", 0,
				0, // Table terminator
				validEndOfTableRaw(t),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mod := newMock64Modifier(t)
			defer mod.CloseMemFile()

			opt := ReplaceSystemInfo(&tt.manufacturer, &tt.productName, &tt.version, &tt.serialNumber, nil, nil, nil, nil)

			if err := mod.Modify(opt); err != nil {
				t.Fatalf("ModifySystemInfo should pass but returned error: %v", err)
			}

			got, err := readFileFromStart(mod.memFile, int64(len(tt.want)))
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Memory content is wrong: (-want +got)\n%v", diff)
			}
		})
	}
}

func TestGetEntries(t *testing.T) {
	tests := []struct {
		name        string
		memFileData []byte
		smbiosBase  func() (int64, int64, error)
		wantE32     *Entry32
		wantE64     *Entry64
	}{
		{
			name:        "Entry32Pass",
			memFileData: mockEntry32Raw(t, 100, 200, 0, 31),
			smbiosBase:  func() (int64, int64, error) { return 0, 31, nil },
			wantE32:     mockEntry32(t, 100, 200, 0, 31),
			wantE64:     nil,
		},
		{
			name:        "Entry64Pass",
			memFileData: mockEntry64Raw(t, 100, 24),
			smbiosBase:  func() (int64, int64, error) { return 0, 24, nil },
			wantE32:     nil,
			wantE64:     mockEntry64(t, 100, 24),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := makeMemFile(t, tt.memFileData)
			m := &Modifier{
				memFile: f,
			}
			defer m.CloseMemFile()
			e32, e64, _, err := getEntries(tt.smbiosBase, f)
			if err != nil {
				t.Fatalf("getEntries should pass but return error: %v", err)
			}
			if !reflect.DeepEqual(tt.wantE32, e32) {
				t.Errorf("Wrong entry32 want %v, got %v", tt.wantE32, e32)
			}
			if !reflect.DeepEqual(tt.wantE64, e64) {
				t.Errorf("Wrong entry64 want %v, got %v", tt.wantE64, e64)
			}
		})
	}
}

func TestReplaceBaseboardInfoMotherboardPass(t *testing.T) {
	assetTag := "newTag"
	serialNumber := "newSerialNumber"
	opt := ReplaceBaseboardInfoMotherboard(nil, nil, nil, &serialNumber, &assetTag, nil, nil, nil, nil, nil)

	oldTable := []*Table{
		{
			Header: Header{
				Type:   TableTypeBaseboardInfo,
				Length: 17,
				Handle: 0,
			},
			data: []byte{
				2, 17, 0, 0, 1, 2, 3, 4, 5, 0, 6, 0, 0,
				byte(BoardTypeMotherboardIncludesProcessorMemoryAndIO),
				1, 10, 0,
			},
			strings: []string{"-", "-", "-", "-", "-", "-"},
		},
		{
			Header: Header{
				Type:   TableTypeBaseboardInfo,
				Length: 17,
				Handle: 0,
			},
			data: []byte{
				2, 17, 0, 0, 1, 2, 3, 4, 5, 0, 6, 0, 0,
				1, // BoardTypeUnknown
				1, 10, 0,
			},
			strings: []string{"-", "-", "-", "-", "-", "-"},
		},
	}
	wantTable := []*Table{
		{
			Header: Header{
				Type:   TableTypeBaseboardInfo,
				Length: 17,
				Handle: 0,
			},
			data: []byte{
				2, 17, 0, 0, 1, 2, 3, 4, 5, 0, 6, 0, 0,
				byte(BoardTypeMotherboardIncludesProcessorMemoryAndIO),
				1, 10, 0,
			},
			strings: []string{"-", "-", "-", "newSerialNumber", "newTag", "-"},
		},
		{
			Header: Header{
				Type:   TableTypeBaseboardInfo,
				Length: 17,
				Handle: 0,
			},
			data: []byte{
				2, 17, 0, 0, 1, 2, 3, 4, 5, 0, 6, 0, 0,
				1, // BoardTypeUnknown
				1, 10, 0,
			},
			strings: []string{"-", "-", "-", "-", "-", "-"},
		},
	}

	newTable, err := opt(oldTable)
	if err != nil {
		t.Fatalf("opt should pass but returned error: %v", err)
	}
	for i := 0; i < len(newTable); i++ {
		if !reflect.DeepEqual(newTable[i], wantTable[i]) {
			t.Errorf("opt return incorrect table, got %+v, want %+v", newTable[i], wantTable[i])
		}
	}
}

func TestReplaceBaseboardInfoMotherboardFail(t *testing.T) {
	malformedTable := []*Table{
		{
			Header: Header{
				Type:   TableTypeBaseboardInfo,
				Length: 17,
				Handle: 0,
			},
			data: []byte{
				2, 17, 0, 0, 1, 2, 3, 4, 5, 0, 6, 0, 0, 1,
				255, // wrong NumberOfContainedObjectHandles
				10, 0,
			},
			strings: []string{"-", "-", "-", "-", "-", "-"},
		},
	}

	opt := ReplaceBaseboardInfoMotherboard(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	if _, err := opt(malformedTable); err == nil {
		t.Fatalf("opt should fail but returned nil error")
	}
}

func TestRemoveBaseboardInfo(t *testing.T) {
	tests := []struct {
		name       string
		oldTables  []*Table
		wantTables []*Table
	}{
		{
			name: "RemoveBaseboardInfoSimple",
			oldTables: []*Table{
				{
					Header: Header{
						Type:   TableTypeBaseboardInfo,
						Length: 17,
						Handle: 0,
					},
					data: []byte{
						2, 17, 0, 0, 1, 2, 3, 4, 5, 0, 6, 0, 0,
						byte(BoardTypeSystemManagementModule),
						1, 10, 0,
					},
					strings: []string{"-", "-", "-", "-", "-", "-"},
				},
			},
			wantTables: []*Table{},
		},
		{
			name: "RemoveMultipleBaseboardInfo",
			oldTables: []*Table{
				{
					Header: Header{
						Type:   TableTypeBaseboardInfo,
						Length: 17,
						Handle: 0,
					},
					data: []byte{
						2, 17, 0, 0, 1, 2, 3, 4, 5, 0, 6, 0, 0,
						byte(BoardTypeSystemManagementModule),
						1, 10, 0,
					},
					strings: []string{"-", "-", "-", "-", "-", "-"},
				},
				{
					Header: Header{
						Type:   TableTypeBaseboardInfo,
						Length: 17,
						Handle: 1,
					},
					data: []byte{
						2, 17, 1, 0, 1, 2, 3, 4, 5, 0, 6, 0, 0,
						byte(BoardTypeSystemManagementModule),
						1, 10, 0,
					},
					strings: []string{"-", "-", "-", "-", "-", "-"},
				},
			},
			wantTables: []*Table{},
		},
		{
			name: "RemoveBaseboardInfoWithGroupAssociation",
			oldTables: []*Table{
				{
					Header: Header{
						Type:   TableTypeBaseboardInfo,
						Length: 17,
						Handle: 0,
					},
					data: []byte{
						2, 17, 0, 0, 1, 2, 3, 4, 5, 0, 6, 0, 0,
						byte(BoardTypeSystemManagementModule),
						1, 10, 0,
					},
					strings: []string{"-", "-", "-", "-", "-", "-"},
				},
				{
					Header: Header{
						Type:   TableTypeBaseboardInfo,
						Length: 17,
						Handle: 1,
					},
					data: []byte{
						2, 17, 1, 0, 1, 2, 3, 4, 5, 0, 6, 0, 0,
						byte(BoardTypeSystemManagementModule),
						1, 10, 0,
					},
					strings: []string{"-", "-", "-", "-", "-", "-"},
				},
				{
					Header: Header{
						Type:   TableTypeGroupAssociation,
						Length: 11,
						Handle: 2,
					},
					data: []byte{
						14, 11, 2, 0, // Header
						1,       // string number
						2, 0, 0, // ItemType, ItemHandle
						2, 1, 0,
					},
					strings: []string{"Group"},
				},
			},
			wantTables: []*Table{{
				Header: Header{
					Type:   TableTypeGroupAssociation,
					Length: 5,
					Handle: 2, // Note: The handle is re-indexed to 0
				},
				data: []byte{
					14, 5, 2, 0, // Header
					1, // string number
				},
				strings: []string{"Group"},
			}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			opt := RemoveBaseboardInfo(BoardTypeSystemManagementModule)
			newTables, err := opt(tc.oldTables)

			if err != nil {
				t.Fatalf("opt should pass but returned an error: %v", err)
			}

			for i := 0; i < len(newTables); i++ {
				// Using reflect.DeepEqual for a comprehensive comparison of the entire slice
				if !reflect.DeepEqual(newTables[i], tc.wantTables[i]) {
					t.Errorf("opt returned incorrect table, got %+v, want %+v", newTables[i], tc.wantTables[i])
				}
			}
		})
	}
}

func TestRemoveBaseboardInfoFail(t *testing.T) {
	opt := RemoveBaseboardInfo(BoardTypeSystemManagementModule)

	malformedTable := []*Table{
		{
			Header: Header{
				Type:   TableTypeBaseboardInfo,
				Length: 17,
				Handle: 0,
			},
			data: []byte{
				2, 17, 0, 0, 1, 2, 3, 4, 5, 0, 6, 0, 0,
				byte(BoardTypeSystemManagementModule),
				255, // wrong NumberOfContainedObjectHandles
				10, 0,
			},
			strings: []string{"-", "-", "-", "-", "-", "-"},
		},
	}

	if _, err := opt(malformedTable); err == nil {
		t.Fatalf("opt should fail but returned nil error")
	}
}
