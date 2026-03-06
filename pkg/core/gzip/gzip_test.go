// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gzip

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestGzipWithKeep(t *testing.T) {
	tmpDir := t.TempDir()
	// Change to the temp directory to avoid path resolution issues
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWd)

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	filePath := "file.txt"
	wantContent := []byte("test file's content\nsecond line")
	err = os.WriteFile(filePath, wantContent, 0o644)
	if err != nil {
		t.Fatalf("os.WriteFile(%v, %v, 0o644) = %v, want nil", filePath, string(wantContent), err)
	}

	// Create a new gzip command
	gzip := New().(*Gzip)

	// Compress the file with force and keep flags
	err = gzip.Run("gzip", "-f", "-k", filePath)
	if err != nil {
		t.Fatalf("gzip.Run(gzip, -f, -k, %v) = %v, want nil", filePath, err)
	}

	// Check that the compressed file exists
	compressedPath := filePath + ".gz"
	_, err = os.Stat(compressedPath)
	if err != nil {
		t.Fatalf("os.Stat(%v) = %v, want nil", compressedPath, err)
	}

	// Check that the original file still exists (using -k)
	_, err = os.Stat(filePath)
	if err != nil {
		t.Fatalf("os.Stat(%v) = %v, want nil", filePath, err)
	}

	// Decompress the file with force and keep flags
	err = gzip.Run("gzip", "-d", "-f", "-k", compressedPath)
	if err != nil {
		t.Fatalf("gzip.Run(gzip, -d, -f, -k, %v) = %v, want nil", compressedPath, err)
	}

	// Check that both files exist
	_, err = os.Stat(filePath)
	if err != nil {
		t.Fatalf("os.Stat(%v) = %v, want nil", filePath, err)
	}
	_, err = os.Stat(compressedPath)
	if err != nil {
		t.Fatalf("os.Stat(%v) = %v, want nil", compressedPath, err)
	}

	// Verify the content
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("os.ReadFile(%v) = %v, want nil", filePath, err)
	}

	if !bytes.Equal(content, wantContent) {
		t.Errorf("os.ReadFile(%v) = %v, want %v", filePath, string(content), string(wantContent))
	}
}

func TestGzipWithWorkingDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a subdirectory
	subDir := filepath.Join(tmpDir, "subdir")
	err := os.Mkdir(subDir, 0o755)
	if err != nil {
		t.Fatalf("os.Mkdir(%v, 0o755) = %v, want nil", subDir, err)
	}

	// Create a file in the subdirectory
	filePath := filepath.Join(subDir, "file.txt")
	wantContent := []byte("test file in subdirectory")
	err = os.WriteFile(filePath, wantContent, 0o644)
	if err != nil {
		t.Fatalf("os.WriteFile(%v, %v, 0o644) = %v, want nil", filePath, string(wantContent), err)
	}

	// Create a new gzip command with working directory set to tmpDir
	gzip := New().(*Gzip)
	gzip.SetWorkingDir(tmpDir)

	// Compress the file using a relative path with force and keep flags
	err = gzip.Run("gzip", "-f", "-k", "subdir/file.txt")
	if err != nil {
		t.Fatalf("gzip.Run(gzip, -f, -k, subdir/file.txt) = %v, want nil", err)
	}

	// Check that the compressed file exists
	compressedPath := filePath + ".gz"
	_, err = os.Stat(compressedPath)
	if err != nil {
		t.Fatalf("os.Stat(%v) = %v, want nil", compressedPath, err)
	}

	// Check that the original file still exists (using -k)
	_, err = os.Stat(filePath)
	if err != nil {
		t.Fatalf("os.Stat(%v) = %v, want nil", filePath, err)
	}

	// Decompress the file with force and keep flags
	err = gzip.Run("gzip", "-d", "-f", "-k", "subdir/file.txt.gz")
	if err != nil {
		t.Fatalf("gzip.Run(gzip, -d, -f, -k, subdir/file.txt.gz) = %v, want nil", err)
	}

	// Verify the content
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("os.ReadFile(%v) = %v, want nil", filePath, err)
	}

	if !bytes.Equal(content, wantContent) {
		t.Errorf("os.ReadFile(%v) = %v, want %v", filePath, string(content), string(wantContent))
	}
}
