// Copyright 2012-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPWD(t *testing.T) {
	dir := t.TempDir()
	err := os.Chdir(dir)
	if err != nil {
		t.Fatal(err)
	}

	rs, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatal(err)
	}

	r, err := pwd(true)
	if err != nil {
		t.Fatal(err)
	}

	if r != rs {
		t.Errorf("expected: %q, got %q", rs, r)
	}
}
