// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util

import (
	"bytes"
	"compress/gzip"
	"io"

	"github.com/u-root/uio/uio"
)

func readGzip(r io.Reader) ([]byte, error) {
	z, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer z.Close()
	return io.ReadAll(z)
}

// TryGzipFilter tries to read from an io.ReaderAt to see if it is a Gzip. If it is not, the
// io.ReaderAt is returned.
// TODO: could be tryCompressionFilters, grub support gz,xz and lzop
// TODO: do we want to keep the filter inside multiboot? This could be the responsibility of the caller...
func TryGzipFilter(r io.ReaderAt) io.ReaderAt {
	b, err := readGzip(uio.Reader(r))
	if err == nil {
		return bytes.NewReader(b)
	}
	return r
}
