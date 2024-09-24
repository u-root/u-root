// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/testutil"
)

var tests = []struct {
	iface  string
	isIPv4 string
	out    string
}{
	{
		iface:  "nosuchanimal",
		isIPv4: "-ipv4=true",
		out:    "no interfaces match nosuchanimal\n",
	},
}

func TestDhclient(t *testing.T) {
	for _, tt := range tests {
		out, err := testutil.Command(t, tt.isIPv4, tt.iface).CombinedOutput()
		if err == nil {
			t.Errorf("%v: got nil, want err", tt)
		}
		if !strings.HasSuffix(string(out), tt.out) {
			t.Errorf("expected:\n%s\ngot:\n%s", tt.out, string(out))
		}
	}
}

func TestMain(m *testing.M) {
	testutil.Run(m, main)
}
