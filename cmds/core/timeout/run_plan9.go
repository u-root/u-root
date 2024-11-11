// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"errors"
	"os/exec"
)

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
