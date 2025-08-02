// Copyright 2012-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package memio

import "testing"

func TestUint8(t *testing.T) {
	for _, tt := range []struct {
		name       string
		valueUint8 Uint8
		value8     uint8
		expString  string
		expSize    int64
	}{
		{
			name:       "uint8",
			valueUint8: Uint8(45),
			value8:     uint8(45),
			expString:  "0x2d",
			expSize:    1,
		},
	} {
		t.Run(tt.name+"Size8", func(t *testing.T) {
			if tt.valueUint8.Size() != tt.expSize {
				t.Errorf("%q failed. Got: %d, Want: %d", tt.name+"Size8", tt.valueUint8.Size(), tt.expSize)
			}
		})

		t.Run(tt.name+"String8", func(t *testing.T) {
			if tt.valueUint8.String() != tt.expString {
				t.Errorf("%q failed. Got: %q, Want: %q", tt.name+"String8", tt.valueUint8.String(), tt.expString)
			}
		})
	}
}

func TestUint16(t *testing.T) {
	for _, tt := range []struct {
		name      string
		value     Uint16
		value16   uint16
		expString string
		expSize   int64
	}{
		{
			name:      "uint16",
			value:     Uint16(45),
			value16:   uint16(45),
			expString: "0x002d",
			expSize:   2,
		},
	} {
		t.Run(tt.name+"Size16", func(t *testing.T) {
			if tt.value.Size() != tt.expSize {
				t.Errorf("%q failed. Got: %d, Want: %d", tt.name+"Size16", tt.value.Size(), tt.expSize)
			}
		})

		t.Run(tt.name+"String16", func(t *testing.T) {
			if tt.value.String() != tt.expString {
				t.Errorf("%q failed. Got: %q, Want: %q", tt.name+"String16", tt.value.String(), tt.expString)
			}
		})
	}
}

func TestUint32(t *testing.T) {
	for _, tt := range []struct {
		name      string
		value     Uint32
		value32   uint32
		expString string
		expSize   int64
	}{
		{
			name:      "uint32",
			value:     Uint32(45),
			value32:   uint32(45),
			expString: "0x0000002d",
			expSize:   4,
		},
	} {
		t.Run(tt.name+"Size32", func(t *testing.T) {
			if tt.value.Size() != tt.expSize {
				t.Errorf("%q failed. Got: %d, Want: %d", tt.name+"Size32", tt.value.Size(), tt.expSize)
			}
		})

		t.Run(tt.name+"String32", func(t *testing.T) {
			if tt.value.String() != tt.expString {
				t.Errorf("%q failed. Got: %q, Want: %q", tt.name+"String32", tt.value.String(), tt.expString)
			}
		})
	}
}

func TestUint64(t *testing.T) {
	for _, tt := range []struct {
		name      string
		value     Uint64
		value64   uint64
		expString string
		expSize   int64
	}{
		{
			name:      "uint64",
			value:     Uint64(45),
			value64:   uint64(45),
			expString: "0x000000000000002d",
			expSize:   8,
		},
	} {
		t.Run(tt.name+"Size64", func(t *testing.T) {
			if tt.value.Size() != tt.expSize {
				t.Errorf("%q failed. Got: %d, Want: %d", tt.name+"Size", tt.value.Size(), tt.expSize)
			}
		})

		t.Run(tt.name+"String64", func(t *testing.T) {
			if tt.value.String() != tt.expString {
				t.Errorf("%q failed. Got: %q, Want: %q", tt.name+"String", tt.value.String(), tt.expString)
			}
		})
	}
}

func TestUintByteSlice(t *testing.T) {
	for _, tt := range []struct {
		name      string
		value     ByteSlice
		valueByte []byte
		expString string
		expSize   int64
	}{
		{
			name:      "ByteSlice1",
			value:     ByteSlice([]byte{1}),
			valueByte: []byte{1},
			expString: "0x01",
			expSize:   1,
		},
		{
			name:      "ByteSlice2",
			value:     ByteSlice([]byte{1, 2}),
			valueByte: []byte{1, 2},
			expString: "0x0102",
			expSize:   2,
		},
		{
			name:      "ByteSlice4",
			value:     ByteSlice([]byte{1, 2, 3, 4}),
			valueByte: []byte{1, 2, 3, 4},
			expString: "0x01020304",
			expSize:   4,
		},
		{
			name:      "ByteSlice8",
			value:     ByteSlice([]byte{1, 2, 3, 4, 5, 6, 7, 8}),
			valueByte: []byte{1, 2, 3, 4, 5, 6, 7, 8},
			expString: "0x0102030405060708",
			expSize:   8,
		},
	} {
		t.Run(tt.name+"SizeByteSlice", func(t *testing.T) {
			if tt.value.Size() != tt.expSize {
				t.Errorf("%q failed. Got: %d, Want: %d", tt.name+"SizeByteSlice", tt.value.Size(), tt.expSize)
			}
		})

		t.Run(tt.name+"StringByteSlice", func(t *testing.T) {
			if tt.value.String() != tt.expString {
				t.Errorf("%q failed. Got: %q, Want: %q", tt.name+"StringByteSlice", tt.value.String(), tt.expString)
			}
		})
	}
}
