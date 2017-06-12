// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpio

import (
	"bytes"
	"fmt"
	"io"
	"time"
)

const Trailer = "TRAILER!!!"

type Record struct {
	io.Reader
	Info
}

var TrailerRecord = EmptyRecord(Info{Name: Trailer})

type RecordReader interface {
	ReadRecord() (Record, error)
}

type RecordWriter interface {
	WriteRecord(Record) error
}

type RecordFormat interface {
	Reader(r io.ReaderAt) RecordReader
	Writer(w io.Writer) RecordWriter
}

func StaticRecord(contents string, info Info) Record {
	info.FileSize = uint64(len(contents))
	return Record{
		bytes.NewReader([]byte(contents)),
		info,
	}
}

type eofReader struct{}

func (eofReader) Read(p []byte) (int, error) {
	return 0, io.EOF
}

func EmptyRecord(info Info) Record {
	return Record{
		eofReader{},
		info,
	}
}

// Info holds metadata about files.
type Info struct {
	Ino      uint64
	Mode     uint64
	UID      uint64
	GID      uint64
	NLink    uint64
	MTime    uint64
	FileSize uint64
	Dev      uint64
	Major    uint64
	Minor    uint64
	Rmajor   uint64
	Rminor   uint64
	Name     string
}

func (i Info) String() string {
	return fmt.Sprintf("%s: Ino %d Mode %#o UID %d GID %d NLink %d MTime %v FileSize %d Major %d Minor %d Rmajor %d Rminor %d",
		i.Name,
		i.Ino,
		i.Mode,
		i.UID,
		i.GID,
		i.NLink,
		time.Unix(int64(i.MTime), 0).UTC(),
		i.FileSize,
		i.Major,
		i.Minor,
		i.Rmajor,
		i.Rminor)
}
