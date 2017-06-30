// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"os/exec"
	"syscall"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

// Ensure 1 is returned.
func TestFalse(t *testing.T) {
	tmpDir, falsePath := testutil.CompileInTempDir(t)
	defer os.RemoveAll(tmpDir)

	out, err := exec.Command(falsePath).CombinedOutput()
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatal("Expected an exit error result")
	}
	retCode := exitErr.Sys().(syscall.WaitStatus).ExitStatus()
	if retCode != 1 {
		t.Fatalf("Expected 1 as the return code; got %v", retCode)
	}
	if len(out) != 0 {
		t.Fatalf("Expected no output; got %#v", string(out))
	}
}
