// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package initramfs

import (
	"fmt"
	"os"

	"github.com/u-root/mkuimage/cpio"
)

// CPIOFile opens a Reader or Writer that reads/writes files from/to a CPIO archive at the given path.
type CPIOFile struct {
	Path string
}

var _ ReadOpener = &CPIOFile{}
var _ WriteOpener = &CPIOFile{}

// OpenWriter opens c.Path for writing.
func (c *CPIOFile) OpenWriter() (Writer, error) {
	if len(c.Path) == 0 {
		return nil, fmt.Errorf("failed to write to CPIO: %w", ErrNoPath)
	}
	f, err := os.OpenFile(c.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, err
	}
	return cpioWriter{cpio.Newc.Writer(f), f}, nil
}

// OpenReader opens c.Path for reading.
func (c *CPIOFile) OpenReader() (cpio.RecordReader, error) {
	if len(c.Path) == 0 {
		return nil, fmt.Errorf("failed to read from CPIO: %w", ErrNoPath)
	}
	f, err := os.Open(c.Path)
	if err != nil {
		return nil, err
	}
	return cpio.Newc.Reader(f), nil
}

// Archive opens a Reader that reads files from an in-memory archive.
type Archive struct {
	*cpio.Archive
}

var _ ReadOpener = &Archive{}

// OpenWriter writes to the archive.
func (a *Archive) OpenWriter() (Writer, error) {
	return cpioWriter{a.Archive, nil}, nil
}

// OpenReader opens the archive for reading.
func (a *Archive) OpenReader() (cpio.RecordReader, error) {
	return a.Archive.Reader(), nil
}

// osWriter implements Writer.
type cpioWriter struct {
	cpio.RecordWriter

	f *os.File
}

// Finish implements Writer.Finish.
func (o cpioWriter) Finish() error {
	err := cpio.WriteTrailer(o)
	if o.f != nil {
		o.f.Close()
	}
	return err
}
