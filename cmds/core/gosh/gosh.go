// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Derived work from Daniel Martí <mvdan@mvdan.cc>

//go:build (!tinygo || tinygo.enable) && !plan9

package main

import (
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

var command = flag.String("c", "", "Command to run")

func main() {
	flag.Parse()

	err := run(os.Stdin, os.Stdout, os.Stderr, *command, flag.Args()...)
	if status, ok := interp.IsExitStatus(err); ok {
		os.Exit(int(status))
	}
	if err != nil {
		log.Fatal(err)
	}
}

var errNotImplemented = errors.New("fancy interactive interpreter not implemented")

func run(stdin io.Reader, stdout, stderr io.Writer, command string, args ...string) error {
	runner, err := interp.New(interp.StdIO(stdin, stdout, stderr))
	if err != nil {
		return err
	}

	if command != "" {
		return runReader(runner, strings.NewReader(command), "")
	}
	if len(args) == 0 {
		if r, ok := stdin.(*os.File); ok && term.IsTerminal(int(r.Fd())) {
			if err := runInteractive(runner, syntax.NewParser(), stdout, stderr); !errors.Is(err, errNotImplemented) {
				return err
			}
			return runInteractiveSimple(runner, stdin, stdout)
		}
		return runReader(runner, stdin, "")
	}
	return runScript(runner, args[0])
}

func runScript(runner *interp.Runner, file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return runReader(runner, f, file)
}

func runReader(runner *interp.Runner, reader io.Reader, name string) error {
	prog, err := syntax.NewParser().Parse(reader, name)
	if err != nil {
		return err
	}
	runner.Reset()
	return runner.Run(context.Background(), prog)
}

func runInteractiveSimple(runner *interp.Runner, stdin io.Reader, stdout io.Writer) error {
	parser := syntax.NewParser()
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
