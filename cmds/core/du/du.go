// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9 && !windows

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

	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

var (
	errNoStatInfo = errors.New("os.FileInfo has no stat_t info")
	errUsage      = errors.New("usage: du [-k] [-H | -L] [-a | -s] [file ...]")
)

type cmd struct {
	stdout            io.Writer
	files             []string
	reportFiles       bool
	kbUnit            bool
	totalSum          bool
	followCMDSymLinks bool
	followSymlinks    bool
}

func command(stdout io.Writer, args []string) *cmd {
	c := cmd{
		stdout: stdout,
	}

	f := flag.NewFlagSet(args[0], flag.ExitOnError)
	f.BoolVar(&c.reportFiles, "a", false, "report the size of each file not of type directory")
	f.BoolVar(&c.kbUnit, "k", false, "write the files sizes in units of 1024 bytes, rather than the default 512-byte units")
	f.BoolVar(&c.totalSum, "s", false, "report only the total sum for each of the specified files")
	f.BoolVar(&c.followCMDSymLinks, "H", false, "follow symlink form [file...]")
	f.BoolVar(&c.followSymlinks, "L", false, "follow all symlinks")

	f.Parse(unixflag.ArgsToGoArgs(args[1:]))
	c.files = f.Args()
	return &c
}

func (c *cmd) run() error {
	if c.totalSum && c.reportFiles {
		return errUsage
	}

	if c.followSymlinks && c.followCMDSymLinks {
		return errUsage
	}

	if len(c.files) == 0 {
		c.files = append(c.files, ".")
	}

	for _, file := range c.files {
		duPath := file
		if c.followCMDSymLinks {
			var err error
			duPath, err = c.evaluateSymlink(file)
			if err != nil {
				return err
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

		if c.followSymlinks && (info.Mode()&os.ModeSymlink != 0) {
			follow, err := c.evaluateSymlink(path)
			if err != nil {
				return err
			}
			symBlocks, err := c.du(follow)
			if err != nil {
				return err
			}
			blocks += symBlocks
			return nil
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

func (c *cmd) evaluateSymlink(path string) (string, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return "", err
	}

	symPath := path
	if fi.Mode()&os.ModeSymlink != 0 {
		symPath, err = os.Readlink(path)
		if err != nil {
			return "", err
		}

		if !filepath.IsAbs(symPath) {
			dir := filepath.Dir(path)
			symPath = filepath.Join(dir, symPath)
		}
	}

	return symPath, nil
}

func main() {
	if err := command(os.Stdout, os.Args).run(); err != nil {
		if errors.Is(err, errUsage) {
			fmt.Fprintln(os.Stderr, errUsage)
			os.Exit(1)
		}
		log.Fatalf("du: %v", err)
	}
}
