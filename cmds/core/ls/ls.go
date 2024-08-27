// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// ls prints the contents of a directory.
//
// Synopsis:
//
//	ls [OPTIONS] [DIRS]...
//
// Options:
//
//	-a: show hidden files
//	-h: show human-readable sizes
//	-d: show directories but not their contents
//	-F: append indicator (, one of */=>@|) to entries
//	-l: long form
//	-Q: quoted
//	-R: equivalent to findutil's find
//	-s: sort by size
//
// Bugs:
//
//	With the `-R` flag, directories are only ever printed once.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"text/tabwriter"

	"github.com/u-root/u-root/pkg/ls"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

type cmd struct {
	w         io.Writer
	all       bool
	human     bool
	directory bool
	long      bool
	quoted    bool
	recurse   bool
	classify  bool
	size      bool
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

func (c cmd) listName(stringer ls.Stringer, d string, prefix bool) error {
	var files []file

	filepath.Walk(d, func(path string, osfi os.FileInfo, err error) error {
		f := file{
			path: path,
			osfi: osfi,
		}

		// error handling that matches standard ls is ... a real joy
		if osfi != nil && !errors.Is(err, os.ErrNotExist) {
			f.lsfi = ls.FromOSFileInfo(path, osfi)
			if err != nil && path == d {
				f.err = err
			}
		} else {
			f.err = err
		}

		files = append(files, f)

		if err != nil {
			return filepath.SkipDir
		}

		if !c.recurse && path == d && c.directory {
			return filepath.SkipDir
		}

		if path != d && f.lsfi.Mode.IsDir() && !c.recurse {
			return filepath.SkipDir
		}

		return nil
	})

	if c.size {
		sort.SliceStable(files, func(i, j int) bool {
			return files[i].lsfi.Size > files[j].lsfi.Size
		})
	}

	for _, f := range files {
		if f.err != nil {
			c.printFile(stringer, f)
			continue
		}
		if c.recurse {
			// Mimic find command
			f.lsfi.Name = f.path
		} else if f.path == d {
			if c.directory {
				fmt.Fprintln(c.w, stringer.FileString(f.lsfi))
				continue
			}

			// Starting directory is a dot when non-recursive
			if f.osfi.IsDir() {
				f.lsfi.Name = "."
				if prefix {
					if c.quoted {
						fmt.Fprintf(c.w, "%q:\n", d)
					} else {
						fmt.Fprintf(c.w, "%v:\n", d)
					}
				}
			}
		}

		c.printFile(stringer, f)
	}

	return nil
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

func (c cmd) list(names []string) error {
	if len(names) == 0 {
		names = []string{"."}
	}
	// Write output in tabular form.
	tw := &tabwriter.Writer{}
	tw.Init(c.w, 0, 0, 1, ' ', 0)
	c.w = tw
	defer tw.Flush()

	var s ls.Stringer = ls.NameStringer{}
	if c.quoted {
		s = ls.QuotedStringer{}
	}
	if c.long {
		s = ls.LongStringer{Human: c.human, Name: s}
	}
	// Is a name a directory? If so, list it in its own section.
	prefix := len(names) > 1
	for _, d := range names {
		if err := c.listName(s, d, prefix); err != nil {
			return fmt.Errorf("error while listing %q: %w", d, err)
		}
		tw.Flush()
	}
	return nil
}

// run runs a command. args are as from os.Args, i.e., args[0] is the command name.
func run(w io.Writer, args []string) error {
	var c cmd
	f := flag.NewFlagSet(args[0], flag.ExitOnError)
	f.BoolVar(&c.all, "a", false, "show hidden files")
	f.BoolVar(&c.human, "h", false, "human readable sizes")
	f.BoolVar(&c.directory, "d", false, "list directories but not their contents")
	f.BoolVar(&c.long, "l", false, "long form")
	f.BoolVar(&c.quoted, "Q", false, "quoted")
	f.BoolVar(&c.recurse, "R", false, "equivalent to findutil's find")
	f.BoolVar(&c.classify, "F", false, "append indicator (, one of */=>@|) to entries")
	f.BoolVar(&c.size, "S", false, "sort by size")
	c.w = w
	f.Parse(unixflag.ArgsToGoArgs(args[1:]))
	return c.list(f.Args())
}

func main() {
	if err := run(os.Stdout, os.Args); err != nil {
		log.Fatal(err)
	}
}
