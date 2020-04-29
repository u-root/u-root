// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package pogosh implements a small POSIX-compatible shell.
package pogosh

import (
	"fmt"
	"io/ioutil"
)

// Run executes the given fragment of shell.
//
// There are three circumstances for this function to return:
//
// 1. There is a shell error (compiler error, file not found, ...)
// 2. The script calls the exit builtin.
// 3. The script reaches the end of input.
//
// In the case the shell script execs another process, this function will not
// return (unless the Exec function is appropriately overriden).
func (s *State) Run(script string) (exitCode int, err error) {
	err = isASCII(script)
	if err != nil {
		return 0, err
	}

	// Lex
	tokens, err := tokenize(script)
	if err != nil {
		return 0, err
	}

	// Parse
	eof := token{script[len(script):], ttEOF}
	// TODO: append newline instead of EOF, or make EOF the only empty value ""
	tokens = append(tokens, eof) // augment
	t := tokenizer{tokens}
	cmd := parseProgram(s, &t)
	if t.ts[0].ttype != ttEOF {
		panic("expected EOF") // TODO: better error message
	}

	// Exit code
	defer func() {
		switch r := recover().(type) {
		case nil:
		case exitError:
			exitCode = r.code
		case error:
			exitCode = 1
			err = r
		default:
			panic(r) // TODO: clobbers stack trace
		}
	}()

	// Execute
	cmd.exec(s)

	return 0, err
}

// RunFile is a convenient wrapper around Run.
func (s *State) RunFile(filename string) (int, error) {
	script, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, err
	}
	return s.Run(string(script))
}

func isASCII(str string) error {
	lineNumber := 0
	lineStart := 0
	for i, v := range []byte(str) {
		if v == '\n' {
			lineNumber++
			lineStart = i + 1
		} else if v > 127 {
			// TODO: include filename if possible
			return fmt.Errorf("<pogosh>:%v:%v: non-ascii character, '\\x%x'", lineNumber+1, i-lineStart+1, v)
		}
	}
	return nil
}
