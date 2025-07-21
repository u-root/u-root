// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package find implements the find core utility.
package find

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/u-root/u-root/pkg/core"
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
	debug      func(string, ...any)
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
func WithDebugLog(l func(string, ...any)) Set {
	return func(f *finder) {
		f.debug = l
	}
}

// Find finds files according to the settings and matchers given.
//
// e.g.
//
//	names := Find(ctx,
//	  WithRoot("/boot"),
//	  WithFilenameMatch("sda[0-9]"),
//	  WithDebugLog(log.Printf),
//	)
func Find(ctx context.Context, opt ...Set) <-chan *File {
	f := &finder{
		root:       "/",
		debug:      func(string, ...any) {},
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

// command implements the find core utility.
type command struct {
	core.Base
}

// New creates a new find command.
func New() core.Command {
	c := &command{}
	c.Init()
	return c
}

type flags struct {
	fileType string
	name     string
	perm     int
	long     bool
	debug    bool
}

// Run executes the command with a `context.Background()`.
func (c *command) Run(args ...string) error {
	return c.RunContext(context.Background(), args...)
}

// Run executes the command.
func (c *command) RunContext(ctx context.Context, args ...string) error {
	var f flags

	fs := flag.NewFlagSet("find", flag.ContinueOnError)
	fs.SetOutput(c.Stderr)

	fs.StringVar(&f.fileType, "type", "", "file type")
	fs.StringVar(&f.name, "name", "", "glob for name")
	fs.IntVar(&f.perm, "mode", -1, "permissions")
	fs.BoolVar(&f.long, "l", false, "long listing")
	fs.BoolVar(&f.debug, "d", false, "enable debugging in the find package")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: find [opts] starting-at-path\n\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(reorderArgs(args)); err != nil {
		return err
	}

	if fs.NArg() != 1 {
		fs.Usage()
		return fmt.Errorf("insufficient arguments")
	}

	root := c.ResolvePath(fs.Args()[0])

	fileTypes := map[string]os.FileMode{
		"f":         0,
		"file":      0,
		"d":         os.ModeDir,
		"directory": os.ModeDir,
		"s":         os.ModeSocket,
		"p":         os.ModeNamedPipe,
		"l":         os.ModeSymlink,
		"c":         os.ModeCharDevice | os.ModeDevice,
		"b":         os.ModeDevice,
	}

	var mask, mode os.FileMode
	if f.perm != -1 {
		mask = os.ModePerm
		mode = os.FileMode(f.perm)
	}
	if f.fileType != "" {
		intType, ok := fileTypes[f.fileType]
		if !ok {
			var keys []string
			for key := range fileTypes {
				keys = append(keys, key)
			}
			return fmt.Errorf("%v is not a valid file type\n valid types are %v", f.fileType, strings.Join(keys, ","))
		}
		mode |= intType
		mask |= os.ModeType
	}

	debugLog := func(string, ...any) {}
	if f.debug {
		debugLog = func(format string, args ...any) {
			fmt.Fprintf(c.Stderr, format+"\n", args...)
		}
	}

	names := Find(ctx,
		WithRoot(root),
		WithModeMatch(mode, mask),
		WithFilenameMatch(f.name),
		WithDebugLog(debugLog),
	)

	for l := range names {
		if l.Err != nil {
			fmt.Fprintf(c.Stderr, "%s: %v\n", l.Name, l.Err)
			continue
		}
		if f.long {
			fmt.Fprintf(c.Stdout, "%s\n", l)
			continue
		}
		fmt.Fprintf(c.Stdout, "%s\n", l.Name)
	}

	return nil
}

// reorderArgs reorders arguments so flags are moved to the front, which is the
// way the "flag" package can parse them.
func reorderArgs(args []string) []string {
	var (
		newArgs []string
		i       int
	)

	for i < len(args) {
		var (
			arg          = args[i]
			expectsValue = arg == "-name" || arg == "-type" || arg == "-mode"
			hasNext      = i+1 < len(args)
			isFlag       = strings.HasPrefix(arg, "-")
		)

		prepend := func(args ...string) {
			newArgs = append(args, newArgs...)
		}

		switch {
		case expectsValue && hasNext:
			next := args[i+1]
			prepend(arg, next)
			i += 2
		case isFlag:
			prepend(arg)
			i++
		default:
			newArgs = append(newArgs, arg)
			i++
		}
	}

	return newArgs
}
