// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tee transcribes the standard input to the standard output and makes copies
// in the files.
//
// Synopsis:
//     tee [-ai] FILES...
//
// Options:
//     -a: append the output to the files rather than rewriting them
//     -i: ignore the SIGINT signal
package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
)

var (
	cat    = flag.Bool("a", false, "append the output to the files rather than rewriting them")
	ignore = flag.Bool("i", false, "ignore the SIGINT signal")
)

// Copy any input from buffer to Stdout and files
func copyinput(files []io.Writer, buf []byte) error {
	for _, v := range files {
		if _, err := v.Write(buf); err != nil {
			return err
		}
	}
	return nil
}

// Parses all the flags and sets variables accordingly
func handleflags() int {
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
	oflags := handleflags()

	files := make([]io.Writer, flag.NArg())

	for i, v := range flag.Args() {
		f, err := os.OpenFile(v, oflags, 0666)
		if err != nil {
			log.Fatalf("error opening %s: %v", v, err)
		}
		files[i] = f
	}

	b := make([]byte, 1048576)
	files = append(files, os.Stdout)
	buf := bufio.NewReader(os.Stdin)

	for {
		if n, err := buf.Read(b[:]); err != nil {
			log.Fatalf("%v", err)
		} else {
			copyinput(files, b[:n])
		}
	}
}
