// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Ls prints the contents of a directory.
//
// Synopsis:
//     ls [OPTIONS] [DIRS]...
//
// Options:
//     -l: long form
//     -Q: quoted
//     -R: equivalent to findutil's find
//
// Bugs:
//     With the `-R` flag, directories are only ever printed once.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"text/tabwriter"
)

var (
	all     = flag.Bool("a", false, "show hidden files")
	human   = flag.Bool("h", false, "human readable sizes")
	long    = flag.Bool("l", false, "long form")
	quoted  = flag.Bool("Q", false, "quoted")
	recurse = flag.Bool("R", false, "equivalent to findutil's find")
)

func stringer(fi fileInfo) fmt.Stringer {
	var s fmt.Stringer = fi
	if *quoted {
		s = quotedStringer{fileInfo: fi}
	}
	if *long {
		s = longStringer{
			fileInfo: fi,
			comp:     s,
			human:    *human,
		}
	}
	return s
}

func listName(d string, w io.Writer, prefix bool) error {
	return filepath.Walk(d, func(path string, osfi os.FileInfo, err error) error {
		// Soft error. Useful when a permissions are insufficient to
		// stat one of the files.
		if err != nil {
			log.Printf("%s: %v\n", path, err)
			return nil
		}

		fi := extractImportantParts(path, osfi)

		if *recurse {
			// Mimic find command
			fi.name = path
		} else if path == d {
			// Starting directory is a dot when non-recursive
			if osfi.IsDir() {
				fi.name = "."
				if prefix {
					fmt.Printf("%q\n", d)
				}
			}
		}

		// Hide .files unless -a was given
		if *all || fi.name[0] != '.' {
			// Print the file in the proper format.
			fmt.Fprintln(w, stringer(fi))
		}

		// Skip directories when non-recursive.
		if path != d && fi.mode.IsDir() && !*recurse {
			return filepath.SkipDir
		}
		return nil
	})
}

func main() {
	flag.Parse()

	// Write output in tabular form.
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 0, 1, ' ', 0)
	defer w.Flush()

	// Array of names to list.
	names := flag.Args()
	if len(names) == 0 {
		names = []string{"."}
	}

	// Is a name a directory? If so, list it in its own section.
	prefix := len(names) > 1
	for _, d := range names {
		if err := listName(d, w, prefix); err != nil {
			log.Printf("error while listing %#v: %v", d, err)
		}
		w.Flush()
	}
}
