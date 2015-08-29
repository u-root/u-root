// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 move (rename) files
 created by Beletti (rhiguita@gmail.com)
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
)

func usage() {
	fmt.Printf("usage: mv source target\n")
	fmt.Printf("       mv source ... directory\n")
	os.Exit(1)
}

func main() {
	todir := false
	flag.Parse()

	if flag.NArg() < 2 {
		usage()
	}

	files := flag.Args()
	lf := files[len(files)-1]
	lfdir, err := os.Stat(lf)
	if err == nil {
		todir = lfdir.IsDir()
	}
	if flag.NArg() > 2 && todir == false {
		fmt.Printf("not a directory: %s\n", lf)
		os.Exit(1)
	}

	if len(files) == 2 && todir == false {
		// rename file
		err := os.Rename(files[0], files[1])
		if err != nil {
			log.Fatalf("%v", err)
		}
	} else {
		// "copying" N files to 1 directory
		for i := 0; i < flag.NArg()-1; i++ {
			ndir := path.Join(lf, files[i])
			err := os.Rename(files[i], ndir)
			if err != nil {
				log.Fatalf("%v", err)
			}
		}
	}
}
