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
)

type newcReader struct {
	io.ReadSeeker
}

type newWriter struct {
	io.Writer
}

func NewcWriter(n string) (RecWriter, error) {
	return nil, fmt.Errorf("Writer: not yet")
}

func NewcReader(r io.ReadSeeker) (RecReader, error) {
	return &newcReader{ReadSeeker: r}, nil
}

func (t *newcReader) RecRead() (*File, error) {
	// There's almost certainly a better way to do this but this
	// will do for now.
	var h = make([]byte, headerLen)
	pos, err := t.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	debug("Next record: pos is %d\n", pos)

	if count, err := t.Read(h[:]); count != len(h) || err != nil {
		return nil, fmt.Errorf("Header: got %d of %d bytes, error %v", len(h), count, err)
	}
	// Make sure it's right.
	magic := string(h[:6])
	if magic != "070701" {
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
	if l, err := t.Read(n); l != int(f.NameSize) || err != nil {
		return nil, fmt.Errorf("Reading name: got %d of %d bytes, err was %v", l, f.NameSize, err)
	}

	// we have to seek to f.NameSize + len(h) rounded up to a multiple of 4.
	seekTo := (int64(pos+int64(f.NameSize)+int64(len(h))+3) / 4) * 4
	if _, err := t.Seek(seekTo, io.SeekStart); err != nil {
		return nil, err
	}

	f.Name = string(n[:f.NameSize-1])
	if f.Name == "TRAILER!!!" {
		debug("AT THE TRAILER!!!\n")
		return nil, io.EOF
	}

	f.Data = &io.LimitedReader{t, int64(f.FileSize)}
	seekTo = (int64(f.FileSize+3) / 4) * 4
	if _, err := t.Seek(seekTo, io.SeekCurrent); err != nil {
		return nil, err
	}
	return &f, nil
}
