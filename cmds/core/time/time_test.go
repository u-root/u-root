// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os/exec"
	"regexp"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

type test struct {
	args   []string
	r      string
	exitok bool
}

func run(c *exec.Cmd) (string, string, error) {
	var o, e bytes.Buffer
	c.Stdout, c.Stderr = &o, &e
	err := c.Run()
	return o.String(), e.String(), err
}

func TestTime(t *testing.T) {
	tests := []test{
		{args: []string{}, r: "real 0.000.*\nuser 0.000.*\nsys 0.000", exitok: true},
		{args: []string{"date"}, r: "real [0-9][0-9]*.*\nuser [0-9][0-9]*.*\nsys [0-9][0-9]*.*", exitok: true},
		{args: []string{"deadbeef"}, r: ".*exec.*deadbeef.*executable file not found .*", exitok: false},
	}

	for _, v := range tests {
		c := testutil.Command(t, v.args...)
		_, e, err := run(c)
		if (err != nil) && v.exitok {
			t.Errorf("%v: got %v, want nil", v, err)
		}
		if (err == nil) && !v.exitok {
			t.Errorf("%v: got nil, want err", v)
		}
		m, err := regexp.MatchString(v.r, e)
		if err != nil {
			t.Errorf("%v: got %v, want nil", v, err)
			continue
		}
		if !m {
			t.Errorf("%v: regexp.MatchString(%s, %s) false, wanted match", v, v.r, e)
		}
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
