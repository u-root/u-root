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
