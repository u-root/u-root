// Copyright 2013-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9 && !windows

package cpio

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateFileInRoot(t *testing.T) {
	tmp := t.TempDir()
	fileName := "file"
	content := "content"
	r := StaticFile(fileName, content, 0o644)
	err := CreateFileInRoot(r, tmp, false)
	if err != nil {
		t.Fatalf("expected nil got %v", err)
	}

	b, err := os.ReadFile(filepath.Join(tmp, "file"))
	if err != nil {
		t.Fatalf("expected nil got %v", err)
	}

	if string(b) != content {
		t.Errorf("expected %q got %q", content, string(b))
	}
}
