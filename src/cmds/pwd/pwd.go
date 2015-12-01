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
)

func pwd() string {

	if *physical {
		if path, err := os.Getwd(); err != nil {
			log.Fatalf("%v", err)
		} else {
			path, _ = filepath.EvalSymlinks(path)
			return path
		}
	} else if *logical {
		if path, err := os.Getwd(); err != nil {
			log.Fatalf("%v", err)
		} else {
			return path
		}
	}

	return "" //It wull never come to this

}

func main() {

	flag.Parse()
	fmt.Printf("$s\n", pwd())

}
