// Copyright 2009-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	dir := t.TempDir()
	f, err := os.CreateTemp(dir, "file")
	if err != nil {
		t.Fatalf("expected nil got %v", err)
	}

	stdout := &bytes.Buffer{}

	err = run(stdout, []string{f.Name()})
	if err != nil {
		t.Fatalf("expected nil got %v", err)
	}

	expected := f.Name() + "\n"

	if stdout.String() != expected {
		t.Errorf("expected %q, got %q", expected, stdout.String())
	}
}
