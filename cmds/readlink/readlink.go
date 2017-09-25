// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// readlink display value of symbolic link file
//
// Synopsis:
//     readlink [OPTIONS] FILE
//
// Options:
//     -f: follow
//     -v: verbose
//
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var (
	follow  = flag.Bool("f", false, "follow recursively")
	verbose = flag.Bool("v", false, "report error messages")
)

func isLink(file string) error {
	var err error
	f, err := os.Lstat(file)

	if err != nil {
		return err
	}

	// If it is not a symbolic link return an error
	if f.Mode()&os.ModeSymlink == 0 {
		return fmt.Errorf("%s Invalid argument", file)
	}

	return err
}

func readFile(file string) error {
	var path string
	var err error

	err = isLink(file)
	if err != nil {
		return err
	}

	// Follow depth
	if *follow {
		path, err = filepath.EvalSymlinks(file)
	} else {
		path, err = os.Readlink(file)
	}

	fmt.Printf("%s\n", path)
	return err
}

func main() {
	flag.Parse()

	exitStatus := 0

	for _, file := range flag.Args() {
		err := readFile(file)

		if err != nil {
			if *verbose {
				fmt.Fprintf(os.Stderr, "readlink: %s\n", err.Error())
			}
			exitStatus = 1
		}
	}

	os.Exit(exitStatus)
}
