// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"bytes"
	"fmt"
	"net"
	"reflect"
	"testing"

	"github.com/vishvananda/netlink"
)

func TestDefaultOptions(t *testing.T) {
	expected := options{
		modes: []string{},
		iKey:  -1,
		oKey:  -1,
		ttl:   -1,
		tos:   -1,
	}

	result := defaultOptions()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("defaultOptions() = %v, want %v", result, expected)
	}
}

func TestParseTunnel(t *testing.T) {
	tests := []struct {
		name     string
		cmd      cmd
		expected options
		wantErr  bool
	}{
		{
			name: "Valid mode and remote",
			cmd: cmd{
				Cursor: 2,
				Out:    new(bytes.Buffer),
				Args:   []string{"ip", "tunnel", "add", "tln", "mode", "gre", "remote", "127.0.0.2", "local", "128.0.0.2", "ttl", "20", "tos", "2", "ikey", "10", "okey", "10", "dev", "lo"},
			},
			expected: options{
				name:   "tln",
				mode:   "gre",
				modes:  []string{"gre", "ip6gre"},
				remote: "127.0.0.2",
				local:  "128.0.0.2",
				iKey:   10,
				oKey:   10,
				ttl:    20,
				tos:    2,
				dev:    "lo",
			},
			wantErr: false,
		},
		{
			name: "invalid tos",
			cmd: cmd{
				Cursor: 2,
				Out:    new(bytes.Buffer),
				Args:   []string{"ip", "tunnel", "add", "tos", "abc"},
			},
			wantErr: true,
		},
		{
			name: "invalid mode",
			cmd: cmd{
				Cursor: 2,
				Out:    new(bytes.Buffer),
				Args:   []string{"ip", "tunnel", "add", "mode", "abc"},
			},
			wantErr: true,
		},
		{
			name: "invalid ttl",
			cmd: cmd{
				Cursor: 2,
				Out:    new(bytes.Buffer),
				Args:   []string{"ip", "tunnel", "add", "ttl", "abc"},
			},
			wantErr: true,
		},
		{
			name: "invalid ikey",
			cmd: cmd{
				Cursor: 2,
				Out:    new(bytes.Buffer),
				Args:   []string{"ip", "tunnel", "add", "ikey", "abc"},
			},
			wantErr: true,
		},
		{
			name: "invalid okey",
			cmd: cmd{
				Cursor: 2,
				Out:    new(bytes.Buffer),
				Args:   []string{"ip", "tunnel", "add", "okey", "abc"},
			},
			wantErr: true,
		},
		{
			name: "ttl inherit & all modes",
			cmd: cmd{
				Cursor: 2,
				Out:    new(bytes.Buffer),
				Args:   []string{"ip", "tunnel", "add", "ttl", "inherit"},
			},
			expected: options{
				modes: allTunnelTypes,
				iKey:  -1,
				oKey:  -1,
				ttl:   0,
				tos:   -1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.cmd.parseTunnel()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTunnel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !reflect.DeepEqual(result, &tt.expected) {
					t.Errorf("parseTunnel() = %v, want %v", result, &tt.expected)
				}
			}
		})
	}
}

