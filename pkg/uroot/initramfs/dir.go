// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package initramfs

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/cpio"
)

// DirArchiver implements Archiver for a directory.
type DirArchiver struct{}

// Reader implements Archiver.Reader.
//
// Currently unsupported for directories.
func (da DirArchiver) Reader(io.ReaderAt) Reader {
	return nil
}

// OpenWriter implements Archiver.OpenWriter.
func (da DirArchiver) OpenWriter(path, goos, goarch string) (Writer, error) {
	if len(path) == 0 {
		var err error
		path, err = ioutil.TempDir("", "u-root")
		if err != nil {
			return nil, err
		}
	} else {
		if _, err := os.Stat(path); os.IsExist(err) {
			return nil, fmt.Errorf("path %q already exists", path)
		}
		if err := os.MkdirAll(path, 0755); err != nil {
			return nil, err
		}
	}
	log.Printf("Path is %s", path)
	return dirWriter{path}, nil
}

// dirWriter implements Writer.
type dirWriter struct {
	dir string
}

// WriteRecord implements Writer.WriteRecord.
func (dw dirWriter) WriteRecord(r cpio.Record) error {
	return cpio.CreateFileInRoot(r, dw.dir)
}

// Finish implements Writer.Finish.
func (dw dirWriter) Finish() error {
	return nil
}
