// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
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
		err   error
		want  string
	}{
		{
			name: "only one arguement",
			args: []string{"onearg"},
			err:  ErrUsage,
		},
		{
			name: "help flag",
			args: []string{"firstarg", "secondarg"},
			help: true,
			err:  ErrUsage,
		},
		{
			name: "first file failed to open",
			args: []string{"firstarg", "secondarg"},
			err:  os.ErrNotExist,
		},
		{
			name: "second file failed to open",
			args: []string{filepath.Join(tmpdir, "file1"), "secondarg"},
			err:  os.ErrNotExist,
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
			if err := os.WriteFile(filepath.Join(tmpdir, "file1"), []byte(tt.file1), 0o644); err != nil {
				t.Errorf("failed to create file1: %v", err)
			}
			if err := os.WriteFile(filepath.Join(tmpdir, "file2"), []byte(tt.file2), 0o644); err != nil {
				t.Errorf("failed to create file1: %v", err)
			}

			buf := &bytes.Buffer{}
			if err := comm(buf, tt.s1, tt.s2, tt.s3, tt.help, tt.args...); !errors.Is(err, tt.err) {
				t.Errorf("comm() = %q, want: %q", err, tt.err)
			}
			if buf.String() != tt.want {
				t.Errorf("comm() = %q, want: %q", buf.String(), tt.want)
			}
		})
	}
}
