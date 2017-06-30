// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

var (
	testPath = "."
	// if true removeAll the testPath on the end
	remove = true
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
	tempDir, mknodPath := testutil.CompileInTempDir(t)

	if remove {
		defer os.RemoveAll(tempDir)
	}

	// Delete a preexisting pipe if it exists, thought it shouldn't
	pipepath := filepath.Join(tempDir, "testpipe")
	_ = os.Remove(pipepath)
	if _, err := os.Stat(pipepath); err != nil && !os.IsNotExist(err) {
		// Couldn't delete the file for reasons other than it didn't exist.
		t.Fatalf("couldn't delete preexisting pipe")
	}

	// Make a pipe and check that it exists.
	fmt.Print(pipepath)
	c := exec.Command(mknodPath, pipepath, "p")
	c.Run()
	if _, err := os.Stat(pipepath); os.IsNotExist(err) {
		t.Errorf("Pipe was not created.")
	}
}

func TestInvocationErrors(t *testing.T) {
	tempDir, mknodPath := testutil.CompileInTempDir(t)

	if remove {
		defer os.RemoveAll(tempDir)
	}

	devpath := filepath.Join(tempDir, "testdev")
	var tests = []test{
		{args: []string{devpath}, expects: "Usage: mknod path type [major minor]\n"},
		{args: []string{""}, expects: "Usage: mknod path type [major minor]\n"},
		{args: []string{devpath, "p", "254", "3"}, expects: "device type p requires no other arguments\n"},
		{args: []string{devpath, "b", "254"}, expects: "Usage: mknod path type [major minor]\n"},
		{args: []string{devpath, "b"}, expects: "device type b requires a major and minor number\n"},
		{args: []string{devpath, "k"}, expects: "device type not recognized: k\n"},
	}

	for _, v := range tests {
		c := exec.Command(mknodPath, v.args...)
		_, e, _ := run(c)
		if e[20:] != v.expects {
			t.Errorf("mknod for '%v' failed: got '%s', want '%s'", v.args, e[20:], v.expects)
		}
	}
}

func TestMknodBlock(t *testing.T) {
	curuser, err := user.Current()
	if err != nil {
		t.Fatalf("can't get current user.")
	}

	if curuser.Uid != "0" {
		t.Logf("not root, uid %v, skipping test\n", curuser.Uid)
		return
	}
	t.Log("root user, proceeding\n")
	//TODO(ganshun): implement block test
	t.Skip("Unimplemented test, need root")
}
