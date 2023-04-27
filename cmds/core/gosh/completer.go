// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !tinygo && !plan9
// +build !tinygo,!plan9

package main

import (
	"io/fs"
	"os"
	"strings"

	"github.com/u-root/prompt"
	"github.com/u-root/prompt/completer"
)

func completerFunc(d prompt.Document) []prompt.Suggest {
	if d.TextBeforeCursor() == "" {
		return []prompt.Suggest{}
	}
	args := strings.Split(d.TextBeforeCursor(), " ")
	w := d.GetWordBeforeCursor()

	// If PIPE is in text before the cursor, returns empty suggestions.
	for i := range args {
		if args[i] == "|" {
			return []prompt.Suggest{}
		}
	}

	if strings.HasPrefix(w, "/") || strings.HasPrefix(w, ".") {
		var filePathCompler = completer.FilePathCompleter{
			IgnoreCase: true,
		}
		return filePathCompler.Complete(d)
	}

	paths := strings.Split(os.ExpandEnv("$PATH"), ":")
	var cmds []fs.DirEntry

	for _, path := range paths {
		entries, err := os.ReadDir(path)
		if err != nil {
			return []prompt.Suggest{}
		}
		cmds = append(cmds, entries...)
	}

	return entryToSuggestion(w, cmds)
}

func entryToSuggestion(w string, e []fs.DirEntry) []prompt.Suggest {
	var s []prompt.Suggest
	for _, entry := range e {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), w) {
			s = append(s, prompt.Suggest{
				Text: entry.Name(),
			})
		}
	}
	return s
}
