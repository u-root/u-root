// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package universalpayload

import (
	"bytes"
	"encoding/binary"
	guid "github.com/google/uuid"
	"github.com/u-root/u-root/pkg/acpi"
	"github.com/u-root/u-root/pkg/boot/kexec"
	"log"
	"os"
	"strings"
	"testing"
	"unsafe"
)

// Main test for constructing Hob list,
// to make sure there are no issues with the overall process
func TestLoadKexecMemWithHobs(t *testing.T) {
	// mock data
	name := "./testdata/upl.dtb"
	fdtLoad, err := getFdtInfo(name, nil)
	if err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}
	data, _ := os.ReadFile(name)

	defer func(old func() (int64, int64, error)) { smbiosSMBIOSBase = old }(smbiosSMBIOSBase)
	smbiosSMBIOSBase = mockGetSMBIOSBase

	defer func(old func() (*acpi.RSDP, error)) { acpiGetRSDP = old }(acpiGetRSDP)
	acpiGetRSDP = mockGetRSDP

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

	mem, loadErr := loadKexecMemWithHobs(fdtLoad, data)
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

func mockKexecMemoryMapFromSysfsMemmap() (kexec.MemoryMap, error) {
	return kexec.MemoryMap{
		{Range: kexec.Range{Start: 0, Size: 50}, Type: kexec.RangeRAM},
		{Range: kexec.Range{Start: 100, Size: 50}, Type: kexec.RangeACPI},
	}, nil
}

func TestAppendMemMapHob(t *testing.T) {
	defer func(old func() (kexec.MemoryMap, error)) { kexecMemoryMapFromSysfsMemmap = old }(kexecMemoryMapFromSysfsMemmap)
	kexecMemoryMapFromSysfsMemmap = mockKexecMemoryMapFromSysfsMemmap
	hobBuf := &bytes.Buffer{}
	var hobLen uint64
	err := appendMemMapHob(hobBuf, &hobLen)
	if err != nil {
		log.Fatal(err)
	}

	var deserializedHob EfiMemoryMapHob
	for hobBuf.Len() > 0 {
		var hob EfiHobResourceDescriptor
		err := binary.Read(hobBuf, binary.LittleEndian, &hob)
		if err != nil {
			t.Fatalf("Unexpected error: %+v", err)
		}
		deserializedHob = append(deserializedHob, hob)
	}

	if len(deserializedHob) != 1 {
		t.Errorf("Unexpected hob size = %d, want = %d", len(deserializedHob), 1)
	}
}

func TestAppendSerialPortHob(t *testing.T) {
	hobBuf := &bytes.Buffer{}
	var hobLen uint64
	err := appendSerialPortHob(hobBuf, &hobLen)
	if err != nil {
		log.Fatal(err)
	}

	var serialPortInfo UniversalPayloadSerialPortInfo
	var efiHobGUID EfiHobGUIDType
	if err := binary.Read(hobBuf, binary.LittleEndian, &efiHobGUID); err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}
	if err := binary.Read(hobBuf, binary.LittleEndian, &serialPortInfo); err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}
	if efiHobGUID.Name.String() != UniversalPayloadSerialPortInfoGUID {
		t.Errorf("Unexpected efiHobGuid = %v, want = %v",
			efiHobGUID.Name.String(), UniversalPayloadSerialPortInfoGUID)
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
	var uplBaseGUIDHob EfiHobGUIDType
	if err := binary.Read(hobBuf, binary.LittleEndian, &uplBaseGUIDHob); err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}
	if err := binary.Read(hobBuf, binary.LittleEndian, &uplBase); err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}
	if hobBuf.Len() > 0 {
		t.Fatalf("Unexpected hobBuf size after deserialization: %v, want = 0", hobBuf.Len())
	}
	if uplBaseGUIDHob.Name.String() != UniversalPayloadBaseGUID {
		t.Errorf("Unexpected uplBaseGUIDHob = %v, want = %v",
			uplBaseGUIDHob.Name.String(), UniversalPayloadBaseGUID)
	}
	if uplBase.Entry != EfiPhysicalAddress(fdtLoad) {
		t.Fatalf("Unexpected UniversalPayloadBase.Entry: %v, want = %v", uplBase.Entry, fdtLoad)
	}
}

