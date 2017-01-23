// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// basic types for cpio

package main

import "io"

type RecReader interface {
	RecRead() (*File, error)
}

type RecWriter interface {
	RecWrite(*File) error
}

type File struct {
	Ino      uint64
	Mode     uint64
	UID      uint64
	GID      uint64
	Nlink    uint64
	Mtime    uint64
	FileSize uint64
	Major    uint64
	Minor    uint64
	Rmajor   uint64
	Rminor   uint64
	NameSize uint64
	Name     string
	Data     io.Reader
}
