// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cpio implements utilities for reading and writing cpio archives.
//
// Currently, only newc-formatted cpio archives are supported through cpio.Newc.
//
// Reading from or writing to a file:
//
//	f, err := os.Open(...)
//	if err ...
//	recReader := cpio.Newc.Reader(f)
//	err := ForEachRecord(recReader, func(r cpio.Record) error {
//
//	})
//
//	// Or...
//	recWriter := cpio.Newc.Writer(f)
//
// Reading from or writing to an in-memory archive:
//
//	a := cpio.InMemArchive()
//	err := a.WriteRecord(...)
//
//	recReader := a.Reader() // Reads from the "beginning."
//
//	if a.Contains("bar/foo") {
//
//	}
package cpio

import (
	"fmt"
	"io"
	"os"
	"time"
)

var (
	formatMap = make(map[string]RecordFormat)

	// Debug can be set e.g. to log.Printf to enable debug prints from
	// marshaling/unmarshaling cpio archives.
	Debug = func(string, ...interface{}) {}
)

// Record represents a CPIO record, which represents a Unix file.
type Record struct {
	// ReaderAt contains the content of this CPIO record.
	io.ReaderAt

	// Info is metadata describing the CPIO record.
	Info

	// metadata about this item's place in the file
	RecPos  int64  // Where in the file this record is
	RecLen  uint64 // How big the record is.
	FilePos int64  // Where in the CPIO the file's contents are.
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

// A RecordReader reads one record from an archive.
type RecordReader interface {
	ReadRecord() (Record, error)
}

// A RecordWriter writes one record to an archive.
type RecordWriter interface {
	WriteRecord(Record) error
}

// A RecordFormat gives readers and writers for dealing with archives from io
// objects.
//
// CPIO files have a number of records, of which newc is the most widely used
// today.
type RecordFormat interface {
	Reader(r io.ReaderAt) RecordReader
	FileReader(f *os.File) RecordReader
	Writer(w io.Writer) RecordWriter
}

// Format returns the RecordFormat with that name, if it exists.
func Format(name string) (RecordFormat, error) {
	op, ok := formatMap[name]
	if !ok {
		return nil, fmt.Errorf("%q is not in cpio format map %v", name, formatMap)
	}
	return op, nil
}
