// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// ls prints the contents of a directory.
//
// Synopsis:
//     ls [OPTIONS] [DIRS]...
//
// Options:
//     -l: long form
//     -Q: quoted
//     -R: equivalent to findutil's find
//     -F: append indicator (one of */=>@|) to entries
//
// Bugs:
//     With the `-R` flag, directories are only ever printed once.
package main

import (
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

func listName(stringer ls.Stringer, d string, w io.Writer, prefix bool) error {
	type file struct {
		path string
		osfi os.FileInfo
		lsfi ls.FileInfo
	}

	var files []file

	filepath.Walk(d, func(path string, osfi os.FileInfo, err error) error {
		// Soft error. Useful when a permissions are insufficient to
		// stat one of the files.
		if err != nil {
			log.Printf("%s: %v\n", path, err)
			return nil
		}

		fi := ls.FromOSFileInfo(path, osfi)

		if !*recurse && path == d && *directory {
			files = append(files, file{
				path: path,
				osfi: osfi,
				lsfi: fi,
			})
			return filepath.SkipDir
		}

		files = append(files, file{
			path: path,
			osfi: osfi,
			lsfi: fi,
		})

		if path != d && fi.Mode.IsDir() && !*recurse {
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

		// Hide .files unless -a was given
		if *all || f.lsfi.Name[0] != '.' {
			// Print the file in the proper format.
			if *classify {
				f.lsfi.Name = f.lsfi.Name + indicator(f.lsfi)
			}
			fmt.Fprintln(w, stringer.FileString(f.lsfi))
		}
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

func main() {
	flag.Parse()

	// Write output in tabular form.
	w := &tabwriter.Writer{}
	w.Init(os.Stdout, 0, 0, 1, ' ', 0)
	defer w.Flush()

	var s ls.Stringer = ls.NameStringer{}
	if *quoted {
		s = ls.QuotedStringer{}
	}
	if *long {
		s = ls.LongStringer{Human: *human, Name: s}
	}

	// Array of names to list.
	names := flag.Args()
	if len(names) == 0 {
		names = []string{"."}
	}

	// Is a name a directory? If so, list it in its own section.
	prefix := len(names) > 1
	for _, d := range names {
		if err := listName(s, d, w, prefix); err != nil {
			log.Printf("error while listing %#v: %v", d, err)
		}
		w.Flush()
	}
}