func TestConstructHobList(t *testing.T) {
	const lenOfEfiTable = uint64(unsafe.Sizeof(EfiHobHandoffInfoTable{}))

	for _, tt := range []struct {
		// mock data
		name   string
		dst    *bytes.Buffer
		src    *bytes.Buffer
		hobLen uint64

		// expected data
		wantHobLen uint64
		efiTable   EfiHobHandoffInfoTable
		endHeader  EfiHobGenericHeader
		wantErr    string
	}{
		{
			name:   "CASE 1: success",
			dst:    &bytes.Buffer{},
			src:    bytes.NewBuffer([]byte("12345678")),
			hobLen: uint64(8),

			wantHobLen: uint64(8 + unsafe.Sizeof(EfiHobGenericHeader{})),
			efiTable: EfiHobHandoffInfoTable{
				Header: EfiHobGenericHeader{
					HobType:   EfiHobTypeHandoff,
					HobLength: uint16(unsafe.Sizeof(EfiHobHandoffInfoTable{})),
				},
				Version:             EfiHobHandoffInfoVersion,
				BootMode:            EfiHobHandoffInfoBootMode,
				EfiMemoryTop:        EfiHobHandoffInfoEfiMemoryTop,
				EfiMemoryBottom:     EfiHobHandoffInfoEfiMemoryBottom,
				EfiFreeMemoryTop:    EfiHobHandoffInfoFreeEfiMemoryTop,
				EfiFreeMemoryBottom: EfiPhysicalAddress(EfiHobHandoffInfoFreeMemoryBottom + 8 + 8 + lenOfEfiTable),
				EfiEndOfHobList:     EfiPhysicalAddress(EfiHobHandoffInfoFreeMemoryBottom + 8 + lenOfEfiTable),
			},
			endHeader: EfiHobGenericHeader{
				HobType:   EfiHobTypeEndOfHobList,
				HobLength: uint16(unsafe.Sizeof(EfiHobGenericHeader{})),
				Reserved:  0,
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := constructHobList(tt.dst, tt.src, &tt.hobLen)
			if err != nil {
				if len(tt.wantErr) == 0 || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatal(err)
				}
			} else if len(tt.wantErr) > 0 {
				t.Errorf("Expected err = %v", tt.wantErr)
			}

			// hobLen should be updated as expected after construction
			if tt.hobLen != tt.wantHobLen {
				t.Errorf("Unexpected hobLen = %v, want = %v", tt.hobLen, tt.wantHobLen)
			}

			// deserialize hobList object from bytes
			var efiTable EfiHobHandoffInfoTable
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
			var endHeader EfiHobGenericHeader
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

func TestAppendEfiCPUHob(t *testing.T) {
	tests := []struct {
		name           string
		cpuInfoContent string
		expectedBits   int
		expectedErr    string
		expectedHob    *EfiHobCPU
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
			expectedErr:  "",
			expectedHob: &EfiHobCPU{
				Header: EfiHobGenericHeader{
					HobType:   EfiHobTypeCPU,
					HobLength: uint16(unsafe.Sizeof(EfiHobCPU{})),
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
			expectedErr:  "address sizes information not found",
		},
		{
			name: "Invalid Address Size",
			// number value out of range
			cpuInfoContent: `
address sizes	: 1000 bits physical, 48 bits virtual
`,
			expectedBits: 0,
			expectedErr:  "phyAddrSize 1000 out of range for uint8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempFile := mockCPUTempInfoFile(t, tt.cpuInfoContent)
			defer os.Remove(tempFile)

			hobBuf := &bytes.Buffer{}
			var hobLen uint64
			err := appendEfiCPUHob(hobBuf, &hobLen)

			if tt.expectedErr == "" {
				// success validation
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}

				// deserialize EfiHobCPU object from bytes
				var efiHobCPU EfiHobCPU
				err = binary.Read(hobBuf, binary.LittleEndian, &efiHobCPU)
				if err != nil {
					t.Fatalf("Unexpected error: %+v", err)
				}
				if *tt.expectedHob != efiHobCPU {
					t.Errorf("Unexpected efiHobCPU = %v, want = %v", efiHobCPU, *tt.expectedHob)
				}
			} else {
				// fault validation
				if err == nil {
					t.Fatalf("Expected error %q, got nil", tt.expectedErr)
				}
				if err.Error()[:len(tt.expectedErr)] != tt.expectedErr {
					t.Errorf("Unxpected error %q, want = %q", err.Error(), tt.expectedErr)
				}
			}
		})
	}
}

func mockGetRSDP() (*acpi.RSDP, error) {
	return &acpi.RSDP{}, nil
}

func mockGetSMBIOSBase() (int64, int64, error) {
	return 100, 200, nil
}

func TestAppendAcpiTableHob(t *testing.T) {
	tests := []struct {
		name               string
		expectedErr        string
		expectedHobLen     uint64
		expectedEfiHobGUID *EfiHobGUIDType
		expectedAcpiTable  *UniversalPayloadAcpiTable
	}{
		{
			name:           "CASE 1: success",
			expectedErr:    "",
			expectedHobLen: uint64(unsafe.Sizeof(EfiHobGUIDType{}) + unsafe.Sizeof(UniversalPayloadSmbiosTable{})),
			expectedEfiHobGUID: &EfiHobGUIDType{
				Header: EfiHobGenericHeader{
					HobType:   EfiHobTypeGUIDExtension,
					HobLength: uint16(unsafe.Sizeof(EfiHobGUIDType{}) + guidToLength[UniversalPayloadAcpiTableGUID]),
				},
				Name: guid.MustParse(UniversalPayloadAcpiTableGUID),
			},
			expectedAcpiTable: &UniversalPayloadAcpiTable{
				Header: UniversalPayloadGenericHeader{
					Revision: UniversalPayloadAcpiTableRevision,
					Length:   uint16(unsafe.Sizeof(UniversalPayloadAcpiTable{})),
				},
				Rsdp: EfiPhysicalAddress(0),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func(old func() (*acpi.RSDP, error)) { acpiGetRSDP = old }(acpiGetRSDP)
			acpiGetRSDP = mockGetRSDP

			hobBuf := &bytes.Buffer{}
			var hobLen uint64
			err := appendAcpiTableHob(hobBuf, &hobLen)

			if tt.expectedErr == "" {
				// success validation
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}

				// hobLen should be updated as expected after construction
				if tt.expectedHobLen != hobLen {
					t.Errorf("Unexpected hobLen = %v, want = %v", hobLen, tt.expectedHobLen)
				}

				// deserialize efiHobGUID object from bytes
				var efiHobGUID EfiHobGUIDType
				err = binary.Read(hobBuf, binary.LittleEndian, &efiHobGUID)
				if err != nil {
					t.Fatalf("Unexpected error: %+v", err)
				}
				if *tt.expectedEfiHobGUID != efiHobGUID {
					t.Errorf("Unexpected efiHobCPU = %v, want = %v", efiHobGUID, *tt.expectedEfiHobGUID)
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
			} else {
				// fault validation
				if err == nil {
					t.Fatalf("Expected error %q, got nil", tt.expectedErr)
				}
				if err.Error()[:len(tt.expectedErr)] != tt.expectedErr {
					t.Errorf("Unxpected error %q, want = %q", err.Error(), tt.expectedErr)
				}
			}
		})
	}
}

func TestAppendSmbiosTableHob(t *testing.T) {
	smbiosBase, _, _ := mockGetSMBIOSBase()
	tests := []struct {
		name                   string
		expectedErr            string
		expectedHobLen         uint64
		expectedEfiHobGUID     *EfiHobGUIDType
		expectedUplSmbiosTable *UniversalPayloadSmbiosTable
	}{
		{
			name:           "CASE 1: success",
			expectedErr:    "",
			expectedHobLen: uint64(unsafe.Sizeof(EfiHobGUIDType{}) + unsafe.Sizeof(UniversalPayloadSmbiosTable{})),
			expectedEfiHobGUID: &EfiHobGUIDType{
				Header: EfiHobGenericHeader{
					HobType:   EfiHobTypeGUIDExtension,
					HobLength: uint16(unsafe.Sizeof(EfiHobGUIDType{}) + guidToLength[UniversalPayloadSmbiosTableGUID]),
				},
				Name: guid.MustParse(UniversalPayloadSmbiosTableGUID),
			},
			expectedUplSmbiosTable: &UniversalPayloadSmbiosTable{
				Header: UniversalPayloadGenericHeader{
					Revision: UniversalPayloadSmbiosTableRevision,
					Length:   uint16(unsafe.Sizeof(UniversalPayloadSmbiosTable{})),
				},
				SmBiosEntryPoint: EfiPhysicalAddress(smbiosBase),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func(old func() (int64, int64, error)) { smbiosSMBIOSBase = old }(smbiosSMBIOSBase)
			smbiosSMBIOSBase = mockGetSMBIOSBase

			hobBuf := &bytes.Buffer{}
			var hobLen uint64
			err := appendSmbiosTableHob(hobBuf, &hobLen)

			if tt.expectedErr == "" {
				// success validation
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}

				// hobLen should be updated as expected after construction
				if tt.expectedHobLen != hobLen {
					t.Errorf("Unexpected hobLen = %v, want = %v", hobLen, tt.expectedHobLen)
				}

				// deserialize efiHobGUID object from bytes
				var efiHobGUID EfiHobGUIDType
				err = binary.Read(hobBuf, binary.LittleEndian, &efiHobGUID)
				if err != nil {
					t.Fatalf("Unexpected error: %+v", err)
				}
				if *tt.expectedEfiHobGUID != efiHobGUID {
					t.Errorf("Unexpected efiHobCPU = %v, want = %v", efiHobGUID, *tt.expectedEfiHobGUID)
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
			} else {
				// fault validation
				if err == nil {
					t.Fatalf("Expected error %q, got nil", tt.expectedErr)
				}
				if err.Error()[:len(tt.expectedErr)] != tt.expectedErr {
					t.Errorf("Unxpected error %q, want = %q", err.Error(), tt.expectedErr)
				}
			}
		})
	}
}
