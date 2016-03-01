// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Cat reads each file from its arguments in sequence and writes it on the standard output.
*/

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
	f = os.Mkdir
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
