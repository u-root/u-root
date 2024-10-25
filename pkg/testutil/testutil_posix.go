// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9

package testutil

import (
	"fmt"
	"os/exec"
	"syscall"
)

func exitStatus(exitErr *exec.ExitError) (int, error) {
	ws, ok := exitErr.Sys().(syscall.WaitStatus)
	if !ok {
		return 0, fmt.Errorf("sys() is not a syscall WaitStatus: %w", exitErr)
	}
	return ws.ExitStatus(), nil
}
