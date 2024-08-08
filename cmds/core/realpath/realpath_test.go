// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestRealpath(t *testing.T) {
	dir := t.TempDir()
	dir, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.CreateTemp(dir, "")
	if err != nil {
		t.Fatal(err)
	}

	err = os.Symlink(f.Name(), filepath.Join(dir, "symlink"))
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		expectedErr    error
		expectedOutput string
		args           []string
	}{
		{
			expectedOutput: f.Name() + "\n",
			args:           []string{f.Name()},
		},
		{
			expectedOutput: f.Name() + "\n",
			args:           []string{filepath.Join(dir, "symlink")},
		},
		{
			expectedOutput: f.Name() + "\n",
			args:           []string{dir + "/../" + filepath.Base(dir) + "/symlink"},
		},
		{
			args:        []string{"filenotexists"},
			expectedErr: os.ErrNotExist,
		},
	}

	for _, test := range tests {
		stdout := &bytes.Buffer{}

		err := run(stdout, test.args...)
		fmt.Println(test.args)
		if !errors.Is(err, test.expectedErr) {
			t.Fatalf("expected %v, got %v", test.expectedErr, err)
		}

		if stdout.String() != test.expectedOutput {
			t.Errorf("expected %v, got %v", test.expectedOutput, stdout.String())
		}
	}
}
