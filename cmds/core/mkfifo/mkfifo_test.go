// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

type test struct {
	name   string
	flags  []string
	stdErr string
}

func TestMkfifo(t *testing.T) {
	tmpDir := t.TempDir()

	// used later in testing
	testDir := filepath.Join(tmpDir, "mkfifoDir")
	if err := os.Mkdir(testDir, 0o700); err != nil {
		t.Error(err)
	}

	tests := []test{
		{
			name:   "no path or mode, error",
			flags:  []string{},
			stdErr: "please provide a path, or multiple, to create a fifo",
		},
		{
			name:   "single path",
			flags:  []string{filepath.Join(testDir, "testfifo")},
			stdErr: "",
		},
		{
			name:   "duplicate path",
			flags:  []string{filepath.Join(testDir, "testfifo1"), filepath.Join(testDir, "testfifo1")},
			stdErr: "file exists",
		},
	}

	for _, tt := range tests {
		var out, stdErr bytes.Buffer
		cmd := testutil.Command(t, tt.flags...)
		cmd.Stdout = &out
		cmd.Stderr = &stdErr
		err := cmd.Run()

		if err != nil && !strings.Contains(stdErr.String(), tt.stdErr) {
			t.Errorf("expected %v got %v", tt.stdErr, stdErr.String())
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
	testutil.Run(m, main)
}
