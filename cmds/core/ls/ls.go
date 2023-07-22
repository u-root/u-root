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
//	-l: long form
//	-Q: quoted
//	-R: equivalent to findutil's find
//	-F: append indicator (one of */=>@|) to entries
//
// Bugs:
//
//	With the `-R` flag, directories are only ever printed once.
package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"text/tabwriter"

	flag "github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/ls"
)

var (
	all       = flag.BoolP("all", "a", false, "show hidden files")
	human     = flag.BoolP("human-readable", "h", false, "human readable sizes")
	directory = flag.BoolP("directory", "d", false, "list directories but not their contents")
	long      = flag.BoolP("long", "l", false, "long form")
	quoted    = flag.BoolP("quote-name", "Q", false, "quoted")
	recurse   = flag.BoolP("recursive", "R", false, "equivalent to findutil's find")
	classify  = flag.BoolP("classify", "F", false, "append indicator (one of */=>@|) to entries")
	size      = flag.BoolP("size", "S", false, "sort by size")
)

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

func listName(stringer ls.Stringer, d string, w io.Writer, prefix bool) error {
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

		if !*recurse && path == d && *directory {
			return filepath.SkipDir
		}

		if path != d && f.lsfi.Mode.IsDir() && !*recurse {
			return filepath.SkipDir
		}

		return nil
	})

	if *size {
		sort.SliceStable(files, func(i, j int) bool {
			return files[i].lsfi.Size > files[j].lsfi.Size
		})
	}

	for _, f := range files {
		if f.err != nil {
			printFile(w, stringer, f)
			continue
		}
		if *recurse {
			// Mimic find command
			f.lsfi.Name = f.path
		} else if f.path == d {
			if *directory {
				fmt.Fprintln(w, stringer.FileString(f.lsfi))
				continue
			}

			// Starting directory is a dot when non-recursive
			if f.osfi.IsDir() {
				f.lsfi.Name = "."
				if prefix {
					if *quoted {
						fmt.Fprintf(w, "%q:\n", d)
					} else {
						fmt.Fprintf(w, "%v:\n", d)
					}
				}
			}
		}

		printFile(w, stringer, f)
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

func list(w io.Writer, names []string) error {
	if len(names) == 0 {
		names = []string{"."}
	}
	// Write output in tabular form.
	tw := &tabwriter.Writer{}
	tw.Init(w, 0, 0, 1, ' ', 0)
	defer tw.Flush()

	var s ls.Stringer = ls.NameStringer{}
	if *quoted {
		s = ls.QuotedStringer{}
	}
	if *long {
		s = ls.LongStringer{Human: *human, Name: s}
	}
	// Is a name a directory? If so, list it in its own section.
	prefix := len(names) > 1
	for _, d := range names {
		if err := listName(s, d, tw, prefix); err != nil {
			return fmt.Errorf("error while listing %q: %w", d, err)
		}
		tw.Flush()
	}
	return nil
}

func main() {
	flag.Parse()
	if err := list(os.Stdout, flag.Args()); err != nil {
		log.Fatal(err)
	}
}
