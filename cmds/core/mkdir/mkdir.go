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
	"context"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/core/mkdir"
)

func main() {
	cmd := mkdir.New()
	exitCode, err := cmd.Run(context.Background(), os.Args...)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(exitCode)
}
