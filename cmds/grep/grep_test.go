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

// GrepTest is a table-driven which spawns grep with a variety of options and inputs.
// We need to look at any output data, as well as exit status for things like the -q switch.
func TestGrep(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "grep")
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
		// BEWARE: the IO package seems to want this to be newline terminated.
		// If you just use hix with no newline the test will fail. Yuck.
		// TODO: Then don't use the IO package.
		//
		// TODO: yeesh. more, and better test cases.
		{
			in:       "hix\n",
			out:      "hix\n",
			exitCode: 0,
			args:     []string{"."},
		},
		{
			in:       "hix\n",
			out:      "",
			exitCode: 0,
			args:     []string{"-q", "."},
		},
		{
			in:       "hix\n",
			out:      "",
			exitCode: 1,
			args:     []string{"-q", "hox"},
		},
	} {
		t.Run(fmt.Sprintf("test%d", i), func(t *testing.T) {
			c := testutil.Command(t, tt.args...)
			c.Stdin = bytes.NewReader([]byte(tt.in))

			o, err := c.CombinedOutput()
			if err := testutil.IsExitCode(err, tt.exitCode); err != nil {
				t.Fatal(err)
			}

			if string(o) != tt.out {
				t.Errorf("grep %v < %v = %v, want %v", tt.args, tt.in, string(o), tt.out)
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
