// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uio

import (
	"io"
	"io/ioutil"
	"math"
	"os"

	"golang.org/x/sys/unix"
)

// InMemReaderAt is an io.ReaderAt for which ReadAtll will return a pointer to
// the bytes in memory without copying them.
type InMemReaderAt interface {
	io.ReaderAt

	Bytes() []byte
}

// ReadAll reads everything that r contains.
//
// Callers *must* not modify bytes in the returned byte slice.
//
// If r is an in-memory representation, ReadAll will attempt to return a
// pointer to those bytes directly.
func ReadAll(r io.ReaderAt) ([]byte, error) {
	if imra, ok := r.(InMemReaderAt); ok {
		return imra.Bytes(), nil
	}
	return ioutil.ReadAll(Reader(r))
}

// Reader generates a Reader from a ReaderAt.
func Reader(ra io.ReaderAt) io.Reader {
	if r, ok := ra.(io.Reader); ok {
		return r
	}
	return io.NewSectionReader(ra, 0, math.MaxInt64)
}

type inMemFile struct {
	*os.File

	bytes []byte
}

// InMemFile mmaps a file.
func InMemFile(f *os.File) (InMemReaderAt, error) {
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	b, err := unix.Mmap(int(f.Fd()), 0, int(fi.Size()), unix.PROT_READ, unix.MAP_PRIVATE)
	if err != nil {
		return nil, err
	}
	return &inMemFile{
		File:  f,
		bytes: b,
	}, nil
}

// Bytes implements InMemReaderAt.
func (i *inMemFile) Bytes() []byte {
	return i.bytes
}
