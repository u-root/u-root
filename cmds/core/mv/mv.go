// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// mv renames files and directories.
//
// Synopsis:
//
//	mv SOURCE [-u] TARGET
//	mv SOURCE... [-u] DIRECTORY
//
// Author:
//
//	Beletti (rhiguita@gmail.com)
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/u-root/u-root/pkg/core/mv"
)

func main() {
	cmd := mv.New()
	exitCode, err := cmd.Run(context.Background(), os.Args[1:]...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mv: %v\n", err)
	}
	os.Exit(exitCode)
}
