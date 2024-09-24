// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"bytes"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

func TestRouteTypeToString(t *testing.T) {
	tests := []struct {
		routeType int
		expected  string
	}{
		{1, "unicast"},
		{2, "local"},
		{3, "broadcast"},
		{5, "multicast"},
		{6, "blackhole"},
		{7, "unreachable"},
		{8, "prohibit"},
		{9, "throw"},
		{10, "nat"},
		{99, "unknown"}, // Test for an unknown route type
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := routeTypeToString(tt.routeType)
			if result != tt.expected {
				t.Errorf("routeTypeToString(%d) = %s; want %s", tt.routeType, result, tt.expected)
			}
		})
	}
}

func TestParseRouteAddAppendReplaceDel(t *testing.T) {
	_, dst, err := net.ParseCIDR("192.0.0.2/24")
	if err != nil {
		t.Fatalf("Failed to parse CIDR: %v", err)
	}

	tests := []struct {
		name         string
		args         []string
		addr         string
		expected     netlink.Route
		expectedLink string
		wantErr      bool
	}{
		{
			name:    "fail to parse dst",
			addr:    "abc",
			wantErr: true,
		},
		{
			name:         "Add route with valid arguments",
			addr:         "192.0.0.2/24",
			args:         []string{"dev", "lo"},
			expectedLink: "lo",
			expected: netlink.Route{
				Dst: dst,
			},
			wantErr: false,
		},
		{
			name:         "all opts",
			addr:         "192.0.0.2/24",
			args:         []string{"dev", "lo", "tos", "1", "table", "1", "proto", "1", "scope", "1", "metric", "1", "mtu", "1", "advmss", "1", "rtt", "1", "rttvar", "1", "reordering", "1", "window", "1", "cwnd", "1", "initcwnd", "1", "ssthresh", "1", "initrwnd", "1", "realms", "1", "src", "127.0.0.2", "rto_min", "1", "hoplimit", "1", "congctl", "a", "features", "1", "quickack", "1", "fastopen_no_cookie", "1"},
			expectedLink: "lo",
			expected: netlink.Route{
				Dst:              dst,
				Tos:              1,
				Table:            1,
				Protocol:         1,
				Scope:            1,
				Priority:         1,
				MTU:              1,
				AdvMSS:           1,
				Rtt:              1,
				RttVar:           1,
				Reordering:       1,
				Window:           1,
				Cwnd:             1,
				InitCwnd:         1,
				Realm:            1,
				Src:              net.ParseIP("127.0.0.2"),
				RtoMin:           1,
				Hoplimit:         1,
				InitRwnd:         1,
				Congctl:          "a",
				Features:         1,
				QuickACK:         1,
				FastOpenNoCookie: 1,
			},
			wantErr: false,
		},
		{
			name:         "quickack 0",
			addr:         "192.0.0.2/24",
			args:         []string{"dev", "lo", "quickack", "0"},
			expectedLink: "lo",
			expected: netlink.Route{
				Dst:      dst,
				QuickACK: 0,
			},
			wantErr: false,
		},
		{
			name:         "fastopen_no_cookie 0",
			addr:         "192.0.0.2/24",
			args:         []string{"dev", "lo", "fastopen_no_cookie", "0"},
			expectedLink: "lo",
			expected: netlink.Route{
				Dst:              dst,
				FastOpenNoCookie: 0,
			},
			wantErr: false,
		},
		{
			name:    "invalid arg",
			addr:    "192.0.0.2/24",
			args:    []string{"dev", "lo", "abc"},
			wantErr: true,
		},
		{
			name:    "fastopen_no_cookie invalid",
			addr:    "192.0.0.2/24",
			args:    []string{"dev", "lo", "fastopen_no_cookie", "2"},
			wantErr: true,
		},
		{
			name:    "quickack invalid",
			addr:    "192.0.0.2/24",
			args:    []string{"dev", "lo", "quickack", "2"},
			wantErr: true,
		},
		{
			name:    "features invalid",
			addr:    "192.0.0.2/24",
			args:    []string{"dev", "lo", "features", "ac"},
			wantErr: true,
		},
		{
			name:    "initrwnd invalid",
			addr:    "192.0.0.2/24",
			args:    []string{"dev", "lo", "initrwnd", "ac"},
			wantErr: true,
		},
		{
			name:    "hoplimit invalid",
			addr:    "192.0.0.2/24",
			args:    []string{"dev", "lo", "hoplimit", "ac"},
			wantErr: true,
		},
		{
			name:    "rto_min invalid",
			addr:    "192.0.0.2/24",
			args:    []string{"dev", "lo", "rto_min", "ac"},
			wantErr: true,
		},
		{
			name:    "src invalid",
			addr:    "192.0.0.2/24",
			args:    []string{"dev", "lo", "src", "ac"},
			wantErr: true,
		},
		{
			name:    "realms invalid",
			addr:    "192.0.0.2/24",
			args:    []string{"dev", "lo", "realms", "ac"},
			wantErr: true,
		},
		{
			name:    "ssthresh invalid",
			addr:    "192.0.0.2/24",
			args:    []string{"dev", "lo", "ssthresh", "ac"},
			wantErr: true,
		},
		{
			name:    "initcwnd invalid",
			addr:    "192.0.0.2/24",
			args:    []string{"dev", "lo", "initcwnd", "ac"},
			wantErr: true,
		},
		{
			name:    "cwnd invalid",
			addr:    "192.0.0.2/24",
			args:    []string{"dev", "lo", "cwnd", "ac"},
			wantErr: true,
		},
		{
			name:    "window invalid",
			addr:    "192.0.0.2/24",
			args:    []string{"dev", "lo", "window", "ac"},
			wantErr: true,
		},
		{
			name:    "reordering invalid",
			addr:    "192.0.0.2/24",
			args:    []string{"dev", "lo", "reordering", "ac"},
			wantErr: true,
		},
		{
			name:    "rttvar invalid",
			addr:    "192.0.0.2/24",
			args:    []string{"dev", "lo", "rttvar", "ac"},
			wantErr: true,
		},
		{
			name:    "rtt invalid",
			addr:    "192.0.0.2/24",
			args:    []string{"dev", "lo", "rtt", "ac"},
			wantErr: true,
		},
		{
			name:    "advmss invalid",
			addr:    "192.0.0.2/24",
			args:    []string{"dev", "lo", "advmss", "ac"},
			wantErr: true,
		},
		{
			name:    "mtu invalid",
			addr:    "192.0.0.2/24",
			args:    []string{"dev", "lo", "mtu", "ac"},
			wantErr: true,
		},
		{
			name:    "metric invalid",
			addr:    "192.0.0.2/24",
			args:    []string{"dev", "lo", "metric", "ac"},
			wantErr: true,
		},
		{
			name:    "scope invalid",
			addr:    "192.0.0.2/24",
			args:    []string{"dev", "lo", "scope", "ac"},
			wantErr: true,
		},
		{
			name:    "proto invalid",
			addr:    "192.0.0.2/24",
			args:    []string{"dev", "lo", "proto", "ac"},
			wantErr: true,
		},
		{
			name:    "table invalid",
			addr:    "192.0.0.2/24",
			args:    []string{"dev", "lo", "table", "ac"},
			wantErr: true,
		},
		{
			name:    "tos invalid",
			addr:    "192.0.0.2/24",
			args:    []string{"dev", "lo", "tos", "ac"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			cmd := cmd{
				Cursor: -1,
				Args:   tt.args,
				Out:    &out,
			}
			route, link, err := cmd.parseRouteAddAppendReplaceDel(tt.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRouteAddAppendReplaceDel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if link != tt.expectedLink {
					t.Errorf("parseRouteAddAppendReplaceDel() = %v, want %v", link, tt.expectedLink)
				}

				if diff := cmp.Diff(*route, tt.expected); diff != "" {
					t.Errorf("parseRouteAddAppendReplaceDel() = %v", diff)
				}
			}
		})
	}
}

