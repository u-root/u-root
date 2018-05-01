// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testutil

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/u-root/u-root/pkg/golang"
)

var binary string

func Command(t testing.TB, args ...string) *exec.Cmd {
	// Skip compilation if EXECPATH is set.
	execPath := os.Getenv("EXECPATH")
	if len(execPath) > 0 {
		exe := strings.Split(os.Getenv("EXECPATH"), " ")
		return exec.Command(exe[0], append(exe[1:], args...)...)
	}

	// Should be cached by PrepareMain if os.Executable is going to fail.
	if len(binary) > 0 {
		t.Logf("binary: %v", binary)
		return exec.Command(binary, args...)
	}

	execPath, err := os.Executable()
	if err != nil {
		// PrepareMain should have prevented this case by caching
		// something in `binary`.
		t.Fatal("You must call testutil.PrepareMain() in your test.")
	}

	c := exec.Command(execPath, args...)
	c.Env = append(c.Env, append(os.Environ(), "UROOT_CALL_MAIN=1")...)
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

func PrepareMain() (func(), bool) {
	eraser := func() {}

	if len(os.Getenv("UROOT_CALL_MAIN")) > 0 {
		return eraser, true
	}

	// Cache the executable. Do this here, so that when testing.M.Run()
	// returns, we can remove the executable using the functor returned.
	_, err := os.Executable()
	if err != nil || len(os.Getenv("UROOT_TEST_BUILD")) > 0 {
		// We can't find ourselves? Probably FreeBSD or something. Try to go
		// build the command.
		//
		// This is NOT build-system-independent, and hence the fallback.
		tmpDir, err := ioutil.TempDir("", "uroot-build")
		if err != nil {
			log.Fatal(err)
		}
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		execPath := filepath.Join(tmpDir, "binary")
		// Build the stuff.
		if err := golang.Default().BuildDir(wd, execPath, golang.BuildOpts{}); err != nil {
			log.Fatal(err)
		}

		// Cache dat.
		binary = execPath
		eraser = func() {
			os.RemoveAll(tmpDir)
		}
	}
	return eraser, false
}
