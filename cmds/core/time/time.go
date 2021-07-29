// Copyright 2012-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Time process execution.
//
// Synopsis:
//     time CMD [ARG]...
//
// Description:
//     After executing CMD, its real, user and system times are printed to
//     stderr in the POSIX format.
//
// Example:
//     $ time sleep 1.23s
//     real 1.230
//     user 0.001
//     sys 0.000
//
// Note:
//     This is different from bash's time command which is built into the shell
//     and can time the entire pipeline.
//
// Bugs:
//     Time is not reported when exiting due to a signal.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func printTime(label string, t time.Duration) {
	fmt.Fprintf(os.Stderr, "%s %.03f\n", label, t.Seconds())
}

func main() {
	flag.Parse()
	a := flag.Args()
	start := time.Now()
	if len(a) == 0 {
		fmt.Fprintf(os.Stderr, "real 0.000\nuser 0.000\nsys 0.000\n")
		os.Exit(0)
	}
	c := exec.Command(a[0], a[1:]...)
	c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stdout, os.Stderr
	defer func(*exec.Cmd, time.Time) {
		realTime := time.Since(start)
		printTime("real", realTime)
		if c.ProcessState != nil {
			printTime("user", c.ProcessState.UserTime())
			printTime("sys", c.ProcessState.SystemTime())
		}
	}(c, start)
	if err := c.Run(); err != nil {
		log.Fatal(err)
	}
}
