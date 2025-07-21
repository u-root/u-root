// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build tinygo

package gzip

import (
	"compress/gzip"
	"io"
)

// Decompress takes gzip compressed input from io.Reader and expands it using pgzip
// to io.Writer. Data is read in blocksize (KB) chunks using upto the number of
// CPU cores specified.
func decompress(r io.Reader, w io.Writer, _, processes int) error {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}

	if _, err := io.Copy(w, zr); err != nil {
		zr.Close()
		return err
	}

	return zr.Close()
}
