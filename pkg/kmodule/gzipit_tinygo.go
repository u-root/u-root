// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build tinygo

package kmodule

import (
	"compress/gzip"
	"io"
	"os"
)

func gzipit(file *os.File) (io.Reader, error) {
	return gzip.NewReader(file)
}
