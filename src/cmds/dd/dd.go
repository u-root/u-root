// Copyright 2013 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

/*
dd is modeled after dd. Each step in the chain is a goroutine that
reads a block and writes a block.
There are two always-there goroutines, in and out. They're actually
the same thing save they have, maybe, different block sizes.
*/
import (
	"flag"
	"fmt"
	"io"
	"os"
)

type passer func(r io.Reader, w io.Writer, ibs, obs int)

var (
	ibs = flag.Int("ibs", 1, "Default input block size")
	obs = flag.Int("obs", 1, "Default output block size")
)

func pass(r io.Reader, w io.WriteCloser, ibs, obs int) {
	b := make([]byte, ibs) 
	defer w.Close()
	for {
		bs := 0
		for bs < ibs {
			n, err := r.Read(b[bs:])
			if err != nil && err != io.EOF {
				return
			}
			if n == 0 {
				break
			}
			bs += n
		}
		if bs == 0 {
			return
		}
		tot := 0
		for tot < bs {
			nn, err := w.Write(b[tot:tot+obs])
			if  err != nil {
				fmt.Fprintf(os.Stderr, "pass: %v\n", err)
				return
			}
			if nn == 0 {
				return
			}
			tot += nn
		}
	}
}

func main() {
	flag.Parse()
	r, w := io.Pipe()
	go pass(os.Stdin, w, *ibs, *ibs)
	pass(r, os.Stdout, *obs, *obs)
}
