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
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/u-root/u-root/pkg/uroot/util"
)

const usage = "mv [ARGS] source target [ARGS] source ... directory"

var (
	update    = flag.Bool("u", false, "move only when the SOURCE file is newer than the destination file or when the destination file is missing")
	noClobber = flag.Bool("n", false, "do not overwrite an existing file")
)

func moveFile(source string, dest string, update bool, noClobber bool) error {
	if noClobber {
		_, err := os.Lstat(dest)
		if !os.IsNotExist(err) {
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

	if err := os.Rename(source, dest); err != nil {
		return err
	}
	return nil
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

func move(files []string, update, noClobber bool) error {
	var todir bool
	dest := files[len(files)-1]
	if destdir, err := os.Lstat(dest); err == nil {
		todir = destdir.IsDir()
	}
	if len(files) > 2 && !todir {
		return fmt.Errorf("not a directory: %s", dest)
	}
	return mv(files, update, noClobber, todir)
}

func main() {
	flag.Usage = util.Usage(flag.Usage, usage)
	flag.Parse()
	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}
	if err := move(flag.Args(), *update, *noClobber); err != nil {
		log.Fatal(err)
	}
}
