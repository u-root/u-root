// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// gzip compresses files using gzip compression.
package main

import (
	"log"
	"os"

	"github.com/u-root/u-root/pkg/core/gzip"
)

func main() {
	program, args := os.Args[0], os.Args[1:]
	cmd := gzip.New(program)
	if err := cmd.Run(args...); err != nil {
		log.Fatalf("%s: %v", program, err)
	}
}
