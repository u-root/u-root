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

	"github.com/u-root/u-root/pkg/cpio"
	_ "github.com/u-root/u-root/pkg/cpio/newc"
)

var (
	debug  = func(string, ...interface{}) {}
	d      = flag.Bool("v", false, "Debug prints")
	format = flag.String("H", "newc", "format")
)

func usage() {
	log.Fatalf("Usage: cpio")
}

func main() {
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

	archiver, err := cpio.Format(*format)
	if err != nil {
		log.Fatalf("Format %q not supported: %v", *format, err)
	}

	switch op {
	case "i":
		rr := archiver.Reader(os.Stdin)
		for {
			rec, err := rr.ReadRecord()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error reading records: %v", err)
			}
			debug("Creating %s\n", rec)
			if err := cpio.CreateFile(rec); err != nil {
				log.Printf("Creating %q failed: %v", rec.Name, err)
			}
		}

	case "o":
		rw := archiver.Writer(os.Stdout)
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			name := scanner.Text()
			rec, err := cpio.GetRecord(name)
			if err != nil {
				log.Fatalf("Getting record of %q failed: %v", name, err)
			}
			if err := rw.WriteRecord(rec); err != nil {
				log.Fatalf("Writing record %q failed: %v", name, err)
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatalf("Error reading stdin: %v", err)
		}
		if err := rw.WriteTrailer(); err != nil {
			log.Fatalf("Error writing trailer record: %v", err)
		}

	case "t":
		rr := archiver.Reader(os.Stdin)
		for {
			rec, err := rr.ReadRecord()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error reading records: %v", err)
			}
			fmt.Println(rec)
		}

	default:
		usage()
	}
}
