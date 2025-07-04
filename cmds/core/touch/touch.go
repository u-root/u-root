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
	"context"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/core/touch"
)

func main() {
	cmd := touch.New()
	exitCode, err := cmd.Run(context.Background(), os.Args...)
	if err != nil {
		log.Fatalf("touch: %v", err)
	}
	os.Exit(exitCode)
}
