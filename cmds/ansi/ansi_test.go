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

// table of tests for each ansi command
// add here more command when needs
// {command, expected_escape}
var tsts = [][]string{
	{"clear", "\033[1;1H\033[2J"},
}

// Test for each ansi command
func TestAnsiCommands(t *testing.T) {
	for _, tst := range tsts {
		cmd, wants := tst[0], []byte(tst[1])
		b := &bytes.Buffer{}
		w := io.Writer(b)
		if err := ansi(w, []string{cmd}); err != nil {
			t.Error(err)
		}

		out := b.Bytes()
		if !reflect.DeepEqual(out, wants) {
			t.Fatalf("'%v' escape code mismatch; got %v, wants %v", cmd, out, wants)
		}
	}
}
