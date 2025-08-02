// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (!tinygo || tinygo.enable) && plan9

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

var (
	defaultKeyFile    = filepath.Join(os.Getenv("home"), "lib/ssh/id_rsa")
	defaultConfigFile = filepath.Join(os.Getenv("home"), "lib/ssh/config")

	consctl *os.File
)

func init() {
	// We have to hold consctl open so we can mess with raw mode.
	var err error
	consctl, err = os.OpenFile("/dev/consctl", os.O_WRONLY, 0o755)
	if err != nil {
		log.Fatal(err)
	}
}

// cleanup turns raw mode back off and closes consctl
func cleanup(in *os.File) {
	consctl.Write([]byte("rawoff"))
	consctl.Close()
}

// raw enters raw mode
func raw(in *os.File) (err error) {
	_, err = consctl.Write([]byte("rawon"))
	return
}

// cooked turns off raw mode
func cooked() (err error) {
	_, err = consctl.Write([]byte("rawoff"))
	return
}

// readPassword prompts the user for a password.
func readPassword(in *os.File, out io.Writer) (string, error) {
	fmt.Fprintf(out, "Password: ")
	raw(in)
	cons, err := os.OpenFile("/dev/cons", os.O_RDWR, 0o755)
	if err != nil {
		return "", err
	}
	defer cons.Close()
	var pw []byte
	for {
		x := make([]byte, 1)
		if _, err := cons.Read(x); err != nil {
			cooked()
			return "", err
		}
		// newline OR carriage return
		if x[0] == '\n' || x[0] == 0x0d {
			break
		}
		pw = append(pw, x[0])
	}
	cooked()
	// output a newline so things look nice
	fmt.Fprintf(out, "\n")
	return string(pw), nil
}

// getSize reads the size of the terminal window.
func getSize(in *os.File) (width, height int, err error) {
	// If we're running vt, there are environment variables to read.
	// If not, we'll just say 80x24
	width, height = 80, 24
	lines := os.Getenv("LINES")
	cols := os.Getenv("COLS")
	if i, err := strconv.Atoi(lines); err == nil {
		height = i
	}
	if i, err := strconv.Atoi(cols); err == nil {
		width = i
	}
	return
}
