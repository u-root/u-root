// Copyright 2012-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package memio

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

var tests = []struct {
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
}

func TestIO(t *testing.T) {
	for _, tt := range tests {
		t.Run(fmt.Sprintf(tt.name), func(t *testing.T) {
			tmpFile, err := ioutil.TempFile("", "io_test")
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

func ExampleRead() {
	var data Uint32
	if err := Read(0x1000000, &data); err != nil {
		log.Fatal(err)
	}
	log.Printf("%v\n", data)
}

func ExampleWrite() {
	data := Uint32(42)
	if err := Write(0x1000000, &data); err != nil {
		log.Fatal(err)
	}
}
