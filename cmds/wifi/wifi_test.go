// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

type WifiTestCase struct {
	name   string
	args   []string
	expect string
}

var testcases = []WifiTestCase{
	{
		name:   "No Flags, No Args",
		args:   nil,
		expect: "Usage",
	},
	{
		name:   "Flags, No Args",
		args:   []string{"-i=123"},
		expect: "Usage",
	},
}

func run(c *exec.Cmd) (string, string, error) {
	var o, e bytes.Buffer
	c.Stdout, c.Stderr = &o, &e
	err := c.Run()
	return o.String(), e.String(), err
}

func TestWifi(t *testing.T) {
	// Set up
	tmpDir, execPath := testutil.CompileInTempDir(t)
	defer os.RemoveAll(tmpDir)

	// Tests
	for _, test := range testcases {
		t.Logf("TEST %v", test.name)
		c := exec.Command(execPath, test.args...)
		_, e, _ := run(c)
		if !strings.Contains(e, test.expect) {
			execStatement := fmt.Sprintf("exec(wifi %s)", strings.Trim(fmt.Sprint(test.args), "[]"))
			t.Errorf("%s\ngot:\n%s\n\nwant:\n%s", execStatement, e, test.expect)
		}
	}
}
