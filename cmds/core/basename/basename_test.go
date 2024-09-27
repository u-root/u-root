// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
)

type test struct {
	err   error
	out   string
	flags []string
}

func TestRun(t *testing.T) {
	tests := []test{
		{
			flags: []string{"foo.h", ".h"},
			out:   "foo\n",
		},
		{
			flags: []string{"/bar/baz/biz/foo.h", "bar/baz/biz"},
			out:   "foo.h\n",
		},
		{
			flags: []string{".h", ".h"},
			out:   ".h\n",
		},
		{
			flags: []string{"/some/path/foo"},
			out:   "foo\n",
		},
		{
			flags: []string{"/some/path/foo"},
			out:   "foo\n",
		},
		{
			flags: []string{},
			err:   errUsage,
		},
	}

	// Table-driven testing
	for _, tt := range tests {
		t.Run(fmt.Sprintf("Using flags %s", tt.flags), func(t *testing.T) {
			var out bytes.Buffer

			err := run(&out, tt.flags)
			if !errors.Is(err, tt.err) {
				t.Fatalf("got %v want %v", err, tt.err)
			}

			if out.String() != tt.out {
				t.Errorf("stdout got:\n%s\nwant:\n%s", out.String(), tt.out)
			}
		})
	}
}