func TestParseRouteShowListFlush(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantFilter *netlink.Route
		wantMask   uint64
		wantRoot   *net.IPNet
		wantMatch  *net.IPNet
		wantExact  *net.IPNet
		wantErr    bool
	}{
		{
			name: "Valid scope and table",
			args: []string{"scope", "2", "table", "2", "proto", "1", "type", "unicast"},
			wantFilter: &netlink.Route{
				Scope:    2,
				Table:    2,
				Protocol: 1,
				Type:     unix.RTN_UNICAST,
			},
			wantMask: netlink.RT_FILTER_SCOPE | netlink.RT_FILTER_TABLE | netlink.RT_FILTER_PROTOCOL | netlink.RT_FILTER_TYPE,
			wantErr:  false,
		},
		{
			name:    "Invalid scope",
			args:    []string{"scope", "invalid"},
			wantErr: true,
		},
		{
			name:       "Valid root prefix",
			args:       []string{"root", "192.168.1.0/24"},
			wantFilter: &netlink.Route{},
			wantRoot: &net.IPNet{
				IP:   net.IPv4(192, 168, 1, 0),
				Mask: net.CIDRMask(24, 32),
			},
			wantErr: false,
		},
		{
			name:    "Invalid root prefix",
			args:    []string{"root", "invalid_prefix"},
			wantErr: true,
		},
		{
			name:    "Invalid table",
			args:    []string{"table", "a"},
			wantErr: true,
		},
		{
			name:    "Invalid proto",
			args:    []string{"proto", "a"},
			wantErr: true,
		},
		{
			name:    "Invalid type",
			args:    []string{"type", "as"},
			wantErr: true,
		},
		{
			name:    "Invalid arg",
			args:    []string{"arg", "as"},
			wantErr: true,
		},
		{
			name:       "Valid match prefix",
			args:       []string{"match", "10.0.0.0/8"},
			wantFilter: &netlink.Route{},
			wantMatch: &net.IPNet{
				IP:   net.IPv4(10, 0, 0, 0),
				Mask: net.CIDRMask(8, 32),
			},
			wantErr: false,
		},
		{
			name:    "Invalid match prefix",
			args:    []string{"match", "invalid_prefix"},
			wantErr: true,
		},
		{
			name:    "Invalid exact prefix",
			args:    []string{"exact", "invalid_prefix"},
			wantErr: true,
		},
		{
			name:       "Valid exact prefix",
			args:       []string{"exact", "172.16.0.0/12"},
			wantFilter: &netlink.Route{},
			wantExact: &net.IPNet{
				IP:   net.IPv4(172, 16, 0, 0),
				Mask: net.CIDRMask(12, 32),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := cmd{
				Cursor: -1,
				Args:   tt.args,
			}
			gotFilter, gotMask, gotRoot, gotMatch, gotExact, err := cmd.parseRouteShowListFlush()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRouteShowListFlush() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if diff := cmp.Diff(gotFilter, tt.wantFilter, cmpopts.IgnoreFields(netlink.Route{}, "Dst")); diff != "" {
					t.Errorf("parseRouteShowListFlush() filter mismatch (-want +got):\n%s", diff)
				}
				if gotMask != tt.wantMask {
					t.Errorf("parseRouteShowListFlush() mask = %v, want %v", gotMask, tt.wantMask)
				}
				if gotRoot != nil && tt.wantRoot != nil && !gotRoot.IP.Equal(tt.wantRoot.IP) {
					t.Errorf("parseRouteShowListFlush() root = %v, want %v", gotRoot, tt.wantRoot)
				}
				if gotMatch != nil && tt.wantMatch != nil && !gotMatch.IP.Equal(tt.wantMatch.IP) {
					t.Errorf("parseRouteShowListFlush() match = %v, want %v", gotMatch, tt.wantMatch)
				}
				if gotExact != nil && tt.wantExact != nil && !gotExact.IP.Equal(tt.wantExact.IP) {
					t.Errorf("parseRouteShowListFlush() exact = %v, want %v", gotExact, tt.wantExact)
				}
			}
		})
	}
}

