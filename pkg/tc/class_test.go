// Copyright 2012-20124 the u-root Authors. All rights reserved
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
		exp     *trafficctl.CArgs
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
				"4",
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
				"70000:4",
			},
			err: trafficctl.ErrOutOfBounds,
		},
		{
			name: "DevClassIDMinExceed",
			cmdline: []string{
				"dev",
				"eth0",
				"classid",
				"4:70000",
			},
			err: trafficctl.ErrOutOfBounds,
		},
		{
			name: "DevQDiscID",
			cmdline: []string{
				"dev",
				"eth0",
				"codel",
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
			_, err := trafficctl.ParseClassArgs(tt.cmdline, &outbuf)
			if !errors.Is(err, tt.err) {
				t.Errorf("ParseClassArgs() = %v", err)
			}
		})
	}
}
