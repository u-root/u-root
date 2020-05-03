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
	"text/tabwriter"

	flag "github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/ls"
)

var (
	all     = flag.BoolP("all", "a", false, "show hidden files")
	human   = flag.BoolP("human-readable", "h", false, "human readable sizes")
	long    = flag.BoolP("long", "l", false, "long form")
	quoted  = flag.BoolP("quote-name", "Q", false, "quoted")
	recurse = flag.BoolP("recursive", "R", false, "equivalent to findutil's find")
)

func listName(stringer ls.Stringer, d string, w io.Writer, prefix bool) error {
	return filepath.Walk(d, func(path string, osfi os.FileInfo, err error) error {
		// Soft error. Useful when a permissions are insufficient to
		// stat one of the files.
		if err != nil {
			log.Printf("%s: %v\n", path, err)
			return nil
		}

		fi := ls.FromOSFileInfo(path, osfi)

		if *recurse {
			// Mimic find command
			fi.Name = path
		} else if path == d {
			// Starting directory is a dot when non-recursive
			if osfi.IsDir() {
				fi.Name = "."
				if prefix {
					fmt.Printf("%q\n", d)
				}
			}
		}

		// Hide .files unless -a was given
		if *all || fi.Name[0] != '.' {
			// Print the file in the proper format.
			fmt.Fprintln(w, stringer.FileString(fi))
		}

		// Skip directories when non-recursive.
		if path != d && fi.Mode.IsDir() && !*recurse {
			return filepath.SkipDir
		}
		return nil
	})
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
