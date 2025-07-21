// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpio

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"path"
	"strings"

	"github.com/u-root/uio/uio"
)

// Trailer is the name of the trailer record.
const Trailer = "TRAILER!!!"

// TrailerRecord is the last record in any CPIO archive.
var TrailerRecord = StaticRecord(nil, Info{Name: Trailer})

// StaticRecord returns a record with the given contents and metadata.
func StaticRecord(contents []byte, info Info) Record {
	info.FileSize = uint64(len(contents))
	return Record{
		ReaderAt: bytes.NewReader(contents),
		Info:     info,
	}
}

// StaticFile returns a normal file record.
func StaticFile(name string, content string, perm uint64) Record {
	return StaticRecord([]byte(content), Info{
		Name: name,
		Mode: S_IFREG | perm,
	})
}

// Symlink returns a symlink record at name pointing to target.
func Symlink(name string, target string) Record {
	return Record{
		ReaderAt: strings.NewReader(target),
		Info: Info{
			FileSize: uint64(len(target)),
			Mode:     S_IFLNK | 0o777,
			Name:     name,
		},
	}
}

// Directory returns a directory record at name.
func Directory(name string, mode uint64) Record {
	return Record{
		Info: Info{
			Name: name,
			Mode: S_IFDIR | mode&^S_IFMT,
		},
	}
}

