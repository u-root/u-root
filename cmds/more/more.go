// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// More pages through files without any terminal trickery.
//
// Synopsis:
//     more [OPTIONS] FILE
//
// Description:
//     Admittedly, this does not follow the conventions of GNU more. Instead,
//     it is built with the goal of not relying on any special ttys, ioctls or
//     special ANSI escapes. This is ideal when your terminal is already
//     borked. For bells and whistles, look at less.
//
// Options:
//     --lines NUMBER: screen size in number of lines
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	flag "github.com/spf13/pflag"
)

var lines = flag.Int("lines", 40, "screen size in number of lines")

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatal("more can only take one file")
	}
	if *lines <= 0 {
		log.Fatal("lines must be positive")
	}

	f, err := os.Open(flag.Args()[0])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for i := 0; scanner.Scan(); i++ {
		if (i+1)%*lines == 0 {
			fmt.Print(scanner.Text())
			c := make([]byte, 1)
			// We expect the OS to echo the newline character.
			if _, err := os.Stdin.Read(c); err != nil {
				return
			}
		} else {
			fmt.Println(scanner.Text())
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
