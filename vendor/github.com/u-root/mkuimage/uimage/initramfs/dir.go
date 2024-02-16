// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package initramfs

import (
	"errors"
	"fmt"
	"os"

	"github.com/u-root/mkuimage/cpio"
)

// Dir opens a Writer that writes all archive files to the given directory.
type Dir struct {
	Path string
}

var _ WriteOpener = &Dir{}

// OpenWriter implements Archiver.OpenWriter.
func (d *Dir) OpenWriter() (Writer, error) {
	if len(d.Path) == 0 {
		return nil, fmt.Errorf("failed to use directory as output: %w", ErrNoPath)
	}
	if err := os.MkdirAll(d.Path, 0o755); err != nil && !errors.Is(err, os.ErrExist) {
		return nil, err
	}
	return dirWriter{d.Path}, nil
}

// dirWriter implements Writer.
type dirWriter struct {
	dir string
}

// WriteRecord implements Writer.WriteRecord.
func (dw dirWriter) WriteRecord(r cpio.Record) error {
	return cpio.CreateFileInRoot(r, dw.dir, false)
}

// Finish implements Writer.Finish.
func (dw dirWriter) Finish() error {
	return nil
}
