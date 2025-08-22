// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package universalpayload

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
	"testing"
	"unsafe"

	guid "github.com/google/uuid"
	"github.com/u-root/u-root/pkg/acpi"
	"github.com/u-root/u-root/pkg/align"
	"github.com/u-root/u-root/pkg/boot/kexec"
)

func mockKexecMemoryMapFromIOMem() (kexec.MemoryMap, error) {
	return kexec.MemoryMap{
		{Range: kexec.Range{Start: 0x1000, Size: 0x400000}, Type: kexec.RangeRAM},
		{Range: kexec.Range{Start: 100, Size: 50}, Type: kexec.RangeACPI},
	}, nil
}

func mockGetRSDP() (uint64, []byte, error) {
	fakeRDSP := acpi.RSDP{}
	return 0, fakeRDSP.AllData(), nil
}

func mockGetSMBIOSBase() (int64, int64, error) {
	return 100, 200, nil
}

// Main test for constructing HOB list,
// to make sure there are no issues with the overall process
func TestLoadKexecMemWithHOBs(t *testing.T) {
	// mock data
	defer func(old func() (int64, int64, error)) { getSMBIOSBase = old }(getSMBIOSBase)
	getSMBIOSBase = mockGetSMBIOSBase

	defer func(old func() (uint64, []byte, error)) { getAcpiRsdpData = old }(getAcpiRsdpData)
	getAcpiRsdpData = mockGetRSDP

	defer func(old func() (kexec.MemoryMap, error)) { kexecMemoryMapFromIOMem = old }(kexecMemoryMapFromIOMem)
	kexecMemoryMapFromIOMem = mockKexecMemoryMapFromIOMem

	tempFile := mockCPUTempInfoFile(t, `
processor	: 0
vendor_id	: GenuineIntel
cpu family	: 6
model		: 142
model name	: Intel(R) Core(TM) i7-8550U CPU @ 1.80GHz
address sizes	: 39 bits physical, 48 bits virtual
`)
	defer os.Remove(tempFile)
	// end of mock data

	// Follow components layout which is defined in utilities.go to place corresponding components
	// |------------------------| <-- Memory Region top
	// |     TRAMPOLINE CODE    |
	// |------------------------| <-- loadAddr + trampolineOffset
	// |      TEMP STACK        |
	// |------------------------| <-- loadAddr + tmpStackOffset
	// |    Device Tree Info    |
	// |------------------------| <-- loadAddr + fdtDtbOffset
	// |  BOOTLOADER PARAMETER  |
	// |  HoBs (Handoff Blocks) |
	// |------------------------| <-- loadAddr + tmpHobOffset
	// |       ACPI DATA        |
	// |------------------------| <-- loadAddr + rsdpTableOffset
	// |     UPL FIT IMAGE      |
	// |------------------------| <-- loadAddr which is 2MB aligned
	// 1 Page size for each component is enough in out positive test case.
	fdtInfo := &FdtLoad{DataOffset: 0x100, DataSize: 0x5000, EntryStart: 0x1000, Load: 0x1800}
	imgSize := fdtInfo.DataOffset + fdtInfo.DataSize
	rsdpOff := (uintptr)(align.UpPage(imgSize))
	hobsOff := rsdpOff + uintptr(pageSize)
	fdtOff := hobsOff + uintptr(pageSize)
	stackOff := fdtOff + uintptr(pageSize)
	trampOff := stackOff + uintptr(pageSize)

	tests := []struct {
		name            string
		fdtLoad         *FdtLoad
		data            []byte
		mem             kexec.Memory
		wantErr         error
		wantAddr        uintptr
		wantMemSegments []kexec.Range
	}{
		{
			name:    "Valid case to relocate FIT image",
			fdtLoad: fdtInfo,
			data: mockWritePeFileBinary(0x100, 0x5100, 0x1000, []*MockSection{
				{".reloc", mockRelocData(0x1000, IMAGE_REL_BASED_DIR64, 0x200)},
			}),
			mem: kexec.Memory{
				Phys: kexec.MemoryMap{
					{Range: kexec.Range{Start: 0x1000, Size: 0x500000}, Type: kexec.RangeRAM},
					{Range: kexec.Range{Start: 100, Size: 50}, Type: kexec.RangeACPI},
				},
			},
			wantErr:  nil,
			wantAddr: 0x200000 + 0xa000,
			wantMemSegments: []kexec.Range{
				{Start: 0x200000, Size: (uint)(imgSize)}, // PeFileBinary
				{Start: 0x200000 + rsdpOff, Size: 0},     // ACPI Data
				{Start: 0x200000 + hobsOff, Size: 0},     // HOBs for bootloader
				{Start: 0x200000 + fdtOff, Size: 0},      // Device Tree Info
				{Start: 0x200000 + stackOff, Size: 0},    // boot env (tmp stack)
				{Start: 0x200000 + trampOff, Size: 0},    // boot env (trampoline)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := loadKexecMemWithHOBs(tt.fdtLoad, tt.data, &tt.mem)
			expectErr(t, err, tt.wantErr)
			t.Logf("mem Segments: %v, want Segments: %v", tt.mem.Segments.Phys(), tt.wantMemSegments)
			if tt.wantAddr != addr {
				t.Fatalf("Unexpected target addr: %v, want: %v", addr, tt.wantAddr)
			}
			if tt.wantMemSegments != nil {
				if len(tt.wantMemSegments) != len(tt.mem.Segments.Phys()) {
					t.Fatalf("Unexpected mem segment size: %v, want: %v", len(tt.mem.Segments.Phys()), len(tt.wantMemSegments))
				}
				for i, r := range tt.wantMemSegments {
					actualRange := tt.mem.Segments.Phys()[i]
					if r.Start != actualRange.Start {
						t.Fatalf("Unexpected mem segment Start: %v, want: %v", actualRange.Start, r.Start)
					}
				}
			}
		})
	}
}

