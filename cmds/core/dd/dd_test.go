// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
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
			name:    "512 MiB zeroed file in 1024 1KiB blocks",
			flags:   []string{"bs=524288", "count=1024", "if=/dev/zero"},
			stdin:   "",
			stdout:  []byte("\x00"),
			count:   1024 * 1024 * 512,
			compare: byteCount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := testutil.Command(t, tt.flags...)
			cmd.Stdin = strings.NewReader(tt.stdin)
			out, err := cmd.StdoutPipe()
			if err != nil {
				t.Fatal(err)
			}
			if err := cmd.Start(); err != nil {
				t.Error(err)
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
	var count int64
	buf := make([]byte, 4096)

	for {
		read, err := i.Read(buf)
		if err != nil || read == 0 {
			break
		}
		for z := 0; z < read; z++ {
			if buf[z] == o[0] {
				count++
			} else {
				return fmt.Errorf("Found non-matching byte: %v != %v, at offset: %d",
					buf[z], o[0], count)
			}
		}

		if count > n {
			break
		}
	}

	if count == n {
		return nil
	}
	return fmt.Errorf("Found %d count of %#v bytes, wanted to find %d count", count, o[0], n)
}

// TestFiles uses `if` and `of` arguments.
func TestFiles(t *testing.T) {
	var tests = []struct {
		name     string
		flags    []string
		inFile   []byte
		outFile  []byte
		expected []byte
	}{
		{
			name:     "Simple copying from input to output",
			flags:    []string{},
			inFile:   []byte("1: defaults"),
			expected: []byte("1: defaults"),
		},
		{
			name:     "Copy from input to output on a non-aligned block size",
			flags:    []string{"bs=8c"},
			inFile:   []byte("2: bs=8c 11b"), // len=12 is not multiple of 8
			expected: []byte("2: bs=8c 11b"),
		},
		{
			name:     "Copy from input to output on an aligned block size",
			flags:    []string{"bs=8"},
			inFile:   []byte("hello world....."), // len=16 is a multiple of 8
			expected: []byte("hello world....."),
		},
		{
			name:     "Use skip and count",
			flags:    []string{"skip=6", "bs=1", "count=5"},
			inFile:   []byte("hello world....."),
			expected: []byte("world"),
		},
		{
			name:     "truncate",
			flags:    []string{"bs=1"},
			inFile:   []byte("1234"),
			outFile:  []byte("abcde"),
			expected: []byte("1234"),
		},
		{
			name:     "no truncate",
			flags:    []string{"bs=1", "conv=notrunc"},
			inFile:   []byte("1234"),
			outFile:  []byte("abcde"),
			expected: []byte("1234e"),
		},
		{
			// Fully testing the file is synchronous would require something more.
			name:     "sync",
			flags:    []string{"oflag=sync"},
			inFile:   []byte("x: defaults"),
			expected: []byte("x: defaults"),
		},
		{
			// Fully testing the file is synchronous would require something more.
			name:     "dsync",
			flags:    []string{"oflag=dsync"},
			inFile:   []byte("y: defaults"),
			expected: []byte("y: defaults"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Write in and out file to temporary dir.
			tmpDir, err := ioutil.TempDir("", "dd-test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tmpDir)
			inFile := filepath.Join(tmpDir, "inFile")
			outFile := filepath.Join(tmpDir, "outFile")
			if err := ioutil.WriteFile(inFile, tt.inFile, 0666); err != nil {
				t.Error(err)
			}
			if err := ioutil.WriteFile(outFile, tt.outFile, 0666); err != nil {
				t.Error(err)
			}

			args := append(tt.flags, "if="+inFile, "of="+outFile)
			if err := testutil.Command(t, args...).Run(); err != nil {
				t.Error(err)
			}
			got, err := ioutil.ReadFile(filepath.Join(tmpDir, "outFile"))
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(tt.expected, got) {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

// BenchmarkDd benchmarks the dd command. Each "op" unit is a 1MiB block.
func BenchmarkDd(b *testing.B) {
	const bytesPerOp = 1024 * 1024
	b.SetBytes(bytesPerOp)

	args := []string{
		"if=/dev/zero",
		"of=/dev/null",
		fmt.Sprintf("count=%d", b.N),
		fmt.Sprintf("bs=%d", bytesPerOp),
	}
	b.ResetTimer()
	if err := testutil.Command(b, args...).Run(); err != nil {
		b.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
