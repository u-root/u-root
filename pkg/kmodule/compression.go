// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kmodule

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/klauspost/compress/zstd"

	"github.com/ulikunitz/xz"
)

// compressionReader returns a reader for the given file based on the file extension.
// This implementation utilizes the optimized compression libraries for each format.
func compressionReader(file *os.File) (reader io.Reader, err error) {
	ext := filepath.Ext(file.Name())

	switch ext {
	case ".ko":
		return file, nil
	case ".xz":
		return xz.NewReader(file)
	case ".gz":
		return gzipit(file)
	case ".zst":
		return zstd.NewReader(file)
	default:
		return nil, fmt.Errorf("compression not supported for %s:%w", ext, os.ErrNotExist)
	}
}
