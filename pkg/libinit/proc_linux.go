// Copyright 2014-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package libinit

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

// WaitOrphans waits for all remaining processes on the system to exit.
func WaitOrphans() uint {
	var numReaped uint
	for {
		var (
			s syscall.WaitStatus
			r syscall.Rusage
		)
		p, err := syscall.Wait4(-1, &s, 0, &r)
		if p == -1 {
			break
		}
		log.Printf("%v: exited with %v, status %v, rusage %v", p, err, s, r)
		numReaped++
	}
	return numReaped
}

// WithTTYControl turns on controlling the TTY on this command.
func WithTTYControl(ctty bool) CommandModifier {
	return func(c *exec.Cmd) {
		if c.SysProcAttr == nil {
			c.SysProcAttr = &syscall.SysProcAttr{}
		}
		c.SysProcAttr.Setctty = ctty
		c.SysProcAttr.Setsid = ctty
	}
}

// WithCloneFlags adds clone(2) flags to the *exec.Cmd.
func WithCloneFlags(flags uintptr) CommandModifier {
	return func(c *exec.Cmd) {
		if c.SysProcAttr == nil {
			c.SysProcAttr = &syscall.SysProcAttr{}
		}
		c.SysProcAttr.Cloneflags = flags
	}
}

func init() {
	osDefault = linuxDefault
}

func linuxDefault(c *exec.Cmd) {
	c.SysProcAttr = &syscall.SysProcAttr{
		Setctty: true,
		Setsid:  true,
	}
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
			var s syscall.WaitStatus
			var r syscall.Rusage
			if p, err := syscall.Wait4(-1, &s, 0, &r); p == cmd.Process.Pid {
				debug("Shell exited, exit status %d", s.ExitStatus())
				break
			} else if p != -1 {
				debug("Reaped PID %d, exit status %d", p, s.ExitStatus())
			} else {
				debug("Error from Wait4 for orphaned child: %v", err)
				break
			}
		}
		if err := cmd.Process.Release(); err != nil {
			log.Printf("Error releasing process %v: %v", cmd, err)
		}
	}
	return cmdCount
}
