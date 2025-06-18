// Copyright 2012-2024 the u-root Authors. All rights reserved
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

func TestParseQdiscArgs(t *testing.T) {
	for _, tt := range []struct {
		name   string
		args   []string
		expBuf string
		err    error
	}{
		{
			name: "help",
			args: []string{
				"help",
			},
			expBuf: trafficctl.QdiscHelp,
		},
		{
			name: "dev",
			args: []string{
				"dev",
				"eth0",
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
			name: "clsact",
			args: []string{
				"clsact",
			},
		},
		{
			name: "parent",
			args: []string{
				"parent",
				"1:10",
			},
		},
		{
			name: "qdisc codel",
			args: []string{
				"codel",
			},
		},
		{
			name: "qdisc qfq",
			args: []string{
				"qfq",
			},
		},
		{
			name: "qdisc codel",
			args: []string{
				"codel",
			},
		},
		{
			name: "qdisc htb",
			args: []string{
				"htb",
			},
		},
		{
			name: "qdisc hfsc",
			args: []string{
				"hfsc",
			},
		},
		{
			name: "qdisc htb",
			args: []string{
				"htb",
			},
		},
		{
			name: "estimator",
			args: []string{
				"estimator",
			},
			err: trafficctl.ErrNotImplemented,
		},
		{
			name: "stab",
			args: []string{
				"stab",
			},
			err: trafficctl.ErrNotImplemented,
		},
		{
			name: "ingress_block",
			args: []string{
				"ingress_block",
			},
			err: trafficctl.ErrNotImplemented,
		},
		{
			name: "egress_block",
			args: []string{
				"egress_block",
			},
			err: trafficctl.ErrNotImplemented,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			var buf bytes.Buffer

			_, err := trafficctl.ParseQdiscArgs(&buf, tt.args)
			if !errors.Is(err, tt.err) {
				t.Errorf("ParseQDiscArgs = %v, not %v", err, tt.err)
			}

			if tt.expBuf != "" {
				if tt.expBuf != buf.String() {
					t.Errorf("output !=  expected output")
				}
			}
		})
	}
}

func TestParseCodelArgs(t *testing.T) {
	for _, tt := range []struct {
		name   string
		args   []string
		err    error
		expBuf string
	}{
		{
			name: "noArgs",
			args: []string{},
		},
		{
			name: "limit",
			args: []string{
				"limit",
				"10",
			},
		},
		{
			name: "limitInvalid",
			args: []string{
				"limit",
				"-1",
			},
			err: strconv.ErrSyntax,
		},
		{
			name: "target",
			args: []string{
				"target",
				"10",
			},
		},
		{
			name: "targetInvalid",
			args: []string{
				"target",
				"-1",
			},
			err: strconv.ErrSyntax,
		},
		{
			name: "targetInvalid",
			args: []string{
				"target",
				"2147483647",
			},
		},
		{
			name: "interval s",
			args: []string{
				"interval",
				"10s",
			},
		},
		{
			name: "interval sec",
			args: []string{
				"interval",
				"10sec",
			},
		},
		{
			name: "interval secs",
			args: []string{
				"interval",
				"10secs",
			},
		},
		{
			name: "interval ms",
			args: []string{
				"interval",
				"10ms",
			},
		},
		{
			name: "interval msec",
			args: []string{
				"interval",
				"10msec",
			},
		},
		{
			name: "interval msecs",
			args: []string{
				"interval",
				"10msecs",
			},
		},
		{
			name: "interval us",
			args: []string{
				"interval",
				"10us",
			},
		},
		{
			name: "interval usec",
			args: []string{
				"interval",
				"10usec",
			},
		},
		{
			name: "interval usecs",
			args: []string{
				"interval",
				"10usecs",
			},
		},
		{
			name: "interval fail",
			args: []string{
				"interval",
				"1asd0usasdasdecs",
			},
			err: strconv.ErrSyntax,
		},
		{
			name: "ce_threshold s",
			args: []string{
				"ce_threshold",
				"10s",
			},
		},
		{
			name: "ce_threshold sec",
			args: []string{
				"ce_threshold",
				"10sec",
			},
		},
		{
			name: "ce_threshold secs",
			args: []string{
				"ce_threshold",
				"10secs",
			},
		},
		{
			name: "ce_threshold ms",
			args: []string{
				"ce_threshold",
				"10ms",
			},
		},
		{
			name: "ce_threshold msec",
			args: []string{
				"ce_threshold",
				"10msec",
			},
		},
		{
			name: "ce_threshold msecs",
			args: []string{
				"ce_threshold",
				"10msecs",
			},
		},
		{
			name: "ce_threshold us",
			args: []string{
				"ce_threshold",
				"10us",
			},
		},
		{
			name: "ce_threshold usec",
			args: []string{
				"ce_threshold",
				"10usec",
			},
		},
		{
			name: "ce_threshold usecs",
			args: []string{
				"interce_thresholdval",
				"10usecs",
			},
		},
		{
			name: "ce_threshold fail",
			args: []string{
				"ce_threshold",
				"1asd0usasdasdecs",
			},
			err: strconv.ErrSyntax,
		},
		{
			name: "ecn",
			args: []string{
				"ecn",
			},
		},
		{
			name: "noecn",
			args: []string{
				"noecn",
			},
		},
		{
			name: "help",
			args: []string{
				"help",
			},
			expBuf: trafficctl.CodelHelp,
			err:    trafficctl.ErrExitAfterHelp,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			var outBuf bytes.Buffer
			_, err := trafficctl.ParseCodelArgs(&outBuf, tt.args)
			if !errors.Is(err, tt.err) {
				t.Errorf("ParseCodelArgs() = %v, not %v", err, tt.err)
			}

			if tt.expBuf != "" {
				if tt.expBuf != outBuf.String() {
					t.Errorf("output !=  expected output")
				}
			}
		})
	}
}

