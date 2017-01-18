// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os/exec"
	"testing"
)

// Test reading from the buffer.
// dmesg
func TestDmesg(t *testing.T) {
	out, err := exec.Command("go", "run", "dmesg.go").Output()
	if err != nil {
		t.Fatalf("Error running dmesg: %v", err)
	}
	// Test passes if anything is read.
	if len(out) == 0 {
		t.Fatalf("Nothing read from dmesg")
	}
}
