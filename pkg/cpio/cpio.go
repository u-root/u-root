// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpio

import (
	"fmt"
	"io"
	"log"
	"path/filepath"
)

var (
	formatMap = make(map[string]RecordFormat)
	formats   []string
	Debug     = func(string, ...interface{}) {}
)

func AddFormat(name string, f RecordFormat) {
	if _, ok := formatMap[name]; ok {
		log.Fatalf("cpio: two requests for format %s", name)
	}

	formatMap[name] = f
	formats = append(formats, name)
}

func Format(name string) (Archiver, error) {
	op, ok := formatMap[name]
	if !ok {
		return Archiver{}, fmt.Errorf("Format %v is not one of %v", name, formats)
	}
	return Archiver{op}, nil
}

type Archiver struct {
	RecordFormat
}

func (a Archiver) Reader(r io.ReaderAt) Reader {
	return Reader{a.RecordFormat.Reader(r)}
}

func (a Archiver) Writer(w io.Writer) Writer {
	return Writer{rw: a.RecordFormat.Writer(w), alreadyWritten: make(map[string]struct{})}
}

type Reader struct {
	rr RecordReader
}

func (r Reader) ReadRecord() (Record, error) {
	rec, err := r.rr.ReadRecord()
	if err != nil {
		return Record{}, err
	}
	// The end of a CPIO archive is marked by a record whose name is "TRAILER!!!".
	if rec.Name == Trailer {
		return Record{}, io.EOF
	}
	return rec, nil
}

func (r Reader) ReadRecords() ([]Record, error) {
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

type Writer struct {
	rw RecordWriter

	// There seems to be no harm done in stripping
	// duplicate names when the record is written,
	// and lots of harm done if we don't do it.
	alreadyWritten map[string]struct{}
}

func (w Writer) WriteRecord(rec Record) error {
	// we do NOT write records with absolute paths.
	if filepath.IsAbs(rec.Name) {
		// There's no constant that means "root".
		// PathSeparator is not really quite right.
		rel, err := filepath.Rel("/", rec.Name)
		if err != nil {
			return fmt.Errorf("Can't make %s relative to /?", rec.Name)
		}
		rec.Name = rel
	}

	if _, ok := w.alreadyWritten[rec.Name]; ok {
		return nil
	}
	w.alreadyWritten[rec.Name] = struct{}{}
	return w.rw.WriteRecord(rec)
}

// WriteRecords writes multiple records.
func (w Writer) WriteRecords(files []Record) error {
	for _, f := range files {
		if err := w.WriteRecord(f); err != nil {
			return fmt.Errorf("WriteRecords: writing %q got %v", f.Info.Name, err)
		}
	}
	return nil
}

// WriteTrailer writes the trailer record.
func (w Writer) WriteTrailer() error {
	return w.WriteRecord(TrailerRecord)
}

// Concat reads files from r one at a time, and writes them to w.
func (w Writer) Concat(r Reader, transform func(Record) Record) error {
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
