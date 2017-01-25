// Copyright 2016-2017 the u-root Authors. All rights reserved
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

//
func TestWc(t *testing.T) {
	var tab = []struct {
		i string
		o string
		s int
		a []string
	}{
		{"simple test count words", "4\n", 0, []string{"-w"}}, // don't fail more
		{"lines\nlines\n", "2\n", 0, []string{"-l"}},
		{"count chars\n", "12\n", 0, []string{"-c"}},
	}

	tmpDir, err := ioutil.TempDir("", "TestWc")
	if err != nil {
		t.Fatal("TempDir failed: ", err)
	}
	defer os.RemoveAll(tmpDir)

	testwcpath := filepath.Join(tmpDir, "testwc.exe")
	out, err := exec.Command("go", "build", "-o", testwcpath, ".").CombinedOutput()
	if err != nil {
		t.Fatalf("go build -o %v cmds/wc: %v\n%s", testwcpath, err, string(out))
	}

	t.Logf("Built %v for test", testwcpath)
	for _, v := range tab {
		c := exec.Command(testwcpath, v.a...)
		c.Stdin = bytes.NewReader([]byte(v.i))
		o, err := c.CombinedOutput()
		s := c.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()

		if s != v.s {
			t.Errorf("Wc %v < %v > %v: want (exit: %v), got (exit %v)", v.a, v.i, v.o, v.s, s)
			continue
		}

		if err != nil && s != v.s {
			t.Errorf("Wc %v < %v > %v: want nil, got %v", v.a, v.i, v.o, err)
			continue
		}
		if string(o) != v.o {
			t.Errorf("Wc %v < %v: want '%v', got '%v'", v.a, v.i, v.o, string(o))
			continue
		}
		t.Logf("[ok] Wc %v < %v: %v", v.a, v.i, v.o)
	}
}
