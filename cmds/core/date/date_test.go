// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//  created by Manoel Vilela (manoel_vilela@engineer.com)

package main

import (
	"bytes"
	"flag"
	"os"
	"regexp"
	"strings"
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

	tests := []dateTest{
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

func TestRun(t *testing.T) {
	testfile, err := os.CreateTemp("", "testfile-for-modtime")
	if err != nil {
		t.Errorf("Unable to create testile-for-modtime: %q", err)
	}
	testfile.Close()
	defer os.RemoveAll(testfile.Name())
	fileStats, err := os.Stat(testfile.Name())
	if err != nil {
		t.Errorf("Unable to get testile-for-modtim stats: %q", err)
	}
	modTime := fileStats.ModTime().In(time.UTC).Format(time.UnixDate)
	for _, tt := range []struct {
		name    string
		arg     []string
		univ    bool
		fileref string
		expExp  string
		wantErr string
	}{
		{
			name:    "Time Now UTC",
			arg:     make([]string, 0),
			univ:    true,
			fileref: "",
		},
		{
			name:    "Time Now Local",
			arg:     make([]string, 0),
			univ:    false,
			fileref: "",
		},
		{
			name:    "Now Format+%C",
			arg:     []string{"+%C"},
			univ:    true,
			fileref: "",
			expExp:  "\\d{1,2}",
		},
		{
			name:    "Now Format+%D",
			arg:     []string{"+%D"},
			univ:    true,
			fileref: "",
			expExp:  "\\d{1,2}\\D\\d{1,2}\\D\\d{1,2}",
		},
		{
			name:    "Now Format+%j",
			arg:     []string{"+%j"},
			univ:    true,
			fileref: "",
			expExp:  "\\d*",
		},
		{
			name:    "Now Format+%r",
			arg:     []string{"+%r"},
			univ:    true,
			fileref: "",
			expExp:  "\\d{1,2}\\D\\d{1,2}\\D\\d{1,2}\\s[A,M|P,M]",
		},
		{
			name:    "Now Format+'%'T",
			arg:     []string{"+%T"},
			univ:    true,
			fileref: "",
			expExp:  "\\d\\d\\D\\d\\d\\D\\d\\d",
		},
		{
			name:    "Now Format+%W",
			arg:     []string{"+%W"},
			univ:    true,
			fileref: "",
			expExp:  "\\d{1,2}",
		},
		{
			name:    "Now Format+'%'w",
			arg:     []string{"+%w"},
			univ:    true,
			fileref: "",
			expExp:  "\\d",
		},
		{
			name:    "Now Format+%V",
			arg:     []string{"+%V"},
			univ:    true,
			fileref: "",
			expExp:  "\\d{1,2}",
		},
		{
			name:    "Now Format+'%'x",
			arg:     []string{"+%x"},
			univ:    true,
			fileref: "",
			expExp:  "\\d{0,2}\\D\\d{0,2}\\D\\d{1,2}",
		},
		{
			name:    "Now Format+'%'F",
			arg:     []string{"+%F"},
			univ:    true,
			fileref: "",
			expExp:  "\\d\\d\\D\\d*\\d\\D\\d\\d",
		},
		{
			name:    "Now Format+'%'X",
			arg:     []string{"+%X"},
			univ:    true,
			fileref: "",
			expExp:  "\\d{0,2}\\D\\d{0,2}\\D\\d{1,2}\\s[A,M|P,M]",
		},
		{
			name:    "Now Format+'%'X'%'t",
			arg:     []string{"+%X%t"},
			univ:    true,
			fileref: "",
			expExp:  "\\d{0,2}\\D\\d{0,2}\\D\\d{1,2}\\s[A,M|P,M]",
		},
		{
			name:    "File modification time",
			univ:    true,
			fileref: testfile.Name(),
			expExp:  modTime,
		},
		{
			name:    "File modification time fail",
			univ:    true,
			fileref: "not-existing-test-file",
			wantErr: "unable to gather stats of file",
		},
		{
			name: "flag usage",
			arg:  []string{"This", "dont", "work"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			rc := RealClock{}
			// Avoid spamming our CI with errors -- we don't
			// look at error output (yet), but when we do, this
			// bytes.Buffer will make it more convenient.
			var stderr bytes.Buffer
			flag.CommandLine.SetOutput(&stderr)
			if err := run(tt.arg, tt.univ, tt.fileref, rc, &buf); err != nil {
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("%q failed: %q", tt.name, err)
				}
				return
			}
			outString := buf.String()
			match, err := regexp.MatchString(tt.expExp, outString)
			if err != nil {
				t.Error(err)
			}
			if !match {
				t.Errorf("%q faile. Format of Got: %q, Want: %q", tt.name, outString, tt.expExp)
			}
		})
	}
}

type fakeClock struct {
	time.Time
}

func (f fakeClock) Now() time.Time {
	return f.Time
}

func TestGetTime(t *testing.T) {
	for _, tt := range []struct {
		name      string
		time      string
		wantYear  int
		wantMonth time.Month
		wantDay   int
		wantHour  int
		wantMin   int
		wantSec   int
		wantNsec  int
		location  *time.Location
		wantErr   string
	}{
		{
			name:      "WithoutOpt",
			time:      "11220405",
			wantYear:  time.Now().Year(),
			wantMonth: time.Month(11),
			wantDay:   22,
			wantHour:  4,
			wantMin:   0o5,
			wantSec:   0,
			location:  time.Local,
		},
		{
			name:      "WithOpt-2",
			time:      "1122040520",
			wantYear:  2020,
			wantMonth: time.Month(11),
			wantDay:   22,
			wantHour:  4,
			wantMin:   0o5,
			wantSec:   0,
			location:  time.Local,
		},
		{
			name:      "WithOpt-3",
			time:      "11220405202",
			wantYear:  time.Now().Year(),
			wantMonth: time.Month(11),
			wantDay:   22,
			wantHour:  4,
			wantMin:   0o5,
			wantSec:   0o2,
			location:  time.Local,
		},
		{
			name:      "WithOpt-4",
			time:      "112204052022",
			wantYear:  2022,
			wantMonth: time.Month(11),
			wantDay:   22,
			wantHour:  4,
			wantMin:   5,
			wantSec:   0,
			location:  time.Local,
		},
		{
			name:      "WithOpt-5",
			time:      "1122040520221",
			wantYear:  2020,
			wantMonth: time.Month(11),
			wantDay:   22,
			wantHour:  4,
			wantMin:   5,
			wantSec:   21,
			location:  time.UTC,
		},
		{
			name:      "WithOpt-all",
			time:      "112204052022.55",
			wantYear:  2022,
			wantMonth: time.Month(11),
			wantDay:   22,
			wantHour:  4,
			wantMin:   5,
			wantSec:   55,
			location:  time.Local,
		},
		{
			name:     "WithOpt-all",
			time:     "11223344201135",
			location: time.Local,
			wantErr:  "instead of [[CC]YY][.ss]",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			fc := fakeClock{}
			fc.Time = time.Date(tt.wantYear, time.Month(tt.wantMonth), tt.wantDay, tt.wantHour, tt.wantMin, tt.wantSec, tt.wantNsec, tt.location)
			testTime, err := getTime(tt.location, tt.time, fc)
			if err != nil {
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("%q failed. Got: %q, Want: %q", tt.name, err, tt.wantErr)
				}
			}
			if err == nil && !strings.Contains(fc.Time.String(), testTime.String()) {
				t.Errorf("test %q failed. Got: %q, Want: %q", tt.name, testTime, fc.Time.String())
			}
		})
	}
}
