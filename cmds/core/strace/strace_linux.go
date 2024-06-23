// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !tinygo && linux && (amd64 || riscv64 || arm64)
// +build !tinygo
// +build linux
// +build amd64 riscv64 arm64

// strace is a simple multi-process syscall & signal tracer.
//
// Synopsis:
//
//	strace <command> [args...]
//
// Description:
//
//	trace a single process given a command name.
package main

import (
	// Don't use spf13 flags. It will not allow commands like
	// strace ls -l
	// it tries to use the -l for strace instead of leaving it alone.
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/u-root/u-root/pkg/strace"
)

var errUsage = errors.New("usage: strace [-o <outputfile>] <command> [args...]")

type params struct {
	output string
}

func run(stdin io.Reader, stdout, stderr io.Writer, p params, args ...string) error {
	if len(args) < 1 {
		return errUsage
	}

	c := exec.Command(args[0], args[1:]...)
	c.Stdin, c.Stdout, c.Stderr = stdin, stdout, stderr

	if p.output != "" {
		f, err := os.Create(p.output)
		if err != nil {
			return fmt.Errorf("creating out file: %s: %w", p.output, err)
		}
		defer f.Close()
		c.Stderr = f
	}

	return strace.Strace(c, c.Stderr)
}

func main() {
	output := flag.String("o", "", "write output to file (if empty, stdout)")
	flag.Parse()

	if err := run(os.Stdin, os.Stdout, os.Stderr, params{output: *output}, flag.Args()...); err != nil {
		log.Fatal(err)
	}
}
