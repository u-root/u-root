// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"

	"github.com/u-root/u-root/pkg/pogosh"
)

func main() {
	// TODO: not standard
	file := "/dev/stdin"
	if len(os.Args) > 1 {
		file = os.Args[1]
	}

	state := pogosh.DefaultState()
	code, err := state.RunFile(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(code)
}
