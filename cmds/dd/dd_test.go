// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
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
	}, {
		// 1 GiB zeroed file in 1024 1KiB blocks.
		flags:  []string{"bs=1048576", "count=1024", "if=/dev/zero"},
		stdin:  "",
		stdout: strings.Repeat("\x00", 1024*1024*1024),
	},
}

// TestDd implements a table-driven test.
func TestDd(t *testing.T) {
	tmpDir, execPath := testutil.CompileInTempDir(t)
	defer os.RemoveAll(tmpDir)

	for _, tt := range tests {
		cmd := exec.Command(execPath, tt.flags...)
		cmd.Stdin = strings.NewReader(tt.stdin)
		out, err := cmd.Output()
		if err != nil {
			t.Errorf("Test %v exited with error: %v", tt.flags, err)
		}
		if string(out) != tt.stdout {
			t.Errorf("Want:\n%#v\nGot:\n%#v", tt.stdout, string(out))
		}
	}
}

// BenchmarkDd benchmarks the dd command. Each "op" unit is a 1MiB block.
func BenchmarkDd(b *testing.B) {
	tmpDir, execPath := testutil.CompileInTempDir(b)
	defer os.RemoveAll(tmpDir)

	const bytesPerOp = 1024 * 1024
	b.SetBytes(bytesPerOp)
	args := []string{
		"if=/dev/zero",
		"of=/dev/null",
		fmt.Sprintf("count=%d", b.N),
		fmt.Sprintf("bs=%d", bytesPerOp),
	}
	b.ResetTimer()
	if err := exec.Command(execPath, args...).Run(); err != nil {
		b.Fatal(err)
	}
}
