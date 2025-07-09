// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// touch changes file access and modification times.
//
// Synopsis:
//
//	touch [-amc] [-d datetime] file...
//
// Description:
//
//	If a file does not exist, it will be created unless -c is specified.
//
// Options:
//
//	-a: change only the access time
//	-m: change only the modification time
//	-c: do not create any file if it does not exist
//	-d: use specified time instead of current time (RFC3339 format)
package main

import (
	"log"
	"os"

	"github.com/u-root/u-root/pkg/core/touch"
)

func init() {
	log.SetFlags(0)
}

func main() {
	cmd := touch.New()
	err := cmd.Run(os.Args[1:]...)
	if err != nil {
		log.Fatal("touch: ", err)
	}
}
