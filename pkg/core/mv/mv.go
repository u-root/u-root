// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package mv implements the mv core utility.
package mv

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/u-root/u-root/pkg/core"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

// Command implements the mv command.
type Command struct {
	core.Base
}

// New creates a new mv command.
func New() core.Command {
	c := &Command{}
	c.Init()
	return c
}

type flags struct {
	update    bool
	noClobber bool
}

func (c *Command) moveFile(source string, dest string, update bool, noClobber bool) error {
	source = c.ResolvePath(source)
	dest = c.ResolvePath(dest)

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

func (c *Command) mv(files []string, update, noClobber, todir bool) error {
	if len(files) == 2 && !todir {
		// Rename/move a single file
		if err := c.moveFile(files[0], files[1], update, noClobber); err != nil {
			return err
		}
	} else {
		// Move one or more files into a directory
		destdir := files[len(files)-1]
		for _, f := range files[:len(files)-1] {
			newPath := filepath.Join(destdir, filepath.Base(f))
			if err := c.moveFile(f, newPath, update, noClobber); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Command) move(files []string, update, noClobber bool) error {
	var todir bool
	dest := files[len(files)-1]
	dest = c.ResolvePath(dest)
	if destdir, err := os.Lstat(dest); err == nil {
		todir = destdir.IsDir()
	}
	if len(files) > 2 && !todir {
		return fmt.Errorf("not a directory: %s", dest)
	}
	return c.mv(files, update, noClobber, todir)
}

// Run executes the mv command.
func (c *Command) Run(ctx context.Context, args ...string) (int, error) {
	var f flags

	fs := flag.NewFlagSet("mv", flag.ContinueOnError)
	fs.SetOutput(c.Stderr)

	fs.BoolVar(&f.update, "u", false, "move only when the SOURCE file is newer than the destination file or when the destination file is missing")
	fs.BoolVar(&f.noClobber, "n", false, "do not overwrite an existing file")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: mv [ARGS] source target [ARGS] source ... directory\n\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(unixflag.ArgsToGoArgs(args)); err != nil {
		return 1, err
	}

	if fs.NArg() < 2 {
		fs.Usage()
		return 1, fmt.Errorf("insufficient arguments")
	}

	if err := c.move(fs.Args(), f.update, f.noClobber); err != nil {
		return 1, err
	}

	return 0, nil
}
