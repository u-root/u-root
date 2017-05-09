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
	long    = flag.Bool("l", false, "long form")
	quoted  = flag.Bool("Q", false, "quoted")
	recurse = flag.Bool("R", false, "equivalent to findutil's find")
)

func stringer(fi fileInfo) fmt.Stringer {
	var s fmt.Stringer = fi
	if *quoted {
		s = quotedStringer{fi}
	}
	if *long {
		s = longStringer{fi, s}
	}
	return s
}

func listDir(d string, w io.Writer) error {
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
			fi.name = "."
		}

		// Print the file in the proper format.
		fmt.Fprintln(w, stringer(fi))

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

	// Array of directories to list.
	dirs := flag.Args()
	if len(dirs) == 0 {
		dirs = []string{"."}
	}

	// List each directory in its own section.
	for _, d := range dirs {
		if len(dirs) > 1 {
			fmt.Printf("%s:\n", d)
		}
		if err := listDir(d, w); err != nil {
			log.Printf("error while listing %#v: %v", d, err)
		}
		w.Flush()
	}
}
