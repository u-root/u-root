// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package mkdir implements the mkdir core utility.
package mkdir

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/u-root/u-root/pkg/core"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

// command implements the mkdir core utility.
type command struct {
	core.Base
}

// New creates a new mkdir command.
func New() core.Command {
	c := &command{}
	c.Init()
	return c
}

type flags struct {
	mode    string
	mkall   bool
	verbose bool
}

const (
	defaultCreationMode = 0o777
	stickyBit           = 0o1000
	sgidBit             = 0o2000
	suidBit             = 0o4000
)

// parseMode parses the mode string and returns the appropriate FileMode.
func (c *command) parseMode(mode string) (os.FileMode, error) {
	var m uint64
	var err error
	if mode == "" {
		m = defaultCreationMode
	} else {
		m, err = strconv.ParseUint(mode, 8, 32)
		if err != nil || m > 0o7777 {
			return 0, fmt.Errorf("invalid mode %q", mode)
		}
	}

	createMode := os.FileMode(m)
	if m&stickyBit != 0 {
		createMode |= os.ModeSticky
	}
	if m&sgidBit != 0 {
		createMode |= os.ModeSetgid
	}
	if m&suidBit != 0 {
		createMode |= os.ModeSetuid
	}

	return createMode, nil
}

// mkdirFiles creates directories according to the flags.
func (c *command) mkdirFiles(f flags, args []string) error {
	mkdirFunc := os.Mkdir
	if f.mkall {
		mkdirFunc = os.MkdirAll
	}

	createMode, err := c.parseMode(f.mode)
	if err != nil {
		return err
	}

	for _, name := range args {
		resolvedName := c.ResolvePath(name)
		if err := mkdirFunc(resolvedName, createMode); err != nil {
			fmt.Fprintf(c.Stderr, "%v: %v\n", name, err)
			continue
		}
		if f.verbose {
			fmt.Fprintf(c.Stdout, "%v\n", name)
		}
		if f.mode != "" {
			os.Chmod(resolvedName, createMode)
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

	fs := flag.NewFlagSet("mkdir", flag.ContinueOnError)
	fs.SetOutput(c.Stderr)

	fs.StringVar(&f.mode, "m", "", "Directory mode")
	fs.BoolVar(&f.mkall, "p", false, "Make all needed directories in the path")
	fs.BoolVar(&f.verbose, "v", false, "Print each directory as it is made")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: mkdir [-m mode] [-v] [-p] DIRECTORY...\n\n")
		fmt.Fprintf(fs.Output(), "mkdir makes a new directory.\n\n")
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(unixflag.ArgsToGoArgs(args)); err != nil {
		return err
	}

	if len(fs.Args()) < 1 {
		fs.Usage()
		return fmt.Errorf("no directories specified")
	}

	if err := c.mkdirFiles(f, fs.Args()); err != nil {
		return err
	}

	return nil
}
