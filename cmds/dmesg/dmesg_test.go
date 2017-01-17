// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os/exec"
	"os/user"
	"testing"
)

// Test reading from the buffer.
// dmesg
func TestDmesg(t *testing.T) {
	_, err := exec.Command("go", "run", "dmesg.go").Output()
	if err != nil {
		t.Fatalf("Error running dmesg: %v", err)
	}
	// FIXME: How can the test verify the output is correct?
}

// Test clearing the buffer.
// dmesg -c
func TestClearDmesg(t *testing.T) {
	// Test requies root priviledges or CAP_SYSLOG capability.
	// FIXME: preferably unit tests do not require root priviledges
	if u, err := user.Current(); err != nil {
		t.Fatal("Cannot get current user", err)
	} else if u.Uid != "0" {
		t.Skipf("Test requires root priviledges (uid == 0), uid = %s", u.Uid)
	}

	// Clear
	out, err := exec.Command("go", "run", "dmesg.go", "-c").Output()
	if err != nil {
		t.Fatalf("Error running dmesg -c: %v", err)
	}

	// Read
	out, err = exec.Command("go", "run", "dmesg.go").Output()
	if err != nil {
		t.Fatalf("Error running dmesg: %v", err)
	}

	// Second run of dmesg.go should be cleared.
	// FIXME: This is actually non-determinstic as the system is free (but
	// unlikely) to write more messages inbetween the syscalls.
	if len(out) > 0 {
		t.Fatalf("The log was not cleared, got %v", out)
	}
}
