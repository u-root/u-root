// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testutil

import (
	"bytes"
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

// ErrorExists is only used to increase readability of future tests
func ErrorExists(err error) bool {
	return err != nil
}

// Helper function for PrintError
func craftPrintMsg(errExists bool, out string) string {
	var msg bytes.Buffer

	if errExists {
		msg.WriteString("Error Status: exists\n")
	} else {
		msg.WriteString("Error Status: not exists\n")
	}
	msg.WriteString("Output:\n")
	msg.WriteString(out)
	return msg.String()
}

// PrintError provides a standard way to print out error message when a test case fails
func PrintError(t *testing.T, funcCallStm string, expOut string, expErrExists bool, actualOut string, actualErr error) {
	expectMsg := craftPrintMsg(expErrExists, expOut)
	actualMsg := craftPrintMsg(ErrorExists(actualErr), actualOut)
	t.Errorf("%s\ngot:\n%s\n\nwant:\n%s", funcCallStm, actualMsg, expectMsg)
}
