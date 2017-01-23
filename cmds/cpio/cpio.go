// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// cpio operates on cpio files using a cpio package
// It only implements basic cpio options.
//
//
// Synopsis:
//     cpio
//
// Description:
//
// Options:
//     o: output an archive to stdout given a pattern
//     i: output files from a stdin stream
//     t: print table of contents
//     -v: debug prints
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	debug = func(string, ...interface{}) {}
	d     = flag.Bool("v", false, "Debug prints")
)

func usage() {
	log.Fatalf("Usage: cpio")
}

func main() {
	var err error
	flag.Parse()
	if *d {
		debug = log.Printf
	}

	a := flag.Args()
	debug("Args %v", a)
	if len(a) < 1 {
		usage()
	}
	op := a[0]

	switch op {
	case "i":
		log.Fatalf("no extract yet")
	case "o":
		log.Fatalf("no extract yet")
	case "t":
		var r RecReader
		if r, err = NewcReader(os.Stdin); err == nil {
			var f *File
			for f, err = r.RecRead(); err == nil; f, err = r.RecRead() {
				fmt.Printf("%s\n", f.String())
			}
		}
	default:
		usage()
	}
	if err != nil && err != io.EOF {
		log.Fatalf("%v", err)
	}
}
