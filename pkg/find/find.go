// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package find searches for files in a directory hierarchy recursively.
//
// find can filter out files by file names, paths, and modes.
package find

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/u-root/u-root/pkg/ls"
)

// File is a found file.
type File struct {
	// Name is the path relative to the root specified in WithRoot.
	Name string

	os.FileInfo
	Err error
}

// String implements a fmt.Stringer for File.
//
// String returns a string long-formatted like `ls` would format it.
func (f *File) String() string {
	s := ls.LongStringer{
		Human: true,
		Name:  ls.NameStringer{},
	}
	rec := ls.FromOSFileInfo(f.Name, f.FileInfo)
	rec.Name = f.Name
	return s.FileString(rec)
}

type finder struct {
	root string

	// Pattern is used with Match.
	pattern string

	// Match is a pattern matching function.
	match      func(pattern string, name string) (bool, error)
	mode       os.FileMode
	modeMask   os.FileMode
	debug      func(string, ...interface{})
	files      chan *File
	sendErrors bool
}

type Set func(*finder)

// WithRoot sets a root path for the file finder. Only descendants of the root
// will be returned on the channel.
func WithRoot(rootPath string) Set {
	return func(f *finder) {
		f.root = rootPath
	}
}

// WithoutError filters out files with errors from being sent on the channel.
func WithoutError() Set {
	return func(f *finder) {
		f.sendErrors = false
	}
}

// WithPathMatch sets up a file path filter.
//
// The file path passed to match will be relative to the finder's root.
func WithPathMatch(pattern string, match func(pattern string, path string) (bool, error)) Set {
	return func(f *finder) {
		f.pattern = pattern
		f.match = match
	}
}

// WithBasenameMatch sets up a file base name filter.
func WithBasenameMatch(pattern string, match func(pattern string, name string) (bool, error)) Set {
	return WithPathMatch(pattern, func(patt string, path string) (bool, error) {
		return match(pattern, filepath.Base(path))
	})
}

// WithRegexPathMatch sets up a path filter using regex.
//
// The file path passed to regexp.Match will be relative to the finder's root.
func WithRegexPathMatch(pattern string) Set {
	return WithPathMatch(pattern, func(pattern, path string) (bool, error) {
		return regexp.Match(pattern, []byte(path))
	})
}

// WithFilenameMatch uses filepath.Match's shell file name matching to filter
// file base names.
func WithFilenameMatch(pattern string) Set {
	return WithBasenameMatch(pattern, filepath.Match)
}

// WithModeMatch ensures only files with fileMode & modeMask == mode are returned.
func WithModeMatch(mode, modeMask os.FileMode) Set {
	return func(f *finder) {
		f.mode = mode
		f.modeMask = modeMask
	}
}

// WithDebugLog logs messages to l.
func WithDebugLog(l func(string, ...interface{})) Set {
	return func(f *finder) {
		f.debug = l
	}
}

// Find finds files according to the settings and matchers given.
//
// e.g.
//
//   names := Find(ctx,
//     WithRoot("/boot"),
//     WithFilenameMatch("sda[0-9]"),
//     WithDebugLog(log.Printf),
//   )
func Find(ctx context.Context, opt ...Set) <-chan *File {
	f := &finder{
		root:       "/",
		debug:      func(string, ...interface{}) {},
		files:      make(chan *File, 128),
		match:      filepath.Match,
		sendErrors: true,
	}

	for _, o := range opt {
		if o != nil {
			o(f)
		}
	}

	go func(f *finder) {
		_ = filepath.Walk(f.root, func(n string, fi os.FileInfo, err error) error {
			if err != nil && !f.sendErrors {
				// Don't send file on channel if user doesn't want them.
				return nil
			}

			file := &File{
				Name:     n,
				FileInfo: fi,
				Err:      err,
			}
			if err == nil {
				// If it matches, then push its name into the result channel,
				// and keep looking.
				f.debug("check pattern %q against name %q", f.pattern, n)
				if f.pattern != "" {
					m, err := f.match(f.pattern, n)
					if err != nil {
						f.debug("%s: err on matching: %v", n, err)
						return nil
					}
					if !m {
						f.debug("%s: name does not match %q", n, f.pattern)
						return nil
					}
				}
				m := fi.Mode()
				f.debug("%s: file mode %v / want mode %s with mask %s", n, m, f.mode, f.modeMask)
				if masked := m & f.modeMask; masked != f.mode {
					f.debug("%s: mode %s (masked %s) does not match expected mode %s", n, m, masked, f.mode)
					return nil
				}
				f.debug("Found: %s", n)
			}
			select {
			case <-ctx.Done():
				return fmt.Errorf("should never be returned to user: stop walking")

			case f.files <- file:
				return nil
			}
		})
		close(f.files)
	}(f)

	return f.files
}
