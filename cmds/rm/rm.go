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
)

var (
	interactive  = flag.Bool("i", false, "Interactive mode.")
	verbose      = flag.Bool("v", false, "Verbose mode.")
	hierarchies  = flag.Bool("R", false, "Remove file hierarchies")
	hierarchiesr = flag.Bool("r", false, "Equivalent to -R.")
	cmd          = "rm [-Rrvi] file..."
)

func init() {
	defUsage := flag.Usage
	flag.Usage = func() {
		os.Args[0] = cmd
		defUsage()
	}
}

func rm(files []string) error {
	f := os.Remove
	//fmt.Printf("\n R: %t \n r: %t \n", *hierarchies, *hierarchiesr)
	if *hierarchies || *hierarchiesr {
		//fmt.Printf("changing value of the function")
		f = os.RemoveAll
	}
	workingPath, err := os.Getwd()
	if err != nil {
		return err
	}

	input := bufio.NewScanner(os.Stdin)
	for _, file := range files {
		if *interactive {
			fmt.Printf("rm: remove '%v'? ", file)
			input.Scan()
			if input.Text()[0] != 'y' {
				continue
			}
		}

		if err := f(file); err != nil {
			return err
		}

		if *verbose {
			toRemove := file
			if !path.IsAbs(file) {
				toRemove = path.Join(workingPath, file)
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
	}

	if err := rm(flag.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
