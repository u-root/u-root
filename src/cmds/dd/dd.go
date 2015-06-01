// Copyright 2013 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

/*
dd is modeled after dd. Each step in the chain is a goroutine that
reads a block and writes a block.
There are two always-the goroutines, in and out. They're actually
the same thing save they have, maybe, different block sizes.
*/
import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type passer func(r io.Reader, w io.Writer, ibs, obs int)

var (
	ibs     = flag.Int("ibs", 1, "Default input block size")
	obs     = flag.Int("obs", 1, "Default output block size")
	skip = flag.Int("skip", 0, "skip n bytes before reading")
	seek = flag.Int("seek", 0, "seek output when writing")
	count = flag.Int("count", max.MaxUInt, "Max output of data to copy")
	inName  = flag.String("if", "", "Input file")
	outName = flag.String("of", "", "Output file")
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
			nn, err := w.Write(b[tot : tot+obs])
			if err != nil {
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

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%v", err)
	os.Exit(1)
}

// rather than, in essence, recreating all the apparatus of flag.xxxx with the if= bits,
// including dup checking, conversion, etc. we just convert the arguments and then
// run flag.Parse. Gross, but hey, it works.
func main() {
<<<<<<< HEAD
	var err error
	flag.Parse()
	for _, v := range flag.Args() {
		l := strings.SplitN(v, "=", 2)
		switch l[0] {
		case "if":
			inFile, err = os.Open(l[1])
		case "of":
			outFile, err = os.OpenFile(l[1], os.O_CREATE|os.O_WRONLY, 0600)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	}
=======
	inFile := os.Stdin
	outFile := os.Stdout
	var err error
	// EVERYTHING in dd follows x=y. So blindly split and convert and sleep well.
	arg := []string{}
	for _, v := range os.Args {
		l := strings.SplitN(v, "=", 2)
		l[0] = "-" + l[0]
		arg = append(arg, l...)
	}
	os.Args = arg
	flag.Parse()
	if *inName != "" {
		inFile, err = os.Open(*inName)
		if err != nil {
			fatal(err)
		}
	}
	if *outName != "" {
		outFile, err = os.OpenFile(*outName, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fatal(err)
		}
	}

>>>>>>> Take dd a bit further: if= and of=
	r, w := io.Pipe()
	go pass(inFile, w, *ibs, *ibs)
	pass(r, outFile, *obs, *obs)
}
