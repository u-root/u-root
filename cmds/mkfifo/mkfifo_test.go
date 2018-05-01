// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

func TestMkfifo(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "mkfifo")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// used later in testing
	testDir := filepath.Join(tmpDir, "mkfifoDir")
	if err := os.Mkdir(testDir, 0700); err != nil {
		t.Error(err)
	}

	for _, tt := range []struct {
		name   string
		flags  []string
		stderr string
	}{
		{
			name:   "no path or mode, error",
			flags:  []string{},
			stderr: "please provide a path, or multiple, to create a fifo",
		},
		{
			name:   "single path",
			flags:  []string{filepath.Join(testDir, "testfifo")},
			stderr: "",
		},
		{
			name:   "duplicate path",
			flags:  []string{filepath.Join(testDir, "testfifo1"), filepath.Join(testDir, "testfifo1")},
			stderr: "file exists",
		},
	} {
		var stderr bytes.Buffer
		cmd := testutil.Command(t, tt.flags...)
		cmd.Stderr = &stderr
		err := cmd.Run()

		if err != nil && !strings.Contains(stderr.String(), tt.stderr) {
			t.Errorf("expected %v got %v", tt.stderr, stderr.String())
		}

		for _, path := range tt.flags {
			testFile, err := os.Stat(path)
			if err != nil {
				t.Errorf("Unable to stat file %s", path)
			}

			mode := testFile.Mode()
			if typ := mode & os.ModeType; typ != os.ModeNamedPipe {
				t.Errorf("got %v, want %v", typ, os.ModeNamedPipe)
			}
		}
	}
}

func TestMain(m *testing.M) {
	if testutil.CallMain() {
		main()
		os.Exit(0)
	}

	os.Exit(m.Run())
}