func TestFilterTunnels(t *testing.T) {
	greTun := netlink.Gretun{LinkAttrs: netlink.LinkAttrs{Name: "link1"}, Local: net.ParseIP("127.0.0.2"), Remote: net.ParseIP("126.0.0.2"), Ttl: 9}
	greTun2 := netlink.Gretun{LinkAttrs: netlink.LinkAttrs{Name: "link1"}, Local: net.ParseIP("127.0.0.2"), Remote: net.ParseIP("126.0.0.3"), Ttl: 9}
	greTun3 := netlink.Gretun{LinkAttrs: netlink.LinkAttrs{Name: "link1"}, Local: net.ParseIP("127.0.0.3"), Remote: net.ParseIP("126.0.0.2"), Ttl: 9}
	greTun4 := netlink.Gretun{LinkAttrs: netlink.LinkAttrs{Name: "link1"}, Local: net.ParseIP("127.0.0.2"), Remote: net.ParseIP("126.0.0.2"), Ttl: 9, IKey: 10}
	greTun5 := netlink.Gretun{LinkAttrs: netlink.LinkAttrs{Name: "link1"}, Local: net.ParseIP("127.0.0.2"), Remote: net.ParseIP("126.0.0.2"), Ttl: 9, OKey: 10}
	greTun6 := netlink.Gretun{LinkAttrs: netlink.LinkAttrs{Name: "link1"}, Local: net.ParseIP("127.0.0.2"), Remote: net.ParseIP("126.0.0.2"), Ttl: 10}
	greTun7 := netlink.Gretun{LinkAttrs: netlink.LinkAttrs{Name: "link1"}, Local: net.ParseIP("127.0.0.2"), Remote: net.ParseIP("126.0.0.2"), Ttl: 9, Tos: 10}
	greTun8 := netlink.Gretun{LinkAttrs: netlink.LinkAttrs{Name: "link2"}, Local: net.ParseIP("127.0.0.2"), Remote: net.ParseIP("126.0.0.2"), Ttl: 9}
	ipTun := netlink.Iptun{LinkAttrs: netlink.LinkAttrs{Name: "link1"}}

	tests := []struct {
		name     string
		links    []netlink.Link
		options  *options
		expected []netlink.Link
	}{
		{
			name: "Filter by opts 1",
			links: []netlink.Link{
				&greTun,
				&greTun8,
			},
			options: &options{
				modes:  []string{"gre"},
				remote: "126.0.0.2",
				local:  "127.0.0.2",
				name:   "link1",
				ttl:    9,
			},
			expected: []netlink.Link{
				&greTun,
			},
		},
		{
			name: "Filter by opts 2",
			links: []netlink.Link{
				&greTun,
				&ipTun,
				&greTun2,
				&greTun3,
				&greTun4,
				&greTun5,
				&greTun6,
				&greTun7,
				&greTun8,
			},
			options: &options{
				modes:  []string{"gre"},
				remote: "126.0.0.2",
				local:  "127.0.0.2",
				name:   "link1",
				ttl:    9,
			},
			expected: []netlink.Link{
				&greTun,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterTunnels(tt.links, tt.options)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("filterTunnels() = %#v, want %#v", result, tt.expected)
			}
		})
	}
}

