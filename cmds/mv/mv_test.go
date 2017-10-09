// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
}

func setup() (string, error) {
	d, err := ioutil.TempDir(os.TempDir(), "hi.dir")
	if err != nil {
		return "", err
	}

	tmpdir := filepath.Join(d, "hi.sub.dir")
	if err := os.Mkdir(tmpdir, 0777); err != nil {
		return "", err
	}

	for i := range tests {
		if err := ioutil.WriteFile(filepath.Join(d, tests[i].n), []byte("hi"), tests[i].m); err != nil {
			return "", err
		}
	}

	return d, nil
}

func Test_mv_1(t *testing.T) {
	d, err := setup()
	if err != nil {
		t.Fatal("err")
	}
	defer os.RemoveAll(d)

	fmt.Println("Renaming file...")
	files1 := []string{filepath.Join(d, "hi1.txt"), filepath.Join(d, "hi4.txt")}
	if err := mv(files1, false); err != nil {
		t.Error(err)
	}

	dsub := filepath.Join(d, "hi.sub.dir")

	fmt.Println("Moving files to directory...")
	files2 := []string{filepath.Join(d, "hi2.txt"), filepath.Join(d, "hi4.txt"), dsub}
	if err := mv(files2, true); err != nil {
		t.Error(err)
	}
}
