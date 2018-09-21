// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpio

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/u-root/u-root/pkg/ls"
	"github.com/u-root/u-root/pkg/uio"
	"golang.org/x/sys/unix"
)

var (
	formatMap = make(map[string]RecordFormat)
	Debug     = func(string, ...interface{}) {}
)

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

// Passthrough copies from a RecordReader to a RecordWriter
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

// Archive implements RecordWriter and is an in-memory list of files.
//
// Archive.Reader returns a RecordReader for this archive.
type Archive struct {
	// Files is a map of relative archive path -> record.
	Files map[string]Record

	// Order is a list of relative archive paths and represents the order
	// in which Files were added.
	Order []string
}

// InMemArchive returns an in-memory file archive.
func InMemArchive() *Archive {
	return &Archive{
		Files: make(map[string]Record),
	}
}

func ArchiveFromRecords(rs []Record) *Archive {
	a := InMemArchive()
	for _, r := range rs {
		a.WriteRecord(r)
	}
	return a
}

// WriteRecord implements RecordWriter and adds a record to the archive.
func (a *Archive) WriteRecord(r Record) error {
	r.Name = Normalize(r.Name)
	a.Files[r.Name] = r
	a.Order = append(a.Order, r.Name)
	return nil
}

// Empty returns whether the archive has any files in it.
func (a *Archive) Empty() bool {
	return len(a.Files) == 0
}

type archiveReader struct {
	a   *Archive
	pos int
}

// Reader returns a RecordReader for the archive that starts at the first
// record.
func (a *Archive) Reader() RecordReader {
	return &EOFReader{&archiveReader{a: a}}
}

func modeFromLinux(mode uint64) os.FileMode {
	m := os.FileMode(mode & 0777)
	switch mode & syscall.S_IFMT {
	case syscall.S_IFBLK:
		m |= os.ModeDevice
	case syscall.S_IFCHR:
		m |= os.ModeDevice | os.ModeCharDevice
	case syscall.S_IFDIR:
		m |= os.ModeDir
	case syscall.S_IFIFO:
		m |= os.ModeNamedPipe
	case syscall.S_IFLNK:
		m |= os.ModeSymlink
	case syscall.S_IFREG:
		// nothing to do
	case syscall.S_IFSOCK:
		m |= os.ModeSocket
	}
	if mode&syscall.S_ISGID != 0 {
		m |= os.ModeSetgid
	}
	if mode&syscall.S_ISUID != 0 {
		m |= os.ModeSetuid
	}
	if mode&syscall.S_ISVTX != 0 {
		m |= os.ModeSticky
	}
	return m
}

// LSInfoFromRecord converts a Record to be usable with the ls package for
// listing files.
func LSInfoFromRecord(rec Record) ls.FileInfo {
	var target string

	mode := modeFromLinux(rec.Mode)
	if mode&os.ModeType == os.ModeSymlink {
		if l, err := uio.ReadAll(rec); err != nil {
			target = err.Error()
		} else {
			target = string(l)
		}
	}

	return ls.FileInfo{
		Name:          rec.Name,
		Mode:          mode,
		Rdev:          unix.Mkdev(uint32(rec.Rmajor), uint32(rec.Rminor)),
		UID:           uint32(rec.UID),
		GID:           uint32(rec.GID),
		Size:          int64(rec.FileSize),
		MTime:         time.Unix(int64(rec.MTime), 0).UTC(),
		SymlinkTarget: target,
	}
}

// String implements fmt.Stringer.
//
// String lists files like ls would.
func (a *Archive) String() string {
	var b strings.Builder
	r := a.Reader()
	for {
		record, err := r.ReadRecord()
		if err != nil {
			return b.String()
		}
		b.WriteString(record.String())
		b.WriteString("\n")
	}
}

// ReadRecord implements RecordReader.
func (ar *archiveReader) ReadRecord() (Record, error) {
	if ar.pos >= len(ar.a.Order) {
		return Record{}, io.EOF
	}

	path := ar.a.Order[ar.pos]
	ar.pos++
	return ar.a.Files[path], nil
}

// Contains returns true if a record matching r is in the archive.
func (a *Archive) Contains(r Record) bool {
	r.Name = Normalize(r.Name)
	if s, ok := a.Files[r.Name]; ok {
		return Equal(r, s)
	}
	return false
}

func (a *Archive) Get(path string) (Record, bool) {
	r, ok := a.Files[Normalize(path)]
	return r, ok
}

// ReadArchive reads the entire archive in-memory and makes it accessible by
// paths.
func ReadArchive(rr RecordReader) (*Archive, error) {
	a := &Archive{
		Files: make(map[string]Record),
	}
	err := ForEachRecord(rr, func(r Record) error {
		return a.WriteRecord(r)
	})
	return a, err
}

// ReadAllRecords returns all records in `r` in the order in which they were
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

// Normalize normalizes path to be relative to /.
func Normalize(path string) string {
	if filepath.IsAbs(path) {
		rel, err := filepath.Rel("/", path)
		if err != nil {
			panic("absolute filepath must be relative to /")
		}
		return rel
	}
	return path
}

// MakeReproducible changes any fields in a Record such that if we run cpio
// again, with the same files presented to it in the same order, and those
// files have unchanged contents, the cpio file it produces will be bit-for-bit
// identical. This is an essential property for firmware-embedded payloads.
func MakeReproducible(r Record) Record {
	r.Ino = 0
	r.Name = Normalize(r.Name)
	r.MTime = 0
	r.UID = 0
	r.GID = 0
	r.Dev = 0
	r.Major = 0
	r.Minor = 0
	r.NLink = 0
	return r
}

// MakeAllReproducible makes all given records reproducible as in
// MakeReproducible.
func MakeAllReproducible(files []Record) {
	for i := range files {
		files[i] = MakeReproducible(files[i])
	}
}
