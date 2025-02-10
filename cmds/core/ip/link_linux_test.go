// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"bytes"
	"math"
	"net"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/vishvananda/netlink"
)

func TestParseLinkAdd(t *testing.T) {
	tests := []struct {
		name      string
		Args      []string
		wantType  string
		wantAttrs netlink.LinkAttrs
		wantErr   bool
	}{
		{
			name:     "Successful parsing",
			Args:     []string{"ip", "link", "add", "name", "eth0", "type", "dummy", "txqueuelen", "1000", "address", "00:0c:29:3e:5c:7f", "mtu", "1500", "index", "1", "numtxqueues", "2", "numrxqueues", "2"},
			wantType: "dummy",
			wantAttrs: netlink.LinkAttrs{
				Name:         "eth0",
				TxQLen:       1000,
				HardwareAddr: net.HardwareAddr{0x00, 0x0c, 0x29, 0x3e, 0x5c, 0x7f},
				MTU:          1500,
				Index:        1,
				NumTxQueues:  2,
				NumRxQueues:  2,
			},
			wantErr: false,
		},
		{
			name:    "invalid txqlen",
			Args:    []string{"ip", "link", "add", "name", "eth0", "type", "dummy", "txqueuelen", "abc", "address", "00:0c:29:3e:5c:7f", "mtu", "1500", "index", "1", "numtxqueues", "2", "numrxqueues", "2"},
			wantErr: true,
		},
		{
			name:    "invalid address",
			Args:    []string{"ip", "link", "add", "name", "eth0", "type", "dummy", "txqueuelen", "1000", "address", "00:0c:29:3e:5c:7f:00", "mtu", "1500", "index", "1", "numtxqueues", "2", "numrxqueues", "2"},
			wantErr: true,
		},
		{
			name:    "invalid mtu",
			Args:    []string{"ip", "link", "add", "name", "eth0", "type", "dummy", "txqueuelen", "1000", "address", "00:0c:29:3e:5c:7f", "mtu", "abc", "index", "1", "numtxqueues", "2", "numrxqueues", "2"},
			wantErr: true,
		},
		{
			name:    "invalid index",
			Args:    []string{"ip", "link", "add", "name", "eth0", "type", "dummy", "txqueuelen", "1000", "address", "00:0c:29:3e:5c:7f", "mtu", "1500", "index", "abc", "numtxqueues", "2", "numrxqueues", "2"},
			wantErr: true,
		},
		{
			name:    "invalid numtxqueues",
			Args:    []string{"ip", "link", "add", "name", "eth0", "type", "dummy", "txqueuelen", "1000", "address", "00:0c:29:3e:5c:7f", "mtu", "1500", "index", "1", "numtxqueues", "abc", "numrxqueues", "2"},
			wantErr: true,
		},
		{
			name:    "invalid numrxqueues",
			Args:    []string{"ip", "link", "add", "name", "eth0", "type", "dummy", "txqueuelen", "1000", "address", "00:0c:29:3e:5c:7f", "mtu", "1500", "index", "1", "numtxqueues", "2", "numrxqueues", "abc"},
			wantErr: true,
		},
		{
			name:    "invalid arg",
			Args:    []string{"ip", "link", "add", "name", "eth0", "abc"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := cmd{
				Cursor: 2,
				Args:   tt.Args,
			}
			gotType, gotAttrs, err := cmd.parseLinkAdd()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseLinkAttrs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotType != tt.wantType {
				t.Errorf("parseLinkAttrs() gotType = %v, want %v", gotType, tt.wantType)
			}
			if !reflect.DeepEqual(gotAttrs, tt.wantAttrs) {
				t.Errorf("parseLinkAttrs() gotAttrs = %v, want %v", gotAttrs, tt.wantAttrs)
			}
		})
	}
}

