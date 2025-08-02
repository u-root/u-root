// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cat implements the cat core utility.
package cat

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/u-root/u-root/pkg/core"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

// command implements the cat core utility.
type command struct {
	core.Base
}

// New creates a new cat command.
func New() core.Command {
	c := &command{}
	c.Init()
	return c
}

type flags struct {
	u bool // ignored flag for compatibility
}

var errCopy = fmt.Errorf("error concatenating stdin to stdout")

// cat copies data from reader to writer.
func (c *command) cat(reader io.Reader, writer io.Writer) error {
	if _, err := io.Copy(writer, reader); err != nil {
		return errCopy
	}
	return nil
}

// runCat processes the files and concatenates them to stdout.
func (c *command) runCat(args []string) error {
	if len(args) == 0 {
		return c.cat(c.Stdin, c.Stdout)
	}

	for _, file := range args {
		if file == "-" {
			err := c.cat(c.Stdin, c.Stdout)
			if err != nil {
				return err
			}
			continue
		}

		resolvedFile := c.ResolvePath(file)
		f, err := os.Open(resolvedFile)
		if err != nil {
			return err
		}

		if err := c.cat(f, c.Stdout); err != nil {
			f.Close()
			return fmt.Errorf("failed to concatenate file %s to given writer", f.Name())
		}
		f.Close()
	}
	return nil
}

// Run executes the command with a `context.Background()`.
func (c *command) Run(args ...string) error {
	return c.RunContext(context.Background(), args...)
}

// RunContext executes the command.
func (c *command) RunContext(ctx context.Context, args ...string) error {
	var f flags

	fs := flag.NewFlagSet("cat", flag.ContinueOnError)
	fs.SetOutput(c.Stderr)

	fs.BoolVar(&f.u, "u", false, "ignored")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: cat [-u] [FILES]...\n\n")
		fmt.Fprintf(fs.Output(), "cat concatenates files and prints them to stdout.\n")
		fmt.Fprintf(fs.Output(), "If no files are specified, read from stdin.\n\n")
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(unixflag.ArgsToGoArgs(args)); err != nil {
		return err
	}

	if err := c.runCat(fs.Args()); err != nil {
		return err
	}

	return nil
}
