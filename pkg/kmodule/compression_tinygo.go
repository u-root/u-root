// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build tinygo
// +build tinygo

package kmodule

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/klauspost/compress/zstd"
	"github.com/ulikunitz/xz"
)

// compressionReader returns a reader for the given file based on the file extension.
// The current status of `tinygo` does not support the use of inline assembly, so for now we can only use pure go libraries.
// Although klauspost/zstd is pure go, interanlly it executes inline assembly to get CPU features.
func compressionReader(file *os.File) (reader io.Reader, err error) {
	ext := filepath.Ext(file.Name())

	switch ext {
	case ".xz":
		// xz is pure go so we can use that here
		return xz.NewReader(file)
	case ".gz":
		return gzip.NewReader(file)
	case ".zst":
		return zstd.NewReader(file)
	default:
		return nil, fmt.Errorf("compression not supported for %s:%w", ext, os.ErrNotExist)
	}
}
