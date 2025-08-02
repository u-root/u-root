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

	"github.com/google/go-cmp/cmp"
	"github.com/u-root/u-root/pkg/boot/kexec"
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
				peSectionHeaderSize*uint32(i+1) + peSectionSize*uint32(i),
		}
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
	// Mock the kexecMemoryMapFromIOMem function to return a test memory map
	defer func(old func() (kexec.MemoryMap, error)) { kexecMemoryMapFromIOMem = old }(kexecMemoryMapFromIOMem)
	kexecMemoryMapFromIOMem = func() (kexec.MemoryMap, error) {
		return kexec.MemoryMap{
			kexec.TypedRange{Range: kexec.Range{Start: 0x1000, Size: 0x400000}, Type: kexec.RangeRAM},
			kexec.TypedRange{Range: kexec.Range{Start: 0x500000, Size: 0x100000}, Type: kexec.RangeReserved},
			kexec.TypedRange{Range: kexec.Range{Start: 0x600000, Size: 0x200000}, Type: kexec.RangeACPI},
		}, nil
	}

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
		"pci0000:00/0000:00:00.0",                           // Out of Bus range
		"pci0000:01/0000:01:02.0",                           // Out of Bus range
		"pci0000:14/0000:14:00.0",                           // Valid Bus
		"pci0000:14/0000:14:02.0",                           // Valid Bus
		"pci0000:16/0000:16:00.0",                           // Valid Bus
		"pci0000:16/0000:16:01.0",                           // Valid Bus
		"pci0000:16/0000:16:01.0/0000:17:00.0",              // Valid Bus
		"pci0000:16/0000:16:01.0/0000:17:00.0/0000:18:01.0", // Valid Bus
		"pci0000:16/0000:16:01.0/0000:17:00.0/0000:18:02.0", // Valid Bus
		"pci0000:16/0000:16:02.0/0000:19:00.1",              // Valid Bus
		"pci0000:16/0000:16:03.0",                           // Valid Bus
		"pci0000:1a/0000:1a:00.0",                           // Valid Bus with Invalid Resource
		"pci0000:1a/0000:1a:00.1",                           // Valid Bus with Invalid Resource
		"pci0000:1a/0000:1a:00.2",                           // Valid Bus with Invalid Resource
		"pci0000:1a/0000:1a:00.3",                           // Valid Bus with Invalid Resource
		"pci0000:1a/0000:1a:00.4",                           // Valid Bus with Invalid Resource
		"pci0000:1a/0000:1a:00.5",                           // Valid Bus with Invalid Resource
		"pci0000:1a/0000:1a:00.6",                           // Valid Bus with Invalid Resource
		"pci0000:1a/0000:1a:00.7",                           // Valid Bus with Invalid Resource
		"pci0000:c8/0000:c8:00.0",                           // Valid Bus
		"pci0000:e1/0000:e1:00.0",                           // Out of Bus range
		"pci0000:ff/0000:ff:00.0",                           // Out of Bus range
		"pci0002:00/0002:00:00.0",                           // Valid Bus
		"pci0002:00/0002:00:02.0",                           // Valid Bus
		"pci0002:1f/0002:1f:00.0",                           // Valid Bus
		"pci0002:af/0002:af:00.0",                           // Valid Bus
		"pci0002:ff/0002:ff:00.0",                           // Valid Bus

	}

	resourceContent := [][]string{
		{
			// Content for "pci0000:00/0000:00:00.0" // Out of Bus range
			"0x00000000DF000000 0x00000000DFFFFFFF 0x0000000000040200\n",
		},
		{
			// Content for "pci0000:01/0000:01:02.0" // Out of Bus range
			"0x00000000E0000000 0x00000000E0FFFFFF 0x0000000000040200\n",
		},
		{
			// Content for "pci0000:14/0000:14:00.0" // Valid Bus
			"0x00000000DE000000 0x00000000DEFFFFFF 0x0000000000040200\n",
			"0x00000000C0000000 0x00000000CFFFFFFF 0x0000000000040200\n",
			"0x000000000000E000 0x000000000000EFFF 0x0000000000040101\n",
			"0x00000000000C0000 0x00000000000DFFFF 0x0000000000000212\n",
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
			"0x0000000800000000 0x0000000800EFFFFF 0x0000000000140204\n",
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
		},
		{
			// Content for "pci0000:14/0000:14:02.0" // Valid Bus
			"0x00000000DC000000 0x00000000DC00FFFF 0x0000000000040200\n",
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
		},
		{
			// Content for "pci0000:16/0000:16:00.0" // Valid Bus
			"0x00000000B8800000 0x00000000B8803FFF 0x0000000000040200\n",
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
		},

		{
			// Content for "pci0000:16/0000:16:01.0" // Valid Bus
			"0x00000000B8804000 0x00000000B880FFFF 0x0000000000040200\n",
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
			"0x0000000000010000 0x0000000000011FFF 0x0000000000040101\n",
		},
		{
			// Content for "pci0000:16/0000:16:01.0/0000:17:00.0" // Valid Bus
			"0x00000000B8804000 0x00000000B880FFFF 0x0000000000040200\n",
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
			"0x0000000000010000 0x0000000000011FFF 0x0000000000040101\n",
		},
		{
			// Content for "pci0000:16/0000:16:01.0/0000:17:00.0/0000:18:01.0" // Valid Bus
			"0x00000000B8804000 0x00000000B8807FFF 0x0000000000040200\n",
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
			"0x00000000B8809000 0x00000000B880EFFF 0x0000000000040200\n",
			"0x0000000000010000 0x0000000000010FFF 0x0000000000040101\n",
		},
		{
			// Content for "pci0000:16/0000:16:01.0/0000:17:00.0/0000:18:02.0" // Valid Bus
			"0x00000000B8808000 0x00000000B8808FFF 0x0000000000040200\n",
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
			"0x00000000B880F000 0x00000000B880FFFF 0x0000000000040200\n",
			"0x0000000000011000 0x0000000000011FFF 0x0000000000040101\n",
		},
		{
			// Content for "pci0000:16/0000:16:02.0/0000:19:00.1" // Valid Bus
			"0x00000000B8810000 0x00000000B8813FFF 0x0000000000040200\n",
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
			"0x0000000000012000 0x0000000000012FFF 0x0000000000040101\n",
		},
		{
			// Content for "pci0000:16/0000:16:03.0" // Valid Bus
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
		},
		{
			// Content for "pci0000:1a/0000:1a:00.0" // Valid Bus with Invalid Resource
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
		},
		{
			// Content for "pci0000:1a/0000:1a:00.1" // Valid Bus with Invalid Resource
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
		},
		{
			// Content for "pci0000:1a/0000:1a:00.2" // Valid Bus with Invalid Resource
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
		},
		{
			// Content for "pci0000:1a/0000:1a:00.3" // Valid Bus with Invalid Resource
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
		},
		{
			// Content for "pci0000:1a/0000:1a:00.4" // Valid Bus with Invalid Resource
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
		},
		{
			// Content for "pci0000:1a/0000:1a:00.5" // Valid Bus with Invalid Resource
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
		},
		{
			// Content for "pci0000:1a/0000:1a:00.6" // Valid Bus with Invalid Resource
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
		},
		{
			// Content for "pci0000:1a/0000:1a:00.7" // Valid Bus with Invalid Resource
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
		},
		{
			// Content for "pci0000:c8/0000:c8:00.0" // Valid Bus
			"0x00000000B0800000 0x00000000B0FFFFFF 0x0000000000040200\n",
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
			"0x000000000000D000 0x000000000000D03F 0x0000000000040101\n",
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
			"0x0000000800F00000 0x0000000800FFFFFF 0x0000000000140204\n",
		},
		{
			// Content for "pci0000:e1/0000:e1:00.0" // Out of Bus range
			"0x00000000B0000000 0x00000000B0FFFFFF 0x0000000000040200\n",
			"0x0000000801000000 0x0000000801FFFFFF 0x0000000000140204\n",
		},
		{
			// Content for "pci0000:ff/0000:ff:00.0" // Out of Bus range
			"0x00000000A8000000 0x00000000AFFFFFFF 0x0000000000040200\n",
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
		},
		{
			// Content for "pci0002:00/0002:00:00.0" // Valid Bus
			"0x0000000810000000 0x000000081000FFFF 0x0000000000140204\n",
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
		},
		{
			// Content for "pci0002:00/0002:00:02.0" // Valid Bus
			"0x00000000A0070000 0x00000000A07FFFFF 0x0000000000040200\n",
			"0x000000000000B000 0x000000000000BFFF 0x0000000000040101\n",
			"0x00000000000C0000 0x00000000000DFFFF 0x0000000000000212\n",
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
			"0x0000000810010000 0x000000081001FFFF 0x0000000000140204\n",
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
		},
		{
			// Content for "pci0002:1f/0002:1f:00.0" // Valid Bus
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
			"0x00000000A0040000 0x00000000A004FFFF 0x0000000000040200\n",
			"0x00000000A0060000 0x00000000A006FFFF 0x0000000000040200\n",
			"0x00000000A0050000 0x00000000A005FFFF 0x0000000000140204\n",
		},
		{
			// Content for "pci0002:af/0002:af:00.0" // Valid Bus
			"0x00000000A0030000 0x00000000A0037FFF 0x0000000000040200\n",
			"0x00000000A0038000 0x00000000A003FFFF 0x0000000000140204\n",
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
		},
		{
			// Content for "pci0002:ff/0002:ff:00.0" // Valid Bus
			"0x00000000A0000000 0x00000000A000FFFF 0x0000000000040200\n",
			"0x000000000000A000 0x000000000000A03F 0x0000000000040101\n",
			"0x0000000000000000 0x0000000000000000 0x0000000000000000\n",
		},
	}

	expectedRbNodes := []*dt.Node{
		dt.NewNode("pci-rb", dt.WithProperty(
			dt.PropertyString("compatible", "pci-rb"),
			dt.PropertyU64("reg", 0xB000_0000),
			dt.PropertyU32Array("bus-range", []uint32{0x14, 0x14}),
			dt.PropertyU32Array("ranges", []uint32{
				0x300_0000,               // 64BITS
				0x0000_0008, 0x0000_0000, // MMIO64 Base high and low
				0x0, 0x0,
				0x0000_0000, 0x00F0_0000, // MMIO64 Limit high and low
				0x200_0000,               // 32BITS
				0x0000_0000, 0xC000_0000, // MMIO32 Base high and low
				0x0, 0x0,
				0x0000_0000, 0x1F00_0000, // MMIO32 Limit high and low
				0x100_0000,               // IOPort
				0x0000_0000, 0x0000_E000, // IOPort Base high and low
				0x0, 0x0,
				0x0000_0000, 0x0000_1000, // IOPort Limit high and low
			}),
		)),
		dt.NewNode("pci-rb", dt.WithProperty(
			dt.PropertyString("compatible", "pci-rb"),
			dt.PropertyU64("reg", 0xB000_0000),
			dt.PropertyU32Array("bus-range", []uint32{0x16, 0x19}),
			dt.PropertyU32Array("ranges", []uint32{
				0x300_0000,               // 64BITS
				0xFFFF_FFFF, 0xFFFF_FFFF, // MMIO64 Base high and low
				0x0, 0x0,
				0x0000_0000, 0x0000_0000, // MMIO64 Limit high and low
				0x200_0000,               // 32BITS
				0x0000_0000, 0xB880_0000, // MMIO32 Base high and low
				0x0, 0x0,
				0x0000_0000, 0x0001_4000, // MMIO32 Limit high and low
				0x100_0000,               // IOPort
				0x0000_0000, 0x0001_0000, // IOPort Base high and low
				0x0, 0x0,
				0x0000_0000, 0x0000_3000, // IOPort Limit high and low
			}),
		)),
		dt.NewNode("pci-rb", dt.WithProperty(
			dt.PropertyString("compatible", "pci-rb"),
			dt.PropertyU64("reg", 0xB000_0000),
			dt.PropertyU32Array("bus-range", []uint32{0x1A, 0x1A}),
			dt.PropertyU32Array("ranges", []uint32{
				0x300_0000,               // 64BITS
				0xFFFF_FFFF, 0xFFFF_FFFF, // MMIO64 Base high and low
				0x0, 0x0,
				0x0000_0000, 0x0000_0000, // MMIO64 Limit high and low
				0x200_0000,               // 32BITS
				0xFFFF_FFFF, 0xFFFF_FFFF, // MMIO32 Base high and low
				0x0, 0x0,
				0x0000_0000, 0x0000_0000, // MMIO32 Limit high and low
				0x100_0000,               // IOPort
				0xFFFF_FFFF, 0xFFFF_FFFF, // IOPort Base high and low
				0x0, 0x0,
				0x0000_0000, 0x0000_0000, // IOPort Limit high and low
			}),
		)),
		dt.NewNode("pci-rb", dt.WithProperty(
			dt.PropertyString("compatible", "pci-rb"),
			dt.PropertyU64("reg", 0xB000_0000),
			dt.PropertyU32Array("bus-range", []uint32{0xC8, 0xC8}),
			dt.PropertyU32Array("ranges", []uint32{
				0x300_0000,               // 64BITS
				0x0000_0008, 0x00F0_0000, // MMIO64 Base high and low
				0x0, 0x0,
				0x0000_0000, 0x0010_0000, // MMIO64 Limit high and low
				0x200_0000,               // 32BITS
				0x0000_0000, 0xB080_0000, // MMIO32 Base high and low
				0x0, 0x0,
				0x0000_0000, 0x0080_0000, // MMIO32 Limit high and low
				0x100_0000,               // IOPort
				0x0000_0000, 0x0000_D000, // IOPort Base high and low
				0x0, 0x0,
				0x0000_0000, 0x0000_0040, // IOPort Limit high and low
			}),
		)),
		dt.NewNode("pci-rb", dt.WithProperty(
			dt.PropertyString("compatible", "pci-rb"),
			dt.PropertyU64("reg", 0xA000_0000),
			dt.PropertyU32Array("bus-range", []uint32{0x00, 0x00}),
			dt.PropertyU32Array("ranges", []uint32{
				0x300_0000,               // 64BITS
				0x0000_0008, 0x1000_0000, // MMIO64 Base high and low
				0x0, 0x0,
				0x0000_0000, 0x0002_0000, // MMIO64 Limit high and low
				0x200_0000,               // 32BITS
				0x0000_0000, 0xA007_0000, // MMIO32 Base high and low
				0x0, 0x0,
				0x0000_0000, 0x0079_0000, // MMIO32 Limit high and low
				0x100_0000,               // IOPort
				0x0000_0000, 0x0000_B000, // IOPort Base high and low
				0x0, 0x0,
				0x0000_0000, 0x0000_1000, // IOPort Limit high and low
			}),
		)),
		dt.NewNode("pci-rb", dt.WithProperty(
			dt.PropertyString("compatible", "pci-rb"),
			dt.PropertyU64("reg", 0xA000_0000),
			dt.PropertyU32Array("bus-range", []uint32{0x1F, 0x1F}),
			dt.PropertyU32Array("ranges", []uint32{
				0x300_0000,               // 64BITS
				0xFFFF_FFFF, 0xFFFF_FFFF, // MMIO64 Base high and low
				0x0, 0x0,
				0x0000_0000, 0x0000_0000, // MMIO64 Limit high and low
				0x200_0000,               // 32BITS
				0x0000_0000, 0xA004_0000, // MMIO32 Base high and low
				0x0, 0x0,
				0x0000_0000, 0x0003_0000, // MMIO32 Limit high and low
				0x100_0000,               // IOPort
				0xFFFF_FFFF, 0xFFFF_FFFF, // IOPort Base high and low
				0x0, 0x0,
				0x0000_0000, 0x0000_0000, // IOPort Limit high and low
			}),
		)),
		dt.NewNode("pci-rb", dt.WithProperty(
			dt.PropertyString("compatible", "pci-rb"),
			dt.PropertyU64("reg", 0xA000_0000),
			dt.PropertyU32Array("bus-range", []uint32{0xAF, 0xAF}),
			dt.PropertyU32Array("ranges", []uint32{
				0x300_0000,               // 64BITS
				0xFFFF_FFFF, 0xFFFF_FFFF, // MMIO64 Base high and low
				0x0, 0x0,
				0x0000_0000, 0x0000_0000, // MMIO64 Limit high and low
				0x200_0000,               // 32BITS
				0x0000_0000, 0xA003_0000, // MMIO32 Base high and low
				0x0, 0x0,
				0x0000_0000, 0x0001_0000, // MMIO32 Limit high and low
				0x100_0000,               // IOPort
				0xFFFF_FFFF, 0xFFFF_FFFF, // IOPort Base high and low
				0x0, 0x0,
				0x0000_0000, 0x0000_0000, // IOPort Limit high and low
			}),
		)),
		dt.NewNode("pci-rb", dt.WithProperty(
			dt.PropertyString("compatible", "pci-rb"),
			dt.PropertyU64("reg", 0xA000_0000),
			dt.PropertyU32Array("bus-range", []uint32{0xFF, 0xFF}),
			dt.PropertyU32Array("ranges", []uint32{
				0x300_0000,               // 64BITS
				0xFFFF_FFFF, 0xFFFF_FFFF, // MMIO64 Base high and low
				0x0, 0x0,
				0x0000_0000, 0x0000_0000, // MMIO64 Limit high and low
				0x200_0000,               // 32BITS
				0x0000_0000, 0xA000_0000, // MMIO32 Base high and low
				0x0, 0x0,
				0x0000_0000, 0x0001_0000, // MMIO32 Limit high and low
				0x100_0000,               // IOPort
				0x0000_0000, 0x0000_A000, // IOPort Base high and low
				0x0, 0x0,
				0x0000_0000, 0x0000_0040, // IOPort Limit high and low
			}),
		)),
	}

	// Create temp folder for testing
	tmpDir := t.TempDir()

	// Write data to files:
	// /$TMPDIR/$DOMAIN_ID:$BUS_ID:$DEVICE_ID.$FUNCTION_ID/resource
	for i, folderName := range subFolder {
		subFolderPath := filepath.Join(tmpDir, folderName)
		if err := os.MkdirAll(subFolderPath, 0o755); err != nil {
			t.Fatalf("Error creating subfolder %s: %v\n", subFolderPath, err)
		}

		filePath := filepath.Join(subFolderPath, "resource")
		data := strings.Join(resourceContent[i], "")

		if err := os.WriteFile(filePath, []byte(data), 0o644); err != nil {
			t.Fatalf("Error writing to file %s: %v\n", filePath, err)
		}
	}

	var idx uint32
	for _, item := range mcfgData {
		rbNodes, err := createPCIRootBridgeNode(tmpDir, item)
		if err != nil {
			t.Fatalf("Failed to create RB node %v\n", err)
		}

		for _, rbNode := range rbNodes {
			expected := expectedRbNodes[idx]
			diff := cmp.Diff(expected, rbNode)
			if diff != "" {
				t.Errorf("Index:%x Mismatch (-expected +rbNode):\n%s\n", idx, diff)
			}
			idx++
		}
	}
}

