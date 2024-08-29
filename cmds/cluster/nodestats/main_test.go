// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"bytes"
	"testing"
)

func TestNodeStat(t *testing.T) {
	var b bytes.Buffer
	if err := run(&b, []string{"nodestat"}); err != nil {
		t.Fatal(err)
	}
	if b.Len() == 0 {
		t.Fatal("output from run: got 0 bytes, expected at least 1")
	}
	if err := run(&b, []string{"nodestat", "bad arg"}); err == nil {
		t.Fatalf("run with one arg: got nil, expect err")
	}
}
