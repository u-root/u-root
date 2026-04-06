// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Delete files.
//
// Synopsis:
//
//	rm [-Rrvif] FILE...
//
// Options:
//
//	-i: interactive mode
//	-v: verbose mode
//	-R: remove file hierarchies
//	-r: equivalent to -R
//	-f: ignore nonexistent files and never prompt
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/u-root/u-root/pkg/uroot/util"
)

var (
	interactive = flag.Bool("i", false, "Interactive mode.")
	verbose     = flag.Bool("v", false, "Verbose mode.")
	recursive   = flag.Bool("r", false, "equivalent to -R")
	r           = flag.Bool("R", false, "Recursive, remove hierarchies")
	force       = flag.Bool("f", false, "Force, ignore nonexistent files and never prompt")
)

const usage = "rm [-Rrvif] file..."

func rm(stdin io.Reader, files []string) error {
	if len(files) < 1 {
		return fmt.Errorf("%v", usage)
	}
	f := os.Remove
	if *recursive || *r {
		f = os.RemoveAll
	}

	if *force {
		*interactive = false
	}

	workingPath, err := os.Getwd()
	if err != nil {
		return err
	}

	input := bufio.NewReader(stdin)
	for _, file := range files {
		if *interactive {
			fmt.Printf("rm: remove '%v'? ", file)
			answer, err := input.ReadString('\n')
			if err != nil || strings.ToLower(answer)[0] != 'y' {
				continue
			}
		}

		if err := f(file); err != nil {
			if *force && os.IsNotExist(err) {
				continue
			}
			return err
		}

		if *verbose {
			toRemove := file
			if !path.IsAbs(file) {
				toRemove = filepath.Join(workingPath, file)
			}
			fmt.Printf("removed '%v'\n", toRemove)
		}
	}
	return nil
}

func main() {
	flag.Usage = util.Usage(flag.Usage, usage)
	flag.Parse()
	if err := rm(os.Stdin, flag.Args()); err != nil {
		log.Fatal(err)
	}
}
