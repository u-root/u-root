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

func errorExists(err error) bool {
	return err != nil
}

func craftPrintMsg(errExists bool, out string) string {
	var msg bytes.Buffer

	if errExists {
		msg.WriteString("Error Status: exists\n")
	} else {
		msg.WriteString("Error Status: not exists\n")
	}
	msg.WriteString("Output:\n")
	msg.WriteString(out)
	return msg.String()
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
		if (test.errExists != errorExists(err)) || !strings.Contains(string(out), test.out) {
			expectMsg := craftPrintMsg(test.errExists, test.out)
			actualMsg := craftPrintMsg(errorExists(err), string(out))
			execStatement := fmt.Sprintf("exec(wifi %s)", strings.Trim(fmt.Sprint(test.args), "[]"))
			t.Errorf("%s\ngot:\n%s\n\nwant:\n%s", execStatement, actualMsg, expectMsg)
		}
	}
}
