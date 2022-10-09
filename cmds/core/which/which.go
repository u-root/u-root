// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Which locates a command.
//
// Synopsis:
//
//	which [-a] [COMMAND]...
//
// Options:
//
//	-a: print all matching pathnames of each argument
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	allPaths = flag.Bool("a", false, "print all matching pathnames of each argument")
	verbose  = flag.Bool("v", false, "verbose output")
)

func which(writer io.Writer, paths []string, cmds []string, allPaths bool) error {
	var foundOne bool
	for _, name := range cmds {
		for _, p := range paths {
			f := filepath.Join(p, name)
			if !canExecute(f) {
				continue
			}

			foundOne = true
			if _, err := writer.Write([]byte(f + "\n")); err != nil {
				return err
			}
			if !allPaths {
				break
			}
		}
	}

	if !foundOne {
		return fmt.Errorf("no suitable executable found")
	}
	return nil
}

func main() {
	flag.Parse()

	p := os.Getenv("PATH")
	if len(p) == 0 {
		// The default bin path of uroot is ubin, fallbacking to it.
		log.Print("No path variable found! Fallbacking to /ubin")
		p = "/ubin"
	}

	if err := which(os.Stdout, strings.Split(p, ":"), flag.Args(), *allPaths); err != nil {
		if *verbose {
			log.Printf("Error: %v", err)
		}
		os.Exit(1)
	}
}
