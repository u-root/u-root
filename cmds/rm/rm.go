// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Delete files.
//
// Synopsis:
//     rm [-Rrvi] FILE...
//
// Options:
//     -i: interactive mode
//     -v: verbose mode
//     -R: remove file hierarchies
//     -r: equivalent to -R
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
)

var (
	flags struct {
		r bool
		v bool
		i bool
	}
	cmd = "rm [-Rrvi] file..."
)

func init() {
	flag.Usage = func(f func()) func() {
		return func() {
			os.Args[0] = cmd
			f()
		}
	}(flag.Usage)
	flag.BoolVar(&flags.i, "i", false, "Interactive mode.")
	flag.BoolVar(&flags.v, "v", false, "Verbose mode.")
	flag.BoolVar(&flags.r, "R", false, "Remove file hierarchies")
	flag.BoolVar(&flags.r, "r", false, "Equivalent to -R.")
	flag.Parse()
}

func rm(files []string) error {
	f := os.Remove
	if flags.r {
		f = os.RemoveAll
	}

	workingPath, err := os.Getwd()
	if err != nil {
		return err
	}

	input := bufio.NewScanner(os.Stdin)
	for _, file := range files {
		if flags.i {
			fmt.Printf("rm: remove '%v'? ", file)
			input.Scan()
			if input.Text()[0] != 'y' {
				continue
			}
		}

		if err := f(file); err != nil {
			return err
		}

		if flags.v {
			toRemove := file
			if !path.IsAbs(file) {
				toRemove = path.Join(workingPath, file)
			}
			fmt.Printf("removed '%v'\n", toRemove)
		}
	}
	return nil
}

func main() {
	if flag.NArg() < 1 {
		flag.Usage()
	}

	if err := rm(flag.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
