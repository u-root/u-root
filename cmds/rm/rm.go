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
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// Used an array of flags in case more flags are added (see rm method signature)
/*
	interactive  = flag.Bool("i", false, "Interactive mode.")
	verbose      = flag.Bool("v", false, "Verbose mode.")
	hierarchies  = flag.Bool("r", false, "Remove file hierarchies")
	flagsRm = [interactive, verbose, hierarchies]
*/
//attempting to use struct
type rmFlags struct {
	recursive   bool
	verbose     bool
	interactive bool
}

/*var {
	cmd = "rm [-Rrvi] file..."
)*/
/*
func init() {
	defUsage := flag.Usage
	flag.Usage = func() {
		os.Args[0] = cmd
		defUsage()
	}
	flag.BoolVar(&flags.verbose, "v", false, "Verbose mode.")
	flag.BoolVar(&flags.recursive, "r", false, "Recursive mode.")
	flag.BoolVar(&flags.interactive, "i", false, "Interactive mode.")
}
*/
func rm(files []string, flags rmFlags) error {
	f := os.Remove
	if flags.recursive {
		f = os.RemoveAll
	}
	workingPath, err := os.Getwd()
	if err != nil {
		return err
	}
	input := bufio.NewScanner(os.Stdin)
	for _, file := range files {
		if flags.interactive {
			fmt.Printf("rm: remove '%v'? ", file)
			input.Scan()
			if input.Text()[0] != 'y' {
				continue
			}
		}

		if err := f(file); err != nil {
			return nil
		}

		if flags.verbose {
			toRemove := file
			if !filepath.IsAbs(file) {
				toRemove = filepath.Join(workingPath, file)
			}
			fmt.Printf("removed '%v'\n", toRemove)
		}
	}
	return nil
}

func main() {
	var flags rmFlags
	flag.BoolVar(&flags.verbose, "v", false, "Verbose mode.")
	flag.BoolVar(&flags.recursive, "r", false, "Recursive mode.")
	flag.BoolVar(&flags.interactive, "i", false, "Interactive mode.")
	fmt.Printf("proof of update")
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
	}

	if err := rm(flag.Args(), flags); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
