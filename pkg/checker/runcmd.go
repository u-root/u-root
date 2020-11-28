// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package checker

import (
	"os"
	"os/exec"
)

// RunCmd executes a command and returns its error status. It can be used for
// checks and remediations.
func RunCmd(prog string, args ...string) error {
	cmd := exec.Command(prog, args...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}
