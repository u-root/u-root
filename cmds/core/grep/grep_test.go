// Copyright 2016-2017 the u-root Authors. All rights reserved
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
		{"hix\n", "hix\n", 0, []string{"-i", "hix"}},
		{"hix\n", "", 0, []string{"-i", "hox"}},
		{"HiX\n", "HiX\n", 0, []string{"-i", "hix"}},
	}

	tmpDir, err := ioutil.TempDir("", "TestGrep")
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
			t.Errorf("Grep %v != %v: want '%v', got '%v'", v.a, v.i, v.o, string(o))
			continue
		}
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
