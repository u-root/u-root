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

var debug = flag.Bool("d", false, "Turn on forth package debugging using log.Printf")

func main() {
	b := make([]byte, 512)
	flag.Parse()
	if *debug {
		forth.Debug = log.Printf
	}
	f := forth.New()
	for {
		fmt.Printf("%vOK\n", f.Stack())
		n, err := os.Stdin.Read(b)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			// Silently exit on EOF. It's the unix way.
			break
		}
		if err := forth.EvalString(f, string(b[:n])); err != nil {
			fmt.Printf("%v\n", err)
		}
	}
}
