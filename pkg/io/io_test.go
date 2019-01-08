// Copyright 2012-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

var tests = []struct {
	name      string
	addr      int64
	writeData interface{}
	readData  interface{}
	err       string
}{
	{
		name:      "uint8",
		addr:      0x10,
		writeData: uint8(0x12),
		readData:  new(uint8),
	},
	{
		name:      "uint16",
		addr:      0x20,
		writeData: uint16(0x1234),
		readData:  new(uint16),
	},
	{
		name:      "uint32",
		addr:      0x30,
		writeData: uint32(0x12345678),
		readData:  new(uint32),
	},
	{
		name:      "uint64",
		addr:      0x40,
		writeData: uint64(0x1234567890abcdef),
		readData:  new(uint64),
	},
	{
		name:      "bad write type",
		addr:      0,
		writeData: int8(0),
		err:       "cannot write type int8",
	},
	{
		name:      "bad read type",
		addr:      0,
		writeData: uint8(0),
		readData:  int8(0),
		err:       "cannot read type int8",
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

			want := fmt.Sprintf("%#016x", tt.writeData)
			got := fmt.Sprintf("%#016x", reflect.Indirect(reflect.ValueOf(tt.readData)).Interface())
			if got != want {
				t.Fatalf("Write(%#016x, %s) = %s; want %s",
					tt.addr, want, got, want)
			}
		})
	}
}

func ExampleRead() {
	var data uint32
	if err := Read(0x1000000, &data); err != nil {
		log.Fatal(err)
	}
	log.Printf("%#08x\n", data)
}

func ExampleWrite() {
	if err := Write(0x1000000, uint32(42)); err != nil {
		log.Fatal(err)
	}
}
