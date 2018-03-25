// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpio

import (
	"fmt"
	"io"
	"path/filepath"
)

var (
	formatMap = make(map[string]RecordFormat)
	Debug     = func(string, ...interface{}) {}
)

// A RecordReader reads one record from an archive.
type RecordReader interface {
	ReadRecord() (Record, error)
}

// A RecordWriter writes on record to an archive.
type RecordWriter interface {
	WriteRecord(Record) error
}

// A RecordFormat gives readers and writers for dealing with archives.
type RecordFormat interface {
	Reader(r io.ReaderAt) RecordReader
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

// EOFReader is a RecordReader that converts the Trailer record to io.EOF.
type EOFReader struct {
	RecordReader
}

// ReadRecord implements RecordReader.
//
// ReadRecord returns io.EOF when the record name is TRAILER!!!.
func (r EOFReader) ReadRecord() (Record, error) {
	rec, err := r.RecordReader.ReadRecord()
	if err != nil {
		return Record{}, err
	}
	// The end of a CPIO archive is marked by a record whose name is
	// "TRAILER!!!".
	if rec.Name == Trailer {
		return Record{}, io.EOF
	}
	return rec, nil
}

// DedupWriter is a RecordWriter that does not write more than one record with
// the same path.
type DedupWriter struct {
	rw RecordWriter

	// There seems to be no harm done in stripping
	// duplicate names when the record is written,
	// and lots of harm done if we don't do it.
	alreadyWritten map[string]struct{}
}

// NewDedupWriter returns a new deduplicating rw.
func NewDedupWriter(rw RecordWriter) RecordWriter {
	return &DedupWriter{
		rw:             rw,
		alreadyWritten: make(map[string]struct{}),
	}
}

// WriteRecord implements RecordWriter.
//
// If rec.Name was already seen once before, it will not be written again and
// WriteRecord returns nil.
func (dw *DedupWriter) WriteRecord(rec Record) error {
	// We do NOT write records with absolute paths.
	if filepath.IsAbs(rec.Name) {
		// There's no constant that means "root". PathSeparator is not
		// really quite right.
		rel, err := filepath.Rel("/", rec.Name)
		if err != nil {
			return fmt.Errorf("Can't make %s relative to /?", rec.Name)
		}
		rec.Name = rel
	}

	if _, ok := dw.alreadyWritten[rec.Name]; ok {
		return nil
	}
	dw.alreadyWritten[rec.Name] = struct{}{}
	return dw.rw.WriteRecord(rec)
}

// WriteRecords writes multiple records.
func WriteRecords(w RecordWriter, files []Record) error {
	for _, f := range files {
		if err := w.WriteRecord(f); err != nil {
			return fmt.Errorf("WriteRecords: writing %q got %v", f.Info.Name, err)
		}
	}
	return nil
}

// WriteTrailer writes the trailer record.
func WriteTrailer(w RecordWriter) error {
	return w.WriteRecord(TrailerRecord)
}

// Concat reads files from r one at a time, and writes them to w.
func Concat(w RecordWriter, r RecordReader, transform func(Record) Record) error {
	// Read and write one file at a time. We don't want all that in memory.
	for {
		f, err := r.ReadRecord()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if transform != nil {
			f = transform(f)
		}
		if err := w.WriteRecord(f); err != nil {
			return err
		}
	}
}

// ReadAllRecords returns all records in `r` in the order in which they were
// read.
func ReadAllRecords(r RecordReader) ([]Record, error) {
	var files []Record
	for {
		f, err := r.ReadRecord()
		if err == io.EOF {
			return files, nil
		}
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
}

// MakeReproducible changes any fields in a Record such that
// if we run cpio again, with the same files presented to it
// in the same order, and those files have unchanged contents,
// the cpio file it produces will be bit-for-bit
// identical. This is an essential property for firmware-embedded
// payloads.
func MakeReproducible(file Record) Record {
	file.MTime = 0
	return file
}

func MakeAllReproducible(files []Record) {
	for i := range files {
		files[i] = MakeReproducible(files[i])
	}
}
