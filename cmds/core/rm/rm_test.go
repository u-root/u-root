// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

func setup(t *testing.T) string {
	d := t.TempDir()
	fbody := []byte("Go is cool!")
	for _, f := range []struct {
		name  string
		mode  os.FileMode
		isdir bool
	}{
		{
			name:  "hi",
			mode:  0o755,
			isdir: true,
		},
		{
			name: "hi/one.txt",
			mode: 0o666,
		},
		{
			name: "hi/two.txt",
			mode: 0o777,
		},
		{
			name: "go.txt",
			mode: 0o555,
		},
	} {
		var (
			err      error
			filepath = filepath.Join(d, f.name)
		)
		if f.isdir {
			err = os.Mkdir(filepath, f.mode)
		} else {
			err = os.WriteFile(filepath, fbody, f.mode)
		}
		if err != nil {
			t.Fatal(err)
		}
	}
	return d
}

func TestRm(t *testing.T) {
	for _, tt := range []struct {
		name    string
		flags   []string
		files   []string
		iString string
		err     error
	}{
		{
			name: "no args",
			err:  errUsage,
		},
		{
			name:  "rm one file",
			files: []string{"go.txt"},
		},
		{
			name:  "rm one file verbose",
			files: []string{"go.txt"},
			flags: []string{"-v"},
		},
		{
			name:  "fail to rm one file",
			files: []string{"go"},
			err:   os.ErrNotExist,
		},
		{
			name:  "fail to rm one file forced to trigger continue",
			files: []string{"go"},
			flags: []string{"-f"},
		},
		{
			name:    "rm one file interactive (y)",
			files:   []string{"go.txt"},
			flags:   []string{"-i"},
			iString: "y\n",
		},
		{
			name:    "rm one file interactive (Y)",
			files:   []string{"go.txt"},
			flags:   []string{"-i"},
			iString: "Y\n",
		},
		{
			name:    "rm one file interactive continue triggered",
			files:   []string{"go.txt"},
			flags:   []string{"-i"},
			iString: "\n",
		},
		{
			name:  "rm two files with verbose and force (unixflags)",
			files: []string{"hi/one.txt", "hi/two.txt"},
			flags: []string{"-vf"},
		},
		{
			name:  "rm two files with verbose and force (goflags)",
			files: []string{"hi/one.txt", "hi/two.txt"},
			flags: []string{"-v", "-f"},
		},
		{
			name:  "rm dir recursively",
			files: []string{"hi"},
			flags: []string{"-r"},
		},
		{
			name:  "rm dir recursively second flag",
			files: []string{"hi"},
			flags: []string{"-R"},
		},
		{
			name:  "rm dir not recursively",
			files: []string{"hi"},
			err:   syscall.ENOTEMPTY,
		},
	} {
		d := setup(t)

		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			buf.WriteString(tt.iString)

			var args []string
			args = append(args, tt.flags...)
			for _, file := range tt.files {
				args = append(args, filepath.Join(d, file))
			}

			err := rm(buf, io.Discard, io.Discard, args...)
			if !errors.Is(err, tt.err) {
				t.Errorf("expected %v, got %v", tt.err, err)
			}
		})
	}
}
