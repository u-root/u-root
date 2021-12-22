// Copyright 2012-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (linux && amd64) || (linux && 386)
// +build linux,amd64 linux,386

package memio

import (
	"os"
	"testing"
)

func newPortMock(f *os.File) (Port, error) {
	memPort, err := NewMemIOPort(f)
	if err != nil {
		return nil, err
	}
	return &LinuxPort{
		MemIO: memPort,
	}, nil
}

func TestLinuxPort(t *testing.T) {
	tmpDir := t.TempDir()
	f, err := os.CreateTemp(tmpDir, "testMem-file-")
	if err != nil {
		t.Errorf("TestPortDev failed: %q", err)
	}
	fdata := make([]byte, 10000)
	if _, err := f.Write(fdata); err != nil {
		t.Errorf("TestPortDev failed: %q", err)
	}
	port, err := newPortMock(f)
	if err != nil {
		t.Fatalf("TestPortDev failed: %q", err)
	}
	for _, tt := range []struct {
		name       string
		valueUint8 Uint8
		expValue   uint8
	}{
		{
			name:       "Uint8",
			valueUint8: Uint8(23),
			expValue:   23,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var data Uint8
			if err := port.Out(0x3f8, &tt.valueUint8); err != nil {
				t.Errorf("klsedfjaghkl;jasd %q", err)
			}

			if err := port.In(0x3f8, &data); err != nil {
				t.Fatal(err)
			}

			if data != Uint8(tt.expValue) {
				t.Errorf("%q failed. Got: %d, Want: %d", tt.name, data, tt.expValue)
			}
		})
	}
	port.Close()
}