func TestParseLinkShow(t *testing.T) {
	tests := []struct {
		name      string
		Args      []string
		wantName  string
		wantTypes []string
	}{
		{
			name:     "Show link with device name",
			Args:     []string{"ip", "link", "show", "eth0"},
			wantName: "eth0",
		},
		{
			name:      "Show link with type",
			Args:      []string{"ip", "link", "show", "type", "dummy"},
			wantTypes: []string{"dummy"},
		},
		{
			name:      "Show link with device name and type",
			Args:      []string{"ip", "link", "show", "eth0", "type", "dummy"},
			wantName:  "eth0",
			wantTypes: []string{"dummy"},
		},
		{
			name:      "Show link with multiple types",
			Args:      []string{"ip", "link", "show", "type", "dummy", "veth"},
			wantTypes: []string{"dummy", "veth"},
		},
		{
			name:     "Show link with device name using 'dev'",
			Args:     []string{"ip", "link", "show", "dev", "eth0"},
			wantName: "eth0",
		},
		{
			name:      "Show link with device name using 'dev' and type",
			Args:      []string{"ip", "link", "show", "dev", "eth0", "type", "dummy"},
			wantName:  "eth0",
			wantTypes: []string{"dummy"},
		},
		{
			name:      "Show link with type and device name using 'dev'",
			Args:      []string{"ip", "link", "show", "type", "dummy", "dev", "eth0"},
			wantName:  "eth0",
			wantTypes: []string{"dummy"},
		},
		{
			name:     "Show link with multiple device names",
			Args:     []string{"ip", "link", "show", "eth0", "eth1"},
			wantName: "eth0",
		},
		{
			name:     "Show link with multiple device names using 'dev'",
			Args:     []string{"ip", "link", "show", "dev", "eth0", "dev", "eth1"},
			wantName: "eth0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := cmd{
				Cursor: 2,
				Args:   tt.Args,
			}
			gotName, gotTypes := cmd.parseLinkShow()
			if gotName != tt.wantName {
				t.Errorf("parseLinkShow() gotName = %v, want %v", gotName, tt.wantName)
			}
			if !reflect.DeepEqual(gotTypes, tt.wantTypes) {
				t.Errorf("parseLinkShow() gotTypes = %v, want %v", gotTypes, tt.wantTypes)
			}
		})
	}
}

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
                "prefixlen": "24",
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
                "prefixlen": "64",
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

