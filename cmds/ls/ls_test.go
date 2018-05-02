// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

func TestLS(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "ls")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create some files.
	os.Create(filepath.Join(tmpDir, "f1"))
	os.Create(filepath.Join(tmpDir, "f2"))
	os.Create(filepath.Join(tmpDir, "f3\nline 2"))
	os.Create(filepath.Join(tmpDir, ".f4"))
	os.Mkdir(filepath.Join(tmpDir, "d1"), 0740)
	os.Create(filepath.Join(tmpDir, "d1/f4"))

	// Table-driven testing
	for _, tt := range []struct {
		args []string
		out  string
		wd   string
	}{
		{
			args: []string{},
			wd:   tmpDir,
			out: `d1
f1
f2
f3?line 2
`,
		}, {
			args: []string{"-Q"},
			wd:   tmpDir,
			out: `"d1"
"f1"
"f2"
"f3\nline 2"
`,
		}, {
			args: []string{"-R"},
			wd:   tmpDir,
			out: `d1
d1/f4
f1
f2
f3?line 2
`,
		}, {
			args: []string{"-a"},
			wd:   tmpDir,
			out: `.
.f4
d1
f1
f2
f3?line 2
`,
		}, {
			args: []string{tmpDir},
			wd:   filepath.Join(tmpDir, "d1"),
			out: `d1
f1
f2
f3?line 2
`,
		},
	} {
		c := testutil.Command(t, tt.args...)
		c.Dir = tt.wd
		out, err := c.Output()
		if err != nil {
			t.Error(err)
		}
		if string(out) != tt.out {
			t.Errorf("got:\n%s\nwant:\n%s", string(out), tt.out)
		}
	}
}

func TestMain(m *testing.M) {
	if testutil.CallMain() {
		main()
		os.Exit(0)
	}

	os.Exit(m.Run())
}
