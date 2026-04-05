// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cp implements routines to copy files.
//
// CopyTree in particular copies entire trees of files.
//
// Only directories, symlinks, and regular files are currently supported.
package cp

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
}

// Default are the default options. Default follows symlinks.
var Default = Options{}

// NoFollowSymlinks is the default options with following symlinks turned off.
var NoFollowSymlinks = Options{
	NoFollowSymlinks: true,
}

func (o Options) stat(path string) (os.FileInfo, error) {
	if o.NoFollowSymlinks {
		return os.Lstat(path)
	}
	return os.Stat(path)
}

// Copy copies a file at src to dst.
func (o Options) Copy(src, dst string) error {
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

// Copy src file to dst file using Default's config.
func Copy(src, dst string) error {
	return Default.Copy(src, dst)
}

// CopyTree recursively copies all files in the src tree to dst using Default's
// config.
func CopyTree(src, dst string) error {
	return Default.CopyTree(src, dst)
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
