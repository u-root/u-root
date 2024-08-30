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

var (
	errNoStatInfo = errors.New("os.FileInfo has no stat_t info")
	errUsage      = errors.New("usage: du [-k] [-H] [-a | -s] [file ...]")
)

type cmd struct {
	stdout            io.Writer
	reportFiles       bool
	kbUnit            bool
	totalSum          bool
	followCMDSymLinks bool
}

func command(stdout io.Writer, reportFiles, kbUnit, totalSum, followCMDSymLinks bool) *cmd {
	return &cmd{
		stdout:            stdout,
		reportFiles:       reportFiles,
		kbUnit:            kbUnit,
		totalSum:          totalSum,
		followCMDSymLinks: followCMDSymLinks,
	}
}

func (c *cmd) run(files ...string) error {
	if c.totalSum && c.reportFiles {
		return errUsage
	}

	if len(files) == 0 {
		files = append(files, ".")
	}

	for _, file := range files {
		duPath := file
		if c.followCMDSymLinks {
			fi, err := os.Lstat(file)
			if err != nil {
				return err
			}

			if fi.Mode()&os.ModeSymlink != 0 {
				duPath, err = os.Readlink(file)
				if err != nil {
					return err
				}
			}
		}

		blocks, err := c.du(duPath)
		if err != nil {
			return err
		}
		c.print(blocks, file)
	}

	return nil
}

func (c *cmd) du(file string) (int64, error) {
	var blocks int64

	filepath.Walk(file, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// report sub-folders and add number of blocks to overall count
		if info.IsDir() && file != path && !c.totalSum {
			dirBlocks, err := c.du(path)
			if err != nil {
				return err
			}

			blocks += dirBlocks

			c.print(dirBlocks, path)
			return fs.SkipDir
		}

		st, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("%v: %w", path, errNoStatInfo)
		}

		if c.reportFiles && !info.IsDir() && !c.totalSum {
			c.print(st.Blocks, path)
		}

		blocks += st.Blocks
		return nil
	})

	return blocks, nil
}

func (c *cmd) print(nblock int64, path string) {
	if c.kbUnit {
		nblock /= 2
	}
	fmt.Fprintf(c.stdout, "%d\t%s\n", nblock, path)
}

func main() {
	var reportFiles = flag.Bool("a", false, "report the size of each file not of type directory")
	var kbUnit = flag.Bool("k", false, "write the files sizes in units of 1024 bytes, rather than the default 512-byte units")
	var totalSum = flag.Bool("s", false, "report only the total sum for each of the specified files")
	var followCMDSymLinks = flag.Bool("H", false, "follow symlink form [file...]")
	flag.Parse()
	if err := command(os.Stdout, *reportFiles, *kbUnit, *totalSum, *followCMDSymLinks).run(flag.Args()...); err != nil {
		if errors.Is(err, errUsage) {
			fmt.Fprintln(os.Stderr, errUsage)
			os.Exit(1)
		}
		log.Fatalf("du: %v", err)
	}
}
