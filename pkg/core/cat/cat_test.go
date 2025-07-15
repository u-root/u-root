// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// by Rafael Campos Nunes <rafaelnunes@engineer.com>

package cat

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"
	"testing/iotest"
)

// setup writes a set of files, putting 1 byte in each file.
func setup(t *testing.T, data []byte) string {
	t.Helper()
	t.Logf(":: Creating simulation data. ")
	dir := t.TempDir()

	for i, d := range data {
		n := fmt.Sprintf("%v%d", filepath.Join(dir, "file"), i)
		if err := os.WriteFile(n, []byte{d}, 0o666); err != nil {
			t.Fatal(err)
		}
	}

	return dir
}

// TestCat test cat function against 4 files, in each file it is written a bit of someData
// array and the test expect the cat to return the exact same bit from someData array with
// the corresponding file.
func TestCat(t *testing.T) {
	var files []string
	someData := []byte{'l', 2, 3, 4, 'd'}

	dir := setup(t, someData)

	for i := range someData {
		files = append(files, fmt.Sprintf("%v%d", filepath.Join(dir, "file"), i))
	}

	cmd := New()
	var stdout, stderr bytes.Buffer
	var stdin bytes.Buffer
	cmd.SetIO(&stdin, &stdout, &stderr)

	exitCode, err := cmd.Run(context.Background(), files...)
	if err != nil {
		t.Fatal(err)
	}
	if exitCode != 0 {
		t.Fatalf("Expected exit code 0, got %d", exitCode)
	}

	if !bytes.Equal(stdout.Bytes(), someData) {
		t.Fatalf("Reading files failed: got %v, want %v", stdout.Bytes(), someData)
	}
}

func TestCatPipe(t *testing.T) {
	cmd := New().(*Command) // Type assertion to access internal methods
	var inputbuf bytes.Buffer
	teststring := "testdata"
	fmt.Fprintf(&inputbuf, "%s", teststring)

	var stdout, stderr bytes.Buffer
	cmd.SetIO(&inputbuf, &stdout, &stderr)

	if err := cmd.cat(&inputbuf, &stdout); err != nil {
		t.Error(err)
	}
	if stdout.String() != teststring {
		t.Errorf("CatPipe: Want %q Got: %q", teststring, stdout.String())
	}
}

func TestRunFiles(t *testing.T) {
	var files []string
	someData := []byte{'l', 2, 3, 4, 'd'}

	dir := setup(t, someData)

	for i := range someData {
		files = append(files, fmt.Sprintf("%v%d", filepath.Join(dir, "file"), i))
	}

	cmd := New()
	var stdout, stderr bytes.Buffer
	var stdin bytes.Buffer
	cmd.SetIO(&stdin, &stdout, &stderr)

	exitCode, err := cmd.Run(context.Background(), files...)
	if err != nil {
		t.Error(err)
	}
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}
	if !bytes.Equal(stdout.Bytes(), someData) {
		t.Fatalf("Reading files failed: got %v, want %v", stdout.Bytes(), someData)
	}
}

func TestRunFilesError(t *testing.T) {
	var files []string
	someData := []byte{'l', 2, 3, 4, 'd'}

	dir := setup(t, someData)

	for i := range someData {
		files = append(files, fmt.Sprintf("%v%d", filepath.Join(dir, "file"), i))
	}
	filenotexist := "testdata/doesnotexist.txt"
	files = append(files, filenotexist)

	cmd := New()
	var stdout, stderr bytes.Buffer
	var stdin bytes.Buffer
	cmd.SetIO(&stdin, &stdout, &stderr)

	exitCode, err := cmd.Run(context.Background(), files...)
	if err == nil {
		t.Error("function run succeeded but should have failed")
	}
	if exitCode == 0 {
		t.Error("Expected non-zero exit code")
	}
}

func TestRunNoArgs(t *testing.T) {
	cmd := New()
	var stdout, stderr bytes.Buffer
	var stdin bytes.Buffer
	inputdata := "teststring"
	fmt.Fprintf(&stdin, "%s", inputdata)
	cmd.SetIO(&stdin, &stdout, &stderr)

	exitCode, err := cmd.Run(context.Background())
	if err != nil {
		t.Error(err)
	}
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}
	if stdout.String() != inputdata {
		t.Errorf("Want: %q Got: %q", inputdata, stdout.String())
	}
}

func TestIOErrors(t *testing.T) {
	cmd := New()
	var stdout, stderr bytes.Buffer
	errReader := iotest.ErrReader(errors.New("read error"))
	cmd.SetIO(errReader, &stdout, &stderr)

	exitCode, err := cmd.Run(context.Background())
	if !errors.Is(err, errCopy) {
		t.Errorf("expected %v, got %v", errCopy, err)
	}
	if exitCode == 0 {
		t.Error("Expected non-zero exit code")
	}

	// Test with dash argument
	cmd2 := New()
	var stdout2, stderr2 bytes.Buffer
	cmd2.SetIO(errReader, &stdout2, &stderr2)

	exitCode, err = cmd2.Run(context.Background(), "-")
	if !errors.Is(err, errCopy) {
		t.Errorf("expected %v, got %v", errCopy, err)
	}
	if exitCode == 0 {
		t.Error("Expected non-zero exit code")
	}
}

func TestCatDash(t *testing.T) {
	tempDir := t.TempDir()

	f1 := path.Join(tempDir, "f1")
	err := os.WriteFile(f1, []byte("line1\nline2\n"), 0o666)
	if err != nil {
		t.Fatal(err)
	}

	f2 := path.Join(tempDir, "f2")
	err = os.WriteFile(f2, []byte("line4\nline5\n"), 0o666)
	if err != nil {
		t.Fatal(err)
	}

	cmd := New()
	var stdout, stderr bytes.Buffer
	var stdin bytes.Buffer
	stdin.WriteString("line3\n")
	cmd.SetIO(&stdin, &stdout, &stderr)

	exitCode, err := cmd.Run(context.Background(), f1, "-", f2)
	if err != nil {
		t.Fatal(err)
	}
	if exitCode != 0 {
		t.Fatalf("Expected exit code 0, got %d", exitCode)
	}

	want := "line1\nline2\nline3\nline4\nline5\n"
	got := stdout.String()

	if got != want {
		t.Errorf("want: %s, got: %s", want, got)
	}
}

func TestCatWorkingDir(t *testing.T) {
	tempDir := t.TempDir()

	// Create a file in the temp directory
	testFile := "test.txt"
	testContent := "test content"
	err := os.WriteFile(filepath.Join(tempDir, testFile), []byte(testContent), 0o666)
	if err != nil {
		t.Fatal(err)
	}

	cmd := New()
	var stdout, stderr bytes.Buffer
	var stdin bytes.Buffer
	cmd.SetIO(&stdin, &stdout, &stderr)
	cmd.SetWorkingDir(tempDir)

	exitCode, err := cmd.Run(context.Background(), testFile)
	if err != nil {
		t.Fatal(err)
	}
	if exitCode != 0 {
		t.Fatalf("Expected exit code 0, got %d", exitCode)
	}

	if stdout.String() != testContent {
		t.Errorf("Want: %q Got: %q", testContent, stdout.String())
	}
}
