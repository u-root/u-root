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
	"os"
	"path/filepath"
	"reflect"
	"strings"
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
	binary.Write(&buf, binary.LittleEndian, make([]byte, offset))

	// write FileHeader
	fileHeader := pe.FileHeader{
		Machine:              pe.IMAGE_FILE_MACHINE_AMD64,
		SizeOfOptionalHeader: uint16(peOptHeaderSize),
		NumberOfSections:     uint16(len(sections)),
	}
	binary.Write(&buf, binary.LittleEndian, fileHeader)

	// write OptionalHeader64
	optionalHeader := pe.OptionalHeader64{
		Magic:               0x20b,
		ImageBase:           imageBase,
		NumberOfRvaAndSizes: 16,
	}
	binary.Write(&buf, binary.LittleEndian, optionalHeader)

	// write Sections
	for i, section := range sections {
		var nameArr [8]uint8
		copy(nameArr[:], section.name)
		sectionHeader := pe.SectionHeader32{
			Name:          nameArr,
			SizeOfRawData: uint32(len(section.data)),
			PointerToRawData: peFileHeaderSize + peOptHeaderSize +
				peSectionHeaderSize*uint32(i+1) + peSectionSize*uint32(i)}
		binary.Write(&buf, binary.LittleEndian, sectionHeader)
		binary.Write(&buf, binary.LittleEndian, section.data)
	}

	// write extra pads ( for relocation)
	if totalSize-buf.Len() > 0 {
		binary.Write(&buf, binary.LittleEndian, make([]byte, totalSize-buf.Len()))
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
	binary.Write(buf, binary.LittleEndian, pageRVA)
	binary.Write(buf, binary.LittleEndian, uint32(10)) // block size including header

	// Type is in the high 4 bits, offset is in the low 12 bits
	entry := (entryType << 12) | entryOffset
	binary.Write(buf, binary.LittleEndian, entry)
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

func mockGetSMBIOS3Base() (int64, int64, error) {
	return 0x100, 0x200, nil
}

func TestConstructSMBIOS3Node(t *testing.T) {
	// mock data
	defer func(old func() (int64, int64, error)) { getSMBIOSBase = old }(getSMBIOSBase)
	getSMBIOSBase = mockGetSMBIOS3Base

	tests := []struct {
		name     string
		wantNode *dt.Node
		wantErr  error
	}{
		{
			name:     "Invalid SMBIOS3 Header size",
			wantNode: nil,
			wantErr:  ErrSMBIOS3NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			smbiosNode, err := constructSMBIOS3Node()
			expectErr(t, err, tt.wantErr)
			if smbiosNode != tt.wantNode {
				t.Fatalf("Unexpected smbios Node: actual(%v) vs. expected(nil)", smbiosNode)
			}
		})
	}
}

