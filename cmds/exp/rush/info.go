// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// info command
//
// Synopsis:
//     info
//
// Description:
//     Print out info about our environment.
//
// Example:
//     $ info
//     Version, goos, etc.
//
// Note:
//
// Bugs:
package main

import (
	"fmt"
	"os"
	"runtime"
)

func init() {
	addBuiltIn("info", infocmd)
}

func infocmd(c *Command) error {
	_, err := fmt.Printf("%s %s %s %q\n", runtime.Version(), runtime.GOOS, runtime.GOARCH, os.Args)
	return err
}
