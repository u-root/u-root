// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Delete files.
//
// Synopsis:
//
//	rm [-Rrvif] FILE...
//
// Options:
//
//	-i: interactive mode
//	-v: verbose mode
//	-R: remove file hierarchies
//	-r: equivalent to -R
//	-f: ignore nonexistent files and never prompt
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

var errUsage = errors.New("usage: rm [-Rrvif] file")

type flags struct {
	interactive bool
	verbose     bool
	recursive   bool
	force       bool
}

func rm(stdin io.Reader, stdout, stderr io.Writer, args ...string) error {
	var f flags
	fs := flag.NewFlagSet("rm", flag.ContinueOnError)
	fs.SetOutput(stderr)

	fs.BoolVar(&f.interactive, "i", false, "Interactive mode.")
	fs.BoolVar(&f.verbose, "v", false, "Verbose mode.")
	fs.BoolVar(&f.recursive, "r", false, "equivalent to -R")
	fs.BoolVar(&f.recursive, "R", false, "Recursive, remove hierarchies")
	fs.BoolVar(&f.force, "f", false, "Force, ignore nonexistent files and never prompt")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "%s...\n", errUsage)
		fs.PrintDefaults()
	}

	if err := fs.Parse(unixflag.ArgsToGoArgs(args)); err != nil {
		return err
	}

	files := fs.Args()

	if len(files) < 1 {
		return errUsage
	}

	rmf := os.Remove
	if f.recursive {
		rmf = os.RemoveAll
	}

	if f.force {
		f.interactive = false
	}

	input := bufio.NewReader(stdin)
	for _, file := range files {
		if f.interactive {
			fmt.Fprintf(stdout, "rm: remove '%v'? ", file)
			answer, err := input.ReadString('\n')
			if err != nil || (answer[0] != 'y' && answer[0] != 'Y') {
				continue
			}
		}

		if err := rmf(file); err != nil {
			if f.force && errors.Is(err, os.ErrNotExist) {
				continue
			}
			return err
		}

		if f.verbose {
			toRemove := file
			if !filepath.IsAbs(file) {
				workingPath, err := os.Getwd()
				if err != nil {
					return err
				}
				toRemove = filepath.Join(workingPath, file)
			}
			fmt.Fprintf(stdout, "removed '%v'\n", toRemove)
		}
	}
	return nil
}

func main() {
	if err := rm(os.Stdin, os.Stdout, os.Stderr, os.Args[1:]...); err != nil {
		log.Fatal(err)
	}
}
