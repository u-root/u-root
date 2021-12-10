// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package multiboot

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"testing"
)

func createFile(hdr *header, offset, size int) (io.Reader, error) {
	buf := bytes.Repeat([]byte{0xDE, 0xAD, 0xBE, 0xEF}, (size+4)/4)
	buf = buf[:size]
	if hdr != nil {
		w := bytes.Buffer{}
		if err := binary.Write(&w, binary.LittleEndian, *hdr); err != nil {
			return nil, err
		}
		copy(buf[offset:], w.Bytes())
	}
	return bytes.NewReader(buf), nil
}

type testFlag string

const (
	flagGood        testFlag = "good"
	flagUnsupported testFlag = "unsup"
	flagBad         testFlag = "bad"
)

func createHeader(fl testFlag) header {
	flags := headerFlag(0x00000002)
	var checksum uint32
	switch fl {
	case flagGood:
		checksum = 0xFFFFFFFF - headerMagic - uint32(flags) + 1
	case flagBad:
		checksum = 0xDEADBEEF
	case flagUnsupported:
		flags = 0x0000FFFC
		checksum = 0xFFFFFFFF - headerMagic - uint32(flags) + 1
	}

	return header{
		mandatory: mandatory{
			Magic:    headerMagic,
			Flags:    flags,
			Checksum: checksum,
		},
		optional: optional{
			HeaderAddr:  1,
			LoadAddr:    2,
			LoadEndAddr: 3,
			BSSEndAddr:  4,
			EntryAddr:   5,

			ModeType: 6,
			Width:    7,
			Height:   8,
			Depth:    9,
		},
	}
}

func TestParseHeader(t *testing.T) {
	mandatorySize := binary.Size(mandatory{})
	optionalSize := binary.Size(optional{})
	sizeofHeader := mandatorySize + optionalSize

	for _, test := range []struct {
		flags  testFlag
		offset int
		size   int
		err    error
	}{
		{flags: flagGood, offset: 0, size: 8192, err: nil},
		{flags: flagGood, offset: 2048, size: 8192, err: nil},
		{flags: flagGood, offset: 8192 - sizeofHeader - 4, size: 8192, err: nil},
		{flags: flagGood, offset: 8192 - sizeofHeader - 1, size: 8192, err: ErrHeaderNotFound},
		{flags: flagGood, offset: 8192 - sizeofHeader, size: 8192, err: nil},
		{flags: flagGood, offset: 8192 - 4, size: 8192, err: ErrHeaderNotFound},
		{flags: flagGood, offset: 8192, size: 16384, err: ErrHeaderNotFound},
		{flags: flagGood, offset: 0, size: 10, err: io.ErrUnexpectedEOF},
		{flags: flagBad, offset: 0, size: 8192, err: ErrHeaderNotFound},
		{flags: flagUnsupported, offset: 0, size: 8192, err: ErrFlagsNotSupported},
		{flags: flagGood, offset: 8192 - mandatorySize, size: 8192, err: nil},
	} {
		t.Run(fmt.Sprintf("flags:%v,off:%v,sz:%v,err:%v", test.flags, test.offset, test.size, test.err), func(t *testing.T) {
			want := createHeader(test.flags)
			r, err := createFile(&want, test.offset, test.size)
			if err != nil {
				t.Fatalf("Cannot create test file: %v", err)
			}
			got, err := parseHeader(r)
			if err != test.err {
				t.Fatalf("parseHeader() got error: %v, want: %v", err, test.err)
			}

			if err != nil {
				return
			}
			if test.size-test.offset > mandatorySize {
				if !reflect.DeepEqual(*got, want) {
					t.Errorf("parseHeader() got %+v, want %+v", *got, want)
				}
			} else {
				if !reflect.DeepEqual(got.mandatory, want.mandatory) {
					t.Errorf("parseHeader() got %+v, want %+v", got.mandatory, want.mandatory)
				}
			}
		})
	}
}
