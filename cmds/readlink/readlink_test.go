// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

type Test struct {
	flags      []string
	out        string
	stdErr     string
	exitStatus int
}

func TestReadlink(t *testing.T) {

	// Create an empty directory
	tmpDir, execPath := testutil.CompileInTempDir(t)
	defer os.RemoveAll(tmpDir)

	// Creating here to utilize path in tests
	testDir := filepath.Join(tmpDir, "readLinkDir")
	os.Mkdir(testDir, 0700)
	os.Chdir(testDir)

	var tests = []Test{
		{
			flags:      []string{},
			out:        "",
			stdErr:     "",
			exitStatus: 0,
		}, {
			flags:      []string{"-v", "f1"},
			out:        "",
			stdErr:     "readlink: f1 Invalid argument\n",
			exitStatus: 1,
		}, {
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
			flags:      []string{"multiLinks"},
			out:        fmt.Sprintf("%s/%s", testDir, "f1Symlink\n"),
			stdErr:     "",
			exitStatus: 0,
		},
		{
			flags:      []string{"-v", "multiLinks", "f1symlink", "f2"},
			out:        fmt.Sprintf("%s/%sf1\n", testDir, "f1Symlink\n"),
			stdErr:     "readlink: f2 Invalid argument\n",
			exitStatus: 1,
		},
	}
	// Createfiles.
	os.Create("f1")
	os.Create("f2")

	// Create symlinks
	f1Symlink := filepath.Join(testDir, "f1Symlink")
	os.Symlink("f1", f1Symlink)

	// Multiple links
	multiLinks := filepath.Join(testDir, "multiLinks")
	os.Symlink(f1Symlink, multiLinks)

	// Table-driven testing
	for _, tt := range tests {

		var out, stdErr bytes.Buffer
		cmd := exec.Command(execPath, tt.flags...)
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
