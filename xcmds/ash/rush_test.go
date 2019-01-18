// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
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
	{"time sleep 0.25\n", "% % ", `real \d+.\d{3}\nuser \d+.\d{3}\nsys \d+.\d{3}\n`, 0},
}

func testRush(t *testing.T) {
	// Create temp directory
	tmpDir, err := ioutil.TempDir("", "TestExit")
	if err != nil {
		t.Fatal("TempDir failed: ", err)
	}
	defer os.RemoveAll(tmpDir)

	// Table-driven testing
	for i, tt := range tests {
		t.Run(fmt.Sprintf("test%d", i), func(t *testing.T) {
			// Run command
			cmd := testutil.Command(t)
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
			if err := testutil.IsExitCode(err, tt.ret); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
