// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/netstat"
)

func TestRun(t *testing.T) {
	for _, tt := range []struct {
		name string
		cmd
		expErr error
	}{
		{
			name:   "SucessSocketsDefault",
			cmd:    cmd{},
			expErr: nil,
		},
		{
			name:   "SucessSocketsIPv4",
			cmd:    cmd{ipv4: true},
			expErr: nil,
		},
		{
			name:   "SucessSocketsIPv6",
			cmd:    cmd{ipv6: true},
			expErr: nil,
		},
		{
			name: "SucessRoute",
			cmd: cmd{
				route: true,
				ipv4:  true,
				ipv6:  true,
			},
			expErr: nil,
		},
		{
			name: "FailRouteCacheIPv4",
			cmd: cmd{
				route:      true,
				routecache: true,
				ipv4:       true,
			},
			expErr: netstat.ErrRouteCacheIPv6only,
		},
		{
			name:   "SuccessInterfaces",
			cmd:    cmd{interfaces: true},
			expErr: nil,
		},
		{
			name:   "SuccessIface_eth0",
			cmd:    cmd{iface: "eth0"},
			expErr: nil,
		},
		{
			name: "SuccessStats",
			cmd: cmd{
				stats: true,
				ipv4:  true,
				ipv6:  true,
			},
			expErr: nil,
		},
		{
			name: "Success_xorFlagsUsage",
			cmd: cmd{
				stats: true,
				route: true,
				ipv4:  true,
			},
			expErr: errMutualExcludeFlags,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var out strings.Builder
			tt.cmd.out = &out
			if err := tt.cmd.run(); !errors.Is(err, tt.expErr) {
				t.Errorf("cmd.run() failed: %v, want: %v", err, tt.expErr)
			}
		})
	}
}

func TestDefaultFlags(t *testing.T) {
	tests := []struct {
		name    string
		cmdline []string
		want    cmd
	}{
		{"No command", []string{"netstat"}, cmd{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cmd := command(nil, test.cmdline)
			if reflect.DeepEqual(cmd, test.want) {
				t.Errorf("\ngot: %+v\nwant: %+v", cmd, test.want)
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
