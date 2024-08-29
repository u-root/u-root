// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trafficctl_test

import (
	"bytes"
	"errors"
	"testing"

	trafficctl "github.com/u-root/u-root/pkg/tc"
)

func TestParseClassArgs(t *testing.T) {
	for _, tt := range []struct {
		name    string
		cmdline []string
		exp     *trafficctl.Args
		err     error
	}{
		{
			name: "Just_Dev",
			cmdline: []string{
				"dev",
				"eth0",
			},
		},
		{
			name: "DevParent",
			cmdline: []string{
				"dev",
				"eth0",
				"parent",
				"1:1",
			},
		},
		{
			name: "DevRoot",
			cmdline: []string{
				"dev",
				"eth0",
				"root",
			},
		},
		{
			name: "DevClassID",
			cmdline: []string{
				"dev",
				"eth0",
				"classid",
				"none",
			},
		},
		{
			name: "DevClassID",
			cmdline: []string{
				"dev",
				"eth0",
				"classid",
				"root",
			},
		},
		{
			name: "DevClassID",
			cmdline: []string{
				"dev",
				"eth0",
				"classid",
				"4:4",
			},
		},
		{
			name: "DevClassIDMajExceed",
			cmdline: []string{
				"dev",
				"eth0",
				"classid",
				"FFFF:4",
			},
		},
		{
			name: "DevClassIDMinExceed",
			cmdline: []string{
				"dev",
				"eth0",
				"classid",
				"4:FFFF",
			},
		},
		{
			name: "DevQDiscInvalid",
			cmdline: []string{
				"dev",
				"eth0",
				"garbage",
			},
			err: trafficctl.ErrInvalidArg,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var outbuf bytes.Buffer
			_, err := trafficctl.ParseClassArgs(&outbuf, tt.cmdline)
			if !errors.Is(err, tt.err) {
				t.Errorf("ParseClassArgs() = %v", err)
			}
		})
	}
}
