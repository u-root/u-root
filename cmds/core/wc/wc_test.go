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

	for _, v := range tab {
		c := testutil.Command(t, v.a...)
		c.Stdin = bytes.NewReader([]byte(v.i))
		o, err := c.CombinedOutput()
		if err := testutil.IsExitCode(err, v.s); err != nil {
			t.Error(err)
			continue
		}
		if string(o) != v.o {
			t.Errorf("Wc %v < %v: want '%v', got '%v'", v.a, v.i, v.o, string(o))
			continue
		}
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
