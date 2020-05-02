// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// readlink display value of symbolic link file.
//
// Synopsis:
//     readlink [OPTIONS] FILE
//
// Options:
//     -f: follow
//     -n: nonewline
//     -v: verbose
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const cmd = "readlink [-fnv] FILE"

var (
	delimiter = "\n"
	follow    = flag.Bool("f", false, "follow recursively")
	nonewline = flag.Bool("n", false, "do not output trailing newline")
	verbose   = flag.Bool("v", false, "report error messages")
)

func init() {
	defUsage := flag.Usage
	flag.Usage = func() {
		os.Args[0] = cmd
		defUsage()
	}
}

func readLink(file string) error {
	path, err := os.Readlink(file)
	if err != nil {
		return err
	}

	if *follow {
		path, err = filepath.EvalSymlinks(file)
	}

	if *nonewline {
		delimiter = ""
	}

	fmt.Printf("%s%s", path, delimiter)
	return err
}

func main() {
	flag.Parse()

	var exitStatus int

	for _, file := range flag.Args() {
		if err := readLink(file); err != nil {
			if *verbose {
				fmt.Fprintf(os.Stderr, "%v\n", err)
			}
			exitStatus = 1
		}
	}

	os.Exit(exitStatus)
}
