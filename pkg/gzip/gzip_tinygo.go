// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build tinygo

package gzip

import (
	"compress/gzip"
	"io"
)

// compress takes input from io.Reader and deflates it using pgzip.
// to io.Writer. Data is compressed in blocksize (KB) chunks using
// up to the number of CPU cores specified.
func compress(r io.Reader, w io.Writer, level, _, _ int) error {
	zw, err := gzip.NewWriterLevel(w, level)
	if err != nil {
		return err
	}

	if _, err := io.Copy(zw, r); err != nil {
		zw.Close()
		return err
	}

	return zw.Close()
}
