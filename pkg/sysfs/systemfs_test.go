package sysfs

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestSystemFS_Open(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()

	// Create test files
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0o644); err != nil {
		t.Fatal(err)
	}

	parentFile := filepath.Join(filepath.Dir(tmpDir), "parent.txt")
	if err := os.WriteFile(parentFile, []byte("parent content"), 0o644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(parentFile)

	sfs := NewSystemFS(tmpDir)

	// Test normal file access
	f, err := sfs.Open("test.txt")
	if err != nil {
		t.Fatalf("Failed to open test.txt: %v", err)
	}
	f.Close()

	// Test access to parent directory (this should work with SystemFS)
	f, err = sfs.Open("../parent.txt")
	if err != nil {
		t.Fatalf("Failed to open ../parent.txt: %v", err)
	}
	f.Close()
}

func TestSystemFS_ReadFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test file
	testContent := "test content"
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte(testContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create parent file
	parentContent := "parent content"
	parentFile := filepath.Join(filepath.Dir(tmpDir), "parent.txt")
	if err := os.WriteFile(parentFile, []byte(parentContent), 0o644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(parentFile)

	sfs := NewSystemFS(tmpDir)

	// Test normal file read
	content, err := sfs.ReadFile("test.txt")
	if err != nil {
		t.Fatalf("Failed to read test.txt: %v", err)
	}
	if string(content) != testContent {
		t.Errorf("Expected %q, got %q", testContent, string(content))
	}

	// Test parent directory access
	content, err = sfs.ReadFile("../parent.txt")
	if err != nil {
		t.Fatalf("Failed to read ../parent.txt: %v", err)
	}
	if string(content) != parentContent {
		t.Errorf("Expected %q, got %q", parentContent, string(content))
	}
}

func TestSystemFS_Stat(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0o644); err != nil {
		t.Fatal(err)
	}

	sfs := NewSystemFS(tmpDir)

	info, err := sfs.Stat("test.txt")
	if err != nil {
		t.Fatalf("Failed to stat test.txt: %v", err)
	}

	if info.Name() != "test.txt" {
		t.Errorf("Expected name test.txt, got %s", info.Name())
	}
}

func TestSystemFS_ReadDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	if err := os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("1"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "file2.txt"), []byte("2"), 0o644); err != nil {
		t.Fatal(err)
	}

	sfs := NewSystemFS(tmpDir)

	entries, err := sfs.ReadDir(".")
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}

	if len(entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(entries))
	}
}

func TestSystemFS_Sub(t *testing.T) {
	tmpDir := t.TempDir()

	// Create subdirectory with file
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0o755); err != nil {
		t.Fatal(err)
	}

	testFile := filepath.Join(subDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("sub content"), 0o644); err != nil {
		t.Fatal(err)
	}

	sfs := NewSystemFS(tmpDir)

	subFS, err := sfs.Sub("subdir")
	if err != nil {
		t.Fatalf("Failed to create sub filesystem: %v", err)
	}

	content, err := fs.ReadFile(subFS, "test.txt")
	if err != nil {
		t.Fatalf("Failed to read from sub filesystem: %v", err)
	}

	if string(content) != "sub content" {
		t.Errorf("Expected 'sub content', got %q", string(content))
	}
}

func TestSystemFS_InvalidPaths(t *testing.T) {
	tmpDir := t.TempDir()
	sfs := NewSystemFS(tmpDir)

	invalidPaths := []string{
		"",
	}

	for _, path := range invalidPaths {
		_, err := sfs.Open(path)
		if err == nil {
			t.Errorf("Expected error for invalid path %q, but got none", path)
		}
	}
}
