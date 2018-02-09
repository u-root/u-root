// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

type WifiTestCase struct {
	name      string
	args      []string
	out       string
	errExists bool
}

var testcases = []WifiTestCase{
	{
		name:      "No Flags, No Args",
		args:      nil,
		out:       "Usage",
		errExists: true,
	},
	{
		name:      "Flags, No Args",
		args:      []string{"-i=123"},
		out:       "Usage",
		errExists: true,
	},
}

func TestWifi(t *testing.T) {
	// Set up
	tmpDir, execPath := testutil.CompileInTempDir(t)
	defer os.RemoveAll(tmpDir)

	// Tests
	for _, test := range testcases {
		t.Logf("TEST %v", test.name)
		c := exec.Command(execPath, test.args...)
		out, err := c.CombinedOutput()
		if (test.errExists != testutil.ErrorExists(err)) || !strings.Contains(string(out), test.out) {
			execStatement := fmt.Sprintf("exec(wifi %s)", strings.Trim(fmt.Sprint(test.args), "[]"))
			testutil.PrintError(t, execStatement, test.out, test.errExists, string(out), err)
		}
	}
}
