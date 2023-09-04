// Copyright 2012-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/u-root/u-root/pkg/gzip"
)

func TestGZIPSmoke(t *testing.T) {
	dir := t.TempDir()
	filePath := dir + "/file.txt"
	wantContent := []byte("test file's content\nsecond line")
	err := os.WriteFile(filePath, wantContent, 0o644)
	if err != nil {
		t.Fatalf("os.WriteFile(%v, %v, 0o644) = %v, want nil", filePath, string(wantContent), err)
	}

	opts := gzip.Options{
		Blocksize: 128,
		Suffix:    ".gz",
		Level:     9,
		Processes: 1,
	}

	err = run(opts, []string{filePath})
	if err != nil {
		t.Fatalf("run(%v, %v) = %v, want nil", opts, []string{dir + "/file.txt"}, err)
	}

	opts.Decompress = true

	err = run(opts, []string{filePath + ".gz"})
	if err != nil {
		t.Fatalf("run(%v, %v) = %v, want nil", opts, []string{dir + "/file.txt.gz"}, err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("os.ReadFile(%v) = %v, want nil", filePath, err)
	}

	if !bytes.Equal(content, wantContent) {
		t.Errorf("os.ReadFile(%v) = %v, want %v", filePath, string(content), string(wantContent))
	}
}
