// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9 && !windows

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"syscall"
)

func run(stdout io.Writer) error {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return err
	}

	s, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return fmt.Errorf("stat stdin:Sys() is empty: %w", os.ErrNotExist)
	}

	return filepath.WalkDir("/dev", func(path string, dir os.DirEntry, _ error) error {
		switch path {
		case "/dev/fd":
			// On some systems, /dev/fd is a directory. On Linux, /dev/fd is a symlink.
			// If /dev/fd is a directory, it should not be walked, so return file path.SkipDir.
			if dir.IsDir() {
				return filepath.SkipDir
			}
		case "/dev/stdin", "/dev/stdout", "/dev/stderr":
			return nil
		}
		if fi, err := os.Stat(path); err == nil {
			stat, ok := fi.Sys().(*syscall.Stat_t)
			if ok {
				if stat.Ino == s.Ino && stat.Dev == s.Dev {
					fmt.Fprintln(stdout, path)
				}
			}

		}

		return nil
	})
}

func main() {
	if err := run(os.Stdin); err != nil {
		log.Fatal(err)
	}
}
