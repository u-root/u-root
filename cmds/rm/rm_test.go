// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
//  created by Manoel Vilela (manoel_vilela@engineer.com)

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type makeit struct {
	n string      // name
	m os.FileMode // mode
	s string      // for symlinks or content
}

var tests = []makeit{
	{
		n: "hi1.txt",
		m: 0666,
		s: "",
	},
	{
		n: "hi2.txt",
		m: 0777,
		s: "",
	},
	{
		n: "go.txt",
		m: 0555,
		s: "",
	},
}

// initial files and folders (for testing)
func setup() (string, error) {
	fmt.Println(":: Creating simulating data...")
	d, err := ioutil.TempDir(os.TempDir(), "hi.dir")
	if err != nil {
		return "", err
	}

	tmpdir := filepath.Join(d, "hi.sub.dir")
	if err := os.Mkdir(tmpdir, 0777); err != nil {
		return "", err
	}

	for i := range tests {
		if err := ioutil.WriteFile(filepath.Join(d, tests[i].n), []byte("Go is cool!"), tests[i].m); err != nil {
			return "", err
		}
	}

	return d, nil
}

// without any flag
func Test_rm_1(t *testing.T) {
	d, err := setup()
	if err != nil {
		t.Fatal("Error on setup of the test: creating files and folders.")
	}
	defer os.RemoveAll(d)

	fmt.Println("== Deleting files and empty folders (no args) ...")
	files := []string{filepath.Join(d, "hi1.txt"), filepath.Join(d, "hi2.txt"), filepath.Join(d, "go.txt")}

	flags.v = true
	if err := rm(files); err != nil {
		t.Error(err)
	}
}

// using r flag
func Test_rm_2(t *testing.T) {
	d, err := setup()
	if err != nil {
		t.Fatal("Error on setup of the test: creating files and folders.")
	}
	defer os.RemoveAll(d)

	flags.v = true
	flags.r = true
	fmt.Println("== Deleting folders recursively (using -r flag) ...")
	files := []string{d}
	if err := rm(files); err != nil {
		t.Error(err)
	}
}
