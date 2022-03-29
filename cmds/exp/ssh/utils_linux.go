//go:build linux
// +build linux

package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

var oldState *term.State

// cleanup returns the terminal to its original state
func cleanup() {
	term.Restore(int(os.Stdin.Fd()), oldState)
}

// raw puts the terminal into raw mode
func raw() (err error) {
	oldState, err = term.MakeRaw(int(os.Stdin.Fd()))
	return
}

// ReadPassword prompts the user for a password.
func ReadPassword() (string, error) {
	fmt.Print("Password: ")
	b, err := term.ReadPassword(int(os.Stdin.Fd()))
	return string(b), err
}

// GetSize reads the size of the terminal.
func GetSize() (width, height int, err error) {
	return term.GetSize(int(os.Stdin.Fd()))
}
