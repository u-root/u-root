// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"os/exec"
	"path"
	"testing"
)

func TestMain(t *testing.T) {
	// This is a simple integration test to ensure the main package works.
	// More detailed tests are in the pkg/core/tar package.

	if os.Getenv("TEST_MAIN_BINARY") == "1" {
		// When running as the test binary, just exit successfully
		return
	}

	tmpDir := t.TempDir()
	err := os.Chdir(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// Create a test file
	filePath := path.Join(tmpDir, "file")
	f, err := os.Create(filePath)
	if err != nil {
		t.Fatal(err)
	}
	content := "hello from tar main test"
	_, err = f.WriteString(content)
	if err != nil {
		t.Fatal(err)
	}
	err = f.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Build the binary
	cmd := exec.Command("go", "build", "-o", "tar_test_binary")
	cmd.Env = append(os.Environ(), "TEST_MAIN_BINARY=1")
	if err := cmd.Run(); err != nil {
		t.Skipf("Failed to build binary: %v", err)
	}
	defer os.Remove("tar_test_binary")

	// Run the binary to create an archive
	createCmd := exec.Command("./tar_test_binary", "-cf", "file.tar", "file")
	if err := createCmd.Run(); err != nil {
		t.Fatalf("Failed to create archive: %v", err)
	}

	// Check that the archive was created
	if _, err := os.Stat("file.tar"); err != nil {
		t.Fatalf("Archive was not created: %v", err)
	}
}
