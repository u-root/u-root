// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testutil

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/u-root/gobusybox/src/pkg/golang"
)

// CheckError is a helper function for tests
// It is common to check if an err is expected in the form of errStr, then
// there should be an actual error reported. This is an if and only if condition
// that needs to be verified.
func CheckError(err error, errStr string) error {
	if err != nil && errStr == "" {
		return fmt.Errorf("no error expected, got: \n%w", err)
	} else if err == nil && errStr != "" {
		return fmt.Errorf("error \n%v\nexpected, got nil error", errStr)
	} else if err != nil && err.Error() != errStr {
		return fmt.Errorf("error \n%v\nexpected, got: \n%w", errStr, err)
	}
	return nil
}

// NowLog returns the current time formatted like the standard log package's
// timestamp.
func NowLog() string {
	return time.Now().Format("2006/01/02 15:04:05")
}

var binary string

// Command returns an exec.Cmd appropriate for testing the u-root command.
//
// Command decides which executable to call based on environment variables:
// - EXECPATH="executable args" overrides any other test subject.
// - UROOT_TEST_BUILD=1 will force compiling the u-root command in question.
func Command(t testing.TB, args ...string) *exec.Cmd {
	// If EXECPATH is set, just use that.
	execPath := os.Getenv("EXECPATH")
	if len(execPath) > 0 {
		exe := strings.Split(os.Getenv("EXECPATH"), " ")
		return exec.Command(exe[0], append(exe[1:], args...)...)
	}

	// Should be cached by Run if os.Executable is going to fail.
	if len(binary) > 0 {
		t.Logf("binary: %v", binary)
		return exec.Command(binary, args...)
	}

	execPath, err := os.Executable()
	if err != nil {
		// Run should have prevented this case by caching something in
		// `binary`.
		t.Fatal("You must call testutil.Run() in your TestMain.")
	}

	c := exec.Command(execPath, args...)
	c.Env = append(c.Env, append(os.Environ(), "UROOT_CALL_MAIN=1")...)
	return c
}

// IsExitCode takes err and checks whether it represents the given process exit
// code.
//
// IsExitCode assumes that `err` is the return value of a successful call to
// exec.Cmd.Run/Output/CombinedOutput and hence an *exec.ExitError.
func IsExitCode(err error, exitCode int) error {
	if err == nil {
		if exitCode != 0 {
			return fmt.Errorf("got code 0, want %d", exitCode)
		}
		return nil
	}

	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		return fmt.Errorf("encountered error other than ExitError: %#w", err)
	}
	es, err := exitStatus(exitErr)
	if err != nil {
		return err
	}
	if es != exitCode {
		return fmt.Errorf("got exit status %d, want %d", es, exitCode)
	}
	return nil
}

func run(m *testing.M, mainFn func()) int {
	// UROOT_CALL_MAIN=1 /proc/self/exe should be the same as just running
	// the command we are testing.
	if len(os.Getenv("UROOT_CALL_MAIN")) > 0 {
		mainFn()
		return 0
	}

	// Normally, /proc/self/exe (and equivalents) are used to test u-root
	// commands.
	//
	// Such a symlink isn't available on Plan 9, OS X, or FreeBSD. On these
	// systems, we compile the u-root command in question on the fly
	// instead.
	//
	// Here, we decide whether to compile or not and cache the executable.
	// Do this here, so that when m.Run() returns, we can remove the
	// executable using the functor returned.
	_, err := os.Executable()
	if err != nil || len(os.Getenv("UROOT_TEST_BUILD")) > 0 {
		// We can't find ourselves? Probably FreeBSD or something. Try to go
		// build the command.
		//
		// This is NOT build-system-independent, and hence the fallback.
		tmpDir, err := os.MkdirTemp("", "uroot-build")
		if err != nil {
			log.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		execPath := filepath.Join(tmpDir, "binary")
		// Build the stuff.
		if err := golang.Default().BuildDir(wd, execPath, nil); err != nil {
			log.Fatal(err)
		}

		// Cache dat.
		binary = execPath
	}

	return m.Run()
}

// Run sets up necessary commands to be compiled, if necessary, and calls
// m.Run.
func Run(m *testing.M, mainFn func()) {
	os.Exit(run(m, mainFn))
}

// SkipIfInVMTest skips a test if it's being executed in a u-root test VM.
func SkipIfInVMTest(t *testing.T) {
	if os.Getenv("VMTEST_IN_GUEST") == "1" {
		t.Skipf("Skipping test since we are in a u-root test VM")
	}
}
