// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

func TestWC(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "wc")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	for i, tt := range []struct {
		in       string
		out      string
		exitCode int
		args     []string
	}{
		{"simple test count words", "4\n", 0, []string{"-w"}}, // don't fail more
		{"lines\nlines\n", "2\n", 0, []string{"-l"}},
		{"count chars\n", "12\n", 0, []string{"-c"}},
	} {
		t.Run(fmt.Sprintf("test%d", i), func(t *testing.T) {
			c := testutil.Command(t, tt.args...)
			c.Stdin = bytes.NewReader([]byte(tt.in))

			o, err := c.CombinedOutput()
			if err := testutil.IsExitCode(err, tt.exitCode); err != nil {
				t.Fatal(err)
			}
			if out := string(o); out != tt.out {
				t.Errorf("wc %v < %v = %v, want %v", tt.args, tt.in, out, tt.out)
			}
		})
	}
}

func TestMain(m *testing.M) {
	if testutil.CallMain() {
		main()
		os.Exit(0)
	}

	os.Exit(m.Run())
}
