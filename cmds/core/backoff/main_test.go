// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os/exec"
	"regexp"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

type test struct {
	args   []string
	stdout string
	stderr string
	exitok bool
}

func run(c *exec.Cmd) (string, string, error) {
	var o, e bytes.Buffer
	c.Stdout, c.Stderr = &o, &e
	err := c.Run()
	return o.String(), e.String(), err
}

// TestOK for now just runs a simple successful test with 0 args or more than one arg.
func TestOK(t *testing.T) {
	var tests = []test{
		{args: []string{}, stdout: "", exitok: true},
		{args: []string{"date"}, stdout: ".*", exitok: true},
		{args: []string{"-t", "wh", "date"}, stdout: ".*", stderr: ".*invalid.*duration.*wh", exitok: false},
		{args: []string{"echo", "hi"}, stdout: ".*hi", exitok: true},
	}

	for _, v := range tests {
		c := testutil.Command(t, v.args...)
		stdout, stderr, err := run(c)
		if (err != nil) && v.exitok {
			t.Errorf("%v: got %v, want nil", v, err)
		}
		if (err == nil) && !v.exitok {
			t.Errorf("%v: got nil, want err", v)
		}
		m, err := regexp.MatchString(v.stderr, stderr)
		if err != nil {
			t.Errorf("stderr: %v: got %v, want nil", v, err)
		} else {
			if !m {
				t.Errorf("%v: regexp.MatchString(%s, %s) false, wanted match", v, v.stderr, stderr)
			}
		}

		m, err = regexp.MatchString(v.stdout, stdout)
		if err != nil {
			t.Errorf("stdout: %v: got %v, want nil", v, err)
			continue
		}
		if !m {
			t.Errorf("%v: regexp.MatchString(%s, %s) false, wanted match", v, v.stdout, stderr)
		}
	}
}
func TestTO(t *testing.T) {
	// The integration test dies after 25s, so do shit for 6s
	var tests = []test{
		{args: []string{"-t", "6s", "false"}, stdout: ".*", stderr: ".*exit.*status.*1", exitok: false},
	}

	for _, v := range tests {
		c := testutil.Command(t, v.args...)
		stdout, stderr, err := run(c)
		if (err != nil) && v.exitok {
			t.Errorf("%v: got %v, want nil", v, err)
		}
		if (err == nil) && !v.exitok {
			t.Errorf("%v: got nil, want err", v)
		}
		m, err := regexp.MatchString(v.stderr, stderr)
		if err != nil {
			t.Errorf("stderr: %v: got %v, want nil", v, err)
		} else {
			if !m {
				t.Errorf("%v: regexp.MatchString(%s, %s) false, wanted match", v, v.stderr, stderr)
			}
		}

		m, err = regexp.MatchString(v.stdout, stdout)
		if err != nil {
			t.Errorf("stdout: %v: got %v, want nil", v, err)
			continue
		}
		if !m {
			t.Errorf("%v: regexp.MatchString(%s, %s) false, wanted match", v, v.stdout, stdout)
		}
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
