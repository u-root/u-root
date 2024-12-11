// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (!tinygo || tinygo.enable) && (!plan9 || !windows)

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/term"
)

var (
	defaultKeyFile    = filepath.Join(os.Getenv("HOME"), ".ssh/id_rsa")
	defaultConfigFile = filepath.Join(os.Getenv("HOME"), ".ssh/config")

	oldState *term.State
)

// cleanup returns the terminal to its original state
func cleanup(in *os.File) {
	if oldState != nil {
		term.Restore(int(in.Fd()), oldState)
	}
}

// raw puts the terminal into raw mode
func raw(in *os.File) (err error) {
	oldState, err = term.MakeRaw(int(in.Fd()))
	return
}

// readPassword prompts the user for a password.
func readPassword(in *os.File, out io.Writer) (string, error) {
	fmt.Fprintf(out, "Password: ")
	b, err := term.ReadPassword(int(in.Fd()))
	return string(b), err
}

// getSize reads the size of the terminal.
func getSize(in *os.File) (width, height int, err error) {
	return term.GetSize(int(in.Fd()))
}
