// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builtin

import (
	"bytes"
	"io"
)

// The Cmd struct is modelled on os/exec.Cmd, with only
// those items need to run a function. Stdin/Stdout/Stderr
// are ReadWriter, to allow functions or returns to be provided with
// data.
type Cmd struct {
	// Path is the BuiltIn's name
	Path string

	// Args holds command line arguments, including the command as Args[0].
	// In typical use, both Path and Args are set by calling Command.
	Args []string

	// Stdin specifies the process's standard input.
	Stdin io.ReadWriter

	// Stdout and Stderr specify the process's standard output and error.
	Stdout io.ReadWriter
	Stderr io.ReadWriter
}

// Command provides a BuiltIn struct with defaults that will behave properly.
// It is recommended that you call this to set up your package builtin,
// but it is not required.
// Stdin has its own bytes.Buffer; Stdout and Stderr share one.
// You can make them separate if you wish.
func Command(path string, args ...string) *Cmd {
	r := &bytes.Buffer{}
	w := &bytes.Buffer{}
	c := &Cmd{Path: path, Args: append([]string{path}, args...), Stdin: r, Stdout: w, Stderr: w}
	return c
}

// Runner is the interface packages must implement.
type Runner interface {
	Run() error
}
