// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package universalpayload

import (
	"bytes"
	"encoding/binary"
	"errors"
	guid "github.com/google/uuid"
	"github.com/u-root/u-root/pkg/acpi"
	"github.com/u-root/u-root/pkg/boot/kexec"
	"log"
	"os"
	"strconv"
	"testing"
	"unsafe"
)

func mockKexecMemoryMapFromSysfsMemmap() (kexec.MemoryMap, error) {
	return kexec.MemoryMap{
		{Range: kexec.Range{Start: 0, Size: 50}, Type: kexec.RangeRAM},
		{Range: kexec.Range{Start: 100, Size: 50}, Type: kexec.RangeACPI},
	}, nil
}

func mockGetRSDP() (*acpi.RSDP, error) {
	return &acpi.RSDP{}, nil
}

func mockGetSMBIOSBase() (int64, int64, error) {
	return 100, 200, nil
}

// Main test for constructing HOB list,
// to make sure there are no issues with the overall process
func TestLoadKexecMemWithHOBs(t *testing.T) {
	// mock data
	name := "./testdata/upl.dtb"
	fdtLoad, err := getFdtInfo(name, nil)
	if err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}
	data, _ := os.ReadFile(name)

	defer func(old func() (int64, int64, error)) { getSMBIOSBase = old }(getSMBIOSBase)
	getSMBIOSBase = mockGetSMBIOSBase

	defer func(old func() (*acpi.RSDP, error)) { getAcpiRSDP = old }(getAcpiRSDP)
	getAcpiRSDP = mockGetRSDP

	defer func(old func() (kexec.MemoryMap, error)) { kexecMemoryMapFromSysfsMemmap = old }(kexecMemoryMapFromSysfsMemmap)
	kexecMemoryMapFromSysfsMemmap = mockKexecMemoryMapFromSysfsMemmap

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

	mem, loadErr := loadKexecMemWithHOBs(fdtLoad, data)
	if loadErr != nil {
		t.Fatalf("Unexpected err: %+v", loadErr)
	}
	if mem.Segments == nil || len(mem.Segments) != 2 {
		t.Fatalf("load memory Segments size = %v, want = 2", len(mem.Segments))
	}
	t.Logf("mem Segments: %v", mem.Segments.Phys())

	fdLoadStart := fdtLoad.Load

	// Check segment 0 (hob list for universal payload)
	// Here we only validate segment.Start and ignore the content (bytes),
	// since it's too difficult to split and deserialize hob list from bytes
	hobListAddr := uintptr(fdLoadStart - 0x100000)
	if mem.Segments.Phys()[0].Start != hobListAddr {
		t.Fatalf("universalpayload hob list segment.Start = %v, want = %v", mem.Segments.Phys()[1], hobListAddr)
	}
	// Check segment 1 (fdLoad raw data)
	var fdLoadRange = kexec.Range{Start: uintptr(fdLoadStart), Size: uint(len(data))}
	if mem.Segments.Phys()[1] != fdLoadRange {
		t.Fatalf("universalpayload fdLoad segment range error, range = %v, want = %v",
			mem.Segments.Phys()[1], fdLoadRange)
	}
}

func TestAppendMemMapHOB(t *testing.T) {
	defer func(old func() (kexec.MemoryMap, error)) { kexecMemoryMapFromSysfsMemmap = old }(kexecMemoryMapFromSysfsMemmap)
	kexecMemoryMapFromSysfsMemmap = mockKexecMemoryMapFromSysfsMemmap
	hobBuf := &bytes.Buffer{}
	var hobLen uint64
	err := appendMemMapHOB(hobBuf, &hobLen)
	if err != nil {
		log.Fatal(err)
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

	if len(deserializedHOB) != 1 {
		t.Fatalf("Unexpected hob size = %d, want = %d", len(deserializedHOB), 1)
	}
}

func TestAppendSerialPortHOB(t *testing.T) {
	hobBuf := &bytes.Buffer{}
	var hobLen uint64
	err := appendSerialPortHOB(hobBuf, &hobLen)
	if err != nil {
		log.Fatal(err)
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
	err := appendUniversalPayloadBase(hobBuf, &hobLen, fdtLoad)
	if err != nil {
		log.Fatal(err)
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
			err = binary.Read(tt.dst, binary.LittleEndian, &efiTable)
			if err != nil {
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
			err = binary.Read(tt.dst, binary.LittleEndian, &endHeader)
			if err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if endHeader != tt.endHeader {
				t.Fatalf("Unexpected efiTable: %v, want = %v", endHeader, tt.endHeader)
			}
		})
	}
}

func TestAppendEFICPUHOB(t *testing.T) {
	tests := []struct {
		name           string
		cpuInfoContent string
		expectedBits   int
		expectedErr    error
		expectedHOB    *EFIHOBCPU
	}{
		{
			name: "Valid Physical Address Bits",
			cpuInfoContent: `
processor	: 0
vendor_id	: GenuineIntel
cpu family	: 6
model		: 142
model name	: Intel(R) Core(TM) i7-8550U CPU @ 1.80GHz
stepping	: 10
microcode	: 0xea
cpu MHz		: 1992.000
cache size	: 8192 KB
physical id	: 0
siblings	: 4
core id		: 0
cpu cores	: 2
apicid		: 0
initial apicid	: 0
address sizes	: 39 bits physical, 48 bits virtual
`,
			expectedBits: 39,
			expectedErr:  nil,
			expectedHOB: &EFIHOBCPU{
				Header: EFIHOBGenericHeader{
					HOBType:   EFIHOBTypeCPU,
					HOBLength: EFIHOBLength(unsafe.Sizeof(EFIHOBCPU{})),
				},
				SizeOfMemorySpace: uint8(39),
				SizeOfIOSpace:     DefaultIOAddressSize,
			},
		},
		{
			name: "No Address Size Info",
			cpuInfoContent: `
processor	: 0
vendor_id	: GenuineIntel
cpu family	: 6
model		: 142value out of range
model name	: Intel(R) Core(TM) i7-8550U CPU @ 1.80GHz
`,
			expectedBits: 0,
			expectedErr:  ErrCPUAddressNotFound,
		},
		{
			name: "Invalid Address Size",
			// number value out of range
			cpuInfoContent: `
address sizes	: 1000 bits physical, 48 bits virtual
`,
			expectedBits: 0,
			expectedErr:  strconv.ErrRange,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempFile := mockCPUTempInfoFile(t, tt.cpuInfoContent)
			defer os.Remove(tempFile)

			hobBuf := &bytes.Buffer{}
			var hobLen uint64
			err := appendEFICPUHOB(hobBuf, &hobLen)

			expectErr(t, err, tt.expectedErr)
			if err != nil { // already checked in expectedErr
				return
			}

			// deserialize EFIHOBCPU object from bytes
			var efiHOBCPU EFIHOBCPU
			err = binary.Read(hobBuf, binary.LittleEndian, &efiHOBCPU)
			if err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if *tt.expectedHOB != efiHOBCPU {
				t.Fatalf("Unexpected efiHOBCPU = %v, want = %v", efiHOBCPU, *tt.expectedHOB)
			}
		})
	}
}

