// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// shasum computes SHA checksums of files.
//
// Synopsis:
//
//	shasum -a <algorithm> <File Name>
//
// Description:
//
//	shasum computes SHA checksums of files using the specified algorithm.
//	If no files are specified, read from stdin.
//
// Options:
//
//	-a, -algorithm: SHA algorithm, valid args are 1, 256 and 512
package main

import (
	"log"
	"os"

	"github.com/u-root/u-root/pkg/core/shasum"
)

func main() {
	cmd := shasum.New()
	err := cmd.Run(os.Args[1:]...)
	if err != nil {
		log.Fatal("shasum: ", err)
	}
}
