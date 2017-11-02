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
	"path/filepath"
	"strings"
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
	defUsage := flag.Usage
	flag.Usage = func() {
		os.Args[0] = cmd
		defUsage()
	}
	flag.BoolVar(&flags.i, "i", false, "Interactive mode.")
	flag.BoolVar(&flags.v, "v", false, "Verbose mode.")
	flag.BoolVar(&flags.r, "R", false, "Remove file hierarchies")
	flag.BoolVar(&flags.r, "r", false, "Equivalent to -R.")
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

	input := bufio.NewReader(os.Stdin)
	for _, file := range files {
		if flags.i {
			fmt.Printf("rm: remove '%v'? ", file)
			answer, err := input.ReadString('\n')
			if err != nil || strings.ToLower(answer)[0] != 'y' {
				continue
			}
		}

		if err := f(file); err != nil {
			return err
		}

		if flags.v {
			toRemove := file
			if !path.IsAbs(file) {
				toRemove = filepath.Join(workingPath, file)
			}
			fmt.Printf("removed '%v'\n", toRemove)
		}
	}
	return nil
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	if err := rm(flag.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
