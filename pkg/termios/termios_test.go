// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !plan9

// Package termios implements basic termios operations including getting
// a termio struct, a winsize struct, and setting raw mode.
// To set raw mode and then restore, one can do:
// t, err := termios.Raw()
// do things
// t.Set()
package termios

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

func TestNew(t *testing.T) {
	if _, err := New(); os.IsNotExist(err) {
		t.Skipf("No /dev/tty here.")
	} else if err != nil {
		t.Errorf("TestNew: want nil, got %v", err)
	}

}

func TestRaw(t *testing.T) {
	// TestRaw no longer works in CircleCi, Restrict to only VM tests.
	testutil.SkipIfNotRoot(t)
	tty, err := New()
	if os.IsNotExist(err) {
		t.Skipf("No /dev/tty here.")
	} else if err != nil {
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

// Test proper unmarshaling and consistent, repeatable output from String()
func TestString(t *testing.T) {
	// This JSON is from a real device.
	j := `{
	"Ispeed": 0,
	"Ospeed": 0,
	"Row": 72,
	"Col": 238,
	"CC": {
		"eof": 4,
		"eol": 255,
		"eol2": 255,
		"erase": 127,
		"intr": 3,
		"kill": 21,
		"lnext": 22,
		"min": 0,
		"quit": 28,
		"start": 17,
		"stop": 19,
		"susp": 26,
		"werase": 23,
		"time": 3
	},
	"Opts": {
		"xcase": false,
		"brkint": true,
		"clocal": false,
		"cread": true,
		"cstopb": false,
		"echo": true,
		"echoctl": true,
		"echoe": true,
		"echok": true,
		"echoke": true,
		"echonl": false,
		"echoprt": false,
		"flusho": false,
		"hupcl": false,
		"icanon": true,
		"icrnl": true,
		"iexten": true,
		"ignbrk": false,
		"igncr": false,
		"ignpar": true,
		"imaxbel": true,
		"inlcr": false,
		"inpck": false,
		"isig": true,
		"istrip": false,
		"iuclc": false,
		"iutf8": true,
		"ixany": false,
		"ixoff": false,
		"ixon": true,
		"noflsh": false,
		"ocrnl": false,
		"ofdel": false,
		"ofill": false,
		"olcuc": false,
		"onlcr": true,
		"onlret": false,
		"onocr": false,
		"opost": true,
		"parenb": false,
		"parmrk": false,
		"parodd": false,
		"pendin": true,
		"tostop": false
	}
}`
	s := `speed:0 rows:72 cols:238 brkint:1 clocal:0 cread:1 cstopb:0 echo:1 echoctl:1 echoe:1 echok:1 echoke:1 echonl:0 echoprt:0 eof:0x04 eol2:0xff eol:0xff erase:0x7f flusho:0 hupcl:0 icanon:1 icrnl:1 iexten:1 ignbrk:0 igncr:0 ignpar:1 imaxbel:1 inlcr:0 inpck:0 intr:0x03 isig:1 istrip:0 iuclc:0 iutf8:1 ixany:0 ixoff:0 ixon:1 kill:0x15 lnext:0x16 min:0x00 noflsh:0 ocrnl:0 ofdel:0 ofill:0 olcuc:0 onlcr:1 onlret:0 onocr:0 opost:1 parenb:0 parmrk:0 parodd:0 pendin:1 quit:0x1c start:0x11 stop:0x13 susp:0x1a time:0x03 tostop:0 werase:0x17 xcase:0`
	g := &TTY{}
	if err := json.Unmarshal([]byte(j), g); err != nil {
		t.Fatalf("stty load: %v", err)
	}

	if g.String() != s {
		t.Errorf("GTTY: want '%v', got '%v'", s, g.String())
		as := strings.Split(s, " ")
		ag := strings.Split(g.String(), " ")
		if len(as) != len(ag) {
			t.Fatalf("Wrong # elements in gtty: want %d, got %d", len(as), len(ag))
		}
		for i := range as {
			t.Errorf("want %s got %s Same %v", as[i], ag[i], as[i] == ag[i])
		}

	}

}