func TestGetReservedMemoryMap(t *testing.T) {
	tests := []struct {
		name           string
		inputMemoryMap kexec.MemoryMap
		expectedResult kexec.MemoryMap
		expectedError  error
	}{
		{
			name: "Success with reserved memory regions",
			inputMemoryMap: kexec.MemoryMap{
				kexec.TypedRange{Range: kexec.Range{Start: 0x1000, Size: 0x400000}, Type: kexec.RangeRAM},
				kexec.TypedRange{Range: kexec.Range{Start: 0x500000, Size: 0x100000}, Type: kexec.RangeReserved},
				kexec.TypedRange{Range: kexec.Range{Start: 0x600000, Size: 0x200000}, Type: kexec.RangeACPI},
				kexec.TypedRange{Range: kexec.Range{Start: 0x800000, Size: 0x50000}, Type: kexec.RangeReserved},
			},
			expectedResult: kexec.MemoryMap{
				kexec.TypedRange{Range: kexec.Range{Start: 0x500000, Size: 0x100000}, Type: kexec.RangeReserved},
				kexec.TypedRange{Range: kexec.Range{Start: 0x800000, Size: 0x50000}, Type: kexec.RangeReserved},
			},
			expectedError: nil,
		},
		{
			name: "No reserved memory regions",
			inputMemoryMap: kexec.MemoryMap{
				kexec.TypedRange{Range: kexec.Range{Start: 0x1000, Size: 0x400000}, Type: kexec.RangeRAM},
				kexec.TypedRange{Range: kexec.Range{Start: 0x500000, Size: 0x100000}, Type: kexec.RangeACPI},
			},
			expectedResult: nil,
			expectedError:  nil,
		},
		{
			name:           "Empty memory map",
			inputMemoryMap: kexec.MemoryMap{},
			expectedResult: nil,
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function under test
			result, err := getReservedMemoryMap(tt.inputMemoryMap)

			// Check error
			if tt.expectedError != nil {
				if err == nil {
					t.Fatalf("Expected error %q, got nil", tt.expectedError)
				}
				if !errors.Is(err, tt.expectedError) {
					t.Errorf("Unexpected error %q, want = %q", err.Error(), tt.expectedError)
				}
			} else if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Check result
			if tt.expectedResult == nil {
				if result != nil {
					t.Fatalf("Expected nil result, got %v", result)
				}
			} else {
				if result == nil {
					t.Fatalf("Expected result %v, got nil", tt.expectedResult)
				}

				// Compare the memory maps
				if len(result) != len(tt.expectedResult) {
					t.Fatalf("Expected %d reserved memory regions, got %d", len(tt.expectedResult), len(result))
				}

				for i, expectedRegion := range tt.expectedResult {
					if i >= len(result) {
						t.Fatalf("Missing memory region at index %d", i)
					}

					actualRegion := result[i]
					if expectedRegion.Range.Start != actualRegion.Range.Start {
						t.Errorf("Memory region %d: expected Start = 0x%x, got 0x%x", i, expectedRegion.Range.Start, actualRegion.Range.Start)
					}
					if expectedRegion.Range.Size != actualRegion.Range.Size {
						t.Errorf("Memory region %d: expected Size = 0x%x, got 0x%x", i, expectedRegion.Range.Size, actualRegion.Range.Size)
					}
					if expectedRegion.Type != actualRegion.Type {
						t.Errorf("Memory region %d: expected Type = %v, got %v", i, expectedRegion.Type, actualRegion.Type)
					}
				}
			}
		})
	}
}

