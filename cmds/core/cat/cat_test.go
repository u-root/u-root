// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// by Rafael Campos Nunes <rafaelnunes@engineer.com>

package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// setup writes a set of files, putting 1 byte in each file.
func setup(t *testing.T, data []byte) (string, error) {
	t.Logf(":: Creating simulation data. ")
	dir := t.TempDir()

	for i, d := range data {
		n := fmt.Sprintf("%v%d", filepath.Join(dir, "file"), i)
		if err := os.WriteFile(n, []byte{d}, 0o666); err != nil {
			return "", err
		}
	}

	return dir, nil
}

// TestCat test cat function against 4 files, in each file it is written a bit of someData
// array and the test expect the cat to return the exact same bit from someData array with
// the corresponding file.
func TestCat(t *testing.T) {
	var files []string
	someData := []byte{'l', 2, 3, 4, 'd'}

	dir, err := setup(t, someData)
	if err != nil {
		t.Fatalf("setup has failed, %v", err)
	}

	for i := range someData {
		files = append(files, fmt.Sprintf("%v%d", filepath.Join(dir, "file"), i))
	}
	var out bytes.Buffer
	if err := run(files, nil, &out); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(out.Bytes(), someData) {
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

	dir, err := setup(t, someData)
	if err != nil {
		t.Fatalf("setup has failed, %v", err)
	}

	for i := range someData {
		files = append(files, fmt.Sprintf("%v%d", filepath.Join(dir, "file"), i))
	}

	var out bytes.Buffer
	if err := run(files, nil, &out); err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(out.Bytes(), someData) {
		t.Fatalf("Reading files failed: got %v, want %v", out.Bytes(), someData)
	}
}

func TestRunFilesError(t *testing.T) {
	var files []string
	someData := []byte{'l', 2, 3, 4, 'd'}

	dir, err := setup(t, someData)
	if err != nil {
		t.Fatalf("setup has failed, %v", err)
	}

	for i := range someData {
		files = append(files, fmt.Sprintf("%v%d", filepath.Join(dir, "file"), i))
	}
	filenotexist := "testdata/doesnotexist.txt"
	files = append(files, filenotexist)
	var in, out bytes.Buffer
	if err := run(files, &in, &out); err == nil {
		t.Error("function run succeeded but should have failed")
	}
}

func TestRunNoArgs(t *testing.T) {
	var in, out bytes.Buffer
	inputdata := "teststring"
	fmt.Fprintf(&in, "%s", inputdata)
	args := make([]string, 0)
	if err := run(args, &in, &out); err != nil {
		t.Error(err)
	}
	if out.String() != inputdata {
		t.Errorf("Want: %q Got: %q", inputdata, out.String())
	}
}
