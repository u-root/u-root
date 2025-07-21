// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package logutil implements utilities for recording log output.
package logutil

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/nanmu42/limitio"
)

// NewWriterToFile creates a Writer that writes out output to a file path up to a maximum limit maxSize.
func NewWriterToFile(maxSize int, path string) (io.Writer, error) {
	logFile, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)
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

// TeeOutput tees out output to a file path specified by env var `UROOT_LOG_PATH` and sets the log output to the newly created writer.
func TeeOutput(writer io.Writer, maxSize int) (io.Writer, error) {
	writer, err := CreateTeeWriter(writer, os.Getenv("UROOT_LOG_PATH"), maxSize)
	if err == nil {
		log.SetOutput(writer)
	}
	return writer, err
}

// CreateTeeWriter tees out output to a file path specified by logPath up to a max limit. Creates necessary directories for the specified logpath if they don't exist.
func CreateTeeWriter(writer io.Writer, logPath string, maxSize int) (io.Writer, error) {
	if logPath == "" {
		return nil, fmt.Errorf("empty log path")
	}
	dir := filepath.Dir(logPath)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}
	lw, err := NewWriterToFile(maxSize, logPath)
	if err != nil {
		return nil, err
	}
	return io.MultiWriter(writer, lw), nil
}
