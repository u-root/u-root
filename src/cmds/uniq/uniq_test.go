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
	"syscall"
	"testing"
)

func TestUniq(t *testing.T) {
	var (
		input1 string = "test\ntest\ngo\ngo\ngo\ncoool\ncoool\ncool\nlegaal\ntest\n"
		input2 string = "u-root\nuniq\nron\nron\nteam\nbinaries\ntest\n\n\n\n\n\n"
		tab           = []struct {
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

	uniqtestpath := filepath.Join(tmpDir, "uniqtest.exe")
	out, err := exec.Command("go", "build", "-o", uniqtestpath, ".").CombinedOutput()
	if err != nil {
		t.Fatalf("go build -o %v cmds/uniq: %v\n%s", uniqtestpath, err, string(out))
	}

	t.Logf("Built %v for test", uniqtestpath)
	for _, v := range tab {
		c := exec.Command(uniqtestpath, v.a...)
		c.Stdin = bytes.NewReader([]byte(v.i))
		o, err := c.CombinedOutput()
		s := c.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()

		if s != v.s {
			t.Errorf("Uniq %v < %v > %v: want (exit: %v), got (exit %v)", v.a, v.i, v.o, v.s, s)
			continue
		}

		if err != nil && s != v.s {
			t.Errorf("Uniq %v < %v > %v: want nil, got %v", v.a, v.i, v.o, err)
			continue
		}
		if string(o) != v.o {
			t.Errorf("Uniq %v < %v: want '%v', got '%v'", v.a, v.i, v.o, string(o))
			continue
		}
		t.Logf("Uniq %v < %v: %v", v.a, v.i, v.o)
	}
}
