// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (!tinygo || tinygo.enable) && !plan9 && !goshsmall && !goshliner

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/knz/bubbline"
	"github.com/knz/bubbline/complete"
	"github.com/knz/bubbline/computil"
	"github.com/knz/bubbline/editline"

	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

// HistFile is the history file.
// This might, possibly, use GetPid to avoid gosh'es writing over each other
var HistFile = filepath.Join(os.TempDir(), "bubble-sh.history")

var completion = flag.Bool("comp", true, "Enable tabcompletion and a more feature rich editline implementation")

type candidate struct {
	repl       string
	moveRight  int
	deleteLeft int
}

func (m candidate) Replacement() string {
	return m.repl
}

func (m candidate) MoveRight() int {
	return m.moveRight
}

func (m candidate) DeleteLeft() int {
	return m.deleteLeft
}

type multiComplete struct {
	complete.Values
	moveRight  int
	deleteLeft int
}

func (m *multiComplete) Candidate(e complete.Entry) editline.Candidate {
	return candidate{e.Title(), m.moveRight, m.deleteLeft}
}

func autocompleteBubb(val [][]rune, line, col int) (msg string, completions editline.Completions) {
	word, wstart, wend := computil.FindWord(val, line, col)
	var candidates []string
	if wstart == 0 && !(strings.HasPrefix(word, ".") || strings.HasPrefix(word, "/")) {
		candidates = commandCompleter(word)
	} else {
		candidates = filepathCompleter(word)
	}

	if len(candidates) != 0 {
		return "", &multiComplete{
			Values:     complete.StringValues("suggestions", candidates),
			moveRight:  wend - col,
			deleteLeft: wend - wstart,
		}
	}
	return "", nil
}

func runInteractive(runner *interp.Runner, parser *syntax.Parser, stdout, stderr io.Writer) error {
	input := bubbline.New()
	// Set default window size to 80x24 in case ioctl isn't able to detect the actual window size
	input.Model.SetSize(80, 24)

	if err := input.LoadHistory(HistFile); err != nil {
		return err
	}

	input.SetAutoSaveHistory(HistFile, true)

	if *completion {
		input.AutoComplete = autocompleteBubb
	}

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
		if runErr != nil {
			fmt.Fprintf(stdout, "error: %s\n", runErr.Error())
			runErr = nil
		}

		line, err := input.GetLine()
		if err != nil {
			if err == io.EOF {
				break // maybe we should continue instead of break
			}
			if errors.Is(err, bubbline.ErrInterrupted) {
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
			input.AutoComplete = nil
			continue
		case "enablecomp":
			input.AutoComplete = autocompleteBubb
			continue
		default:
		}

		if line != "" {
			if err := input.AddHistory(line); err != nil {
				fmt.Fprintf(stdout, "unable to add %s to history: %v\n", line, err)
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
