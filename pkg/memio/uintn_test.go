// Copyright 2012-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package memio

import (
	"fmt"
	"reflect"
	"testing"
)

var uintn_tests = []struct {
	name      string
	addr      int64
	writeData UintN
	want      string
	err       string
}{
	{
		name:      "uint8",
		addr:      0x10,
		writeData: &[]Uint8{0x12}[0],
		want:      "0x12",
	},
	{
		name:      "uint16",
		addr:      0x20,
		writeData: &[]Uint16{0x1234}[0],
		want:      "0x1234",
	},
	{
		name:      "uint32",
		addr:      0x30,
		writeData: &[]Uint32{0x12345678}[0],
		want:      "0x12345678",
	},
	{
		name:      "uint64",
		addr:      0x40,
		writeData: &[]Uint64{0x1234567890abcdef}[0],
		want:      "0x1234567890abcdef",
	},
	{
		name:      "byte slice",
		addr:      0x50,
		writeData: &[]ByteSlice{[]byte("Hello")}[0],
		want:      "0x48656c6c6f",
	},
}

func TestString(t *testing.T) {
	for _, tt := range uintn_tests {
		t.Run(fmt.Sprintf(tt.name), func(t *testing.T) {
			got := tt.writeData.String()
			want := tt.want
			if !reflect.DeepEqual(want, got) {
				t.Errorf("Want %v, got %v", want, got)
			}
		},
		)
	}
}
