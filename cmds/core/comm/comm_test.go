// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestComm(t *testing.T) {
	tmpdir := t.TempDir()
	for _, tt := range []struct {
		name  string
		args  []string
		file1 string
		file2 string
		s1    bool
		s2    bool
		s3    bool
		help  bool
		want  string
	}{
		{
			name: "only one arguement",
			args: []string{"onearg"},
			want: ErrUsage.Error(),
		},
		{
			name: "help flag",
			args: []string{"firstarg", "secondarg"},
			help: true,
			want: ErrUsage.Error(),
		},
		{
			name: "first file failed to open",
			args: []string{"firstarg", "secondarg"},
			want: "can't open firstarg: open firstarg: no such file or directory",
		},
		{
			name: "second file failed to open",
			args: []string{filepath.Join(tmpdir, "file1"), "secondarg"},
			want: "can't open secondarg: open secondarg: no such file or directory",
		},
		{
			name:  "comm case s1 > s2",
			args:  []string{filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2")},
			file1: "line1\nline2\nline3\n",
			file2: "line1\nline2\nline4\n",
			want:  "\t\tline1\n\t\tline2\nline3\n\tline4\n",
		},
		{
			name:  "comm case s1 < s2",
			args:  []string{filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2")},
			file1: "line1\nline2\nline4\n",
			file2: "line1\nline2\nline3\n",
			want:  "\t\tline1\n\t\tline2\n\tline3\nline4\n",
		},
		{
			name:  "comm flag s1 true",
			args:  []string{filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2")},
			file1: "line1\nline2\nline4\n",
			file2: "line1\nline2\nline3\n",
			s1:    true,
			want:  "\t\tline1\n\t\tline2\n\tline3\n",
		},
		{
			name:  "comm flag s2 true",
			args:  []string{filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2")},
			file1: "line1\nline2\nline4\n",
			file2: "line1\nline2\nline3\n",
			s2:    true,
			want:  "\t\tline1\n\t\tline2\nline4\n",
		},
		{
			name:  "comm flag s3 true",
			args:  []string{filepath.Join(tmpdir, "file1"), filepath.Join(tmpdir, "file2")},
			file1: "line1\nline2\nline4\n",
			file2: "line1\nline2\nline3\n",
			s3:    true,
			want:  "\tline3\nline4\n",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// Create files
			file1, err := os.Create(filepath.Join(tmpdir, "file1"))
			if err != nil {
				t.Errorf("failed to create file1: %v", err)
			}
			file2, err := os.Create(filepath.Join(tmpdir, "file2"))
			if err != nil {
				t.Errorf("failed to create file2: %v", err)
			}
			// Write data in file that should be compared
			_, err = file1.WriteString(tt.file1)
			if err != nil {
				t.Errorf("failed to write to file1: %v", err)
			}
			_, err = file2.WriteString(tt.file2)
			if err != nil {
				t.Errorf("failed to write to file2: %v", err)
			}

			// Setting flags
			*s1 = tt.s1
			*s2 = tt.s2
			*s3 = tt.s3
			*help = tt.help

			buf := &bytes.Buffer{}
			if got := comm(buf, tt.args...); got != nil {
				if got.Error() != tt.want {
					t.Errorf("comm() = %q, want: %q", got.Error(), tt.want)
				}
			} else {
				if buf.String() != tt.want {
					t.Errorf("comm() = %q, want: %q", buf.String(), tt.want)
				}
			}
		})
	}
}
