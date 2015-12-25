// Copyright 2016 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// By Manoel Vilela <manoel_vilela@engineer.com>

package main

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

// Simple test for ansi command: clear call
func TestAnsiClear(t *testing.T) {
	wants := []byte("\033[1;1H\033[2J")
	b := &bytes.Buffer{}
	w := io.Writer(b)

	var tst = []string{"clear"}
	if err := ansi(w, tst); err != nil {
		t.Error(err)
	}

	out := b.Bytes()
	if !reflect.DeepEqual(out, wants) {
		t.Fatalf("'Clear' escape code mismatch; got %v, wants %v", out, wants)
	}
}
