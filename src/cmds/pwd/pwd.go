// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 print name of current/working directory
 created by Beletti (rhiguita@gmail.com)
*/

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

func usage() {
	fmt.Printf("Usage: %v\n", cmd)
	flag.PrintDefaults()
}

func init() {
	args := os.Args[1:]
	flag.Usage = usage
	flag.Parse()
	for _, flag := range args {
		switch flag {
		case "-L":
			*physical = false
		case "-P":
			*physical = true
		}
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
	if err := pwd(); err != nil {
		log.Fatalf("%v", err)
	}
}
