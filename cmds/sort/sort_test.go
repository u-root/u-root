// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

type sortTest struct {
	name  string
	flags []string
	in    string
	out   string
}

var tests = []sortTest{
	{
		name:  "empty",
		flags: []string{},
		in:    "",
		out:   "",
	},
	{
		name:  "already sorted, in == out",
		flags: []string{},
		in:    "a\nb\nc\n",
		out:   "a\nb\nc\n",
	},
	{
		name:  "sort letters",
		flags: []string{},
		in:    "c\na\nb\n",
		out:   "a\nb\nc\n",
	},
	{
		name:  "sort lexicographic",
		flags: []string{},
		in:    "abc \nab\na bc\n",
		out:   "a bc\nab\nabc \n",
	},
	{
		name:  "sort without terminating newline",
		flags: []string{},
		in:    "a\nb\nc",
		out:   "a\nb\nc\n",
	},
	{
		name:  "sort with utf-8 characters",
		flags: []string{},
		in:    "γ\nα\nβ\n",
		out:   "α\nβ\nγ\n",
	},
	{
		name:  "reverse sort",
		flags: []string{"-r"},
		in:    "c\na\nb\n",
		out:   "c\nb\na\n",
	},
	{
		name:  "reverse sort without terminating newline",
		flags: []string{"-r"},
		in:    "a\nb\nc",
		out:   "c\nb\na\n",
	},
}

// sort < in > out
func TestSortWithPipes(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := testutil.Command(t, tt.flags...)
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
		})
	}
}

// Helper function to create input files, run sort and compare the output.
func sortWithFiles(t *testing.T, tt sortTest, tmpDir string, inFiles []string, outFile string) {
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

	args := append(append(tt.flags, "-o", outPath), inPaths...)
	out, err := testutil.Command(t, args...).CombinedOutput()
	if err != nil {
		t.Errorf("sort %s: %v\n%s", strings.Join(args, " "), err, out)
		return
	}

	out, err = ioutil.ReadFile(outPath)
	if err != nil {
		t.Errorf("Cannot open out file: %v", err)
		return
	}
	if string(out) != tt.out {
		t.Errorf("sort %s = %#v; want %#v", strings.Join(args, " "), string(out), tt.out)
	}
}

// sort -o in out
func TestSortWithFiles(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "sort_files")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortWithFiles(t, tt, tmpDir, []string{"in"}, "out")
		})
	}
}

// sort -o file file
func TestInplaceSort(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "sort_files")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortWithFiles(t, tt, tmpDir, []string{"file"}, "file")
		})
	}
}

// sort -o out in1 in2 in3 in4
func TestMultipleFileInputs(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "sort_files")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tt := sortTest{
		name:  "multiple ins",
		flags: []string{},
		in:    "a\nb\nc\n",
		out:   "a\na\na\na\nb\nb\nb\nb\nc\nc\nc\nc\n",
	}
	sortWithFiles(t, tt, tmpDir,
		[]string{"in1", "in2", "in3", "in4"}, "out")

	// Run the test again without newline terminators.
	tt = sortTest{
		name:  "multiple ins with terminators",
		flags: []string{},
		in:    "a\nb\nc",
		out:   "a\na\na\na\nb\nb\nb\nb\nc\nc\nc\nc\n",
	}
	sortWithFiles(t, tt, tmpDir,
		[]string{"in1", "in2", "in3", "in4"}, "out")
}

func TestMain(m *testing.M) {
	if testutil.CallMain() {
		main()
		os.Exit(0)
	}

	os.Exit(m.Run())
}
