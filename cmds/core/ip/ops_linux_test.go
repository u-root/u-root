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
	"github.com/vishvananda/netlink"
)

func TestPrintLinkJSON(t *testing.T) {
	tests := []struct {
		name     string
		links    []linkData
		opts     flags
		expected string
	}{
		{
			name: "Single link with IPv4 address",
			links: []linkData{
				{
					attrs: &netlink.LinkAttrs{
						Name:         "eth0",
						Flags:        net.FlagUp,
						OperState:    netlink.OperUp,
						HardwareAddr: net.HardwareAddr{0x00, 0x1a, 0x2b, 0x3c, 0x4d, 0x5e},
						Index:        1,
						MTU:          1500,
						Group:        0,
						TxQLen:       1000,
					},
					typeName: "device",
					addresses: []netlink.Addr{
						{
							IPNet: &net.IPNet{
								IP:   net.IPv4(192, 168, 1, 1),
								Mask: net.CIDRMask(24, 32),
							},
							Broadcast:   net.IPv4(192, 168, 1, 255),
							Scope:       int(netlink.SCOPE_HOST),
							Label:       "eth0",
							PreferedLft: 0,
							ValidLft:    0,
						},
					},
				},
			},
			opts: flags{JSON: true, Prettify: true},
			expected: `[
    {
        "ifindex": 1,
        "ifname": "eth0",
        "flags": [
            "up"
        ],
        "mtu": 1500,
        "operstate": "up",
        "group": "default",
        "txqlen": 1000,
        "link_type": "device",
        "address": "00:1a:2b:3c:4d:5e",
        "addr_info": [
            {
                "ip": "inet",
                "local": "192.168.1.1",
                "prefixlen": "ffffff00",
                "broadcast": "192.168.1.255",
                "scope": "host",
                "label": "eth0",
                "valid_life_time": "0sec",
                "preferred_life_time": "0sec"
            }
        ]
    }
]`,
		},
		{
			name: "Single link with IPv6 address",
			links: []linkData{
				{
					attrs: &netlink.LinkAttrs{
						Name:         "eth0",
						Flags:        net.FlagUp,
						OperState:    netlink.OperUp,
						HardwareAddr: net.HardwareAddr{0x00, 0x1a, 0x2b, 0x3c, 0x4d, 0x5e},
						Index:        1,
						MTU:          1500,
						Group:        0,
						TxQLen:       1000,
					},
					typeName: "device",
					addresses: []netlink.Addr{
						{
							IPNet: &net.IPNet{
								IP:   net.ParseIP("2001:db8::1"),
								Mask: net.CIDRMask(64, 128),
							},
							Broadcast:   nil,
							Scope:       int(netlink.SCOPE_HOST),
							Label:       "eth0",
							PreferedLft: 3600,
							ValidLft:    7200,
						},
					},
				},
			},
			opts: flags{JSON: true, Prettify: true},
			expected: `[
    {
        "ifindex": 1,
        "ifname": "eth0",
        "flags": [
            "up"
        ],
        "mtu": 1500,
        "operstate": "up",
        "group": "default",
        "txqlen": 1000,
        "link_type": "device",
        "address": "00:1a:2b:3c:4d:5e",
        "addr_info": [
            {
                "ip": "inet6",
                "local": "2001:db8::1",
                "prefixlen": "ffffffffffffffff0000000000000000",
                "scope": "host",
                "label": "eth0",
                "valid_life_time": "7200sec",
                "preferred_life_time": "3600sec"
            }
        ]
    }
]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			cmd := cmd{
				Out:  &out,
				Opts: tt.opts,
			}

			err := cmd.printLinkJSON(tt.links)
			if err != nil {
				t.Fatalf("printLinkJSON() error = %v", err)
			}

			if c := cmp.Diff(out.String(), tt.expected); c != "" {
				t.Errorf("printLinkJSON() = %v", c)
			}
		})
	}
}

func TestShowLinks(t *testing.T) {
	tests := []struct {
		name     string
		links    []linkData
		opts     flags
		expected string
	}{
		{
			name: "Single link with IPv4 address JSON",
			links: []linkData{
				{
					attrs: &netlink.LinkAttrs{
						Name:         "eth0",
						Flags:        net.FlagUp,
						OperState:    netlink.OperUp,
						HardwareAddr: net.HardwareAddr{0x00, 0x1a, 0x2b, 0x3c, 0x4d, 0x5e},
						Index:        1,
						MTU:          1500,
						Group:        0,
						TxQLen:       1000,
						Statistics: &netlink.LinkStatistics{
							RxPackets: 100,
							TxPackets: 200,
							RxBytes:   1000,
							TxBytes:   2000,
							RxErrors:  10,
							TxErrors:  20,
							RxDropped: 1,
							TxDropped: 2,
						},
					},
					typeName: "device",
					addresses: []netlink.Addr{
						{
							IPNet: &net.IPNet{
								IP:   net.IPv4(192, 168, 1, 1),
								Mask: net.CIDRMask(24, 32),
							},
							Broadcast:   net.IPv4(192, 168, 1, 255),
							Scope:       int(netlink.SCOPE_HOST),
							Label:       "eth0",
							PreferedLft: 0,
							ValidLft:    0,
						},
					},
				},
			},
			opts: flags{JSON: true, Prettify: true},
			expected: `[
    {
        "ifindex": 1,
        "ifname": "eth0",
        "flags": [
            "up"
        ],
        "mtu": 1500,
        "operstate": "up",
        "group": "default",
        "txqlen": 1000,
        "link_type": "device",
        "address": "00:1a:2b:3c:4d:5e",
        "addr_info": [
            {
                "ip": "inet",
                "local": "192.168.1.1",
                "prefixlen": "ffffff00",
                "broadcast": "192.168.1.255",
                "scope": "host",
                "label": "eth0",
                "valid_life_time": "0sec",
                "preferred_life_time": "0sec"
            }
        ]
    }
]`,
		},
		{
			name: "Single link with IPv4 address",
			links: []linkData{
				{
					attrs: &netlink.LinkAttrs{
						Name:         "eth0",
						Flags:        net.FlagUp,
						OperState:    netlink.OperUp,
						HardwareAddr: net.HardwareAddr{0x00, 0x1a, 0x2b, 0x3c, 0x4d, 0x5e},
						Index:        1,
						MTU:          1500,
						Group:        0,
						TxQLen:       1000,
					},
					typeName: "device",
					addresses: []netlink.Addr{
						{
							IPNet: &net.IPNet{
								IP:   net.IPv4(192, 168, 1, 1),
								Mask: net.CIDRMask(24, 32),
							},
							Broadcast:   net.IPv4(192, 168, 1, 255),
							Scope:       int(netlink.SCOPE_HOST),
							Label:       "eth0",
							PreferedLft: 0,
							ValidLft:    0,
						},
					},
				},
			},
			opts:     flags{},
			expected: "1: eth0: <UP> mtu 1500 state UP group default\n    link/ 00:1a:2b:3c:4d:5e\n    inet 192.168.1.1 brd 192.168.1.255 scope host eth0\n       valid_lft 0sec preferred_lft 0sec\n",
		},
		{
			name: "Single link with IPv4 address brief",
			links: []linkData{
				{
					attrs: &netlink.LinkAttrs{
						Name:         "eth0",
						Flags:        net.FlagUp,
						OperState:    netlink.OperUp,
						HardwareAddr: net.HardwareAddr{0x00, 0x1a, 0x2b, 0x3c, 0x4d, 0x5e},
						Index:        1,
						MTU:          1500,
						Group:        0,
						TxQLen:       1000,
					},
					typeName: "device",
					addresses: []netlink.Addr{
						{
							IPNet: &net.IPNet{
								IP:   net.IPv4(192, 168, 1, 1),
								Mask: net.CIDRMask(24, 32),
							},
							Broadcast:   net.IPv4(192, 168, 1, 255),
							Scope:       int(netlink.SCOPE_HOST),
							Label:       "eth0",
							PreferedLft: 0,
							ValidLft:    0,
						},
					},
				},
			},
			opts:     flags{Brief: true},
			expected: "eth0                 up         192.168.1.1\n",
		},
		{
			name: "Single link brief",
			links: []linkData{
				{
					attrs: &netlink.LinkAttrs{
						Name:         "eth0",
						Flags:        net.FlagUp,
						OperState:    netlink.OperUp,
						HardwareAddr: net.HardwareAddr{0x00, 0x1a, 0x2b, 0x3c, 0x4d, 0x5e},
						Index:        1,
						MTU:          1500,
						Group:        0,
						TxQLen:       1000,
					},
					typeName: "device",
				},
			},
			opts:     flags{Brief: true},
			expected: "eth0                      up         00:1a:2b:3c:4d:5e   <UP>\n",
		},
		{
			name: "Stats",
			links: []linkData{
				{
					attrs: &netlink.LinkAttrs{
						Name:         "eth0",
						Flags:        net.FlagUp,
						OperState:    netlink.OperUp,
						HardwareAddr: net.HardwareAddr{0x00, 0x1a, 0x2b, 0x3c, 0x4d, 0x5e},
						Index:        1,
						MTU:          1500,
						Group:        0,
						TxQLen:       1000,
						Statistics: &netlink.LinkStatistics{
							RxPackets: 100,
							TxPackets: 200,
							RxBytes:   1000,
							TxBytes:   2000,
							RxErrors:  10,
							TxErrors:  20,
							RxDropped: 1,
							TxDropped: 2,
						},
					},
					typeName: "device",
					addresses: []netlink.Addr{
						{
							IPNet: &net.IPNet{
								IP:   net.IPv4(192, 168, 1, 1),
								Mask: net.CIDRMask(24, 32),
							},
							Broadcast:   net.IPv4(192, 168, 1, 255),
							Scope:       int(netlink.SCOPE_HOST),
							Label:       "eth0",
							PreferedLft: 0,
							ValidLft:    0,
						},
					},
				},
			},
			opts: flags{Stats: true},
			expected: `1: eth0: <UP> mtu 1500 state UP group default
    link/ 00:1a:2b:3c:4d:5e
    RX: bytes 1000 packets 100 errors 10 dropped 1 missed 0 mcast 0
    TX: bytes 2000 packets 200 errors 20 dropped 2 carrier 0 collsns 0
    inet 192.168.1.1 brd 192.168.1.255 scope host eth0
       valid_lft 0sec preferred_lft 0sec
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			cmd := cmd{
				Out:  &out,
				Opts: tt.opts,
			}

			err := cmd.printLinks(false, tt.links)
			if err != nil {
				t.Fatalf("showLinks() error = %v", err)
			}

			if diff := cmp.Diff(out.String(), tt.expected); diff != "" {
				t.Errorf("showLinks() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
