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
	cmd := gzip.New()
	if err := cmd.Run(os.Args...); err != nil {
		log.Fatalf("gzip: %v", err)
	}
}
