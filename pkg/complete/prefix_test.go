// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package complete

import (
	"testing"
)

func TestPrefix(t *testing.T) {
	var tests = []string{
		"/etc/hosts.allow",
		"/etc/hosts.deny",
		"/etc/hosts",
		"/etc/host.conf",
		"/etc/hostname",
		"/etc",
		"/a",
		"",
	}

	var testsout = []string{
		"/etc/hosts.allow",
		"/etc/hosts.",
		"/etc/hosts",
		"/etc/host",
		"/etc/host",
		"/etc",
		"/",
		"",
	}
	for i := range tests {
		test := tests[:i+1]
		p := Prefix(test)
		if p != testsout[i] {
			t.Errorf("Test %d: %v: got %v, want %v", i, test, p, testsout[i])
		}
	}
}
