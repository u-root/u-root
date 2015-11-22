// Copyright 2013 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//Tee transcribes the standard input to the standard output and makes copies in the files.
package main

import (
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
)

var (
	append = flag.Bool("a", false, "append the output to the files rather than rewriting them")
	ignore = flag.Bool("i", false, "ignore the SIGINT signal")
)

//Copy any input up to n bytes from buffer to Stdout and files
func copyinput(files []*os.File, buf [8192]byte, n int) {

	os.Stdout.Write(buf[:n])
	for _, v := range files {
		v.Write(buf[:n])
	}

}

//Parses all the flags and sets variables accordingly
func handleflags() int {

	flag.Parse()

	oflags := os.O_WRONLY | os.O_CREATE

	if *append {
		oflags |= os.O_APPEND
	}

	if *ignore {
		signal.Ignore(os.Interrupt)
	}

	return oflags
}

func main() {

	var buf [8192]byte

	oflags := handleflags()

	files := make([]*os.File, flag.NArg())

	for i, v := range flag.Args() {
		f, err := os.OpenFile(v, oflags, 0666)
		if err != nil {
			log.Fatalf("error opening %s: %v", v, err)
		}
		files[i] = f
	}

	for {
		n, err := os.Stdin.Read(buf[:])
		if err != nil {
			if err != io.EOF {
				log.Fatalf("error reading stdin: %v\n", err)
			}
			break
		}
		copyinput(files, buf, n)
	}

}
