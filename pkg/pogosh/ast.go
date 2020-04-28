// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pogosh

type command interface {
	exec(*State)
}

// ! cmd
type not struct {
	cmd command
}

// cmd1 && cmd2
type and struct {
	cmd1 command
	cmd2 command
}

// cmd1 || cmd2
type or struct {
	cmd1 command
	cmd2 command
}

// cmd &
type async struct {
	cmd command
}

// cmds[0]; cmds[1]; ... cmds[n-1];
type compoundList struct {
	cmds []command
}

// cmds[0] | cmds[1] | ... | cmds[n-1]
type pipeline struct {
	cmds []command
}

// ( cmd )
type subshell struct {
	cmd command
}

// for name in wordlist; do cmd; done
type forClause struct {
	name     []byte
	wordlist [][]byte
	cmd      command
}

// case word in cases esac
type caseClause struct {
	word  []byte
	cases []caseItem
}

// pattern ) cmd
type caseItem struct {
	pattern []byte
	cmd     command
}

// if cmdPred then cmdThen else cmdElse done
type ifClause struct {
	cmdPred command
	cmdThen command
	cmdElse command
}

// while cmdPred; do cmd; done
type whileClause struct {
	cmdPred command
	cmd     command
}

// name() cmd
type function struct {
	name []byte
	cmd  command
}

// name args... redirects...
type simpleCommand struct {
	name      []byte
	args      [][]byte
	redirects []redirect
}

// ioOp filename
type redirect struct {
	ioOp     []byte
	filename []byte
}
