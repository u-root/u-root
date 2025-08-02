// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/u-root/u-root/pkg/testutil"
)

func TestBadInvocation(t *testing.T) {
	tests := []struct {
		cmd   cmd
		err   error
		errno int
	}{
		{cmd: cmd{args: []string{}, signal: "KILL"}, err: errNoArgs, errno: 1},
		{cmd: cmd{args: []string{"echo", "WOO"}, signal: "WOO"}, err: os.ErrInvalid, errno: 1},
	}

	for _, v := range tests {
		t.Logf("Run %v", v.cmd)
		if errno, err := v.cmd.run(); !errors.Is(err, v.err) || errno != v.errno {
			t.Errorf("run %v: got (%d, %v), want (%d, %v)", v.cmd, errno, err, v.errno, v.err)
		}
	}
}

func TestRun(t *testing.T) {
	if _, err := exec.LookPath("sleep"); err != nil {
		t.Skipf("Skipping this test as sleep is not in the path")
	}

	tests := []struct {
		cmd cmd
		ok  bool
	}{
		{cmd: cmd{args: []string{"sleep", "4"}, timeout: 2 * time.Minute, signal: "KILL"}, ok: true},
		{cmd: cmd{args: []string{"sleep", "30"}, timeout: time.Second, signal: "KILL"}},
	}
	for _, v := range tests {
		t.Logf("Run %v", v.cmd)
		// Return errors from running sleep are not guaranteed across all kernels.
		// Just see if it succeeded or not.
		if _, err := v.cmd.run(); err == nil != v.ok {
			t.Errorf("run %v: got %v, want %v", v.cmd, err == nil, v.ok)
		}
	}
}

// Test real execution. Why do this if we covered all the code above?
// Because not long ago, someone working on a different u-root
// command got good coverage numbers but never tried
// running the program, and they had completely broken it.
// Tests worked, coverage was good, program was broken for real use.
// It pays to test real operation..
func TestProg(t *testing.T) {
	if _, err := exec.LookPath("sleep"); err != nil {
		t.Skipf("Skipping this test as sleep is not in the path")
	}

	c := testutil.Command(t, "-t=30s", "sleep", "4")
	c.Stdout, c.Stderr = &bytes.Buffer{}, &bytes.Buffer{}
	if err := c.Run(); err != nil {
		t.Errorf("Running -t=30 sleep 4: got %v, want nil", err)
	}

	c = testutil.Command(t, "-t=1s", "sleep", "30")
	c.Stdout, c.Stderr = &bytes.Buffer{}, &bytes.Buffer{}
	if err := c.Run(); err == nil {
		t.Errorf("Running -t=1 sleep 30: got nil, want err")
	}
}

func TestBashExit(t *testing.T) {
	if _, err := exec.LookPath("bash"); err != nil {
		t.Skipf("Skipping test because there is no bash")
	}
	c := testutil.Command(t, "bash", "-c", "exit 20")
	c.Stdout, c.Stderr = &bytes.Buffer{}, &bytes.Buffer{}

	err := c.Run()
	if err == nil {
		t.Fatalf(`Running "bash", "-c", "exit 20": got nil, want err`)
	}

	var errno int
	var e *exec.ExitError
	if errors.As(err, &e) {
		errno = e.ExitCode()
	} else {
		t.Fatalf(`Running "bash", "-c", "exit 20": got %T, want *exec.ExitError`, err)
	}
	if errno != 20 {
		t.Fatalf(`Running "bash", "-c", "exit 20": got %d, want 20`, errno)
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
