// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Print ansi escape sequences.
//
// Synopsis:
//
//	ansi COMMAND
//
// Options:
//
//	COMMAND must be one of:
//	    - clear: clear the screen and reset the cursor position
//
// Author:
//
//	Manoel Vilela <manoel_vilela@engineer.com>
package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

// using ansi escape codes /033 => escape code
// The "\033[1;1H" part moves the cursor to position (1,1)
// "\033[2J" part clears the screen.
// if you wants add more escape codes, append on map below
// arg:escape_code
var commands = map[string]string{
	"clear": "\033[1;1H\033[2J",
}

func ansi(w io.Writer, args []string) error {
	for _, arg := range args {
		c, ok := commands[arg]
		if !ok {
			return fmt.Errorf("command ANSI %q don't exists", arg)
		}
		fmt.Fprint(w, c)
	}
	return nil
}

func main() {
	if err := ansi(os.Stdout, os.Args[1:]); err != nil {
		log.Fatalf("%v", err)
	}
}