func TestFetchACPIMCFGDataNegative(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr error
	}{
		{
			name: "MCFG Data length short",
			data: []byte{
				'M', 'C', 'F', 'G', // Signature
				0x3c, 0x00, 0x00, 0x00, // Length
				0x01,                         // Revision
				0x00,                         // Checksum (ignored)
				'u', '-', 'r', 'o', 'o', 't', // OemId
				'F', 'A', 'K', 'E', 'T', 'a', 'b', 'l', // OemTableId
				0x01, 0x00, 0x00, 0x00, // OemRevision
				'U', 'R', 'O', 'T', // CreatorId
				0x01, 0x00, 0x00, 0x00, // CreatorRevision
				0x00, 0x00, 0x00, 0x00, // Reserved
			},
			wantErr: ErrMcfgDataLenthTooShort,
		},
		{
			name: "MCFG Data magic mismatch",
			data: []byte{
				'G', 'F', 'C', 'M', // Signature (Invalid Magic)
				0x3c, 0x00, 0x00, 0x00, // Length
				0x01,                         // Revision
				0x00,                         // Checksum (ignored)
				'u', '-', 'r', 'o', 'o', 't', // OemId
				'F', 'A', 'K', 'E', 'T', 'a', 'b', 'l', // OemTableId
				0x01, 0x00, 0x00, 0x00, // OemRevision
				'U', 'R', 'O', 'T', // CreatorId
				0x01, 0x00, 0x00, 0x00, // CreatorRevision
				0x00, 0x00, 0x00, 0x00, // Reserved
				0x00, 0x00, 0x00, 0xE0, // Base Address low parts
				0x00, 0x00, 0x00, 0x00, // Base Address high parts
				0x00, 0x00, // Pci Segment Group Number
				0x00,                   // Start Bus Number
				0xFF,                   // End Bus Number
				0x00, 0x00, 0x00, 0x00, // Reserved
				0x00, 0x00, 0x00, 0x00, // Reserved
			},
			wantErr: ErrMcfgSignatureMismatch,
		},
		{
			name: "MCFG Data Base Address Allocation corrupt",
			data: []byte{
				'M', 'C', 'F', 'G', // Signature
				0x3c, 0x00, 0x00, 0x00, // Length
				0x01,                         // Revision
				0x00,                         // Checksum (ignored)
				'u', '-', 'r', 'o', 'o', 't', // OemId
				'F', 'A', 'K', 'E', 'T', 'a', 'b', 'l', // OemTableId
				0x01, 0x00, 0x00, 0x00, // OemRevision
				'U', 'R', 'O', 'T', // CreatorId
				0x01, 0x00, 0x00, 0x00, // CreatorRevision
				0x00, 0x00, 0x00, 0x00, // Reserved
				0x00, 0x00, 0x00, 0x00, // Reserved
				0x00, 0x00, 0x00, 0xE0, // Base Address low parts
				0x00, 0x00, 0x00, 0x00, // Base Address high parts
				0x00, 0x00, // Pci Segment Group Number
				0x00,                   // Start Bus Number
				0xFF,                   // End Bus Number
				0x00, 0x00, 0x00, 0x00, // Reserved
				0x00, 0x00, 0x08, 0xE0, // Base Address low parts -- Invalid
			},
			wantErr: ErrMcfgBaseAddrAllocCorrupt,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mcfgData, err := fetchACPIMCFGData(tt.data)
			expectErr(t, err, tt.wantErr)
			if mcfgData != nil {
				t.Fatalf("Invalid return value with %s, expected(nil)", tt.name)
			}
		})
	}
}

