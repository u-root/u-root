// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package logutil implements utilities for recording log output.
package logutil

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/nanmu42/limitio"
)

// NewWriterToFile creates a Writer that writes out output to a file path up to a maximum limit maxSize.
func NewWriterToFile(maxSize int, path string) (io.Writer, error) {
	logFile, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
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

// TeeOutput tees out output to a file path specified by env var `UROOT_LOG_PATH` up to a max limit. Creates necessary directories for the specified logpath if they don't exist.
func TeeOutput(writer io.Writer, maxSize int) (io.Writer, error) {
	logPath := os.Getenv("UROOT_LOG_PATH")
	if logPath != "" {
		dir := filepath.Dir(logPath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			log.Printf("Log directory %s doesn't exist, creating...", dir)
			if err := os.MkdirAll(dir, 0700); err != nil {
				return nil, err
			}
		}
		lw, err := NewWriterToFile(maxSize, logPath)
		if err != nil {
			return nil, err
		}
		writer = io.MultiWriter(writer, lw)
		log.SetOutput(writer)
	}
	return writer, nil
}
