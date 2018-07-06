// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

type test struct {
	flags      []string
	out        string
	stdErr     string
	exitStatus int
}

func TestBasename(t *testing.T) {

	var tests = []test{
		{
			flags:      []string{"foo.h", ".h"},
			out:        "foo\n",
			stdErr:     "",
			exitStatus: 0,
		},
		{
			flags:      []string{"/bar/baz/biz/foo.h", "bar/baz/biz"},
			out:        "foo.h\n",
			stdErr:     "",
			exitStatus: 0,
		},
		{
			// Found this behavior kind of odd, if you provide just the suffix
			// name you are also asking to remove basename will still print. The
			// way I have implemented here it will print a blank line instead, that
			// seems to make more sense to me.
			flags:      []string{"-s", ".h", "/bar/biz/foo.h", ".h"},
			out:        "foo\n\n",
			stdErr:     "",
			exitStatus: 0,
		},
		{
			flags:      []string{"-s", ".h", "/bar/biz/foo.h", "/foo/bar/thing.h"},
			out:        "foo\nthing\n",
			stdErr:     "",
			exitStatus: 0,
		},
	}

	// Table-driven testing
	for _, tt := range tests {
		t.Run(fmt.Sprintf("Using flags %s", tt.flags), func(t *testing.T) {
			var out, stdErr bytes.Buffer
			cmd := testutil.Command(t, tt.flags...)
			cmd.Stdout = &out
			cmd.Stderr = &stdErr
			err := cmd.Run()

			if out.String() != tt.out {
				t.Errorf("stdout got:\n%s\nwant:\n%s", out.String(), tt.out)
			}

			if stdErr.String() != tt.stdErr {
				t.Errorf("stderr got:\n%s\nwant:\n%s", stdErr.String(), tt.stdErr)
			}

			if tt.exitStatus == 0 && err != nil {
				t.Errorf("expected to exit with %d, but exited with err %s", tt.exitStatus, err)
			}
		})
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
