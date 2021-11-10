// Copyright 2012-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package memio

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
