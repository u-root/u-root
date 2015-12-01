// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
//  created by Manoel Vilela (manoel_vilela@engineer.com)

package main

import (
	"testing"
	"time"
)

// without any flag
func TestDateNoFlags(t *testing.T) {
	t.Log("== Printing date with default location (no flags)...")
	msg, err := date()
	t.Logf("%v\n", msg)
	if err != nil {
		t.Error(err)
	}
	d, err := time.Parse(time.UnixDate, msg)
	if err != nil {
		t.Error(err)
	}
	out := d.Format(time.UnixDate)
	if msg != out {
		t.Error(msg, "!=", d)
		t.Error("Mismatch parsing!")
	}

}

// using u flag
func TestDateUniversal(t *testing.T) {
	flags.universal = true
	t.Log("== Printing date with UTC (using -u flag)...")
	msg, err := date()
	t.Logf("%v\n", msg)
	if err != nil {
		t.Error(err)
	}
	d, err := time.Parse(time.UnixDate, msg)
	if err != nil {
		t.Error(err)
	}
	out := d.Format(time.UnixDate)
	if msg != out {
		t.Error(msg, "!=", d)
		t.Error("Mismatch parsing!")
	}
}