func TestAppendAcpiTableHOB(t *testing.T) {
	tests := []struct {
		name               string
		expectedErr        error
		expectedHOBLen     uint64
		expectedEFIHOBGUID *EFIHOBGUIDType
		expectedAcpiTable  *UniversalPayloadAcpiTable
	}{
		{
			name:           "CASE 1: success",
			expectedErr:    nil,
			expectedHOBLen: uint64(unsafe.Sizeof(EFIHOBGUIDType{}) + unsafe.Sizeof(UniversalPayloadSmbiosTable{})),
			expectedEFIHOBGUID: &EFIHOBGUIDType{
				Header: EFIHOBGenericHeader{
					HOBType:   EFIHOBTypeGUIDExtension,
					HOBLength: EFIHOBLength(unsafe.Sizeof(EFIHOBGUIDType{}) + guidToLength[UniversalPayloadAcpiTableGUID]),
				},
				Name: guid.MustParse(UniversalPayloadAcpiTableGUID),
			},
			expectedAcpiTable: &UniversalPayloadAcpiTable{
				Header: UniversalPayloadGenericHeader{
					Revision: UniversalPayloadAcpiTableRevision,
					Length:   uint16(unsafe.Sizeof(UniversalPayloadAcpiTable{})),
				},
				Rsdp: EFIPhysicalAddress(0),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func(old func() (*acpi.RSDP, error)) { getAcpiRSDP = old }(getAcpiRSDP)
			getAcpiRSDP = mockGetRSDP

			hobBuf := &bytes.Buffer{}
			var hobLen uint64
			err := appendAcpiTableHOB(hobBuf, &hobLen)

			expectErr(t, err, tt.expectedErr)
			if err != nil { // already checked in expectedErr
				return
			}

			// hobLen should be updated as expected after construction
			if tt.expectedHOBLen != hobLen {
				t.Fatalf("Unexpected hobLen = %v, want = %v", hobLen, tt.expectedHOBLen)
			}

			// deserialize efiHOBGUID object from bytes
			var efiHOBGUID EFIHOBGUIDType
			err = binary.Read(hobBuf, binary.LittleEndian, &efiHOBGUID)
			if err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if *tt.expectedEFIHOBGUID != efiHOBGUID {
				t.Fatalf("Unexpected efiHOBCPU = %v, want = %v", efiHOBGUID, *tt.expectedEFIHOBGUID)
			}

			// deserialize acpiTable object from bytes
			var acpiTable UniversalPayloadAcpiTable
			err = binary.Read(hobBuf, binary.LittleEndian, &acpiTable)
			if err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if *tt.expectedAcpiTable != acpiTable {
				t.Fatalf("Unexpected acpiTable = %v, want = %v", acpiTable, *tt.expectedAcpiTable)
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
			err = binary.Read(hobBuf, binary.LittleEndian, &efiHOBGUID)
			if err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if *tt.expectedEFIHOBGUID != efiHOBGUID {
				t.Fatalf("Unexpected efiHOBCPU = %v, want = %v", efiHOBGUID, *tt.expectedEFIHOBGUID)
			}

			// deserialize smbiosTable object from bytes
			var smbiosTable UniversalPayloadSmbiosTable
			err = binary.Read(hobBuf, binary.LittleEndian, &smbiosTable)
			if err != nil {
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
