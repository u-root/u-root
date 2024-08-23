// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"syscall"
)

var errNoStatInfo = errors.New("os.FileInfo has no stat_t info")

type cmd struct {
	stdout      io.Writer
	reportFiles bool
}

func command(stdout io.Writer, reportFiles bool) *cmd {
	return &cmd{
		stdout:      stdout,
		reportFiles: reportFiles,
	}
}

func (c *cmd) run(files ...string) error {
	if len(files) == 0 {
		files = append(files, ".")
	}

	for _, file := range files {
		blocks, err := c.du(file)
		if err != nil {
			return err
		}
		fmt.Fprintf(c.stdout, "%d\t%s\n", blocks, file)
	}

	return nil
}

func (c *cmd) du(file string) (int64, error) {
	var blocks int64

	filepath.Walk(file, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		st, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("%v: %w", path, errNoStatInfo)
		}

		if c.reportFiles && !info.IsDir() {
			fmt.Fprintf(c.stdout, "%d\t%s\n", st.Blocks, path)
		}

		blocks += st.Blocks
		return nil
	})

	return blocks, nil
}

func main() {
	var reportFiles = flag.Bool("a", false, "report the size of each file not of type directory")
	flag.Parse()
	if err := command(os.Stdout, *reportFiles).run(flag.Args()...); err != nil {
		log.Fatalf("du: %v", err)
	}
}
