// Copyright 2013 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"flag"
)

var recursive = flag.Bool("R", false, "Remove file hierarchies.")
var recursive_too = flag.Bool("r", false, "Equivalent to -R.")
var verbose = flag.Bool("v", false, "Verbose mode.")

func usage (){
	fmt.Fprintf(os.Stderr, "usage: %s [-Rrv] file...\n", os.Args[0])
	os.Exit(1)
}

// rm function 
func rm(args []string, do_recursive bool, verbose bool) error {
	
	// recursive or no depend do_recursive value
	var f = os.Remove
	if do_recursive {
		f = os.RemoveAll
	}

	// looping for removing files and folders
	for _,arg := range(args) {
		if arg == "-r" || arg == "-R" {
			continue
		}

		err := f(arg)
		if err != nil {
			fmt.Printf("%v: %v\n", arg, err)
			return err
		}

		if verbose {
			fmt.Printf("Deleting: %v\n", arg)
		}
	}
	return nil
}

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		usage()
	}

	rm(os.Args[:1], *recursive || *recursive_too, *verbose)
}
