// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestMore(t *testing.T) {
	t.Run("files is not exist", func(t *testing.T) {
		err := run(nil, nil, 40, []string{"file-is-not-exists"})
		if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("expected %v, got %v", os.ErrNotExist, err)
		}
	})
	t.Run("negative lines", func(t *testing.T) {
		err := run(nil, nil, -1, []string{"file1"})
		if !errors.Is(err, errLinesMustBePositive) {
			t.Errorf("expected %v, got %v", errLinesMustBePositive, err)
		}
	})
	t.Run("one screen file", func(t *testing.T) {
		content := "line1\nline2\nline3\n"
		path := filepath.Join(t.TempDir(), "file1")
		err := os.WriteFile(path, []byte(content), 0o644)
		if err != nil {
			t.Fatalf("failed to write file: %v", err)
		}

		stdout := &bytes.Buffer{}

		err = run(nil, stdout, 10, []string{path})
		if err != nil {
			t.Fatalf("failed to run more: %v", err)
		}

		if stdout.String() != content {
			t.Errorf("expected %q, got %q", content, stdout.String())
		}
	})
	t.Run("two screen file", func(t *testing.T) {
		content := "line1\nline2\nline3\nline4\nline5\n"
		// third '\n' missing, 'enter' is emulated by stdin buffer
		// '\n' comes from 'enter'
		expectedOutput := "line1\nline2\nline3line4\nline5\n"
		path := filepath.Join(t.TempDir(), "file1")
		err := os.WriteFile(path, []byte(content), 0o644)
		if err != nil {
			t.Fatalf("failed to write file: %v", err)
		}

		stdin := bytes.NewBufferString("a")
		stdout := &bytes.Buffer{}

		if err := run(stdin, stdout, 3, []string{path}); err != nil {
			t.Fatalf("failed to run more: %v", err)
		}

		if stdout.String() != expectedOutput {
			t.Errorf("expected %q, got %q", expectedOutput, stdout.String())
		}
	})
}
