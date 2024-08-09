// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9

// Package termios implements basic termios operations including getting
// a termio struct, a winsize struct, and setting raw mode.
// To set raw mode and then restore, one can do:
// t, err := termios.Raw()
// do things
// t.Set()
package termios

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"syscall"
	"testing"

	"github.com/hugelgupf/vmtest/guest"
	"github.com/u-root/u-root/pkg/testutil"
)

var (
	// This JSON is from a real device.
	j = `{
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
	s = `speed:0 rows:72 cols:238 eof:0x04 eol2:0xff eol:0xff erase:0x7f intr:0x03 kill:0x15 lnext:0x16 min:0x00 quit:0x1c start:0x11 stop:0x13 susp:0x1a time:0x03 werase:0x17 brkint cread echo echoctl echoe echok echoke icanon icrnl iexten ignpar imaxbel isig iutf8 ixon onlcr opost pendin ~clocal ~cstopb ~echonl ~echoprt ~flusho ~hupcl ~ignbrk ~igncr ~inlcr ~inpck ~istrip ~iuclc ~ixany ~ixoff ~noflsh ~ocrnl ~ofdel ~ofill ~olcuc ~onlret ~onocr ~parenb ~parmrk ~parodd ~tostop ~xcase`
)

func TestNew(t *testing.T) {
	if _, err := New(); os.IsNotExist(err) || errors.Is(err, syscall.ENXIO) {
		t.Skipf("No /dev/tty here.")
	} else if err != nil {
		t.Errorf("TestNew: want nil, got %v", err)
	}
}

func TestChangeTermios(t *testing.T) {
	tty, err := New()
	if os.IsNotExist(err) || errors.Is(err, syscall.ENXIO) {
		t.Skipf("No /dev/tty here.")
	} else if err != nil {
		t.Fatalf("TestRaw new: want nil, got %v", err)
	}
	term, err := tty.Get()
	if err != nil {
		t.Fatalf("TestRaw get: want nil, got %v", err)
	}
	raw := MakeRaw(term)
	if reflect.DeepEqual(raw, term) {
		t.Fatalf("reflect.DeepEqual(%v, %v): true != false", term, raw)
	}
}

func TestRaw(t *testing.T) {
	// TestRaw no longer works in CircleCi, Restrict to only VM tests.
	guest.SkipIfNotInVM(t)

	tty, err := New()
	if os.IsNotExist(err) || errors.Is(err, syscall.ENXIO) {
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

func TestSet(t *testing.T) {
	g := &TTY{}
	if err := json.Unmarshal([]byte(j), g); err != nil {
		t.Fatalf("load from JSON: got %v, want nil", err)
	}
	sets := [][]string{
		{"speed", "0"},
		{"rows", "72"},
		{"cols", "238"},
		{"brkint"},
		{"~clocal"},
		{"cread"},
		{"~cstopb"},
		{"echo"},
		{"echoctl"},
		{"echoe"},
		{"echok"},
		{"echoke"},
		{"~echonl"},
		{"~echoprt"},
		{"eof", "0x04"},
		{"eol2", "0xff"},
		{"eol", "0xff"},
		{"erase", "0x7f"},
		{"~flusho"},
		{"~hupcl"},
		{"icanon"},
		{"icrnl"},
		{"iexten"},
		{"~ignbrk"},
		{"~igncr"},
		{"ignpar"},
		{"imaxbel"},
		{"~inlcr"},
		{"~inpck"},
		{"intr", "0x03"},
		{"isig"},
		{"~istrip"},
		{"~ixany"},
		{"~ixoff"},
		{"ixon"},
		{"kill", "0x15"},
		{"lnext", "0x16"},
		{"min", "0x00"},
		{"~noflsh"},
		{"~ocrnl"},
		{"onlcr"},
		{"~onlret"},
		{"~onocr"},
		{"opost"},
		{"~parenb"},
		{"~parmrk"},
		{"~parodd"},
		{"pendin"},
		{"quit", "0x1c"},
		{"start", "0x11"},
		{"stop", "0x13"},
		{"susp", "0x1a"},
		{"time", "0x03"},
		{"~tostop"},
		{"werase", "0x17"},
	}

	if runtime.GOOS == "linux" {
		sets = append(sets, []string{"~iuclc"}, []string{"~olcuc"}, []string{"~xcase"})
	}
	if runtime.GOOS != "freebsd" {
		sets = append(sets, []string{"iutf8"}, []string{"~ofdel"}, []string{"~ofill"})
	}
	for _, set := range sets {
		if err := g.SetOpts(set); err != nil {
			t.Errorf("Setting %q: got %v, want nil", set, err)
		}
	}
	bad := [][]string{
		{"hi", "1"},
		{"rows"},
		{"rows", "z"},
		{"erase"},
		{"erase", "z"},
		{"hi"},
		{"~hi"},
	}
	for _, set := range bad {
		if err := g.SetOpts(set); err == nil {
			t.Errorf("Setting %q: got nil, want err", set)
		}
	}
}

// This test tries to prevent people from breaking other operating systems.
//
// Compare:
//
//	GOOS=linux   GOARCH=amd64 go doc golang.org/x/sys/unix.Termios
//	GOOS=darwin  GOARCH=amd64 go doc golang.org/x/sys/unix.Termios
//	GOOS=openbsd GOARCH=amd64 go doc golang.org/x/sys/unix.Termios
//
// It's all a mess.
//
// This at least makes sure that the package compiles on all platforms.
func TestCrossCompile(t *testing.T) {
	testutil.SkipIfInVMTest(t)

	if testing.Short() {
		t.Skip("skipping in short mode")
	}
	platforms := []string{
		"linux/386",
		"linux/amd64",
		"linux/arm64",
		"darwin/amd64",
		"darwin/arm64",
		"freebsd/amd64",
		"openbsd/amd64",
	}
	// As of Go 1.19+, running "go" selects the right one:
	//
	// https://go.dev/doc/go1.19
	//
	// "go test and go generate now place GOROOT/bin at the
	// beginning of the PATH used for the subprocess, so tests and
	// generators that execute the go command will resolve it to
	// same GOROOT"
	//
	// And this project currently (2023-01-07) uses Go 1.19 as its
	// CI minimum, so it should be fine to just use "go" here.
	td := t.TempDir()
	for _, platform := range platforms {
		goos, goarch, _ := strings.Cut(platform, "/")
		t.Run(platform, func(t *testing.T) {
			t.Parallel()
			outFile := filepath.Join(td, goos+"-"+goarch+".test")
			cmd := exec.Command("go", "test", "-c", "-o", outFile, ".")
			cmd.Env = append(os.Environ(), "GOOS="+goos, "GOARCH="+goarch, "CGO_ENABLED=0")
			if out, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("Failed: %v, %s", err, out)
			}
		})
	}
}
