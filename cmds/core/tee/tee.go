// Copyright 2013-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

// Tee transcribes the standard input to the standard output and makes copies
// in the files.
//
// Synopsis:
//
//	tee [-ai] FILES...
//
// Options:
//
//	-a, --append: append the output to the files rather than rewriting them
//	-i, --ignore-interrupts: ignore the SIGINT signal
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"

	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

type cmd struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
	args   []string
	cat    bool
	ignore bool
}

func (c *cmd) run() error {
	oflags := os.O_WRONLY | os.O_CREATE
	if c.cat {
		oflags |= os.O_APPEND
	}

	if c.ignore {
		signal.Ignore(os.Interrupt)
	}

	files := make([]*os.File, 0, len(c.args))
	writers := make([]io.Writer, 0, len(c.args)+1)
	for _, fname := range c.args {
		f, err := os.OpenFile(fname, oflags, 0o666)
		if err != nil {
			return fmt.Errorf("error opening %s: %w", fname, err)
		}
		files = append(files, f)
		writers = append(writers, f)
	}
	writers = append(writers, c.stdout)

	mw := io.MultiWriter(writers...)
	if _, err := io.Copy(mw, c.stdin); err != nil {
		return fmt.Errorf("error: %w", err)
	}

	for _, f := range files {
		if err := f.Close(); err != nil {
			fmt.Fprintf(c.stderr, "tee: error closing file %q: %v\n", f.Name(), err)
		}
	}

	return nil
}

func command(args []string) *cmd {
	c := &cmd{
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}

	f := flag.NewFlagSet(args[0], flag.ExitOnError)

	f.BoolVar(&c.cat, "append", false, "append the output to the files rather than rewriting them")
	f.BoolVar(&c.cat, "a", false, "append the output to the files rather than rewriting them")

	f.BoolVar(&c.ignore, "ignore-interrupts", false, "ignore the SIGINT signal")
	f.BoolVar(&c.ignore, "i", false, "ignore the SIGINT signal")

	f.Parse(unixflag.ArgsToGoArgs(args[1:]))
	c.args = f.Args()

	return c
}

func main() {
	if err := command(os.Args).run(); err != nil {
		log.Fatalf("tee: %v", err)
	}
}
