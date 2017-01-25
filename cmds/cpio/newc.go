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

type newcWriter struct {
	pos int64
	io.Writer
}

// Write implements the write interface for newcWriter.
// It allows us to track the position and round it up
// if needed. This allows us to use function such as
// io.copy and Fprintf
func (t *newcWriter) Write(b []byte) (int, error) {
	amt, err := t.Writer.Write(b)
	if err != nil {
		return -1, err
	}
	t.pos += int64(amt)
	return amt, err
}

func (t *newcWriter) advance(amt int64) error {
	pad := make([]byte, 5)
	o := round4(t.pos, amt)
	if o == t.pos {
		return nil
	}
	_, err := t.Write(pad[:o-t.pos])
	return err
}

func NewcWriter(w io.Writer) (RecWriter, error) {
	return &newcWriter{Writer: w}, nil
}

// RecWrite writes cpio records. It pads the header+name write to
// 4 byte alignment and pads the data write as well.
func (t *newcWriter) RecWrite(f *File) error {
	if _, err := t.Write([]byte(newcMagic)); err != nil {
		return err
	}

	v := reflect.ValueOf(&f.Header)
	for i := 0; i < 13; i++ {
		n := v.Elem().Field(i)
		if _, err := fmt.Fprintf(t, "%08x", n.Uint()); err != nil {
			return err
		}
	}

	if _, err := t.Write([]byte(f.Name)); err != nil {
		return err
	}
	// round to at least one byte past the name.
	if err := t.advance(1); err != nil {
		return err
	}

	if f.Data != nil {
		_, err := io.Copy(t, f.Data)
		if err != nil {
			return err
		}
		if err := t.advance(0); err != nil {
			return err
		}
	}

	return nil
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
	v := reflect.ValueOf(&f.Header)
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
	t.pos = int64(round4(t.pos, int64(f.NameSize)))

	f.Name = string(n[:f.NameSize-1])
	if f.Name == "TRAILER!!!" {
		debug("AT THE TRAILER!!!\n")
		return nil, io.EOF
	}

	f.Data = io.NewSectionReader(t, t.pos, int64(f.FileSize))
	t.pos = int64(round4(t.pos + int64(f.FileSize)))
	return &f, nil
}
