// Copyright 2013-2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// readlink display value of symbolic link file.
//
// Synopsis:
//
//	readlink [OPTIONS] [FILE...]
//
// Options:
//
//	-f: follow
//	-n: nonewline
//	-v: verbose
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const cmd = "readlink [-fnv] FILE"

var (
	follow    = flag.Bool("f", false, "follow recursively")
	noNewLine = flag.Bool("n", false, "do not output trailing newline")
	verbose   = flag.Bool("v", false, "report error messages")
)

func init() {
	defUsage := flag.Usage
	flag.Usage = func() {
		os.Args[0] = cmd
		defUsage()
	}
}

func readLink(stdout io.Writer, file string) error {
	path, err := os.Readlink(file)
	if err != nil {
		return err
	}

	if *follow {
		path, err = filepath.EvalSymlinks(file)
	}

	delimiter := "\n"
	if *noNewLine {
		delimiter = ""
	}

	fmt.Fprintf(stdout, "%s%s", path, delimiter)
	return err
}

func run(stdout io.Writer, stderr io.Writer, args []string) error {
	if len(args) == 0 {
		if *verbose {
			fmt.Fprintf(stderr, "missing operand")
		}
		return fmt.Errorf("missing operand")
	}

	var runErr error
	for _, file := range args {
		err := readLink(stdout, file)
		if err != nil {
			if *verbose {
				fmt.Fprintf(stderr, "%v\n", err)
			}
			runErr = err
		}
	}

	return runErr
}

func main() {
	flag.Parse()
	if err := run(os.Stdout, os.Stderr, flag.Args()); err != nil {
		os.Exit(1)
	}
}
