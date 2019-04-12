// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sh

import (
	"log"
	"os"
	"os/exec"
)

// Run runs a command with stdin, stdout and stderr.
func Run(arg0 string, args ...string) error {
	cmd := exec.Command(arg0, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunOrDie runs a commands with stdin, stdout and stderr. If there is a an
// error, it is fatally logged.
func RunOrDie(arg0 string, args ...string) {
	if err := Run(arg0, args...); err != nil {
		log.Fatal(err)
	}
}
