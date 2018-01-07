// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpio

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"golang.org/x/sys/unix"
)

const Trailer = "TRAILER!!!"

type Record struct {
	io.ReadCloser
	Info
}

var TrailerRecord = StaticRecord(nil, Info{Name: Trailer})

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

func StaticRecord(contents []byte, info Info) Record {
	info.FileSize = uint64(len(contents))
	return Record{
		ReadCloser: ioutil.NopCloser(bytes.NewReader(contents)),
		Info:       info,
	}
}

func StaticFile(name string, content string, perm uint64) Record {
	return StaticRecord([]byte(content), Info{
		Name: name,
		Mode: unix.S_IFREG | perm,
	})
}

// Symlink returns a symlink record at name pointing to target.
func Symlink(name string, target string) Record {
	return Record{
		ReadCloser: ioutil.NopCloser(bytes.NewReader([]byte(target))),
		Info: Info{
			FileSize: uint64(len(target)),
			Mode:     unix.S_IFLNK | 0777,
			Name:     name,
		},
	}
}

// Directory returns a directory record at `name`.
func Directory(name string, mode uint64) Record {
	return Record{
		Info: Info{
			Name: name,
			Mode: unix.S_IFDIR | mode&^unix.S_IFMT,
		},
	}
}

// CharDev returns a character device record at `name`.
func CharDev(name string, perm uint64, rmajor, rminor uint64) Record {
	return Record{
		Info: Info{
			Name:   name,
			Mode:   unix.S_IFCHR | perm,
			Rmajor: rmajor,
			Rminor: rminor,
		},
	}
}

func NewBytesReadCloser(contents []byte) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader(contents))
}

func NewReadCloser(r io.Reader) io.ReadCloser {
	return ioutil.NopCloser(r)
}

type LazyOpen struct {
	Name string
	File *os.File
}

func (r *LazyOpen) Read(p []byte) (int, error) {
	if r.File == nil {
		f, err := os.Open(r.Name)
		if err != nil {
			return -1, err
		}
		r.File = f
	}
	return r.File.Read(p)
}

func (r *LazyOpen) Close() error {
	return r.File.Close()
}

func NewDeferReadCloser(name string) io.ReadCloser {
	return &LazyOpen{Name: name}
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

func Equal(r Record, s Record) bool {
	if r.Info != s.Info {
		return false
	}
	return ReadCloserEqual(r.ReadCloser, s.ReadCloser)
}

func ReadCloserEqual(r1, r2 io.ReadCloser) bool {
	var c, d []byte
	var err error
	if r1 != nil {
		c, err = ioutil.ReadAll(r1)
		if err != nil {
			return false
		}
	}

	if r2 != nil {
		d, err = ioutil.ReadAll(r2)
		if err != nil {
			return false
		}
	}
	return bytes.Equal(c, d)
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
