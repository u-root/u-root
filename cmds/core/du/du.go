// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9

package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"syscall"
)

func run(stdout io.Writer, files ...string) error {
	if len(files) == 0 {
		files = append(files, ".")
	}

	for _, file := range files {
		blocks, err := du(file)
		if err != nil {
			return err
		}
		fmt.Fprintf(stdout, "%d\t%s\n", blocks, file)
	}

	return nil
}

func du(file string) (int64, error) {
	var blocks int64

	filepath.Walk(file, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		st, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("can not stat: %v", path)
		}

		blocks += st.Blocks
		return nil
	})

	return blocks, nil
}

func main() {
	if err := run(os.Stdout, os.Args[1:]...); err != nil {
		log.Fatalf("du: %v", err)
	}
}
