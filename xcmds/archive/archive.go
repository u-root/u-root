// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Archive archives files.
//
//
// Synopsis:
//     archive d|e|t [-d] [args...]
//
// Description:
//     The VTOC is at the front; we're not modeling tape drives or streams as
//     in tar and cpio. This will greatly speed up listing the archive,
//     modifying it, and so on. We think. Why a new tool?
//
// Options:
//     d: decode
//     e: toc (table of contents)
//     t: toc (table of contents)
//     -d: debug prints
package main

import (
	"flag"
	"log"
	"os"
)

// You'll see the name VTOC used a lot.
// The Volume Table Of Contents (extra points for looking this
// up) is an array of structs listing the file names, and
// their info. We'll see how much is needed -- almost
// certainly more than we think.  No attempt is made hee to be
// space efficient. Disks are big, and metadata is much the
// smallest part of it all. I think.
// The VTOC goes after the data, because the VTOC size
// is actually dependent on the values of the offsets in the
// individual structs. We have to write the data, fill in the
// VTOC, and write the encoded VTOC out.
// VTOC is a []file.

var (
	debug = func(string, ...interface{}) {}
	d     = flag.Bool("d", false, "Debug prints")
)

func usage() {
	log.Fatalf("Usage: archive d|e|t [args...]")
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
	case "e":
		err = encode(os.Stdout, a[1:]...)
	case "t":
		err = toc(a[1:]...)
	case "d":
		err = decode(a[1:]...)
	default:
		usage()
	}
	if err != nil {
		log.Fatalf("%v", err)
	}
}