func TestSkipReservedRange(t *testing.T) {
	// Create a test memory map with reserved regions only
	testMemoryMap := kexec.MemoryMap{
		kexec.TypedRange{Range: kexec.Range{Start: 0x500000, Size: 0x100000}, Type: kexec.RangeReserved},
		kexec.TypedRange{Range: kexec.Range{Start: 0x800000, Size: 0x50000}, Type: kexec.RangeReserved},
	}

	tests := []struct {
		name           string
		memoryMap      kexec.MemoryMap
		base           uintptr
		attr           uint64
		expectedResult bool
		description    string
	}{
		{
			name:           "IOPort resource should not skip",
			memoryMap:      testMemoryMap,
			base:           0x1000,
			attr:           0x100, // PCIIOPortRes
			expectedResult: false,
			description:    "IOPort resources should never be skipped regardless of memory map",
		},
		{
			name:           "IOPort resource with other bitsshould not skip",
			memoryMap:      testMemoryMap,
			base:           0x500000,
			attr:           0x40100, // PCIIOPortAttr (includes PCIIOPortRes)
			expectedResult: false,
			description:    "IOPort resources should never be skipped even if in reserved memory",
		},
		{
			name:           "ReadOnly MMIO should skip",
			memoryMap:      testMemoryMap,
			base:           0x1000,
			attr:           0x4000, // PCIMMIOReadOnly
			expectedResult: true,
			description:    "ReadOnly MMIO should always be skipped",
		},
		{
			name:           "ReadOnly MMIO with other bits should skip",
			memoryMap:      testMemoryMap,
			base:           0x1000,
			attr:           0x44000, // PCIMMIOReadOnly + other bits
			expectedResult: true,
			description:    "ReadOnly MMIO should always be skipped even with other attribute bits",
		},
		{
			name:           "Base in reserved memory region should skip",
			memoryMap:      testMemoryMap,
			base:           0x500000,
			attr:           0x40200, // PCIMMIO32Attr
			expectedResult: true,
			description:    "Base address within reserved memory region should be skipped",
		},
		{
			name:           "Base at start of reserved memory region should skip",
			memoryMap:      testMemoryMap,
			base:           0x500000,
			attr:           0x40200, // PCIMMIO32Attr
			expectedResult: true,
			description:    "Base address at start of reserved memory region should be skipped",
		},
		{
			name:           "Base at end of reserved memory region should skip",
			memoryMap:      testMemoryMap,
			base:           0x5FFFFF, // 0x500000 + 0x100000 - 1
			attr:           0x40200,  // PCIMMIO32Attr
			expectedResult: true,
			description:    "Base address at end of reserved memory region should be skipped",
		},
		{
			name:           "Base in middle of reserved memory region should skip",
			memoryMap:      testMemoryMap,
			base:           0x550000, // Middle of 0x500000-0x600000 range
			attr:           0x40200,  // PCIMMIO32Attr
			expectedResult: true,
			description:    "Base address in middle of reserved memory region should be skipped",
		},
		{
			name:           "Base in second reserved memory region should skip",
			memoryMap:      testMemoryMap,
			base:           0x800000,
			attr:           0x40200, // PCIMMIO32Attr
			expectedResult: true,
			description:    "Base address within second reserved memory region should be skipped",
		},
		{
			name:           "Base not in any reserved memory region should not skip",
			memoryMap:      testMemoryMap,
			base:           0x1000,
			attr:           0x40200, // PCIMMIO32Attr
			expectedResult: false,
			description:    "Base address not in any reserved memory region should not be skipped",
		},
		{
			name:           "Base between reserved regions should not skip",
			memoryMap:      testMemoryMap,
			base:           0x700000, // Between 0x500000-0x600000 and 0x800000-0x850000
			attr:           0x40200,  // PCIMMIO32Attr
			expectedResult: false,
			description:    "Base address between reserved memory regions should not be skipped",
		},
		{
			name:           "Base outside all reserved memory regions should not skip",
			memoryMap:      testMemoryMap,
			base:           0x900000,
			attr:           0x40200, // PCIMMIO32Attr
			expectedResult: false,
			description:    "Base address outside all reserved memory regions should not be skipped",
		},
		{
			name:           "Empty memory map should not skip",
			memoryMap:      kexec.MemoryMap{},
			base:           0x500000,
			attr:           0x40200, // PCIMMIO32Attr
			expectedResult: false,
			description:    "With empty memory map, should not skip any addresses",
		},
		{
			name:           "Base at boundary of reserved region should not skip",
			memoryMap:      testMemoryMap,
			base:           0x600000, // End of first reserved region (0x500000 + 0x100000)
			attr:           0x40200,  // PCIMMIO32Attr
			expectedResult: false,
			description:    "Base address at boundary of reserved region should not be skipped if not in reserved region",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := skipReservedRange(tt.memoryMap, tt.base, tt.attr)

			if result != tt.expectedResult {
				t.Errorf("skipReservedRange:\n%v, (base: 0x%x, attr: 0x%x) = %v, want %v\nDescription: %s",
					tt.memoryMap, tt.base, tt.attr, result, tt.expectedResult, tt.description)
			}
		})
	}
}