func TestAppendMemMapHOB(t *testing.T) {
	defer func(old func() (kexec.MemoryMap, error)) { kexecMemoryMapFromIOMem = old }(kexecMemoryMapFromIOMem)
	kexecMemoryMapFromIOMem = mockKexecMemoryMapFromIOMem

	var memMap kexec.MemoryMap
	if ioMem, err := kexecMemoryMapFromIOMem(); err != nil {
		t.Fatal(err)
	} else {
		memMap = ioMem
	}

	hobBuf := &bytes.Buffer{}
	var hobLen uint64
	if err := appendMemMapHOB(hobBuf, &hobLen, memMap); err != nil {
		t.Fatal(err)
	}

	var deserializedHOB EFIMemoryMapHOB
	for hobBuf.Len() > 0 {
		var hob EFIHOBResourceDescriptor
		err := binary.Read(hobBuf, binary.LittleEndian, &hob)
		if err != nil {
			t.Fatalf("Unexpected error: %+v", err)
		}
		deserializedHOB = append(deserializedHOB, hob)
	}

	// We will pass all non system memory regions info to UPL, update to the actual
	// memory region numbers provided in test case.
	if len(deserializedHOB) != 1 {
		t.Fatalf("Unexpected hob size = %d, want = %d", len(deserializedHOB), 1)
	}
}

func TestAppendSerialPortHOB(t *testing.T) {
	hobBuf := &bytes.Buffer{}
	var hobLen uint64
	if err := appendSerialPortHOB(hobBuf, &hobLen); err != nil {
		t.Fatal(err)
	}

	var serialPortInfo UniversalPayloadSerialPortInfo
	var efiHOBGUID EFIHOBGUIDType
	if err := binary.Read(hobBuf, binary.LittleEndian, &efiHOBGUID); err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}
	if err := binary.Read(hobBuf, binary.LittleEndian, &serialPortInfo); err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}
	if efiHOBGUID.Name.String() != UniversalPayloadSerialPortInfoGUID {
		t.Fatalf("Unexpected efiHOBGuid = %v, want = %v",
			efiHOBGUID.Name.String(), UniversalPayloadSerialPortInfoGUID)
	}
}

