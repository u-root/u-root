// Copyright 2012-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Run a command, repeatedly, until it succeeds or we are out of time
//
// Synopsis:
//	backoff [-t duration-string] command [args...]
//
// Description:
//	backoff will run the command until it succeeds or a timeout has passed.
//	The default timeout is 30s.
//	If no args are given, it will just return.
//
// Example:
//	$ backoff echo hi
//	hi
//	$

//go:build !test
// +build !test

package main

import (
	"flag"
	"log"
	"os"
)

var timeout = flag.String("t", "", "Timeout for command")

func main() {
	if err := runit(*timeout, os.Args[0], os.Args[1:]...); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
