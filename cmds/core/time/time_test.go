// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"os/exec"
	"regexp"
	"testing"
)

func TestTime(t *testing.T) {
	tests := []struct {
		args []string
		want string
		err  error
	}{
		{
			want: "real 0.000.*\nuser 0.000.*\nsys 0.000",
		},
		{
			args: []string{"date"},
			want: "real [0-9][0-9]*.*\nuser [0-9][0-9]*.*\nsys [0-9][0-9]*.*",
		},
		{
			args: []string{"deadbeef"},
			err:  exec.ErrNotFound,
		},
	}

	for _, test := range tests {
		var stdin, stdout, stderr bytes.Buffer
		err := run(test.args, &stdin, &stdout, &stderr)
		if !errors.Is(err, test.err) {
			t.Errorf("got %v, want %v", err, test.err)
			continue
		}

		res := stderr.String()
		m, err := regexp.MatchString(test.want, res)
		if err != nil {
			t.Fatal(err)
		}
		if !m {
			t.Errorf("regexp.MatchString(%q, %q) false, wanted match", test.want, res)
		}
	}
}
