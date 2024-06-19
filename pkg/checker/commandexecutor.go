// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package checker

import (
	"log"
	"os"
	"os/exec"
)

// DefaultShell is used by EmergencyShell
var DefaultShell = "gosh"

func runCmd(prog string, args ...string) error {
	cmd := exec.Command(prog, args...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}

// CommandExecutor returns a check that runs the provided command and arguments.
func CommandExecutor(prog string, args ...string) Checker {
	return func() error {
		return runCmd(prog, args...)
	}
}

// CommandExecutorRemediation is like CommandExecutor, but returns a Remediator.
func CommandExecutorRemediation(prog string, args ...string) Remediator {
	return func() error {
		return runCmd(prog, args...)
	}
}

// EmergencyShell is a remediation that prints the given banner, and then calls
// an emergency shell.
func EmergencyShell(banner string) Remediator {
	return func() error {
		log.Print(green("Running emergency shell: %s", DefaultShell))
		if banner != "" {
			log.Print(banner)
		}
		return CommandExecutorRemediation(DefaultShell)()
	}
}
