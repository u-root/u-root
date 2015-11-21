// Copyright 2013 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"path"
)

var (
	recursive_flag  = flag.Bool("R", false, "Remove file hierarchies.")
	recursive_alias = flag.Bool("r", false, "Equivalent to -R.")
	verbose         = flag.Bool("v", false, "Verbose mode.")
	cmd             = struct{ name, flags string }{
		"rm",
		"[-Rrv] file...",
	}
)

func rm(files []string, do_recursive bool, verbose bool) error {
	f := os.Remove
	if do_recursive {
		f = os.RemoveAll
	}
	working_path, _ := os.Getwd()

	// loop for remove files and folders
	for _, file := range files {
		err := f(file)
		if err != nil {
			fmt.Printf("%v: %v\n", file, err)
			return err
		}

		if verbose {
			deleted := path.Join(working_path, file)
			fmt.Printf("Deleting: %v\n", deleted)
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
	recursive := *recursive_flag || *recursive_alias
	flag.Parse()

	if flag.NArg() < 1 {
		usage()
	}

	rm(flag.Args(), recursive, *verbose)
}
