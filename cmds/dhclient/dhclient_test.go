// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"os/exec"
	"testing"

	"github.com/u-root/u-root/shared/testutil"
)

const testFlash = "fake_test.flash"

var tests = []struct {
	cmd string
	out string
}{
	{
		cmd: "nosuchanimal",
		out: "Can't get mac for nosuchanimal\n",
	},
}

func TestDhclient(t *testing.T) {
	tmpDir, execPath := testutil.CompileInTempDir(t)
	defer os.RemoveAll(tmpDir)

	for _, tt := range tests {
		out, err := exec.Command(execPath, tt.cmd).CombinedOutput()
		if err != nil {
			t.Error(err)
		}
		out = bytes.Replace(out, []byte{0}, []byte{}, -1)
		if string(out) != tt.out {
			t.Errorf("expected:\n%s\ngot:\n%s", tt.out, string(out))
		}
	}
}
