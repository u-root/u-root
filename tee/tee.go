// Copyright 2013 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Tee transcribes the standard input to the standard output and makes copies in the files.

The options are:
      –a    Append the output to the files rather than rewriting them.
*/

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

var append = flag.Bool("a", false, "append the output to the files rather than rewriting them")

func main() {
	var buf [8192]byte

	flag.Parse()

	oflags := os.O_WRONLY | os.O_CREATE
	if *append {
		oflags |= os.O_APPEND
	}

	files := make([]*os.File, flag.NArg())
	for i, v := range flag.Args() {
		f, err := os.OpenFile(v, oflags, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error opening %s: %v", v, err)
			os.Exit(1)
		}
		files[i] = f
	}

	for {
		n, err := os.Stdin.Read(buf[:])
		if err != nil {
			if err != io.EOF {
				fmt.Fprintf(os.Stderr, "error reading stdin: %v\n", err)
				os.Exit(1)
			}
			break
		}

		os.Stdout.Write(buf[:n])
		for _, v := range files {
			v.Write(buf[:n])
		}
	}
}
