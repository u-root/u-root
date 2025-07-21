// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCRC(t *testing.T) {
	t.Run("test unknown function", func(t *testing.T) {
		err := run(nil, nil, "unknown", nil)
		if !errors.Is(err, os.ErrInvalid) {
			t.Errorf("expected %v, got %v", os.ErrInvalid, err)
		}
	})

	t.Run("test stdin", func(t *testing.T) {
		stdin := strings.NewReader("test\n")
		stdout := &bytes.Buffer{}
		err := run(stdin, stdout, "crc32-ieee", nil)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		expected := "3bb935c6\n"
		if stdout.String() != expected {
			t.Errorf("expected = %q, want %q", stdout.String(), expected)
		}
	})

	t.Run("test file", func(t *testing.T) {
		tmp := t.TempDir()
		path := filepath.Join(tmp, "file")
		err := os.WriteFile(path, []byte("test\n"), 0o644)
		if err != nil {
			t.Fatalf("cannot create file: %v", err)
		}

		stdout := &bytes.Buffer{}
		err = run(nil, stdout, "crc32-ieee", []string{path})
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}

		expected := "3bb935c6\n"
		if stdout.String() != expected {
			t.Errorf("expected = %q, want %q", stdout.String(), expected)
		}
	})
}
