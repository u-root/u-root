// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Forth is a forth interpreter.
// It reads a line at a time and puts it through the interpreter.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/forth"
)

var debug = flag.Bool("d", false, "Turn on stack dump after each Eval")

func main() {
	var b = make([]byte, 512)
	flag.Parse()
	f := forth.New()
	for {
		n, err := os.Stdin.Read(b)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			// Silently exit on EOF. It's the unix way.
			break
		}
		// NOTE: should be f.Eval. Why did I not do that? There was a reason ...
		// I don't remember what it was
		s, err := forth.Eval(f, string(b[:n]))
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		if *debug {
			fmt.Printf("%v", f.Stack())
		}
		fmt.Printf("%s\n", s)
		// And push it back. It's much more convenient to have it
		// always on TOS.
		f.Push(s)
	}
}
