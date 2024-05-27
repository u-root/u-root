// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/netstat"
)

func TestRun(t *testing.T) {
	for _, tt := range []struct {
		name string
		Flags
		expErr error
	}{
		{
			name:   "SucessSocketsDefault",
			Flags:  Flags{},
			expErr: nil,
		},
		{
			name:   "SucessSocketsIPv4",
			Flags:  Flags{ipv4: true},
			expErr: nil,
		},
		{
			name:   "SucessSocketsIPv6",
			Flags:  Flags{ipv6: true},
			expErr: nil,
		},
		{
			name: "SucessRoute",
			Flags: Flags{
				route: true,
				ipv4:  true,
				ipv6:  true,
			},
			expErr: nil,
		},
		{
			name: "FailRouteCacheIPv4",
			Flags: Flags{
				route:      true,
				routecache: true,
				ipv4:       true,
			},
			expErr: netstat.ErrRouteCacheIPv6only,
		},
		{
			name:   "SuccessInterfaces",
			Flags:  Flags{interfaces: true},
			expErr: nil,
		},
		{
			name:   "SuccessIface_eth0",
			Flags:  Flags{iface: "eth0"},
			expErr: nil,
		},
		{
			name: "SuccessStats",
			Flags: Flags{
				stats: true,
				ipv4:  true,
				ipv6:  true,
			},
			expErr: nil,
		},
		{
			name: "Success_xorFlagsUsage",
			Flags: Flags{
				stats: true,
				route: true,
				ipv4:  true,
			},
			expErr: nil,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var out strings.Builder
			if err := run(tt.Flags, &out); !errors.Is(err, tt.expErr) {
				t.Errorf("evalFlags() failed: %v, want: %v", err, tt.expErr)
			}
		})
	}
}

func TestXorFlags(t *testing.T) {
	tests := []struct {
		name   string
		flags  []bool
		result bool
	}{
		{"None set", []bool{false, false, false}, true},
		{"One set", []bool{true, false, false}, true},
		{"Two set", []bool{true, true, false}, false},
		{"All set", []bool{true, true, true}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := xorFlags(test.flags...); got != test.result {
				t.Errorf("xorFlags() = %v, want %v", got, test.result)
			}
		})
	}
}

func TestEvalProtocols(t *testing.T) {
	tests := []struct {
		name          string
		tcp, udp      bool
		udpl, raw     bool
		unix, ipv4    bool
		ipv6          bool
		wantProtocols int
	}{
		{"TCP and IPv4", true, false, false, false, false, true, false, 1},
		{"UDP and IPv6", false, true, false, false, false, false, true, 1},
		{"TCP and UDP with IPv4", true, true, false, false, false, true, false, 2},
		{"All protocols with IPv4 and IPv6", true, true, true, true, true, true, true, 9},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			protocols, err := evalProtocols(test.tcp, test.udp, test.udpl, test.raw, test.unix, test.ipv4, test.ipv6)
			if err != nil {
				t.Errorf("evalProtocols() error = %v", err)
				return
			}
			if len(protocols) != test.wantProtocols {
				t.Errorf("evalProtocols() len = %v, want %v", len(protocols), test.wantProtocols)
			}
		})
	}
}
