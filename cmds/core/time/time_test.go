// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"regexp"
	"testing"
)

func TestTimeRewrite(t *testing.T) {
	var tests = []struct {
		args      []string
		want      string
		wantError bool
	}{
		{
			want:      "real 0.000.*\nuser 0.000.*\nsys 0.000",
			wantError: false,
		},
		{
			args:      []string{"date"},
			want:      "real [0-9][0-9]*.*\nuser [0-9][0-9]*.*\nsys [0-9][0-9]*.*",
			wantError: false,
		},
		{
			args:      []string{"deadbeef"},
			wantError: true,
		},
	}

	for _, test := range tests {
		var stdin, stdout, stderr bytes.Buffer
		err := run(test.args, &stdin, &stdout, &stderr)
		if test.wantError && err == nil {
			t.Error("want error but got nil")
			continue
		}
		if !test.wantError && err != nil {
			t.Errorf("want nil got: %v", err)
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
