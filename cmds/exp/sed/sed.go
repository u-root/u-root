// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// By Ahmed Kamal <email.ahmedkamal@googlemail.com>

package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Writes to a temp file first
// Renames to the final file upon closing
type tmpWriter struct {
	filename string
	ftmp     os.File
}

func newTmpWriter(filename string) (*tmpWriter, error) {
	ftmp, err := os.CreateTemp("/tmp", ".sed*.txt")
	if err != nil {
		return nil, fmt.Errorf("unable to create temp file: %w", err)
	}
	return &tmpWriter{filename: filename, ftmp: *ftmp}, nil
}

func (tw *tmpWriter) Write(b []byte) (int, error) {
	return tw.ftmp.Write(b)
}

func (tw *tmpWriter) Close() error {
	err := tw.ftmp.Close()
	if err != nil {
		return err
	}
	os.Rename(tw.ftmp.Name(), tw.filename)
	return nil
}

func main() {
	tw, err := newTmpWriter("/tmp/tw")
	if err != nil {
		log.Fatalf("unable to create temp writer: %#v", err)
	}
	defer tw.Close()
	fmt.Fprintf(tw, "some data\n")
	time.Sleep(1 * time.Second)
}
