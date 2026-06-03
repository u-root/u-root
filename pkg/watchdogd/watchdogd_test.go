// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package watchdogd

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDaemonStartStop(t *testing.T) {
	// Create a temp file to act as the watchdog device.
	tmpFile, err := os.CreateTemp("", "watchdog-test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close() // Close it so Daemon can open it.

	opts := &DaemonOpts{
		Dev:        tmpFile.Name(),
		Timeout:    timeoutIgnore,
		PreTimeout: timeoutIgnore,
		KeepAlive:  10 * time.Millisecond, // Fast petting for test
		UDS:        filepath.Join(t.TempDir(), "watchdogd.sock"),
	}

	d := NewDaemon(opts)

	// Test Arm
	if r := d.ArmWatchdog(); r != OpResultOk {
		t.Fatalf("ArmWatchdog failed: %c", r)
	}

	if d.CurrentWd == nil {
		t.Fatal("Expected CurrentWd to be non-nil after Arm")
	}
	if !d.PettingOn {
		t.Fatal("Expected PettingOn to be true after Arm (due to auto-start)")
	}

	// Let it pet a few times
	time.Sleep(50 * time.Millisecond)

	// Test Disarm
	if r := d.DisarmWatchdog(); r != OpResultOk {
		t.Fatalf("DisarmWatchdog failed: %c", r)
	}

	if d.CurrentWd != nil {
		t.Fatal("Expected CurrentWd to be nil after Disarm")
	}
	if d.PettingOn {
		t.Fatal("Expected PettingOn to be false after Disarm")
	}

	// Read the temp file to verify writes
	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}

	t.Logf("Watchdog file content: %q", string(data))

	// It should contain some '1's and end with 'V' (magic close)
	if len(data) == 0 {
		t.Fatal("Expected data to be written to watchdog file")
	}
	if data[len(data)-1] != 'V' {
		t.Fatalf("Expected last char to be 'V' (magic close), got %q", data[len(data)-1])
	}
	for i := 0; i < len(data)-1; i++ {
		if data[i] != '1' {
			t.Fatalf("Expected char at %d to be '1' (pet), got %q", i, data[i])
		}
	}
}
