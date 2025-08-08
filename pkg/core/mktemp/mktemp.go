// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package mktemp implements the mktemp core utility.
package mktemp

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/u-root/u-root/pkg/core"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

// command implements the mktemp core utility.
type command struct {
	core.Base
}

// New creates a new mktemp command.
func New() core.Command {
	c := &command{}
	c.Init()
	return c
}

type flags struct {
	d      bool
	u      bool
	q      bool
	prefix string
	suffix string
	dir    string
}

func (c *command) mktemp(f flags) (string, error) {
	dir := f.dir
	if dir == "" {
		if tmpdir := c.Getenv("TMPDIR"); tmpdir != "" {
			dir = tmpdir
		} else {
			dir = "/tmp"
		}
	}

	if f.u {
		if !f.q {
			log.Printf("Not doing anything but dry-run is an inherently unsafe concept")
		}
		return "", nil
	}

	if f.d {
		d, err := os.MkdirTemp(dir, f.prefix)
		return d, err
	}
	file, err := os.CreateTemp(dir, f.prefix)
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}

// Run executes the command with a `context.Background()`.
func (c *command) Run(args ...string) error {
	return c.RunContext(context.Background(), args...)
}

// RunContext executes the command.
func (c *command) RunContext(ctx context.Context, args ...string) error {
	var f flags

	fs := flag.NewFlagSet("mktemp", flag.ContinueOnError)
	fs.SetOutput(c.Stderr)

	fs.BoolVar(&f.d, "directory", false, "Make a directory")
	fs.BoolVar(&f.d, "d", false, "Make a directory (shorthand)")

	fs.BoolVar(&f.u, "dry-run", false, "Do everything save the actual create")
	fs.BoolVar(&f.u, "u", false, "Do everything save the actual create (shorthand)")

	fs.BoolVar(&f.q, "quiet", false, "Quiet: show no errors")
	fs.BoolVar(&f.q, "q", false, "Quiet: show no errors (shorthand)")

	fs.StringVar(&f.prefix, "prefix", "", "add a prefix")
	fs.StringVar(&f.prefix, "s", "", "add a prefix (shorthand, 's' is for compatibility with GNU mktemp")

	fs.StringVar(&f.suffix, "suffix", "", "add a suffix to the prefix (rather than the end of the mktemp file)")

	fs.StringVar(&f.dir, "tmpdir", "", "Tmp directory to use. If this is not set, TMPDIR is used, else /tmp")
	fs.StringVar(&f.dir, "p", "", "Tmp directory to use. If this is not set, TMPDIR is used, else /tmp (shorthand)")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: mktemp [options] [template]\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(unixflag.ArgsToGoArgs(args)); err != nil {
		return err
	}

	switch fs.NArg() {
	case 1:
		f.prefix = f.prefix + strings.Split(fs.Args()[0], "X")[0] + f.suffix
	case 0:
	default:
		fs.Usage()
		return fmt.Errorf("too many arguments")
	}

	fileName, err := c.mktemp(f)
	if err != nil && !f.q {
		return err
	}

	fmt.Fprintf(c.Stdout, "%s\n", fileName)
	return nil
}
