// Copyright 2012-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
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
