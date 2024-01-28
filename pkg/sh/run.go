// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sh

import (
	"io"
	"log"
	"os"
	"os/exec"
)

// Run runs a command with stdin, stdout and stderr.
func Run(arg0 string, args ...string) error {
	return RunWithIO(os.Stdin, os.Stdout, os.Stderr, arg0, args...)
}

// RunWithLogs runs a command with stdin, stdout and stderr. This function is
// more verbose than log.Run.
func RunWithLogs(arg0 string, args ...string) error {
	log.Printf("Executing command %q with args %q...", arg0, args)
	err := RunWithIO(os.Stdin, os.Stdout, os.Stderr, arg0, args...)
	if err != nil {
		log.Printf("Command %q with args %q failed: %v", arg0, args, err)
	}
	return err
}

// RunWithIO runs a command with the given input, output and error.
func RunWithIO(in io.Reader, out, err io.Writer, arg0 string, args ...string) error {
	cmd := exec.Command(arg0, args...)
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = err
	return cmd.Run()
}

// RunOrDie runs a commands with stdin, stdout and stderr. If there is a an
// error, it is fatally logged.
func RunOrDie(arg0 string, args ...string) {
	if err := Run(arg0, args...); err != nil {
		log.Fatal(err)
	}
}
