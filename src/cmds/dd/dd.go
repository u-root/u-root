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
	"math"
	"os"
	"strings"
)

type passer func(r io.Reader, w io.Writer, ibs, obs int, conv string)

var (
	ibs     = flag.Int("ibs", 1, "Default input block size")
	obs     = flag.Int("obs", 1, "Default output block size")
	bs      = flag.Int("bs", 0, "Default input and output block size")
	skip    = flag.Int64("skip", 0, "skip n bytes before reading")
	seek    = flag.Int64("seek", 0, "seek output when writing")
	conv    = flag.String("conv", "", "Convert the file on a specific way, like notrunc")
	count   = flag.Int64("count", math.MaxInt64, "Max output of data to copy")
	inName  = flag.String("if", "", "Input file")
	outName = flag.String("of", "", "Output file")
)

// The 'close' thing is a real hack, but needed for proper
// operation in single-process mode.
func pass(r io.Reader, w io.WriteCloser, ibs, obs int, conv string, close bool) {
	var err error
	var nn int
	b := make([]byte, ibs)
	defer func() {
		if close {
			w.Close()
		}
	}()
	for {
		bsc := 0
		for bsc < ibs {
			n, err := r.Read(b[bsc:])
			if err != nil && err != io.EOF {
				return
			}
			if n == 0 {
				break
			}
			bsc += n
		}
		if bsc == 0 {
			return
		}
		for tot := 0; tot < bsc; tot += nn {
			switch conv {
			case "ucase":
				nn, err = w.Write([]byte(strings.ToUpper(string(b[tot : tot+obs]))))
			case "lcase":
				nn, err = w.Write([]byte(strings.ToLower(string(b[tot : tot+obs]))))
			default:
				nn, err = w.Write(b[tot : tot+obs])
			}

			if err != nil {
				fmt.Fprintf(os.Stderr, "pass: %v\n", err)
				return
			}
			if nn == 0 {
				return
			}
		}
	}
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%v", err)
	os.Exit(1)
}

func SplitArgs() []string {
	// EVERYTHING in dd follows x=y. So blindly split and convert sleep well
	arg := []string{}
	for _, v := range os.Args {
		l := strings.SplitN(v, "=", 2)
		// We only fix the exact case for x=y.
		if len(l) == 2 {
			l[0] = "-" + l[0]
			arg = append(arg, l...)
		} else {
			arg = append(arg, l...)
		}
	}
	return arg
}

func OpenFiles() (os.File, os.File) {
	inFile := os.Stdin
	outFile := os.Stdout
	var err error

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

	// position things.
	if *skip > 0 {
		if _, err = inFile.Seek(*skip, 0); err != nil {
			fatal(err)
		}
	}
	if *seek > 0 {
		if _, err = outFile.Seek(*seek, 0); err != nil {
			fatal(err)
		}
	}
	// bs = both 'ibs' and 'obs' (IEEE Std 1003.1 - 2013)
	if *bs > 0 {
		*ibs = *bs
		*obs = *bs
	}

	return *inFile, *outFile
}

func InOut(inFile, outFile *os.File) {
	r, w := io.Pipe()
	go pass(inFile, w, *ibs, *ibs, *conv, true)
	// push other filters here as needed.
	pass(r, outFile, *obs, *obs, *conv, false)
}

// rather than, in essence, recreating all the apparatus of flag.xxxx with the if= bits,
// including dup checking, conversion, etc. we just convert the arguments and then
// run flag.Parse. Gross, but hey, it works.
func main() {
	os.Args = SplitArgs()
	flag.Parse()
	inFile, outFile := OpenFiles()
	InOut(&inFile, &outFile)
}