func TestPrintLinks(t *testing.T) {
	// Avoid problems building this test on 32-bit archs. See field netlink.Addr.ValidLft in
	// test case "Single link with ValidLft set to max uint32 (forever)" below.
	uint32fix := func(i uint32) int { return int(i) }

	tests := []struct {
		name           string
		withAddresses  bool
		links          []linkData
		opts           flags
		expectedValues []string
	}{
		// Default Format
		{
			name:          "Single link without addresses, default format",
			withAddresses: false,
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
			opts:           flags{},
			expectedValues: []string{"eth0", "UP", "mtu 1500"},
		},
		{
			name:          "Single link with addresses, default format",
			withAddresses: true,
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
							Label:       "eth0-label",
							PreferedLft: 0,
							ValidLft:    0,
						},
					},
				},
			},
			opts:           flags{},
			expectedValues: []string{"eth0", "UP", "mtu 1500", "192.168.1.1/24", "eth0-label"},
		},
		{
			name:          "Single link with addresses, default format, IPv4 address with broadcast",
			withAddresses: true,
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
							Label:       "eth0-label",
							PreferedLft: 0,
							ValidLft:    0,
						},
					},
				},
			},
			opts:           flags{},
			expectedValues: []string{"eth0", "UP", "mtu 1500", "192.168.1.1/24", "brd 192.168.1.255", "eth0-label"},
		},
		{
			name:          "Single link with addresses, default format, IPv6 address without broadcast",
			withAddresses: true,
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
							Label:       "eth0-label",
							PreferedLft: 3600,
							ValidLft:    7200,
						},
					},
				},
			},
			opts:           flags{},
			expectedValues: []string{"eth0", "UP", "mtu 1500", "2001:db8::1/64", "eth0-label"},
		},
		{
			name:          "Single link with master name, default format",
			withAddresses: false,
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
					typeName:   "device",
					masterName: "bridge0",
				},
			},
			opts:           flags{},
			expectedValues: []string{"eth0", "UP", "mtu 1500", "master bridge0"},
		},
		{
			name:          "Single link with numeric option, default format",
			withAddresses: false,
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
			opts:           flags{Numeric: true},
			expectedValues: []string{"eth0", "UP", "mtu 1500", "group 0"},
		},
		{
			name:          "Single link with details, default format",
			withAddresses: false,
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
					typeName:       "device",
					specificDevice: &netlink.Dummy{},
				},
			},
			opts:           flags{Details: true},
			expectedValues: []string{"eth0", "UP", "mtu 1500"},
		},
		{
			name:          "Single link with statistics, default format",
			withAddresses: false,
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
				},
			},
			opts:           flags{Stats: true},
			expectedValues: []string{"eth0", "UP", "mtu 1500", "RX:  bytes  packets errors dropped  missed   mcast", "1000", "100", "10", "1", "TX:  bytes  packets errors dropped carrier collsns", "2000", "200", "20", "2"},
		},
		{
			name:          "Multiple links, default format",
			withAddresses: false,
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
				{
					attrs: &netlink.LinkAttrs{
						Name:         "eth1",
						Flags:        net.FlagUp,
						OperState:    netlink.OperUp,
						HardwareAddr: net.HardwareAddr{0x00, 0x1a, 0x2b, 0x3c, 0x4d, 0x5f},
						Index:        2,
						MTU:          1500,
						Group:        0,
						TxQLen:       1000,
					},
					typeName: "device",
				},
			},
			opts:           flags{},
			expectedValues: []string{"eth0", "UP", "mtu 1500", "eth1", "UP", "mtu 1500"},
		},
		// Brief Format
		{
			name:          "Single link without addresses, brief format",
			withAddresses: false,
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
			opts:           flags{Brief: true},
			expectedValues: []string{"eth0", "UP", "00:1a:2b:3c:4d:5e"},
		},
		{
			name:          "Single link with addresses, brief format",
			withAddresses: true,
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
							Label:       "eth0-label",
							PreferedLft: 0,
							ValidLft:    0,
						},
					},
				},
			},
			opts:           flags{Brief: true},
			expectedValues: []string{"eth0", "UP", "192.168.1.1/24"},
		},
		{
			name:          "Multiple links, brief format",
			withAddresses: false,
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
				{
					attrs: &netlink.LinkAttrs{
						Name:         "eth1",
						Flags:        net.FlagUp,
						OperState:    netlink.OperUp,
						HardwareAddr: net.HardwareAddr{0x00, 0x1a, 0x2b, 0x3c, 0x4d, 0x5f},
						Index:        2,
						MTU:          1500,
						Group:        0,
						TxQLen:       1000,
					},
					typeName: "device",
				},
			},
			opts:           flags{Brief: true},
			expectedValues: []string{"eth0", "UP", "00:1a:2b:3c:4d:5e", "eth1", "UP", "00:1a:2b:3c:4d:5f"},
		},
		// Oneline Format
		{
			name:          "Single link without addresses, oneline format",
			withAddresses: false,
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
			opts:           flags{Oneline: true},
			expectedValues: []string{"eth0", "UP", "mtu 1500"},
		},
		{
			name:          "Single link with addresses, oneline format",
			withAddresses: true,
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
							Label:       "eth0-label",
							PreferedLft: 0,
							ValidLft:    0,
						},
					},
				},
			},
			opts:           flags{Oneline: true},
			expectedValues: []string{"eth0", "192.168.1.1/24", "eth0-label"},
		},
		{
			name:          "Single link with addresses, oneline format, IPv4 address with broadcast",
			withAddresses: true,
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
							Label:       "eth0-label",
							PreferedLft: 0,
							ValidLft:    0,
						},
					},
				},
			},
			opts:           flags{Oneline: true},
			expectedValues: []string{"eth0", "192.168.1.1/24", "brd 192.168.1.255", "eth0-label"},
		},
		{
			name:          "Single link with addresses, oneline format, IPv6 address without broadcast",
			withAddresses: true,
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
							Label:       "eth0-label",
							PreferedLft: 3600,
							ValidLft:    7200,
						},
					},
				},
			},
			opts:           flags{Oneline: true},
			expectedValues: []string{"eth0", "2001:db8::1/64", "eth0-label"},
		},
		{
			name:          "Single link with master name, oneline format",
			withAddresses: false,
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
					typeName:   "device",
					masterName: "bridge0",
				},
			},
			opts:           flags{Oneline: true},
			expectedValues: []string{"eth0", "UP", "mtu 1500", "master bridge0"},
		},
		{
			name:          "Single link with numeric option, oneline format",
			withAddresses: false,
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
			opts:           flags{Oneline: true, Numeric: true},
			expectedValues: []string{"eth0", "UP", "mtu 1500", "group 0"},
		},
		{
			name:          "Single link with details, oneline format",
			withAddresses: false,
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
					typeName:       "device",
					specificDevice: &netlink.Dummy{},
				},
			},
			opts:           flags{Oneline: true, Details: true},
			expectedValues: []string{"eth0", "UP", "mtu 1500"},
		},
		{
			name:          "Single link with statistics, oneline format",
			withAddresses: false,
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
				},
			},
			opts:           flags{Oneline: true, Stats: true},
			expectedValues: []string{"eth0", "UP", "mtu 1500", "RX:  bytes  packets errors dropped  missed   mcast", "1000", "100", "10", "1", "TX:  bytes  packets errors dropped carrier collsns", "2000", "200", "20", "2"},
		},
		{
			name:          "Multiple links, oneline format",
			withAddresses: false,
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
				{
					attrs: &netlink.LinkAttrs{
						Name:         "eth1",
						Flags:        net.FlagUp,
						OperState:    netlink.OperUp,
						HardwareAddr: net.HardwareAddr{0x00, 0x1a, 0x2b, 0x3c, 0x4d, 0x5f},
						Index:        2,
						MTU:          1500,
						Group:        0,
						TxQLen:       1000,
					},
					typeName: "device",
				},
			},
			opts:           flags{Oneline: true},
			expectedValues: []string{"eth0", "UP", "mtu 1500", "eth1", "UP", "mtu 1500"},
		},

		// Edge Cases
		{
			name:          "Single link with ValidLft set to max uint32 (forever)",
			withAddresses: true,
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
							Label:       "eth0-label",
							PreferedLft: 0,
							ValidLft:    uint32fix(math.MaxUint32), // fix vishnavanda/netlink. *Lft should be uint32, not int.
						},
					},
				},
			},
			opts:           flags{},
			expectedValues: []string{"eth0", "UP", "mtu 1500", "192.168.1.1/24", "eth0-label", "forever"},
		},
		{
			name:          "Single link with dummy specificDevice",
			withAddresses: false,
			links: []linkData{
				{
					attrs: &netlink.LinkAttrs{
						Name:         "dummy0",
						Flags:        net.FlagUp,
						OperState:    netlink.OperUp,
						HardwareAddr: net.HardwareAddr{0x00, 0x1a, 0x2b, 0x3c, 0x4d, 0x5e},
						Index:        1,
						MTU:          1500,
						Group:        0,
						TxQLen:       1000,
					},
					typeName:       "dummy",
					specificDevice: &netlink.Dummy{},
				},
			},
			opts:           flags{},
			expectedValues: []string{"dummy0", "UP", "mtu 1500"},
		},
		{
			name:          "Single link with master name and numeric group",
			withAddresses: false,
			links: []linkData{
				{
					attrs: &netlink.LinkAttrs{
						Name:         "eth0",
						Flags:        net.FlagUp,
						OperState:    netlink.OperUp,
						HardwareAddr: net.HardwareAddr{0x00, 0x1a, 0x2b, 0x3c, 0x4d, 0x5e},
						Index:        1,
						MTU:          1500,
						Group:        1,
						TxQLen:       1000,
					},
					typeName:   "device",
					masterName: "bridge0",
				},
			},
			opts:           flags{Numeric: true},
			expectedValues: []string{"eth0", "UP", "mtu 1500", "master bridge0", "group 1"},
		},
		{
			name:          "Single link without addresses, brief",
			withAddresses: false,
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
			opts:           flags{Brief: true},
			expectedValues: []string{"eth0", "UP", "00:1a:2b:3c:4d:5e"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			cmd := cmd{
				Out:  &out,
				Opts: tt.opts,
			}

			err := cmd.printLinks(tt.withAddresses, tt.links)
			if err != nil {
				t.Fatalf("printLinks() error = %v", err)
			}

			output := out.String()
			for _, expectedValue := range tt.expectedValues {
				if !strings.Contains(output, expectedValue) {
					t.Errorf("printLinks() missing expected value: %v", expectedValue)
				}
			}
		})
	}
}

