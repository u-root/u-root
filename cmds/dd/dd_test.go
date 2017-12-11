// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

// TestDd implements a table-driven test.
func TestDd(t *testing.T) {
	var tests = []struct {
		name    string
		flags   []string
		stdin   string
		stdout  []byte
		count   int64
		compare func(io.Reader, []byte, int64) error
	}{

		{
			name:    "Simple copying from input to output",
			flags:   []string{},
			stdin:   "1: defaults",
			stdout:  []byte("1: defaults"),
			compare: stdoutEqual,
		},
		{
			name:    "Copy from input to output on a non-aligned block size",
			flags:   []string{"bs=8c"},
			stdin:   "2: bs=8c 11b", // len=12 is not multiple of 8
			stdout:  []byte("2: bs=8c 11b"),
			compare: stdoutEqual,
		},
		{
			name:    "case lower change",
			flags:   []string{"bs=8", "conv=lcase"},
			stdin:   "3: Bs=8 11B", // len=11 is not multiple of 8
			stdout:  []byte("3: bs=8 11b"),
			compare: stdoutEqual,
		},
		{
			name:    "case upper change",
			flags:   []string{"bs=8", "conv=ucase"},
			stdin:   "3: Bs=8 11B", // len=11 is not multiple of 8
			stdout:  []byte("3: BS=8 11B"),
			compare: stdoutEqual,
		},
		{
			name:    "Copy from input to output on an aligned block size",
			flags:   []string{"bs=8"},
			stdin:   "hello world.....", // len=16 is a multiple of 8
			stdout:  []byte("hello world....."),
			compare: stdoutEqual,
		},
		{
			name:    "Create a 64KiB zeroed file in 1KiB blocks",
			flags:   []string{"if=/dev/zero", "bs=1K", "count=64"},
			stdin:   "",
			stdout:  []byte("\x00"),
			count:   64 * 1024,
			compare: byteCount,
		},
		{
			name:    "Create a 64KiB zeroed file in 1 byte blocks",
			flags:   []string{"if=/dev/zero", "bs=1", "count=65536"},
			stdin:   "",
			stdout:  []byte("\x00"),
			count:   64 * 1024,
			compare: byteCount,
		},
		{
			name:    "Create a 64KiB zeroed file in one 64KiB block",
			flags:   []string{"if=/dev/zero", "bs=64K", "count=1"},
			stdin:   "",
			stdout:  []byte("\x00"),
			count:   64 * 1024,
			compare: byteCount,
		},
		{
			name:    "Use skip and count",
			flags:   []string{"skip=6", "bs=1", "count=5"},
			stdin:   "hello world.....",
			stdout:  []byte("world"),
			compare: stdoutEqual,
		},
		{
			name:    "Count clamps to end of stream",
			flags:   []string{"bs=2", "skip=3", "count=100000"},
			stdin:   "hello world.....",
			stdout:  []byte("world....."),
			compare: stdoutEqual,
		},
		{
			name:    "1 GiB zeroed file in 1024 1KiB blocks",
			flags:   []string{"bs=1048576", "count=1024", "if=/dev/zero"},
			stdin:   "",
			stdout:  []byte("\x00"),
			count:   1024 * 1024 * 1024,
			compare: byteCount,
		},
	}
	tmpDir, execPath := testutil.CompileInTempDir(t)
	defer os.RemoveAll(tmpDir)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(execPath, tt.flags...)
			cmd.Stdin = strings.NewReader(tt.stdin)
			out, err := cmd.StdoutPipe()
			if err != nil {
				t.Errorf("Test %v exited with error: %v", tt.flags, err)
			}
			if err := cmd.Start(); err != nil {
				t.Errorf("Test %v exited with error: %v", tt.flags, err)
			}
			err = tt.compare(out, tt.stdout, tt.count)
			if err != nil {
				t.Errorf("Test compare function returned: %v", err)
			}
			if err := cmd.Wait(); err != nil {
				t.Errorf("Test %v exited with error: %v", tt.flags, err)
			}
		})
	}
}

// stdoutEqual creates a bufio Reader from io.Reader, then compares a byte at a time input []byte.
// The third argument (int64) is ignored and only exists to make the function signature compatible
// with func byteCount.
// Returns an error if mismatch is found with offset.
func stdoutEqual(i io.Reader, o []byte, _ int64) error {
	var count int64
	b := bufio.NewReader(i)

	for {
		z, err := b.ReadByte()
		if err != nil {
			break
		}
		if o[count] != z {
			return fmt.Errorf("Found mismatch at offset %d, wanted %s, found %s", count, string(o[count]), string(z))
		}
		count++
	}
	return nil
}

// byteCount creates a bufio Reader from io.Reader, then counts the number of sequential bytes
// that match the first byte in the input []byte. If the count matches input n int64, nil error
// is returned. Otherwise an error is returned for a non-matching byte or if the count doesn't
// match.
func byteCount(i io.Reader, o []byte, n int64) error {
	b := bufio.NewReader(i)
	var count int64

	for {
		z, err := b.ReadByte()
		if err != nil {
			break
		}
		if z == o[0] {
			count++
		} else {
			return fmt.Errorf("Found non-matching byte: %v, at offset: %d", o[0], count)
		}
	}

	if count == n {
		return nil
	}
	return fmt.Errorf("Found %d count of %#v bytes, wanted to find %d count", count, o[0], n)
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
