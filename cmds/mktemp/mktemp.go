// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Mktemp makes a temporary file (or directory)
//
// Synopsis:
//     mktemp [-d] [-p]
//
// Options:
//     -d: make a temp directory instead of file
//     -p: change prefix used for template
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var (
	directory = flag.Bool("d", false, "Create a directory instead of file")
	prefix    = flag.String("p", "", "Add prefix to end of temp creation")
)

func init() {
	// Usage Definition
	defUsage := flag.Usage
	flag.Usage = func() {
		os.Args[0] = "mktemp [-dp]"
		defUsage()
	}
}

func main() {
	flag.Parse()
	if len(flag.Args()) > 0 {
		flag.Usage()
		os.Exit(1)
	}

	if *directory {
		dirName, err := ioutil.TempDir("", *prefix)
		if err != nil {
			log.Fatalf("Error creating temp dir: %v\n", err)
		}
		fmt.Println(dirName)
		os.Exit(0)
	}

	file, err := ioutil.TempFile("", *prefix)
	if err != nil {
		log.Fatalf("Error creating temp file: %v\n", err)
	}
	fmt.Println(file.Name())
}
