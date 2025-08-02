// Copyright 2012-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (linux && amd64) || (linux && 386)

package memio

import (
	"errors"
	"os"
	"strings"
	"testing"
)

func newPortMock(f *os.File) (PortReadWriter, error) {
	memPort := NewMemIOPort(f)
	return &LinuxPort{
		ReadWriteCloser: memPort,
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
		name        string
		valueUint8  Uint8
		valueUint16 Uint16
		expValue    uint8
		wantErr     string
	}{
		{
			name:       "Uint8",
			valueUint8: Uint8(23),
			expValue:   23,
		},
		{
			name:        "Uint16",
			valueUint16: Uint16(42),
			wantErr:     "/dev/port data must be 8 bits on Linux",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var data Uint8
			if err := port.Out(0x3f8, &tt.valueUint8); err != nil {
				t.Errorf("%q failed at port.Out: %q", tt.name, err)
			}

			if err := port.In(0x3f8, &data); !errors.Is(err, nil) {
				t.Errorf("%q failed at port.In: %q", tt.name, err)
			}

			if data != Uint8(tt.expValue) {
				t.Errorf("%q failed. Got: %d, Want: %d", tt.name, data, tt.expValue)
			}
		})

		t.Run(tt.name+"16Bit_fail", func(t *testing.T) {
			var data Uint16
			if err := port.Out(0x3f8, &tt.valueUint16); err != nil {
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("%q failed at port.Out: %q", tt.name, err)
				}
			}
			if err := port.In(0x3f8, &data); err != nil {
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("%q failed at port.Out: %q", tt.name, err)
				}
			}
		})

		t.Run(tt.name+"_Deprecated", func(t *testing.T) {
			linuxPath = f.Name()
			defer func() { linuxPath = "/dev/port" }()
			data := Uint8(34)
			if err := In(0x23, &data); !errors.Is(err, nil) {
				t.Errorf("%q_Deprecated failed at port.In(): %q", tt.name, err)
			}
			if err := Out(0x23, &data); !errors.Is(err, nil) {
				t.Errorf("%q_Deprecated failed at port.In(): %q", tt.name, err)
			}
		})
		t.Run(tt.name+"_DeprecatedFail", func(t *testing.T) {
			linuxPath = "file-does-not-exist"
			defer func() { linuxPath = "/dev/port" }()
			data := Uint8(34)
			if err := In(0x23, &data); !errors.Is(err, os.ErrNotExist) {
				t.Errorf("%q_Deprecated failed at port.In(): %q", tt.name, err)
			}
		})
		t.Run(tt.name+"_DeprecatedFail", func(t *testing.T) {
			linuxPath = "file-does-not-exist"
			defer func() { linuxPath = "/dev/port" }()
			data := Uint8(34)
			if err := Out(0x23, &data); !errors.Is(err, os.ErrNotExist) {
				t.Errorf("%q_Deprecated failed at port.In(): %q", tt.name, err)
			}
		})
	}
	port.Close()
}

func TestNewPortFail(t *testing.T) {
	linuxPath = "file-does-not-exist"
	defer func() { linuxPath = "/dev/port" }()
	if _, err := NewPort(); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("TestNewPortFail failed: %q", err)
	}
}

func TestNewPortSucceed(t *testing.T) {
	tmpDir := t.TempDir()
	file, err := os.CreateTemp(tmpDir, "tmpfile-")
	if err != nil {
		t.Errorf("TestNewPortSucceed failed: %q", err)
	}
	fdata := make([]byte, 10000)
	if _, err := file.Write(fdata); err != nil {
		t.Errorf("TestPortDev failed: %q", err)
	}

	linuxPath = file.Name()
	defer func() { linuxPath = "/dev/port" }()
	port, err := NewPort()
	if !errors.Is(err, nil) {
		t.Errorf("TestNewPortFail failed: %q", err)
	}
	defer port.Close()
}
