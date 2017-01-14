// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Mkdir makes a new directory.
//
// Synopsis:
//     mkdir [-m mode] [-v] [-p] DIRECTORY...
//
// Options:
//     -m: make all needed directories in the path
//     -v: directory mode (ex: 666)
//     -p: print each directory as it is made
package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	mkall   = flag.Bool("p", false, "Make all needed directories in the path")
	mode    = flag.Int("m", 0666, "Directory mode")
	verbose = flag.Bool("v", false, "Print each directory as it is made")
	f       = os.Mkdir
)

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Printf("Usage: mkdir [-m mode] [-v] [-p] <directory> [more directories]\n")
		os.Exit(1)
	}
	if *mkall {
		f = os.MkdirAll
	}
	for _, name := range flag.Args() {
		if err := f(name, os.FileMode(*mode)); err != nil {
			fmt.Printf("%v: %v\n", name, err)
		} else {
			if *verbose {
				fmt.Printf("%v\n", name)
			}
		}
	}
}
