// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// by Rafael Campos Nunes <rafaelnunes@engineer.com>

package main

import (
	"bytes"
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
	var out bytes.Buffer
	if err := run(nil, &out, files...); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(out.Bytes(), someData) {
		t.Fatalf("Reading files failed: got %v, want %v", out.Bytes(), someData)
	}
}

func TestCatPipe(t *testing.T) {
	var inputbuf bytes.Buffer
	teststring := "testdata"
	fmt.Fprintf(&inputbuf, "%s", teststring)

	var out bytes.Buffer

	if err := cat(&inputbuf, &out); err != nil {
		t.Error(err)
	}
	if out.String() != teststring {
		t.Errorf("CatPipe: Want %q Got: %q", teststring, out.String())
	}
}

func TestRunFiles(t *testing.T) {
	var files []string
	someData := []byte{'l', 2, 3, 4, 'd'}

	dir := setup(t, someData)

	for i := range someData {
		files = append(files, fmt.Sprintf("%v%d", filepath.Join(dir, "file"), i))
	}

	var out bytes.Buffer
	if err := run(nil, &out, files...); err != nil {
		t.Error(err)
	}
	if !bytes.Equal(out.Bytes(), someData) {
		t.Fatalf("Reading files failed: got %v, want %v", out.Bytes(), someData)
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
	var in, out bytes.Buffer
	if err := run(&in, &out, files...); err == nil {
		t.Error("function run succeeded but should have failed")
	}
}

func TestRunNoArgs(t *testing.T) {
	var in, out bytes.Buffer
	inputdata := "teststring"
	fmt.Fprintf(&in, "%s", inputdata)
	if err := run(&in, &out); err != nil {
		t.Error(err)
	}
	if out.String() != inputdata {
		t.Errorf("Want: %q Got: %q", inputdata, out.String())
	}
}

func TestIOErrors(t *testing.T) {
	stdout := bytes.Buffer{}
	errReader := iotest.ErrReader(errors.New("read error"))
	err := run(errReader, &stdout)
	if !errors.Is(err, errCopy) {
		t.Errorf("expected %v, got %v", errCopy, err)
	}

	err = run(errReader, &stdout, "-")
	if !errors.Is(err, errCopy) {
		t.Errorf("expected %v, got %v", errCopy, err)
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

	var stdin, stdout bytes.Buffer
	stdin.WriteString("line3\n")

	if err = run(&stdin, &stdout, f1, "-", f2); err != nil {
		t.Fatal(err)
	}

	want := "line1\nline2\nline3\nline4\nline5\n"
	got := stdout.String()

	if got != want {
		t.Errorf("want: %s, got: %s", want, got)
	}
}
