// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"os"
	"syscall"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

var (
	uskip = len("2018/08/10 21:20:42 ")
)

func TestSimple(t *testing.T) {
	if os.Getuid() != 0 {
		t.Skip("Must be root for this test")
	}

	tmpDir, err := ioutil.TempDir("", "pox")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	var tests = []struct {
		args   []string
		name   string
		status int
		out    string
		skip   int
		stdin  *testutil.FakeStdin
	}{
		//		  -c, --create          create it (default true)
		//  -d, --debug           enable debug prints
		//  -o, --output string   Output file (default "/tmp/pox.tcz")
		//  -t, --test            run a test with the first argument
		{
			args:   []string{"-o", "/tmp/x/a/g/c/d/e/f/g", "bin/bash"},
			name:   "Bad executable",
			status: 1,
			out:    "open bin/bash: no such file or directory\n",
			skip:   uskip,
		},
		{
			args:   []string{"-o", "/tmp/x/a/g/c/d/e/f/g", "/bin/bash"},
			name:   "Bad output file",
			status: 1,
			out:    "/tmp/x/a/g/c/d/e/f/g -noappend]: Could not stat destination file: Not a directory\n: exit status 1\n",
			skip:   uskip + len("[mksquashfs /tmp/pox373051153 "), // the tempname varies so skip it.
		},
		{
			args:  []string{"-t", "/bin/bash"},
			name:  "shellexit",
			out:   "shell-init: error retrieving current directory: getcwd: cannot access parent directories: No such file or directory\n",
			skip:  0,
			stdin: testutil.NewFakeStdin("exit"),
		},
	}

	// Table-driven testing
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := testutil.Command(t, tt.args...)
			// ignore the error, we deal with it via process status,
			// and most of these commands are supposed to get an error.
			out, _ := c.CombinedOutput()
			status := c.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
			if tt.status != status {
				t.Errorf("err got: %v want %v", status, tt.status)
			}
			if len(out) < tt.skip {
				t.Errorf("err got: %v wanted at least %d bytes", string(out), tt.skip)
				return
			}
			m := string(out[tt.skip:])
			if m != tt.out {
				t.Errorf("got:'%q'(%d bytes) want:'%q'(%d bytes)", m, len(m), tt.out, len(tt.out))
			}
		})
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
