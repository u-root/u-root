// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package termios implements basic termios operations including getting
// a termio struct, a winsize struct, and setting raw mode.
// To set raw mode and then restore, one can do:
// t, err := termios.Raw()
// do things
// t.Set()
package termios

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	if _, err := New(); err != nil {
		t.Errorf("TestNew: want nil, got %v", err)
	}

}

func TestRaw(t *testing.T) {
	tty, err := New()
	if err != nil {
		t.Fatalf("TestRaw new: want nil, got %v", err)
	}
	term, err := tty.Get()
	if err != nil {
		t.Fatalf("TestRaw get: want nil, got %v", err)
	}

	n, err := tty.Raw()
	if err != nil {
		t.Fatalf("TestRaw raw: want nil, got %v", err)
	}
	if !reflect.DeepEqual(term, n) {
		t.Fatalf("TestRaw: New(%v) and Raw(%v) should be equal, are not", t, n)
	}
	if err := tty.Set(n); err != nil {
		t.Fatalf("TestRaw restore mode: want nil, got %v", err)
	}
	n, err = tty.Get()
	if err != nil {
		t.Fatalf("TestRaw second call to New(): want nil, got %v", err)
	}
	if !reflect.DeepEqual(term, n) {
		t.Fatalf("TestRaw: After Raw restore: New(%v) and check(%v) should be equal, are not", term, n)
	}
}
