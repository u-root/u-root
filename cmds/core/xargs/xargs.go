// Copyright 2013-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

const defaultMaxArgs = 5000
const tty = "/dev/tty"

type params struct {
	maxArgs int
	trace   bool
	prompt  bool
}

type cmd struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
	tty    string
	params
}

func command(stdin io.Reader, stdout, stderr io.Writer, p params) *cmd {
	return &cmd{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
		tty:    tty,
		params: p,
	}
}

func parseParams() params {
	var maxArgs = flag.Int("n", defaultMaxArgs, "max number of arguments per command")
	var trace = flag.Bool("t", false, "enable trace mode, each command is written to stderr")
	var prompt = flag.Bool("p", false, "the user is asked whether to execute utility at each invocation")

	flag.Parse()
	p := params{
		maxArgs: *maxArgs,
		trace:   *trace || *prompt,
		prompt:  *prompt,
	}

	return p
}

func main() {
	c := command(os.Stdin, os.Stdout, os.Stderr, parseParams())
	if err := c.run(flag.Args()...); err != nil {
		log.Fatal(err)
	}
}

func (c *cmd) run(args ...string) error {
	if len(args) == 0 {
		args = append(args, "echo")
	}

	var xArgs []string
	scanner := bufio.NewScanner(c.stdin)
	for scanner.Scan() {
		sp := strings.Fields(scanner.Text())
		xArgs = append(xArgs, sp...)
	}

	argsLen := len(args)
	var ttyScanner *bufio.Scanner
	if c.prompt {
		var err error
		f, err := os.Open(c.tty)
		if err != nil {
			return err
		}
		ttyScanner = bufio.NewScanner(f)
	}

	for i := 0; i < len(xArgs); i += c.maxArgs {
		m := len(xArgs)
		if i+c.maxArgs < m {
			m = i + c.maxArgs
		}
		args = append(args, xArgs[i:m]...)

		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdin = c.stdin
		cmd.Stdout = c.stdout
		cmd.Stderr = c.stderr

		if c.prompt {
			fmt.Fprintf(c.stderr, "%s...?", strings.Join(args, " "))
		} else if c.trace {
			fmt.Fprintf(c.stderr, "%s\n", strings.Join(args, " "))
		}

		if c.prompt && ttyScanner.Scan() {
			input := ttyScanner.Text()
			if !strings.HasPrefix(input, "y") && !strings.HasPrefix(input, "Y") {
				args = args[:argsLen]
				continue
			}
		}

		if err := cmd.Run(); err != nil {
			return err
		}

		args = args[:argsLen]
	}

	return nil
}
