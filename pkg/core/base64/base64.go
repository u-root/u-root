// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package base64 implements the base64 core utility.
package base64

import (
	"context"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/u-root/u-root/pkg/core"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

// command implements the base64 core utility.
type command struct {
	core.Base
}

// New creates a new base64 command.
func New() core.Command {
	c := &command{}
	c.Init()
	return c
}

type flags struct {
	decode bool
}

var errBadUsage = errors.New("usage: base64 [-d] [file]")

// do performs the actual base64 encoding or decoding operation.
func (c *command) do(r io.Reader, w io.Writer, decode bool) error {
	if decode {
		r = base64.NewDecoder(base64.StdEncoding, r)
		if _, err := io.Copy(w, r); err != nil {
			return fmt.Errorf("base64: error decoding %w", err)
		}
		return nil
	}

	// WriteCloser is important here, from NewEncoder documentation:
	// when finished writing, the caller must Close the returned encoder
	// to flush any partially written blocks.
	wc := base64.NewEncoder(base64.StdEncoding, w)
	defer wc.Close()
	if _, err := io.Copy(wc, r); err != nil {
		return fmt.Errorf("base64: error encoding %w", err)
	}
	if err := wc.Close(); err != nil { // flush any remaining data
		return fmt.Errorf("base64: error closing encoder %w", err)
	}
	if _, err := fmt.Fprintln(w); err != nil { // add trailing newline
		return fmt.Errorf("base64: error writing newline %w", err)
	}
	return nil
}

// runBase64 processes the input and performs base64 encoding/decoding.
func (c *command) runBase64(decode bool, names []string) error {
	reader := c.Stdin

	switch len(names) {
	case 0:
		// Use stdin
	case 1:
		resolvedFile := c.ResolvePath(names[0])
		f, err := os.Open(resolvedFile)
		if err != nil {
			return err
		}
		defer f.Close()
		reader = f
	default:
		return errBadUsage
	}

	return c.do(reader, c.Stdout, decode)
}

// Run executes the command with a `context.Background()`.
func (c *command) Run(args ...string) error {
	return c.RunContext(context.Background(), args...)
}

// RunContext executes the command.
func (c *command) RunContext(ctx context.Context, args ...string) error {
	var f flags

	fs := flag.NewFlagSet("base64", flag.ContinueOnError)
	fs.SetOutput(c.Stderr)

	fs.BoolVar(&f.decode, "d", false, "Decode")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: base64 [-d] [FILE]\n\n")
		fmt.Fprintf(fs.Output(), "Encode or decode a file to or from base64 encoding.\n")
		fmt.Fprintf(fs.Output(), "For stdin, on standard Unix systems, you can use /dev/stdin\n\n")
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(unixflag.ArgsToGoArgs(args)); err != nil {
		return err
	}

	return c.runBase64(f.decode, fs.Args())
}
