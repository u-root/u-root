// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package universalpayload

import (
	"bytes"
	"debug/pe"
	"encoding/binary"
	"errors"
	"io"
	"reflect"
	"testing"
	"unsafe"

	"github.com/u-root/u-root/pkg/dt"
)

func fdtReader(t *testing.T, fdt *dt.FDT) io.ReaderAt {
	t.Helper()
	var b bytes.Buffer
	fdt.Header.Magic = dt.Magic
	fdt.Header.Version = 17
	if _, err := fdt.Write(&b); err != nil {
		t.Fatal(err)
	}
	return bytes.NewReader(b.Bytes())
}

func TestGetFdtInfo(t *testing.T) {
	for _, tt := range []struct {
		// Inputs
		name string
		fdt  io.ReaderAt

		// Results
		fdtLoad *FdtLoad
		err     error
	}{
		// CASE 1: normal case, a FdtLoad object is returned with expected values
		{
			name: "testdata/upl.dtb",
			fdt: fdtReader(t, &dt.FDT{
				RootNode: dt.NewNode("/", dt.WithChildren(
					dt.NewNode("images", dt.WithChildren(
						dt.NewNode("tianocore", dt.WithProperty(
							dt.PropertyString("arch", "x86_64"),
							dt.PropertyU64("entry-start", 0x00805ac3),
							dt.PropertyU64("load", 0x00800000),
							dt.PropertyU32("data-offset", 0x00001000),
							dt.PropertyU32("data-size", 0x0000c000),
						)),
					)),
				)),
			}),
			fdtLoad: &FdtLoad{
				Load:       uint64(0x0000800000),
				EntryStart: uint64(0x0000805ac3),
			},
			err: nil,
		},
		// CASE 2: dtb file not found
		{
			name:    "testdata/not_exist_file.dtb",
			fdt:     nil,
			fdtLoad: nil,
			err:     ErrFailToReadFdtFile,
		},
		// CASE 3: missing first level node: /images
		{
			name: "testdata/missing_first_node_images.dtb",
			fdt: fdtReader(t, &dt.FDT{
				RootNode: dt.NewNode("/", dt.WithChildren(
					dt.NewNode("description", dt.WithProperty(
						dt.PropertyString("arch", "x86_64"),
					)),
				)),
			}),
			fdtLoad: nil,
			err:     ErrNodeImagesNotFound,
		},
		// CASE 4: missing second level node: /images/tianocore
		{
			name: "testdata/missing_second_node_images.dtb",
			fdt: fdtReader(t, &dt.FDT{
				RootNode: dt.NewNode("/", dt.WithChildren(
					dt.NewNode("images", dt.WithProperty(
						dt.PropertyString("arch", "x86_64"),
					)),
				)),
			}),
			fdtLoad: nil,
			err:     ErrNodeTianocoreNotFound,
		},
		// CASE 5: failed to get /images/tianocore/load property
		{
			name: "testdata/missing_property_load.dtb",
			fdt: fdtReader(t, &dt.FDT{
				RootNode: dt.NewNode("/", dt.WithChildren(
					dt.NewNode("images", dt.WithChildren(
						dt.NewNode("tianocore", dt.WithProperty(
							dt.PropertyString("arch", "x86_64"),
							dt.PropertyU64("entry-start", 0x00805ac3),
						)),
					)),
				)),
			}),
			fdtLoad: nil,
			err:     ErrNodeLoadNotFound,
		},
		// CASE 6: failed to convert /images/tianocore/load property (type error)
		{
			name: "testdata/missing_property_load.dtb",
			fdt: fdtReader(t, &dt.FDT{
				RootNode: dt.NewNode("/", dt.WithChildren(
					dt.NewNode("images", dt.WithChildren(
						dt.NewNode("tianocore", dt.WithProperty(
							dt.PropertyString("arch", "x86_64"),
							dt.PropertyString("load", "0x00800000"),
						)),
					)),
				)),
			}),
			fdtLoad: nil,
			err:     ErrFailToConvertLoad,
		},
		// CASE 7: failed to get /images/tianocore/entry-start property
		{
			name: "testdata/missing_property_entry_start.dtb",
			fdt: fdtReader(t, &dt.FDT{
				RootNode: dt.NewNode("/", dt.WithChildren(
					dt.NewNode("images", dt.WithChildren(
						dt.NewNode("tianocore", dt.WithProperty(
							dt.PropertyString("arch", "x86_64"),
							dt.PropertyU64("load", 0x00800000),
						)),
					)),
				)),
			}),
			fdtLoad: nil,
			err:     ErrNodeEntryStartNotFound,
		},
		// CASE 8: failed to convert /images/tianocore/entry-start property (type error)
		{
			name: "testdata/fail_convert_property_entry_start.dtb",
			fdt: fdtReader(t, &dt.FDT{
				RootNode: dt.NewNode("/", dt.WithChildren(
					dt.NewNode("images", dt.WithChildren(
						dt.NewNode("tianocore", dt.WithProperty(
							dt.PropertyString("arch", "x86_64"),
							dt.PropertyU64("load", 0x00800000),
							dt.PropertyString("entry-start", "0x00800000"),
						)),
					)),
				)),
			}),
			fdtLoad: nil,
			err:     ErrFailToConvertEntryStart,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getFdtInfo(tt.name, tt.fdt)
			if tt.err != nil {
				if err == nil {
					t.Fatalf("Expected error %q, got nil", tt.err)
				}
				if !errors.Is(err, tt.err) {
					t.Errorf("Unxpected error %q, want = %q", err.Error(), tt.err)
				}
			} else if err != nil {
				t.Fatal(err)
			}

			if tt.fdtLoad != nil && got == nil {
				t.Fatalf("getFdtInfo fdtLoad = nil, want = %v", tt.fdtLoad)
			}

			if tt.fdtLoad != nil {
				if tt.fdtLoad.Load != got.Load {
					t.Fatalf("getFdtInfo fdtLoad.Load = %d, want = %v", got.Load, tt.fdtLoad.Load)
				}
				if tt.fdtLoad.EntryStart != got.EntryStart {
					t.Fatalf("getFdtInfo fdtLoad.EntryStart = %v, want = %v", got.EntryStart, tt.fdtLoad.EntryStart)
				}
			}
		})
	}
}

