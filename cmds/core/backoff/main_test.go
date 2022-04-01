// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

func TestRunIt(t *testing.T) {
	for _, tt := range []struct {
		name    string
		timeout string
		cmd     string
		args    []string
		wantErr error
	}{
		{
			name:    "_date",
			timeout: "3s",
			cmd:     "date",
		},
		{
			name:    "noCmd",
			timeout: "3s",
			cmd:     "",
			wantErr: ErrNoCmd,
		},
		{
			name:    "echo",
			timeout: "3s",
			cmd:     "echo",
			args:    []string{"hi"},
		},
		{
			name:    "echo_missing_unit",
			timeout: "3",
			cmd:     "echo",
			args:    []string{"hi"},
			wantErr: fmt.Errorf("time: missing unit in duration \"3\""),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if err := runit(tt.timeout, tt.cmd, tt.args...); !errors.Is(err, tt.wantErr) {
				if err != nil {
					if !strings.Contains(err.Error(), tt.wantErr.Error()) {
						t.Errorf("runit(%s, %s, %s)= %q, want %q", tt.timeout, tt.cmd, tt.args, err, tt.wantErr)
					}
				}
			}
		})
	}
}

// TestOK for now just runs a simple successful test with 0 args or more than one arg.
func TestOK(t *testing.T) {
	var tests = []struct {
		args   []string
		stdout string
		stderr string
		exitok bool
	}{
		{args: []string{}, stdout: "", exitok: false},
		{args: []string{"date"}, stdout: ".*", exitok: true},
		{args: []string{"-t", "wh", "date"}, stdout: ".*", stderr: ".*invalid.*duration.*wh", exitok: false},
		{args: []string{"echo", "hi"}, stdout: ".*hi", exitok: true},
		{args: []string{"-t", "3s", "false"}, exitok: false},
	}

	for _, v := range tests {
		c := testutil.Command(t, v.args...)
		stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
		c.Stdout, c.Stderr = stdout, stderr
		err := c.Run()
		if (err != nil) && v.exitok {
			t.Errorf("%v: got %v, want nil", v, err)
		}
		if (err == nil) && !v.exitok {
			t.Errorf("%v: got nil, want err", v)
		}
		m, err := regexp.MatchString(v.stderr, stderr.String())
		if err != nil {
			t.Errorf("stderr: %v: got %v, want nil", v, err)
		} else {
			if !m {
				t.Errorf("%v: regexp.MatchString(%s, %s) false, wanted match", v, v.stderr, stderr)
			}
		}

		m, err = regexp.MatchString(v.stdout, stdout.String())
		if err != nil {
			t.Errorf("stdout: %v: got %v, want nil", v, err)
		}
		if !m {
			t.Errorf("%v: regexp.MatchString(%s, %s) false, wanted match", v, v.stdout, stderr.String())
		}
	}
}

// If you really like fork-bombing your machine, remove these lines :-)
func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
