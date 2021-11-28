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

	"github.com/u-root/u-root/pkg/testutil"
)

type test struct {
	flags      []string
	out        string
	stdErr     string
	exitStatus int
}

func TestReadlink(t *testing.T) {
	tmpDir := t.TempDir()

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
			flags:      []string{},
			out:        "",
			stdErr:     "",
			exitStatus: 1,
		},
		{
			flags:      []string{"-v", "f1"},
			out:        "",
			stdErr:     "readlink f1: invalid argument\n",
			exitStatus: 1,
		},
		{
			flags:      []string{"-f", "f2"},
			out:        "",
			stdErr:     "",
			exitStatus: 1,
		},
		{
			flags:      []string{"f1symlink"},
			out:        "f1\n",
			stdErr:     "",
			exitStatus: 0,
		},
		{
			flags:      []string{"multilinks"},
			out:        fmt.Sprintf("%s/%s", testDir, "f1symlink\n"),
			stdErr:     "",
			exitStatus: 0,
		},
		{
			flags:      []string{"-v", "multilinks", "f1symlink", "f2"},
			out:        fmt.Sprintf("%s/%sf1\n", testDir, "f1symlink\n"),
			stdErr:     "readlink f2: invalid argument\n",
			exitStatus: 1,
		},
		{
			flags:      []string{"-v", testDir},
			out:        "",
			stdErr:     fmt.Sprintf("readlink %s: invalid argument\n", testDir),
			exitStatus: 1,
		},
		{
			flags:      []string{"-v", "foo.bar"},
			out:        "",
			stdErr:     "readlink foo.bar: no such file or directory\n",
			exitStatus: 1,
		},
	}
	// Createfiles.
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
		var out, stdErr bytes.Buffer
		cmd := testutil.Command(t, tt.flags...)
		cmd.Stdout = &out
		cmd.Stderr = &stdErr
		err := cmd.Run()

		if out.String() != tt.out {
			t.Errorf("stdout got:\n%s\nwant:\n%s", out.String(), tt.out)
		}

		if stdErr.String() != tt.stdErr {
			t.Errorf("stderr got:\n%s\nwant:\n%s", stdErr.String(), tt.stdErr)
		}

		if tt.exitStatus == 0 && err != nil {
			t.Errorf("expected to exit with %d, but exited with err %s", tt.exitStatus, err)
		}
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
