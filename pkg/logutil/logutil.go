// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package logutil implements utilities for recording log output.
package logutil

import (
	"io"
	"os"
	"path/filepath"

	"github.com/nanmu42/limitio"
)

// NewWriterToFile creates a Writer that writes out output to a file path up to a maximum limit maxSize.
func NewWriterToFile(maxSize int, dirPath, name string) (io.Writer, error) {
	logFile, err := os.OpenFile(filepath.Join(dirPath, name), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	fi, err := logFile.Stat()
	if err != nil {
		return nil, err
	}

	lw := limitio.NewWriter(logFile, maxSize-(int)(fi.Size()), true)

	return lw, nil

}
