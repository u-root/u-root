// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package forth

import (
	"errors"
	"fmt"

	"github.com/u-root/u-root/pkg/builtin"
)

// Cmd supports using the forth package as a builtin.
type Cmd struct {
	*builtin.Cmd
	f     Forth
	cells []Cell
}

// Command returns a new Command, with a forth interpreter attached.
func Command(path string, args ...string) *Cmd {
	var cells []Cell
	for _, s := range args {
		cells = append(cells, Cell(s))
	}

	return &Cmd{
		Cmd:   builtin.Command(path, args...),
		f:     New(),
		cells: cells,
	}
}

var _ builtin.Runner = (*Cmd)(nil)

// Run implements builtin.Run. For forth, it takes the args
// and runs them, returning results c.Stdout (which, by default,
// is a bytes.Buffer)
func (c *Cmd) Run() error {
	evalErr := Eval(c.f, c.cells...)
	s := c.f.Stack()
	Debug("stack:%v", s)
	if len(s) == 0 {
		return errors.Join(evalErr, ErrEmptyStack)
	}

	_, err := fmt.Fprintf(c.Stdout, "%s", s[0])

	return errors.Join(evalErr, err)
}
