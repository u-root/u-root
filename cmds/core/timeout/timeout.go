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

//go:build !test

package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/exec"
	"time"
)

type cmd struct {
	args         []string
	timeout      time.Duration
	in, out, err *os.File
}

var (
	timeout   = flag.Duration("t", 30*time.Second, "Timeout for command")
	errNoArgs = errors.New("need at least a command to run")
)

func main() {
	flag.Parse()
	c := &cmd{args: flag.Args(), in: os.Stdin, out: os.Stdout, err: os.Stderr, timeout: *timeout}
	if errno, err := c.run(); err != nil || errno != 0 {
		log.Printf("timeout(%v):%v", *timeout, err)
		os.Exit(errno)
	}
}

func (c *cmd) run() (int, error) {
	if len(c.args) == 0 {
		return 1, errNoArgs
	}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	proc := exec.CommandContext(ctx, c.args[0], c.args[1:]...)
	proc.Stdin, proc.Stdout, proc.Stderr = c.in, c.out, c.err
	if err := proc.Run(); err != nil {
		errno := 1
		var e *exec.ExitError
		if errors.As(err, &e) {
			errno = e.ExitCode()
		}
		return errno, err
	}
	return 0, nil
}
