// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (!tinygo || tinygo.enable) && !plan9 && !goshsmall && goshliner

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/peterh/liner"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

// HistFile is the history file.
// This might, possibly, use GetPid to avoid gosh'es writing over each other
var HistFile = filepath.Join(os.TempDir(), "gosh.history")

var completion = flag.Bool("comp", true, "Enable tabcompletion and a more feature rich editline implementation")

func runInteractive(runner *interp.Runner, parser *syntax.Parser, stdout, stderr io.Writer) error {
	input := liner.NewLiner()
	defer input.Close()

	f, err := os.OpenFile(HistFile, os.O_RDWR|os.O_CREATE, 0)
	if err == nil {
		input.ReadHistory(f)
	} else if f, err = os.Open(HistFile); err != nil {
		log.Printf("Failed to open or create history file: %v", err)
	}
	if f != nil {
		defer f.Close()
	}

	input.SetCtrlCAborts(true)
	if *completion {
		input.SetCompleter(autocompleteLiner(parser))
	}

	var runErr error
	for {
		if runErr != nil {
			fmt.Fprintf(stdout, "error: %s\n", runErr.Error())
			runErr = nil
		}

		line, err := input.Prompt("$ ")
		if err != nil {
			if err == io.EOF {
				break // maybe we should continue instead of break
			}
			if errors.Is(err, liner.ErrPromptAborted) {
				fmt.Fprintf(stdout, "^C\n")
			} else {
				fmt.Fprintf(stderr, "error: %s\n", err.Error())
			}
			err = nil
			continue
		}

		switch line {
		case "exit":
			goto exit
		case "disablecomp":
			input.SetCompleter(nil)
			continue
		case "enablecomp":
			input.SetCompleter(autocompleteLiner(parser))
			continue
		default:
		}

		if line != "" {
			input.AppendHistory(line)
			if f != nil {
				input.WriteHistory(f)
			}
		}
		if err := parser.Stmts(strings.NewReader(line), func(stmt *syntax.Stmt) bool {
			if parser.Incomplete() {
				fmt.Fprintf(stdout, "-> ")
				return true
			}

			runErr = runner.Run(context.Background(), stmt)
			return !runner.Exited()
		}); err != nil {
			fmt.Fprintf(stderr, "error: %s\n", err.Error())
		}
	}
exit:
	return nil
}
