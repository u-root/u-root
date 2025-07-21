// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rm implements the rm core utility.
package rm

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/u-root/u-root/pkg/core"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

// command implements the rm core utility.
type command struct {
	core.Base
}

// New creates a new rm command.
func New() core.Command {
	c := &command{}
	c.Init()
	return c
}

type flags struct {
	interactive bool
	verbose     bool
	recursive   bool
	r           bool
	force       bool
}

const usage = "rm [-Rrvif] file..."

// promptRemove asks the user if they want to remove the file.
func (c *command) promptRemove(file string) (bool, error) {
	fmt.Fprintf(c.Stderr, "rm: remove '%v'? ", file)
	reader := bufio.NewReader(c.Stdin)
	answer, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	return strings.ToLower(answer)[0] == 'y', nil
}

// removeFiles removes the specified files according to the flags.
func (c *command) removeFiles(files []string, f flags) error {
	if len(files) < 1 {
		return fmt.Errorf("%v", usage)
	}

	removeFunc := os.Remove
	if f.recursive || f.r {
		removeFunc = os.RemoveAll
	}

	if f.force {
		f.interactive = false
	}

	workingPath := c.WorkingDir
	if workingPath == "" {
		var err error
		workingPath, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	for _, file := range files {
		resolvedFile := c.ResolvePath(file)

		if f.interactive {
			shouldRemove, err := c.promptRemove(file)
			if err != nil {
				return err
			}
			if !shouldRemove {
				continue
			}
		}

		if err := removeFunc(resolvedFile); err != nil {
			if f.force && os.IsNotExist(err) {
				continue
			}
			return err
		}

		if f.verbose {
			toRemove := file
			if !path.IsAbs(file) {
				toRemove = filepath.Join(workingPath, file)
			}
			fmt.Fprintf(c.Stdout, "removed '%v'\n", toRemove)
		}
	}
	return nil
}

// Run executes the command with a `context.Background()`.
func (c *command) Run(args ...string) error {
	return c.RunContext(context.Background(), args...)
}

// Run executes the command.
func (c *command) RunContext(ctx context.Context, args ...string) error {
	var f flags

	fs := flag.NewFlagSet("rm", flag.ContinueOnError)
	fs.SetOutput(c.Stderr)

	fs.BoolVar(&f.interactive, "i", false, "Interactive mode.")
	fs.BoolVar(&f.verbose, "v", false, "Verbose mode.")
	fs.BoolVar(&f.recursive, "r", false, "equivalent to -R")
	fs.BoolVar(&f.r, "R", false, "Recursive, remove hierarchies")
	fs.BoolVar(&f.force, "f", false, "Force, ignore nonexistent files and never prompt")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: %s\n", usage)
		fs.PrintDefaults()
	}

	if err := fs.Parse(unixflag.ArgsToGoArgs(args)); err != nil {
		return err
	}

	if err := c.removeFiles(fs.Args(), f); err != nil {
		return err
	}

	return nil
}
