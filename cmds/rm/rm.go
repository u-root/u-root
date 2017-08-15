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
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	//"io"
)

// You can add more flags to this struct
type rmFlags struct {
	recursive   bool
	verbose     bool
	interactive bool
	//force       bool
}

func recursiveDelete(file string, flags rmFlags) error {
	input := bufio.NewScanner(os.Stdin)
	statval, err := os.Stat(file)
	if err != nil {
		return err
	}
	if statval.IsDir() {
		//Throws error if flag is not recursive
		if !flags.recursive {
			//because it is a path error, the print statement is different and cannot be returned
			newError := os.PathError{Op: "\nrm:", Path: file, Err: syscall.EISDIR}
			fmt.Fprintf(os.Stderr, "%v\n", newError.Error())
			return nil
		}
		//At this point, the recursive flag is on
		if flags.interactive || flags.verbose {
			//TODO: sort this list by (unknown parameter)
			fileList, err := ioutil.ReadDir(file)
			if err != nil {
				return err
			}
			if len(fileList) == 0 {
				if flags.interactive {
					fmt.Printf("rm: remove directory '%v'? ", file)
					input.Scan()
					if input.Text()[0] != 'y' {
						return nil
					}
				}
			} else if err != nil {
				return err
			} else {
				if flags.interactive {
					fmt.Printf("rm: descend into directory '%v'? ", file)
					input.Scan()
					if input.Text()[0] != 'y' {
						return nil
					}
				}
				for _, each := range fileList {
					recursiveDelete(filepath.Join(file, each.Name()), flags)
				}
			}

		} else {
			os.RemoveAll(file)
		}
		//This step is for post processing after all the files inside a directory are removed
		if flags.interactive {
			fmt.Printf("rm: remove directory '%v'? ", file)
			input.Scan()
			if input.Text()[0] != 'y' {
				return nil
			}
		}
		if flags.verbose {
			fmt.Printf("removed directory '%v'\n", file)
		}

	} else {

		if flags.interactive {
			if statval.Size() == 0 {
				fmt.Printf("rm: remove regular empty file '%v'? ", file)
			} else {
				fmt.Printf("rm: remove '%v'? ", file)
			}
			input.Scan()
			if input.Text()[0] != 'y' {
				return nil
			}
		}
		if err := os.Remove(file); err != nil {
			return err
		}
		if flags.verbose {
			fmt.Printf("removed '%v'\n", file)
		}
	}

	return nil
}

func rm(files []string, flags rmFlags) error {
	for _, file := range files {
		if err := recursiveDelete(file, flags); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	var flags rmFlags
	flag.BoolVar(&flags.verbose, "v", false, "Verbose mode.")
	flag.BoolVar(&flags.recursive, "r", false, "Recursive mode.")
	flag.BoolVar(&flags.interactive, "i", false, "Interactive mode.")
	//flag.BoolVar(&flags.interactive, "f", false, "Force mode")
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
	}

	if err := rm(flag.Args(), flags); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
