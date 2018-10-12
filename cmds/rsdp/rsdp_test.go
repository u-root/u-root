// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

func TestRSDP(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "rsdp")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	n := filepath.Join(tmpDir, "f")
	var tests = []struct {
		name   string
		args   []string
		data   string
		status int
		out    string
		skip   int
	}{
		{name: "bad file", args: []string{"-f", ""}, data: "", status: 1, out: "open : no such file or directory\n", skip: 20},
		{name: "too many args (fla + 1)", args: []string{"-f", "", ""}, data: "", status: 1, out: "Usage: rsdp [-f file]\n", skip: 20},
		{name: "too many args (1)", args: []string{""}, data: "", status: 1, out: "Usage: rsdp [-f file]\n", skip: 20},
		{name: "rsdp", args: []string{"-f", n}, data: "a b c\n6,209,0,-;ACPI: RSDP 0x00000000000F6A10 000024 (v02 PTLTD )\nc d\n", out: " acpi_rsdp=0x00000000000F6A10 \n"},
		{name: "no rsdp", args: []string{"-f", n}, data: "a b c\n6,209,0,-;ACPI: SDP 0x00000000000F6A10 000024 (v02 PTLTD )\nc d\n", status: 1, out: "Could not find RSDP\n", skip: 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.data != "" {
				if err := ioutil.WriteFile(n, []byte(tt.data), 0666); err != nil {
					t.Error(err)
					return
				}
			}
			c := testutil.Command(t, tt.args...)
			out, _ := c.CombinedOutput()
			status := c.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
			if tt.status != status {
				t.Errorf("err got: %v want %v", status, tt.status)
			}
			m := string(out[tt.skip:])
			if m != tt.out {
				t.Errorf("got: '%q', want '%q'", m, tt.out)
			}
		})
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
