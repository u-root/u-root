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
	"log"
	"os"

	"github.com/u-root/u-root/pkg/core/rm"
)

func init() {
	log.SetFlags(0)
}

func main() {
	cmd := rm.New()
	err := cmd.Run(os.Args[1:]...)
	if err != nil {
		log.Fatal("rm: ", err)
	}
}