func TestFetchACPIMCFGData(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		mcfgData []MCFGBaseAddressAllocation
	}{
		{
			name: "MCFG Data single segment",
			data: []byte{
				'M', 'C', 'F', 'G', // Signature
				0x3c, 0x00, 0x00, 0x00, // Length
				0x01,                         // Revision
				0x00,                         // Checksum (ignored)
				'u', '-', 'r', 'o', 'o', 't', // OemId
				'F', 'A', 'K', 'E', 'T', 'a', 'b', 'l', // OemTableId
				0x01, 0x00, 0x00, 0x00, // OemRevision
				'U', 'R', 'O', 'T', // CreatorId
				0x01, 0x00, 0x00, 0x00, // CreatorRevision
				0x00, 0x00, 0x00, 0x00, // Reserved
				0x00, 0x00, 0x00, 0x00, // Reserved
				0x00, 0x00, 0x00, 0xE0, // Base Address low parts
				0x00, 0x00, 0x00, 0x00, // Base Address high parts
				0x00, 0x00, // Pci Segment Group Number
				0x00,                   // Start Bus Number
				0xFF,                   // End Bus Number
				0x00, 0x00, 0x00, 0x00, // Reserved
			},
			mcfgData: []MCFGBaseAddressAllocation{
				{
					BaseAddr:  0xE000_0000,
					PCISegGrp: 0x00,
					StartBus:  0x00,
					EndBus:    0xFF,
					Reserved:  0x00,
				},
			},
		},
		{
			name: "MCFG Data multiple segments",
			data: []byte{
				'M', 'C', 'F', 'G', // Signature
				0x3c, 0x00, 0x00, 0x00, // Length
				0x01,                         // Revision
				0x00,                         // Checksum (ignored)
				'u', '-', 'r', 'o', 'o', 't', // OemId
				'F', 'A', 'K', 'E', 'T', 'a', 'b', 'l', // OemTableId
				0x01, 0x00, 0x00, 0x00, // OemRevision
				'U', 'R', 'O', 'T', // CreatorId
				0x01, 0x00, 0x00, 0x00, // CreatorRevision
				0x00, 0x00, 0x00, 0x00, // Reserved
				0x00, 0x00, 0x00, 0x00, // Reserved
				0x00, 0x00, 0x00, 0xE0, // Base Address low parts
				0x00, 0x00, 0x00, 0x00, // Base Address high parts
				0x00, 0x00, // Pci Segment Group Number
				0x00,                   // Start Bus Number
				0xFF,                   // End Bus Number
				0x00, 0x00, 0x00, 0x00, // Reserved
				0x00, 0x00, 0x00, 0xE8, // Base Address low parts
				0x00, 0x00, 0x00, 0x00, // Base Address high parts
				0x01, 0x00, // Pci Segment Group Number
				0x10,                   // Start Bus Number
				0xE0,                   // End Bus Number
				0x00, 0x00, 0x00, 0x00, // Reserved
				0x00, 0x00, 0x00, 0xEC, // Base Address low parts
				0x00, 0x00, 0x00, 0x00, // Base Address high parts
				0x04, 0x00, // Pci Segment Group Number
				0x18,                   // Start Bus Number
				0xC8,                   // End Bus Number
				0x00, 0x00, 0x00, 0x00, // Reserved
			},
			mcfgData: []MCFGBaseAddressAllocation{
				{
					BaseAddr:  0xE000_0000,
					PCISegGrp: 0x00,
					StartBus:  0x00,
					EndBus:    0xFF,
					Reserved:  0x00,
				},
				{
					BaseAddr:  0xE800_0000,
					PCISegGrp: 0x01,
					StartBus:  0x10,
					EndBus:    0xE0,
					Reserved:  0x00,
				},
				{
					BaseAddr:  0xEC00_0000,
					PCISegGrp: 0x04,
					StartBus:  0x18,
					EndBus:    0xC8,
					Reserved:  0x00,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if mcfgData, err := fetchACPIMCFGData(tt.data); err != nil {
				t.Fatalf("Unexpected error:%v", err)
			} else {
				if !reflect.DeepEqual(mcfgData, tt.mcfgData) {
					t.Errorf("got %+v, want %+v", mcfgData, tt.mcfgData)
				}

			}
		})
	}
}

