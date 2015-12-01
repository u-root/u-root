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
	logical  = flag.Bool("L", true, "Follow symlinks") //this is the default behavior
	physical = flag.Bool("P", false, "Don't follow symlinks")
	cmd      = "pwd [-LP]"
)

func usage() {
	fmt.Printf("Usage: %v", cmd)
	flag.PrintDefaults()
}

func init() {
	args := os.Args[1:] // saving the args
	flag.Usage = usage
	flag.Parse()
	var L, P int
	for index, flag := range args {
		if flag == "-L" {
			L = index
		} else if flag == "-P" {
			P = index
		}
	}
	if P > L {
		*physical = true // if P appears after
	} else if L > P {
		*physical = false // if L appears after
	}
}

func pwd() {
	if *physical {
		if path, err := os.Getwd(); err != nil {
			log.Fatalf("%v", err)
		} else {
			path, _ = filepath.EvalSymlinks(path)
			fmt.Println(path)
		}
	} else if *logical {
		if path, err := os.Getwd(); err != nil {
			log.Fatalf("%v", err)
		} else {
			fmt.Println(path)
		}
	}

}

func main() {
	pwd()
}
