// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !tinygo

package gzip

import (
	"io"

	"github.com/klauspost/pgzip"
)

// Decompress takes gzip compressed input from io.Reader and expands it using pgzip
// to io.Writer. Data is read in block size (KB) chunks using up to the number of
// CPU cores specified.
func decompress(r io.Reader, w io.Writer, blocksize int, processes int) error {
	zr, err := pgzip.NewReaderN(r, blocksize*1024, processes)
	if err != nil {
		return err
	}

	if _, err := io.Copy(w, zr); err != nil {
		zr.Close()
		return err
	}

	return zr.Close()
}
