// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testutil

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// CompileInTempDir creates a temp directory and compiles the main package of
// the current directory. Remember to delete the directory after the test:
//
//     defer os.RemoveAll(tmpDir)
//
// The first argument of the environment variable EXECPATH overrides execPath.
func CompileInTempDir(t testing.TB) (tmpDir string, execPath string) {
	// Create temp directory
	tmpDir, err := ioutil.TempDir("", "Test")
	if err != nil {
		t.Fatal("TempDir failed: ", err)
	}

	// Skip compilation if EXECPATH is set.
	execPath = os.Getenv("EXECPATH")
	if execPath != "" {
		execPath = strings.SplitN(execPath, " ", 2)[0]
		return
	}

	// Compile the program
	execPath = filepath.Join(tmpDir, "exec")
	out, err := exec.Command("go", "build", "-o", execPath).CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build: %v\n%s", err, string(out))
	}
	return
}
