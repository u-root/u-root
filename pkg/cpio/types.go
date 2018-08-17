// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpio

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/u-root/u-root/pkg/ls"
	"github.com/u-root/u-root/pkg/uio"
	"golang.org/x/sys/unix"
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

// String implements a fmt.Stringer for Record.
//
// String returns a string formatted like `ls` would format it.
func (r Record) String() string {
	s := ls.LongStringer{
		Human: true,
		Name:  ls.NameStringer{},
	}
	return s.FileString(LSInfoFromRecord(r))
}

// Trailer is the name of the trailer record.
const Trailer = "TRAILER!!!"

// TrailerRecord is the last record in any CPIO archive.
var TrailerRecord = StaticRecord(nil, Info{Name: Trailer})

func StaticRecord(contents []byte, info Info) Record {
	info.FileSize = uint64(len(contents))
	return Record{
		ReaderAt: bytes.NewReader(contents),
		Info:     info,
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
		ReaderAt: strings.NewReader(target),
		Info: Info{
			FileSize: uint64(len(target)),
			Mode:     unix.S_IFLNK | 0777,
			Name:     name,
		},
	}
}

// Directory returns a directory record at name.
func Directory(name string, mode uint64) Record {
	return Record{
		Info: Info{
			Name: name,
			Mode: unix.S_IFDIR | mode&^unix.S_IFMT,
		},
	}
}

// CharDev returns a character device record at name.
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

func NewLazyFile(name string) io.ReaderAt {
	return uio.NewLazyOpenerAt(func() (io.ReaderAt, error) {
		return os.Open(name)
	})
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

func Equal(r Record, s Record) bool {
	if r.Info != s.Info {
		return false
	}
	return ReaderAtEqual(r.ReaderAt, s.ReaderAt)
}

func ReaderAtEqual(r1, r2 io.ReaderAt) bool {
	var c, d []byte
	var err error
	if r1 != nil {
		c, err = uio.ReadAll(r1)
		if err != nil {
			return false
		}
	}

	if r2 != nil {
		d, err = uio.ReadAll(r2)
		if err != nil {
			return false
		}
	}
	return bytes.Equal(c, d)
}
