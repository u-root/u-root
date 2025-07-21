// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

type ws struct {
	io.Writer
}

// Write implements os.Write
func (w *ws) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// Seek implements a limited form of io.Seek:
// whence is ignored. That is ok, because
// all the tests start from the start of the file.
func (w *ws) Seek(offset int64, _ int) (int64, error) {
	for range offset {
		if _, err := w.Write([]byte{1}[:]); err != nil {
			return -1, err
		}
	}
	return offset, nil
}

func TestRead(t *testing.T) {
	tests := []struct {
		name      string
		offset    int64
		maxRead   int64
		expected  []byte
		wantError bool
	}{
		{
			name:     "read one byte from offset 0",
			offset:   0,
			maxRead:  10,
			expected: []byte("A"),
		},
		{
			name:     "read one byte from offset 3",
			offset:   3,
			maxRead:  10,
			expected: []byte("D"),
		},
		{
			name:      "read out of bounds",
			offset:    11,
			maxRead:   10,
			expected:  []byte{},
			wantError: true,
		},
		{
			name:      "Read EOF",
			offset:    0,
			maxRead:   0,
			expected:  []byte{},
			wantError: true,
		},
	}

	p, cleanup := setupDatafile(t, "datafile")
	defer cleanup()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := make([]byte, len(tt.expected))

			file, err := os.Open(p)
			if err != nil {
				t.Errorf("Unable to open mock file: %v", err)
			}

			defer file.Close()

			reader := &sectionReader{tt.offset, 0, tt.maxRead, file}
			_, err = reader.Read(buffer)
			if err != nil && !tt.wantError {
				t.Errorf("Unable to read from sectionReader: %v", err)
			}

			if !reflect.DeepEqual(buffer, tt.expected) {
				t.Errorf("Got: %v - Want: %v", buffer, tt.expected)
			}
		})
	}
}

func TestInFile(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		outputBytes int64
		seek        int64
		count       int64
		wantErr     bool
	}{
		{
			name:        "Seek to first byte",
			filename:    "datafile",
			outputBytes: 1,
			seek:        0,
			count:       1,
			wantErr:     false,
		},
		{
			name:        "Seek to second byte",
			filename:    "datafile",
			outputBytes: 1,
			seek:        1,
			count:       1,
			wantErr:     false,
		},
		{
			name:        "no filename",
			filename:    "",
			outputBytes: 1,
			seek:        0,
			count:       1,
			wantErr:     false,
		},
		{
			name:        "unknown file",
			filename:    "/something/something",
			outputBytes: 1,
			seek:        0,
			count:       1,
			wantErr:     true,
		},
		{
			name:        "no filename and seek to nowhere",
			filename:    "",
			outputBytes: 8,
			seek:        8,
			count:       1,
			wantErr:     true,
		},
	}

	p, cleanup := setupDatafile(t, "datafile")
	defer cleanup()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := inFile(&bytes.Buffer{}, p, tt.outputBytes, tt.seek, tt.count)
			if err != nil && !tt.wantErr {
				t.Errorf("outFile failed with %v", err)
			}
		})
	}
}

func setupDatafile(t *testing.T, name string) (string, func()) {
	t.Helper()

	testDir := t.TempDir()
	dataFilePath := filepath.Join(testDir, name)

	if err := os.WriteFile(dataFilePath, []byte("ABCDEFG"), 0o644); err != nil {
		t.Errorf("unable to mockup file: %v", err)
	}

	return dataFilePath, func() { os.Remove(dataFilePath) }
}

func TestOutFile(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		outputBytes int64
		seek        int64
		flags       int
		wantErr     bool
	}{
		{
			name:        "Seek to first byte",
			filename:    "datafile",
			outputBytes: 1,
			seek:        0,
			flags:       0,
			wantErr:     false,
		},
		{
			name:        "Seek to second byte",
			filename:    "datafile",
			outputBytes: 1,
			seek:        1,
			flags:       0,
			wantErr:     false,
		},
		{
			name:        "no filename",
			filename:    "",
			outputBytes: 1,
			seek:        0,
			flags:       0,
			wantErr:     false,
		},
		{
			name:        "unknown file",
			filename:    "/something/something",
			outputBytes: 1,
			seek:        0,
			flags:       0,
			wantErr:     true,
		},
		{
			name:        "no filename and seek to nowhere",
			filename:    "",
			outputBytes: 8,
			seek:        8,
			flags:       0,
			wantErr:     true,
		},
	}

	p, cleanup := setupDatafile(t, "datafile")
	defer cleanup()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := outFile(&ws{Writer: &bytes.Buffer{}}, p, tt.outputBytes, tt.seek, tt.flags)
			if err != nil && !tt.wantErr {
				t.Errorf("outFile failed with %v", err)
			}
		})
	}
}

