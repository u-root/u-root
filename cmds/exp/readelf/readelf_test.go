// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	_ "embed"
)

//go:embed testdata/a
var elfFile []byte

//go:embed testdata/out
var goodOut string

func TestRead(t *testing.T) {
	d := t.TempDir()
	f := filepath.Join(d, "data")

	if err := run(nil, io.Discard, f); err == nil {
		t.Errorf("Opening nonexistent file: got nil, want err")
	}

	if err := os.WriteFile(f, elfFile, 0o666); err != nil {
		t.Fatalf("Writing data: %v", err)
	}

	var out bytes.Buffer
	if err := run(nil, &out, f); err != nil {
		t.Errorf("reading good elfFile: got %v, want nil", err)
	}

	if out.Len() != len(goodOut) {
		t.Fatalf("elf output: read %d bytes, want %d bytes", out.Len(), len(goodOut))
	}

	if out.String() != goodOut {
		t.Fatalf("elf output: want %q, got %q", out.String(), goodOut)
	}

	// short file
	if err := os.WriteFile(f, elfFile[:64], 0o666); err != nil {
		t.Fatalf("Writing data: %v", err)
	}

	if err := run(nil, io.Discard, f); err == nil {
		t.Errorf("Opening short file: got nil, want err")
	}

	// corrupt header
	elfFile[1] = ^elfFile[1]
	if err := os.WriteFile(f, elfFile, 0o666); err != nil {
		t.Fatalf("Writing data: %v", err)
	}

	if err := run(nil, io.Discard, f); err == nil {
		t.Errorf("Opening corrupt file: got nil, want err")
	}
}
