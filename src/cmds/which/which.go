// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Which locates a command.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

var (
	flags struct {
		a bool
	}
)

const L = os.ModeSymlink

func init() {
	flag.BoolVar(&flags.a, "a", false, "print all matching pathnames of each argument")
}

func which(p string, writer io.Writer, cmds []string) {
	pathArray := strings.Split(p, ":")

	for _, name := range cmds {
		for i := range pathArray {
			f := path.Join(pathArray[i], name)
			if info, err := os.Stat(f); err == nil {
				if m := info.Mode(); !m.IsDir() && m&0111 != 0 && !(m&L == L) {
					writer.Write([]byte(f + "\n"))
				}
			}
		}
	}
}

func main() {
	flag.Parse()

	p := os.Getenv("PATH")
	if len(p) == 0 {
		// The default bin path of uroot is ubin, fallbacking to it.
		fmt.Println("No path variable found! Fallbacking to /ubin")
		p = "/ubin"
	}

	which(p, os.Stdout, flag.Args())
}
