// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux
// +build linux

package main

import (
	"fmt"
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
func cleanup() {
	if oldState != nil {
		term.Restore(int(os.Stdin.Fd()), oldState)
	}
}

// raw puts the terminal into raw mode
func raw() (err error) {
	oldState, err = term.MakeRaw(int(os.Stdin.Fd()))
	return
}

// readPassword prompts the user for a password.
func readPassword() (string, error) {
	fmt.Print("Password: ")
	b, err := term.ReadPassword(int(os.Stdin.Fd()))
	return string(b), err
}

// getSize reads the size of the terminal.
func getSize() (width, height int, err error) {
	return term.GetSize(int(os.Stdin.Fd()))
}
