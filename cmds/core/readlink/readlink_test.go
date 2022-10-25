// Copyright 2017 the u-root Authors. All rights reserved
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

type test struct {
	args      []string
	out       string
	stdErr    string
	hasError  bool
	follow    bool
	noNewLine bool
	verbose   bool
}

func TestReadlink(t *testing.T) {
	tmpDir := t.TempDir()
	defer os.Remove(tmpDir)

	// Creating here to utilize path in tests
	testDir := filepath.Join(tmpDir, "readLinkDir")
	if err := os.Mkdir(testDir, 0o700); err != nil {
		t.Error(err)
	}

	if err := os.Chdir(testDir); err != nil {
		t.Error(err)
	}

	tests := []test{
		{
			args:     []string{},
			out:      "",
			stdErr:   "",
			hasError: true,
		},
		{
			args:     []string{},
			out:      "",
			stdErr:   "missing operand",
			hasError: true,
			verbose:  true,
		},
		{
			args:     []string{"f1"},
			out:      "",
			stdErr:   "readlink f1: invalid argument\n",
			hasError: true,
			verbose:  true,
		},
		{
			args:     []string{"f2"},
			out:      "",
			stdErr:   "",
			hasError: true,
			follow:   true,
		},
		{
			args:     []string{"f1symlink"},
			out:      "f1\n",
			stdErr:   "",
			hasError: false,
		},
		{
			args:      []string{"f1symlink"},
			out:       "f1",
			stdErr:    "",
			hasError:  false,
			noNewLine: true,
		},
		{
			args:     []string{"multilinks"},
			out:      fmt.Sprintf("%s/%s", testDir, "f1symlink\n"),
			stdErr:   "",
			hasError: false,
		},
		{
			args:     []string{"multilinks", "f1symlink", "f2"},
			out:      fmt.Sprintf("%s/%sf1\n", testDir, "f1symlink\n"),
			stdErr:   "readlink f2: invalid argument\n",
			hasError: true,
			verbose:  true,
		},
		{
			args:     []string{testDir},
			out:      "",
			stdErr:   fmt.Sprintf("readlink %s: invalid argument\n", testDir),
			hasError: true,
			verbose:  true,
		},
		{
			args:     []string{"foo.bar"},
			out:      "",
			stdErr:   "readlink foo.bar: no such file or directory\n",
			hasError: true,
			verbose:  true,
		},
	}
	// Create files.
	if _, err := os.Create("f1"); err != nil {
		t.Error(err)
	}

	if _, err := os.Create("f2"); err != nil {
		t.Error(err)
	}

	// Create symlinks
	f1Symlink := filepath.Join(testDir, "f1symlink")
	if err := os.Symlink("f1", f1Symlink); err != nil {
		t.Error(err)
	}

	// Multiple links
	multiLinks := filepath.Join(testDir, "multilinks")
	if err := os.Symlink(f1Symlink, multiLinks); err != nil {
		t.Error(err)
	}

	// Table-driven testing
	for _, tt := range tests {
		*verbose = tt.verbose
		*noNewLine = tt.noNewLine
		*follow = tt.follow

		var out, stdErr bytes.Buffer
		err := run(&out, &stdErr, tt.args)

		if out.String() != tt.out {
			t.Errorf("stdout got:\n%s\nwant:\n%s", out.String(), tt.out)
		}

		if stdErr.String() != tt.stdErr {
			t.Errorf("stderr got:\n%s\nwant:\n%s", stdErr.String(), tt.stdErr)
		}

		if tt.hasError && err == nil {
			t.Error("expected to exit with error")
		}

		if !tt.hasError && err != nil {
			t.Error("expected to exit without error")
		}
	}
}
