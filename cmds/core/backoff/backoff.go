// Copyright 2012-2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Run a command, repeatedly, until it succeeds or we are out of time
//
// Synopsis:
//	backoff [-v] [-t duration-string] command [args...]
//
// Description:
//	backoff will run the command until it succeeds or a timeout has passed.
//	The default timeout is 30s.
//	If -v is set, it will show what it is running, each time it is tried.
//	If no args are given, it will print command help.
//
// Example:
//	$ backoff echo hi
//	hi
//	$
//	$ backoff -v -t=2s false
//	  2022/03/31 14:29:37 Run ["false"]
//	  2022/03/31 14:29:37 Set timeout to 2s
//	  2022/03/31 14:29:37 "false" []:exit status 1
//	  2022/03/31 14:29:38 "false" []:exit status 1
//	  2022/03/31 14:29:39 "false" []:exit status 1
//	  2022/03/31 14:29:39 Error: exit status 1

//go:build !test

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/cenkalti/backoff/v4"
)

var (
	timeout = flag.Duration("t", 30*time.Second, "Timeout for command")
	verbose = flag.Bool("v", false, "Log each attempt to run the command")
	v       = func(string, ...any) {}

	errNoCmd = fmt.Errorf("no command passed")
)

func run(timeout time.Duration, c string, a ...string) error {
	if c == "" {
		return errNoCmd
	}
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = timeout
	f := func() error {
		cmd := exec.Command(c, a...)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		err := cmd.Run()
		v("%q %q:%v", c, a, err)
		return err
	}

	return backoff.Retry(f, b)
}

func main() {
	flag.Parse()
	if *verbose {
		v = log.Printf
	}
	a := flag.Args()
	if len(a) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	v("Run %q", a)
	if err := run(*timeout, a[0], a[1:]...); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
