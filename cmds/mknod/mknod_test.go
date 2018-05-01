// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

var logPrefixLength = len("2009/11/10 23:00:00 ")

func TestMknodFifo(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "mknod_fifo")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	pipe := filepath.Join(tmpDir, "pipe")

	c := testutil.Command(t, pipe, "p")
	if err := c.Run(); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(pipe); os.IsNotExist(err) {
		t.Error("Pipe was not created.")
	}
}

func TestInvocationErrors(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "mknod_fifo")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	devpath := filepath.Join(tmpDir, "testdev")

	for _, tt := range []struct {
		args []string
		want string
	}{
		{args: []string{devpath}, want: "Usage: mknod path type [major minor]\n"},
		{args: []string{""}, want: "Usage: mknod path type [major minor]\n"},
		{args: []string{devpath, "p", "254", "3"}, want: "device type p requires no other arguments\n"},
		{args: []string{devpath, "b", "254"}, want: "Usage: mknod path type [major minor]\n"},
		{args: []string{devpath, "b"}, want: "device type b requires a major and minor number\n"},
		{args: []string{devpath, "k"}, want: "device type not recognized: k\n"},
	} {
		c := testutil.Command(t, tt.args...)
		var stderr bytes.Buffer
		c.Stderr = &stderr
		c.Run()

		if e := stderr.String(); e[logPrefixLength:] != tt.want {
			t.Errorf("mknod for %q failed: got %q, want %q", tt.args, e[logPrefixLength:], tt.want)
		}
	}
}

func TestMain(m *testing.M) {
	if testutil.CallMain() {
		main()
		os.Exit(0)
	}

	os.Exit(m.Run())
}
