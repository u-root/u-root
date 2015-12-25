// Copyright 2015 the u-root Authors. All rights reserved
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
	expected := "\033[1;1H\033[2J"
	b := &bytes.Buffer{}
	w := io.Writer(b)
	var test = []string{"clear"}
	if err := ansi(w, test); err != nil {
		t.Error(err)
	}
	value, err := b.ReadString('J')
	if err != nil {
		t.Error(err)
	}

	if value != expected {
		t.Fatalf("Clear escape code buffering; got %v, wants %v", []byte(value), []byte(expected[:]))
	}
}