func TestAppendUniversalPayloadBase(t *testing.T) {
	hobBuf := &bytes.Buffer{}
	var hobLen uint64
	const fdtLoad = uint64(0x700000)
	if err := appendUniversalPayloadBase(hobBuf, &hobLen, fdtLoad); err != nil {
		t.Fatal(err)
	}
	var uplBase UniversalPayloadBase
	var uplBaseGUIDHOB EFIHOBGUIDType
	if err := binary.Read(hobBuf, binary.LittleEndian, &uplBaseGUIDHOB); err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}
	if err := binary.Read(hobBuf, binary.LittleEndian, &uplBase); err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}
	if hobBuf.Len() > 0 {
		t.Fatalf("Unexpected hobBuf size after deserialization: %v, want = 0", hobBuf.Len())
	}
	if uplBaseGUIDHOB.Name.String() != UniversalPayloadBaseGUID {
		t.Fatalf("Unexpected uplBaseGUIDHOB = %v, want = %v",
			uplBaseGUIDHOB.Name.String(), UniversalPayloadBaseGUID)
	}
	if uplBase.Entry != EFIPhysicalAddress(fdtLoad) {
		t.Fatalf("Unexpected UniversalPayloadBase.Entry: %v, want = %v", uplBase.Entry, fdtLoad)
	}
}

func TestConstructHOBList(t *testing.T) {
	const lenOfEFITable = uint64(unsafe.Sizeof(EFIHOBHandoffInfoTable{}))

	for _, tt := range []struct {
		// mock data
		name   string
		dst    *bytes.Buffer
		src    *bytes.Buffer
		hobLen uint64

		// expected data
		wantHOBLen uint64
		efiTable   EFIHOBHandoffInfoTable
		endHeader  EFIHOBGenericHeader
		wantErr    error
	}{
		{
			name:   "CASE 1: success",
			dst:    &bytes.Buffer{},
			src:    bytes.NewBuffer([]byte("12345678")),
			hobLen: uint64(8),

			wantHOBLen: uint64(8 + unsafe.Sizeof(EFIHOBGenericHeader{})),
			efiTable: EFIHOBHandoffInfoTable{
				Header: EFIHOBGenericHeader{
					HOBType:   EFIHOBTypeHandoff,
					HOBLength: EFIHOBLength(unsafe.Sizeof(EFIHOBHandoffInfoTable{})),
				},
				Version:          EFIHOBHandoffInfoVersion,
				BootMode:         EFIHOBHandoffInfoBootMode,
				MemoryTop:        EFIHOBHandoffInfoEFIMemoryTop,
				MemoryBottom:     EFIHOBHandoffInfoEFIMemoryBottom,
				FreeMemoryTop:    EFIHOBHandoffInfoFreeEFIMemoryTop,
				FreeMemoryBottom: EFIHOBHandoffInfoFreeMemoryBottom + EFIPhysicalAddress(8+8+lenOfEFITable),
				EndOfHOBList:     EFIHOBHandoffInfoFreeMemoryBottom + EFIPhysicalAddress(8+lenOfEFITable),
			},
			endHeader: EFIHOBGenericHeader{
				HOBType:   EFIHOBTypeEndOfHOBList,
				HOBLength: EFIHOBLength(unsafe.Sizeof(EFIHOBGenericHeader{})),
				Reserved:  0,
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := constructHOBList(tt.dst, tt.src, &tt.hobLen)

			expectErr(t, err, tt.wantErr)
			if err != nil { // already checked in expectedErr
				return
			}

			// hobLen should be updated as expected after construction
			if tt.hobLen != tt.wantHOBLen {
				t.Fatalf("Unexpected hobLen = %v, want = %v", tt.hobLen, tt.wantHOBLen)
			}

			// deserialize hobList object from bytes
			var efiTable EFIHOBHandoffInfoTable
			if err := binary.Read(tt.dst, binary.LittleEndian, &efiTable); err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if efiTable != tt.efiTable {
				t.Fatalf("Unexpected efiTable: %v, want = %v", efiTable, tt.efiTable)
			}

			// src data
			src := make([]byte, len(tt.src.Bytes()))
			n, rerr := tt.dst.Read(src)
			if rerr != nil {
				t.Fatal(rerr)
			}
			if n != len(tt.src.Bytes()) {
				t.Fatalf("Unexpected src data len: %v, want = %v", n, len(tt.src.Bytes()))
			}
			if !bytes.Equal(src, tt.src.Bytes()) {
				t.Fatalf("Unexpected src data: %v, want = %v", src, tt.src.Bytes())
			}

			// end header
			var endHeader EFIHOBGenericHeader
			if err := binary.Read(tt.dst, binary.LittleEndian, &endHeader); err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if endHeader != tt.endHeader {
				t.Fatalf("Unexpected efiTable: %v, want = %v", endHeader, tt.endHeader)
			}
		})
	}
}

