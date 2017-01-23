// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// newc implements the interface for new type cpio files.
package main

import (
	"fmt"
	"io"
	"reflect"
)

const (
	magicLen  = 6
	headerLen = 13*8 + magicLen
	// bad news. These are not defined on
	// whatever version of Go travis is running.
	SeekSet     = 0
	SeekCurrent = 1
	newcMagic   = "070701"
)

type newcReader struct {
	pos int64
	io.ReaderAt
}

type newWriter struct {
	io.Writer
}

func NewcWriter(n string) (RecWriter, error) {
	return nil, fmt.Errorf("Writer: not yet")
}

func NewcReader(r io.ReaderAt) (RecReader, error) {
	m := io.NewSectionReader(r, 0, 6)
	var magic [6]byte
	if _, err := m.Read(magic[:]); err != nil {
		return nil, fmt.Errorf("NewcReader: unable to read magic: %v", err)
	}
	if string(magic[:]) != newcMagic {
		return nil, fmt.Errorf("NewcReader: magic is '%s' and must be '%s'", magic, newcMagic)
	}
	return &newcReader{ReaderAt: r}, nil
}

func (t *newcReader) RecRead() (*File, error) {
	// There's almost certainly a better way to do this but this
	// will do for now.
	var h = make([]byte, headerLen)

	debug("Next record: pos is %d\n", t.pos)

	if count, err := t.ReadAt(h[:], t.pos); count != len(h) || err != nil {
		return nil, fmt.Errorf("Header: got %d of %d bytes, error %v", len(h), count, err)
	}
	t.pos += int64(len(h))
	// Make sure it's right.
	magic := string(h[:6])
	if magic != newcMagic {
		return nil, fmt.Errorf("Reader: magic '%s' not a newc file", magic)
	}

	debug("Header is %v\n", h)
	var f File
	v := reflect.ValueOf(&f)
	for i := 0; i < 12; i++ {
		var n uint64
		f := v.Elem().Field(i)
		_, err := fmt.Sscanf(string(h[i*8+6:(i+1)*8+6]), "%x", &n)
		if err != nil {
			return nil, err
		}
		f.SetUint(n)
	}
	debug("f is %s\n", (&f).String())
	var n = make([]byte, f.NameSize)
	if l, err := t.ReadAt(n, t.pos); l != int(f.NameSize) || err != nil {
		return nil, fmt.Errorf("Reading name: got %d of %d bytes, err was %v", l, f.NameSize, err)
	}

	// we have to seek to f.NameSize + len(h) rounded up to a multiple of 4.
	t.pos = int64(round4(uint64(t.pos), f.NameSize))

	f.Name = string(n[:f.NameSize-1])
	if f.Name == "TRAILER!!!" {
		debug("AT THE TRAILER!!!\n")
		return nil, io.EOF
	}

	f.Data = io.NewSectionReader(t, t.pos, int64(f.FileSize))
	t.pos = int64(round4(uint64(t.pos) + uint64(f.FileSize)))
	return &f, nil
}
