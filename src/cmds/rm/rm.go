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
	recursive       = flag.Bool("R", false, "Remove file hierarchies.")
	recursiveAlias = flag.Bool("r", false, "Equivalent to -R.")
	verbose         = flag.Bool("v", false, "Verbose mode.")
	interactive     = flag.Bool("i", false, "Interactive mode.")
	cmd             = struct{ name, flags string }{
		"rm",
		"[-Rrvi] file...",
	}
)

func rm(files []string, recursive bool, verbose bool, interactive bool) error {
	f := os.Remove
	if recursive {
		f = os.RemoveAll
	}
	workingPath, _ := os.Getwd()

	// loop for remove files and folders
	for _, file := range files {
		if interactive {
			fmt.Printf("%v: remove '%v'?: ", cmd.name, file)
			input := bufio.NewScanner(os.Stdin)
			input.Scan()
			if input.Text() != "y" {
				continue
			}
		}

		if verbose {
			toRemove := path.Join(workingPath, file)
			fmt.Printf("Deleting: %v\n", toRemove)
		}

		if err := f(file); err != nil {
			fmt.Fprintf(os.Stderr, "%v: %v\n", file, err)
			return err
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

	rm(flag.Args(), recursive, *verbose, *interactive)
}
