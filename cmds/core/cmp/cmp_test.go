// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
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
		want   string
	}{
		{
			name: "empty args",
			args: []string{},
			want: fmt.Sprintf("expected two filenames (and one to two optional offsets), got %d", 0),
		},
		{
			name: "err in open file",
			args: []string{"filedoesnotexist", "filealsodoesnotexist"},
			want: fmt.Sprintf("failed to open %s: %s", "filedoesnotexist", "open filedoesnotexist: no such file or directory"),
		},
		{
			name:  "cmp two files without flags, without offsets",
			args:  []string{filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2")},
			file1: "hello\nthis is a test\n",
			file2: "hello\nthiz is a text",
			want:  fmt.Sprintf("%s %s differ: char %d", filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2"), 10),
		},
		{
			name: "cmp two files without flags, with wrong first offset",
			args: []string{filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2"), ""},
			want: fmt.Sprintf("bad offset1: %s: %v", "", "invalid size \"\""),
		},
		{
			name:  "cmp two files without flags, with correct first offset",
			args:  []string{filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2"), "4"},
			file1: "hello\nthis is a test\n",
			file2: "hello\nthiz is a text",
			want:  fmt.Sprintf("%s %s differ: char %d", filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2"), 1),
		},
		{
			name:  "cmp two files without flags, with both offsets correct",
			args:  []string{filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2"), "4", "4"},
			file1: "hello\nthis is a test\n",
			file2: "hello\nthiz is a text",
			want:  fmt.Sprintf("%s %s differ: char %d", filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2"), 6),
		},
		{
			name:  "cmp two files without flags, with both offset set but first not valid",
			args:  []string{filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2"), "", "4"},
			file1: "hello\nthis is a test\n",
			file2: "hello\nthiz is a text",
			want:  fmt.Sprintf("bad offset1: %s: %v", "", "invalid size \"\""),
		},
		{
			name:  "cmp two files without flags, with both offset set but second not valid",
			args:  []string{filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2"), "4", ""},
			file1: "hello\nthis is a test\n",
			file2: "hello\nthiz is a text",
			want:  fmt.Sprintf("bad offset2: %s: %v", "", "invalid size \"\""),
		},
		{
			name:  "cmp two files, flag line = true",
			args:  []string{filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2")},
			file1: "hello\nthis is a test\n",
			file2: "hello\nthiz is a text",
			line:  true,
			want:  fmt.Sprintf("%s %s differ: char %d line %d", filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2"), 10, 2),
		},
		{
			name:   "cmp two files, flag silent = true",
			args:   []string{filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2")},
			file1:  "hello\nthis is a test\n",
			file2:  "hello\nthiz is a text",
			silent: true,
			want:   "",
		},
		{
			name:  "cmp two files, flag long = true",
			args:  []string{filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2")},
			file1: "hello\nthis is a test",
			file2: "hello\nthiz is a text",
			long:  true,
			want:  fmt.Sprintf("%8d %#.2o %#.2o\n%8d %#.2o %#.2o\n", 10, 0163, 0172, 19, 0163, 0170),
		},
		{
			name:  "cmp two files, flag long = true, first file ends first",
			args:  []string{filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2")},
			file1: "hello\nthis is a tes",
			file2: "hello\nthiz is a text",
			long:  true,
			want:  fmt.Sprintf("EOF on %s", filepath.Join(tmpdir, "file1")),
		},
		{
			name:  "cmp two files, flag long = true, second file ends first",
			args:  []string{filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2")},
			file1: "hello\nthis is a test",
			file2: "hello\nthiz is a tex",
			long:  true,
			want:  fmt.Sprintf("EOF on %s", filepath.Join(tmpdir, "file2")),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// Write data in file that should be compared
			file1, err := os.Create(filepath.Join(tmpdir, "file1"))
			if err != nil {
				t.Errorf("failed to create file1: %v", err)
			}
			file2, err := os.Create(filepath.Join(tmpdir, "file2"))
			if err != nil {
				t.Errorf("failed to create file2: %v", err)
			}
			_, err = file1.WriteString(tt.file1)
			if err != nil {
				t.Errorf("failed to write to file1: %v", err)
			}
			_, err = file2.WriteString(tt.file2)
			if err != nil {
				t.Errorf("failed to write to file2: %v", err)
			}

			// Set flags
			*long = tt.long
			*line = tt.line
			*silent = tt.silent

			// Start tests
			buf := &bytes.Buffer{}
			if got := cmp(buf, tt.args...); got != nil {
				if got.Error() != tt.want {
					t.Errorf("cmp() = %q, want: %q", got.Error(), tt.want)
				}
			} else {
				if buf.String() != tt.want {
					t.Errorf("cmp() = %q, want: %q", buf.String(), tt.want)
				}
			}

		})
	}
}
