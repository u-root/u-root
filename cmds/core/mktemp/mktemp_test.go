// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type test struct {
	cmdline []string
	out     string
	err     error
}

func TestMkTemp(t *testing.T) {
	tmpDir := os.TempDir()
	tests := []test{
		{
			cmdline: []string{"mktemp"},
			out:     tmpDir,
			err:     nil,
		},
		{
			cmdline: []string{"mktemp", "-d"},
			out:     tmpDir,
			err:     nil,
		},
		{
			cmdline: []string{"mktemp", "foofoo.XXXX"},
			out:     filepath.Join(tmpDir, "foofoo"),
			err:     nil,
		},
		{
			cmdline: []string{"mktemp", "--suffix", "baz", "foo.XXXX"},
			out:     filepath.Join(tmpDir, "foo.baz"),
			err:     nil,
		},
		{
			cmdline: []string{"mktemp", "-u"},
			out:     "",
			err:     nil,
		},
	}

	// Table-driven testing
	for _, tt := range tests {
		var stdout bytes.Buffer
		cmd := command(&stdout, tt.cmdline)
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

func TestDefaultFlags(t *testing.T) {
	f := flags{}
	f.register(flag.CommandLine)

	if f.d {
		t.Error("directory should be false by default")
	}
	if f.u {
		t.Error("dry-run should be false by default")
	}
	if f.q {
		t.Error("quiet should be false by default")
	}
	if f.prefix != "" {
		t.Error("prefix should be empty string by default")
	}
	if f.suffix != "" {
		t.Error("suffix should be empty string by default")
	}
	if f.dir != "" {
		t.Error("tmpdir should be empty string by default")
	}
}
