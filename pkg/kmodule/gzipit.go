// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !tinygo

package kmodule

import (
	"io"
	"os"

	"github.com/klauspost/pgzip"
)

func gzipit(file *os.File) (io.Reader, error) {
	return pgzip.NewReader(file)
}