func TestAlignHOBLength(t *testing.T) {
	tests := []struct {
		name        string
		expectLen   uint64
		bufLen      int
		expectedErr error
	}{
		{
			name:      "Exact Length",
			expectLen: 5,
			bufLen:    5,
		},
		{
			name:      "Padding Required",
			expectLen: 10,
			bufLen:    5,
		},
		{
			name:        "Negative Padding",
			expectLen:   3,
			bufLen:      5,
			expectedErr: ErrAlignPadRange,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBufferString("12345")
			err := alignHOBLength(tt.expectLen, tt.bufLen, buf)

			if tt.expectedErr == nil {
				// success validation
				if err != nil {
					t.Fatalf("Unexpected error: %+v", err)
				}
				if uint64(buf.Len()) != tt.expectLen {
					t.Fatalf("alignHOBLength() got = %d, want = %d", buf.Len(), tt.expectLen)
				}
			} else {
				// fault validation
				if err == nil {
					t.Fatalf("Expected error %q, got nil", tt.expectedErr)
				}
				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("Unxpected error %+v, want = %q", err, tt.expectedErr)
				}
			}
		})
	}
}

type MockSection struct {
	name string
	data []byte
}

const (
	peFileHeaderSize    = uint32(unsafe.Sizeof(pe.FileHeader{}))
	peOptHeaderSize     = uint32(unsafe.Sizeof(pe.OptionalHeader64{}))
	peSectionHeaderSize = uint32(unsafe.Sizeof(pe.SectionHeader32{}))
	peSectionSize       = uint32(10) // for simplicity, use constant section size
)

func mockWritePeFileBinary(offset int, totalSize int, imageBase uint64, sections []*MockSection) []byte {
	// Serialize the PE file data to bytes
	var buf bytes.Buffer

	// write offset pads
	_ = binary.Write(&buf, binary.LittleEndian, make([]byte, offset))

	// write FileHeader
	fileHeader := pe.FileHeader{
		Machine:              pe.IMAGE_FILE_MACHINE_AMD64,
		SizeOfOptionalHeader: uint16(peOptHeaderSize),
		NumberOfSections:     uint16(len(sections)),
	}
	_ = binary.Write(&buf, binary.LittleEndian, fileHeader)

	// write OptionalHeader64
	optionalHeader := pe.OptionalHeader64{
		Magic:               0x20b,
		ImageBase:           imageBase,
		NumberOfRvaAndSizes: 16,
	}
	_ = binary.Write(&buf, binary.LittleEndian, optionalHeader)

	// write Sections
	for i, section := range sections {
		var nameArr [8]uint8
		copy(nameArr[:], section.name)
		sectionHeader := pe.SectionHeader32{
			Name:          nameArr,
			SizeOfRawData: uint32(len(section.data)),
			PointerToRawData: peFileHeaderSize + peOptHeaderSize +
				peSectionHeaderSize*uint32(i+1) + peSectionSize*uint32(i)}
		_ = binary.Write(&buf, binary.LittleEndian, sectionHeader)
		_ = binary.Write(&buf, binary.LittleEndian, section.data)
	}

	// write extra pads ( for relocation)
	if totalSize-buf.Len() > 0 {
		_ = binary.Write(&buf, binary.LittleEndian, make([]byte, totalSize-buf.Len()))
	}

	return buf.Bytes()
}

