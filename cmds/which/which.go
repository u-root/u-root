// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Which locates a command.
//
// Synopsis:
//     which [-a] [COMMAND]...
//
// Options:
//     -a: print all matching pathnames of each argument
package main

import (
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	flags struct {
		allPaths bool
	}
)

func init() {
	flag.BoolVar(&flags.allPaths, "a", false, "print all matching pathnames of each argument")
}

func which(p string, writer io.Writer, cmds []string) {
	pathArray := strings.Split(p, ":")

	// If no matches are found will exit 1, else 0
	exitValue := 1
	for _, name := range cmds {
		for i := range pathArray {
			f := filepath.Join(pathArray[i], name)
			if info, err := os.Stat(f); err == nil {
				// TODO: this test (0111) is not quite right.
				// Consider a file executable only by root (0100)
				// when I'm not root. I can't run it.
				if m := info.Mode(); m&0111 != 0 && !(m&os.ModeType == os.ModeSymlink) {
					exitValue = 0
					writer.Write([]byte(f + "\n"))
					if !flags.allPaths {
						break
					}
				}
			}
		}
	}
	os.Exit(exitValue)
}

func main() {
	flag.Parse()

	p := os.Getenv("PATH")
	if len(p) == 0 {
		// The default bin path of uroot is ubin, fallbacking to it.
		log.Print("No path variable found! Fallbacking to /ubin")
		p = "/ubin"
	}

	which(p, os.Stdout, flag.Args())
}