func TestRetrieveRootBridgeResources(t *testing.T) {
	mcfgData := []MCFGBaseAddressAllocation{
		{
			BaseAddr:  0xB000_0000,
			PCISegGrp: 0x00,
			StartBus:  0x10,
			EndBus:    0xE0,
			Reserved:  0x00,
		},
		{
			BaseAddr:  0xA000_0000,
			PCISegGrp: 0x02,
			StartBus:  0x00,
			EndBus:    0xFF,
			Reserved:  0x00,
		},
	}

	subFolder := []string{
		"0000:00:00.0", // Out of Bus range
		"0000:01:02.0", // Out of Bus range
		"0000:15:00.0", // Valid Bus
		"0000:15:02.0", // Valid Bus
		"0000:c8:00.0", // Valid Bus
		"0000:e1:00.0", // Out of Bus range
		"0000:ff:00.0", // Out of Bus range
		"0002:00:00.0", // Valid Bus
		"0002:00:02.0", // Valid Bus
		"0002:1f:00.0", // Valid Bus
		"0002:af:00.0", // Valid Bus
		"0002:ff:00.0", // Valid Bus
	}

	resourceContent := [][]string{
		{
			// Content for "0000:00:00.0"
			"0x00000000DF000000 0x00000000DFFFFFFF 0x0000000000040200\n",
		},
		{
			// Content for "0000:01:02.0"
			"0x00000000E0000000 0x00000000E0FFFFFF 0x0000000000040200\n",
		},
		{
			// Content for "0000:15:00.0"
			"0x00000000DE000000 0x00000000DEFFFFFF 0x0000000000040200\n",
			"0x00000000C0000000 0x00000000CFFFFFFF 0x0000000000040200\n",
			"0x000000000000F000 0x000000000000FFFF 0x0000000000040101\n",
			"0x00000000000C0000 0x00000000000DFFFF 0x0000000000000212\n",
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
			"0x0000000800000000 0x0000000800EFFFFF 0x0000000000140204\n",
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
		},
		{
			// Content for "0000:15:02.0"
			"0x00000000B8000000 0x00000000B87FFFFF 0x0000000000040200\n",
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
		},
		{
			// Content for "0000:C8:00.0"
			"0x00000000B8800000 0x00000000B8FFFFFF 0x0000000000040200\n",
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
			"0x000000000000E000 0x000000000000E03F 0x0000000000040101\n",
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
			"0x0000000800F00000 0x0000000800FFFFFF 0x0000000000140204\n",
		},
		{
			// Content for "0000:E1:00.0"
			"0x00000000B0000000 0x00000000B0FFFFFF 0x0000000000040200\n",
			"0x0000000801000000 0x0000000801FFFFFF 0x0000000000140204\n",
		},
		{
			// Content for "0000:FF:00.0"
			"0x00000000A8000000 0x00000000AFFFFFFF 0x0000000000040200\n",
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
		},
		{
			// Content for "0002:00:00.0"
			"0x0000000810000000 0x000000081000FFFF 0x0000000000140204\n",
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
		},
		{
			// Content for "0002:00:02.0"
			"0x00000000A0070000 0x00000000A07FFFFF 0x0000000000040200\n",
			"0x000000000000B000 0x000000000000BFFF 0x0000000000040101\n",
			"0x00000000000C0000 0x00000000000DFFFF 0x0000000000000212\n",
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
			"0x0000000810010000 0x000000081001FFFF 0x0000000000140204\n",
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
		},
		{
			// Content for "0002:1f:00.0"
			"0x00000000A0080000 0x00000000A09FFFFF 0x0000000000046200\n",
			"0x00000000A0040000 0x00000000A004FFFF 0x0000000000040200\n",
			"0x00000000A0060000 0x00000000A006FFFF 0x0000000000040200\n",
			"0x0000000070000000 0x00000008101FFFFF 0x0000000000140204\n",
		},
		{
			// Content for "0002:af:00.0"
			"0x00000000A0030000 0x00000000A003FFFF 0x0000000000040200\n",
			"0x0000000810020000 0x00000008103FFFFF 0x0000000000140204\n",
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
		},
		{
			// Content for "0002:ff:00.0"
			"0x00000000A0000000 0x00000000A000FFFF 0x0000000000040200\n",
			"0x000000000000A000 0x000000000000A03F 0x0000000000040101\n",
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
		},
	}

	expectedResourceRegion := []ResourceRegions{
		{
			MMIO64Base:  0x0000_0008_0000_0000,
			MMIO64Limit: 0x0000_0000_0100_0000,
			MMIO32Base:  0x0000_0000_B800_0000,
			MMIO32Limit: 0x0000_0000_2700_0000,
			IOPortBase:  0x0000_0000_0000_E000,
			IOPortLimit: 0x0000_0000_0000_2000,
		},
		{
			MMIO64Base:  0x0000_0008_1000_0000,
			MMIO64Limit: 0x0000_0000_0040_0000,
			MMIO32Base:  0x0000_0000_A000_0000,
			MMIO32Limit: 0x0000_0000_0080_0000,
			IOPortBase:  0x0000_0000_0000_A000,
			IOPortLimit: 0x0000_0000_0000_2000,
		},
	}

	expectedRbNodes := []*dt.Node{
		dt.NewNode("pci-rb", dt.WithProperty(
			dt.PropertyString("compatible", "pci-rb"),
			dt.PropertyU64("reg", 0xB000_0000),
			dt.PropertyU32Array("bus-range", []uint32{0x10, 0xE0}),
			dt.PropertyU32Array("ranges", []uint32{
				0x300_0000,               // 64BITS
				0x0000_0008, 0x0000_0000, // MMIO64 Base high and low
				0x0, 0x0,
				0x0000_0000, 0x0100_0000, // MMIO64 Limit high and low
				0x200_0000,               // 32BITS
				0x0000_0000, 0xB800_0000, // MMIO32 Base high and low
				0x0, 0x0,
				0x0000_0000, 0x2700_0000, // MMIO32 Limit high and low
				0x100_0000,               // IOPort
				0x0000_0000, 0x0000_E000, // IOPort Base high and low
				0x0, 0x0,
				0x0000_0000, 0x0000_2000, // IOPort Limit high and low
			}),
		)),
		dt.NewNode("pci-rb", dt.WithProperty(
			dt.PropertyString("compatible", "pci-rb"),
			dt.PropertyU64("reg", 0xA000_0000),
			dt.PropertyU32Array("bus-range", []uint32{0x00, 0xFF}),
			dt.PropertyU32Array("ranges", []uint32{
				0x300_0000,               // 64BITS
				0x0000_0008, 0x1000_0000, // MMIO64 Base high and low
				0x0, 0x0,
				0x0000_0000, 0x0040_0000, // MMIO64 Limit high and low
				0x200_0000,               // 32BITS
				0x0000_0000, 0xA000_0000, // MMIO32 Base high and low
				0x0, 0x0,
				0x0000_0000, 0x0080_0000, // MMIO32 Limit high and low
				0x100_0000,               // IOPort
				0x0000_0000, 0x0000_A000, // IOPort Base high and low
				0x0, 0x0,
				0x0000_0000, 0x0000_2000, // IOPort Limit high and low
			}),
		)),
	}

	// Create temp folder for testing
	tmpDir := t.TempDir()

	// Write data to files:
	// /$TMPDIR/$DOMAIN_ID:$BUS_ID:$DEVICE_ID.$FUNCTION_ID/resource
	for i, folderName := range subFolder {
		subFolderPath := filepath.Join(tmpDir, folderName)
		if err := os.MkdirAll(subFolderPath, 0755); err != nil {
			t.Fatalf("Error creating subfolder %s: %v\n", subFolderPath, err)
		}

		filePath := filepath.Join(subFolderPath, "resource")
		data := strings.Join(resourceContent[i], "")

		if err := os.WriteFile(filePath, []byte(data), 0644); err != nil {
			t.Fatalf("Error writing to file %s: %v\n", filePath, err)
		}
	}

	for idx, item := range mcfgData {
		resource, err := retrieveRootBridgeResources(tmpDir, item)
		if err != nil {
			t.Fatalf("Failed to retrieve RB resource %v\n", err)
		}

		if !reflect.DeepEqual(*resource, expectedResourceRegion[idx]) {
			t.Errorf("got %+v, want %+v", resource, expectedResourceRegion[idx])
		}

		rbNode, err := createPCIRootBridgeNode(tmpDir, item)
		if err != nil {
			t.Fatalf("Failed to create RB node %v\n", err)
		}

		if !reflect.DeepEqual(rbNode, expectedRbNodes[idx]) {
			t.Errorf("\ngot %+v, want %+v", rbNode, expectedRbNodes[idx])
		}
	}

}
