// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

func prepareDirLayout(t *testing.T) {
	t.Helper()
	tmpDir := t.TempDir()
	err := os.Chdir(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Create("file1")
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Create("file2")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Mkdir("dir1", os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Create("dir1/file1")
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Create("dir1/file2")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Mkdir("dir2", os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Create("dir2/file1")
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Create("dir2/file3")
	if err != nil {
		t.Fatal(err)
	}
}

func TestFind(t *testing.T) {
	prepareDirLayout(t)

	var tests = []struct {
		wantStdout string
		commandErr error
		runErr     error
		args       []string
	}{
		{
			wantStdout: "file1\n",
			args:       []string{"file1"},
		},
		{
			wantStdout: "dir1\ndir1/file1\ndir1/file2\n",
			args:       []string{"dir1"},
		},
		{
			wantStdout: "dir1/file1\ndir2/file1\nfile1\n",
			args:       []string{"-name=file1", "."},
		},
		{
			wantStdout: ".\ndir1\ndir2\n",
			args:       []string{"-type=d", "."},
		},
		{
			wantStdout: ".\ndir1\ndir2\n",
			args:       []string{"-type=directory", "."},
		},
		{
			wantStdout: "dir1/file1\ndir1/file2\ndir2/file1\ndir2/file3\nfile1\nfile2\n",
			args:       []string{"-type=f", "."},
		},
		{
			wantStdout: "dir1/file1\ndir1/file2\ndir2/file1\ndir2/file3\nfile1\nfile2\n",
			args:       []string{"-type=file", "."},
		},
		{
			args:   []string{"-type=notvalid", "."},
			runErr: errNotValidType,
		},
		{
			wantStdout: "file1\n",
			args:       []string{"-mode=0644", "file1"},
		},
		{
			args:       []string{"-mode=0644"},
			commandErr: errUsage,
		},
	}

	for _, tt := range tests {
		var stdout bytes.Buffer
		c, err := command(&stdout, io.Discard, tt.args)
		if !errors.Is(err, tt.commandErr) {
			t.Fatalf("expected %v, got %v", tt.commandErr, err)
		}
		if err != nil {
			continue
		}

		err = c.run()
		if !errors.Is(err, tt.runErr) {
			t.Fatalf("expected %v, got %v", tt.runErr, err)
		}
		if err != nil {
			continue
		}

		resStdout := stdout.String()
		if resStdout != tt.wantStdout {
			t.Errorf("want\n %s, got\n %s", tt.wantStdout, resStdout)
		}
	}
}

func TestFindLong(t *testing.T) {
	prepareDirLayout(t)

	var stdout bytes.Buffer
	c, err := command(&stdout, nil, []string{"-l", "file1"})
	if err != nil {
		t.Fatal(err)
	}

	err = c.run()
	if err != nil {
		t.Fatal(err)
	}

	res := strings.TrimSpace(stdout.String())

	if !strings.HasPrefix(res, "-rw-r--r--") {
		t.Errorf("want prefix: -rw-r--r--, got prefix: %s", res[:10])
	}
	if !strings.HasSuffix(res, "file1") {
		t.Errorf("want suffix: file1, got suffix: %s", res[len(res)-5:])
	}
}
