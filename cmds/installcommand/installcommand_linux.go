// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"os/exec"
	"syscall"
)

func exitWithStatus(err *exec.ExitError) {
	os.Exit(err.Sys().(syscall.WaitStatus).ExitStatus())
}
