// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package multiboot

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"

	"github.com/u-root/u-root/pkg/uio"
)

func readGzip(r io.Reader) ([]byte, error) {
	z, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer z.Close()
	return ioutil.ReadAll(z)
}

// TODO: could be tryCompressionFilters, grub support gz,xz and lzop
// TODO: do we want to keep the filter inside multiboot? This could be the responsibility of the caller...
func tryGzipFilter(r io.ReaderAt) io.ReaderAt {
	b, err := readGzip(uio.Reader(r))
	if err == nil {
		return bytes.NewReader(b)
	}
	return r
}
