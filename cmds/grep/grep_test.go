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

// GrepTest is a table-driven which spawns grep with a variety of options and inputs.
// We need to look at any output data, as well as exit status for things like the -q switch.
func TestGrep(t *testing.T) {
	var tab = []struct {
		i string
		o string
		s int
		a []string
	}{
		// BEWARE: the IO package seems to want this to be newline terminated.
		// If you just use hix with no newline the test will fail. Yuck.
		{"hix\n", "hix\n", 0, []string{"."}},
		{"hix\n", "", 0, []string{"-q", "."}},
		{"hix\n", "", 1, []string{"-q", "hox"}},
	}

	tmpDir, err := ioutil.TempDir("", "TestGrep")
	if err != nil {
		t.Fatal("TempDir failed: ", err)
	}
	defer os.RemoveAll(tmpDir)

	testgreppath := filepath.Join(tmpDir, "testgrep.exe")
	out, err := exec.Command("go", "build", "-o", testgreppath, ".").CombinedOutput()
	if err != nil {
		t.Fatalf("go build -o %v cmds/grep: %v\n%s", testgreppath, err, string(out))
	}

	t.Logf("Built %v for test", testgreppath)
	for _, v := range tab {
		t.Logf("Run %v args %v", testgreppath, v)
		c := exec.Command(testgreppath, v.a...)
		c.Stdin = bytes.NewReader([]byte(v.i))
		o, err := c.CombinedOutput()
		s := c.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()

		if s != v.s {
			t.Errorf("Grep %v < %v > %v: want (exit: %v), got (exit %v)", v.a, v.i, v.o, v.s, s)
			continue
		}

		if err != nil && s != v.s {
			t.Errorf("Grep %v < %v > %v: want nil, got %v", v.a, v.i, v.o, err)
			continue
		}
		if string(o) != v.o {
			t.Errorf("Grep %v < %v: want '%v', got '%v'", v.a, v.i, v.o, string(o))
			continue
		}
		t.Logf("Grep %v < %v: %v", v.a, v.i, v.o)
	}
}
