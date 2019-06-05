// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

func TestUniq(t *testing.T) {
	var (
		input1 = "test\ntest\ngo\ngo\ngo\ncoool\ncoool\ncool\nlegaal\ntest\n"
		input2 = "u-root\nuniq\nron\nron\nteam\nbinaries\ntest\n\n\n\n\n\n"
		tab    = []struct {
			i string
			o string
			s int
			a []string
		}{
			{input1, "test\ngo\ncoool\ncool\nlegaal\ntest\n", 0, nil},
			{input1, "2\ttest\n3\tgo\n2\tcoool\n1\tcool\n1\tlegaal\n1\ttest\n", 0, []string{"-c"}},
			{input1, "cool\nlegaal\ntest\n", 0, []string{"-u"}},
			{input1, "test\ngo\ncoool\n", 0, []string{"-d"}},
			{input2, "u-root\nuniq\nron\nteam\nbinaries\ntest\n\n", 0, nil},
			{input2, "1\tu-root\n1\tuniq\n2\tron\n1\tteam\n1\tbinaries\n1\ttest\n5\t\n", 0, []string{"-c"}},
			{input2, "u-root\nuniq\nteam\nbinaries\ntest\n", 0, []string{"-u"}},
			{input2, "ron\n\n", 0, []string{"-d"}},
		}
	)

	tmpDir, err := ioutil.TempDir("", "UniqTest")
	if err != nil {
		t.Fatal("TempDir failed: ", err)
	}
	defer os.RemoveAll(tmpDir)

	for _, v := range tab {
		c := testutil.Command(t, v.a...)
		c.Stdin = bytes.NewReader([]byte(v.i))
		o, err := c.CombinedOutput()
		if err := testutil.IsExitCode(err, v.s); err != nil {
			t.Error(err)
			continue
		}
		if string(o) != v.o {
			t.Errorf("Uniq %v < %v: want '%v', got '%v'", v.a, v.i, v.o, string(o))
			continue
		}
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
