// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sfdp

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/u-root/u-root/pkg/flash/spimock"
)

// fakeSFDPPrettyPrint corresponds to the pretty print of spimock.FakeSFDP.
var fakeSFDPPrettyPrint = `BlockSectorEraseSize           0x1
WriteGranularity               0x1
WriteEnableInstructionRequired 0x0
WriteEnableOpcodeSelect        0x0
4KBEraseOpcode                 0x20
112FastRead                    0x1
AddressBytesNumberUsed         0x1
DoubleTransferRateClocking     0x0
122FastReadSupported           0x1
144FastReadSupported           0x1
114FastReadSupported           0x1
FlashMemoryDensity             0x1fffffff
144FastReadNumberOfWaitStates  0x4
144FastReadNumberOfModeBits    0x2
144FastReadOpcode              0xeb
114FastReadNumberOfWaitStates  0x8
114FastReadNumberOfModeBits    0x0
114FastReadOpcode              0x6b
112FastReadNumberOfWaitStates  0x8
112FastReadNumberOfModeBits    0x0
112FastReadOpcode              0x3b
122FastReadNumberOfWaitStates  0x4
122FastReadNumberOfModeBits    0x0
122FastReadOpcode              0xbb
222FastReadSupported           0x0
444FastReadSupported           0x1
`

// errorLookupParams contains Params which will return an error when read.
var errorLookupParams = []ParamLookupEntry{
	{"ExistingParam", Param{0, 0, 0x00, 0x20}},
	{"NonExistingTable", Param{0x999, 0, 0x00, 0x20}},
	{"NonExistingDword", Param{0, 0x999, 0x00, 0x20}},
}

// errorPrettyPrint corresponds to the pretty print of spimock.FakeSFDP for the
// Params in errorLookupParams.
var errorPrettyPrint = `ExistingParam    0xfff320e5
NonExistingTable Error: could not find table 0x999
NonExistingDword Error: could not find dword 0x999 in table 0x0
`

// TestPrettyPrint prints the SFDP to a beatiful string.
func TestPrettyPrint(t *testing.T) {
	for _, tt := range []struct {
		name        string
		fakeSFDP    []byte
		paramLookup []ParamLookupEntry
		prettyPrint string
	}{
		{
			name:        "MX66L51235F",
			fakeSFDP:    spimock.FakeSFDP,
			paramLookup: BasicTableLookup,
			prettyPrint: fakeSFDPPrettyPrint,
		},
		{
			name:        "error prints",
			fakeSFDP:    spimock.FakeSFDP,
			paramLookup: errorLookupParams,
			prettyPrint: errorPrettyPrint,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewReader(tt.fakeSFDP)
			sfdp, err := Read(r)
			if err != nil {
				t.Fatal(err)
			}

			w := &bytes.Buffer{}
			if err := sfdp.PrettyPrint(w, tt.paramLookup); err != nil {
				t.Fatal(err)
			}

			if w.String() != tt.prettyPrint {
				t.Errorf("sfdp.PrettyPrint() =\n%s\n;want\n%s", w.String(), tt.prettyPrint)
			}
		})
	}
}

func TestRead(t *testing.T) {
	for _, tt := range []struct {
		name      string
		fakeSFDP  []byte
		wantError error
	}{
		{
			name:     "successful",
			fakeSFDP: spimock.FakeSFDP,
		},
		{
			name:      "magic not found",
			fakeSFDP:  spimock.FakeSFDP[4:],
			wantError: &UnsupportedError{},
		},
		{
			// Induce an EOF reader error with an empty buffer.
			name:      "reader error",
			fakeSFDP:  []byte{},
			wantError: io.EOF,
		},
		{
			// Induce a reader error on the tables with a truncated
			// buffer.
			name:      "reader table error",
			fakeSFDP:  spimock.FakeSFDP[:0x30],
			wantError: io.EOF,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewReader(tt.fakeSFDP)
			_, err := Read(r)
			if !errors.Is(err, tt.wantError) {
				t.Errorf("sfdp.Read() err = %v; want %v", err, tt.wantError)
			}
		})
	}
}

// TestRead16BitId checks that the v1.5 16-bit table IDs are supported.
func TestRead16BitId(t *testing.T) {
	sfdpV1_5 := []byte{
		// 0x00: Magic
		0x53, 0x46, 0x44, 0x50,
		// 0x04: Version v1.5
		0x05, 0x01,
		// 0x06: Number of headers, 0 means there's 1 header
		0x00,
		// 0x07: Unused
		0x00,
		// 0x08: Parameter header ID LSB
		0xcd,
		// 0x09: Version v1.0
		0x00, 0x01,
		// 0x0b: Number of DWORDS in the table
		0x01,
		// 0x0c: Pointer
		0x10, 0x00, 0x00,
		// 0x0f: ID MSB
		0xab,
		// 0x10: DWORD 0
		0x78, 0x56, 0x34, 0x12,
	}

	r := bytes.NewReader(sfdpV1_5)
	sfdp, err := Read(r)
	if err != nil {
		t.Fatal(err)
	}

	testParam := Param{
		Table: 0xabcd,
		Dword: 0,
		Shift: 0,
		Bits:  32,
	}

	val, err := sfdp.Param(testParam)
	if err != nil {
		t.Error(err)
	}
	var want int64 = 0x12345678
	if val != want {
		t.Errorf("sfdp.Param() = %#x; want %#x", val, want)
	}
}
