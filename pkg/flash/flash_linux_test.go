// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flash

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/u-root/u-root/pkg/flash/spimock"
)

// TestSFDPReader tests reading arbitrary offsets from the SFDP.
func TestSFDPReader(t *testing.T) {
	fakeErr := errors.New("fake transfer error")
	for _, tt := range []struct {
		name             string
		readOffset       int64
		readSize         int
		forceTransferErr error
		wantData         []byte
		wantNewErr       error
		wantReadAtErr    error
	}{
		{
			name:       "read sfdp data",
			readOffset: 0x10,
			readSize:   4,
			wantData:   []byte{0xc2, 0x00, 0x01, 0x04},
		},
		{
			name:          "invalid offset",
			readOffset:    sfdpMaxAddress + 1,
			readSize:      4,
			wantReadAtErr: io.EOF,
		},
		{
			name:             "transfer error",
			readOffset:       0x10,
			readSize:         4,
			forceTransferErr: fakeErr,
			wantNewErr:       fakeErr,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			s := spimock.New()
			s.ForceTransferErr = tt.forceTransferErr
			f, err := New(s)
			if !errors.Is(err, tt.wantNewErr) {
				t.Errorf("flash.New() err = %v; want %v", err, tt.forceTransferErr)
			}
			if err != nil {
				return
			}

			data := make([]byte, tt.readSize)
			n, err := f.SFDPReader().ReadAt(data, tt.readOffset)
			if gotErrString, wantErrString := fmt.Sprint(err), fmt.Sprint(tt.wantReadAtErr); gotErrString != wantErrString {
				t.Errorf("SFDPReader().ReadAt() err = %q; want %q", gotErrString, wantErrString)
			}
			if err == nil && n != len(data) {
				t.Errorf("SFDPReader().ReadAt() n = %d; want %d", n, len(data))
			}

			if err == nil && string(data) != string(tt.wantData) {
				t.Errorf("SFDPReader().ReadAt() data = %#02x; want %#02x", data, tt.wantData)
			}
		})
	}
}

// TestSFDPReadDWORD checks a DWORD can be parsed from the SFDP tables.
func TestSFDPReadDWORD(t *testing.T) {
	s := spimock.New()
	f, err := New(s)
	if err != nil {
		t.Fatal(err)
	}

	sfdp := f.SFDP()
	if sfdp == nil {
		t.Fatalf("f.SFDP: got nil, want value")
	}

	dword, err := sfdp.Dword(0, 0)
	if err != nil {
		t.Error(err)
	}
	var want uint32 = 0xfff320e5
	if dword != want {
		t.Errorf("sfdp.TableDword() = %#08x; want %#08x", dword, want)
	}
}
