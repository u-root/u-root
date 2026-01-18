// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (!tinygo || tinygo.enable) && !plan9 && !goshsmall && !goshliner

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/knz/bubbline/complete"
	"github.com/knz/bubbline/computil"
	"github.com/knz/bubbline/editline"

	"github.com/peterh/liner"
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
	var runErr error

	// Set up liner
	line := liner.NewLiner()
	defer line.Close()

	// Cache commands from PATH
	pathCommands := getPathCommands()

	// Set up completion with both commands and files
	line.SetCompleter(func(input string) []string {
		// Find the position where we should insert completions
		words := strings.Fields(input)
		if len(words) == 0 {
			return nil
		}

		// Get the last word being typed
		lastWord := words[len(words)-1]

		// Calculate prefix (everything before the last word)
		lastWordStart := strings.LastIndex(input, lastWord)
		prefix := ""
		if lastWordStart > 0 {
			prefix = input[:lastWordStart]
		}

		var matches []string

		// If it's the first word and doesn't contain /, complete with commands from PATH
		if len(words) == 1 && !strings.Contains(lastWord, "/") {
			for _, cmd := range pathCommands {
				if strings.HasPrefix(cmd, lastWord) {
					matches = append(matches, prefix+cmd)
				}
			}
		}

		// Always try file completion for paths
		fileMatches := completeFiles(lastWord)
		for _, match := range fileMatches {
			matches = append(matches, prefix+match)
		}

		return matches
	})

	// Enable Ctrl-C handling
	line.SetCtrlCAborts(true)

	// Set up signal handling
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func(ch chan os.Signal) {
		for {
			<-ch
		}
	}(ch)

	prompt := "$ "
	var inputBuffer strings.Builder

	for {
		// Read a line with the current prompt
		input, err := line.Prompt(prompt)
		if err != nil {
			if err == liner.ErrPromptAborted {
				// Ctrl-C was pressed
				inputBuffer.Reset()
				prompt = "$ "
				continue
			}
			if err == io.EOF {
				// Ctrl-D was pressed
				return nil
			}
			return err
		}

		// Add to history if non-empty
		if input != "" {
			line.AppendHistory(input)
		}

		// Accumulate multi-line input
		if inputBuffer.Len() > 0 {
			inputBuffer.WriteString("\n")
		}
		inputBuffer.WriteString(input)

		// Try to parse accumulated input
		r := strings.NewReader(inputBuffer.String() + "\n")

		incomplete := false
		fn := func(stmts []*syntax.Stmt) bool {
			if parser.Incomplete() {
				incomplete = true
				prompt = "> "
				return true
			}

			for _, stmt := range stmts {
				runErr = runner.Run(context.Background(), stmt)
				if runner.Exited() {
					return false
				}
			}

			prompt = "$ "
			return true
		}

		if err := parser.Interactive(r, fn); err != nil {
			fmt.Fprintf(stderr, "error: %s\n", err.Error())
			inputBuffer.Reset()
			prompt = "$ "
			continue
		}

		// If incomplete, keep accumulating
		if incomplete {
			continue
		}

		// Clear buffer after successful parse
		inputBuffer.Reset()

		if runErr != nil {
			fmt.Fprintf(stderr, "error: %s\n", runErr.Error())
			runErr = nil
		}

		if runner.Exited() {
			return nil
		}
	}
}

// getPathCommands returns a list of all executable commands in PATH
func getPathCommands() []string {
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return nil
	}

	paths := filepath.SplitList(pathEnv)
	commandSet := make(map[string]bool)

	for _, dir := range paths {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			// Check if executable
			info, err := entry.Info()
			if err != nil {
				continue
			}

			mode := info.Mode()
			if mode&0111 != 0 { // Has any execute bit set
				commandSet[entry.Name()] = true
			}
		}
	}

	// Convert set to sorted slice
	commands := make([]string, 0, len(commandSet))
	for cmd := range commandSet {
		commands = append(commands, cmd)
	}

	return commands
}

// completeFiles returns file/directory completions for the given prefix
func completeFiles(prefix string) []string {
	// Handle empty prefix - complete in current directory
	if prefix == "" {
		prefix = "./"
	}

	// Expand home directory
	if strings.HasPrefix(prefix, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			prefix = filepath.Join(home, prefix[2:])
		}
	}

	// Get directory and file prefix
	dir := filepath.Dir(prefix)
	filePrefix := filepath.Base(prefix)

	// If prefix ends with /, we want to complete in that directory
	if strings.HasSuffix(prefix, "/") {
		dir = prefix
		filePrefix = ""
	}

	// Read directory entries
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var matches []string
	for _, entry := range entries {
		name := entry.Name()

		// Skip hidden files unless explicitly requested
		if !strings.HasPrefix(filePrefix, ".") && strings.HasPrefix(name, ".") {
			continue
		}

		// Check if name matches prefix
		if strings.HasPrefix(name, filePrefix) {
			// Build full path
			var fullPath string
			if dir == "." {
				fullPath = name
			} else {
				fullPath = filepath.Join(dir, name)
			}

			// Add trailing slash for directories
			if entry.IsDir() {
				fullPath += "/"
			}

			matches = append(matches, fullPath)
		}
	}

	return matches
}

//```
//
//## Key Changes
//
//1. **`getPathCommands()`** - Scans all directories in `$PATH` and finds executable files
//2. **Caches commands** - Only scans PATH once at startup for performance
//3. **Checks execute permissions** - Only includes files with execute bit set
//4. **Deduplicates** - Uses a map to handle commands that appear in multiple PATH directories
//
//## Features
//
//- Completes commands from `$PATH` for the first argument
//- Completes files/directories for subsequent arguments or paths with `/`
//- Fast completion (PATH scanned only once)
//- Works with any command in your PATH
//
//Example:
//```
//$ gre<TAB>  → completes to grep, grep-changelog, etc.
//$ /usr/bi<TAB> → completes files in /usr/bin/
//$ ls Do<TAB> → completes to Documents/, Downloads/, etc.
