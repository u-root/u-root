// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpio

import (
	"fmt"
	"io"
	"log"
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
	return Writer{a.RecordFormat.Writer(w)}
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
	return files, nil
}

type Writer struct {
	rw RecordWriter
}

func (w Writer) WriteRecord(rec Record) error {
	return w.rw.WriteRecord(rec)
}

func (w Writer) WriteRecords(files []Record) error {
	for _, f := range files {
		if err := w.WriteRecord(f); err != nil {
			return fmt.Errorf("WriteRecords: writing %q got %v", f.Info.Name, err)
		}
	}
	return nil
}

func (w Writer) WriteTrailer() error {
	return w.WriteRecord(TrailerRecord)
}
