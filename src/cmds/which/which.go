// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

var (
	flags struct {
		a bool
	}
	cmd = "which [-a] filename ..."
)

func usage() {
	fmt.Fprintln(os.Stderr, "Usage:", cmd)
	flag.PrintDefaults()
	flag.Usage = usage
}

func init() {
	flag.BoolVar(&flags.a, "a", false, "print all matching pathnames of each argument")
}

func which(p string, cmds []string) {
	pathArray := strings.Split(p, ":")

	for _, name := range cmds {
		for i := range pathArray {
			f := path.Join(pathArray[i], name)
			if info, err := os.Stat(f); err == nil {
				if (info.Mode() & 0111) == 0111 {
					fmt.Println(f)
				}
			}
		}
	}
}

func main() {
	flag.Parse()

	p := os.Getenv("PATH")
	if len(p) == 0 {
		fmt.Println("No path variable found!")
		os.Exit(1)
	}

	which(p, flag.Args())
}
