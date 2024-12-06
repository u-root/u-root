// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9 && !windows

package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

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
