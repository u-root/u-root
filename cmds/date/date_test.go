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
	t.Log("::  Printing date with default location (no flags)...")
	d, err := date()
	t.Logf("Date: %v\n", d)
	if err != nil {
		t.Error(err)
	}
	dParsed, err := time.Parse(time.UnixDate, d)
	if err != nil {
		t.Error(err)
	}
	dTest := dParsed.Format(time.UnixDate)
	if d != dTest {
		t.Errorf("Mismatched dates; want %v, got %v\n", d, dTest)
	}

}

// using u flag
func TestDateUniversal(t *testing.T) {
	flags.universal = true
	t.Log("::  Printing date with UTC (using -u flag)...")
	d, err := date()
	t.Logf("Date: %v\n", d)
	if err != nil {
		t.Error(err)
	}
	dParsed, err := time.Parse(time.UnixDate, d)
	if err != nil {
		t.Error(err)
	}
	dTest := dParsed.Format(time.UnixDate)
	if d != dTest {
		t.Errorf("Mismatched dates; want %v, got %v\n", d, dTest)
	}
}
