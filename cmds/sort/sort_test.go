// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

type test struct {
	flags []string
	in    string
	out   string
}

var sortTests = []test{
	// empty
	{[]string{}, "", ""},
	// already sorted, in == out
	{[]string{}, "a\nb\nc\n", "a\nb\nc\n"},
	// sort letters
	{[]string{}, "c\na\nb\n", "a\nb\nc\n"},
	// sort lexicographic
	{[]string{}, "abc \nab\na bc\n", "a bc\nab\nabc \n"},
	// sort without terminating newline
	{[]string{}, "a\nb\nc", "a\nb\nc\n"},
	// sort with utf-8 characters
	{[]string{}, "γ\nα\nβ\n", "α\nβ\nγ\n"},
	// reverse sort
	{[]string{"-r"}, "c\na\nb\n", "c\nb\na\n"},
	// reverse sort without terminating newline
	{[]string{"-r"}, "a\nb\nc", "c\nb\na\n"},
}

// sort < in > out
func TestSortWithPipes(t *testing.T) {
	tmpDir, sortPath := testutil.CompileInTempDir(t)
	defer os.RemoveAll(tmpDir)

	// Table-driven testing
	for _, tt := range sortTests {
		cmd := exec.Command(sortPath, tt.flags...)
		cmd.Stdin = strings.NewReader(tt.in)
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			t.Errorf("sort(%#v): %v", tt.in, err)
		}
		if out.String() != tt.out {
			t.Errorf("sort(%#v) = %#v; want %#v", tt.in,
				out.String(), tt.out)
		}
	}
}

// Helper function to create input files, run sort and compare the output.
func sortWithFiles(t *testing.T, tt test, tmpDir string, sortPath string,
	inFiles []string, outFile string) {
	// Create input files
	inPaths := make([]string, len(inFiles))
	for i, inFile := range inFiles {
		inPaths[i] = filepath.Join(tmpDir, inFile)
		if err := ioutil.WriteFile(inPaths[i], []byte(tt.in), 0600); err != nil {
			t.Error(err)
			return
		}
	}
	outPath := filepath.Join(tmpDir, outFile)

	args := append(append(append([]string{}, tt.flags...), "-o",
		outPath), inPaths...)
	out, err := exec.Command(sortPath, args...).CombinedOutput()
	if err != nil {
		t.Errorf("sort %s: %v\n%s", strings.Join(args, " "),
			err, out)
		return
	}

	out, err = ioutil.ReadFile(outPath)
	if err != nil {
		t.Errorf("Cannot open out file: %v", err)
		return
	}
	if string(out) != tt.out {
		t.Errorf("sort %s = %#v; want %#v", strings.Join(args, " "),
			string(out), tt.out)
	}
}

// sort -o in out
func TestSortWithFiles(t *testing.T) {
	tmpDir, sortPath := testutil.CompileInTempDir(t)
	defer os.RemoveAll(tmpDir)

	// Table-driven testing
	for _, tt := range sortTests {
		sortWithFiles(t, tt, tmpDir, sortPath, []string{"in"}, "out")
	}
}

// sort -o file file
func TestInplaceSort(t *testing.T) {
	tmpDir, sortPath := testutil.CompileInTempDir(t)
	defer os.RemoveAll(tmpDir)

	// Table-driven testing
	for _, tt := range sortTests {
		sortWithFiles(t, tt, tmpDir, sortPath, []string{"file"}, "file")
	}
}

// sort -o out in1 in2 in3 in4
func TestMultipleFileInputs(t *testing.T) {
	tmpDir, sortPath := testutil.CompileInTempDir(t)
	defer os.RemoveAll(tmpDir)

	tt := test{[]string{}, "a\nb\nc\n",
		"a\na\na\na\nb\nb\nb\nb\nc\nc\nc\nc\n"}
	sortWithFiles(t, tt, tmpDir, sortPath,
		[]string{"in1", "in2", "in3", "in4"}, "out")

	// Run the test again without newline terminators.
	tt = test{[]string{}, "a\nb\nc",
		"a\na\na\na\nb\nb\nb\nb\nc\nc\nc\nc\n"}
	sortWithFiles(t, tt, tmpDir, sortPath,
		[]string{"in1", "in2", "in3", "in4"}, "out")
}
