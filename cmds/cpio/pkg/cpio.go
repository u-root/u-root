// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpio

import (
	"fmt"
	"log"
	"io"
	"os"
)

const (
	typeMask    = 0170000 //This masks the file type bits.
	typeSocket  = 0140000 //File type value for sockets.
	typeSymLink = 0120000 //File type value for symbolic links.  For symbolic links, the link body is stored as file data.
	typeFile    = 0100000 //File type value for regular files.
	typeBlock   = 0060000 // File type value for block special devices.
	typeDir     = 0040000 // File type value for directories.
	typeChar    = 0020000 // File type value for character special devices.
	typeFIFO    = 0010000 // File type value for named pipes or FIFOs.
	SUID        = 0004000 // SUID bit.
	SGID        = 0002000 // SGID bit.
	sticky      = 0001000 // Sticky bit.  On some systems, this modifies the behavior
	// of executables and/or directories.
	mode = 0000777 //The lower 9 bits specify read/write/execute permissions
	//for world, group, and user following standard POSIX con-
	//ventions.
)

var (
	ModeMap = map[uint64]os.FileMode{
		typeSocket:  os.ModeSocket,
		typeSymLink: os.ModeNamedPipe,
		typeFile:    os.FileMode(0),
		typeBlock:   os.ModeDevice,
		typeDir:     os.ModeDir,
		typeChar:    os.ModeCharDevice,
		typeFIFO:    os.ModeNamedPipe,
	}

	FormatMap = map[string] *ops {
	}
	Formats string
	Debug = func(string, ...interface{}) {}
)

func AddMap(n string, r NewReader, w NewWriter) {
	if _, ok := FormatMap[n]; ok {
		log.Fatalf("cpio: two requests for format %s", n)
	}

	FormatMap[n] = &ops{r, w}
	Formats += n + " "
}

func Reader(n string, r io.ReaderAt) (RecReader, error) {
	op, ok := FormatMap[n]
	if ! ok {
		return nil, fmt.Errorf("Format %v is not one of %v", n, Formats)
	}
	return op.NewReader(r)
}

func Writer(n string, w io.Writer) (RecWriter, error) {
	op, ok := FormatMap[n]
	if ! ok {
		return nil, fmt.Errorf("Format %v is not one of %v", n, Formats)
	}
	return op.NewWriter(w)
}

func perm(f *File) uint32 {
	return uint32(f.Mode) & mode
}

func dev(f *File) int {
	return int(f.Rmajor<<8 | f.Rminor)
}

func cpioModetoMode(m uint64) (os.FileMode, error) {
	if t, ok := ModeMap[m&typeMask]; ok {
		return t, nil
	}
	return os.FileMode(0), fmt.Errorf("Invalid file type %#x", m&typeMask)
}

func (f *File) String() string {
	return fmt.Sprintf("%s: Ino %d Mode %#o UID %d GID %d Nlink %d Mtime %#x FileSize %d Major %d Minor %d RMajor %d Rminor %d",
		f.Name,
		f.Ino,
		f.Mode,
		f.UID,
		f.GID,
		f.Nlink,
		// what a mess. This fails on travis.
		//time.Unix(int64(f.Mtime), 0).String(),
		f.Mtime,
		f.FileSize,
		f.Major,
		f.Minor,
		f.Rmajor,
		f.Rminor)
}

func RecWriteAll(w RecWriter, f []*File) (int, error) {
	var tot int 
	for _, wf := range f {
		n, err := w.RecWrite(wf)
		if err != nil {
			return -1, fmt.Errorf("TestReadWrite: writing got %v, want nil", err)
		}
		tot += n
	}
	return tot, nil
}

// RecReadAll reads all the File records.
// This interface setup seems broken. This should work
// for all record types. If we ever get more than one
// we will have to revisit this.
func RecReadAll(r RecReader) ([]*File, error) {
	var f []*File
	for {
		nf, err := r.RecRead()
		if err == io.EOF {
			return f, nil
		}
		if err != nil {
			return nil, err
		}
		f = append(f, nf)
	}
}

