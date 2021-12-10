// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Delete files.
//
// Synopsis:
//     rm [-Rrvif] FILE...
//
// Options:
//     -i: interactive mode
//     -v: verbose mode
//     -R: remove file hierarchies
//     -r: equivalent to -R
//     -f: ignore nonexistent files and never prompt
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"
)

var (
	flags struct {
		r bool
		v bool
		i bool
		f bool
	}
	cmd = "rm [-Rrvif] file..."
)

func init() {
	defUsage := flag.Usage
	flag.Usage = func() {
		os.Args[0] = cmd
		defUsage()
	}
	flag.BoolVarP(&flags.i, "interactive", "i", false, "Interactive mode.")
	flag.BoolVarP(&flags.v, "verbose", "v", false, "Verbose mode.")
	flag.BoolVarP(&flags.r, "recursive", "r", false, "remove hierarchies")
	flag.BoolVarP(&flags.r, "RECURSIVE", "R", false, "remove hierarchies")
	flag.BoolVarP(&flags.f, "force", "f", false, "Ignore nonexistent files and never prompt")
}

func rm(stdin io.Reader, files []string) error {
	f := os.Remove
	if flags.r {
		f = os.RemoveAll
	}

	if flags.f {
		flags.i = false
	}

	workingPath, err := os.Getwd()
	if err != nil {
		return err
	}

	input := bufio.NewReader(stdin)
	for _, file := range files {
		if flags.i {
			fmt.Printf("rm: remove '%v'? ", file)
			answer, err := input.ReadString('\n')
			if err != nil || strings.ToLower(answer)[0] != 'y' {
				continue
			}
		}

		if err := f(file); err != nil {
			if flags.f && os.IsNotExist(err) {
				continue
			}
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

	if err := rm(os.Stdin, flag.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
