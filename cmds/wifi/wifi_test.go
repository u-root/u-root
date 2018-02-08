// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
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
	exp_error bool
}

var testcases = []WifiTestCase{
	{
		name:      "No Flags, No Args",
		args:      nil,
		out:       "Usage",
		exp_error: true,
	},
	{
		name:      "Flags, No Args",
		args:      []string{"-i=123"},
		out:       "Usage",
		exp_error: true,
	},
}

var ERROR_MSG_FORMAT = "\nEXPECTED:\n%s\n\nACTUAL:\n%s\n"

func errorExists(err error) bool {
	return err != nil
}

func craftPrintMsg(err_exists bool, out string) string {
	var msg bytes.Buffer

	if err_exists {
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
		if (test.exp_error != errorExists(err)) || !strings.Contains(string(out), test.out) {
			expectMsg := craftPrintMsg(test.exp_error, test.out)
			actualMsg := craftPrintMsg(errorExists(err), string(out))
			t.Errorf(ERROR_MSG_FORMAT, expectMsg, actualMsg)
		}
	}
}
