// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// runtime runs the command and prints the time it took.
// The command can be a builtin, e.g.
// time time time time time time date
// works fine.

package main

import (
	"fmt"
	"os/exec"
	"time"
)

func init() {
	addBuiltIn("time", runtime)
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
	cost := time.Since(start)
	fmt.Fprintf(c.Stderr, "%v\n", cost)
	return err
}
