// Copyright 2012-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (linux && amd64) || (linux && 386)
// +build linux,amd64 linux,386

package memio

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

func ExampleIn() {
	var data Uint8
	if err := In(0x3f8, &data); err != nil {
		log.Fatal(err)
	}
	fmt.Println(data)
}

func ExampleOut() {
	data := Uint8('A')
	if err := Out(0x3f8, &data); err != nil {
		log.Fatal(err)
	}
}

func ExampleArchIn() {
	var data Uint8
	if err := In(0x80, &data); err != nil {
		log.Fatal(err)
	}
	fmt.Println(data)
}

func ExampleArchOut() {
	data := Uint8('A')
	if err := Out(0x80, &data); err != nil {
		log.Fatal(err)
	}
}

var testsUint16 = []struct {
	name                string
	addr                uint16
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
}

func TestIn(t *testing.T) {
	// In calls pathRead(portPath,...) therefore a tmpfile shall be used
	tmpFile, err := ioutil.TempFile("", "TestIn")
	if err != nil {
		t.Errorf("Failed to create tempfile: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	portPath = tmpFile.Name()
	defer func() { portPath = "/dev/port" }()
	_, err = tmpFile.Write(make([]byte, 10000))
	if err != nil {
		t.Errorf("Failed to write to tempfile: %v", err)
	}
	tmpFile.Close()
	for _, tt := range testsUint16 {
		t.Run(fmt.Sprintf("In(%v)", tt.name), func(t *testing.T) {
			if err := In(tt.addr, tt.readData); err != nil {
				switch err.Error() {
				case "/dev/port data must be 8 bits on Linux":
					return
					// Catches pathRead failing to access "/dev/port" or tmpfile, but tmpfile is set in this test
				case "/dev/port: permission denied":
					return
					// Catches empty tmpfile, but tmpfile is not empty here
				case "EOF":
					return
				default:
					t.Errorf("Want %v, got %v", nil, err)
				}
			}
		})
	}
}
func TestOut(t *testing.T) {
	// Out calls pathWrite(portPath,...) therefore a tmpfile shall be used
	tmpFile, err := ioutil.TempFile("", "TestOut")
	if err != nil {
		t.Errorf("Failed to create tempfile: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	portPath = tmpFile.Name()
	defer func() { portPath = "/dev/port" }()
	_, err = tmpFile.Write(make([]byte, 10000))
	if err != nil {
		t.Errorf("Failed to write to tempfile: %v", err)
	}
	tmpFile.Close()
	for _, tt := range testsUint16 {
		t.Run(fmt.Sprintf("In(%v)", tt.name), func(t *testing.T) {
			if err := Out(tt.addr, tt.writeData); err != nil {
				switch err.Error() {
				case "/dev/port data must be 8 bits on Linux":
					return
					// Catches pathRead failing to access "/dev/port" or tmpfile, but tmpfile is set in this test
				case "/dev/port: permission denied":
					return
					// Catches empty tmpfile, but tmpfile is not empty here
				case "EOF":
					return
				default:
					t.Errorf("Want %v, got %v", nil, err)
				}
			}
		})
	}
}

func TestUintErr(t *testing.T) {
	oSyscallIopl := syscallIopl
	syscallIopl = func(int) error { return nil }
	defer func() { syscallIopl = oSyscallIopl }()

	// Mock assembly implementations
	oAIL := archInLong
	archInLong = func(uint16) uint32 { return uint32(0x11) }
	defer func() { archInLong = oAIL }()
	oAIW := archInWord
	archInWord = func(uint16) uint16 { return uint16(0x11) }
	defer func() { archInWord = oAIW }()
	oAIB := archInb
	archInByte = func(uint16) uint8 { return uint8(0x11) }
	defer func() { archInByte = oAIB }()
	oAOL := archOutLong
	archOutLong = func(uint16, uint32) {}
	defer func() { archOutLong = oAOL }()
	oAOW := archOutWord
	archOutWord = func(uint16, uint16) {}
	defer func() { archOutWord = oAOW }()
	oAOB := archOutByte
	archOutByte = func(uint16, uint8) {}
	defer func() { archOutByte = oAOB }()

	for _, tt := range testsUint16 {
		t.Run(fmt.Sprintf("ArchIn(%v)", tt.name), func(t *testing.T) {
			if err := ArchIn(tt.addr, tt.readData); err != nil {
				switch err.Error() {
				case "port data must be 8, 16 or 32 bits":
					return
				// catching the failed sysopcall due to insufficient permissions here
				case "operation not permitted":
					return

				default:
					t.Errorf("ArchIn failed: %v", err)
				}
			}
			got := tt.readData.Size()
			want := tt.writeData.Size()
			if !reflect.DeepEqual(got, want) {
				t.Errorf("ArchIn(%#016x) got size %v, want %v", tt.addr, got, want)
			}
		})

		t.Run(fmt.Sprintf("ArchOut(%v)", tt.name), func(t *testing.T) {
			t.Logf("%T, %T, size: %v", tt.addr, tt.writeData, tt.writeData.Size())
			if err := ArchOut(tt.addr, tt.writeData); err != nil {
				switch err.Error() {
				case "port data must be 8, 16 or 32 bits":
					return
				// catching the failed sysopcall due to insufficient permissions here
				case "operation not permitted":
					return

				default:
					t.Errorf("ArchOut failed: %v", err)
				}
			}
			got := tt.readData.Size()
			want := tt.writeData.Size()
			if !reflect.DeepEqual(got, want) {
				t.Errorf("ArchOut(%#016x) got size %v, want %v", tt.addr, got, want)
			}
		})
	}
}
