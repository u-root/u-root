// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Note from author, xplshn: I am not going to lie; I don't know how to write a test for Go. I think this covers it, but I may be doing random BS here.
package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// setup creates a temporary file with the given content for testing.
func setup(t *testing.T, content string) string {
	t.Helper()
	t.Logf(":: Creating simulation data.")
	dir := t.TempDir()
	filePath := filepath.Join(dir, "testfile")
	if err := os.WriteFile(filePath, []byte(content), 0o666); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	return filePath
}

// mockIssueSequences sets the mock values for the issue sequences map.
func mockIssueSequences() map[string]string {
	return map[string]string{
		"\\H": "Welcome To The Machine",
		"\\d": time.Now().Format("Monday, 02 Jan 2006"),
		"\\t": time.Now().Format("15:04:05"),
	}
}

// TestIssue verifies the main functionality of the issue command.
func TestIssue(t *testing.T) {
	content := "Welcome to \\H\nCurrent date: \\d\nCurrent time: \\t\n"
	filePath := setup(t, content)

	mockSequences := mockIssueSequences()
	expectedOutput := fmt.Sprintf("Welcome to %s\nCurrent date: %s\nCurrent time: %s\n",
		mockSequences["\\H"], mockSequences["\\d"], mockSequences["\\t"])

	// Capture the output
	var out bytes.Buffer
	run(&out, filePath, mockSequences)

	if out.String() != expectedOutput {
		t.Errorf("Expected %q but got %q", expectedOutput, out.String())
	}
}

// run executes the main function with a given file and sequences map, capturing output.
func run(out *bytes.Buffer, filePath string, sequences map[string]string) {
	// Backup original stdout and replace with buffer
	origStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Set the issueSequences map to mock values for the test
	issueSequences = sequences

	// Run the main function
	os.Args = []string{"issue", filePath}
	main()

	// Restore stdout and capture the output
	w.Close()
	_, _ = out.ReadFrom(r)
	os.Stdout = origStdout
}

func TestFileNotFound(t *testing.T) {
	var out bytes.Buffer
	run(&out, "nonexistentfile", mockIssueSequences())

	expected := "Error opening file: open nonexistentfile: no such file or directory\n"
	if out.String() != expected {
		t.Errorf("Expected %q but got %q", expected, out.String())
	}
}

func TestNoArgs(t *testing.T) {
	content := "Welcome to \\H\n"
	filePath := setup(t, content)

	mockSequences := mockIssueSequences()
	expectedOutput := fmt.Sprintf("Welcome to %s\n", mockSequences["\\H"])

	var out bytes.Buffer
	run(&out, filePath, mockSequences)

	if out.String() != expectedOutput {
		t.Errorf("Expected %q but got %q", expectedOutput, out.String())
	}
}

func TestEmptyFile(t *testing.T) {
	filePath := setup(t, "")

	mockSequences := mockIssueSequences()
	expectedOutput := ""

	var out bytes.Buffer
	run(&out, filePath, mockSequences)

	if out.String() != expectedOutput {
		t.Errorf("Expected %q but got %q", expectedOutput, out.String())
	}
}

func TestCustomSequences(t *testing.T) {
	content := "This is a custom sequence: \\u\n"
	filePath := setup(t, content)

	mockSequences := mockIssueSequences()
	mockSequences["\\u"] = "Custom User"

	expectedOutput := "This is a custom sequence: Custom User\n"

	var out bytes.Buffer
	run(&out, filePath, mockSequences)

	if out.String() != expectedOutput {
		t.Errorf("Expected %q but got %q", expectedOutput, out.String())
	}
}
