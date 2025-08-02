// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// cat concatenates files and prints them to stdout.
//
// Synopsis:
//
//	cat [-u] [FILES]...
//
// Description:
//
//	If no files are specified, read from stdin.
//
// Options:
//
//	-u: ignored flag
package main

import (
	"log"
	"os"

	"github.com/u-root/u-root/pkg/core/cat"
)

func init() {
	log.SetFlags(0)
}

func main() {
	cmd := cat.New()
	err := cmd.Run(os.Args[1:]...)
	if err != nil {
		log.Fatal("cat: ", err)
	}
}
