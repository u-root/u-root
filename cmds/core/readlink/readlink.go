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
//     -v: verbose
//
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const cmd = "readlink [-fv] FILE"

var (
	follow  = flag.Bool("f", false, "follow recursively")
	verbose = flag.Bool("v", false, "report error messages")
)

func init() {
	defUsage := flag.Usage
	flag.Usage = func() {
		os.Args[0] = cmd
		defUsage()
	}
	flag.Parse()
}

func readLink(file string) error {
	path, err := os.Readlink(file)
	if err != nil {
		return err
	}

	if *follow {
		path, err = filepath.EvalSymlinks(file)
	}

	fmt.Printf("%s\n", path)
	return err
}

func main() {
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
