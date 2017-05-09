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
//
// Bugs: in i mode, it can't use non-seekable stdin, i.e. a pipe. Yep, this sucks.
// But if we implement seek on such things, we have to do it by reading, which
// really sucks. It's doable, we'll do it if we have to, but for now I'd like
// to avoid the complexity. cpio is a 40 year old concept. If you want something
// better, see ../archive which has a VTOC and separates data from metadata (unlike cpio).
// We could test for ESPIPE and fix it that way ... later.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/u-root/u-root/cmds/cpio/pkg"
	_ "github.com/u-root/u-root/cmds/cpio/pkg/newc"
)

var (
	debug = func(string, ...interface{}) {}
	d     = flag.Bool("v", false, "Debug prints")
	format = flag.String("H", "newc", "format")
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
		var r cpio.RecReader
		if r, err = cpio.Reader(*format, os.Stdin); err == nil {
			var f *cpio.File
			for f, err = r.RecRead(); err == nil; f, err = r.RecRead() {
				fmt.Printf("%s\n", f.String())
				err = cpio.Create(f)
				if err != nil {
					fmt.Printf("%v: %v", f, err)
				}
			}
		}

	case "o":
		var w cpio.RecWriter
		if w, err = cpio.Writer(*format, os.Stdout); err != nil {
			log.Fatal(err)
		}

		b := bufio.NewReader(os.Stdin)

		for {
			var name string
			if name, err = b.ReadString('\n'); err != nil {
				if err == io.EOF {
					err = w.Finish()
				}
				break
			}
			name = strings.TrimRight(name, "\r\n")
			fi, err := os.Lstat(name)
			if err != nil {
				break
			}
			f, err := cpio.FIToFile(name, fi)
			if err != nil {
				break
			}
			_, err = w.RecWrite(f)
			if err != nil {
				break
			}
		}
	case "t":
		var r cpio.RecReader
		if r, err = cpio.Reader(*format, os.Stdin); err == nil {
			var f *cpio.File
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
