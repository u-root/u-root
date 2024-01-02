// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hugelgupf/vmtest/guest"
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
	{"echo|wc\n", ".*", "", 0},
	{"true\n", "% % ", "", 0},
	{"false\n", "% % ", "wait: exit status 1\n", 0},
}

func TestRush(t *testing.T) {
	guest.SkipIfNotInVM(t)

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
			strout := stdout.String()
			// If you need the ^$ anchor put it in the test array.
			if !regexp.MustCompile(tt.stdout).MatchString(strout) {
				t.Errorf("Want: %#v; Got: %#v", tt.stdout, strout)
			}

			// Check stderr
			strerr := stderr.String()
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
