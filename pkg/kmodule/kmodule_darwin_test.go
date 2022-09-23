// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin || linux

package kmodule

import (
	"testing"
)

func TestUnsupported(t *testing.T) {
	l, err := New()
	if err != nil {
		t.Fatalf("New(): got %v, want nil", err)
	}
	// For now, just call things that are not supported,
	// ignore return.
	l.Init(nil, "")
	l.FileInit(nil, "", 0)
	l.Delete("", 0)
	l.Probe("", "")

}