func TestParseQFQArgs(t *testing.T) {
	for _, tt := range []struct {
		name   string
		args   []string
		expBuf string
		err    error
	}{
		{
			name: "No Args",
			args: []string{},
		},
		{
			name: "Weight",
			args: []string{
				"weight",
				"100",
			},
		},
		{
			name: "maxpkt",
			args: []string{
				"maxpkt",
				"100",
			},
		},
		{
			name: "Help",
			args: []string{
				"help",
			},
			err: trafficctl.ErrExitAfterHelp,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var outBuf bytes.Buffer
			_, err := trafficctl.ParseQFQArgs(&outBuf, tt.args)
			if !errors.Is(err, tt.err) {
				t.Errorf("ParseQFQArgs() = %v, not %v", err, tt.err)
			}

			if tt.expBuf != "" {
				if tt.expBuf != outBuf.String() {
					t.Errorf("output !=  expected output")
				}
			}
		})
	}
}

func TestParseHTBArgs(t *testing.T) {
	for _, tt := range []struct {
		name   string
		args   []string
		err    error
		expBuf string
	}{
		{
			name: "default",
			args: []string{
				"default",
				"1",
			},
		},
		{
			name: "r2q",
			args: []string{
				"r2q",
				"10",
			},
		},
		{
			name: "debug",
			args: []string{
				"debug",
			},
			err: trafficctl.ErrNotImplemented,
		},
		{
			name: "offload",
			args: []string{
				"offload",
			},
			err: trafficctl.ErrNotImplemented,
		},
		{
			name: "direct_qlen",
			args: []string{
				"direct_qlen",
				"10",
			},
		},
		{
			name: "help",
			args: []string{
				"help",
			},
			expBuf: trafficctl.HTBHelp,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			var outBuf bytes.Buffer
			_, err := trafficctl.ParseHTBQDiscArgs(&outBuf, tt.args)
			if !errors.Is(err, tt.err) {
				t.Errorf("ParseHTBArgs(%q) = %v, not %v", tt.args, err, tt.err)
			}

			if tt.expBuf != "" {
				if tt.expBuf != outBuf.String() {
					t.Errorf("output !=  expected output")
				}
			}
		})
	}
}

func TestParseHFSCQdiscArgs(t *testing.T) {
	for _, tt := range []struct {
		name string
		args []string
		err  error
	}{
		{
			name: "default",
			args: []string{
				"default",
				"10",
			},
		},
		{
			name: "invalid",
			args: []string{
				"invalid",
			},
			err: trafficctl.ErrInvalidArg,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			var outbuf bytes.Buffer
			if _, err := trafficctl.ParseHFSCQDiscArgs(&outbuf, tt.args); !errors.Is(err, tt.err) {
				t.Errorf("ParseHFSCQdiscArgs() = %v, not %v", err, tt.err)
			}
		})
	}
}
