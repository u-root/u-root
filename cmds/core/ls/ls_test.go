// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/ls"
	"github.com/u-root/u-root/pkg/testutil"
)

var tests = []struct {
	flags []string
	out   string
}{
	{
		flags: []string{},
		out: `d1
f1
f2
f3?line 2
`,
	}, {
		flags: []string{"-Q"},
		out: `"d1"
"f1"
"f2"
"f3\nline 2"
`,
	}, {
		flags: []string{"-aR"},
		out: `.
.f4
d1
d1/f4
f1
f2
f3?line 2
`,
	}, {
		flags: []string{"-R"},
		out: `d1
d1/f4
f1
f2
f3?line 2
`,
	}, {
		flags: []string{"-a"},
		out: `.
.f4
d1
f1
f2
f3?line 2
`,
	}, {
		flags: []string{"f1", "d1", "f2"},
		out: `f1
d1:
f4
f2
`,
	}, {
		flags: []string{"-d"},
		out: `.
`,
	}, {
		flags: []string{"-d", "f1", "d1", "f2"},
		out: `f1
d1
f2
`,
	},
}

func TestLs(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "ls")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create an empty directory.
	testDir := filepath.Join(tmpDir, "testDir")
	os.Mkdir(testDir, 0700)

	// Create some files.
	os.Create(filepath.Join(testDir, "f1"))
	os.Create(filepath.Join(testDir, "f2"))
	os.Create(filepath.Join(testDir, "f3\nline 2"))
	os.Create(filepath.Join(testDir, ".f4"))
	os.Mkdir(filepath.Join(testDir, "d1"), 0740)
	os.Create(filepath.Join(testDir, "d1/f4"))

	// Table-driven testing
	for _, tt := range tests {
		c := testutil.Command(t, tt.flags...)
		c.Dir = testDir
		out, err := c.Output()
		if err != nil {
			t.Error(err)
		}
		if string(out) != tt.out {
			t.Errorf("got:\n%s\nwant:\n%s", string(out), tt.out)
		}
	}
}

func TestIndicator(t *testing.T) {
	var tests = []struct {
		lsInfo ls.FileInfo
		symbol string
	}{
		{ls.FileInfo{Mode: os.ModeDir}, "/"},
		{ls.FileInfo{Mode: os.ModeNamedPipe}, "|"},
		{ls.FileInfo{Mode: os.ModeSymlink}, "@"},
		{ls.FileInfo{Mode: os.ModeSocket}, "="},
		{ls.FileInfo{Mode: 0b110110100}, ""},
		{ls.FileInfo{Mode: 0b111111101}, "*"},
	}

	for _, test := range tests {
		r := indicator(test.lsInfo)
		if r != test.symbol {
			t.Errorf("for mode %b expected %q, got %q", test.lsInfo.Mode, test.symbol, r)
		}
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