func TestEqualRemotes(t *testing.T) {
	tests := []struct {
		name     string
		link     netlink.Link
		remote   string
		expected bool
	}{
		{
			name:     "Empty remote",
			link:     &netlink.Gretun{Remote: net.ParseIP("192.168.1.1")},
			remote:   "",
			expected: true,
		},
		{
			name:     "Remote is 'any'",
			link:     &netlink.Gretun{Remote: net.ParseIP("192.168.1.1")},
			remote:   "any",
			expected: true,
		},
		{
			name:     "Matching remote IP for Gretun",
			link:     &netlink.Gretun{Remote: net.ParseIP("192.168.1.1")},
			remote:   "192.168.1.1",
			expected: true,
		},
		{
			name:     "Non-matching remote IP for Gretun",
			link:     &netlink.Gretun{Remote: net.ParseIP("192.168.1.1")},
			remote:   "192.168.1.2",
			expected: false,
		},
		{
			name:     "Matching remote IP for Iptun",
			link:     &netlink.Iptun{Remote: net.ParseIP("192.168.1.1")},
			remote:   "192.168.1.1",
			expected: true,
		},
		{
			name:     "Non-matching remote IP for Iptun",
			link:     &netlink.Iptun{Remote: net.ParseIP("192.168.1.1")},
			remote:   "192.168.1.2",
			expected: false,
		},
		{
			name:     "Matching remote IP for Ip6tnl",
			link:     &netlink.Ip6tnl{Remote: net.ParseIP("2001:db8::1")},
			remote:   "2001:db8::1",
			expected: true,
		},
		{
			name:     "Non-matching remote IP for Ip6tnl",
			link:     &netlink.Ip6tnl{Remote: net.ParseIP("2001:db8::1")},
			remote:   "2001:db8::2",
			expected: false,
		},
		{
			name:     "Matching remote IP for Vti",
			link:     &netlink.Vti{Remote: net.ParseIP("192.168.1.1")},
			remote:   "192.168.1.1",
			expected: true,
		},
		{
			name:     "Non-matching remote IP for Vti",
			link:     &netlink.Vti{Remote: net.ParseIP("192.168.1.1")},
			remote:   "192.168.1.2",
			expected: false,
		},
		{
			name:     "Matching remote IP for Sittun",
			link:     &netlink.Sittun{Remote: net.ParseIP("192.168.1.1")},
			remote:   "192.168.1.1",
			expected: true,
		},
		{
			name:     "Non-matching remote IP for Sittun",
			link:     &netlink.Sittun{Remote: net.ParseIP("192.168.1.1")},
			remote:   "192.168.1.2",
			expected: false,
		},
		{
			name:     "Unsupported link type",
			link:     &netlink.Dummy{LinkAttrs: netlink.LinkAttrs{Name: "dummy0"}},
			remote:   "192.168.1.1",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := equalRemotes(tt.link, tt.remote)
			if result != tt.expected {
				t.Errorf("equalRemotes() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEqualLocals(t *testing.T) {
	tests := []struct {
		name     string
		link     netlink.Link
		local    string
		expected bool
	}{
		{
			name:     "Empty local",
			link:     &netlink.Gretun{Local: net.ParseIP("192.168.1.1")},
			local:    "",
			expected: true,
		},
		{
			name:     "Local is 'any'",
			link:     &netlink.Gretun{Local: net.ParseIP("192.168.1.1")},
			local:    "any",
			expected: true,
		},
		{
			name:     "Matching local IP for Gretun",
			link:     &netlink.Gretun{Local: net.ParseIP("192.168.1.1")},
			local:    "192.168.1.1",
			expected: true,
		},
		{
			name:     "Non-matching local IP for Gretun",
			link:     &netlink.Gretun{Local: net.ParseIP("192.168.1.1")},
			local:    "192.168.1.2",
			expected: false,
		},
		{
			name:     "Matching local IP for Iptun",
			link:     &netlink.Iptun{Local: net.ParseIP("192.168.1.1")},
			local:    "192.168.1.1",
			expected: true,
		},
		{
			name:     "Non-matching local IP for Iptun",
			link:     &netlink.Iptun{Local: net.ParseIP("192.168.1.1")},
			local:    "192.168.1.2",
			expected: false,
		},
		{
			name:     "Matching local IP for Ip6tnl",
			link:     &netlink.Ip6tnl{Local: net.ParseIP("2001:db8::1")},
			local:    "2001:db8::1",
			expected: true,
		},
		{
			name:     "Non-matching local IP for Ip6tnl",
			link:     &netlink.Ip6tnl{Local: net.ParseIP("2001:db8::1")},
			local:    "2001:db8::2",
			expected: false,
		},
		{
			name:     "Matching local IP for Vti",
			link:     &netlink.Vti{Local: net.ParseIP("192.168.1.1")},
			local:    "192.168.1.1",
			expected: true,
		},
		{
			name:     "Non-matching local IP for Vti",
			link:     &netlink.Vti{Local: net.ParseIP("192.168.1.1")},
			local:    "192.168.1.2",
			expected: false,
		},
		{
			name:     "Matching local IP for Sittun",
			link:     &netlink.Sittun{Local: net.ParseIP("192.168.1.1")},
			local:    "192.168.1.1",
			expected: true,
		},
		{
			name:     "Non-matching local IP for Sittun",
			link:     &netlink.Sittun{Local: net.ParseIP("192.168.1.1")},
			local:    "192.168.1.2",
			expected: false,
		},
		{
			name:     "Unsupported link type",
			link:     &netlink.Dummy{LinkAttrs: netlink.LinkAttrs{Name: "dummy0"}},
			local:    "192.168.1.1",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := equalLocals(tt.link, tt.local)
			if result != tt.expected {
				t.Errorf("equalLocals() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEqualTTL(t *testing.T) {
	tests := []struct {
		name     string
		link     netlink.Link
		ttl      int
		expected bool
	}{
		{
			name:     "TTL is -1",
			link:     &netlink.Gretun{Ttl: 64},
			ttl:      -1,
			expected: true,
		},
		{
			name:     "TTL is 0",
			link:     &netlink.Gretun{Ttl: 64},
			ttl:      0,
			expected: true,
		},
		{
			name:     "TTL is 255",
			link:     &netlink.Gretun{Ttl: 64},
			ttl:      255,
			expected: true,
		},
		{
			name:     "Matching TTL for Gretun",
			link:     &netlink.Gretun{Ttl: 64},
			ttl:      64,
			expected: true,
		},
		{
			name:     "Non-matching TTL for Gretun",
			link:     &netlink.Gretun{Ttl: 64},
			ttl:      128,
			expected: false,
		},
		{
			name:     "Matching TTL for Iptun",
			link:     &netlink.Iptun{Ttl: 64},
			ttl:      64,
			expected: true,
		},
		{
			name:     "Non-matching TTL for Iptun",
			link:     &netlink.Iptun{Ttl: 64},
			ttl:      128,
			expected: false,
		},
		{
			name:     "Matching TTL for Ip6tnl",
			link:     &netlink.Ip6tnl{Ttl: 64},
			ttl:      64,
			expected: true,
		},
		{
			name:     "Non-matching TTL for Ip6tnl",
			link:     &netlink.Ip6tnl{Ttl: 64},
			ttl:      128,
			expected: false,
		},
		{
			name:     "Vti link type",
			link:     &netlink.Vti{},
			ttl:      64,
			expected: true,
		},
		{
			name:     "Matching TTL for Sittun",
			link:     &netlink.Sittun{Ttl: 64},
			ttl:      64,
			expected: true,
		},
		{
			name:     "Non-matching TTL for Sittun",
			link:     &netlink.Sittun{Ttl: 64},
			ttl:      128,
			expected: false,
		},
		{
			name:     "Unsupported link type",
			link:     &netlink.Dummy{LinkAttrs: netlink.LinkAttrs{Name: "dummy0"}},
			ttl:      64,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := equalTTL(tt.link, tt.ttl)
			if result != tt.expected {
				t.Errorf("equalTTL() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEqualTOS(t *testing.T) {
	tests := []struct {
		name     string
		link     netlink.Link
		tos      int
		expected bool
	}{
		{
			name:     "TOS is -1",
			link:     &netlink.Gretun{Tos: 16},
			tos:      -1,
			expected: true,
		},
		{
			name:     "Matching TOS for Gretun",
			link:     &netlink.Gretun{Tos: 16},
			tos:      16,
			expected: true,
		},
		{
			name:     "Non-matching TOS for Gretun",
			link:     &netlink.Gretun{Tos: 16},
			tos:      32,
			expected: false,
		},
		{
			name:     "Matching TOS for Iptun",
			link:     &netlink.Iptun{Tos: 16},
			tos:      16,
			expected: true,
		},
		{
			name:     "Non-matching TOS for Iptun",
			link:     &netlink.Iptun{Tos: 16},
			tos:      32,
			expected: false,
		},
		{
			name:     "Matching TOS for Ip6tnl",
			link:     &netlink.Ip6tnl{Tos: 16},
			tos:      16,
			expected: true,
		},
		{
			name:     "Non-matching TOS for Ip6tnl",
			link:     &netlink.Ip6tnl{Tos: 16},
			tos:      32,
			expected: false,
		},
		{
			name:     "Vti link type",
			link:     &netlink.Vti{},
			tos:      16,
			expected: true,
		},
		{
			name:     "Matching TOS for Sittun",
			link:     &netlink.Sittun{Tos: 16},
			tos:      16,
			expected: true,
		},
		{
			name:     "Non-matching TOS for Sittun",
			link:     &netlink.Sittun{Tos: 16},
			tos:      32,
			expected: false,
		},
		{
			name:     "Unsupported link type",
			link:     &netlink.Dummy{LinkAttrs: netlink.LinkAttrs{Name: "dummy0"}},
			tos:      16,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := equalTOS(tt.link, tt.tos)
			if result != tt.expected {
				t.Errorf("equalTOS() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEqualIKey(t *testing.T) {
	tests := []struct {
		name     string
		link     netlink.Link
		iKey     int
		expected bool
	}{
		{
			name:     "IKey is -1",
			link:     &netlink.Gretun{IKey: 1234},
			iKey:     -1,
			expected: true,
		},
		{
			name:     "Matching IKey for Gretun",
			link:     &netlink.Gretun{IKey: 1234},
			iKey:     1234,
			expected: true,
		},
		{
			name:     "Non-matching IKey for Gretun",
			link:     &netlink.Gretun{IKey: 1234},
			iKey:     5678,
			expected: false,
		},
		{
			name:     "Matching IKey for Vti",
			link:     &netlink.Vti{IKey: 1234},
			iKey:     1234,
			expected: true,
		},
		{
			name:     "Non-matching IKey for Vti",
			link:     &netlink.Vti{IKey: 1234},
			iKey:     5678,
			expected: false,
		},
		{
			name:     "Unsupported link type",
			link:     &netlink.Dummy{LinkAttrs: netlink.LinkAttrs{Name: "dummy0"}},
			iKey:     1234,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := equalIKey(tt.link, tt.iKey)
			if result != tt.expected {
				t.Errorf("equalIKey() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEqualOKey(t *testing.T) {
	tests := []struct {
		name     string
		link     netlink.Link
		oKey     int
		expected bool
	}{
		{
			name:     "OKey is -1",
			link:     &netlink.Gretun{OKey: 1234},
			oKey:     -1,
			expected: true,
		},
		{
			name:     "Matching OKey for Gretun",
			link:     &netlink.Gretun{OKey: 1234},
			oKey:     1234,
			expected: true,
		},
		{
			name:     "Non-matching OKey for Gretun",
			link:     &netlink.Gretun{OKey: 1234},
			oKey:     5678,
			expected: false,
		},
		{
			name:     "Matching OKey for Vti",
			link:     &netlink.Vti{OKey: 1234},
			oKey:     1234,
			expected: true,
		},
		{
			name:     "Non-matching OKey for Vti",
			link:     &netlink.Vti{OKey: 1234},
			oKey:     5678,
			expected: false,
		},
		{
			name:     "Unsupported link type",
			link:     &netlink.Dummy{LinkAttrs: netlink.LinkAttrs{Name: "dummy0"}},
			oKey:     1234,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := equalOKey(tt.link, tt.oKey)
			if result != tt.expected {
				t.Errorf("equalOKey() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNormalizeOptsForAddingTunnel(t *testing.T) {
	tests := []struct {
		name     string
		op       *options
		expected *options
		err      error
	}{
		{
			name: "Missing mode",
			op:   &options{},
			expected: &options{
				mode: "",
				name: "",
				iKey: 0,
				oKey: 0,
				ttl:  0,
				tos:  0,
			},
			err: fmt.Errorf("tunnel mode is required"),
		},
		{
			name: "Missing name for gre mode",
			op:   &options{mode: "gre"},
			expected: &options{
				mode: "gre",
				name: "gre0",
				iKey: 0,
				oKey: 0,
				ttl:  0,
				tos:  0,
			},
			err: nil,
		},
		{
			name: "Negative iKey and oKey",
			op:   &options{mode: "ipip", iKey: -1, oKey: -1},
			expected: &options{
				mode: "ipip",
				name: "tuln0",
				iKey: 0,
				oKey: 0,
				ttl:  0,
				tos:  0,
			},
			err: nil,
		},
		{
			name: "Negative ttl and tos",
			op:   &options{mode: "sit", ttl: -1, tos: -1},
			expected: &options{
				mode: "sit",
				name: "sit0",
				iKey: 0,
				oKey: 0,
				ttl:  0,
				tos:  0,
			},
			err: nil,
		},
		{
			name: "All fields provided",
			op:   &options{mode: "vti", name: "custom0", iKey: 123, oKey: 456, ttl: 64, tos: 16},
			expected: &options{
				mode: "vti",
				name: "custom0",
				iKey: 123,
				oKey: 456,
				ttl:  64,
				tos:  16,
			},
			err: nil,
		},
		{
			name: "ip6tln",
			op:   &options{mode: "ip6tln", name: "", iKey: 123, oKey: 456, ttl: 64, tos: 16},
			expected: &options{
				mode: "ip6tln",
				name: "ip6tnl0",
				iKey: 123,
				oKey: 456,
				ttl:  64,
				tos:  16,
			},
			err: nil,
		},
		{
			name: "vti",
			op:   &options{mode: "vti", name: "", iKey: 123, oKey: 456, ttl: 64, tos: 16},
			expected: &options{
				mode: "vti",
				name: "ip_vti0",
				iKey: 123,
				oKey: 456,
				ttl:  64,
				tos:  16,
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := normalizeOptsForAddingTunnel(tt.op)
			if (err != nil) && (err.Error() != tt.err.Error()) {
				t.Errorf("normalizeOptsForAddingTunnel() error = %v, want %v", err, tt.err)
			}
			if !reflect.DeepEqual(tt.op, tt.expected) {
				t.Errorf("normalizeOptsForAddingTunnel() = %v, want %v", tt.op, tt.expected)
			}
		})
	}
}

func TestPrintTunnels(t *testing.T) {
	tests := []struct {
		name    string
		tunnels []netlink.Link
		json    bool
		want    string
		wantErr bool
	}{
		{
			name: "Single GRE tunnel",
			tunnels: []netlink.Link{
				&netlink.Gretun{
					LinkAttrs: netlink.LinkAttrs{Name: "gre0"},
					Local:     net.ParseIP("192.168.1.1"),
					Remote:    net.ParseIP("192.168.1.2"),
					Ttl:       64,
				},
			},
			json: false,
			want: "gre0: gre/ip remote 192.168.1.2 local 192.168.1.1 ttl 64\n",
		},
		{
			name: "Single IP tunnel",
			tunnels: []netlink.Link{
				&netlink.Iptun{
					LinkAttrs: netlink.LinkAttrs{Name: "ip0"},
					Local:     net.ParseIP("192.168.1.1"),
					Remote:    net.ParseIP("192.168.1.2"),
					Ttl:       64,
				},
			},
			json: false,
			want: "ip0: any/ip remote 192.168.1.2 local 192.168.1.1 ttl 64\n",
		},
		{
			name: "Single IPv6 tunnel",
			tunnels: []netlink.Link{
				&netlink.Ip6tnl{
					LinkAttrs: netlink.LinkAttrs{Name: "ipv60"},
					Local:     net.ParseIP("::1"),
					Remote:    net.ParseIP("::2"),
					Ttl:       64,
				},
			},
			json: false,
			want: "ipv60: ip6tln/ip remote ::2 local ::1 ttl 64\n",
		},
		{
			name: "Single VTI tunnel",
			tunnels: []netlink.Link{
				&netlink.Vti{
					LinkAttrs: netlink.LinkAttrs{Name: "vti0"},
					Local:     net.ParseIP("192.168.1.1"),
					Remote:    net.ParseIP("192.168.1.2"),
				},
			},
			json: false,
			want: "vti0: ip/ip remote 192.168.1.2 local 192.168.1.1 ttl inherit\n",
		},
		{
			name: "Single SIT tunnel",
			tunnels: []netlink.Link{
				&netlink.Sittun{
					LinkAttrs: netlink.LinkAttrs{Name: "sit0"},
					Local:     net.ParseIP("192.168.1.1"),
					Remote:    net.ParseIP("192.168.1.2"),
					Ttl:       64,
				},
			},
			json: false,
			want: "sit0: sit/ip remote 192.168.1.2 local 192.168.1.1 ttl 64\n",
		},
		{
			name: "Unsupported tunnel type",
			tunnels: []netlink.Link{
				&netlink.Dummy{
					LinkAttrs: netlink.LinkAttrs{Name: "dummy0"},
				},
			},
			json:    false,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			cmd := cmd{
				Out: &out,
				Opts: flags{
					JSON: tt.json,
				},
			}
			cmd.Opts.JSON = tt.json

			err := cmd.printTunnels(tt.tunnels)
			if err != nil && !tt.wantErr {
				t.Fatalf("printTunnels() error = %v", err)
			}

			if !tt.wantErr {
				got := out.String()
				if got != tt.want {
					t.Errorf("printTunnels() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
