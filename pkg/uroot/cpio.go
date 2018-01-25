// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uroot

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/cpio"
	_ "github.com/u-root/u-root/pkg/cpio/newc"
)

// CPIOArchiver is an implementation of Archiver for the cpio format.
type CPIOArchiver struct {
	// Format is the name of the cpio format to use.
	Format string
}

// OpenWriter opens `path` as the correct file type and returns an
// ArchiveWriter pointing to `path`.
//
// If `path` is empty, a default path of /tmp/initramfs.GOOS_GOARCH.cpio is
// used.
func (ca CPIOArchiver) OpenWriter(path, goos, goarch string) (ArchiveWriter, error) {
	if len(path) == 0 {
		path = fmt.Sprintf("/tmp/initramfs.%s_%s.cpio", goos, goarch)
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return nil, err
	}
	log.Printf("Filename is %s", path)
	archiver, err := cpio.Format(ca.Format)
	if err != nil {
		return nil, err
	}

	return osWriter{archiver.Writer(f), f}, nil
}

// osWriter implements ArchiveWriter.
type osWriter struct {
	cpio.Writer
	f *os.File
}

func (o osWriter) Finish() error {
	err := o.WriteTrailer()
	o.f.Close()
	return err
}

// Reader implements Archiver.Reader.
func (ca CPIOArchiver) Reader(r io.ReaderAt) ArchiveReader {
	archiver, err := cpio.Format(ca.Format)
	if err != nil {
		panic(err)
	}

	return archiver.Reader(r)
}
