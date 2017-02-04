// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testutil

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"testing"
)

// CompileInTempDir creates a temp directory and compiles the main package of
// the current directory. Remember to delete the directory after the test:
//     defer os.RemoveAll(tmpDir)
func CompileInTempDir(t *testing.T) (tmpDir string, execPath string) {
	// Create temp directory
	tmpDir, err := ioutil.TempDir("", "Test")
	if err != nil {
		t.Fatal("TempDir failed: ", err)
	}

	// Compile the program
	execPath = filepath.Join(tmpDir, "exec")
	out, err := exec.Command("go", "build", "-o", execPath).CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build: %v\n%s", err, string(out))
	}
	return
}
