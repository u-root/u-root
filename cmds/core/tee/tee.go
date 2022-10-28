// Copyright 2013-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tee transcribes the standard input to the standard output and makes copies
// in the files.
//
// Synopsis:
//
//	tee [-ai] FILES...
//
// Options:
//
//	-a, --append: append the output to the files rather than rewriting them
//	-i, --ignore-interrupts: ignore the SIGINT signal
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"

	flag "github.com/spf13/pflag"
)

var (
	cat    = flag.BoolP("append", "a", false, "append the output to the files rather than rewriting them")
	ignore = flag.BoolP("ignore-interrupts", "i", false, "ignore the SIGINT signal")
)

func run(stdin io.Reader, stdout io.Writer, stderr io.Writer, args []string) error {
	oflags := os.O_WRONLY | os.O_CREATE
	if *cat {
		oflags |= os.O_APPEND
	}

	if *ignore {
		signal.Ignore(os.Interrupt)
	}

	files := make([]*os.File, 0, len(args))
	writers := make([]io.Writer, 0, len(args)+1)
	for _, fname := range args {
		f, err := os.OpenFile(fname, oflags, 0o666)
		if err != nil {
			return fmt.Errorf("error opening %s: %v", fname, err)
		}
		files = append(files, f)
		writers = append(writers, f)
	}
	writers = append(writers, stdout)

	mw := io.MultiWriter(writers...)
	if _, err := io.Copy(mw, stdin); err != nil {
		return fmt.Errorf("error: %v", err)
	}

	for _, f := range files {
		if err := f.Close(); err != nil {
			fmt.Fprintf(stderr, "tee: error closing file %q: %v\n", f.Name(), err)
		}
	}

	return nil
}

func main() {
	flag.Parse()
	if err := run(os.Stdin, os.Stdout, os.Stderr, flag.Args()); err != nil {
		log.Fatalf("tee: %v", err)
	}
}
