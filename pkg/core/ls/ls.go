// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ls implements the ls core utility.
package ls

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"text/tabwriter"

	"github.com/u-root/u-root/pkg/core"
	"github.com/u-root/u-root/pkg/ls"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

// command implements the ls command.
type command struct {
	core.Base
}

// New creates a new ls command.
func New() core.Command {
	c := &command{}
	c.Init()
	return c
}

type flags struct {
	all       bool
	human     bool
	directory bool
	long      bool
	quoted    bool
	recurse   bool
	classify  bool
	size      bool
	final     bool // Plan9/Windows specific
}

// file describes a file, its name, attributes, and the error
// accessing it, if any.
//
// Any such description must take into account the inherently
// racy nature of a file system. Can a file which exists in one
// instant vanish in another instant? Yes. Can we get into situations
// in which ls might never terminate? Yes (seen in HPC systems).
// If our consumer (ls) is slow enough, and our producer (thousands of
// compute nodes) is fast enough, an ls can take *hours*.
//
// Hence, file must include the path name (since a file can vanish,
// the stat might then fail, so using the fileinfo will not work)
// and must include an error (since the file may cease to exist).
// It is possible, for example, to do
// ls /a /b /c
// and between the time the command is typed, some or all of these
// files might vanish. Users wish to know of this situation:
// $ ls /a /b /tmp
// ls: /a: No such file or directory
// ls: /b: No such file or directory
// ls: /c: No such file or directory
// ls is more complex than it appears at first.
// TODO: do we really need BOTH osfi and lsfi?
// This may be required on non-unix systems like Plan 9 but it
// would be nice to make sure.
type file struct {
	path string
	osfi os.FileInfo
	lsfi ls.FileInfo
	err  error
}

func (c *command) listName(stringer ls.Stringer, d string, prefix bool, f flags) {
	var files []file
	resolvedPath := c.ResolvePath(d)

	filepath.Walk(resolvedPath, func(path string, osfi os.FileInfo, err error) error {
		file := file{
			path: path,
			osfi: osfi,
		}

		// error handling that matches standard ls is ... a real joy
		if osfi != nil && !errors.Is(err, os.ErrNotExist) {
			file.lsfi = ls.FromOSFileInfo(path, osfi)
			if err != nil && path == resolvedPath {
				file.err = err
			}
		} else {
			file.err = err
		}

		files = append(files, file)

		if err != nil {
			return filepath.SkipDir
		}

		if !f.recurse && path == resolvedPath && f.directory {
			return filepath.SkipDir
		}

		if path != resolvedPath && file.lsfi.Mode.IsDir() && !f.recurse {
			return filepath.SkipDir
		}

		return nil
	})

	if f.size {
		sort.SliceStable(files, func(i, j int) bool {
			return files[i].lsfi.Size > files[j].lsfi.Size
		})
	}

	for _, file := range files {
		if file.err != nil {
			c.printFile(stringer, file, f)
			continue
		}
		if f.recurse {
			// Mimic find command
			file.lsfi.Name = file.path
		} else if file.path == resolvedPath {
			if f.directory {
				fmt.Fprintln(c.Stdout, stringer.FileString(file.lsfi))
				continue
			}

			// Starting directory is a dot when non-recursive
			if file.osfi.IsDir() {
				file.lsfi.Name = "."
				if prefix {
					if f.quoted {
						fmt.Fprintf(c.Stdout, "%q:\n", d)
					} else {
						fmt.Fprintf(c.Stdout, "%v:\n", d)
					}
				}
			}
		}

		c.printFile(stringer, file, f)
	}
}

func indicator(fi ls.FileInfo) string {
	if fi.Mode.IsRegular() && fi.Mode&0o111 != 0 {
		return "*"
	}
	if fi.Mode&os.ModeDir != 0 {
		return "/"
	}
	if fi.Mode&os.ModeSymlink != 0 {
		return "@"
	}
	if fi.Mode&os.ModeSocket != 0 {
		return "="
	}
	if fi.Mode&os.ModeNamedPipe != 0 {
		return "|"
	}
	return ""
}

func (c *command) list(names []string, f flags) error {
	if len(names) == 0 {
		names = []string{"."}
	}
	// Write output in tabular form.
	tw := &tabwriter.Writer{}
	tw.Init(c.Stdout, 0, 0, 1, ' ', 0)
	stdout := c.Stdout
	c.Stdout = tw
	defer func() {
		tw.Flush()
		c.Stdout = stdout
	}()

	var s ls.Stringer = ls.NameStringer{}
	if f.quoted {
		s = ls.QuotedStringer{}
	}
	if f.long {
		s = ls.LongStringer{Human: f.human, Name: s}
	}
	// Is a name a directory? If so, list it in its own section.
	prefix := len(names) > 1
	for _, d := range names {
		c.listName(s, d, prefix, f)
		tw.Flush()
	}
	return nil
}

// Run executes the command with a `context.Background()`.
func (c *command) Run(args ...string) error {
	return c.RunContext(context.Background(), args...)
}

// Run executes the command.
func (c *command) RunContext(ctx context.Context, args ...string) error {
	var f flags

	fs := flag.NewFlagSet("ls", flag.ContinueOnError)
	fs.SetOutput(c.Stderr)

	fs.BoolVar(&f.all, "a", false, "show hidden files")
	fs.BoolVar(&f.human, "h", false, "human readable sizes")
	fs.BoolVar(&f.directory, "d", false, "list directories but not their contents")
	fs.BoolVar(&f.long, "l", false, "long form")
	fs.BoolVar(&f.quoted, "Q", false, "quoted")
	fs.BoolVar(&f.recurse, "R", false, "equivalent to findutil's find")
	fs.BoolVar(&f.classify, "F", false, "append indicator (, one of */=>@|) to entries")
	fs.BoolVar(&f.size, "S", false, "sort by size")

	// OS-specific flags
	c.addOSSpecificFlags(fs, &f)

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: ls [OPTIONS] [DIRS]...\n\n")
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(unixflag.ArgsToGoArgs(args)); err != nil {
		return err
	}

	if err := c.list(fs.Args(), f); err != nil {
		return err
	}

	return nil
}

// TestIndicator exposes the indicator function for testing.
func (c *command) TestIndicator(fi ls.FileInfo) string {
	return indicator(fi)
}

// TestIndicator exposes the indicator function for external testing.
func TestIndicator(fi ls.FileInfo) string {
	return indicator(fi)
}
