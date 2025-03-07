// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin || (linux && !arm && !386 && !mips && !mipsle)

package main

import (
	"os"
)

func main() {
	os.Exit(run(os.Stdout, os.Stderr, os.Args[1:]...))
}
