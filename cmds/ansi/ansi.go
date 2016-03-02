// Copyright 2015 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// By Manoel Vilela <manoel_vilela@engineer.com>

package main

import (
	"errors"
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
		_, exists := commands[arg]
		if exists {
			fmt.Fprintf(w, commands[arg])
		} else {
			return errors.New(fmt.Sprintf("Command ANSI '%v' don't exists", arg))
		}
	}
	return nil
}

func main() {
	if err := ansi(os.Stdout, os.Args[1:]); err != nil {
		log.Fatalf("%v", err)
	}

}
