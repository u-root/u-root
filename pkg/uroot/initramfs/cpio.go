// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package initramfs

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/cpio"
)

// CPIOArchiver is an implementation of Archiver for the cpio format.
type CPIOArchiver struct {
	cpio.RecordFormat
}

// OpenWriter opens `path` as the correct file type and returns an
// Writer pointing to `path`.
//
// If `path` is empty, a default path of /tmp/initramfs.GOOS_GOARCH.cpio is
// used.
func (ca CPIOArchiver) OpenWriter(path, goos, goarch string) (Writer, error) {
	if len(path) == 0 && len(goos) == 0 && len(goarch) == 0 {
		return nil, fmt.Errorf("passed no path, GOOS, and GOARCH to CPIOArchiver.OpenWriter")
	}
	if len(path) == 0 {
		path = fmt.Sprintf("/tmp/initramfs.%s_%s.cpio", goos, goarch)
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return nil, err
	}
	log.Printf("Filename is %s", path)
	return osWriter{ca.RecordFormat.Writer(f), f}, nil
}

// osWriter implements Writer.
type osWriter struct {
	cpio.RecordWriter

	f *os.File
}

// Finish implements Writer.Finish.
func (o osWriter) Finish() error {
	err := cpio.WriteTrailer(o)
	o.f.Close()
	return err
}

// Reader implements Archiver.Reader.
func (ca CPIOArchiver) Reader(r io.ReaderAt) Reader {
	return ca.RecordFormat.Reader(r)
}
