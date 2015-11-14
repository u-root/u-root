// Copyright 2013 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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

func usage() {
	fmt.Fprintln(os.Stderr, "Usage:", cmd)
	flag.PrintDefaults()
	os.Exit(1)
}

func init() {
	flag.BoolVar(&flags.i, "i", false, "Interactive mode.")
	flag.BoolVar(&flags.v, "v", false, "Verbose mode.")
	flag.BoolVar(&flags.r, "R", false, "Remove file hierarchies")
	flag.BoolVar(&flags.r, "r", false, "Equivalent to -R.")
	flag.Parse()
	flag.Usage = usage
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
			fmt.Printf("%v: remove '%v'? ", cmd, file)
			input.Scan()
			if input.Text() != "y" {
				continue
			}
		}

		if flags.v {
			toRemove := file
			if file[0] != '/' {
				toRemove = path.Join(workingPath, file)
			}
			fmt.Printf("Deleting: %v\n", toRemove)
		}

		if err := f(file); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if flag.NArg() < 1 {
		usage()
	}

	if err := rm(flag.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