func TestRelocateFdtData(t *testing.T) {
	tests := []struct {
		name        string
		dst         uint64
		fdtLoad     *FdtLoad
		data        []byte
		wantErr     error
		wantFdtLoad *FdtLoad
	}{
		{
			name:    "Valid PE relocation using mock struct",
			dst:     0x2000,
			fdtLoad: &FdtLoad{DataOffset: 0x100, DataSize: 0x5000, EntryStart: 0x1000, Load: 0x1800},
			data: mockWritePeFileBinary(0x100, 0x5100, 0x1000, []*MockSection{
				{".reloc", mockRelocData(0x1000, IMAGE_REL_BASED_DIR64, 0x200)},
			}),
			wantErr:     nil,
			wantFdtLoad: &FdtLoad{DataOffset: 0x100, DataSize: 0x5000, EntryStart: 0x2000 + (0x1000 - 0x1800), Load: 0x2000},
		},
		{
			name:    "Relocation address out of bounds",
			dst:     0x6000,
			fdtLoad: &FdtLoad{DataOffset: 0x100, DataSize: 0x5000, EntryStart: 0x1000, Load: 0x1800},
			data: mockWritePeFileBinary(0x100, 0x5100, 0x1000, []*MockSection{
				{".reloc", mockRelocData(0x6000, IMAGE_REL_BASED_DIR64, 0x200)},
			}),
			wantErr:     ErrPeRelocOutOfBound,
			wantFdtLoad: nil,
		},
		{
			name:        "No .reloc section found in PE file",
			dst:         0x2000,
			fdtLoad:     &FdtLoad{DataOffset: 0x100, DataSize: 0x5000, EntryStart: 0x1000, Load: 0x1800},
			data:        mockWritePeFileBinary(0x100, 0x5100, 0x1000, []*MockSection{}),
			wantErr:     nil,
			wantFdtLoad: &FdtLoad{DataOffset: 0x100, DataSize: 0x5000, EntryStart: 0x2000 + (0x1000 - 0x1800), Load: 0x2000},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := relocateFdtData(tt.dst, tt.fdtLoad, tt.data)
			expectErr(t, err, tt.wantErr)
			if tt.wantFdtLoad != nil && !reflect.DeepEqual(tt.fdtLoad, tt.wantFdtLoad) {
				t.Fatalf("Unexpected relocated FdtLoad: %v, want: %v", *tt.fdtLoad, *tt.wantFdtLoad)
			}
		})
	}

}

// Helper function to mock relocData for test cases
func mockRelocData(pageRVA uint32, entryType, entryOffset uint16) []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, pageRVA)
	_ = binary.Write(buf, binary.LittleEndian, uint32(10)) // block size including header

	// Type is in the high 4 bits, offset is in the low 12 bits
	entry := (entryType << 12) | entryOffset
	_ = binary.Write(buf, binary.LittleEndian, entry)
	return buf.Bytes()
}

// Helper function to mock data with specific content
func mockData(size int, offset int, value uint64) []byte {
	data := make([]byte, size)
	binary.LittleEndian.PutUint64(data[offset:], value)
	return data
}

// Helper function to mock expected data after relocation
func mockExpectedData(size int, offset int, relocatedValue uint64) []byte {
	data := make([]byte, size)
	binary.LittleEndian.PutUint64(data[offset:], relocatedValue)
	return data
}

func TestRelocatePE(t *testing.T) {
	tests := []struct {
		name      string
		relocData []byte
		delta     uint64
		data      []byte
		expected  []byte
		wantErr   error
	}{
		{
			name:      "Valid relocation",
			relocData: mockRelocData(0x1000, IMAGE_REL_BASED_DIR64, 0x200),
			delta:     0x1000,
			data:      mockData(0x10000, 0x1000+0x200, 0x4000),
			expected:  mockExpectedData(0x10000, 0x1000+0x200, 0x4000+0x1000),
			wantErr:   nil,
		},
		{
			name:      "Out of bounds relocation",
			relocData: mockRelocData(0x1000, IMAGE_REL_BASED_DIR64, uint16(0x800)), // relocation out of bounds
			delta:     0x1000,
			data:      mockData(0x1500, 0x1000+0x200, 0x4000),
			expected:  nil,
			wantErr:   ErrPeRelocOutOfBound,
		},
		{
			name:      "Fail to get block size",
			relocData: []byte{0x01, 0x02, 0x03, 0x04}, // insufficient data
			delta:     0x1000,
			data:      mockData(0x1000, 0x200, 0x4000),
			expected:  nil,
			wantErr:   ErrPeFailToGetBlockSize,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataCopy := make([]byte, len(tt.data))
			copy(dataCopy, tt.data)

			err := relocatePE(tt.relocData, tt.delta, dataCopy)

			expectErr(t, err, tt.wantErr)

			if tt.wantErr == nil && !bytes.Equal(dataCopy, tt.expected) {
				t.Errorf("expected data %v, got %v", tt.expected, dataCopy)
			}
		})
	}
}
