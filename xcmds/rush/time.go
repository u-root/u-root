// Copyright 2012 the u-root Authors. All rights reserved
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
//     CMD can be a builtin, e.g. time cd
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
	"fmt"
	"os"
	"os/exec"
	"time"
)

func init() {
	addBuiltIn("time", runtime)
}

func printTime(label string, t time.Duration) {
	fmt.Fprintf(os.Stderr, "%s %.03f\n", label, t.Seconds())
}

func runtime(c *Command) error {
	var err error
	start := time.Now()
	if len(c.argv) > 0 {
		c.cmd = c.argv[0]
		c.argv = c.argv[1:]
		c.args = c.args[1:]
		// If we are in a builtin, then the lookup failed.
		// The result of the failed lookup remains in
		// c.Cmd and will make start fail. We have to make
		// a new Cmd.
		nCmd := exec.Command(c.cmd, c.argv[:]...)
		nCmd.Stdin = c.Stdin
		nCmd.Stdout = c.Stdout
		nCmd.Stderr = c.Stderr
		c.Cmd = nCmd
		err = runit(c)
	}
	realTime := time.Since(start)
	printTime("real", realTime)
	if c.ProcessState != nil {
		printTime("user", c.ProcessState.UserTime())
		printTime("sys", c.ProcessState.SystemTime())
	}
	return err
}
