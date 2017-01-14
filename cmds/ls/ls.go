// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Ls prints the contents of a directory.
//
// Synopsis:
//     ls [OPTIONS] [DIRS]...
//
// Options:
//     -l: Long form.
//     -r: raw (%v) form
//     -R: recurse
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var (
	long      = flag.Bool("l", false, "Long form")
	raw       = flag.Bool("r", false, "raw struct")
	recursive = flag.Bool("R", false, "Recurse")
)

func show(fullpath string, fi os.FileInfo) {
	switch {
	case *raw == true:
		fmt.Printf("%v\n", fi)
	case *long == false:
		fmt.Printf("%v\n", fi.Name())
	// -rw-r--r-- 1 root root 174 Aug 18 17:18 /etc/hosts
	case *long == true:
		fmt.Printf("%v\t%v\t%v\t%v", fi.Mode(), fi.Size(), fi.Name(), fi.ModTime())
		if link, err := os.Readlink(fullpath); err == nil {
			fmt.Printf(" -> %v", link)
		}
		fmt.Printf("\n")
	}

}

func main() {
	flag.Parse()

	dirs := flag.Args()

	if len(dirs) == 0 {
		dirs = []string{"."}
	}
	for _, v := range dirs {
		if len(dirs) > 1 {
			fmt.Printf("%v:\n", v)
		}
		err := filepath.Walk(v, func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("%v: %v\n", path, err)
				return err
			}
			show(path, fi)
			if fi.IsDir() && !*recursive && path != v {
				return filepath.SkipDir
			}

			return err
		})
		if err != nil {
			fmt.Printf("%s: %v\n", v, err)
		}
	}
}
