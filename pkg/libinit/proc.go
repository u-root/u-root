// Copyright 2014-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package libinit

import (
	"io"
	"os"
	"os/exec"

	"github.com/u-root/u-root/pkg/upath"
)

var osDefault = func(*exec.Cmd) {}

// CommandModifier makes *exec.Cmd construction modular.
type CommandModifier func(c *exec.Cmd)

// WithArguments adds command-line arguments to a command.
func WithArguments(arg ...string) CommandModifier {
	return func(c *exec.Cmd) {
		if len(arg) > 0 {
			c.Args = append(c.Args, arg...)
		}
	}
}

// WithStdin changes the command's stdin to r.
func WithStdin(r io.Reader) CommandModifier {
	return func(c *exec.Cmd) {
		c.Stdin = r
	}
}

// WithStdout changes the command's stdout to w.
func WithStdout(w io.Writer) CommandModifier {
	return func(c *exec.Cmd) {
		c.Stdout = w
	}
}

// WithStderr changes the command's stderr to w.
func WithStderr(w io.Writer) CommandModifier {
	return func(c *exec.Cmd) {
		c.Stderr = w
	}
}

// Command constructs an *exec.Cmd object.
func Command(bin string, m ...CommandModifier) *exec.Cmd {
	bin = upath.UrootPath(bin)
	cmd := exec.Command(bin)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	osDefault(cmd)
	for _, mod := range m {
		mod(cmd)
	}
	return cmd
}
