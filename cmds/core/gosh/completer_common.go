// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (!tinygo || tinygo.enable) && !plan9 && !goshsmall

package main

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"unicode"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/syntax"
)

func lastStmt(parser *syntax.Parser, line string) *syntax.Stmt {
	var s *syntax.Stmt
	parser.Stmts(strings.NewReader(line), func(stmt *syntax.Stmt) bool {
		s = stmt
		return true
	})
	return s
}

// word returns the word to complete, whether to complete it as a command, and
// the positions it was found at.
func word(stmt *syntax.Stmt, trailingSpaces int) (isCmd bool, pos int, word string) {
	cfg := &expand.Config{
		ReadDir:  ioutil.ReadDir,
		GlobStar: true,
	}
	callExpr, ok := stmt.Cmd.(*syntax.CallExpr)
	if !ok {
		// Not a callexpr, don't expand.
		return false, -1, ""
	}
	if len(callExpr.Args) == 0 {
		// CallExpr with assignment but no args yet, e.g. "FOO=bar".
		// Expand with a command if there was a space, i.e. expand
		// "FOO=bar " but not "FOO=bar"
		if trailingSpaces == 0 {
			return false, -1, ""
		}
		return true, 0, ""
	}

	lastWord := callExpr.Args[len(callExpr.Args)-1]
	args, err := expand.Fields(cfg, lastWord)
	if err != nil || len(args) != 1 {
		return false, -1, ""
	}

	pos = int(lastWord.Pos().Offset())
	if trailingSpaces == 0 {
		return len(callExpr.Args) == 1, pos, args[0]
	}
	return false, pos + trailingSpaces, ""
}

// lastWord returns whether to auto-complete a command-name (true) or file name
// (false), the position of the auto-completion in the line (-1 if nothing to
// complete), and the word to auto-complete.
func lastWord(parser *syntax.Parser, line string) (bool, int, string) {
	withoutSpaces := strings.TrimRightFunc(line, unicode.IsSpace)
	if len(withoutSpaces) == 0 {
		return true, len(line), ""
	}
	stmt := lastStmt(parser, line)
	if stmt == nil {
		return false, -1, ""
	}

	trailingSpaces := len(line) - len(withoutSpaces)
	if stmt.Semicolon.IsValid() {
		if stmt.Background || stmt.Coprocess {
			// Don't autocomplete after "<stmt> &"
			return false, -1, ""
		}

		// We're at "<stmt>;  " with an arbitrary number of
		// spaces after the semicolon, which we want to
		// preserve.
		pos := int(stmt.Semicolon.Offset()) + trailingSpaces + 1
		return true, pos, ""
	}

	// syntax.DebugPrint(os.Stderr, stmt)
	isCmd, pos, word := word(stmt, trailingSpaces)
	if pos == -1 {
		return false, -1, ""
	}
	// When we end with spaces, we're always auto-completing at the end and
	// a completely new word.
	if len(line) != len(withoutSpaces) {
		return isCmd, len(line), ""
	}
	return isCmd, pos, word
}

func addPrefix(s string, t []string) []string {
	var q []string
	for _, tt := range t {
		q = append(q, s+tt)
	}
	return q
}

func autocompleteLiner(parser *syntax.Parser) func(line string) []string {
	return func(line string) []string {
		isCmd, pos, word := lastWord(parser, line)
		if pos == -1 {
			return nil
		}
		prefix := line[:pos]

		if isCmd && !strings.HasPrefix(word, ".") && !strings.HasPrefix(word, "/") {
			return addPrefix(prefix, commandCompleter(word))
		}
		return addPrefix(prefix, filepathCompleter(word))
	}
}

func join(path, entry string) string {
	if path == "" {
		return entry
	}
	if strings.HasSuffix(path, "/") {
		return path + entry
	}
	return path + "/" + entry
}

func filepathCompleter(input string) []string {
	var candidates []string

	path, trail := path.Split(input)
	if path == "" {
		path = "./"
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		if trail == "" || strings.HasPrefix(entry.Name(), trail) {
			candidates = append(candidates, join(path, entry.Name()))
		}
	}
	return candidates
}

func commandCompleter(input string) []string {
	var candidates []string

	for path := range strings.SplitSeq(os.Getenv("PATH"), ":") {
		if err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
			if d != nil && !d.IsDir() && strings.HasPrefix(d.Name(), input) {
				// Is executable?
				if fi, err := d.Info(); err == nil && fi.Mode().Perm()&0o111 != 0 {
					candidates = append(candidates, d.Name())
				}
			}
			return nil
		}); err != nil {
			continue
		}
	}

	return candidates
}
