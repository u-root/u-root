// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pogosh

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func (c *not) exec(s *State) {
	c.cmd.exec(s)
	if s.varExitStatus == 0 {
		s.varExitStatus = 1
	} else {
		s.varExitStatus = 0
	}
}

func (c *and) exec(s *State) {
	c.cmd1.exec(s)
	if s.varExitStatus == 0 {
		c.cmd2.exec(s)
	}
}

func (c *or) exec(s *State) {
	c.cmd1.exec(s)
	if s.varExitStatus != 0 {
		c.cmd2.exec(s)
	}
}

func (c *async) exec(s *State) {
	// TODO
}

func (c *compoundList) exec(s *State) {
	for _, cmd := range c.cmds {
		cmd.exec(s)
	}
}

func (c *pipeline) exec(s *State) {
	// TODO: wire redirects
	for _, cmd := range c.cmds {
		cmd.exec(s)
	}
}

func (c *subshell) exec(s *State) {
	// TODO
	// cmd *command
}

func (c *forClause) exec(s *State) {
	// TODO
	// name     []byte
	// wordlist [][]byte
	// cmd      *command
}

func (c *caseClause) exec(s *State) {
	// TODO
	// word  []byte
	// cases []caseItem

	// func (c *caseItem) exec(s *State) {
	// pattern []byte
	// cmd     *command
	// }
}

func (c *ifClause) exec(s *State) {
	// TODO
	// cmdPred *command
	// cmdThen *command
	// cmdElse *command
}

func (c *whileClause) exec(s *State) {
	// TODO
	// cmdPred *command
	// cmd     *command
}

func (c *function) exec(s *State) {
	// TODO
	// name []byte
	// cmd  *command
}

func searchPath(env string, cmdName string) (string, error) {
	if strings.Contains(cmdName, "/") {
		return cmdName, nil
	}

	for _, prefix := range filepath.SplitList(env) {
		if prefix == "" {
			prefix = "."
		}
		if prefix[len(prefix)-1] != '/' {
			prefix += "/"
		}
		path := prefix + cmdName

		fi, err := os.Stat(path)
		// TODO: could the permission check be more strick?
		if err == nil && fi.Mode().IsRegular() && fi.Mode()&0111 != 0 {
			return path, nil
		}
	}
	return "", fmt.Errorf("Could not find command '%s'", cmdName)
}

// TODO: move to expansion.go file
func wordExpansion(s *State, word string) []string {
	word = tildeExpansion(s, word)
	word = recursiveExpansion(s, word)
	words := fieldSplitting(s, word)
	for i := range words {
		words[i] = pathnameExpansion(s, words[i])
		words[i] = quoteRemoval(s, words[i])
	}
	return words
}

// Contains:
// - Parameter substitution
// - Command substitution
// - Arithmetic substitution
func recursiveExpansion(s *State, word string) string {
	for i := 0; i < len(word); i++ {
		// TODO: check for EOF
		switch {
		case word[i:i+3] == "$((":
			for j := i + 3; j < len(word)-1; j++ {
				if word[j:j+2] == "))" {
					inside := word[i+2 : j]
					inside = recursiveExpansion(s, inside)
					inside = arithmeticSubstitution(s, inside)
					word = word[:i] + inside + word[j+2:]
					i += len(inside) - 1
					break
				}
			}
		case word[i:i+2] == "$(":
			for j := i + 2; j < len(word); j++ {
				if word[j:j+2] == ")" {
					inside := word[i+2 : j]
					inside = recursiveExpansion(s, inside)
					inside = commandSubstitution(s, inside)
					word = word[:i] + inside + word[j+2:]
					i += len(inside) - 1
					break
				}
			}
		case word[i] == '`':
			for j := i + 1; j < len(word); j++ {
				if word[j] == '`' {
					inside := word[i+1 : j]
					inside = recursiveExpansion(s, inside)
					inside = commandSubstitution(s, inside)
					word = word[:i] + inside + word[j+1:]
					i += len(inside) - 1
					break
				}
			}
		case word[i:i+2] == "${":
			for j := i; j < len(word); j++ {

			}
		default:
			word = parameterExpansion(s, word)
		}
	}
	return word
}

