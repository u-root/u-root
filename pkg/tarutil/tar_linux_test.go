// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tarutil

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/sys/unix"
)

// TestCreateTarManyFiles archives more files than the file descriptor limit
// allows to be open at once. It is a regression test for CreateTar leaking
// one descriptor per archived file.
func TestCreateTarManyFiles(t *testing.T) {
	dir := t.TempDir()
	for i := range 80 {
		if err := os.WriteFile(filepath.Join(dir, fmt.Sprintf("file-%03d", i)), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	var old unix.Rlimit
	if err := unix.Getrlimit(unix.RLIMIT_NOFILE, &old); err != nil {
		t.Fatal(err)
	}
	limit := old
	limit.Cur = min(limit.Cur, 64)
	if err := unix.Setrlimit(unix.RLIMIT_NOFILE, &limit); err != nil {
		t.Fatal(err)
	}
	defer unix.Setrlimit(unix.RLIMIT_NOFILE, &old)

	// CreateTar only derives in-archive names for paths relative to
	// ChangeDirectory, so archive the temporary directory like
	// "tar -C / DIR".
	var archive bytes.Buffer
	relDir := strings.TrimPrefix(dir, string(filepath.Separator))
	if err := CreateTar(&archive, []string{relDir}, &Opts{ChangeDirectory: string(filepath.Separator)}); err != nil {
		t.Fatalf("CreateTar(%q) = %v, want nil", relDir, err)
	}
}
