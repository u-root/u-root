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
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
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

var sigmap = map[string]syscall.Signal{
	"KILL": syscall.SIGKILL,
	"TERM": syscall.SIGTERM,
}

func (c *cmd) run() (int, error) {
	if len(c.args) == 0 {
		return 1, errNoArgs
	}

	sig, ok := sigmap[c.signal]
	if !ok {
		return 1, fmt.Errorf("unknown signal: %q: %w", c.signal, os.ErrInvalid)
	}

	cmd := exec.Command(c.args[0], c.args[1:]...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = c.in, c.out, c.err
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := cmd.Start(); err != nil {
		return 1, err
	}

	time.AfterFunc(c.timeout, func() {
		syscall.Kill(-cmd.Process.Pid, sig)
	})

	if err := cmd.Wait(); err != nil {
		errno := 1
		var e *exec.ExitError
		if errors.As(err, &e) {
			errno = e.ExitCode()
		}
		return errno, err
	}

	return 0, nil
}

func main() {
	flag.Parse()
	c := &cmd{args: flag.Args(), in: os.Stdin, out: os.Stdout, err: os.Stderr, timeout: *timeout, signal: *signal}
	if errno, err := c.run(); err != nil || errno != 0 {
		log.Printf("timeout(%v):%v", *timeout, err)
		os.Exit(errno)
	}
}
