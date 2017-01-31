// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

type test struct {
	in  string
	out string
	ret int
}

var tests = []test{
	{"exit\n", "", 0},
	{"exit 77\n", "", 77},
	{"exit 1 2 3\n", "Too many arguments\n", 0},
	{"exit abcd\n", "Non numeric argument\n", 0},
}

func TestExit(t *testing.T) {
	// Create temp directory
	tmpDir, err := ioutil.TempDir("", "TestExit")
	if err != nil {
		t.Fatal("TempDir failed: ", err)
	}
	defer os.RemoveAll(tmpDir)

	// Compile rush
	rushPath := filepath.Join(tmpDir, "rush")
	out, err := exec.Command("go", "build", "-o", rushPath).CombinedOutput()
	if err != nil {
		t.Fatalf("go build -o %v cmds/rush: %v\n%s", rushPath, err, string(out))
	}

	// Table-driven testing
	for _, tt := range tests {
		// Run command
		cmd := exec.Command(rushPath)
		cmd.Stdin = strings.NewReader(tt.in)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err := cmd.Run()

		// Check stderr
		out := string(stderr.Bytes())
		if out != tt.out {
			t.Errorf("Want: %#v; Got: %#v", tt.out, out)
		}

		// Check return code
		retCode := 0
		if err != nil {
			exitErr, ok := err.(*exec.ExitError)
			if !ok {
				t.Errorf("Error running rush: %v", err)
				continue
			}
			retCode = exitErr.Sys().(syscall.WaitStatus).ExitStatus()
		}
		if retCode != tt.ret {
			t.Errorf("Want: %d; Got: %d", tt.ret, retCode)
		}
	}
}
