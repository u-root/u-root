// Copyright 2013 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"flag"
)

var (
	recursive = flag.Bool("R", false, "Remove file hierarchies.")
	recursive_too = flag.Bool("r", false, "Equivalent to -R.")
	verbose = flag.Bool("v", false, "Verbose mode.")
	cmd = struct { name, flags string } {
		"rm",
		"[-Rrv] file...",
	}
)

// rm function 
func rm(args []string, do_recursive bool, verbose bool) error {
	f := os.Remove
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

func usage() {
	fmt.Fprintln(os.Stderr, "Usage:", cmd.name, cmd.flags)
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	flag.Usage = usage
	flag.Parse()


	if flag.NArg() < 1 {
		usage()
	}

	rm(os.Args[:1], *recursive || *recursive_too, *verbose)
}
