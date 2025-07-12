// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// mkdir makes a new directory.
//
// Synopsis:
//
//	mkdir [-m mode] [-v] [-p] DIRECTORY...
//
// Options:
//
//	-m: directory mode (ex: 755)
//	-v: print each directory as it is made
//	-p: make all needed directories in the path
package main

import (
	"log"
	"os"

	"github.com/u-root/u-root/pkg/core/mkdir"
)

func init() {
	log.SetFlags(0)
}

func main() {
	cmd := mkdir.New()
	err := cmd.Run(os.Args[1:]...)
	if err != nil {
		log.Fatal("mkdir: ", err)
	}
}
