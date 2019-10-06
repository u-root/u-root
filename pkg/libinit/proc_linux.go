// Copyright 2014-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package libinit

import (
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/u-root/u-root/pkg/upath"
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

// RunCommands runs commands in sequence.
//
// commands must refer to absolute paths at the moment.
func RunCommands(debug func(string, ...interface{}), commands ...*exec.Cmd) {
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
	if cmdCount == 0 {
		log.Printf("No suitable executable found in %v", commands)
	}
}

// CommandModifier makes *exec.Cmd construction modular.
type CommandModifier func(c *exec.Cmd)

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

// WithArguments adds command-line arguments to a command.
func WithArguments(arg ...string) CommandModifier {
	return func(c *exec.Cmd) {
		c.Args = append(c.Args, arg...)
	}
}

// Command constructs an *exec.Cmd object.
func Command(bin string, m ...CommandModifier) *exec.Cmd {
	bin = upath.UrootPath(bin)
	cmd := exec.Command(bin)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	// By default, this stuff is on.
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setctty: true,
		Setsid:  true,
	}
	for _, mod := range m {
		mod(cmd)
	}
	return cmd
}
