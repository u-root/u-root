// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !tinygo

package gzip

import (
	"io"

	"github.com/klauspost/pgzip"
)

// Compress takes input from io.Reader and deflates it using pgzip
// to io.Writer. Data is compressed in block size (KB) chunks using
// up to the number of CPU cores specified.
func compress(r io.Reader, w io.Writer, level int, blocksize int, processes int) error {
	zw, err := pgzip.NewWriterLevel(w, level)
	if err != nil {
		return err
	}

	if err := zw.SetConcurrency(blocksize*1024, processes); err != nil {
		zw.Close()
		return err
	}

	if _, err := io.Copy(zw, r); err != nil {
		zw.Close()
		return err
	}

	return zw.Close()
}
