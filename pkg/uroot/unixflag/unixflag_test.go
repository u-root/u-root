// Copyright 2014-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unixflag_test

import (
	"os"
	"slices"
	"testing"

	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

func TestArgsToGoArgs(t *testing.T) {
	for _, tt := range []struct {
		name string
		args []string
		out  []string
	}{
		{name: "no args", args: []string{}, out: []string{}},
		{name: "-l", args: []string{"-l"}, out: []string{"-l"}},
		{name: "-l etc", args: []string{"-l", "etc"}, out: []string{"-l", "etc"}},
		{name: "-long etc", args: []string{"-long", "etc"}, out: []string{"-l", "-o", "-n", "-g", "etc"}},
		{name: "-long --short etc", args: []string{"-long", "--short", "etc"}, out: []string{"-l", "-o", "-n", "-g", "-short", "etc"}},
		{name: "-long --short etc -long ", args: []string{"-long", "--short", "etc", "-long"}, out: []string{"-l", "-o", "-n", "-g", "-short", "etc", "-long"}},
		{name: "-aux", args: []string{"-aux"}, out: []string{"-a", "-u", "-x"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			out := unixflag.ArgsToGoArgs(tt.args)
			if !slices.Equal(out, tt.out) {
				t.Fatalf("%v: got %v, want %v", tt.args, out, tt.out)
			}
		})
	}
}

func TestOSArgsToGoArgs(t *testing.T) {
	// because this test has to set os.Args, it is racy, so either only
	// do one case or don't run the tests concurrently.
	os.Args = []string{"sh", "-xe", "--baz", "ls", "-l", "--foo"}
	out := unixflag.OSArgsToGoArgs()
	xargs := []string{"-x", "-e", "-baz", "ls", "-l", "--foo"}
	if !slices.Equal(out, xargs) {
		t.Fatalf("%v:got %v, want %v", os.Args, out, xargs)
	}
}

func TestStringArray(t *testing.T) {
	var s unixflag.StringArray
	if err := s.Set("foo"); err != nil {
		t.Fatal(err)
	}
	if err := s.Set("bar"); err != nil {
		t.Fatal(err)
	}
	if got, want := len(s), 2; got != want {
		t.Fatalf("got slice of length %d, want %d", got, want)
	}
	if got, want := s.String(), "foo,bar"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestStringSlice(t *testing.T) {
	var s unixflag.StringSlice
	if err := s.Set("foo,bar"); err != nil {
		t.Fatal(err)
	}
	if err := s.Set("baz"); err != nil {
		t.Fatal(err)
	}
	if got, want := len(s), 3; got != want {
		t.Fatalf("got slice of length %d, want %d", got, want)
	}
	if got, want := s.String(), "foo,bar,baz"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}
