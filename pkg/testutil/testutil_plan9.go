// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testutil

import (
	"fmt"
	"os/exec"
	"syscall"
)

func exitStatus(exitErr *exec.ExitError) (int, error) {
	ws, ok := exitErr.Sys().(syscall.Waitmsg)
	if !ok {
		return 0, fmt.Errorf("sys() is not a syscall Waitmsg: %v", exitErr)
	}
	return ws.ExitStatus(), nil
}
