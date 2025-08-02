// Copyright 2013-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// xargs reads space, tab, newline and end-of-file delimited strings from the
// standard input and executes utility with the strings as arguments.
//
// Synopsis:
//
//	xargs [OPTIONS] [COMMAND [ARGS]...]
//
// Options:
//
//	-n: max number of arguments per command
//	-t: enable trace mode, each command is written to stderr
//	-p: the user is asked whether to execute utility at each invocation
//	-0: use a null byte as the input argument delimiter
package main

import (
	"log"
	"os"

	"github.com/u-root/u-root/pkg/core/xargs"
)

func init() {
	log.SetFlags(0)
}

func main() {
	cmd := xargs.New()
	err := cmd.Run(os.Args[1:]...)
	if err != nil {
		log.Fatal("xargs: ", err)
	}
}
