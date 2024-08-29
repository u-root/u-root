// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"bytes"
	"testing"

	"github.com/vishvananda/netlink"
)

func TestParseAddrAddReplace(t *testing.T) {
	tests := []struct {
		name             string
		cmd              cmd
		wantValidLft     int
		wantPreferredLft int
		wantErr          bool
	}{
		{
			name: "default",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "addr", "add", "127.0.0.1/24", "dev", "lo"},
				Out:    new(bytes.Buffer),
			},
		},
		{
			name: "frv lfts",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "addr", "add", "127.0.0.1/24", "dev", "lo", "valid_lft", "forever", "preferred_lft", "forever"},
				Out:    new(bytes.Buffer),
			},
		},
		{
			name: "10 lfts",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "addr", "add", "127.0.0.1/24", "dev", "lo", "valid_lft", "10", "preferred_lft", "10"},
				Out:    new(bytes.Buffer),
			},
			wantValidLft:     10,
			wantPreferredLft: 10,
		},
		{
			name: "invalid valid_lft",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "addr", "add", "127.0.0.1/24", "dev", "lo", "valid_lft", "abc", "preferred_lft", "abc"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "invalid lft",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "addr", "add", "127.0.0.1/24", "dev", "lo", "valid_lft", "10", "preferred_lft", "abc"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "invalid addr",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "addr", "add", "abcde"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "invalid dev",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "addr", "add", "127.0.0.1/24", "dev", "fjghyy"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, addr, err := tt.cmd.parseAddrAddReplace()
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

func TestParseAddrShow(t *testing.T) {
	tests := []struct {
		name     string
		cmd      cmd
		dev      string
		typeName string
		wantErr  bool
	}{
		{
			name: "default",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "addr", "show"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "values",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "addr", "show", "dev", "lo", "type", "bridge"},
				Out:    new(bytes.Buffer),
			},
			dev:      "lo",
			typeName: "bridge",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			link, typeStr, err := tt.cmd.parseAddrShow()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAddrShow() error = %v, wantErr %t", err, tt.wantErr)
			}

			if !tt.wantErr {
				if link.Attrs().Name != tt.dev {
					t.Errorf("link.Name = %v, want %s", link.Attrs().Name, tt.dev)
				}
				if typeStr != tt.typeName {
					t.Errorf("type = %v, want %s", typeStr, tt.typeName)
				}
			}
		})
	}
}

func TestParseAddrFlush(t *testing.T) {
	tests := []struct {
		name    string
		cmd     cmd
		dev     string
		wantErr bool
	}{
		{
			name: "default",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "addr", "flush", "dev", "lo"},
				Out:    new(bytes.Buffer),
			},
			dev: "lo",
		},
		{
			name: "values",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "addr", "flush", "dev", "lo", "scope", "host", "label", "label"},
				Out:    new(bytes.Buffer),
			},
			dev: "lo",
		},
		{
			name: "fail on dev",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "addr", "flush", "deva"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "integer scope",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "addr", "flush", "dev", "lo", "scope", "2", "label", "label"},
				Out:    new(bytes.Buffer),
			},
			dev: "lo",
		},
		{
			name: "scope wrong arg",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "addr", "flush", "dev", "lo", "scope", "abcdef", "label", "label"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			link, _, err := tt.cmd.parseAddrFlush()
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
