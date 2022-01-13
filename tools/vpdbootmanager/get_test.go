// Copyright 2017-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"testing"

	"github.com/u-root/u-root/pkg/vpd"
)

var getter *Getter

func TestNewgetter(t *testing.T) {
	if g := NewGetter(); g == nil {
		t.Fatalf(`NewGetter() = %v, want not nil`, g)
	} else if g.R == nil || g.Out == nil {
		t.Errorf(`g.R, g.Out = %v, %v, want not nil, not nil`, g.R, g.Out)
	}
}

// testGetOne is failing
func testGetOne(t *testing.T) {
	var buf bytes.Buffer
	getter := NewGetter()
	getter.Out = &buf
	getter.R = vpd.NewReader()
	getter.R.VpdDir = "./tests"
	if err := getter.Print("firmware_version"); err != nil {
		t.Error(err)
	}
	if buf.String() != "firmware_version(RO) => 1.2.3\n\nfirmware_version(RW) => 3.2.1\n\n" {
		t.Error("buffer contents is not correct")
	}
}

// testGetAll is failing
func testGetAll(t *testing.T) {
	var buf bytes.Buffer
	getter := NewGetter()
	getter.Out = &buf
	getter.R = vpd.NewReader()
	getter.R.VpdDir = "./tests"
	if err := getter.Print(""); err != nil {
		t.Error(err)
	}
	out := buf.String()
	if out != "firmware_version(RO) => 1.2.3\n\nsomething(RO) => else\n\nfirmware_version(RW) => 3.2.1\n\n" {
		t.Error("buffer contents is not correct")
	}
}
