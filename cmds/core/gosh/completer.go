// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !tinygo && !plan9
// +build !tinygo,!plan9

package main

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/knz/bubbline/complete"
	"github.com/knz/bubbline/computil"
	"github.com/knz/bubbline/editline"
)

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

func autocomplete(val [][]rune, line, col int) (msg string, completions editline.Completions) {
	word, wstart, wend := computil.FindWord(val, line, col)

	if wstart == 0 && !(strings.HasPrefix(word, ".") || strings.HasPrefix(word, "/")) {
		return commandCompleter(word, col, wstart, wend)
	}

	return filepathCompleter(word, col, wstart, wend)
}

func filepathCompleter(input string, col, wstart, wend int) (msg string, completions editline.Completions) {
	candidates := []string{}

	path, trail := path.Split(input)
	if path == "" {
		pwd, err := os.Getwd()
		if err != nil {
			return msg, completions
		}

		path = pwd
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return msg, completions
	}

	for _, entry := range entries {
		if trail == "" || strings.Contains(entry.Name(), trail) {
			candidates = append(candidates, filepath.Join(path, entry.Name()))
		}
	}

	if len(candidates) != 0 {
		completions = &multiComplete{
			Values:     complete.StringValues("suggestions", candidates),
			moveRight:  wend - col,
			deleteLeft: wend - wstart,
		}
	}

	return msg, completions
}

func commandCompleter(input string, col, wstart, wend int) (msg string, completions editline.Completions) {
	candidates := []string{}

	for _, path := range strings.Split(os.Getenv("PATH"), ":") {
		filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
			if d != nil && !d.IsDir() && strings.HasPrefix(d.Name(), input) {
				candidates = append(candidates, d.Name())
			}

			return nil
		})
	}

	if len(candidates) != 0 {
		completions = &multiComplete{
			Values:     complete.StringValues("suggestions", candidates),
			moveRight:  wend - col,
			deleteLeft: wend - wstart,
		}
	}

	return msg, completions
}