func TestDeviceDetailsLine(t *testing.T) {
	uint32Ptr := func(v uint32) *uint32 { return &v }
	boolPtr := func(v bool) *bool { return &v }

	tests := []struct {
		name     string
		device   any
		expected string
	}{
		{
			name: "Bridge device",
			device: &netlink.Bridge{
				HelloTime:     uint32Ptr(2),
				AgeingTime:    uint32Ptr(300),
				VlanFiltering: boolPtr(true),
				LinkAttrs: netlink.LinkAttrs{
					NumTxQueues: 1,
					NumRxQueues: 1,
					GSOMaxSize:  65536,
					GSOMaxSegs:  64,
				},
			},
			expected: "    bridge hello_time 2 ageing_time 300 vlan_filtering 1 numtxqueues 1 numrxqueues 1 gso_max_size 65536 gso_max_segs 64",
		},
		{
			name: "Vlan device",
			device: &netlink.Vlan{
				VlanProtocol: netlink.VLAN_PROTOCOL_8021Q,
				VlanId:       100,
				LinkAttrs: netlink.LinkAttrs{
					NumTxQueues: 1,
					NumRxQueues: 1,
					GSOMaxSize:  65536,
					GSOMaxSegs:  64,
				},
			},
			expected: "    vlan 802.1q vlan-id 100 numtxqueues 1 numrxqueues 1 gso_max_size 65536 gso_max_segs 64",
		},
		{
			name: "Macvlan device",
			device: &netlink.Macvlan{
				Mode: netlink.MACVLAN_MODE_PRIVATE,
				LinkAttrs: netlink.LinkAttrs{
					NumTxQueues: 1,
					NumRxQueues: 1,
					GSOMaxSize:  65536,
					GSOMaxSegs:  64,
				},
			},
			expected: "    macvlan mode 1 numtxqueues 1 numrxqueues 1 gso_max_size 65536 gso_max_segs 64",
		},
		{
			name: "Macvtap device",
			device: &netlink.Macvtap{
				Macvlan: netlink.Macvlan{
					Mode: netlink.MACVLAN_MODE_PRIVATE,
					LinkAttrs: netlink.LinkAttrs{
						NumTxQueues: 1,
						NumRxQueues: 1,
						GSOMaxSize:  65536,
						GSOMaxSegs:  64,
					},
				},
			},
			expected: "    macvtap mode 1 numtxqueues 1 numrxqueues 1 gso_max_size 65536 gso_max_segs 64",
		},
		{
			name: "Tuntap device",
			device: &netlink.Tuntap{
				Mode:  netlink.TUNTAP_MODE_TUN,
				Owner: 1000,
				Group: 1000,
				LinkAttrs: netlink.LinkAttrs{
					NumTxQueues: 1,
					NumRxQueues: 1,
					GSOMaxSize:  65536,
					GSOMaxSegs:  64,
				},
			},
			expected: "    tuntap mode tun owner 1000 group 1000 numtxqueues 1 numrxqueues 1 gso_max_size 65536 gso_max_segs 64",
		},
		{
			name: "Veth device",
			device: &netlink.Veth{
				PeerName:         "veth-peer",
				PeerHardwareAddr: net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
				LinkAttrs: netlink.LinkAttrs{
					NumTxQueues: 1,
					NumRxQueues: 1,
					GSOMaxSize:  65536,
					GSOMaxSegs:  64,
				},
			},
			expected: "    peer veth-peer peer-address 00:11:22:33:44:55 numtxqueues 1 numrxqueues 1 gso_max_size 65536 gso_max_segs 64",
		},
		{
			name: "Vxlan device",
			device: &netlink.Vxlan{
				VxlanId:  42,
				SrcAddr:  net.ParseIP("192.168.1.1"),
				Group:    net.ParseIP("239.0.0.1"),
				TTL:      64,
				TOS:      0,
				Learning: true,
				Proxy:    false,
				RSC:      false,
				Age:      300,
				Limit:    100,
				Port:     4789,
				PortLow:  4789,
				PortHigh: 4790,
				LinkAttrs: netlink.LinkAttrs{
					NumTxQueues: 1,
					NumRxQueues: 1,
					GSOMaxSize:  65536,
					GSOMaxSegs:  64,
				},
			},
			expected: "    vxlan id 42 src 192.168.1.1 group 239.0.0.1 ttl 64 tos 0 learning true proxy false rsc false age 300 limit 100 port 4789 port-low 4789 port-high 4790 numtxqueues 1 numrxqueues 1 gso_max_size 65536 gso_max_segs 64",
		},
		{
			name: "IPVlan device",
			device: &netlink.IPVlan{
				Mode: netlink.IPVLAN_MODE_L2,
				LinkAttrs: netlink.LinkAttrs{
					Flags:       0,
					NumTxQueues: 1,
					NumRxQueues: 1,
					GSOMaxSize:  65536,
					GSOMaxSegs:  64,
				},
			},
			expected: "    ipvlan mode 0 flags 0 numtxqueues 1 numrxqueues 1 gso_max_size 65536 gso_max_segs 64",
		},
		{
			name: "IPVtap device",
			device: &netlink.IPVtap{
				IPVlan: netlink.IPVlan{
					Mode: netlink.IPVLAN_MODE_L2,
					LinkAttrs: netlink.LinkAttrs{
						Flags:       0,
						NumTxQueues: 1,
						NumRxQueues: 1,
						GSOMaxSize:  65536,
						GSOMaxSegs:  64,
					},
				},
			},
			expected: "    ipvtap mode 0 flags 0 numtxqueues 1 numrxqueues 1 gso_max_size 65536 gso_max_segs 64",
		},
		{
			name: "Bond device",
			device: &netlink.Bond{
				Mode:            netlink.BOND_MODE_ACTIVE_BACKUP,
				ActiveSlave:     1,
				Miimon:          100,
				UpDelay:         200,
				DownDelay:       200,
				UseCarrier:      1,
				ArpInterval:     1000,
				ArpValidate:     netlink.BOND_ARP_VALIDATE_NONE,
				ArpAllTargets:   netlink.BOND_ARP_ALL_TARGETS_ANY,
				Primary:         1,
				PrimaryReselect: netlink.BOND_PRIMARY_RESELECT_ALWAYS,
				FailOverMac:     netlink.BOND_FAIL_OVER_MAC_NONE,
				XmitHashPolicy:  netlink.BOND_XMIT_HASH_POLICY_LAYER2,
				ResendIgmp:      1,
				NumPeerNotif:    1,
				AllSlavesActive: 1,
				MinLinks:        1,
				LpInterval:      1,
				PacketsPerSlave: 1,
				LacpRate:        netlink.BOND_LACP_RATE_SLOW,
				AdSelect:        netlink.BOND_AD_SELECT_STABLE,
				LinkAttrs: netlink.LinkAttrs{
					NumTxQueues: 1,
					NumRxQueues: 1,
					GSOMaxSize:  65536,
					GSOMaxSegs:  64,
				},
			},
			expected: "    bond mode active slave 1 1 miimon 100 updelay 200 downdelay 200 use_carrier 1 arp_interval 1000 arp_validate none arp_all_targets any primary 1 primary_reselect always fail_over_mac none layer2 resend_igmp 1 num_peer_notif 1 all_slaves_active 1 min_links 1 lp_interval 1 packets_per_slave 1 lacp_rate slow ad_select stable numtxqueues 1 numrxqueues 1 gso_max_size 65536 gso_max_segs 64",
		},
		{
			name: "Geneve device",
			device: &netlink.Geneve{
				ID:             1,
				Remote:         net.ParseIP("192.168.1.1"),
				Ttl:            64,
				Tos:            0,
				Dport:          6081,
				UdpCsum:        1,
				UdpZeroCsum6Tx: 1,
				UdpZeroCsum6Rx: 1,
				Link:           1,
				FlowBased:      true,
				LinkAttrs: netlink.LinkAttrs{
					NumTxQueues: 1,
					NumRxQueues: 1,
					GSOMaxSize:  65536,
					GSOMaxSegs:  64,
				},
			},
			expected: "    geneve id 1 remote 192.168.1.1 ttl 64 tos 0 dport 6081 udpcsum 1 udp_zero_csum_6TX 1 udp_zero_csum_6RX 1 link 1 flow_based true numtxqueues 1 numrxqueues 1 gso_max_size 65536 gso_max_segs 64",
		},
		{
			name: "Gretap device",
			device: &netlink.Gretap{
				IKey:       1,
				OKey:       1,
				EncapSport: 6081,
				EncapDport: 6081,
				Local:      net.ParseIP("192.168.1.1"),
				Remote:     net.ParseIP("192.168.1.2"),
				IFlags:     1,
				OFlags:     1,
				PMtuDisc:   1,
				Ttl:        64,
				Tos:        0,
				EncapType:  1,
				EncapFlags: 1,
				Link:       1,
				FlowBased:  true,
				LinkAttrs: netlink.LinkAttrs{
					NumTxQueues: 1,
					NumRxQueues: 1,
					GSOMaxSize:  65536,
					GSOMaxSegs:  64,
				},
			},
			expected: "    gretap i_key 1 o_key 1 encap_src_port 6081 encap_dst_port 6081 local 192.168.1.1 remote 192.168.1.2 iflags 1 oflags 1 pmtudisc 1 ttl 64 tos 0 encap_type 1 encap_flags 1 link 1 flow_based true numtxqueues 1 numrxqueues 1 gso_max_size 65536 gso_max_segs 64",
		},
		{
			name: "Iptun device",
			device: &netlink.Iptun{
				Local:      net.ParseIP("192.168.1.1"),
				Remote:     net.ParseIP("192.168.1.2"),
				EncapType:  1,
				EncapFlags: 1,
				Link:       1,
				FlowBased:  true,
				LinkAttrs: netlink.LinkAttrs{
					NumTxQueues: 1,
					NumRxQueues: 1,
					GSOMaxSize:  65536,
					GSOMaxSegs:  64,
				},
			},
			expected: "    iptun local 192.168.1.1 remote 192.168.1.2 encap_type 1 encap_flags 1 link 1 flow_based true numtxqueues 1 numrxqueues 1 gso_max_size 65536 gso_max_segs 64",
		},
		{
			name: "Ip6tnl device",
			device: &netlink.Ip6tnl{
				Local:      net.ParseIP("2001:db8::1"),
				Remote:     net.ParseIP("2001:db8::2"),
				Ttl:        64,
				Tos:        0,
				Proto:      41,
				FlowInfo:   0,
				EncapLimit: 4,
				EncapType:  1,
				EncapSport: 6081,
				EncapDport: 6081,
				EncapFlags: 1,
				LinkAttrs: netlink.LinkAttrs{
					NumTxQueues: 1,
					NumRxQueues: 1,
					GSOMaxSize:  65536,
					GSOMaxSegs:  64,
				},
			},
			expected: "    ip6tnl local 2001:db8::1 remote 2001:db8::2 ttl 64 tos 0 proto 41 flow_info 0 encap_limit 4 encap_type 1 encap_src_port 6081 encap_dst_port 6081 encap_flags 1 numtxqueues 1 numrxqueues 1 gso_max_size 65536 gso_max_segs 64",
		},
		{
			name: "Sittun device",
			device: &netlink.Sittun{
				Local:      net.ParseIP("192.168.1.1"),
				Remote:     net.ParseIP("192.168.1.2"),
				Ttl:        64,
				Tos:        0,
				Proto:      41,
				EncapLimit: 4,
				EncapType:  1,
				EncapSport: 6081,
				EncapDport: 6081,
				EncapFlags: 1,
				LinkAttrs: netlink.LinkAttrs{
					NumTxQueues: 1,
					NumRxQueues: 1,
					GSOMaxSize:  65536,
					GSOMaxSegs:  64,
				},
			},
			expected: "    sittun local 192.168.1.1 remote 192.168.1.2 ttl 64 tos 0 proto 41 encap_limit 4 encap_type 1 encap_src_port 6081 encap_dst_port 6081 encap_flags 1 numtxqueues 1 numrxqueues 1 gso_max_size 65536 gso_max_segs 64",
		},
		{
			name: "Vti device",
			device: &netlink.Vti{
				Local:  net.ParseIP("192.168.1.1"),
				Remote: net.ParseIP("192.168.1.2"),
				IKey:   1,
				OKey:   1,
				LinkAttrs: netlink.LinkAttrs{
					NumTxQueues: 1,
					NumRxQueues: 1,
					GSOMaxSize:  65536,
					GSOMaxSegs:  64,
				},
			},
			expected: "    vti local 192.168.1.1 remote 192.168.1.2 ikey 1 okey 1 numtxqueues 1 numrxqueues 1 gso_max_size 65536 gso_max_segs 64",
		},
		{
			name: "Gretun device",
			device: &netlink.Gretun{
				Local:      net.ParseIP("192.168.1.1"),
				Remote:     net.ParseIP("192.168.1.2"),
				Ttl:        64,
				Tos:        0,
				PMtuDisc:   1,
				EncapType:  1,
				EncapSport: 6081,
				EncapDport: 6081,
				EncapFlags: 1,
				IKey:       1,
				OKey:       1,
				LinkAttrs: netlink.LinkAttrs{
					NumTxQueues: 1,
					NumRxQueues: 1,
					GSOMaxSize:  65536,
					GSOMaxSegs:  64,
				},
			},
			expected: "    gretun local 192.168.1.1 remote 192.168.1.2 ttl 64 tos 0 ptmudisc 1 encap_type 1 encap_src_port 6081 encap_dst_port 6081 encap_flags 1 ikey 1 okey 1 numtxqueues 1 numrxqueues 1 gso_max_size 65536 gso_max_segs 64",
		},
		{
			name: "Xfrmi device",
			device: &netlink.Xfrmi{
				Ifid: 1,
			},
			expected: "    xfrmi if_id 1",
		},
		{
			name: "Can device",
			device: &netlink.Can{
				State:              1,
				BitRate:            500000,
				SamplePoint:        875,
				TimeQuanta:         125,
				PropagationSegment: 1,
				PhaseSegment1:      2,
				PhaseSegment2:      3,
			},
			expected: "    can state 1 bitrate 500000 sample-point 875 tq 125 prop-seg 1 phase-seg1 2 phase-seg2 3",
		},
		{
			name: "IPoIB device",
			device: &netlink.IPoIB{
				Pkey:   1,
				Mode:   1,
				Umcast: 1,
			},
			expected: "    ipoib pkey 1 mode 1 umcast 1",
		},
		{
			name: "BareUDP device",
			device: &netlink.BareUDP{
				Port:       4789,
				EtherType:  0x0800,
				SrcPortMin: 1024,
				MultiProto: true,
			},
			expected: "    port 4789 ethertype 2048 srcport 1024 min multi_proto true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := linkPrinter{}
			line := p.deviceDetailsLine(tt.device)
			if diff := cmp.Diff(line, tt.expected); diff != "" {
				t.Errorf("deviceDetailsLine() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFilterLinks(t *testing.T) {
	links := []netlink.Link{
		&netlink.Dummy{LinkAttrs: netlink.LinkAttrs{Name: "dummy0"}},
		&netlink.Veth{LinkAttrs: netlink.LinkAttrs{Name: "veth0"}},
		&netlink.Bridge{LinkAttrs: netlink.LinkAttrs{Name: "bridge0"}},
		&netlink.Dummy{LinkAttrs: netlink.LinkAttrs{Name: "dummy1"}},
	}

	tests := []struct {
		name     string
		filters  []linkfilter
		expected []netlink.Link
	}{
		{
			name:     "Filter by type dummy",
			filters:  []linkfilter{linkTypeFilter([]string{"dummy"})},
			expected: []netlink.Link{links[0], links[3]},
		},
		{
			name:     "Filter by name dummy0",
			filters:  []linkfilter{linkNameFilter([]string{"dummy0"})},
			expected: []netlink.Link{links[0]},
		},
		{
			name:     "Filter by type dummy and name dummy0",
			filters:  []linkfilter{linkTypeFilter([]string{"dummy"}), linkNameFilter([]string{"dummy0"})},
			expected: []netlink.Link{links[0]},
		},
		{
			name:     "Filter by name dummy0 and type dummy",
			filters:  []linkfilter{linkNameFilter([]string{"dummy0"}), linkTypeFilter([]string{"dummy"})},
			expected: []netlink.Link{links[0]},
		},
		{
			name:     "Filter by type veth",
			filters:  []linkfilter{linkTypeFilter([]string{"veth"})},
			expected: []netlink.Link{links[1]},
		},
		{
			name:     "Filter by type bridge",
			filters:  []linkfilter{linkTypeFilter([]string{"bridge"})},
			expected: []netlink.Link{links[2]},
		},
		{
			name:     "Filter by type dummy and veth",
			filters:  []linkfilter{linkTypeFilter([]string{"dummy", "veth"})},
			expected: []netlink.Link{links[0], links[1], links[3]},
		},
		{
			name:     "Filter by type dummy and bridge",
			filters:  []linkfilter{linkTypeFilter([]string{"dummy", "bridge"})},
			expected: []netlink.Link{links[0], links[2], links[3]},
		},
		{
			name:     "No filters (nil)",
			filters:  nil,
			expected: links,
		},
		{
			name:     "No filters (empty)",
			filters:  []linkfilter{},
			expected: links,
		},
		{
			name:     "Filter by type (nil input)",
			filters:  []linkfilter{linkTypeFilter(nil)},
			expected: links,
		},
		{
			name:     "Filter by type (empty input)",
			filters:  []linkfilter{linkTypeFilter([]string{})},
			expected: links,
		},
		{
			name:     "Filter by name (nil input)",
			filters:  []linkfilter{linkNameFilter(nil)},
			expected: links,
		},
		{
			name:     "Filter by name (empty input)",
			filters:  []linkfilter{linkNameFilter([]string{})},
			expected: links,
		},
		{
			name: "Custom filter: filter links with '0' in name",
			filters: []linkfilter{
				func(link netlink.Link) bool {
					return strings.Contains(link.Attrs().Name, "0")
				},
			},
			expected: []netlink.Link{links[0], links[1], links[2]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filteredLinks := filterLinks(links, tt.filters)
			if diff := cmp.Diff(filteredLinks, tt.expected); diff != "" {
				t.Errorf("filterLinks() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGetLinkDevices(t *testing.T) {
	cmd := &cmd{
		Family: netlink.FAMILY_ALL,
	}

	// getLinkDevices() has a dependency on netlink.LinkList() which is not mocked.
	// Therefore, only test basic error cases here. It is not this package's responsibility
	// to test the correctness of the netlink package.
	// Further testing can be done in the integration tests / VM tests.

	tests := []struct {
		name           string
		withAddresses  bool
		expectedErrNil bool
	}{
		{
			name:           "Get all link devices without addresses",
			withAddresses:  false,
			expectedErrNil: true,
		},
		{
			name:           "Get all link devices with addresses",
			withAddresses:  true,
			expectedErrNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := cmd.getLinkDevices(tt.withAddresses)
			if (err == nil) != tt.expectedErrNil {
				t.Errorf("getLinkDevices() error = %v, expectedErrNil %v", err, tt.expectedErrNil)
			}
		})
	}
}
