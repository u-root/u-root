// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

// GrepTest is a table-driven which spawns grep with a variety of options and inputs.
// We need to look at any output data, as well as exit status for things like the -q switch.
func TestGrep(t *testing.T) {
	tab := []struct {
		input  string
		output string
		status int
		args   []string
	}{
		// BEWARE: the IO package seems to want this to be newline terminated.
		// If you just use hix with no newline the test will fail. Yuck.
		{"hix\n", "hix\n", 0, []string{"."}},
		{"hix\n", "", 0, []string{"-q", "."}},
		{"hix\n", "hix\n", 0, []string{"-i", "hix"}},
		{"hix\n", "", 0, []string{"-i", "hox"}},
		{"HiX\n", "HiX\n", 0, []string{"-i", "hix"}},
		{"hix\n", ":0:hix\n", 0, []string{"-n", "hix"}},
		{"hix\n", "hix\n", 0, []string{"-e", "hix"}},
		{"hix\n", "1\n", 0, []string{"-c", "hix"}},
		// These tests don't make a lot of sense the way we're running it, but
		// hopefully it'll make codecov shut up.
		{"hix\n", "hix\n", 0, []string{"-h", "hix"}},
		{"hix\n", "hix\n", 0, []string{"-r", "hix"}},
		{"hix\nfoo\n", "foo\n", 0, []string{"-v", "hix"}},
		{"hix\n", "\n", 0, []string{"-l", "hix"}}, // no filename, so it just prints a newline
	}

	for _, v := range tab {
		c := testutil.Command(t, v.args...)
		c.Stdin = bytes.NewReader([]byte(v.input))
		o, err := c.CombinedOutput()
		if err := testutil.IsExitCode(err, v.status); err != nil {
			t.Error(err)
			continue
		}
		if string(o) != v.output {
			t.Errorf("grep %v != %v: want '%v', got '%v'", v.args, v.input, v.output, string(o))
			continue
		}
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
