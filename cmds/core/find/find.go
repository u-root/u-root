// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Find finds files. It is similar to the Unix command. It uses REs, not globs,
// for matching.
//
// OPTIONS:
//
//	-d: enable debugging in the find package
//	-mode integer-arg: match against mode, e.g. -mode 0755
//	-type: match against a file type, e.g. -type f will match files
//	-name: glob to match against file
//	-l: long listing. It's not very good, yet, but it's useful enough.
package main

import (
	"log"
	"os"

	"github.com/u-root/u-root/pkg/core/find"
)

func init() {
	log.SetFlags(0)
}

func main() {
	cmd := find.New()
	err := cmd.Run(os.Args[1:]...)
	if err != nil {
		log.Fatal("find: ", err)
	}
}
