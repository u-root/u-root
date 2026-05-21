// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// mv renames files and directories.
//
// Synopsis:
//
//	mv SOURCE [-u] TARGET
//	mv SOURCE... [-u] DIRECTORY
//
// Author:
//
//	Beletti (rhiguita@gmail.com)
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"syscall"

	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

type flags struct {
	update    bool
	noClobber bool
}

var errUsage = errors.New("insufficient arguments")

func moveFile(source string, dest string, update bool, noClobber bool) error {
	if noClobber {
		_, err := os.Lstat(dest)
		if !errors.Is(err, os.ErrNotExist) {
			// This is either a real error if something unexpected happen during Lstat or nil
			return err
		}
	}

	if update {
		sourceInfo, err := os.Lstat(source)
		if err != nil {
			return err
		}

		destInfo, err := os.Lstat(dest)
		if err != nil {
			return err
		}

		// Check if the destination already exists and was touched later than the source
		if destInfo.ModTime().After(sourceInfo.ModTime()) {
			// Source is older and we don't want to "downgrade"
			return nil
		}
	}

	return os.Rename(source, dest)
}

func mv(files []string, update, noClobber, todir bool) error {
	if len(files) == 2 && !todir {
		// Rename/move a single file
		if err := moveFile(files[0], files[1], update, noClobber); err != nil {
			return err
		}
	} else {
		// Move one or more files into a directory
		destdir := files[len(files)-1]
		for _, f := range files[:len(files)-1] {
			newPath := filepath.Join(destdir, filepath.Base(f))
			if err := moveFile(f, newPath, update, noClobber); err != nil {
				return err
			}
		}
	}
	return nil
}

func move(stderr io.Writer, args []string) error {
	var f flags

	fs := flag.NewFlagSet("mv", flag.ContinueOnError)
	fs.SetOutput(stderr)

	fs.BoolVar(&f.update, "u", false, "move only when the SOURCE file is newer than the destination file or when the destination file is missing")
	fs.BoolVar(&f.noClobber, "n", false, "do not overwrite an existing file")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: mv [ARGS] source target [ARGS] source ... directory\n\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(unixflag.ArgsToGoArgs(args)); err != nil {
		return err
	}

	if fs.NArg() < 2 {
		fs.Usage()
		return errUsage
	}

	files := fs.Args()

	var todir bool
	dest := files[len(files)-1]

	if destdir, err := os.Lstat(dest); err == nil {
		todir = destdir.IsDir()
	}

	if len(files) > 2 && !todir {
		return fmt.Errorf("%s: %w", dest, syscall.ENOTDIR)
	}

	return mv(files, f.update, f.noClobber, todir)
}

func main() {
	if err := move(os.Stderr, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
