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
	n := time.Now()
	d := date(n, time.Local)
	t.Logf("Date: %v\n", d)
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
	t.Log("::  Printing date with UTC (using -u flag)...")
	n := time.Now()
	d := date(n, time.UTC)
	t.Logf("Date: %v\n", d)
	dParsed, err := time.Parse(time.UnixDate, d)
	if err != nil {
		t.Error(err)
	}
	dTest := dParsed.Format(time.UnixDate)
	if d != dTest {
		t.Errorf("Mismatched dates; want %v, got %v\n", d, dTest)
	}
}

func TestFormatParser(t *testing.T) {
	test := "%d %w %x %D %% asdf qwer qwe s sd fqwer % qwer"
	expected := []string{"%d", "%w", "%x", "%D"}
	t.Log(":: Test of FormatParser greping the n flags")
	for index, match := range formatParser(test) {
		t.Logf(":: Parsed on iteration %d: %v", index, match)
		if match != expected[index] {
			t.Errorf("Parsing Error; Want %v, got %v\n", expected[index], match)
		}
	}
}

func TestDateMap(t *testing.T) {
	t.Log(":: Test of DateMap formatting")
	posixFormat := "%a %b %e %H:%M:%S %Z %Y"
	n := time.Now()
	test := dateMap(n, time.Local, posixFormat)
	expected := n.Format(time.UnixDate)

	if test != expected {
		t.Errorf("Mismatch outputs; \nwant %v, \n got %v", expected, test)
	}
}

func TestDateMapExamples(t *testing.T) {
	type dateTest struct {
		format  string // format flags
		example string // correct example
	}

	var tests = []dateTest{
		{
			"%a %b %e %H:%M:%S %Z %Y",
			"Tue Jun 26 09:58:10 PDT 1990",
		},
		{
			"DATE: %m/%d/%y%nTIME: %H:%M:%S",
			"DATE: 11/02/91\nTIME: 13:36:16",
		},
		{
			"TIME: %r",
			"TIME: 01:36:32 PM",
		},
	}

	t.Log(":: Sequence of examples for dateMap")
	n := time.Now()
	for _, test := range tests {
		t.Logf(" Format: \n%v\n", test.format)
		t.Logf("Example: \n%v\n", test.example)
		t.Logf(" Output: \n%v\n", dateMap(n, time.Local, test.format))
	}
}