func TestMatchRoutes(t *testing.T) {
	tests := []struct {
		name    string
		routes  []netlink.Route
		root    *net.IPNet
		match   *net.IPNet
		exact   *net.IPNet
		want    []netlink.Route
		wantErr bool
	}{
		{
			name: "Match root prefix",
			routes: []netlink.Route{
				{Dst: &net.IPNet{IP: net.IPv4(192, 168, 1, 1), Mask: net.CIDRMask(24, 32)}},
				{Dst: &net.IPNet{IP: net.IPv4(10, 0, 0, 1), Mask: net.CIDRMask(8, 32)}},
			},
			root: &net.IPNet{
				IP:   net.IPv4(192, 168, 1, 0),
				Mask: net.CIDRMask(24, 32),
			},
			want: []netlink.Route{
				{Dst: &net.IPNet{IP: net.IPv4(192, 168, 1, 1), Mask: net.CIDRMask(24, 32)}},
			},
			wantErr: false,
		},
		{
			name: "Match exact prefix",
			routes: []netlink.Route{
				{Dst: &net.IPNet{IP: net.IPv4(192, 168, 1, 1), Mask: net.CIDRMask(24, 32)}},
				{Dst: &net.IPNet{IP: net.IPv4(10, 0, 0, 1), Mask: net.CIDRMask(8, 32)}},
			},
			exact: &net.IPNet{
				IP:   net.IPv4(10, 0, 0, 1),
				Mask: net.CIDRMask(8, 32),
			},
			want: []netlink.Route{
				{Dst: &net.IPNet{IP: net.IPv4(10, 0, 0, 1), Mask: net.CIDRMask(8, 32)}},
			},
			wantErr: false,
		},
		{
			name: "Match prefix",
			routes: []netlink.Route{
				{Dst: &net.IPNet{IP: net.IPv4(192, 168, 1, 1), Mask: net.CIDRMask(24, 32)}},
				{Dst: &net.IPNet{IP: net.IPv4(10, 0, 0, 1), Mask: net.CIDRMask(8, 32)}},
			},
			match: &net.IPNet{
				IP:   net.IPv4(10, 0, 0, 0),
				Mask: net.CIDRMask(8, 32),
			},
			want: []netlink.Route{
				{Dst: &net.IPNet{IP: net.IPv4(10, 0, 0, 1), Mask: net.CIDRMask(8, 32)}},
			},
			wantErr: false,
		},
		{
			name: "No match",
			routes: []netlink.Route{
				{Dst: &net.IPNet{IP: net.IPv4(192, 168, 1, 1), Mask: net.CIDRMask(24, 32)}},
				{Dst: &net.IPNet{IP: net.IPv4(10, 0, 0, 1), Mask: net.CIDRMask(8, 32)}},
			},
			root: &net.IPNet{
				IP:   net.IPv4(172, 16, 0, 0),
				Mask: net.CIDRMask(12, 32),
			},
			want:    []netlink.Route{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := matchRoutes(tt.routes, tt.root, tt.match, tt.exact)
			if (err != nil) != tt.wantErr {
				t.Errorf("matchRoutes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("matchRoutes() = %v, want %v", got, tt.want)
				return
			}
			for i := range got {
				if !got[i].Dst.IP.Equal(tt.want[i].Dst.IP) || got[i].Dst.Mask.String() != tt.want[i].Dst.Mask.String() {
					t.Errorf("matchRoutes() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestDefaultRoute(t *testing.T) {
	tests := []struct {
		name     string
		cmd      cmd
		route    netlink.Route
		linkName string
		expected string
	}{
		{
			name: "Numeric and Details false",
			cmd: cmd{
				Opts: flags{},
				Out:  new(bytes.Buffer),
			},
			route: netlink.Route{
				Gw:       net.IPv4(192, 168, 1, 1),
				Protocol: 1,
				Priority: 100,
			},
			linkName: "eth0",
			expected: "default via 192.168.1.1 dev eth0 proto redirect metric 100\n",
		},
		{
			name: "Numeric true and Details false",
			cmd: cmd{
				Opts: flags{Numeric: true},

				Out: new(bytes.Buffer),
			},
			route: netlink.Route{
				Gw:       net.IPv4(192, 168, 1, 1),
				Protocol: 2,
				Priority: 200,
			},
			linkName: "eth1",
			expected: "default via 192.168.1.1 dev eth1 proto 2 metric 200\n",
		},
		{
			name: "Numeric false and Details true",
			cmd: cmd{
				Opts: flags{Details: true},
				Out:  new(bytes.Buffer),
			},
			route: netlink.Route{
				Gw:       net.IPv4(192, 168, 1, 1),
				Protocol: 1,
				Priority: 300,
				Type:     1,
			},
			linkName: "eth2",
			expected: "unicast default via 192.168.1.1 dev eth2 proto redirect metric 300\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			tt.cmd.Out = &out

			tt.cmd.defaultRoute(tt.route, tt.linkName)
			if got := out.String(); got != tt.expected {
				t.Errorf("defaultRoute() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestShowRoute(t *testing.T) {
	_, dst, err := net.ParseCIDR("192.0.0.4/24")
	if err != nil {
		t.Fatalf("Failed to parse CIDR: %v", err)
	}

	_, ipv6Dst, err := net.ParseCIDR("2001:db8::1/64")
	if err != nil {
		t.Fatalf("Failed to parse CIDR: %v", err)
	}

	tests := []struct {
		name     string
		cmd      cmd
		route    netlink.Route
		linkName string
		expected string
	}{
		{
			name: "IPv4 route with FAMILY_V4",
			cmd: cmd{
				Family: netlink.FAMILY_V4,
				Out:    new(bytes.Buffer),
			},
			route: netlink.Route{
				Dst:      dst,
				Protocol: 1,
				Scope:    netlink.SCOPE_HOST,
				Src:      net.IPv4(127, 0, 0, 1),
			},
			linkName: "eth0",
			expected: "192.0.0.0/24 dev eth0 proto redirect scope host src 127.0.0.1 metric 0\n",
		},
		{
			name: "IPv4 route with FAMILY_V6",
			cmd: cmd{
				Family: netlink.FAMILY_V6,
				Out:    new(bytes.Buffer),
			},
			route: netlink.Route{
				Dst: dst,
			},
			linkName: "eth1",
			expected: "",
		},
		{
			name: "IPv6 route with FAMILY_V4",
			cmd: cmd{
				Family: netlink.FAMILY_V4,
				Out:    new(bytes.Buffer),
			},
			route: netlink.Route{
				Dst: ipv6Dst,
			},
			linkName: "eth1",
			expected: "",
		},
		{
			name: "IPv6 route with FAMILY_V6",
			cmd: cmd{
				Family: netlink.FAMILY_V6,
				Out:    new(bytes.Buffer),
			},
			route: netlink.Route{
				Dst: ipv6Dst,
			},
			linkName: "eth1",
			expected: "2001:db8::/64 dev eth1 proto unspec metric 0\n",
		},
		{
			name: "Mixed family with FAMILY_ALL (IPv4)",
			cmd: cmd{
				Family: netlink.FAMILY_ALL,
				Out:    new(bytes.Buffer),
			},
			route: netlink.Route{
				Dst:      dst,
				Protocol: 1,
				Scope:    netlink.SCOPE_HOST,
				Src:      net.IPv4(127, 0, 0, 1),
			},
			linkName: "eth2",
			expected: "192.0.0.0/24 dev eth2 proto redirect scope host src 127.0.0.1 metric 0\n",
		},
		{
			name: "IPv4 route with numeric",
			cmd: cmd{
				Family: netlink.FAMILY_V4,
				Out:    new(bytes.Buffer),
				Opts:   flags{Numeric: true},
			},
			route: netlink.Route{
				Dst:      dst,
				Protocol: 1,
				Scope:    netlink.SCOPE_HOST,
				Src:      net.IPv4(127, 0, 0, 1),
			},
			linkName: "eth0",
			expected: "192.0.0.0/24 dev eth0 proto 1 scope 254 src 127.0.0.1 metric 0\n",
		},
		{
			name: "IPv6 route with numeric",
			cmd: cmd{
				Family: netlink.FAMILY_V6,
				Out:    new(bytes.Buffer),
				Opts:   flags{Numeric: true},
			},
			route: netlink.Route{
				Dst: ipv6Dst,
			},
			linkName: "eth1",
			expected: "2001:db8::/64 dev eth1 proto 0 metric 0\n",
		},
		{
			name: "IPv4 route with details",
			cmd: cmd{
				Family: netlink.FAMILY_V4,
				Out:    new(bytes.Buffer),
				Opts:   flags{Details: true},
			},
			route: netlink.Route{
				Dst:      dst,
				Protocol: 1,
				Scope:    netlink.SCOPE_HOST,
				Src:      net.IPv4(127, 0, 0, 1),
				Type:     unix.RTN_UNICAST,
			},
			linkName: "eth0",
			expected: "unicast 192.0.0.0/24 dev eth0 proto redirect scope host src 127.0.0.1 metric 0\n",
		},
		{
			name: "IPv6 route with details",
			cmd: cmd{
				Family: netlink.FAMILY_V6,
				Out:    new(bytes.Buffer),
				Opts:   flags{Details: true},
			},
			route: netlink.Route{
				Dst:  ipv6Dst,
				Type: unix.RTN_UNICAST,
			},
			linkName: "eth1",
			expected: "unicast 2001:db8::/64 dev eth1 proto unspec metric 0\n",
		},
		{
			name: "IPv6 route with Gateway",
			cmd: cmd{
				Family: netlink.FAMILY_V6,
				Out:    new(bytes.Buffer),
			},
			route: netlink.Route{
				Dst:  ipv6Dst,
				Type: unix.RTN_UNICAST,
				Gw:   net.IPv6loopback,
			},
			linkName: "eth1",
			expected: "2001:db8::/64 via ::1 dev eth1 proto unspec metric 0\n",
		},
	}

	for _, tt := range tests {
		var out bytes.Buffer
		t.Run(tt.name, func(t *testing.T) {
			tt.cmd.Out = &out
			tt.cmd.showRoute(tt.route, tt.linkName)
			if got := out.String(); got != tt.expected {
				t.Errorf("showRoute() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseRouteGet(t *testing.T) {
	tests := []struct {
		name    string
		cmd     cmd
		want    netlink.RouteGetOptions
		wantErr bool
	}{
		{
			name:    "Valid input with all options",
			cmd:     cmd{Cursor: -1, Args: []string{"from", "127.0.0.1", "oif", "1", "iif", "2", "vrf", "vrf0"}},
			want:    netlink.RouteGetOptions{SrcAddr: net.ParseIP("127.0.0.1"), Oif: "1", Iif: "2", VrfName: "vrf0"},
			wantErr: false,
		},
		{
			name:    "Invalid input",
			cmd:     cmd{Cursor: -1, Args: []string{"arg"}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cmd.parseRouteGet()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRouteGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if diff := cmp.Diff(*got, tt.want); diff != "" {
					t.Errorf("parseRouteGet() = %v", diff)
				}
			}
		})
	}
}

func TestShowRoutes(t *testing.T) {
	tests := []struct {
		name       string
		opts       flags
		routes     []netlink.Route
		ifaceNames []string
		wantOutput string
		wantErr    bool
	}{
		{
			name: "JSON output",
			opts: flags{JSON: true},
			routes: []netlink.Route{
				{
					Dst: &net.IPNet{
						IP:   net.ParseIP("192.168.1.0"),
						Mask: net.CIDRMask(24, 32),
					},
					Scope:    netlink.SCOPE_UNIVERSE,
					Protocol: 2,
					Src:      net.ParseIP("127.0.0.3"),
					Flags:    unix.RTNH_F_ONLINK,
				},
			},
			ifaceNames: []string{"eth0"},
			wantOutput: `[{"dst":"192.168.1.0/24","dev":"eth0","protocol":"kernel","scope":"universe","prefsrc":"127.0.0.3","flags":["onlink"]}]`,
			wantErr:    false,
		},
		{
			name: "JSON output with numeric",
			opts: flags{JSON: true, Numeric: true},
			routes: []netlink.Route{
				{
					Dst: &net.IPNet{
						IP:   net.ParseIP("192.168.1.0"),
						Mask: net.CIDRMask(24, 32),
					},
					Scope:    netlink.SCOPE_UNIVERSE,
					Protocol: 2,
				},
			},
			ifaceNames: []string{"eth0"},
			wantOutput: `[{"dst":"192.168.1.0/24","dev":"eth0","protocol":"2","scope":"0","prefsrc":""}]`,
			wantErr:    false,
		},
		{
			name: "normal output",
			routes: []netlink.Route{
				{
					Dst: &net.IPNet{
						IP:   net.ParseIP("192.168.1.0"),
						Mask: net.CIDRMask(24, 32),
					},
					Scope:    netlink.SCOPE_UNIVERSE,
					Protocol: 2,
					Src:      net.ParseIP("127.0.0.3"),
					Flags:    unix.RTNH_F_ONLINK,
				},
			},
			ifaceNames: []string{"eth0"},
			wantOutput: `192.168.1.0/24 dev eth0 proto kernel scope global src 127.0.0.3 metric 0
`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			cmd := cmd{
				Opts: tt.opts,
				Out:  &out,
			}

			err := cmd.showRoutes(tt.routes, tt.ifaceNames)
			if (err != nil) != tt.wantErr {
				t.Errorf("showRoutes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOutput := out.String(); gotOutput != tt.wantOutput {
				t.Errorf("showRoutes() output = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}
