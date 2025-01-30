// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"reflect"
	"testing"

	"github.com/vishvananda/netlink"
)

func TestParseAddrAddReplace(t *testing.T) {
	tests := []struct {
		name             string
		Args             []string
		wantValidLft     int
		wantPreferredLft int
		wantErr          bool
	}{
		{
			name: "default",
			Args: []string{"ip", "addr", "add", "127.0.0.1/24", "dev", "lo"},
		},
		{
			name: "frv lfts",
			Args: []string{"ip", "addr", "add", "127.0.0.1/24", "dev", "lo", "valid_lft", "forever", "preferred_lft", "forever"},
		},
		{
			name:             "10 lfts",
			Args:             []string{"ip", "addr", "add", "127.0.0.1/24", "dev", "lo", "valid_lft", "10", "preferred_lft", "10"},
			wantValidLft:     10,
			wantPreferredLft: 10,
		},
		{
			name:    "invalid valid_lft",
			Args:    []string{"ip", "addr", "add", "127.0.0.1/24", "dev", "lo", "valid_lft", "abc", "preferred_lft", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid lft",
			Args:    []string{"ip", "addr", "add", "127.0.0.1/24", "dev", "lo", "valid_lft", "10", "preferred_lft", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid addr",
			Args:    []string{"ip", "addr", "add", "abcde"},
			wantErr: true,
		},
		{
			name:    "invalid dev",
			Args:    []string{"ip", "addr", "add", "127.0.0.1/24", "dev", "fjghyy"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := cmd{
				Cursor: 2,
				Args:   tt.Args,
			}
			_, addr, err := cmd.parseAddrAddReplace()
			if (err != nil) != tt.wantErr {
				t.Errorf("vrf() = %v, want %t", err, tt.wantErr)
			}

			if !tt.wantErr {
				if addr.ValidLft != tt.wantValidLft {
					t.Errorf("valid_lft = %v, want %v", addr.ValidLft, tt.wantValidLft)
				}
				if addr.PreferedLft != tt.wantPreferredLft {
					t.Errorf("preferred_lft = %v, want %v", addr.PreferedLft, tt.wantPreferredLft)
				}

			}
		})
	}
}

func TestParseAddrFlush(t *testing.T) {
	tests := []struct {
		name    string
		Args    []string
		dev     string
		wantErr bool
	}{
		{
			name: "default",
			Args: []string{"ip", "addr", "flush", "dev", "lo"},
			dev:  "lo",
		},
		{
			name: "values",
			Args: []string{"ip", "addr", "flush", "dev", "lo", "scope", "host", "label", "label"},
			dev:  "lo",
		},
		{
			name:    "fail on dev",
			Args:    []string{"ip", "addr", "flush", "deva"},
			wantErr: true,
		},
		{
			name: "integer scope",
			Args: []string{"ip", "addr", "flush", "dev", "lo", "scope", "2", "label", "label"},
			dev:  "lo",
		},
		{
			name:    "scope wrong arg",
			Args:    []string{"ip", "addr", "flush", "dev", "lo", "scope", "abcdef", "label", "label"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := cmd{
				Cursor: 2,
				Args:   tt.Args,
			}
			link, _, err := cmd.parseAddrFlush()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAddrShow() error = %v, wantErr %t", err, tt.wantErr)
			}

			if !tt.wantErr {
				if link.Attrs().Name != tt.dev {
					t.Errorf("link.Name = %v, want %s", link.Attrs().Name, tt.dev)
				}
			}
		})
	}
}

func TestParseAddrShow(t *testing.T) {
	tests := []struct {
		name      string
		Args      []string
		wantName  string
		wantTypes []string
	}{
		{
			name:     "Show address with device name",
			Args:     []string{"ip", "addr", "show", "eth0"},
			wantName: "eth0",
		},
		{
			name:      "Show address with type",
			Args:      []string{"ip", "addr", "show", "type", "dummy"},
			wantTypes: []string{"dummy"},
		},
		{
			name:      "Show address with device name and type",
			Args:      []string{"ip", "addr", "show", "eth0", "type", "dummy"},
			wantName:  "eth0",
			wantTypes: []string{"dummy"},
		},
		{
			name:      "Show address with multiple types",
			Args:      []string{"ip", "addr", "show", "type", "dummy", "veth"},
			wantTypes: []string{"dummy", "veth"},
		},
		{
			name:     "Show address with device name using 'dev'",
			Args:     []string{"ip", "addr", "show", "dev", "eth0"},
			wantName: "eth0",
		},
		{
			name:      "Show address with device name using 'dev' and type",
			Args:      []string{"ip", "addr", "show", "dev", "eth0", "type", "dummy"},
			wantName:  "eth0",
			wantTypes: []string{"dummy"},
		},
		{
			name:      "Show address with type and device name using 'dev'",
			Args:      []string{"ip", "addr", "show", "type", "dummy", "dev", "eth0"},
			wantName:  "eth0",
			wantTypes: []string{"dummy"},
		},
		{
			name:     "Show address with multiple device names",
			Args:     []string{"ip", "addr", "show", "eth0", "eth1"},
			wantName: "eth0",
		},
		{
			name:     "Show address with multiple device names using 'dev'",
			Args:     []string{"ip", "addr", "show", "dev", "eth0", "dev", "eth1"},
			wantName: "eth0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := cmd{
				Cursor: 2,
				Args:   tt.Args,
			}
			gotName, gotTypes := cmd.parseAddrShow()
			if gotName != tt.wantName {
				t.Errorf("parseAddrShow() gotName = %v, want %v", gotName, tt.wantName)
			}
			if !reflect.DeepEqual(gotTypes, tt.wantTypes) {
				t.Errorf("parseAddrShow() gotTypes = %v, want %v", gotTypes, tt.wantTypes)
			}
		})
	}
}

func TestSkipAddr(t *testing.T) {
	tests := []struct {
		name     string
		addr     netlink.Addr
		a        netlink.Addr
		expected bool
	}{
		{
			name:     "Different Scope",
			addr:     netlink.Addr{Scope: 1},
			a:        netlink.Addr{Scope: 2},
			expected: true,
		},
		{
			name:     "Same Scope",
			addr:     netlink.Addr{Scope: 1},
			a:        netlink.Addr{Scope: 1},
			expected: false,
		},
		{
			name:     "Different Label",
			addr:     netlink.Addr{Label: "eth0"},
			a:        netlink.Addr{Label: "eth1"},
			expected: true,
		},
		{
			name:     "Same Label",
			addr:     netlink.Addr{Label: "eth0"},
			a:        netlink.Addr{Label: "eth0"},
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := skipAddr(test.addr, test.a)
			if result != test.expected {
				t.Errorf("skipAddr(%v, %v) = %v; want %v", test.addr, test.a, result, test.expected)
			}
		})
	}
}
