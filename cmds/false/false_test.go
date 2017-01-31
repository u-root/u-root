// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os/exec"
	"syscall"
	"testing"
)

// Ensure 1 is returned.
func TestFalse(t *testing.T) {
	err := exec.Command("go", "build", "false.go").Run()
	if err != nil {
		t.Fatal("Cannot build false.go:", err)
	}
	out, err := exec.Command("./false").CombinedOutput()
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
