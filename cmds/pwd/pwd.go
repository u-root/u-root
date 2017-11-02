// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Print name of current directory.
//
// Synopsis:
//     pwd [-LP]
//
// Options:
//     -L: follow symlinks (default)
//     -P: don't follow symlinks
//
// Author:
//     created by Beletti (rhiguita@gmail.com)
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	logical  = flag.Bool("L", true, "Follow symlinks") // this is the default behavior
	physical = flag.Bool("P", false, "Don't follow symlinks")
	cmd      = "pwd [-LP]"
)

func init() {
	defUsage := flag.Usage
	flag.Usage = func() {
		os.Args[0] = cmd
		defUsage()
	}
}

func pwd() error {
	path, err := os.Getwd()
	if err == nil && *physical {
		path, err = filepath.EvalSymlinks(path)
	}

	if err == nil {
		fmt.Println(path)
	}

	return err
}

func main() {
	args := os.Args[1:]
	flag.Parse()
	for _, flag := range args {
		switch flag {
		case "-L":
			*physical = false
		case "-P":
			*physical = true
		}
	}

	if err := pwd(); err != nil {
		log.Fatalf("%v", err)
	}
}
