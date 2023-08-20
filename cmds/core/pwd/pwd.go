// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Print name of current directory.
//
// Synopsis:
//
//	pwd [-LP]
//
// Options:
//
//	-P: don't follow symlinks
//
// Author:
//
//	created by Beletti (rhiguita@gmail.com)
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	// This is the default. Setting it to false doesn't do anything in GNU
	// or zsh pwd, because you just can't even set it to false.
	_ = flag.Bool("L", true, "don't follow any symlinks")

	physical = flag.Bool("P", false, "follow all symlinks (avoid all symlinks)")
)

func pwd(followSymlinks bool) (string, error) {
	path, err := os.Getwd()
	if err == nil && followSymlinks {
		path, err = filepath.EvalSymlinks(path)
	}
	return path, err
}

func main() {
	flag.Parse()

	path, err := pwd(*physical)
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Println(path)
}
