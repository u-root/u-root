// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/u-root/u-root/shared/testutil"
)

var tests = []struct {
	flags  []string
	stdin  string
	stdout string
}{
	{
		// Simple copying from input to output.
		flags:  []string{},
		stdin:  "1: defaults",
		stdout: "1: defaults",
	}, {
		// Copy from input to output on a non-aligned block size.
		flags:  []string{"bs=8"},
		stdin:  "2: bs=8 11b", // len=11 is not multiple of 8
		stdout: "2: bs=8 11b",
	}, {
		//  case change
		flags:  []string{"bs=8", "conv=lcase"},
		stdin:  "3: Bs=8 11B", // len=11 is not multiple of 8
		stdout: "3: bs=8 11b",
	}, {
		//  case change
		flags:  []string{"bs=8", "conv=ucase"},
		stdin:  "3: Bs=8 11B", // len=11 is not multiple of 8
		stdout: "3: BS=8 11B",
	}, {
		// Copy from input to output on an aligned block size.
		flags:  []string{"bs=8"},
		stdin:  "hello world.....", // len=16 is a multiple of 8
		stdout: "hello world.....",
	}, {
		// Create a 64KiB zeroed file in 1KiB blocks
		flags:  []string{"if=/dev/zero", "bs=1024", "count=64"},
		stdin:  "",
		stdout: strings.Repeat("\x00", 64*1024),
	}, {
		// Create a 64KiB zeroed file in 1 byte blocks
		flags:  []string{"if=/dev/zero", "bs=1", "count=65536"},
		stdin:  "",
		stdout: strings.Repeat("\x00", 64*1024),
	}, {
		// Create a 64KiB zeroed file in one 64KiB block
		flags:  []string{"if=/dev/zero", "bs=65536", "count=1"},
		stdin:  "",
		stdout: strings.Repeat("\x00", 64*1024),
	}, {
		// Use skip and count.
		flags:  []string{"skip=6", "bs=1", "count=5"},
		stdin:  "hello world.....",
		stdout: "world",
	}, {
		// Count clamps to end of stream.
		flags:  []string{"bs=2", "skip=3", "count=100000"},
		stdin:  "hello world.....",
		stdout: "world.....",
	},
}

// TestDd implements a table-drivent test.
func TestDd(t *testing.T) {
	tmpDir, execPath := testutil.CompileInTempDir(t)
	defer os.RemoveAll(tmpDir)

	for _, tt := range tests {
		cmd := exec.Command(execPath, tt.flags...)
		cmd.Stdin = strings.NewReader(tt.stdin)
		out, err := cmd.Output()
		if err != nil {
			t.Errorf("Exited with error: %v", err)
		}
		if string(out) != tt.stdout {
			t.Errorf("Want:\n%#v\nGot:\n%#v", tt.stdout, string(out))
		}
	}
}
