// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// basic types for cpio

package cpio

import (
	"io"
)

const (
	// bad news. These are not defined on
	// whatever version of Go travis is running.
	SeekSet     = 0
	SeekCurrent = 1
)

type CpioReader struct {
	Pos int64
	io.ReaderAt
}

type CpioWriter struct {
	Pos int64
	io.Writer
}

type Info struct {
	Name string
	Ino      uint64
	Mode     uint64
	UID      uint64
	GID      uint64
	Nlink    uint64
	Mtime    uint64
	FileSize uint64
	Major    uint64
	Minor    uint64
	Rmajor   uint64
	Rminor   uint64
}

type File struct {
	Info
	Data io.Reader
}

type RecReader interface {
       RecRead() (*File, error)
}

type RecWriter interface {
	RecWrite(*File) (int, error)
	Finish() error
}

type NewReader func (io.ReaderAt) (RecReader, error)
type NewWriter func (io.Writer) (RecWriter, error)
type ops struct {
	NewReader
	NewWriter
}

