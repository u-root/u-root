// Copyright 2016-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// chmod changes mode bits (e.g. permissions) of a file.
//
// Synopsis:
//
//	chmod MODE FILE...
//
// Desription:
//
//	MODE is a three character octal value or a string like a=rwx
package main

import (
	"context"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/core/chmod"
)

func main() {
	cmd := chmod.New()
	exitCode, err := cmd.Run(context.Background(), os.Args...)
	if err != nil {
		log.Fatalf("%v", err)
	}
	os.Exit(exitCode)
}
