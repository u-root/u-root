// Copyright 2013 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"time"
)

type file struct {
	Name    string
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
	IsDir   bool
	Link    string

	// Yes, this is Unix-specific. But u-root in general is very Unix-specific
	// and I can't get that worried about it.
	Uid int
	Gid int
	Dev uint64
	// The offset is not exported.
	// We want the VTOC at the head of the file so we can stream the file.
	// We don't save the offset that is not known until we know the size
	// of the VTOC, which then changes its size with many encodings
	// (gob, JSON, etc.)
	// So we compute it once we read in the VTOC -- from the
	// file structs themselves.
	offset int64
}

type vtoc struct {
	f    os.File
	vtoc []*file
}

type VTOCOpt func(*vtoc) error
