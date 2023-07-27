// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Derived work from Daniel Martí <mvdan@mvdan.cc>

//go:build !tinygo && !plan9
// +build !tinygo,!plan9

package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"

	"golang.org/x/term"

	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

func main() {
	flag.Parse()
	err := run(os.Stdin, os.Stdout, os.Stderr, flag.Args()...)

	if status, ok := interp.IsExitStatus(err); ok {
		os.Exit(int(status))
	}

	if err != nil {
		log.Fatal(err)
	}
}

var errNotImplemented = errors.New("fancy interactive interpreter not implemented")

func run(stdin io.Reader, stdout, stderr io.Writer, args ...string) error {
	runner, err := interp.New(interp.StdIO(stdin, stdout, stderr))
	if err != nil {
		return err
	}

	parser := syntax.NewParser()

	if len(args) > 0 {
		if strings.HasSuffix(args[0], "sh") {
			return runScript(runner, parser, args[0])
		}

		return runCmd(runner, parser, strings.NewReader(strings.Join(args, " ")), args[0])
	}

	if r, ok := stdin.(*os.File); ok && term.IsTerminal(int(r.Fd())) {
		if err := runInteractive(runner, parser, stdout, stderr); !errors.Is(err, errNotImplemented) {
			return err
		}
		return runInteractiveSimple(runner, parser, stdin, stdout)
	}

	return runCmd(runner, parser, stdin, "")
}

func runScript(runner *interp.Runner, parser *syntax.Parser, file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	prog, err := parser.Parse(f, file)
	if err != nil {
		return err
	}

	runner.Reset()

	return runner.Run(context.Background(), prog)
}

func runCmd(runner *interp.Runner, parser *syntax.Parser, command io.Reader, name string) error {
	scanner := bufio.NewScanner(command)
	defer runner.Reset()

	for scanner.Scan() {
		h := scanner.Text()
		if err := scanner.Err(); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		prog, err := parser.Parse(bytes.NewBuffer([]byte(h)), name)
		if err != nil {
			return err
		}

		if err := runner.Run(context.Background(), prog); err != nil {
			return err
		}
	}
	return nil
}

func runInteractiveSimple(runner *interp.Runner, parser *syntax.Parser, stdin io.Reader, stdout io.Writer) error {
	fmt.Fprintf(stdout, "$ ")

	var runErr error

	// The following code is used to intercept SIGINT signals.
	// Calling signal.Ignore wouldn't work as child prcesses inherit this trait.
	// We only want to catch SIGINTs that are propagated from a child,
	// the child itself should get the signal as per usual.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func(ch chan os.Signal) {
		for {
			<-ch
		}
	}(ch)

	for {
		fn := func(stmts []*syntax.Stmt) bool {
			if parser.Incomplete() {
				fmt.Fprintf(stdout, "> ")
				return true
			}
			for _, stmt := range stmts {
				runErr = runner.Run(context.Background(), stmt)
				if runner.Exited() {
					return false
				}
			}
			fmt.Fprintf(stdout, "$ ")
			return true
		}

		if err := parser.Interactive(stdin, fn); err != nil {
			return err
		}

		if runErr != nil {
			fmt.Fprintf(stdout, "error: %s\n", runErr.Error())
			runErr = nil
		} else {
			return nil
		}
	}
}
