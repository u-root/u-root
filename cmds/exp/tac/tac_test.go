// Copyright 2012-2023 the u-root Authors. All rights reserved
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

func TestTac(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tac1")
	err := os.WriteFile(path, []byte("hello\nworld\n"), 0o644)
	if err != nil {
		t.Fatalf(`os.WriteFile(%q, []byte("hello\nworld\n"), 0644) = %v, want nil`, path, err)
	}

	stdout := &bytes.Buffer{}
	err = tac(stdout, []string{path})
	if err != nil {
		t.Fatalf(`tac(stdout, []string{f.Name(), f.Name()}) = %v, want nil`, err)
	}

	expected := "world\nhello\n"
	if stdout.String() != expected {
		t.Errorf("expected %s, got %s", expected, stdout.String())
	}
}

func TestTacStdin(t *testing.T) {
	err := tac(nil, nil)
	if !errors.Is(err, errStdin) {
		t.Errorf("expected %v, got %v", errStdin, err)
	}
}
