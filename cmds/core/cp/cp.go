// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// cp copies files.
//
// Synopsis:
//
//	cp [-rRfivwP] FROM... TO
//
// Options:
//
//	-w n: number of worker goroutines
//	-R: copy file hierarchies
//	-r: alias to -R recursive mode
//	-i: prompt about overwriting file
//	-f: force overwrite files
//	-v: verbose copy mode
//	-P: don't follow symlinks
package main

import (
	"context"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/core/cp"
)

func main() {
	cmd := cp.New()
	exitCode, err := cmd.Run(context.Background(), os.Args[1:]...)
	if err != nil {
		log.Fatalf("%v", err)
	}
	os.Exit(exitCode)
}
