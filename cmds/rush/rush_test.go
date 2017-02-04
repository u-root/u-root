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
	"regexp"
	"strings"
	"syscall"
	"testing"
)

var tests = []struct {
	stdin  string // input
	stdout string // output (regular expression)
	stderr string // output (regular expression)
	ret    int    // output
}{
	// TODO: Create a `-c` flag for rush so stdout does not contain
	// prompts, or have the prompt be derived from $PS1.
	{"exit\n", "% ", "", 0},
	{"exit 77\n", "% ", "", 77},
	{"exit 1 2 3\n", "% % ", "Too many arguments\n", 0},
	{"exit abcd\n", "% % ", "Non numeric argument\n", 0},
	{"time cd .\n", "% % ", `real 0.0\d\d\n`, 0},
	{"time sleep 0.25\n", "% % ", `real 0.2\d\d\nuser 0.00\d\nsys 0.00\d\n`, 0},
}

func TestRush(t *testing.T) {
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
		cmd.Stdin = strings.NewReader(tt.stdin)
		var stdout bytes.Buffer
		cmd.Stdout = &stdout
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err := cmd.Run()

		// Check stdout
		strout := string(stdout.Bytes())
		if !regexp.MustCompile("^" + tt.stdout + "$").MatchString(strout) {
			t.Errorf("Want: %#v; Got: %#v", tt.stdout, strout)
		}

		// Check stderr
		strerr := string(stderr.Bytes())
		if !regexp.MustCompile("^" + tt.stderr + "$").MatchString(strerr) {
			t.Errorf("Want: %#v; Got: %#v", tt.stderr, strerr)
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