// CharDev returns a character device record at name.
func CharDev(name string, perm uint64, rmajor, rminor uint64) Record {
	return Record{
		Info: Info{
			Name:   name,
			Mode:   S_IFCHR | perm,
			Rmajor: rmajor,
			Rminor: rminor,
		},
	}
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
//
// There seems to be no harm done in stripping duplicate names when the record
// is written, and lots of harm done if we don't do it.
type DedupWriter struct {
	rw RecordWriter

	// alreadyWritten keeps track of paths already written to rw.
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
	rec.Name = Normalize(rec.Name)

	if _, ok := dw.alreadyWritten[rec.Name]; ok {
		return nil
	}
	dw.alreadyWritten[rec.Name] = struct{}{}
	return dw.rw.WriteRecord(rec)
}

// WriteRecords writes multiple records to w.
func WriteRecords(w RecordWriter, files []Record) error {
	for _, f := range files {
		if err := w.WriteRecord(f); err != nil {
			return fmt.Errorf("WriteRecords: writing %q got %w", f.Info.Name, err)
		}
	}
	return nil
}

// WriteRecordsAndDirs writes records to w, with a slight difference from WriteRecords:
// the record path is split and all the
// directories are written first, in order, mimic'ing what happens with
// find . -print
//
// When is this function needed?
// Most cpio programs will create directories as needed for paths such as a/b/c/d
// The cpio creation process for Linux uses find, and will create a
// record for each directory in a/b/c/d
//
// But when code programatically generates a cpio for the Linux kernel,
// the cpio is not generated via find, and Linux will not create
// intermediate directories. The result, seen in practice, is that a path,
// such as a/b/c/d, when unpacked by the linux kernel, will be ignored if
// a/b/c does not exist!
//
// Again, this function is very rarely needed, save when we programatically generate
// an initramfs for Linux.
// This code only works with a deduplicating writer. Further, it will not accept a
// Record if the full pathname of that Record already exists. This is arguably
// overly restrictive but, at the same, avoids some very unpleasant programmer
// errors.
// There is overlap here with DedupWriter but given that this is a Special Snowflake
// function, it seems best to leave the DedupWriter code alone.
func WriteRecordsAndDirs(rw RecordWriter, files []Record) error {
	w, ok := rw.(*DedupWriter)
	if !ok {
		return fmt.Errorf("WriteRecordsAndDirs(%T,...): only DedupWriter allowed:%w", rw, os.ErrInvalid)
	}
	for _, f := range files {
		// This redundant Normalize does no harm, but, yes, it is redundant.
		// Signed
		// The Department of Redundancy Department.
		f.Name = Normalize(f.Name)
		if r, ok := w.alreadyWritten[f.Name]; ok {
			return fmt.Errorf("WriteRecordsAndDirs: %q already in the archive: %v:%w", f.Name, r, os.ErrExist)
		}

		var recs []Record
		// Paths must be written to the archive in the order in which they
		// need to be created, i.e., a/b/c/d must be written as
		// a, a/b/, a/b/c, a/b/c/d
		// Note: do not use os.Separator here: cpio is a Unix standard, and hence
		// / is used.
		// do NOT use filepath, use path for the same reason.
		// Things you learn the hard way when you run on Windows.
		els := strings.Split(path.Dir(f.Name), "/")
		for i := range els {
			d := path.Join(els[:i+1]...)
			recs = append(recs, Directory(d, 0o777))
		}
		recs = append(recs, f)
		if err := WriteRecords(rw, recs); err != nil {
			return fmt.Errorf("WriteRecords: writing %q got %w", f.Info.Name, err)
		}
	}
	return nil
}

// Passthrough copies from a RecordReader to a RecordWriter.
//
// Passthrough writes a trailer record.
//
// It processes one record at a time to minimize the memory footprint.
func Passthrough(r RecordReader, w RecordWriter) error {
	if err := Concat(w, r, nil); err != nil {
		return err
	}
	if err := WriteTrailer(w); err != nil {
		return err
	}
	return nil
}

// WriteTrailer writes the trailer record.
func WriteTrailer(w RecordWriter) error {
	return w.WriteRecord(TrailerRecord)
}

// Concat reads files from r one at a time, and writes them to w.
//
// Concat does not write a trailer record and applies transform to every record
// before writing it. transform may be nil.
func Concat(w RecordWriter, r RecordReader, transform func(Record) Record) error {
	return ForEachRecord(r, func(f Record) error {
		if transform != nil {
			f = transform(f)
		}
		return w.WriteRecord(f)
	})
}

// ReadAllRecords returns all records in r in the order in which they were
// read.
func ReadAllRecords(rr RecordReader) ([]Record, error) {
	var files []Record
	err := ForEachRecord(rr, func(r Record) error {
		files = append(files, r)
		return nil
	})
	return files, err
}

// ForEachRecord reads every record from r and applies f.
func ForEachRecord(rr RecordReader, fun func(Record) error) error {
	for {
		rec, err := rr.ReadRecord()
		switch err {
		case io.EOF:
			return nil

		case nil:
			if err := fun(rec); err != nil {
				return err
			}

		default:
			return err
		}
	}
}

// Normalize normalizes namepath to be relative to /.
func Normalize(name string) string {
	// do not use filepath.IsAbs, it will not work on Windows.
	// do not use filepath.Rel, that will not work
	// sensibly on windows.
	// The only thing one can do is strip all leading
	// /
	name = strings.TrimLeft(name, "/")
	// do not use filepath.Clean here.
	// This will result in paths with \\ on windows, and
	// / is the cpio standard.
	return path.Clean(name)
}

// MakeReproducible changes any fields in a Record such that if we run cpio
// again, with the same files presented to it in the same order, and those
// files have unchanged contents, the cpio file it produces will be bit-for-bit
// identical. This is an essential property for firmware-embedded payloads.
func MakeReproducible(r Record) Record {
	// Do NOT zero Ino. The Ino is created in a reproducible manner
	// and a non-zero value is critical for creating hard links when
	// reading the archive.
	// r.Ino = 0
	r.Name = Normalize(r.Name)
	r.MTime = 0
	r.UID = 0
	r.GID = 0
	r.Dev = 0
	r.Major = 0
	r.Minor = 0
	// Consider that a file may have 10 links,
	// but we are only including 1: NLink will be incorrect. In the
	// general case, it is almost impossible to set NLink correctly.
	if r.NLink > 1 {
		r.NLink = math.MaxUint64
	}
	return r
}

// MakeAllReproducible makes all given records reproducible as in
// MakeReproducible.
func MakeAllReproducible(files []Record) {
	for i := range files {
		files[i] = MakeReproducible(files[i])
	}
}

// AllEqual compares all metadata and contents of r and s.
func AllEqual(r []Record, s []Record) bool {
	if len(r) != len(s) {
		return false
	}
	for i := range r {
		if !Equal(r[i], s[i]) {
			return false
		}
	}
	return true
}

// Equal compares the metadata and contents of r and s.
func Equal(r Record, s Record) bool {
	if r.Info != s.Info {
		return false
	}
	return uio.ReaderAtEqual(r.ReaderAt, s.ReaderAt)
}
