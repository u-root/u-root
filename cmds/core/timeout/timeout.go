// Copyright 2012-2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Run a command and kill it if it runs more than a specified duration
//
// Synopsis:
//	timeout [-t duration-string] command [args...]
//
// Description:
//	timeout will run the command until it succeeds or too much time has passed.
//	The default timeout is 30s.
//	If no args are given, it will print a usage error.
//
// Example:
//	$ timeout echo hi
//	hi
//	$
//	$./timeout -t 5s bash -c 'sleep 40'
//	$ 2022/03/31 14:47:32 signal: killed
//	$./timeout  -t 5s bash -c 'sleep 40'
//	$ 2022/03/31 14:47:40 signal: killed
//	$./timeout  -t 5s bash -c 'sleep 1'
//	$

//go:build !test && !windows

package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"time"
)

type cmd struct {
	args         []string
	timeout      time.Duration
	signal       string
	in, out, err *os.File
}

var (
	timeout   = flag.Duration("t", 30*time.Second, "Timeout for command")
	signal    = flag.String("signal", "TERM", "specify the signal to be sent on timeout")
	errNoArgs = errors.New("need at least a command to run")
)

func main() {
	flag.Parse()
	c := &cmd{args: flag.Args(), in: os.Stdin, out: os.Stdout, err: os.Stderr, timeout: *timeout, signal: *signal}
	if errno, err := c.run(); err != nil || errno != 0 {
		log.Printf("timeout(%v):%v", *timeout, err)
		os.Exit(errno)
	}
}
