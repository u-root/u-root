// Copyright 2012-20124 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trafficctl_test

import (
	"bytes"
	"errors"
	"strconv"
	"testing"

	trafficctl "github.com/u-root/u-root/pkg/tc"
)

func TestParseFilterArgs(t *testing.T) {
	for _, tt := range []struct {
		name   string
		args   []string
		err    error
		expBuf string
	}{
		{
			name: "help",
			args: []string{
				"help",
			},
			expBuf: trafficctl.Filterhelp,
		},
		{
			name: "dev",
			args: []string{
				"dev",
				"eth0",
			},
		},
		{
			name: "protocol",
			args: []string{
				"protocol",
				"ip",
			},
		},
		{
			name: "parent",
			args: []string{
				"parent",
				"1:",
			},
		},
		{
			name: "parentRoot",
			args: []string{
				"parent",
				"root",
			},
		},
		{
			name: "parentnone",
			args: []string{
				"parent",
				"none",
			},
		},
		{
			name: "handle",
			args: []string{
				"handle",
				"1:",
			},
		},
		{
			name: "handleInvalid",
			args: []string{
				"handle",
				"66000:",
			},
			err: strconv.ErrRange,
		},
		{
			name: "handleInvalid",
			args: []string{
				"handle",
				"dawawd",
			},
			err: trafficctl.ErrInvalidArg,
		},
		{
			name: "preference",
			args: []string{
				"preference",
				"2040",
			},
		},
		{
			name: "preferenceInvalid",
			args: []string{
				"preference",
				"-1",
			},
			err: trafficctl.ErrOutOfBounds,
		},
		{
			name: "preferenceInvalid",
			args: []string{
				"preference",
				"2147483647",
			},
			err: trafficctl.ErrOutOfBounds,
		},
		{
			name: "root",
			args: []string{
				"root",
			},
		},
		{
			name: "ingress",
			args: []string{
				"ingress",
			},
		},
		{
			name: "egress",
			args: []string{
				"egress",
			},
		},
		{
			name: "block",
			args: []string{
				"block",
			},
			err: trafficctl.ErrNotImplemented,
		},
		{
			name: "chain",
			args: []string{
				"chain",
			},
			err: trafficctl.ErrNotImplemented,
		},
		{
			name: "estimator",
			args: []string{
				"estimator",
			},
			err: trafficctl.ErrNotImplemented,
		},
		{
			name: "basic filter",
			args: []string{
				"basic",
				"drop",
			},
		},
		{
			name: "bpf filter",
			args: []string{
				"bpf",
			},
			err: trafficctl.ErrInvalidArg,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			var outbuf bytes.Buffer
			if _, err := trafficctl.ParseFilterArgs(&outbuf, tt.args); !errors.Is(err, tt.err) {
				t.Errorf("ParseFilterArgs() = %v, not %v", err, tt.err)
			}

			if tt.expBuf != "" {
				if tt.expBuf != outbuf.String() {
					t.Errorf("output !=  expected output")
				}
			}
		})
	}
}
