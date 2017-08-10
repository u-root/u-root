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
	"syscall"
)

// Used an array of flags in case more flags are added (see rm method signature)
/*
	interactive  = flag.Bool("i", false, "Interactive mode.")
	verbose      = flag.Bool("v", false, "Verbose mode.")
	hierarchies  = flag.Bool("r", false, "Remove file hierarchies")
	flagsRm = [interactive, verbose, hierarchies]
*/
// You can add more flags to this struct
type rmFlags struct {
	recursive   bool
	verbose     bool
	interactive bool
}

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
		//Throw an error if the file is a directory
		statval, err := os.Stat(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		if statval.IsDir() {
			newError := os.PathError{Op: "rm:", Path: file, Err: syscall.EISDIR}
			fmt.Fprintf(os.Stderr, "%v\n", newError.Error())
			continue
		}
		if flags.interactive {
			fmt.Printf("rm: remove '%v'? ", file)
			input.Scan()
			if input.Text()[0] != 'y' {
				continue
			}
		}
		if err := f(file); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
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
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
	}

	if err := rm(flag.Args(), flags); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
