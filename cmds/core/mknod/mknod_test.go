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
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

type test struct {
	args    []string
	expects string
}

func run(c *exec.Cmd) (string, string, error) {
	var o, e bytes.Buffer
	c.Stdout, c.Stderr = &o, &e
	err := c.Run()
	return o.String(), e.String(), err
}

func TestMknodFifo(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "ls")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Delete a preexisting pipe if it exists, thought it shouldn't
	pipepath := filepath.Join(tmpDir, "testpipe")
	_ = os.Remove(pipepath)
	if _, err := os.Stat(pipepath); err != nil && !os.IsNotExist(err) {
		// Couldn't delete the file for reasons other than it didn't exist.
		t.Fatalf("couldn't delete preexisting pipe")
	}

	// Make a pipe and check that it exists.
	c := testutil.Command(t, pipepath, "p")
	c.Run()
	if _, err := os.Stat(pipepath); os.IsNotExist(err) {
		t.Errorf("Pipe was not created.")
	}
}

func TestInvocationErrors(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "ls")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	devpath := filepath.Join(tmpDir, "testdev")
	var tests = []test{
		{args: []string{devpath}, expects: "mknod: usage: mknod path type [major minor]\n"},
		{args: []string{""}, expects: "mknod: usage: mknod path type [major minor]\n"},
		{args: []string{devpath, "p", "254", "3"}, expects: "mknod: device type p requires no other arguments\n"},
		{args: []string{devpath, "b", "254"}, expects: "mknod: usage: mknod path type [major minor]\n"},
		{args: []string{devpath, "b"}, expects: "mknod: device type b requires a major and minor number\n"},
		{args: []string{devpath, "k"}, expects: "mknod: device type not recognized: k\n"},
	}

	for _, v := range tests {
		c := testutil.Command(t, v.args...)
		_, e, _ := run(c)
		if e[20:] != v.expects {
			t.Errorf("mknod for '%v' failed: got '%s', want '%s'", v.args, e[20:], v.expects)
		}
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
