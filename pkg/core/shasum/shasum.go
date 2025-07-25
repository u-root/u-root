// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package shasum implements the shasum core utility.
package shasum

import (
	"bufio"
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"hash"
	"io"
	"os"

	"github.com/u-root/u-root/pkg/core"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

// command implements the shasum core utility.
type command struct {
	core.Base
}

// New creates a new shasum command.
func New() core.Command {
	c := &command{}
	c.Init()
	return c
}

type flags struct {
	algorithm int
}

// shaGenerator generates SHA hash of given data. The
// value of algorithm is expected to be 1 for SHA1
// 256 for SHA256
// and 512 for SHA512
func (c *command) shaGenerator(r io.Reader, algo int) ([]byte, error) {
	var h hash.Hash
	switch algo {
	case 1:
		h = sha1.New()
	case 256:
		h = sha256.New()
	case 512:
		h = sha512.New()
	default:
		return nil, fmt.Errorf("invalid algorithm, only 1, 256 or 512 are valid: %w", os.ErrInvalid)
	}
	if _, err := io.Copy(h, r); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

// runShasum processes the files and generates their SHA hashes.
func (c *command) runShasum(algorithm int, args []string) error {
	var hashbytes []byte
	var err error
	if len(args) == 0 {
		buf := bufio.NewReader(c.Stdin)
		if hashbytes, err = c.shaGenerator(buf, algorithm); err != nil {
			return err
		}
		fmt.Fprintf(c.Stdout, "%x -\n", hashbytes)
		return nil
	}
	for _, arg := range args {
		resolvedPath := c.ResolvePath(arg)
		file, err := os.Open(resolvedPath)
		if err != nil {
			return err
		}
		defer file.Close()
		if hashbytes, err = c.shaGenerator(file, algorithm); err != nil {
			return err
		}
		fmt.Fprintf(c.Stdout, "%x %s\n", hashbytes, arg)
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

	fs := flag.NewFlagSet("shasum", flag.ContinueOnError)
	fs.SetOutput(c.Stderr)

	fs.IntVar(&f.algorithm, "algorithm", 1, "SHA algorithm, valid args are 1, 256 and 512")
	fs.IntVar(&f.algorithm, "a", 1, "SHA algorithm, valid args are 1, 256 and 512")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: shasum -a <algorithm> <File Name>\n\n")
		fmt.Fprintf(fs.Output(), "shasum computes SHA checksums of files.\n")
		fmt.Fprintf(fs.Output(), "If no files are specified, read from stdin.\n\n")
		fmt.Fprintf(fs.Output(), "Options:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(unixflag.ArgsToGoArgs(args)); err != nil {
		return err
	}

	if err := c.runShasum(f.algorithm, fs.Args()); err != nil {
		return err
	}

	return nil
}
