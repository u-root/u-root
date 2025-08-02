// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package touch implements the touch core utility.
package touch

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/u-root/u-root/pkg/core"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

// command implements the touch core utility.
type command struct {
	core.Base
}

// New creates a new touch command.
func New() core.Command {
	c := &command{}
	c.Init()
	return c
}

type flags struct {
	access       bool
	modification bool
	create       bool
	dateTime     string
}

type params struct {
	time         time.Time
	access       bool
	modification bool
	create       bool
}

// parseParams parses the command parameters and returns a params struct.
func (c *command) parseParams(dateTime string, access, modification, create bool) (params, error) {
	t := time.Now()
	if dateTime != "" {
		var err error
		t, err = time.Parse(time.RFC3339, dateTime)
		if err != nil {
			return params{}, err
		}
	}
	return params{
		access:       access || (!access && !modification),
		modification: modification || (!access && !modification),
		create:       create,
		time:         t,
	}, nil
}

// touchFiles processes the files according to the parameters.
func (c *command) touchFiles(p params, args []string) error {
	var errs error
	for _, arg := range args {
		resolvedArg := c.ResolvePath(arg)
		_, existsErr := os.Stat(resolvedArg)
		notExist := os.IsNotExist(existsErr)
		if notExist {
			if p.create {
				continue
			}

			f, err := os.Create(resolvedArg)
			if err != nil {
				errs = errors.Join(errs, err)
				continue
			}
			f.Close()
		}

		accessTime := time.Time{}
		if p.access || notExist {
			accessTime = p.time
		}
		modificationTime := time.Time{}
		if p.modification || notExist {
			modificationTime = p.time
		}

		err := os.Chtimes(resolvedArg, accessTime, modificationTime)
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}

	return errs
}

// Run executes the command with a `context.Background()`.
func (c *command) Run(args ...string) error {
	return c.RunContext(context.Background(), args...)
}

// Run executes the command.
func (c *command) RunContext(ctx context.Context, args ...string) error {
	var f flags

	fs := flag.NewFlagSet("touch", flag.ContinueOnError)
	fs.SetOutput(c.Stderr)

	fs.BoolVar(&f.access, "a", false, "change only the access time")
	fs.BoolVar(&f.modification, "m", false, "change only the modification time")
	fs.BoolVar(&f.create, "c", false, "do not create any file if it does not exist")
	fs.StringVar(&f.dateTime, "d", "", "use specified time instead of current time RFC3339")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: touch [-amc] [-d datetime] file...\n\n")
		fmt.Fprintf(fs.Output(), "touch changes file access and modification times.\n")
		fmt.Fprintf(fs.Output(), "If a file does not exist, it will be created unless -c is specified.\n\n")
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(unixflag.ArgsToGoArgs(args)); err != nil {
		return err
	}

	if len(fs.Args()) == 0 {
		fs.Usage()
		return fmt.Errorf("no files specified")
	}

	p, err := c.parseParams(f.dateTime, f.access, f.modification, f.create)
	if err != nil {
		return err
	}

	if err := c.touchFiles(p, fs.Args()); err != nil {
		return err
	}

	return nil
}