func TestIsValidPCIDeviceName(t *testing.T) {
	tests := []struct {
		name           string
		deviceName     string
		expectedResult bool
		description    string
	}{
		// Positive test cases - valid PCI device names
		{
			name:           "Valid PCI device name with all zeros",
			deviceName:     "0000:00:00.0",
			expectedResult: true,
			description:    "Standard valid PCI device name with all zero values",
		},
		{
			name:           "Valid PCI device name with hex values",
			deviceName:     "0001:0a:1f.3",
			expectedResult: true,
			description:    "Valid PCI device name with non-zero hex values",
		},
		{
			name:           "Valid PCI device name with maximum values",
			deviceName:     "ffff:ff:ff.f",
			expectedResult: true,
			description:    "Valid PCI device name with maximum hex values",
		},
		{
			name:           "Valid PCI device name with mixed case hex",
			deviceName:     "aBcD:1f:2e.5",
			expectedResult: true,
			description:    "Valid PCI device name with mixed case hex characters",
		},
		{
			name:           "Valid PCI device name with function 0",
			deviceName:     "0000:01:02.0",
			expectedResult: true,
			description:    "Valid PCI device name with function number 0",
		},
		{
			name:           "Valid PCI device name with function 7",
			deviceName:     "0000:01:02.7",
			expectedResult: true,
			description:    "Valid PCI device name with function number 7",
		},
		{
			name:           "Valid PCI device name with function f",
			deviceName:     "0000:01:02.f",
			expectedResult: true,
			description:    "Valid PCI device name with function number f (15)",
		},

		// Negative test cases - invalid PCI device names
		{
			name:           "Empty string",
			deviceName:     "",
			expectedResult: false,
			description:    "Empty string should be invalid",
		},
		{
			name:           "Too short - missing parts",
			deviceName:     "0000:00",
			expectedResult: false,
			description:    "Device name too short, missing device.function part",
		},
		{
			name:           "Too long - extra parts",
			deviceName:     "0000:00:00.0:extra",
			expectedResult: false,
			description:    "Device name too long, has extra parts",
		},
		{
			name:           "Wrong length - 11 characters",
			deviceName:     "000:00:00.0",
			expectedResult: false,
			description:    "Device name with wrong total length (11 instead of 12)",
		},
		{
			name:           "Wrong length - 13 characters",
			deviceName:     "0000:00:00.00",
			expectedResult: false,
			description:    "Device name with wrong total length (13 instead of 12)",
		},
		{
			name:           "Domain too short",
			deviceName:     "000:00:00.0",
			expectedResult: false,
			description:    "Domain part too short (3 chars instead of 4)",
		},
		{
			name:           "Domain too long",
			deviceName:     "00000:00:00.0",
			expectedResult: false,
			description:    "Domain part too long (5 chars instead of 4)",
		},
		{
			name:           "Bus too short",
			deviceName:     "0000:0:00.0",
			expectedResult: false,
			description:    "Bus part too short (1 char instead of 2)",
		},
		{
			name:           "Bus too long",
			deviceName:     "0000:000:00.0",
			expectedResult: false,
			description:    "Bus part too long (3 chars instead of 2)",
		},
		{
			name:           "Device part too short",
			deviceName:     "0000:00:0.0",
			expectedResult: false,
			description:    "Device part too short (1 char instead of 2)",
		},
		{
			name:           "Device part too long",
			deviceName:     "0000:00:000.0",
			expectedResult: false,
			description:    "Device part too long (3 chars instead of 2)",
		},
		{
			name:           "Function part too short",
			deviceName:     "0000:00:00.",
			expectedResult: false,
			description:    "Function part missing (0 chars instead of 1)",
		},
		{
			name:           "Function part too long",
			deviceName:     "0000:00:00.00",
			expectedResult: false,
			description:    "Function part too long (2 chars instead of 1)",
		},
		{
			name:           "Missing first colon",
			deviceName:     "000000:00.0",
			expectedResult: false,
			description:    "Missing first colon separator",
		},
		{
			name:           "Missing second colon",
			deviceName:     "0000:0000.0",
			expectedResult: false,
			description:    "Missing second colon separator",
		},
		{
			name:           "Missing dot separator",
			deviceName:     "0000:00:000",
			expectedResult: false,
			description:    "Missing dot separator between device and function",
		},
		{
			name:           "Extra separators",
			deviceName:     "0000::00:00.0",
			expectedResult: false,
			description:    "Extra colon separator",
		},
		{
			name:           "PCI bridge device info",
			deviceName:     "0000:00:1c.0:pcie002",
			expectedResult: false,
			description:    "Not a valid PCI device name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidPCIDeviceName(tt.deviceName)

			if result != tt.expectedResult {
				t.Errorf("isValidPCIDeviceName(%q) = %v, want %v\nDescription: %s",
					tt.deviceName, result, tt.expectedResult, tt.description)
			}
		})
	}
}

