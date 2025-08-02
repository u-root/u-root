// Copyright 2012-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package memio

import (
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestIORealSyscalls(t *testing.T) {
	for _, tt := range []struct {
		name                string
		addr                int64
		writeData, readData UintN
		err                 string
	}{
		{
			name:      "uint8",
			addr:      0x10,
			writeData: &[]Uint8{0x12}[0],
			readData:  new(Uint8),
		},
		{
			name:      "uint16",
			addr:      0x20,
			writeData: &[]Uint16{0x1234}[0],
			readData:  new(Uint16),
		},
		{
			name:      "uint32",
			addr:      0x30,
			writeData: &[]Uint32{0x12345678}[0],
			readData:  new(Uint32),
		},
		{
			name:      "uint64",
			addr:      0x40,
			writeData: &[]Uint64{0x1234567890abcdef}[0],
			readData:  new(Uint64),
		},
		{
			name:      "byte slice",
			addr:      0x50,
			writeData: &[]ByteSlice{[]byte("Hello")}[0],
			readData:  &[]ByteSlice{make([]byte, 5)}[0],
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "io_test")
			if err != nil {
				t.Fatal(err)
			}
			tmpFile.Write(make([]byte, 10000))
			tmpFile.Close()
			defer os.Remove(tmpFile.Name())
			m, err := NewMMap(tmpFile.Name())
			if err != nil {
				t.Errorf("%q failed at NewMMap: %q", tt.name, err)
			}
			defer m.Close()
			// Write to the file.
			if err := m.WriteAt(tt.addr, tt.writeData); err != nil {
				if err.Error() == tt.err {
					return
				}
				t.Fatal(err)
			}

			// Read back the value.
			if err := m.ReadAt(tt.addr, tt.readData); err != nil {
				if err.Error() == tt.err {
					return
				}
				t.Fatal(err)
			}

			want := tt.writeData
			got := tt.readData
			if !reflect.DeepEqual(want, got) {
				t.Fatalf("Write(%#016x, %v) = %v; want %v",
					tt.addr, want, got, want)
			}
		})
		t.Run(tt.name+"Deprecated", func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "io_test")
			if err != nil {
				t.Fatal(err)
			}
			tmpFile.Write(make([]byte, 10000))
			tmpFile.Close()
			defer os.Remove(tmpFile.Name())
			memPath = tmpFile.Name()
			defer func() { memPath = "/dev/mem" }()

			// Write to the file.
			if err := Write(tt.addr, tt.writeData); err != nil {
				if err.Error() == tt.err {
					return
				}
				t.Fatal(err)
			}

			// Read back the value.
			if err := Read(tt.addr, tt.readData); err != nil {
				if err.Error() == tt.err {
					return
				}
				t.Fatal(err)
			}

			want := tt.writeData
			got := tt.readData
			if !reflect.DeepEqual(want, got) {
				t.Fatalf("Write(%#016x, %v) = %v; want %v",
					tt.addr, want, got, want)
			}
		})
	}
}

func TestNetMMapFail(t *testing.T) {
	_, err := NewMMap("file-does-not-exist")
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("TestNetMapFail failed: %q", err)
	}
}

func TestReadWriteErrorWrongPath(t *testing.T) {
	memPath = "file-does-not-exist"
	defer func() { memPath = "/dev/mem" }()

	var data UintN
	if err := Write(0x35, data); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("TestReadWriteErrorWrongPath failed at Write(..): %q", err)
	}
	if err := Read(0x35, data); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("TestReadWriteErrorWrongPath failed at Read(..): %q", err)
	}
}

type fakeSyscalls struct {
	errMmap   error
	errMunMap error
	retBytes  []byte
}

func (f *fakeSyscalls) Mmap(fd int, page int64, mapSize int, prot int, callid int) ([]byte, error) {
	return f.retBytes, f.errMmap
}

func (f *fakeSyscalls) Munmap(mem []byte) error {
	return f.errMunMap
}

func TestMemIOAbstractSyscalls(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "io_test")
	if err != nil {
		t.Fatal(err)
	}
	tmpFile.Write(make([]byte, 10000))
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())
	m, err := NewMMap(tmpFile.Name())
	if err != nil {
		t.Errorf("TestMemIOAbstractSyscalls failed at NewMMap: %q", err)
	}
	defer m.Close()
	for _, tt := range []struct {
		name      string
		errMmap   string
		errMunMap string
		retbyte   []byte
		data      UintN
	}{
		{
			name:    "TestMmapError",
			errMmap: "force mmap error",
			data:    &[]Uint8{0x12}[0],
		},
	} {
		{
			m.syscalls = &fakeSyscalls{
				errMmap:   errors.New(tt.errMmap),
				errMunMap: errors.New(tt.errMunMap),
				retBytes:  tt.retbyte,
			}
			t.Run(tt.name, func(t *testing.T) {
				if err := m.ReadAt(0x23, tt.data); !strings.Contains(err.Error(), tt.errMmap) {
					t.Errorf("%q_ReadAt failed. Want: %q, Got: %q", tt.name, tt.errMmap, err)
				}
			})
			t.Run(tt.name, func(t *testing.T) {
				if err := m.WriteAt(0x23, tt.data); !strings.Contains(err.Error(), tt.errMmap) {
					t.Errorf("%q_WriteAt failed. Want: %q, Got: %q", tt.name, tt.errMmap, err)
				}
			})
		}
	}
}
