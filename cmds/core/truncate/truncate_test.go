// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rck/unit"
)

func TestTruncate(t *testing.T) {
	tmpdir := t.TempDir()
	for _, tt := range []struct {
		name     string
		args     []string
		create   bool
		size     unit.Value
		sizeWant int64
		rfile    string
		want     string
	}{
		{
			name: "!size.IsSet && *rfile == \"\"",
			want: "you need to specify size via -s <number> or -r <rfile>",
		},
		{
			name: "size.IsSet && *rfile == \"\"",
			size: unit.Value{
				IsSet: true,
			},
			rfile: "testfile",
			want:  "you need to specify size via -s <number> or -r <rfile>",
		},
		{
			name:  "len(args) == 0",
			rfile: "testfile",
			want:  "you need to specify one or more files as argument",
		},
		{
			name:   "create non existing file with err in rfile path",
			create: false,
			rfile:  "testfile",
			args:   []string{filepath.Join(tmpdir, "file1")},
			want:   "could not stat reference file: stat testfile: no such file or directory",
		},
		{
			name:   "create non existing file without err in rfile path",
			create: false,
			rfile:  filepath.Join(tmpdir, "file1"),
			args:   []string{filepath.Join(tmpdir, "file1")},
		},
		{
			name:   "truncate existing file without err and positive size",
			create: false,
			size: unit.Value{
				Value:        -50,
				IsSet:        true,
				ExplicitSign: unit.Sign(2),
			},
			sizeWant: 0,
			args:     []string{filepath.Join(tmpdir, "file1")},
		},
		{
			name:   "truncate existing file without err and negative size",
			create: false,
			size: unit.Value{
				Value:        50,
				IsSet:        true,
				ExplicitSign: unit.Sign(1),
			},
			sizeWant: 50,
			args:     []string{filepath.Join(tmpdir, "file1")},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			*create = tt.create
			*size = tt.size
			*rfile = tt.rfile
			if got := truncate(tt.args...); got != nil {
				if got.Error() != tt.want {
					t.Errorf("truncate() = %q, want: %q", got.Error(), tt.want)
				}
			} else {
				st, err := os.Stat(filepath.Join(tmpdir, "file1"))
				if err != nil {
					t.Errorf("failed to get file stats")
				}
				if st.Size() != tt.sizeWant {
					t.Errorf("the file was not truncated right, got file size %v, want: %v", st.Size(), tt.sizeWant)
				}
			}
		})
	}
}
