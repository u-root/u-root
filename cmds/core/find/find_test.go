// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/find"
)

// create creates a file with a standard mode.
// Do not use os.Create alone; it is sensitive to umask and
// can fail at times.
func create(t *testing.T, name string) {
	t.Helper()
	f, err := os.Create(name)
	if err != nil {
		t.Fatalf("creating test file: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("closing test file: %v", err)
	}
	if err := os.Chmod(name, 0644); err != nil {
		t.Fatalf("setting test file mode: %v", err)
	}
}

func prepareDirLayout(t *testing.T) {
	t.Helper()
	d := t.TempDir()
	// this Chdir applies to the entire test, including
	// the caller. If you remove it, you will need to rewrite
	// the mkdir/create calls below, as well as all args in the
	// test.
	if err := os.Chdir(d); err != nil {
		t.Fatal(err)
	}
	for _, n := range []string{"1", "2"} {
		if err := os.Mkdir(n, os.ModePerm); err != nil {
			t.Fatal(err)
		}
	}

	for _, n := range []string{"file1", "file2", filepath.Join("1", "file1"), filepath.Join("1", "file2"), filepath.Join("2", "file1"), filepath.Join("2", "file3")} {
		create(t, n)
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
			wantStdout: "1\n1/file1\n1/file2\n",
			args:       []string{"1"},
		},
		{
			wantStdout: "1/file1\n2/file1\nfile1\n",
			args:       []string{"-name=file1", "."},
		},
		{
			wantStdout: ".\n1\n2\n",
			args:       []string{"-type=d", "."},
		},
		{
			wantStdout: ".\n1\n2\n",
			args:       []string{"-type=directory", "."},
		},
		{
			wantStdout: "1/file1\n1/file2\n2/file1\n2/file3\nfile1\nfile2\n",
			args:       []string{"-type=f", "."},
		},
		{
			wantStdout: "1/file1\n1/file2\n2/file1\n2/file3\nfile1\nfile2\n",
			args:       []string{"-type=file", "."},
		},
		{
			wantStdout: "1/file1\n2/file1\n",
			args:       []string{"-regex=[12]/file1", "."},
		},
		{
			wantStdout: "1/file2\n",
			args:       []string{"-regex=1/file2", "."},
		},
		{
			wantStdout: "1/file1\n2/file1\nfile1\n",
			args:       []string{"-regex=file1", "."},
		},
		{
			wantStdout: "1/file1\n2/file1\n",
			args:       []string{"-name=file1", "-regex=[12]/file1", "."},
		},
		{
			args:       []string{"-regex=[", "."},
			commandErr: find.ErrInvalidRegexp,
		},
		{
			args:       []string{"-type=notvalid", "."},
			commandErr: errNotValidType,
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

	failed := []int{}
	for i, tt := range tests {
		var stdout bytes.Buffer
		t.Logf("%d:%v", i, tt)
		c, err := command(&stdout, io.Discard, tt.args)
		t.Logf("%v, %v", c, err)
		if !errors.Is(err, tt.commandErr) {
			failed = append(failed, i)
			t.Errorf("expected %v, got %v", tt.commandErr, err)
			continue
		}
		if err != nil {
			continue
		}

		t.Logf("Now run test %d", i)
		err = c.run()
		if !errors.Is(err, tt.runErr) {
			failed = append(failed, i)
			t.Errorf("expected %v, got %v", tt.runErr, err)
			continue
		}
		if err != nil {
			continue
		}

		resStdout := stdout.String()
		if resStdout != tt.wantStdout {
			t.Errorf("want\n %s, got\n %s", tt.wantStdout, resStdout)
			failed = append(failed, i)
		}
	}
	if len(failed) > 0 {
		t.Logf("failing tests: %v", failed)
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
