// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os/exec"
	"testing"
)

var tests = []struct {
	flag string
	out  string
}{
	// Test summary
	{
		flag: "-s",
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
		flag: "-u",
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
	for _, tt := range tests {
		testFlash := "fake_test.flash"
		out, err := exec.Command("go", "run", "fmap.go", tt.flag, testFlash).CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		// Filter out null characters which may be present in fmap strings.
		out = bytes.Replace(out, []byte{0}, []byte{}, -1)
		if string(out) != tt.out {
			t.Errorf("expected:\n%s\ngot:\n%s", tt.out, string(out))
		}
	}
}
