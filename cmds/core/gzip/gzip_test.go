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

func TestGZIPSmoke(t *testing.T) {
	dir := t.TempDir()
	filePath := dir + "/file.txt"
	wantContent := []byte("test file's content\nsecond line")
	err := os.WriteFile(filePath, wantContent, 0o644)
	if err != nil {
		t.Fatalf("os.WriteFile(%v, %v, 0o644) = %v, want nil", filePath, string(wantContent), err)
	}

	err = run([]string{"-9", "-b", "128", "-p", "1", filePath})
	if err != nil {
		t.Fatalf("run(%v) = %v, want nil", []string{dir + "/file.txt"}, err)
	}

	err = run([]string{"-d", filePath + ".gz"})
	if err != nil {
		t.Fatalf("run(%v) = %v, want nil", []string{dir + "/file.txt.gz"}, err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("os.ReadFile(%v) = %v, want nil", filePath, err)
	}

	if !bytes.Equal(content, wantContent) {
		t.Errorf("os.ReadFile(%v) = %v, want %v", filePath, string(content), string(wantContent))
	}
}

func TestGZIPExistingOutput(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(in, []byte("content"), 0o644); err != nil {
		t.Fatalf("os.WriteFile(%v) = %v, want nil", in, err)
	}
	if err := run([]string{in}); err != nil {
		t.Fatalf("run(%v) = %v, want nil", []string{in}, err)
	}

	// Someone else's symlink where the decompressed file would go.
	target := filepath.Join(t.TempDir(), "target")
	if err := os.Symlink(target, in); err != nil {
		t.Fatalf("os.Symlink(%v, %v) = %v, want nil", target, in, err)
	}
	args := []string{"-d", in + ".gz"}
	if err := run(args); err != nil {
		t.Fatalf("run(%v) = %v, want nil", args, err)
	}
	if _, err := os.Lstat(target); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("after run(%v): os.Lstat(%v) = %v, want %v", args, target, err, os.ErrNotExist)
	}
	if _, err := os.Stat(in + ".gz"); err != nil {
		t.Errorf("after run(%v): os.Stat(%v) = %v, want nil", args, in+".gz", err)
	}

	// If the output can't be created even with -f, the input stays.
	if err := os.Remove(in); err != nil {
		t.Fatalf("os.Remove(%v) = %v, want nil", in, err)
	}
	if err := os.Mkdir(in, 0o755); err != nil {
		t.Fatalf("os.Mkdir(%v) = %v, want nil", in, err)
	}
	args = []string{"-d", "-f", in + ".gz"}
	if err := run(args); err != nil {
		t.Fatalf("run(%v) = %v, want nil", args, err)
	}
	if _, err := os.Stat(in + ".gz"); err != nil {
		t.Errorf("after run(%v): os.Stat(%v) = %v, want nil", args, in+".gz", err)
	}
}
