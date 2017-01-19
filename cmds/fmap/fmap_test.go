// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os/exec"
	"testing"
)

func TestFlashSummary(t *testing.T) {
	testFlash := "fake_test.flash"
	expected := `Fmap found at 0x9f4:
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
`
	out, err := exec.Command("go", "run", "fmap.go", "-s", testFlash).CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	out = bytes.Replace(out, []byte{0}, []byte{}, -1)
	if string(out) != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, string(out))
	}
}
