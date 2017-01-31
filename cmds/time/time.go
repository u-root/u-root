// Copyright 2017 the u-root Authors. All rights reserved
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
	"syscall"
	"time"
)

func printTime(label string, t time.Duration) {
	fmt.Fprintf(os.Stderr, "%s %.03f\n", label, t.Seconds())
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		log.Fatal("too few arguments")
	}

	// Run command
	cmd := exec.Command(flag.Arg(0), flag.Args()[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	start := time.Now()
	err := cmd.Run()
	realTime := time.Since(start)

	// Check for io errors.
	exitErr, ok := err.(*exec.ExitError)
	if err != nil && !ok {
		log.Fatal(err)
	}

	// Print usage.
	printTime("real", realTime)
	printTime("user", cmd.ProcessState.UserTime())
	printTime("sys", cmd.ProcessState.SystemTime())

	// Propagate return value.
	if exitErr != nil {
		os.Exit(exitErr.Sys().(syscall.WaitStatus).ExitStatus())
	}
}
