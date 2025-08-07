// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// base64 - encode and decode base64 from stdin or file to stdout
//
// Synopsis:
//
//	base64 [-d] [FILE]
//
// Description:
//
//	Encode or decode a file to or from base64 encoding.
//	-d   decode data (default is to encode)
//	For stdin, on standard Unix systems, you can use /dev/stdin
package main

import (
	"log"
	"os"

	"github.com/u-root/u-root/pkg/core/base64"
)

func init() {
	log.SetFlags(0)
}

func main() {
	cmd := base64.New()
	err := cmd.Run(os.Args[1:]...)
	if err != nil {
		log.Fatal("base64: ", err)
	}
}
