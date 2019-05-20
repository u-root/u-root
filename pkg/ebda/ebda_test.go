// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ebda

import (
	"bytes"
	"encoding/binary"
	"strings"
	"testing"
)

func checkErrPrefix(t *testing.T, testName string, prefix string, err error) {
	if err != nil && prefix == "" {
		t.Errorf("test %s expected no error, got %v", testName, err)
	} else if err == nil && prefix != "" {
		t.Errorf("test %s expected error beginning with %s, got nil", testName, prefix)
	} else if err != nil && prefix != "" && !strings.HasPrefix(err.Error(), prefix) {
		t.Errorf("test %s expected error beginning with %s, got %v", testName, prefix, err)
	}
}

func checkEqualsEBDA(t *testing.T, testName string, e, g *EBDA) {
	if e != nil && g == nil {
		t.Errorf("test %s expected EBDA %v, got nil", testName, *e)
	} else if e == nil && g != nil {
		t.Errorf("test %s expected no EBDA, got %v", testName, *g)
	} else if e != nil && g != nil {
		if e.BaseOffset != g.BaseOffset || e.Length != g.Length || !bytes.Equal(e.Data, g.Data) {
			t.Errorf("test %s expected EBDA: \n%v,\ngot:\n%v", testName, *e, *g)
		}
	}
}

func fakeEBDA(sizeKB int) []byte {
	b := make([]byte, sizeKB*1024)
	if sizeKB == 0 {
		return append(b, make([]byte, 1024)...) // Helps us test weird edge cases
	}
	b[0] = byte(sizeKB)
	copy(b[16:], []byte("UROOTEBDA"))
	return b
}

func fakeDevMemEBDA(offset, sizeKB int) []byte {
	b := make([]byte, offset)
	binOffset := uint16(offset >> 4)
	binary.LittleEndian.PutUint16(b[EBDAAddressOffset:EBDAAddressOffset+2], binOffset)
	b = append(b, fakeEBDA(sizeKB)...)
	return b
}

func TestFindEBDAOffset(t *testing.T) {
	for _, tt := range []struct {
		name      string
		fakeMem   []byte
		offset    int64
		errPrefix string
	}{
		{
			name:      "ReadFrom40EFail",
			fakeMem:   make([]byte, 1),
			errPrefix: "unable to read EBDA Pointer:",
		},
		{
			name:      "EmptyEBDAPointerFail",
			fakeMem:   make([]byte, EBDAAddressOffset+2),
			errPrefix: "ebda offset is 0! unable to proceed",
		},
		{
			name:    "FindEBDAOffsetSuccess",
			fakeMem: fakeDevMemEBDA(0x9000, 1),
			offset:  0x9000,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			o, err := findEBDAOffset(bytes.NewReader(tt.fakeMem))
			checkErrPrefix(t, tt.name, tt.errPrefix, err)
			if o != tt.offset {
				t.Errorf("test %s expected offset %v, got %v", tt.name, tt.offset, o)
			}
		})
	}
}

func TestReadEBDA(t *testing.T) {
	for _, tt := range []struct {
		name      string
		fakeMem   []byte
		ebda      *EBDA
		errPrefix string
	}{
		{
			name:      "ReadSizeFail",
			fakeMem:   fakeDevMemEBDA(0x9000, 1)[:0x9000], // Slice before length
			errPrefix: "error reading EBDA length, got:",
		},
		{
			name:      "ReadZeroSizeFail",
			fakeMem:   fakeDevMemEBDA(0x9000, 0), // Fails because size is zero, so it tries to read to end of segment
			errPrefix: "error reading EBDA region, tried to read from",
		},
		{
			name:    "ReadEBDA",
			fakeMem: fakeDevMemEBDA(0x9000, 1), // Fails because size is zero, so it tries to read to end of segment
			ebda: &EBDA{
				BaseOffset: 0x9000,
				Length:     1024,
				Data:       fakeEBDA(1),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			g, err := ReadEBDA(bytes.NewReader(tt.fakeMem))
			checkErrPrefix(t, tt.name, tt.errPrefix, err)
			checkEqualsEBDA(t, tt.name, tt.ebda, g)
		})
	}
}

type MockByteReadWriteSeeker struct {
	*bytes.Reader
}

func (MockByteReadWriteSeeker) Write(p []byte) (int, error) {
	return len(p), nil
}

func TestWriteEBDA(t *testing.T) {
	for _, tt := range []struct {
		name      string
		fakeMem   []byte
		ebda      *EBDA
		errPrefix string
	}{
		{
			name:    "NoBaseOffsetFail",
			fakeMem: make([]byte, 1),
			ebda: &EBDA{
				Length: 1024,
				Data:   fakeEBDA(1),
			},
			errPrefix: "unable to read EBDA Pointer:",
		},
		{
			name:    "UnalignedSizeFail",
			fakeMem: fakeDevMemEBDA(0x9000, 1),
			ebda: &EBDA{
				Length: 1023,
				Data:   fakeEBDA(1),
			},
			errPrefix: "length is not an integer multiple of 1 KiB, got",
		},
		{
			name:    "MismatchedSizeFail",
			fakeMem: fakeDevMemEBDA(0x9000, 1),
			ebda: &EBDA{
				BaseOffset: 0x9000,
				Length:     1024,
				Data:       fakeEBDA(2),
			},
			errPrefix: "length field is not equal to buffer length",
		},
		{
			name:    "SizeTooBigFail",
			fakeMem: fakeDevMemEBDA(0x9000, 1),
			ebda: &EBDA{
				BaseOffset: 0x9000,
				Length:     256 * 1024,
				Data:       fakeEBDA(256),
			},
			errPrefix: "length is greater than 255 KiB",
		},
		{
			name:    "WriteEBDA",
			fakeMem: fakeDevMemEBDA(0x9000, 1),
			ebda: &EBDA{
				BaseOffset: 0x9000,
				Length:     1024,
				Data:       fakeEBDA(1),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := WriteEBDA(tt.ebda, MockByteReadWriteSeeker{bytes.NewReader(tt.fakeMem)})
			checkErrPrefix(t, tt.name, tt.errPrefix, err)
		})
	}
}
