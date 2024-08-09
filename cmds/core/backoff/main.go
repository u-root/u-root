// Copyright 2012-2021 the u-root Authors. All rights reserved
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
//	If no args are given, it will just return.
//
// Example:
//	$ backoff echo hi
//	hi
//	$
//	$ backoff -v -t 2s false
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
	"log"
	"os"
)

var (
	timeout = flag.String("t", "", "Timeout for command")
	verbose = flag.Bool("v", true, "Log each attempt to run the command")
	v       = func(string, ...interface{}) {}
)

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
	if err := runit(*timeout, a[0], a[1:]...); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