func tildeExpansion(s *State, word string) string {
	// TODO: there are more rules than this
	if len(word) > 0 && word[0] == '~' {
		word = s.variables["HOME"].Value + word[1:]
	}
	return word
}

func parameterExpansion(s *State, word string) string {
	return word // TODO
}

func commandSubstitution(s *State, word string) string {
	return word // TODO
}

func arithmeticSubstitution(s *State, word string) string {
	return word // TODO
}

func fieldSplitting(s *State, word string) []string {
	return strings.Split(word, " ") // TODO
}

func pathnameExpansion(s *State, word string) string {
	return word // TODO
}

func quoteRemoval(s *State, word string) string {
	return word // TODO
}

func (c *simpleCommand) exec(s *State) {
	// We are using the lower-level process API to have more control over file
	// descriptors.
	cmd := Cmd{
		ProcAttr: os.ProcAttr{
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		},
	}

	// Clear the exit status variable.
	s.varExitStatus = 0

	// Name
	cmd.name = string(c.name)

	// Arguments
	for _, arg := range c.args {
		cmd.argv = append(cmd.argv, string(arg))
	}

	// Redirects
	for _, redirect := range c.redirects {
		switch string(redirect.ioOp) {
		case "<":
			// Redirect input
			f, err := os.OpenFile(string(redirect.filename), os.O_RDONLY, 0666)
			defer f.Close()
			if err != nil {
				panic(err)
			}
			cmd.Files[0] = f
		case "<&":
			// Duplicating an input file descriptor
			fd, err := strconv.Atoi(string(redirect.filename))
			if err != nil {
				panic(err)
			}
			cmd.Files[0] = os.NewFile(uintptr(fd), "TODO") // TODO: make part of state
			// TODO: closing files with -
		case ">":
			// Redirect output
			f, err := os.OpenFile(string(redirect.filename), os.O_CREATE|os.O_WRONLY, 0666)
			defer f.Close()
			if err != nil {
				panic(err)
			}
			cmd.Files[1] = f
		case ">&":
			// Duplicating an output file descriptor
			fd, err := strconv.Atoi(string(redirect.filename))
			if err != nil {
				panic(err)
			}
			cmd.Files[1] = os.NewFile(uintptr(fd), "TODO") // TODO: make part of state
			// TODO: closing files with -
		case ">>":
			// Appending redirected output
			f, err := os.OpenFile(string(redirect.filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
			defer f.Close()
			if err != nil {
				panic(err)
			}
			cmd.Files[1] = f
		case "<>":
			// Open file descriptor for reading and writing
			f, err := os.OpenFile(string(redirect.filename), os.O_CREATE|os.O_RDWR, 0666)
			defer f.Close()
			if err != nil {
				panic(err)
			}
			cmd.Files[0] = f
		case ">|":
			// TODO
		}
	}

	// First, resolve and execute builtins
	if builtin, ok := s.Builtins[cmd.name]; ok {
		builtin(s, &cmd)
		return
	}

	// Second, resolve PATH
	var err error
	cmd.name, err = searchPath(os.Getenv("PATH"), string(c.name))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err) // TODO: better error handling
		s.varExitStatus = 127
		return
	}

	// Finally, execute the command
	proc, err := os.StartProcess(cmd.name, cmd.argv, &cmd.ProcAttr)
	if err != nil {
		// TODO: check other error types
		fmt.Fprintf(os.Stderr, "Cannot find command %s, error: %s\n", cmd.name, err)
		s.varExitStatus = 127
		return
	}

	processState, err := proc.Wait()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running command %s, error: %s\n", cmd.name, err)
		s.varExitStatus = 127
		return
	}

	// TODO: syscall.WaitStatus not same on all systems
	s.varExitStatus = processState.Sys().(syscall.WaitStatus).ExitStatus()
}
