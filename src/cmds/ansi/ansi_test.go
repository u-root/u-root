// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// By Manoel Vilela <manoel_vilela@engineer.com>

package main

import (
	"bytes"
	"io"
	"testing"
)

// Simple test for ansi command: clear call
func TestAnsiClear(t *testing.T) {
	wants := "\033[1;1H\033[2J"
	b := &bytes.Buffer{}
	w := io.Writer(b)

	var tst = []string{"clear"}
	if err := ansi(w, tst); err != nil {
		t.Error(err)
	}

	v, err := b.ReadString('J')
	if err != nil {
		t.Error(err)
	}

	if v != wants {
		t.Fatalf("Clear escape code buffering; got %v, wants %v", []byte(value), []byte(expected))
	}
}
