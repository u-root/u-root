// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

const testFlash = "fake_test.flash"

var tests = []struct {
	cmd string
	out string
}{
	// Test summary
	{
		cmd: "summary",
		out: `Fmap found at 0x5f74:
	Signature:  __FMAP__
	VerMajor:   1
	VerMinor:   0
	Base:       0xcafebabedeadbeef
	Size:       0x44332211
	Name:       Fake flash
	NAreas:     2
	Areas[0]:
		Offset:  0xdeadbeef
		Size:    0x11111111
		Name:    Area Number 1Hello
		Flags:   0x1013 (STATIC|COMPRESSED|0x1010)
	Areas[1]:
		Offset:  0xcafebabe
		Size:    0x22222222
		Name:    Area Number 2xxxxxxxxxxxxxxxxxxx
		Flags:   0x0 (0x0)
`,
	},
	// Test usage
	{
		cmd: "usage",
		out: `Legend: '.' - full (0xff), '0' - zero (0x00), '#' - mixed
0x00000000: 0..###
Blocks:       6 (100.0%)
Full (0xff):  2 (33.3%)
Empty (0x00): 1 (16.7%)
Mixed:        3 (50.0%)
`,
	},
}

// Table driven testing
func TestFmap(t *testing.T) {
	tmpDir, execPath := testutil.CompileInTempDir(t)
	defer os.RemoveAll(tmpDir)

	for _, tt := range tests {
		out, err := exec.Command(execPath, tt.cmd, testFlash).CombinedOutput()
		if err != nil {
			t.Error(err)
		}
		// Filter out null characters which may be present in fmap strings.
		out = bytes.Replace(out, []byte{0}, []byte{}, -1)
		if string(out) != tt.out {
			t.Errorf("expected:\n%s\ngot:\n%s", tt.out, string(out))
		}
	}
}

func TestJson(t *testing.T) {
	tmpDir, execPath := testutil.CompileInTempDir(t)
	defer os.RemoveAll(tmpDir)

	jsonFile := filepath.Join(tmpDir, "tmp.json")
	if err := exec.Command(execPath, "jget", jsonFile, testFlash).Run(); err != nil {
		t.Fatal(err)
	}
	got, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		t.Fatal(err)
	}
	want := `{
	"FMap": {
		"Signature": [
			95,
			95,
			70,
			77,
			65,
			80,
			95,
			95
		],
		"VerMajor": 1,
		"VerMinor": 0,
		"Base": 14627333968688430831,
		"Size": 1144201745,
		"Name": "Fake flash",
		"NAreas": 2,
		"Areas": [
			{
				"Offset": 3735928559,
				"Size": 286331153,
				"Name": "Area Number 1\u0000\u0000\u0000Hello",
				"Flags": 4115
			},
			{
				"Offset": 3405691582,
				"Size": 572662306,
				"Name": "Area Number 2xxxxxxxxxxxxxxxxxxx",
				"Flags": 0
			}
		]
	},
	"Metadata": {
		"Start": 24436
	}
}
`
	if string(got) != want {
		t.Errorf("want:%s; got:%s", string(want), got)
	}
}
