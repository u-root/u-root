// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// TODO: test multi-file input

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type Test struct {
	flag string
	in   string
	out  string
}

var sortTests = []Test{
	// already sorted, in == out
	{"", "a\nb\nc\n", "a\nb\nc\n"},
	// sort letters
	{"", "c\na\nb\n", "a\nb\nc\n"},
	// sort lexicographic
	{"", "abc \nab\na bc\n", "a bc\nab\nabc \n"},
	// sort without terminating newline
	{"", "a\nb\nc", "a\nb\nc\n"},
	// sort with utf-8 characters
	{"", "γ\nα\nβ\n", "α\nβ\nγ\n"},
	// reverse sort
	{"-r", "c\na\nb\n", "c\nb\na\n"},
	// reverse sort without terminating newline
	{"-r", "a\nb\nc", "c\nb\na\n"},
}

func TestSort(t *testing.T) {
	// Create temp directory
	tmpDir, err := ioutil.TempDir("", "TestSort")
	if err != nil {
		t.Fatal("TempDir failed: ", err)
	}
	defer os.RemoveAll(tmpDir)

	// Compile sort program
	sortPath := filepath.Join(tmpDir, "sort")
	out, err := exec.Command("go", "build", "-o", sortPath, "sort.go").CombinedOutput()
	if err != nil {
		t.Fatalf("go build -o %v cmds/sort: %v\n%s", sortPath, err, string(out))
	}

	// Table-driven testing
	for _, test := range sortTests {
		// Two ways to test:
		pipes(t, test, tmpDir, sortPath)
		files(t, test, tmpDir, sortPath)
	}
}

func pipes(t *testing.T, test Test, tmpDir string, sortPath string) {
	cmd := fmt.Sprintf("printf %#v | %s %s", test.in, sortPath, test.flag)
	out, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	if err != nil {
		t.Errorf("%s: %v", cmd, err)
		return
	}
	if string(out) != test.out {
		t.Errorf("%s = %#v; want %#v", cmd, string(out), test.out)
	}
}

func files(t *testing.T, test Test, tmpDir string, sortPath string) {
	inFile := filepath.Join(tmpDir, "in")
	outFile := filepath.Join(tmpDir, "out")
	if err := ioutil.WriteFile(inFile, []byte(test.in), 0600); err != nil {
		t.Error(err)
		return
	}
	cmd := []string{sortPath, "-o", outFile, inFile}
	if test.flag != "" {
		cmd = []string{sortPath, test.flag, "-o", outFile, inFile}
	}
	out, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
	if err != nil {
		t.Errorf("%s: %v\n%s", strings.Join(cmd, " "), err, out)
		return
	}
	out, err = ioutil.ReadFile(outFile)
	if err != nil {
		t.Errorf("Cannot open out file: %v", err)
		return
	}
	if string(out) != test.out {
		t.Errorf("%s = %#v; want %#v", strings.Join(cmd, " "),
			string(out), test.out)
	}
}
