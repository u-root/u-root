// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mkfifo

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

type test struct {
	name  string
	paths []string
	mode  uint32
	err   error
}

func TestMkfifo(t *testing.T) {
	tmpDir, _ := testutil.CompileInTempDir(t)
	defer os.RemoveAll(tmpDir)

	// used later in testing
	testDir := filepath.Join(tmpDir, "mkfifoDir")
	err := os.Mkdir(testDir, 0700)
	if err != nil {
		t.Error(err)
	}

	var tests = []test{
		{
			name:  "no path or mode, no errors",
			paths: []string{},
			mode:  syscall.S_IFIFO,
			err:   nil,
		},
		{
			name:  "single path",
			paths: []string{filepath.Join(testDir, "testfifo")},
			mode:  syscall.S_IFIFO,
			err:   nil,
		},
		{
			name:  "duplicate path",
			paths: []string{filepath.Join(testDir, "testfifo1"), filepath.Join(testDir, "testfifo1")},
			mode:  syscall.S_IFIFO,
			err:   errors.New("file exists"),
		},
	}

	for _, tt := range tests {
		mk := Mkfifo{Paths: tt.paths, Mode: tt.mode}
		err := mk.Exec()

		if err != tt.err && err.Error() != tt.err.Error() {
			t.Errorf("expected %v got %v", err, tt.err)
		}

		// This might be overkill
		for _, path := range tt.paths {
			testFile, err := os.Lstat(path)

			if err != nil {
				t.Errorf("Unable to stat file %s", path)
			}

			mode := testFile.Mode()
			if mode.Perm() != os.FileMode(tt.mode).Perm() {
				t.Errorf("mode incorrect. expected %v got %v", mode, tt.mode)
			}
		}
	}
}