func TestUpdateResourceRanges(t *testing.T) {
	tests := []struct {
		name           string
		initialRegion  *ResourceRegions
		resType        string
		base           uint64
		end            uint64
		expectedRegion *ResourceRegions
		description    string
	}{
		// MMIO64 tests
		{
			name: "MMIO64 First resource with invalid base",
			initialRegion: &ResourceRegions{
				MMIO64Base: PCIInvalidBase,
				MMIO64End:  0,
				MMIO32Base: PCIInvalidBase,
				MMIO32End:  0,
				IOPortBase: PCIInvalidBase,
				IOPortEnd:  0,
			},
			resType: PCIMMIO64Type,
			base:    0x1000000,
			end:     0x1000FFF,
			expectedRegion: &ResourceRegions{
				MMIO64Base: 0x1000000,
				MMIO64End:  0x1000FFF, // align.UpPage(0x1000FFF) - 1 = 0x1001FFF
				MMIO32Base: PCIInvalidBase,
				MMIO32End:  0,
				IOPortBase: PCIInvalidBase,
				IOPortEnd:  0,
			},
			description: "First MMIO64 resource should set both base and end",
		},
		{
			name: "MMIO64 Update base to lower value",
			initialRegion: &ResourceRegions{
				MMIO64Base: 0x2000000,
				MMIO64End:  0x2001FFF,
				MMIO32Base: PCIInvalidBase,
				MMIO32End:  0,
				IOPortBase: PCIInvalidBase,
				IOPortEnd:  0,
			},
			resType: PCIMMIO64Type,
			base:    0x1000000,
			end:     0x1000FFF,
			expectedRegion: &ResourceRegions{
				MMIO64Base: 0x1000000, // min(0x1000000, 0x2000000)
				MMIO64End:  0x2001FFF, // max(0x1001FFF, 0x2001FFF)
				MMIO32Base: PCIInvalidBase,
				MMIO32End:  0,
				IOPortBase: PCIInvalidBase,
				IOPortEnd:  0,
			},
			description: "MMIO64 base should be updated to lower value, end to higher value",
		},
		{
			name: "MMIO64 Update end to higher value",
			initialRegion: &ResourceRegions{
				MMIO64Base: 0x1000000,
				MMIO64End:  0x1001FFF,
				MMIO32Base: PCIInvalidBase,
				MMIO32End:  0,
				IOPortBase: PCIInvalidBase,
				IOPortEnd:  0,
			},
			resType: PCIMMIO64Type,
			base:    0x2000000,
			end:     0x2000FFF,
			expectedRegion: &ResourceRegions{
				MMIO64Base: 0x1000000, // min(0x2000000, 0x1000000)
				MMIO64End:  0x2000FFF, // max(0x2000FFF, 0x1001FFF)
				MMIO32Base: PCIInvalidBase,
				MMIO32End:  0,
				IOPortBase: PCIInvalidBase,
				IOPortEnd:  0,
			},
			description: "MMIO64 end should be updated to higher value, base remains lower",
		},
		{
			name: "MMIO64 Page alignment test",
			initialRegion: &ResourceRegions{
				MMIO64Base: PCIInvalidBase,
				MMIO64End:  0,
				MMIO32Base: PCIInvalidBase,
				MMIO32End:  0,
				IOPortBase: PCIInvalidBase,
				IOPortEnd:  0,
			},
			resType: PCIMMIO64Type,
			base:    0x1000000,
			end:     0x1000ABC, // Not page aligned
			expectedRegion: &ResourceRegions{
				MMIO64Base: 0x1000000,
				MMIO64End:  0x1000FFF, // align.UpPage(0x1000ABC) - 1
				MMIO32Base: PCIInvalidBase,
				MMIO32End:  0,
				IOPortBase: PCIInvalidBase,
				IOPortEnd:  0,
			},
			description: "MMIO64 end should be page aligned up then decremented by 1",
		},

		// MMIO32 tests
		{
			name: "MMIO32 First resource with invalid base",
			initialRegion: &ResourceRegions{
				MMIO64Base: PCIInvalidBase,
				MMIO64End:  0,
				MMIO32Base: PCIInvalidBase,
				MMIO32End:  0,
				IOPortBase: PCIInvalidBase,
				IOPortEnd:  0,
			},
			resType: PCIMMIO32Type,
			base:    0x8000000,
			end:     0x8000FFF,
			expectedRegion: &ResourceRegions{
				MMIO64Base: PCIInvalidBase,
				MMIO64End:  0,
				MMIO32Base: 0x8000000,
				MMIO32End:  0x8000FFF, // align.UpPage(0x8000FFF) - 1
				IOPortBase: PCIInvalidBase,
				IOPortEnd:  0,
			},
			description: "First MMIO32 resource should set both base and end",
		},
		{
			name: "MMIO32 Update base to lower value",
			initialRegion: &ResourceRegions{
				MMIO64Base: PCIInvalidBase,
				MMIO64End:  0,
				MMIO32Base: 0x9000000,
				MMIO32End:  0x9001FFF,
				IOPortBase: PCIInvalidBase,
				IOPortEnd:  0,
			},
			resType: PCIMMIO32Type,
			base:    0x8000000,
			end:     0x8000FFF,
			expectedRegion: &ResourceRegions{
				MMIO64Base: PCIInvalidBase,
				MMIO64End:  0,
				MMIO32Base: 0x8000000,
				MMIO32End:  0x9001FFF,
				IOPortBase: PCIInvalidBase,
				IOPortEnd:  0,
			},
			description: "MMIO32 base should be updated to lower value, end to higher value",
		},

		// IOPort tests
		{
			name: "IOPort First resource with invalid base",
			initialRegion: &ResourceRegions{
				MMIO64Base: PCIInvalidBase,
				MMIO64End:  0,
				MMIO32Base: PCIInvalidBase,
				MMIO32End:  0,
				IOPortBase: PCIInvalidBase,
				IOPortEnd:  0,
			},
			resType: PCIIOPortType,
			base:    0x1000,
			end:     0x10FF,
			expectedRegion: &ResourceRegions{
				MMIO64Base: PCIInvalidBase,
				MMIO64End:  0,
				MMIO32Base: PCIInvalidBase,
				MMIO32End:  0,
				IOPortBase: 0x1000,
				IOPortEnd:  0x10FF, // No page alignment for IOPort
			},
			description: "First IOPort resource should set both base and end without page alignment",
		},
		{
			name: "IOPort Update base to lower value",
			initialRegion: &ResourceRegions{
				MMIO64Base: PCIInvalidBase,
				MMIO64End:  0,
				MMIO32Base: PCIInvalidBase,
				MMIO32End:  0,
				IOPortBase: 0x2000,
				IOPortEnd:  0x20FF,
			},
			resType: PCIIOPortType,
			base:    0x1000,
			end:     0x10FF,
			expectedRegion: &ResourceRegions{
				MMIO64Base: PCIInvalidBase,
				MMIO64End:  0,
				MMIO32Base: PCIInvalidBase,
				MMIO32End:  0,
				IOPortBase: 0x1000, // min(0x1000, 0x2000)
				IOPortEnd:  0x20FF, // max(0x10FF, 0x20FF)
			},
			description: "IOPort base should be updated to lower value, end to higher value",
		},
		{
			name: "IOPort Update end to higher value",
			initialRegion: &ResourceRegions{
				MMIO64Base: PCIInvalidBase,
				MMIO64End:  0,
				MMIO32Base: PCIInvalidBase,
				MMIO32End:  0,
				IOPortBase: 0x1000,
				IOPortEnd:  0x10FF,
			},
			resType: PCIIOPortType,
			base:    0x2000,
			end:     0x20FF,
			expectedRegion: &ResourceRegions{
				MMIO64Base: PCIInvalidBase,
				MMIO64End:  0,
				MMIO32Base: PCIInvalidBase,
				MMIO32End:  0,
				IOPortBase: 0x1000,
				IOPortEnd:  0x20FF,
			},
			description: "IOPort end should be updated to higher value, base remains lower",
		},

		// Unknown resource type test
		{
			name: "Unknown resource type",
			initialRegion: &ResourceRegions{
				MMIO64Base: 0x1000000,
				MMIO64End:  0x1001FFF,
				MMIO32Base: 0x8000000,
				MMIO32End:  0x8001FFF,
				IOPortBase: 0x1000,
				IOPortEnd:  0x10FF,
			},
			resType: "UNKNOWN",
			base:    0x3000000,
			end:     0x3000FFF,
			expectedRegion: &ResourceRegions{
				MMIO64Base: 0x1000000, // Unchanged
				MMIO64End:  0x1001FFF, // Unchanged
				MMIO32Base: 0x8000000, // Unchanged
				MMIO32End:  0x8001FFF, // Unchanged
				IOPortBase: 0x1000,    // Unchanged
				IOPortEnd:  0x10FF,    // Unchanged
			},
			description: "Unknown resource type should not modify any fields",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a copy of the initial region to avoid modifying the test data
			resourceRegion := &ResourceRegions{
				MMIO64Base: tt.initialRegion.MMIO64Base,
				MMIO64End:  tt.initialRegion.MMIO64End,
				MMIO32Base: tt.initialRegion.MMIO32Base,
				MMIO32End:  tt.initialRegion.MMIO32End,
				IOPortBase: tt.initialRegion.IOPortBase,
				IOPortEnd:  tt.initialRegion.IOPortEnd,
				StartBus:   tt.initialRegion.StartBus,
				EndBus:     tt.initialRegion.EndBus,
			}

			// Call the function under test
			updateResourceRanges(resourceRegion, tt.resType, tt.base, tt.end)

			// Compare results
			if resourceRegion.MMIO64Base != tt.expectedRegion.MMIO64Base {
				t.Errorf("MMIO64Base = 0x%x, want 0x%x\nDescription: %s",
					resourceRegion.MMIO64Base, tt.expectedRegion.MMIO64Base, tt.description)
			}
			if resourceRegion.MMIO64End != tt.expectedRegion.MMIO64End {
				t.Errorf("MMIO64End = 0x%x, want 0x%x\nDescription: %s",
					resourceRegion.MMIO64End, tt.expectedRegion.MMIO64End, tt.description)
			}
			if resourceRegion.MMIO32Base != tt.expectedRegion.MMIO32Base {
				t.Errorf("MMIO32Base = 0x%x, want 0x%x\nDescription: %s",
					resourceRegion.MMIO32Base, tt.expectedRegion.MMIO32Base, tt.description)
			}
			if resourceRegion.MMIO32End != tt.expectedRegion.MMIO32End {
				t.Errorf("MMIO32End = 0x%x, want 0x%x\nDescription: %s",
					resourceRegion.MMIO32End, tt.expectedRegion.MMIO32End, tt.description)
			}
			if resourceRegion.IOPortBase != tt.expectedRegion.IOPortBase {
				t.Errorf("IOPortBase = 0x%x, want 0x%x\nDescription: %s",
					resourceRegion.IOPortBase, tt.expectedRegion.IOPortBase, tt.description)
			}
			if resourceRegion.IOPortEnd != tt.expectedRegion.IOPortEnd {
				t.Errorf("IOPortEnd = 0x%x, want 0x%x\nDescription: %s",
					resourceRegion.IOPortEnd, tt.expectedRegion.IOPortEnd, tt.description)
			}
		})
	}
}
