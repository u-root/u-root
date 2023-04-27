// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Derived work from Daniel Martí <mvdan@mvdan.cc>

//go:build !tinygo && !plan9
// +build !tinygo,!plan9

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/u-root/prompt"
	"github.com/u-root/prompt/completer"
	"golang.org/x/term"

	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

type input interface {
	Input(prefix string, completer prompt.Completer, opts ...prompt.Option) string
}

type inputPrompt struct{}

func (i inputPrompt) Input(prefix string, completer prompt.Completer, opts ...prompt.Option) string {
	return prompt.Input(prefix, completer, opts...)
}

type shell struct {
	input
}

func main() {
	flag.Parse()

	sh := shell{
		input: inputPrompt{},
	}

	err := sh.runAll(flag.NArg())

	if e, ok := interp.IsExitStatus(err); ok {
		os.Exit(int(e))
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (s shell) runAll(narg int) error {
	r, err := interp.New(interp.StdIO(os.Stdin, os.Stdout, os.Stderr))
	if err != nil {
		return err
	}

	if narg > 0 {
		return s.run(r, strings.NewReader(strings.Join(flag.Args(), " ")), "")
	}

	if narg == 0 {
		if term.IsTerminal(int(os.Stdin.Fd())) {
			return s.runInteractiveTabCompletion(r, os.Stdout)
		}

		return s.run(r, os.Stdin, "")
	}

	return nil
}

func (s shell) run(r *interp.Runner, reader io.Reader, name string) error {
	prog, err := syntax.NewParser().Parse(reader, name)
	if err != nil {
		return err
	}

	r.Reset()

	return r.Run(context.Background(), prog)
}

func (s shell) runInteractiveTabCompletion(r *interp.Runner, stdout io.Writer) error {
	parser := syntax.NewParser()

	if s.input == nil {
		s.input = inputPrompt{}
	}

	for {
		in := s.Input(
			"$ ",
			completerFunc,
			prompt.OptionCompletionWordSeparator(completer.FilePathCompletionSeparator),
		)

		if in == "exit" {
			break
		}

		if err := parser.Stmts(strings.NewReader(in), func(stmt *syntax.Stmt) bool {
			if parser.Incomplete() {
				fmt.Fprintf(stdout, "> ")

				return true
			}

			_ = r.Run(context.Background(), stmt)

			return !r.Exited()
		}); err != nil {
			return err
		}
	}

	return nil
}
