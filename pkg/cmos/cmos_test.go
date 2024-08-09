// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (linux && amd64) || (linux && 386)

package cmos

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/hugelgupf/vmtest/guest"
	"github.com/u-root/u-root/pkg/memio"
)

func newMock(errStr string, inBuf, outBuf io.ReadWriter, f *os.File) *Chip {
	memPort := memio.NewMemIOPort(f)
	return &Chip{
		PortReadWriter: &memio.LinuxPort{
			ReadWriteCloser: memPort,
		},
	}
}

func TestCMOS(t *testing.T) {
	for _, tt := range []struct {
		name                string
		addr                memio.Uint8
		writeData, readData memio.UintN
		err                 string
	}{
		{
			name:      "uint8",
			addr:      0x10,
			writeData: &[]memio.Uint8{0x12}[0],
		},
		{
			name:      "uint16",
			addr:      0x20,
			writeData: &[]memio.Uint16{0x1234}[0],
		},
		{
			name:      "uint32",
			addr:      0x30,
			writeData: &[]memio.Uint32{0x12345678}[0],
		},
		{
			name:      "uint64",
			addr:      0x40,
			writeData: &[]memio.Uint64{0x1234567890abcdef}[0],
			err:       "/dev/port data must be 8 bits on Linux",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			f, err := os.CreateTemp(tmpDir, "testfile-")
			if err != nil {
				t.Errorf(`os.CreateTemp(tmpDir, "testfile-") = f, %q, not f, nil`, err)
			}
			data := make([]byte, 10000)
			if _, err := f.Write(data); err != nil {
				t.Errorf(`f.Write(data) = _, %q, not _, nil`, err)
			}
			var in, out bytes.Buffer
			c := newMock(tt.err, &in, &out, f)
			// Set internal function to dummy but save old state for reset later
			if err := c.Write(tt.addr, tt.writeData); err != nil {
				if !strings.Contains(err.Error(), tt.err) {
					t.Errorf(`c.Write(tt.addr, tt.writeData) = %q, not nil`, err)
				}
			}
			err = c.Read(tt.addr, tt.readData)
			if err != nil {
				if !strings.Contains(err.Error(), tt.err) {
					t.Errorf(`c.Read(tt.addr, tt.readData) = %q, not nil`, err)
				}
			}
			// We can only progress if error is nil.
			if err == nil {
				got := in.String()
				want := tt.writeData.String()
				if got != want {
					t.Errorf("%s, not %s", want, got)
				}
			}
		})
	}
}

// This is just for coverage percentage. This test does nothing of any other value.
func TestNew(t *testing.T) {
	guest.SkipIfNotInVM(t)

	if _, err := New(); err != nil {
		t.Errorf(`New() = %q, not nil`, err)
	}
}
