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
	"path"
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

	tmpdir := path.Join(d, "hi.sub.dir")
	if err := os.Mkdir(tmpdir, 0777); err != nil {
		return "", err
	}

	for i := range tests {
		if err := ioutil.WriteFile(path.Join(d, tests[i].n), []byte("Go is cool!"), tests[i].m); err != nil {
			return "", err
		}
	}

	return d, nil
}

// withouth any flag
func Test_rm_1(t *testing.T) {
	d, err := setup()
	if err != nil {
		t.Fatal("Error on setup of the test: creating files and folders.")
	}
	defer os.RemoveAll(d)

	fmt.Println("== Deleting files and empty folders (no args) ...")
	files := []string{path.Join(d, "hi1.txt"), path.Join(d, "hi2.txt"), path.Join(d, "go.txt")}
	if err := rm(files, false, true, false); err != nil {
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

	fmt.Println("== Deleting folders recursively (using -r flag) ...")
	files := []string{d}
	if err := rm(files, true, true, false); err != nil {
		t.Error(err)
	}
}
