// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"os/exec"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

type test struct {
	args []string
	out  string
	err  string
}

var dirnameTests = []test{
	// For no args it seems we have to print an error.
	// It should be missing operand[s] but that's not the standard.
	{args: []string{}, err: "dirname: missing operand\n"},
	{args: []string{""}, out: ".\n"},
	{args: []string{"/this/that"}, out: "/this\n"},
	{args: []string{"/this/that", "/other"}, out: "/this\n/\n"},
	{args: []string{"/this/that", "/other thing/space"}, out: "/this\n/other thing\n"},
}

func TestDirName(t *testing.T) {
	tmpDir, dirname := testutil.CompileInTempDir(t)
	defer os.RemoveAll(tmpDir)

	// Table-driven testing
	for _, tt := range dirnameTests {
		c := exec.Command(dirname, tt.args...)
		stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
		c.Stdout, c.Stderr = stdout, stderr
		err := c.Run()
		if err != nil && tt.err == "" {
			t.Errorf("Test %v: got %q, want nil", tt.args, err)
			continue
		}

		t.Logf("RUN: %v: got %q, %q", tt.args, stdout, stderr)
		if stdout.String() != tt.out {
			t.Errorf("%v: stdout got %q, wants %q", tt.args, stdout.String(), tt.out)
		}
		if err != nil && stderr.String() != "" && stderr.String()[len("yyyy/mm/dd hh:mm:ss "):] != tt.err {
			t.Errorf("%v: stderr got %q, wants %q", tt.args, stderr.String(), tt.err)
		}

	}
}
