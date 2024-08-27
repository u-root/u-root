// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build tinygo

package main

import (
	"fmt"
	"io"
	"os/exec"
)

func printProcessState(w io.Writer, c *exec.Cmd) {
	if c.ProcessState == nil {
		return
	}
	fmt.Fprintf(w, "%v", c.ProcessState)
}
