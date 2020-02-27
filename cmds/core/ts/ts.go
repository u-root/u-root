// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// ts prepends each line of stdin with a timestamp.
//
// Synopsis:
//     ts
package main

import (
	"io"
	"log"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/ts"
)

func main() {
	flag.Parse()
	if flag.NArg() != 0 {
		log.Fatal("Usage: ts")
	}

	_, err := io.Copy(os.Stdout, ts.New(os.Stdin))
	if err != nil {
		log.Fatal(err)
	}
}
