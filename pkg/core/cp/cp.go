// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cp implements the cp core utility.
package cp

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/u-root/u-root/pkg/core"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

// ErrSkip can be returned by PreCallback to skip a file.
var ErrSkip = errors.New("skip")

// Options are configuration options for how copying files should behave.
type Options struct {
	// If NoFollowSymlinks is set, Copy copies the symlink itself rather
	// than following the symlink and copying the file it points to.
	NoFollowSymlinks bool

	// PreCallback is called on each file to be copied before it is copied
	// if specified.
	//
	// If PreCallback returns ErrSkip, the file is skipped and Copy returns
	// nil.
	//
	// If PreCallback returns another non-nil error, the file is not copied
	// and Copy returns the error.
	PreCallback func(src, dst string, srcfi os.FileInfo) error

	// PostCallback is called on each file after it is copied if specified.
	PostCallback func(src, dst string)

	// WorkingDir is the working directory for relative path resolution.
	WorkingDir string
}

// Default are the default options. Default follows symlinks.
var Default = Options{}

// NoFollowSymlinks is the default options with following symlinks turned off.
var NoFollowSymlinks = Options{
	NoFollowSymlinks: true,
}

// command implements the cp core utility.
type command struct {
	core.Base
}

// New creates a new cp command.
func New() core.Command {
	c := &command{}
	c.Init()
	return c
}

type flags struct {
	recursive        bool
	ask              bool
	force            bool
	verbose          bool
	noFollowSymlinks bool
}

func (o Options) stat(path string) (os.FileInfo, error) {
	if o.NoFollowSymlinks {
		return os.Lstat(path)
	}
	return os.Stat(path)
}

// resolvePath resolves a path relative to the working directory.
func (o Options) resolvePath(path string) string {
	if filepath.IsAbs(path) || o.WorkingDir == "" {
		return path
	}
	return filepath.Join(o.WorkingDir, path)
}

// Copy copies a file at src to dst.
func (o Options) Copy(src, dst string) error {
	src = o.resolvePath(src)
	dst = o.resolvePath(dst)

	srcInfo, err := o.stat(src)
	if err != nil {
		return err
	}

	if o.PreCallback != nil {
		if err := o.PreCallback(src, dst, srcInfo); err == ErrSkip {
			return nil
		} else if err != nil {
			return err
		}
	}
	if err := copyFile(src, dst, srcInfo); err != nil {
		return err
	}
	if o.PostCallback != nil {
		o.PostCallback(src, dst)
	}
	return nil
}

// CopyTree recursively copies all files in the src tree to dst.
func (o Options) CopyTree(src, dst string) error {
	src = o.resolvePath(src)
	dst = o.resolvePath(dst)

	return filepath.Walk(src, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		return o.Copy(path, filepath.Join(dst, rel))
	})
}

func copyFile(src, dst string, srcInfo os.FileInfo) error {
	m := srcInfo.Mode()
	switch {
	case m.IsDir():
		return os.MkdirAll(dst, srcInfo.Mode().Perm())

	case m.IsRegular():
		return copyRegularFile(src, dst, srcInfo)

	case m&os.ModeSymlink == os.ModeSymlink:
		// Yeah, this may not make any sense logically. But this is how
		// cp does it.
		target, err := os.Readlink(src)
		if err != nil {
			return err
		}
		return os.Symlink(target, dst)

	default:
		return &os.PathError{
			Op:   "copy",
			Path: src,
			Err:  fmt.Errorf("unsupported file mode %s", m),
		}
	}
}

func copyRegularFile(src, dst string, srcfi os.FileInfo) error {
	srcf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcf.Close()

	dstf, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, srcfi.Mode().Perm())
	if err != nil {
		return err
	}
	defer dstf.Close()

	_, err = io.Copy(dstf, srcf)
	return err
}

// promptOverwrite ask if the user wants overwrite file
func (c *command) promptOverwrite(dst string) (bool, error) {
	fmt.Fprintf(c.Stderr, "cp: overwrite %q? ", dst)
	reader := bufio.NewReader(c.Stdin)
	answer, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	if strings.ToLower(answer)[0] != 'y' {
		return false, nil
	}

	return true, nil
}

