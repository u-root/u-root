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

func TestParseActionGAT(t *testing.T) {
	for _, tt := range []struct {
		name string
		args []string
		err  error
	}{
		{
			name: "continue",
			args: []string{
				"continue",
			},
		},
		{
			name: "drop",
			args: []string{
				"drop",
			},
		},
		{
			name: "shot",
			args: []string{
				"shot",
			},
		},
		{
			name: "pass",
			args: []string{
				"pass",
			},
		},
		{
			name: "ok",
			args: []string{
				"ok",
			},
		},
		{
			name: "reclassify",
			args: []string{
				"reclassify",
			},
		},
		{
			name: "pipe",
			args: []string{
				"pipe",
			},
		},
		{
			name: "goto",
			args: []string{
				"goto",
			},
		},
		{
			name: "jump",
			args: []string{
				"jump",
			},
		},
		{
			name: "trap",
			args: []string{
				"trap",
			},
		},
		{
			name: "invalid",
			args: []string{
				"invalid",
			},
			err: trafficctl.ErrInvalidActionControl,
		},
		{
			name: "invalidNumArgs",
			args: []string{},
			err:  trafficctl.ErrNotEnoughArgs,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			var outBuf bytes.Buffer
			if _, err := trafficctl.ParseActionGAT(&outBuf, tt.args); !errors.Is(err, tt.err) {
				t.Errorf("ParseActionGAT(%q) = %v, not %v", tt.args, err, tt.err)
			}
		})
	}
}
