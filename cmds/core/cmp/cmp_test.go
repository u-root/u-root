// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestCmp(t *testing.T) {
	tmpdir := t.TempDir()
	for _, tt := range []struct {
		name   string
		args   []string
		file1  string
		file2  string
		long   bool
		line   bool
		silent bool
		err    error
		stderr string
		stdout string
	}{
		{
			name:   "empty args",
			args:   []string{},
			err:    ErrArgCount,
			stderr: usage,
		},
		{
			name: "err in open file",
			args: []string{"filedoesnotexist", "filealsodoesnotexist"},
			err:  os.ErrNotExist,
		},
		{
			name:   "cmp two files without flags, without offsets",
			args:   []string{filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2")},
			file1:  "hello\nthis is a test\n",
			file2:  "hello\nthiz is a text",
			err:    ErrDiffer,
			stdout: fmt.Sprintf("%s %s: char %d", filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2"), 10),
		},
		{
			name:   "cmp two files without flags, with wrong first offset",
			args:   []string{filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2"), ""},
			err:    ErrBadOffset,
			stderr: fmt.Sprintf("bad offset1: %s: %v", "", "invalid size \"\""),
		},
		{
			name:   "cmp two files without flags, with correct first offset",
			args:   []string{filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2"), "4"},
			file1:  "hello\nthis is a test\n",
			file2:  "hello\nthiz is a text",
			err:    ErrDiffer,
			stdout: fmt.Sprintf("%s %s: char %d", filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2"), 1),
		},
		{
			name:   "cmp two files without flags, with both offsets correct",
			args:   []string{filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2"), "4", "4"},
			file1:  "hello\nthis is a test\n",
			file2:  "hello\nthiz is a text",
			err:    ErrDiffer,
			stdout: fmt.Sprintf("%s %s: char %d", filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2"), 6),
		},
		{
			name:   "cmp two files without flags, with both offset set but first not valid",
			args:   []string{filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2"), "", "4"},
			file1:  "hello\nthis is a test\n",
			file2:  "hello\nthiz is a text",
			err:    ErrBadOffset,
			stderr: fmt.Sprintf("bad offset1: %s: %v", "", "invalid size \"\""),
		},
		{
			name:   "cmp two files without flags, with both offset set but second not valid",
			args:   []string{filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2"), "4", ""},
			file1:  "hello\nthis is a test\n",
			file2:  "hello\nthiz is a text",
			err:    ErrBadOffset,
			stderr: fmt.Sprintf("bad offset2: %s: %v", "", "invalid size \"\""),
		},
		{
			name:   "cmp two files, flag line = true",
			args:   []string{filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2")},
			file1:  "hello\nthis is a test\n",
			file2:  "hello\nthiz is a text",
			line:   true,
			err:    ErrDiffer,
			stdout: fmt.Sprintf("%s %s: char %d line %d", filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2"), 10, 2),
		},
		{
			name:   "cmp two files, flag silent = true",
			args:   []string{filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2")},
			file1:  "hello\nthis is a test\n",
			file2:  "hello\nthiz is a text",
			silent: true,
		},
		{
			name:   "cmp two files, flag long = true",
			args:   []string{filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2")},
			file1:  "hello\nthis is a test",
			file2:  "hello\nthiz is a text",
			long:   true,
			stdout: fmt.Sprintf("%8d %#.2o %#.2o\n%8d %#.2o %#.2o\n", 10, 0o163, 0o172, 19, 0o163, 0o170),
		},
		{
			name:   "cmp two files, flag long = true, first file ends first",
			args:   []string{filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2")},
			file1:  "hello\nthis is a tes",
			file2:  "hello\nthiz is a text",
			long:   true,
			err:    io.EOF,
			stdout: fmt.Sprintf("%8d %#.2o %#.2o\n%8d %#.2o %#.2o\n", 10, 0o163, 0o172, 19, 0o163, 0o170),
			stderr: fmt.Sprintf("%s:%v", filepath.Join(tmpdir, "file1"), io.EOF),
		},
		{
			name:   "cmp two files, flag long = true, second file ends first",
			args:   []string{filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2")},
			file1:  "hello\nthis is a test",
			file2:  "hello\nthiz is a tex",
			long:   true,
			err:    io.EOF,
			stdout: fmt.Sprintf("%8d %#.2o %#.2o\n%8d %#.2o %#.2o\n", 10, 0o163, 0o172, 19, 0o163, 0o170),
			stderr: fmt.Sprintf("%s:%v", filepath.Join(tmpdir, "file2"), io.EOF),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f1 := filepath.Join(tmpdir, "file1")
			if err := os.WriteFile(f1, []byte(tt.file1), 0o666); err != nil {
				t.Fatal(err)
			}
			f2 := filepath.Join(tmpdir, "file2")
			if err := os.WriteFile(f2, []byte(tt.file2), 0o666); err != nil {
				t.Fatal(err)
			}

			// Start tests
			var stdout, stderr bytes.Buffer
			if err := cmp(&stdout, &stderr, tt.long, tt.line, tt.silent, tt.args...); !errors.Is(err, tt.err) {
				t.Errorf("cmp(): got %v, want %v", err, tt.err)
			}
			if stdout.String() != tt.stdout {
				t.Errorf("cmp():stdout:got %s, want %s", stdout.String(), tt.stdout)
			}
			if stderr.String() != tt.stderr {
				t.Errorf("cmp():stderr:got %s, want %s", stderr.String(), tt.stderr)
			}
		})
	}
}
