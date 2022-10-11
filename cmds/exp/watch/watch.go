// Copyright 2015-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// watch periodically executes the executable specified in argument.
//
// Synopsis:
//
//	watch [-n] SEC [-t] cmd-exec
//
// Description:
//
//	cmd-exec is executed every n seconds
//	example, watch -n 5 dmesg >> log.txt
//	: executes dmesg every 5 sec and stores the log in log.txt
//
// Options:
//
//	-n: time in seconds
//	-t: do not print header
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	t = flag.Bool("t", false, "Don't print header")
	n = flag.Int64("n", 2, "Loop period in SEC, default 2")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "watch [-n SEC] [-t] PROG ARGS\n")
		fmt.Fprintf(os.Stderr, "Run PROG Periodically\n")

		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	argRem := flag.Args()

	// argRem is the remaining args(non flag) after parsing.
	if len(argRem) == 0 {
		flag.Usage()
		os.Exit(0)
	}

	seconds := uint64(*n)
	if seconds <= 0 {
		seconds = 2
	}

	for {
		fmt.Print("\033[0;0H")
		fmt.Print("\033[J")
		if !*t {
			fmt.Printf("Every %d : %v \n\n", seconds, argRem)
		}
		cmd := exec.Command(argRem[0], argRem[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			if strings.Contains(err.Error(), "executable file not found") {
				fmt.Print(err)
			}
		}

		time.Sleep(time.Second * time.Duration(seconds))
	}
}