func (c *command) setupPreCallback(recursive, ask, force bool) func(string, string, os.FileInfo) error {
	return func(src, dst string, srcfi os.FileInfo) error {
		// check if src is dir
		if !recursive && srcfi.IsDir() {
			fmt.Fprintf(c.Stderr, "cp: -r not specified, omitting directory %s\n", src)
			return ErrSkip
		}

		dstfi, err := os.Stat(dst)
		if err != nil && !os.IsNotExist(err) {
			fmt.Fprintf(c.Stderr, "cp: %q: can't handle error %v\n", dst, err)
			return ErrSkip
		} else if err != nil {
			// dst does not exist.
			return nil
		}

		// dst does exist.

		if os.SameFile(srcfi, dstfi) {
			fmt.Fprintf(c.Stderr, "cp: %q and %q are the same file\n", src, dst)
			return ErrSkip
		}
		if ask && !force {
			overwrite, err := c.promptOverwrite(dst)
			if err != nil {
				return err
			}
			if !overwrite {
				return ErrSkip
			}
		}
		return nil
	}
}

func (c *command) setupPostCallback(verbose bool) func(src, dst string) {
	return func(src, dst string) {
		if verbose {
			fmt.Fprintf(c.Stdout, "%q -> %q\n", src, dst)
		}
	}
}

// Run executes the command with a `context.Background()`.
func (c *command) Run(args ...string) error {
	return c.RunContext(context.Background(), args...)
}

// Run executes the command.
func (c *command) RunContext(ctx context.Context, args ...string) error {
	var f flags

	fs := flag.NewFlagSet("cp", flag.ContinueOnError)
	fs.SetOutput(c.Stderr)

	fs.BoolVar(&f.recursive, "RECURSIVE", false, "copy file hierarchies")
	fs.BoolVar(&f.recursive, "R", false, "copy file hierarchies (shorthand)")

	fs.BoolVar(&f.recursive, "recursive", false, "alias to -R recursive mode")
	fs.BoolVar(&f.recursive, "r", false, "alias to -R recursive mode (shorthand)")

	fs.BoolVar(&f.ask, "interactive", false, "prompt about overwriting file")
	fs.BoolVar(&f.ask, "i", false, "prompt about overwriting file (shorthand)")

	fs.BoolVar(&f.force, "force", false, "force overwrite files")
	fs.BoolVar(&f.force, "f", false, "force overwrite files (shorthand)")

	fs.BoolVar(&f.verbose, "verbose", false, "verbose copy mode")
	fs.BoolVar(&f.verbose, "v", false, "verbose copy mode (shorthand)")

	fs.BoolVar(&f.noFollowSymlinks, "no-dereference", false, "don't follow symlinks")
	fs.BoolVar(&f.noFollowSymlinks, "P", false, "don't follow symlinks (shorthand)")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: cp [-RrifvP] file[s] ... dest\n\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(unixflag.ArgsToGoArgs(args)); err != nil {
		return err
	}

	if fs.NArg() < 2 {
		fs.Usage()
		return fmt.Errorf("insufficient arguments")
	}

	todir := false
	from, to := fs.Args()[:fs.NArg()-1], fs.Args()[fs.NArg()-1]

	toStat, err := os.Stat(to)
	if err == nil {
		todir = toStat.IsDir()
	}
	if fs.NArg() > 2 && !todir {
		return fmt.Errorf("target %q is not a directory", to)
	}

	opts := Options{
		NoFollowSymlinks: f.noFollowSymlinks,
		WorkingDir:       c.WorkingDir,
		PreCallback:      c.setupPreCallback(f.recursive, f.ask, f.force),
		PostCallback:     c.setupPostCallback(f.verbose),
	}

	var lastErr error
	for _, file := range from {
		dst := to
		if todir {
			dst = filepath.Join(dst, filepath.Base(file))
		}
		if f.recursive {
			lastErr = opts.CopyTree(file, dst)
		} else {
			lastErr = opts.Copy(file, dst)
		}
	}

	if lastErr != nil {
		return lastErr
	}
	return nil
}

// Copy src file to dst file using Default's config.
func Copy(src, dst string) error {
	return Default.Copy(src, dst)
}

// CopyTree recursively copies all files in the src tree to dst using Default's
// config.
func CopyTree(src, dst string) error {
	return Default.CopyTree(src, dst)
}
