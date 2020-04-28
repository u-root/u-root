// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pogosh

import "os"

// TODO: rename file to types.go ?

// Cmd is a builtin command to run.
type Cmd struct {
	os.ProcAttr
	name string
	argv []string
}

// Builtin is the function used to run a builtin.
type Builtin = func(*State, *Cmd)

// State holds data on the current interpreter execution.
type State struct {
	IsInteractive bool

	Builtins  map[string]func(*State, *Cmd)
	Aliases   map[string]string
	variables map[string]Var

	// Special variables
	varExitStatus int // $?

	Overrides Overrides

	parent *State
}

// Var holds information on shell variable.
type Var struct {
	Value string
}

// Overrides change the behaviour of builtins.
type Overrides struct {
	Chdir   func(dir string) error
	Environ func() []string
	Exit    func(code int)
}

// DefaultState creates a new default state.
func DefaultState() State {
	return State{
		Builtins: DefaultBuiltins(),
		Aliases:  map[string]string{},

		Overrides: DefaultOverrides(),
	}
}

// DefaultOverrides creates a new default overrides.
func DefaultOverrides() Overrides {
	return Overrides{
		Chdir:   os.Chdir,
		Environ: os.Environ,
		Exit:    func(code int) { panic(exitError{code}) },
	}
}

// Errors propagated using panic
type exitError struct{ code int }
