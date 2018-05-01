// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

func TestReadlink(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "mkdir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create files.
	f1 := filepath.Join(tmpDir, "f1")
	if _, err := os.Create(f1); err != nil {
		t.Fatal(err)
	}
	f2 := filepath.Join(tmpDir, "f2")
	if _, err := os.Create(f2); err != nil {
		t.Fatal(err)
	}

	// Create symlinks
	f1symlink := filepath.Join(tmpDir, "f1symlink")
	if err := os.Symlink("f1", f1symlink); err != nil {
		t.Fatal(err)
	}

	// Multiple links
	multilinks := filepath.Join(tmpDir, "multilinks")
	if err := os.Symlink(f1symlink, multilinks); err != nil {
		t.Fatal(err)
	}

	for i, tt := range []struct {
		flags      []string
		out        string
		stderr     string
		exitStatus int
	}{
		{
			flags:      []string{},
			out:        "",
			stderr:     "",
			exitStatus: 1,
		}, {
			flags:      []string{"-v", f1},
			out:        "",
			stderr:     fmt.Sprintf("readlink %s: invalid argument\n", f1),
			exitStatus: 1,
		}, {
			flags:      []string{"-f", f2},
			out:        "",
			stderr:     "",
			exitStatus: 1,
		}, {
			flags:      []string{f1symlink},
			out:        "f1\n",
			stderr:     "",
			exitStatus: 0,
		}, {
			flags:      []string{multilinks},
			out:        fmt.Sprintf("%s\n", f1symlink),
			stderr:     "",
			exitStatus: 0,
		}, {
			flags:      []string{"-v", multilinks, f1symlink, f2},
			out:        fmt.Sprintf("%s\nf1\n", f1symlink),
			stderr:     fmt.Sprintf("readlink %s: invalid argument\n", f2),
			exitStatus: 1,
		}, {
			flags:      []string{"-v", tmpDir},
			out:        "",
			stderr:     fmt.Sprintf("readlink %s: invalid argument\n", tmpDir),
			exitStatus: 1,
		}, {
			flags:      []string{"-v", "foo.bar"},
			out:        "",
			stderr:     fmt.Sprintf("readlink foo.bar: no such file or directory\n"),
			exitStatus: 1,
		},
	} {
		t.Run(fmt.Sprintf("test%d", i), func(t *testing.T) {
			var out, stderr bytes.Buffer
			cmd := testutil.Command(t, tt.flags...)
			cmd.Stdout = &out
			cmd.Stderr = &stderr
			err := cmd.Run()

			if out.String() != tt.out {
				t.Errorf("stdout got: %q, want: %q", out.String(), tt.out)
			}

			if stderr.String() != tt.stderr {
				t.Errorf("stderr got: %q, want: %q", stderr.String(), tt.stderr)
			}

			if tt.exitStatus == 0 && err != nil {
				t.Errorf("expected to exit with %d, but exited with err %v", tt.exitStatus, err)
			}
		})
	}
}

func TestMain(m *testing.M) {
	if testutil.CallMain() {
		main()
		os.Exit(0)
	}

	os.Exit(m.Run())
}
