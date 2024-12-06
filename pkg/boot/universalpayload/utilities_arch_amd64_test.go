// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build amd64

package universalpayload

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
	"strconv"
	"testing"
	"unsafe"
)

func TestGetPhysicalAddressSizes(t *testing.T) {
	tests := []struct {
		name           string
		cpuInfoContent string
		expectedBits   uint8
		expectedErr    error
	}{
		{
			name: "Valid Address Size",
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

			physicalBits, err := getPhysicalAddressSizes()

			if tt.expectedErr == nil {
				// success validation
				if err != nil {
					t.Fatalf("Unexpected error: %+v", err)
				}
				if physicalBits != tt.expectedBits {
					t.Errorf("Unexpected physical address size %d, want = %d", physicalBits, tt.expectedBits)
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
			if err := binary.Read(hobBuf, binary.LittleEndian, &efiHOBCPU); err != nil {
				t.Fatalf("Unexpected error: %+v", err)
			}
			if *tt.expectedHOB != efiHOBCPU {
				t.Fatalf("Unexpected efiHOBCPU = %v, want = %v", efiHOBCPU, *tt.expectedHOB)
			}
		})
	}
}