func TestConvertArgs(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		expectedArgs []string
	}{
		{
			name:         "Empty Args",
			args:         []string{""},
			expectedArgs: []string{""},
		},
		{
			name:         "One Arg",
			args:         []string{"if=somefile"},
			expectedArgs: []string{"-if", "somefile"},
		},
		{
			name:         "Two Args",
			args:         []string{"if=somefile", "conv=none"},
			expectedArgs: []string{"-if", "somefile", "-conv", "none"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotArgs := convertArgs(tt.args)

			if !reflect.DeepEqual(gotArgs, tt.expectedArgs) {
				t.Errorf("Args not equal. Got %v, want %v", gotArgs, tt.expectedArgs)
			}
		})
	}
}

// TestDd implements a table-driven test.
func TestDd(t *testing.T) {
	tests := []struct {
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
			name:    "256 MiB zeroed file in 1024 1KiB blocks",
			flags:   []string{"bs=524288", "count=256", "if=/dev/zero"},
			stdin:   "",
			stdout:  []byte("\x00"),
			count:   256 * 1024 * 512,
			compare: byteCount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdin := strings.NewReader(tt.stdin)
			stdout := &bytes.Buffer{}
			stderr := &ws{Writer: io.Discard}
			if err := run(stdin, &ws{Writer: stdout}, stderr, tt.name, tt.flags); err != nil {
				t.Errorf("run: got %v, want nil", err)
			}
			if err := tt.compare(stdout, tt.stdout, tt.count); err != nil {
				t.Errorf("Test compare function returned: %v", err)
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
		for z := range read {
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
	tests := []struct {
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
		// This test only works on Linux.
		{
			// Fully testing the file is synchronous would require something more.
			name:     "dsync",
			flags:    []string{"oflag=dsync"},
			inFile:   []byte("y: defaults"),
			expected: []byte("y: defaults"),
		},
	}

	// This is yucky. But it's simple.
	lastTest := len(tests)
	if runtime.GOOS != "linux" {
		lastTest--
	}

	for _, tt := range tests[:lastTest] {
		t.Run(tt.name, func(t *testing.T) {
			// Write in and out file to temporary dir.
			tmpDir := t.TempDir()
			inFile := filepath.Join(tmpDir, "inFile")
			outFile := filepath.Join(tmpDir, "outFile")
			if err := os.WriteFile(inFile, tt.inFile, 0o666); err != nil {
				t.Error(err)
			}
			if err := os.WriteFile(outFile, tt.outFile, 0o666); err != nil {
				t.Error(err)
			}

			args := append(tt.flags, "if="+inFile, "of="+outFile)
			if err := run(&bytes.Buffer{}, &ws{Writer: io.Discard}, &ws{Writer: io.Discard}, tt.name, args); err != nil {
				t.Error(err)
			}
			got, err := os.ReadFile(filepath.Join(tmpDir, "outFile"))
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(tt.expected, got) {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

type testBS struct {
	name string
	dat  *bytes.Buffer
	sz   int64
	ibs  int64
	obs  int64
	err  error
}

// Write implements write.
func (t *testBS) Write(b []byte) (int, error) {
	if t.sz < t.obs && int64(len(b)) == t.sz {
		return len(b), nil
	}
	if int64(len(b)) != t.obs {
		return -1, fmt.Errorf("%s: write len is %d want %d:%w", t.name, len(b), t.obs, os.ErrInvalid)
	}
	return len(b), nil
}

func (t *testBS) Read(b []byte) (int, error) {
	if int64(len(b)) != t.ibs {
		return -1, fmt.Errorf("%s: read len is %d want %d:%w", t.name, len(b), t.obs, os.ErrInvalid)
	}
	return t.dat.Read(b)
}

func TestBS(t *testing.T) {
	for _, tt := range []testBS{
		{
			name: "dat 128 ibs 128 obs 128",
			dat:  bytes.NewBuffer(make([]byte, 128)),
			sz:   128,
			ibs:  128,
			obs:  128,
			err:  nil,
		},
		{
			name: "dat 26 ibs 26 obs 13",
			dat:  bytes.NewBuffer(make([]byte, 26)),
			sz:   26,
			ibs:  26,
			obs:  13,
			err:  nil,
		},
		{
			name: "dat 26 ibs 13 obs 26",
			dat:  bytes.NewBuffer(make([]byte, 26)),
			sz:   26,
			ibs:  13,
			obs:  26,
			err:  nil,
		},
		{
			name: "dat 14 ibs 14 obs 26",
			dat:  bytes.NewBuffer(make([]byte, 14)),
			sz:   14,
			ibs:  14,
			obs:  26,
			err:  nil,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var bw int64
			if err := dd(&tt, &tt, tt.ibs, tt.obs, &bw); !errors.Is(err, tt.err) {
				t.Fatalf("got %v, want %v", err, tt.err)
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
	if err := run(&bytes.Buffer{}, &ws{Writer: io.Discard}, &ws{Writer: io.Discard}, "dd", args); err != nil {
		b.Fatal(err)
	}
}
