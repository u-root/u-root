// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// gzcat cats gzip compressed files.
package main

import (
	"log"
	"os"

	"github.com/u-root/u-root/pkg/core/gzip"
)

func main() {
	cmd := gzip.NewGzcat()
	if err := cmd.Run(os.Args[1:]...); err != nil {
		log.Fatalf("gzcat: %v", err)
	}
}
