// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

var logPrefixLength = len("2009/11/10 23:00:00 ")

func TestInvocation(t *testing.T) {
	for _, test := range []struct {
		opt []string
		out string
	}{
		{opt: []string{"-n"}, out: "id: cannot print only names in default format\n"},
		{opt: []string{"-G", "-g"}, out: "id: cannot print \"only\" of more than one choice\n"},
		{opt: []string{"-G", "-u"}, out: "id: cannot print \"only\" of more than one choice\n"},
		{opt: []string{"-g", "-u"}, out: "id: cannot print \"only\" of more than one choice\n"},
		{opt: []string{"-g", "-u", "-G"}, out: "id: cannot print \"only\" of more than one choice\n"},
	} {
		c := testutil.Command(t, test.opt...)
		stderr := &bytes.Buffer{}
		c.Stderr = stderr
		// TODO: expect the exit status to be 1.
		c.Run()

		e := stderr.String()
		// Ignore the date and time because we're using Log.Fatalf
		if e[logPrefixLength:] != test.out {
			t.Errorf("id for '%v' failed: got '%s', want '%s'", test.opt, e, test.out)
		}
	}
}

func TestMain(m *testing.M) {
	if testutil.CallMain() {
		main()
		os.Exit(0)
	}

	os.Exit(m.Run())
}
