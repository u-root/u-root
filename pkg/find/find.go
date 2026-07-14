// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package find searches for files in a directory hierarchy recursively.
//
// find can filter out files by file names, paths, and modes.
package find

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/u-root/u-root/pkg/ls"
)

var ErrInvalidRegexp = errors.New("invalid regular expression")

// File is a found file.
type File struct {
	// Name is the path relative to the root specified in WithRoot.
	Name string
	Err  error
}

// String implements a fmt.Stringer for File.
//
// String returns a string long-formatted like `ls` would format it.
func (f *File) String() string {
	s := ls.LongStringer{
		Human: true,
		Name:  ls.NameStringer{},
	}

	// If the walk already captured an error, show that.
	// This avoids an unneeded stat, which will also
	// take an error.
	if f.Err != nil {
		return fmt.Sprintf("%s: %v", f.Name, f.Err)
	}

	// If the stat fails, show that it failed.
	fi, err := os.Lstat(f.Name)
	if err != nil {
		return fmt.Sprintf("%s: %v", f.Name, err)
	}

	rec := ls.FromOSFileInfo(f.Name, fi)
	rec.Name = f.Name
	return s.FileString(rec)
}

// A Matcher returns true if the string/fs.DirEntry match some criteria.
// It will return an error if it had an error of some type.
type Matcher func(string, fs.DirEntry) (bool, error)

// Finder controls show files are walked.
type Finder struct {
	root         string
	matchers     []Matcher
	files        chan *File
	debug        func(string, ...any)
	ignoreErrors bool
}

func New() *Finder {
	return &Finder{
		debug: func(string, ...any) {},
		files: make(chan *File, 128),
	}
}

type Set func(*Finder)

// WithRoot sets a root path for the file finder. Only descendants of the root
// will be returned on the channel.
func WithRoot(rootPath string) Set {
	return func(f *Finder) {
		f.root = rootPath
	}
}

// WithoutError filters out files with errors from being sent on the channel.
func WithoutError() Set {
	return func(f *Finder) {
		f.ignoreErrors = true
	}
}

// Find finds files according to the settings and matchers given.
// It can be called with no opts, but requires a context.
// e.g. names := Find(context.Background())
// or
//
//	names := Find(ctx,
//	  WithRoot("/boot"),
//	  WithFilenameMatch("sda[0-9]"),
//	  WithDebugLog(log.Printf),
//	)
func Find(ctx context.Context, opt ...Set) <-chan *File {
	f := New()
	return RunFind(ctx, f, opt...)
}

// RunFind runs a Finder
func RunFind(ctx context.Context, f *Finder, opt ...Set) <-chan *File {

	for _, o := range opt {
		if o != nil {
			o(f)
		}
	}

	go func(f *Finder) {
		_ = filepath.WalkDir(f.root, func(n string, de fs.DirEntry, err error) error {
			f.debug("walk to %v, de %v, err %v", f, de, err)
			if err != nil && f.ignoreErrors {
				// Don't send file on channel if user doesn't want them.
				return nil
			}

			file := &File{
				Name: n,
				Err:  err,
			}

			if err == nil {
				// If it matches, then push its name into the result channel,
				// and keep looking.
				f.debug("check name %q", n)
				var ok bool
				var err error
				for _, m := range f.matchers {
					ok, err = m(n, de)
					if err != nil {
						f.debug("%s: err on matching: %v", n, err)
						return nil
					}
					if !ok {
						f.debug("%q:", n)
						return nil
					}

				}
				f.debug("Found: %s", n)
			}
			f.debug("pushing file")
			select {
			case <-ctx.Done():
				f.debug("ctx says done")
				return fmt.Errorf("should never be returned to user: stop walking")

			case f.files <- file:
				f.debug("pushed file")
				return nil
			}
		})
		f.debug("Done, close chan")
		close(f.files)
	}(f)

	f.debug("return from Find")
	return f.files
}
