// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"os/exec"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

var tests = []struct {
	iface  string
	isIPv4 string
	test   string
	out    string
}{
	{
		iface:  "nosuchanimal",
		isIPv4: "-ipv4=true",
		test:   "-test=true",
		out:    "cannot get interface by name nosuchanimal: Link not found\n",
	},
}

func TestDhclient(t *testing.T) {
	tmpDir, execPath := testutil.CompileInTempDir(t)
	defer os.RemoveAll(tmpDir)

	for _, tt := range tests {
		out, err := exec.Command(execPath, tt.isIPv4, tt.test, tt.iface).CombinedOutput()
		if err != nil {
			t.Error(err)
		}
		out = bytes.Replace(out, []byte{0}, []byte{}, -1)
		if string(out) != tt.out {
			t.Errorf("expected:\n%s\ngot:\n%s", tt.out, string(out))
		}
	}
}
