// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
//  move (rename) files
// created by Beletti (rhiguita@gmail.com)
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s source target\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "       %s source ... directory\n", os.Args[0])
	os.Exit(1)
}


func mv(files []string, todir bool) (error) {
	if len(files) == 2 && todir == false {
		if err := os.Rename(files[0], files[1]); err != nil {
			return err
		}
	} else {
		lf := files[len(files)-1]
		// "copying" N files to 1 directory
		for i := range files[:len(files)-1] {
			ndir := path.Join(lf, path.Base(files[i]))
			if err := os.Rename(files[i], ndir); err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	var todir bool
	flag.Parse()

	if flag.NArg() < 2 {
		usage()
	}

	files := flag.Args()
	lf := files[len(files)-1]
	if lfdir, err := os.Lstat(lf); err == nil {
		todir = lfdir.IsDir()
	}
	if flag.NArg() > 2 && todir == false {
		fmt.Printf("not a directory: %s\n", lf)
		os.Exit(1)
	}

	if err := mv(files,todir); err != nil {
		log.Fatalf("%v", err)
	}
}
