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

func init() {
	flag.Parse()
	flag.BoolVar(&flags.i, "i", false, "Interactive mode.")
	flag.BoolVar(&flags.v, "v", false, "Verbose mode.")
	flag.BoolVar(&flags.r, "R", false, "Remove file hierarchies")
	flag.BoolVar(&flags.r, "r", false, "Equivalent to -R.")
}

func rm(files []string, recursive bool, verbose bool, interactive bool) error {
	f := os.Remove
	if recursive {
		f = os.RemoveAll
	}

	workingPath, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}

	// loop for remove files and folders
	for _, file := range files {
		// -i ask for remove file or no
		if interactive {
			fmt.Printf("%v: remove '%v'?: ", cmd.name, file)
			input := bufio.NewScanner(os.Stdin)
			input.Scan()
			if input.Text() != "y" {
				continue
			}
		}

		// -v print the deleting named path/files
		if verbose {
			toRemove := path.Join(workingPath, file)
			fmt.Printf("Deleting: %v\n", toRemove)
		}

		// try remove the file and return err if exists
		if err := f(file); err != nil {
			fmt.Fprintf(os.Stderr, "%v: %v\n", file, err)
			return err
		}
	}
	return nil
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage:", cmd)
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	flag.Usage = usage

	if flag.NArg() < 1 {
		usage()
	}

	if err := rm(flag.Args(), flags.r, flags.v, flags.i); err != nil {
		os.Exit(1)
	}
}
