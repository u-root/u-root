// Copyright 2013 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Ls reads the directories in the command line and prints out the names.

The options are:
	–l		Long form.
*/

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var (
	long      = flag.Bool("l", false, "Long form")
	recursive = flag.Bool("R", false, "Recurse")
)

func main() {
	flag.Parse()

	dirs := flag.Args()

	if len(dirs) == 0 {
		dirs = []string{"."}
	}
	for _, v := range dirs {
		err := filepath.Walk(v, func(path string, fi os.FileInfo, err error) error {
			fmt.Printf("%v: %v\n", v, fi)
			if fi.IsDir() && !*recursive {
				return filepath.SkipDir
			}
			return err
		})
		if err != nil {
			fmt.Printf("%s: %v\n", v, err)
		}
	}
}
