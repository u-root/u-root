// Copyright 2014-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package libinit

import (
	"golang.org/x/sys/plan9"
	"log"
	"os"
	"os/exec"
)

// WaitOrphans waits for all remaining processes on the system to exit.
func WaitOrphans() uint {
	var numReaped uint
	for {
		var w plan9.Waitmsg
		err := plan9.Await(&w)
		if err != nil {
			break
		}
		log.Printf("Exited with %v", w)
		numReaped++
	}
	return numReaped
}

// WithRforkFlags adds rfork flags to the *exec.Cmd.
func WithRforkFlags(flags uintptr) CommandModifier {
	return func(c *exec.Cmd) {
		if c.SysProcAttr == nil {
			c.SysProcAttr = &plan9.Rfork{}
		}
		c.SysProcAttr.Rfork = int(flags)
	}
}

func init() {
	osDefault = plan9Default
}

func plan9Default(c *exec.Cmd) {
	c.SysProcAttr = &plan9.Rfork{}
}

// FIX ME: make it not linux-specific
// RunCommands runs commands in sequence.
//
// RunCommands returns how many commands existed and were attempted to run.
//
// commands must refer to absolute paths at the moment.
func RunCommands(debug func(string, ...interface{}), commands ...*exec.Cmd) int {
	var cmdCount int
	for _, cmd := range commands {
		if _, err := os.Stat(cmd.Path); os.IsNotExist(err) {
			debug("%v", err)
			continue
		}

		cmdCount++
		debug("Trying to run %v", cmd)
		if err := cmd.Start(); err != nil {
			log.Printf("Error starting %v: %v", cmd, err)
			continue
		}

		for {
			var w plan9.Waitmsg
			if err := plan9.Await(&w); err != nil {
				debug("Error from Await: %v", err)
				break
			}
			if w.Pid == cmd.Process.Pid {
				debug("Shell exited, exit status %v", w)
				break
			}
			debug("Reaped PID %d, exit status %v", w.Pid, w)
		}
		if err := cmd.Process.Release(); err != nil {
			log.Printf("Error releasing process %v: %v", cmd, err)
		}
	}
	return cmdCount
}