func TestAppendSmbiosTableHOB(t *testing.T) {
	smbiosBase, _, _ := mockGetSMBIOSBase()
	tests := []struct {
		name                   string
		expectedErr            error
		expectedHOBLen         uint64
		expectedEFIHOBGUID     *EFIHOBGUIDType
		expectedUplSmbiosTable *UniversalPayloadSmbiosTable
	}{
		{
			name:           "CASE 1: success",
			expectedErr:    nil,
			expectedHOBLen: uint64(unsafe.Sizeof(EFIHOBGUIDType{}) + unsafe.Sizeof(UniversalPayloadSmbiosTable{})),
			expectedEFIHOBGUID: &EFIHOBGUIDType{
				Header: EFIHOBGenericHeader{
					HOBType:   EFIHOBTypeGUIDExtension,
					HOBLength: EFIHOBLength(unsafe.Sizeof(EFIHOBGUIDType{}) + guidToLength[UniversalPayloadSmbiosTableGUID]),
				},
				Name: guid.MustParse(UniversalPayloadSmbiosTableGUID),
			},
			expectedUplSmbiosTable: &UniversalPayloadSmbiosTable{
				Header: UniversalPayloadGenericHeader{
					Revision: UniversalPayloadSmbiosTableRevision,
					Length:   uint16(unsafe.Sizeof(UniversalPayloadSmbiosTable{})),
				},
				SmBiosEntryPoint: EFIPhysicalAddress(smbiosBase),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func(old func() (int64, int64, error)) { getSMBIOSBase = old }(getSMBIOSBase)
			getSMBIOSBase = mockGetSMBIOSBase

			hobBuf := &bytes.Buffer{}
			var hobLen uint64
			err := appendSmbiosTableHOB(hobBuf, &hobLen)

			expectErr(t, err, tt.expectedErr)
			if err != nil { // already checked in expectedErr
				return
			}

			// hobLen should be updated as expected after construction
			if tt.expectedHOBLen != hobLen {
				t.Fatalf("Unexpected hobLen = %v, want = %v", hobLen, tt.expectedHOBLen)
			}

			// deserialize EFIHOBGUID object from bytes
			var efiHOBGUID EFIHOBGUIDType
			if err := binary.Read(hobBuf, binary.LittleEndian, &efiHOBGUID); err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if *tt.expectedEFIHOBGUID != efiHOBGUID {
				t.Fatalf("Unexpected efiHOBCPU = %v, want = %v", efiHOBGUID, *tt.expectedEFIHOBGUID)
			}

			// deserialize smbiosTable object from bytes
			var smbiosTable UniversalPayloadSmbiosTable
			if err := binary.Read(hobBuf, binary.LittleEndian, &smbiosTable); err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if *tt.expectedUplSmbiosTable != smbiosTable {
				t.Fatalf("Unexpected smbiosTable = %v, want = %v", smbiosTable, *tt.expectedUplSmbiosTable)
			}
		})
	}
}

func expectErr(t *testing.T, err error, expectedErr error) {
	if expectedErr == nil {
		if err != nil {
			t.Fatalf("Unexpected error: %+v", err)
		}
	} else {
		if err == nil {
			t.Fatalf("Expected error %q, got nil", expectedErr)
		}
		if !errors.Is(err, expectedErr) {
			t.Fatalf("Unxpected error %+v, want = %q", err, expectedErr)
		}
	}
}
