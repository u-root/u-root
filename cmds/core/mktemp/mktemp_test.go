// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type test struct {
	flags mktempflags
	args  []string
	out   string
	err   error
}

func TestMkTemp(t *testing.T) {
	tmpDir := os.TempDir()
	tests := []test{
		{
			flags: mktempflags{},
			out:   tmpDir,
			err:   nil,
		},
		{
			flags: mktempflags{d: true},
			out:   tmpDir,
			err:   nil,
		},
		{
			flags: mktempflags{},
			args:  []string{"foofoo.XXXX"},
			out:   filepath.Join(tmpDir, "foofoo"),
			err:   nil,
		},
		{
			flags: mktempflags{suffix: "baz"},
			args:  []string{"foo.XXXX"},
			out:   filepath.Join(tmpDir, "foo.baz"),
			err:   nil,
		},
		{
			flags: mktempflags{u: true, q: true},
			out:   "",
			err:   nil,
		},
	}

	// Table-driven testing
	for _, tt := range tests {
		var stdout bytes.Buffer
		cmd := command(&stdout, tt.flags, tt.args)
		err := cmd.run()
		if err != tt.err {
			t.Errorf("expected %v, got %v", tt.err, err)
		}

		r := stdout.String()
		if !strings.HasPrefix(r, tt.out) {
			t.Errorf("stdout got:\n%s\nwant starting with:\n%s", r, tt.out)
		}
	}
}
