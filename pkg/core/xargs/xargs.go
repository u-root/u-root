// Copyright 2013-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xargs implements the xargs core utility.
package xargs

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/u-root/u-root/pkg/core"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

const (
	defaultMaxArgs = 5000
	defaultTTY     = "/dev/tty"
)

// command implements the xargs core utility.
type command struct {
	core.Base
	tty string
}

// New creates a new xargs command.
func New() core.Command {
	c := &command{
		tty: defaultTTY,
	}
	c.Init()
	return c
}

type flags struct {
	maxArgs int
	trace   bool
	prompt  bool
	null    bool
}

// Run executes the command with a `context.Background()`.
func (c *command) Run(args ...string) error {
	return c.RunContext(context.Background(), args...)
}

// Run executes the command.
func (c *command) RunContext(ctx context.Context, args ...string) error {
	var f flags

	fs := flag.NewFlagSet("xargs", flag.ContinueOnError)
	fs.SetOutput(c.Stderr)

	fs.IntVar(&f.maxArgs, "n", defaultMaxArgs, "max number of arguments per command")
	fs.BoolVar(&f.trace, "t", false, "enable trace mode, each command is written to stderr")
	fs.BoolVar(&f.prompt, "p", false, "the user is asked whether to execute utility at each invocation")
	fs.BoolVar(&f.null, "0", false, "use a null byte as the input argument delimiter and do not treat any other input bytes as special")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: xargs [OPTIONS] [COMMAND [ARGS]...]\n\n")
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(unixflag.ArgsToGoArgs(args)); err != nil {
		return err
	}

	// Enable trace if prompt is enabled
	if f.prompt {
		f.trace = true
	}

	cmdArgs := fs.Args()
	if len(cmdArgs) == 0 {
		cmdArgs = append(cmdArgs, "echo")
	}

	var xArgs []string

	if f.null {
		r := bufio.NewReader(c.Stdin)
		for {
			b, err := r.ReadBytes(0x00)
			if err != nil && err != io.EOF {
				return err
			}
			if len(b) != 0 {
				if b[len(b)-1] == 0x00 {
					xArgs = append(xArgs, string(b[:len(b)-1]))
				} else {
					xArgs = append(xArgs, string(b))
				}
			}
			if err == io.EOF {
				break
			}
		}
	} else {
		scanner := bufio.NewScanner(c.Stdin)
		for scanner.Scan() {
			sp := strings.Fields(scanner.Text())
			xArgs = append(xArgs, sp...)
		}
	}

	argsLen := len(cmdArgs)
	var ttyScanner *bufio.Scanner
	if f.prompt {
		ttyFile, err := os.Open(c.tty)
		if err != nil {
			return err
		}
		defer ttyFile.Close()
		ttyScanner = bufio.NewScanner(ttyFile)
	}

	for i := 0; i < len(xArgs); i += f.maxArgs {
		m := min(i+f.maxArgs, len(xArgs))
		cmdArgs = append(cmdArgs, xArgs[i:m]...)

		cmd := exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)
		cmd.Stdin = c.Stdin
		cmd.Stdout = c.Stdout
		cmd.Stderr = c.Stderr

		if f.prompt {
			fmt.Fprintf(c.Stderr, "%s...?", strings.Join(cmdArgs, " "))
		} else if f.trace {
			fmt.Fprintf(c.Stderr, "%s\n", strings.Join(cmdArgs, " "))
		}

		if f.prompt && ttyScanner.Scan() {
			input := ttyScanner.Text()
			if !strings.HasPrefix(input, "y") && !strings.HasPrefix(input, "Y") {
				cmdArgs = cmdArgs[:argsLen]
				continue
			}
		}

		if err := cmd.Run(); err != nil {
			return err
		}

		cmdArgs = cmdArgs[:argsLen]
	}

	return nil
}

// SetTTY sets the TTY device path for prompt mode.
func SetTTY(c core.Command, tty string) {
	cmd := c.(*command)
	cmd.tty = tty
}
