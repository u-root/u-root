// Copyright 2013-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tee transcribes the standard input to the standard output and makes copies
// in the files.
//
// Synopsis:
//     tee [-ai] FILES...
//
// Options:
//     -a, --append: append the output to the files rather than rewriting them
//     -i, --ignore-interrupts: ignore the SIGINT signal
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"

	flag "github.com/spf13/pflag"
)

const name = "tee"

var (
	cat    = flag.BoolP("append", "a", false, "append the output to the files rather than rewriting them")
	ignore = flag.BoolP("ignore-interrupts", "i", false, "ignore the SIGINT signal")
)

// handeFlags parses all the flags and sets variables accordingly
func handleFlags() int {
	flag.Parse()

	oflags := os.O_WRONLY | os.O_CREATE

	if *cat {
		oflags |= os.O_APPEND
	}

	if *ignore {
		signal.Ignore(os.Interrupt)
	}

	return oflags
}

func main() {
	oflags := handleFlags()

	files := make([]*os.File, 0, flag.NArg())
	writers := make([]io.Writer, 0, flag.NArg()+1)
	for _, fname := range flag.Args() {
		f, err := os.OpenFile(fname, oflags, 0666)
		if err != nil {
			log.Fatalf("%s: error opening %s: %v", name, fname, err)
		}
		files = append(files, f)
		writers = append(writers, f)
	}
	writers = append(writers, os.Stdout)

	mw := io.MultiWriter(writers...)
	if _, err := io.Copy(mw, os.Stdin); err != nil {
		log.Fatalf("%s: error: %v", name, err)
	}

	for _, f := range files {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "%s: error closing file %q: %v\n", name, f.Name(), err)
		}
	}
}
