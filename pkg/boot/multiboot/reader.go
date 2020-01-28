// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package multiboot

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type kernelReader struct {
	buf []byte
	off int
}

func (kr kernelReader) ReadAt(p []byte, off int64) (n int, err error) {
	if off < 0 || off > int64(len(kr.buf)) {
		return 0, fmt.Errorf("bad offset %v", off)
	}
	if n = copy(p, kr.buf[off:]); n < len(p) {
		err = io.EOF
	}
	return n, err
}

func (kr *kernelReader) Read(p []byte) (n int, err error) {
	if n = copy(p, kr.buf[kr.off:]); n < len(p) {
		err = io.EOF
	}
	kr.off += n
	return n, err
}

func readGzip(r io.Reader) ([]byte, error) {
	z, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer z.Close()
	return ioutil.ReadAll(z)
}

func readFile(name string) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := readGzip(f)
	if err == nil {
		return b, err
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("cannot rewind file: %v", err)
	}

	return ioutil.ReadAll(f)
}
