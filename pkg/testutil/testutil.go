// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testutil

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"
)

func Command(t testing.TB, args ...string) *exec.Cmd {
	// Skip compilation if EXECPATH is set.
	execPath := os.Getenv("EXECPATH")
	if len(execPath) > 0 {
		exe := strings.Split(os.Getenv("EXECPATH"), " ")
		return exec.Command(exe[0], append(exe[1:], args...)...)
	}

	execPath, err := os.Executable()
	if err != nil {
		t.Errorf("on strange system: cannot find executable path? %v", err)
	}
	c := exec.Command(execPath, args...)
	c.Env = append(c.Env, "UROOT_CALL_MAIN=1")
	return c
}

func IsExitCode(err error, exitCode int) error {
	if err == nil {
		if exitCode != 0 {
			return fmt.Errorf("got code 0, want %d", exitCode)
		}
		return nil
	}

	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		return fmt.Errorf("encountered error other than ExitError: %#v", err)
	}
	ws, ok := exitErr.Sys().(syscall.WaitStatus)
	if !ok {
		return fmt.Errorf("sys() is not a syscall WaitStatus: %v", err)
	}
	if es := ws.ExitStatus(); es != exitCode {
		return fmt.Errorf("got exit status %d, want %d", es, exitCode)
	}
	return nil
}

func CallMain() bool {
	return len(os.Getenv("UROOT_CALL_MAIN")) > 0
}
