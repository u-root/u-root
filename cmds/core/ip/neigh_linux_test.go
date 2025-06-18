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
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/vishvananda/netlink"
)

func TestParseNeighAddDelReplaceParam(t *testing.T) {
	tests := []struct {
		name      string
		cmd       cmd
		wantNeigh netlink.Neigh
		wantErr   bool
	}{
		{
			name: "all opts",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "neigh", "add", "address", "127.0.0.2", "lladdr", "00:00:00:00:00:01", "dev", "lo", "nud", "1", "router", "extern_learn"},
				Out:    new(bytes.Buffer),
			},
			wantNeigh: netlink.Neigh{
				LinkIndex:    1,
				State:        netlink.NUD_INCOMPLETE,
				Family:       netlink.FAMILY_V4,
				Flags:        netlink.NTF_ROUTER | netlink.NTF_EXT_LEARNED,
				HardwareAddr: net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
				IP:           net.ParseIP("127.0.0.2"),
			},
		},
		{
			name: "wrong addr",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "neigh", "add", "bcx"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "wrong dev",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "neigh", "add", "address", "127.0.0.2", "dev", "byzxa"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "wrong hwAddr",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "neigh", "add", "address", "127.0.0.2", "lladdr", "b"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "wrong hwAddr",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "neigh", "add", "address", "127.0.0.2", "nud", "a"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "invalid option",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "neigh", "add", "address", "127.0.0.2", "x"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "device not specified",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "neigh", "add", "address", "127.0.0.2"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "all opts ipv6",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "neigh", "add", "address", "::ff", "lladdr", "00:00:00:00:00:01", "dev", "lo", "nud", "1", "router", "extern_learn"},
				Out:    new(bytes.Buffer),
			},
			wantNeigh: netlink.Neigh{
				LinkIndex:    1,
				State:        netlink.NUD_INCOMPLETE,
				Family:       netlink.FAMILY_V6,
				Flags:        netlink.NTF_ROUTER | netlink.NTF_EXT_LEARNED,
				HardwareAddr: net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
				IP:           net.ParseIP("::ff"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			neigh, err := tt.cmd.parseNeighAddDelReplaceParams()
			if (err != nil) != tt.wantErr {
				t.Errorf("vrf() = %v, want %t", err, tt.wantErr)
			}

			if !tt.wantErr {
				diff := cmp.Diff(*neigh, tt.wantNeigh)
				if diff != "" {
					t.Errorf("unexpected result (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestParseNeighShowFlush(t *testing.T) {
	validIP, validSubNet, err := net.ParseCIDR("192.168.187.0/8")
	if err != nil {
		t.Fatalf("failed to parse CIDR: %v", err)
	}

	tests := []struct {
		name         string
		cmd          cmd
		wantAddr     net.IP
		wantSubNet   *net.IPNet
		wantLinkName string
		wantProxy    bool
		wantNud      int
		wantErr      bool
	}{
		{
			name: "all opts",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "neigh", "show", "to", "192.6.6.6", "dev", "lo", "nud", "none", "proxy"},
				Out:    new(bytes.Buffer),
			},
			wantAddr:     net.ParseIP("192.6.6.6"),
			wantLinkName: "lo",
			wantProxy:    true,
			wantNud:      netlink.NUD_NONE,
		},
		{
			name: "subnet",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "neigh", "show", "to", "192.168.187.0/8", "dev", "lo", "nud", "none", "proxy"},
				Out:    new(bytes.Buffer),
			},
			wantAddr:     validIP,
			wantSubNet:   validSubNet,
			wantLinkName: "lo",
			wantProxy:    true,
			wantNud:      netlink.NUD_NONE,
		},
		{
			name: "subnet with IPv6",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "neigh", "show", "to", "::/64", "dev", "lo", "nud", "none", "proxy"},
				Out:    new(bytes.Buffer),
			},
			wantAddr:     net.ParseIP("::"),
			wantSubNet:   &net.IPNet{IP: net.ParseIP("::"), Mask: net.CIDRMask(64, 128)},
			wantLinkName: "lo",
			wantProxy:    true,
			wantNud:      netlink.NUD_NONE,
		},
		{
			name: "invalid ip",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "neigh", "show", "to", "abc"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "invalid nud",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "neigh", "show", "nud", "abc"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "invalid dev",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "neigh", "show", "dev", "xyz"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "invalid opt",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "neigh", "show", "nid"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, subNet, iface, proxy, nud, err := tt.cmd.parseNeighShowFlush()
			if (err != nil) != tt.wantErr {
				t.Errorf("vrf() = %v, want %t", err, tt.wantErr)
			}

			if !tt.wantErr {
				if addr != nil && addr.String() != tt.wantAddr.String() {
					t.Errorf("unexpected result (-want +got):\n%s", cmp.Diff(addr, tt.wantAddr))
				}

				if subNet != nil && subNet.String() != tt.wantSubNet.String() {
					t.Errorf("unexpected result (-want +got):\n%s", cmp.Diff(subNet, tt.wantSubNet))
				}

				if iface.Attrs().Name != tt.wantLinkName {
					t.Errorf("unexpected result (-want +got):\n%s", cmp.Diff(iface, tt.wantLinkName))
				}

				if proxy != tt.wantProxy {
					t.Errorf("unexpected result (-want +got):\n%s", cmp.Diff(proxy, tt.wantProxy))
				}

				if nud != tt.wantNud {
					t.Errorf("unexpected result (-want +got):\n%s", cmp.Diff(nud, tt.wantNud))
				}
			}
		})
	}
}

func TestGetState(t *testing.T) {
	tests := []struct {
		state    int
		expected string
	}{
		{0x01, "INCOMPLETE"},
		{0x02, "REACHABLE"},
		{0x04, "STALE"},
		{0x08, "DELAY"},
		{0x10, "PROBE"},
		{0x20, "FAILED"},
		{0x40, "NOARP"},
		{0x80, "PERMANENT"},
		{0x00, "NONE"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := getState(tt.state)
			if result != tt.expected {
				t.Errorf("getState(%d) = %s; want %s", tt.state, result, tt.expected)
			}
		})
	}
}

func TestFilterNeighsByAddr(t *testing.T) {
	_, validSubNet, err := net.ParseCIDR("192.168.1.0/24")
	if err != nil {
		t.Fatalf("failed to parse CIDR: %v", err)
	}

	tests := []struct {
		name              string
		neighs            []netlink.Neigh
		address           net.IP
		subNet            *net.IPNet
		expected          []netlink.Neigh
		linkNames         []string
		expectedLinkNames []string
	}{
		{
			name: "Filter by specific IP",
			neighs: []netlink.Neigh{
				{IP: net.ParseIP("192.168.1.1")},
				{IP: net.ParseIP("192.168.1.2")},
			},
			address:           net.ParseIP("192.168.1.1"),
			expected:          []netlink.Neigh{{IP: net.ParseIP("192.168.1.1")}},
			linkNames:         []string{"eth0", "eth1"},
			expectedLinkNames: []string{"eth0"},
		},
		{
			name: "Filter out NUD_NOARP state",
			neighs: []netlink.Neigh{
				{IP: net.ParseIP("192.168.1.1"), State: netlink.NUD_NOARP},
				{IP: net.ParseIP("192.168.1.2")},
			},
			address:           nil,
			expected:          []netlink.Neigh{{IP: net.ParseIP("192.168.1.2")}},
			linkNames:         []string{"eth0", "eth1"},
			expectedLinkNames: []string{"eth1"},
		},
		{
			name: "No address filter",
			neighs: []netlink.Neigh{
				{IP: net.ParseIP("192.168.1.1")},
				{IP: net.ParseIP("192.168.1.2")},
			},
			address: nil,
			expected: []netlink.Neigh{
				{IP: net.ParseIP("192.168.1.1")},
				{IP: net.ParseIP("192.168.1.2")},
			},
			linkNames:         []string{"eth0", "eth1"},
			expectedLinkNames: []string{"eth0", "eth1"},
		},
		{
			name:              "Empty neighbors list",
			neighs:            []netlink.Neigh{},
			address:           nil,
			expected:          []netlink.Neigh{},
			linkNames:         []string{},
			expectedLinkNames: []string{},
		},
		{
			name: "Filter by subnet /24",
			neighs: []netlink.Neigh{
				{IP: net.ParseIP("192.168.1.1")},
				{IP: net.ParseIP("192.168.1.2")},
				{IP: net.ParseIP("10.0.0.1")},
			},
			subNet:            validSubNet,
			expected:          []netlink.Neigh{{IP: net.ParseIP("192.168.1.1")}, {IP: net.ParseIP("192.168.1.2")}},
			linkNames:         []string{"eth0", "eth1", "eth2"},
			expectedLinkNames: []string{"eth0", "eth1"},
		},
		{
			name: "Filter by subnet /16",
			neighs: []netlink.Neigh{
				{IP: net.ParseIP("192.168.1.1")},
				{IP: net.ParseIP("192.168.2.1")},
				{IP: net.ParseIP("10.0.0.1")},
			},
			subNet:            mustParseCIDR("192.168.0.0/16"),
			expected:          []netlink.Neigh{{IP: net.ParseIP("192.168.1.1")}, {IP: net.ParseIP("192.168.2.1")}},
			linkNames:         []string{"eth0", "eth1", "eth2"},
			expectedLinkNames: []string{"eth0", "eth1"},
		},
		{
			name: "Filter by IPv6 subnet",
			neighs: []netlink.Neigh{
				{IP: net.ParseIP("2001:db8::1")},
				{IP: net.ParseIP("2001:db8::2")},
				{IP: net.ParseIP("fe80::1")},
			},
			subNet:            mustParseCIDR("2001:db8::/64"),
			expected:          []netlink.Neigh{{IP: net.ParseIP("2001:db8::1")}, {IP: net.ParseIP("2001:db8::2")}},
			linkNames:         []string{"eth0", "eth1", "eth2"},
			expectedLinkNames: []string{"eth0", "eth1"},
		},
		{
			name: "Filter by subnet with no matches",
			neighs: []netlink.Neigh{
				{IP: net.ParseIP("192.168.1.1")},
				{IP: net.ParseIP("192.168.1.2")},
			},
			subNet:            mustParseCIDR("10.0.0.0/8"),
			expected:          []netlink.Neigh{},
			linkNames:         []string{"eth0", "eth1"},
			expectedLinkNames: []string{},
		},
		{
			name: "Mix of NUD_NOARP and subnet filter",
			neighs: []netlink.Neigh{
				{IP: net.ParseIP("192.168.1.1"), State: netlink.NUD_NOARP},
				{IP: net.ParseIP("192.168.1.2")},
				{IP: net.ParseIP("192.168.1.3")},
			},
			subNet:            mustParseCIDR("192.168.1.0/24"),
			expected:          []netlink.Neigh{{IP: net.ParseIP("192.168.1.2")}, {IP: net.ParseIP("192.168.1.3")}},
			linkNames:         []string{"eth0", "eth1", "eth2"},
			expectedLinkNames: []string{"eth1", "eth2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, linkNames := filterNeighsByAddr(tt.neighs, tt.linkNames, &tt.address, tt.subNet)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Test %s failed: expected neighbors %v, got %v", tt.name, tt.expected, result)
			}
			if !reflect.DeepEqual(linkNames, tt.expectedLinkNames) {
				t.Errorf("Test %s failed: expected link names %v, got %v", tt.name, tt.linkNames, tt.expectedLinkNames)
			}
		})
	}
}

func TestPrintNeighs(t *testing.T) {
	tests := []struct {
		name        string
		neighs      []netlink.Neigh
		ifacesNames []string
		opts        flags
		expected    string
	}{
		{
			name: "Print neighbors in brief format",
			neighs: []netlink.Neigh{
				{IP: net.ParseIP("192.168.1.1"), HardwareAddr: net.HardwareAddr{0x00, 0x0c, 0x29, 0x3e, 0x1e, 0x4c}},
				{IP: net.ParseIP("192.168.1.2"), HardwareAddr: net.HardwareAddr{0x00, 0x0c, 0x29, 0x3e, 0x1e, 0x4d}},
			},
			ifacesNames: []string{"eth0", "eth1"},
			opts:        flags{Brief: true},
			expected:    "192.168.1.1                             eth0          00:0c:29:3e:1e:4c\n192.168.1.2                             eth1          00:0c:29:3e:1e:4d\n",
		},
		{
			name: "Print neighbors in detailed format",
			neighs: []netlink.Neigh{
				{IP: net.ParseIP("192.168.1.1"), HardwareAddr: net.HardwareAddr{0x00, 0x0c, 0x29, 0x3e, 0x1e, 0x4c}, State: netlink.NUD_REACHABLE, Flags: netlink.NTF_ROUTER},
				{IP: net.ParseIP("192.168.1.2"), HardwareAddr: net.HardwareAddr{0x00, 0x0c, 0x29, 0x3e, 0x1e, 0x4d}, State: netlink.NUD_STALE},
			},
			ifacesNames: []string{"eth0", "eth1"},
			opts:        flags{Brief: false},
			expected:    "192.168.1.1 dev eth0 lladdr 00:0c:29:3e:1e:4c router REACHABLE\n192.168.1.2 dev eth1 lladdr 00:0c:29:3e:1e:4d STALE\n",
		},
		{
			name: "Print neighbors in JSON format (brief)",
			neighs: []netlink.Neigh{
				{IP: net.ParseIP("192.168.1.1"), HardwareAddr: net.HardwareAddr{0x00, 0x0c, 0x29, 0x3e, 0x1e, 0x4c}},
				{IP: net.ParseIP("192.168.1.2"), HardwareAddr: net.HardwareAddr{0x00, 0x0c, 0x29, 0x3e, 0x1e, 0x4d}},
			},
			ifacesNames: []string{"eth0", "eth1"},
			opts:        flags{JSON: true, Brief: true},
			expected:    `[{"dst":"192.168.1.1","dev":"eth0","lladdr":"00:0c:29:3e:1e:4c"},{"dst":"192.168.1.2","dev":"eth1","lladdr":"00:0c:29:3e:1e:4d"}]`,
		},
		{
			name: "Print neighbors in JSON format (detailed)",
			neighs: []netlink.Neigh{
				{IP: net.ParseIP("192.168.1.1"), HardwareAddr: net.HardwareAddr{0x00, 0x0c, 0x29, 0x3e, 0x1e, 0x4c}, State: netlink.NUD_REACHABLE},
				{IP: net.ParseIP("192.168.1.2"), HardwareAddr: net.HardwareAddr{0x00, 0x0c, 0x29, 0x3e, 0x1e, 0x4d}, State: netlink.NUD_STALE},
			},
			ifacesNames: []string{"eth0", "eth1"},
			opts:        flags{JSON: true, Brief: false},
			expected:    `[{"dst":"192.168.1.1","dev":"eth0","lladdr":"00:0c:29:3e:1e:4c","state":"REACHABLE"},{"dst":"192.168.1.2","dev":"eth1","lladdr":"00:0c:29:3e:1e:4d","state":"STALE"}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			cmd := cmd{
				Opts: tt.opts,
				Out:  &out,
			}

			err := cmd.printNeighs(tt.neighs, tt.ifacesNames)
			if err != nil {
				t.Fatalf("printNeighs() error = %v", err)
			}

			if got := out.String(); got != tt.expected {
				t.Errorf("printNeighs() diff:\n%v", cmp.Diff(got, tt.expected))
			}
		})
	}
}

func TestNeighFlagState(t *testing.T) {
	tests := []struct {
		name      string
		cmd       cmd
		proxy     bool
		nud       int
		wantFlags uint8
		wantState uint16
		wantErr   bool
	}{
		{
			name:      "Valid family, proxy true, valid nud",
			cmd:       cmd{Family: 2},
			proxy:     true,
			nud:       100,
			wantFlags: netlink.NTF_PROXY,
			wantState: 100,
			wantErr:   false,
		},
		{
			name:      "Valid family, proxy false, valid nud",
			cmd:       cmd{Family: 2},
			proxy:     false,
			nud:       100,
			wantFlags: 0,
			wantState: 100,
			wantErr:   false,
		},
		{
			name:      "Invalid family",
			cmd:       cmd{Family: 300},
			proxy:     false,
			nud:       100,
			wantFlags: 0,
			wantState: 0,
			wantErr:   true,
		},
		{
			name:      "Valid family, proxy true, nud -1",
			cmd:       cmd{Family: 2},
			proxy:     true,
			nud:       -1,
			wantFlags: netlink.NTF_PROXY,
			wantState: 0,
			wantErr:   false,
		},
		{
			name:      "Valid family, proxy false, nud max uint16",
			cmd:       cmd{Family: 2},
			proxy:     false,
			nud:       math.MaxUint16,
			wantFlags: 0,
			wantState: math.MaxUint16,
			wantErr:   false,
		},
		{
			name:      "Valid family, proxy false, nud greater than max uint16",
			cmd:       cmd{Family: 2},
			proxy:     false,
			nud:       math.MaxUint16 + 1,
			wantFlags: 0,
			wantState: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFlags, gotState, err := tt.cmd.neighFlagState(tt.proxy, tt.nud)
			if (err != nil) != tt.wantErr {
				t.Errorf("neighFlagState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotFlags != tt.wantFlags {
				t.Errorf("neighFlagState() gotFlags = %v, want %v", gotFlags, tt.wantFlags)
			}
			if gotState != tt.wantState {
				t.Errorf("neighFlagState() gotState = %v, want %v", gotState, tt.wantState)
			}
		})
	}
}

func TestParseNeighGet(t *testing.T) {
	tests := []struct {
		name    string
		token   []string
		wantIP  net.IP
		wantIF  netlink.Link
		wantErr bool
	}{
		{
			name:    "Valid IP and interface",
			token:   []string{"address", "192.168.0.2", "dev", "lo"},
			wantErr: false,
		},
		{
			name:  "Error in parseAddress",
			token: []string{"abc"},

			wantIP:  nil,
			wantIF:  nil,
			wantErr: true,
		},
		{
			name:    "Error in parseDeviceName",
			token:   []string{"address", "192.168.0.2", "xyz"},
			wantIF:  nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := cmd{
				Cursor: -1,
				Args:   tt.token,
			}
			_, _, err := cmd.parseNeighGet()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseNeighGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestParseNUD(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{
			name:    "valid text state - none",
			input:   "none",
			want:    netlink.NUD_NONE,
			wantErr: false,
		},
		{
			name:    "valid text state - incomplete",
			input:   "incomplete",
			want:    netlink.NUD_INCOMPLETE,
			wantErr: false,
		},
		{
			name:    "valid text state - reachable",
			input:   "reachable",
			want:    netlink.NUD_REACHABLE,
			wantErr: false,
		},
		{
			name:    "valid text state - stale",
			input:   "stale",
			want:    netlink.NUD_STALE,
			wantErr: false,
		},
		{
			name:    "valid text state - delay",
			input:   "delay",
			want:    netlink.NUD_DELAY,
			wantErr: false,
		},
		{
			name:    "valid text state - probe",
			input:   "probe",
			want:    netlink.NUD_PROBE,
			wantErr: false,
		},
		{
			name:    "valid text state - failed",
			input:   "failed",
			want:    netlink.NUD_FAILED,
			wantErr: false,
		},
		{
			name:    "valid text state - noarp",
			input:   "noarp",
			want:    netlink.NUD_NOARP,
			wantErr: false,
		},
		{
			name:    "valid text state - permanent",
			input:   "permanent",
			want:    netlink.NUD_PERMANENT,
			wantErr: false,
		},
		{
			name:    "valid text state - mixed case",
			input:   "PeRmAnEnT",
			want:    netlink.NUD_PERMANENT,
			wantErr: false,
		},
		{
			name:    "valid numeric state - 0",
			input:   "0",
			want:    netlink.NUD_NONE,
			wantErr: false,
		},
		{
			name:    "valid numeric state - 2",
			input:   "2",
			want:    netlink.NUD_REACHABLE,
			wantErr: false,
		},
		{
			name:    "valid numeric state - 128",
			input:   "128",
			want:    netlink.NUD_PERMANENT,
			wantErr: false,
		},
		{
			name:    "invalid text state",
			input:   "unknown",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid numeric state - negative",
			input:   "-1",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid numeric state - not defined",
			input:   "3",
			want:    3,
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseNUD(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseNUD() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseNUD() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortNeighs(t *testing.T) {
	tests := []struct {
		name            string
		input           []netlink.Neigh
		ifacesNames     []string
		wantNeighs      []netlink.Neigh
		wantIfacesNames []string
	}{
		{
			name: "Sort by IP address",
			input: []netlink.Neigh{
				{IP: net.ParseIP("192.168.1.2")},
				{IP: net.ParseIP("192.168.1.1")},
				{IP: net.ParseIP("10.0.0.1")},
			},
			ifacesNames: []string{"eth1", "eth0", "eth2"},
			wantNeighs: []netlink.Neigh{
				{IP: net.ParseIP("10.0.0.1")},
				{IP: net.ParseIP("192.168.1.1")},
				{IP: net.ParseIP("192.168.1.2")},
			},
			wantIfacesNames: []string{"eth2", "eth0", "eth1"},
		},
		{
			name: "Sort with IPv4 and IPv6 addresses",
			input: []netlink.Neigh{
				{IP: net.ParseIP("192.168.1.1")},
				{IP: net.ParseIP("::1")},
				{IP: net.ParseIP("10.0.0.1")},
				{IP: net.ParseIP("2001:db8::1")},
			},
			ifacesNames: []string{"eth1", "eth0", "eth2", "eth3"},
			wantNeighs: []netlink.Neigh{
				{IP: net.ParseIP("10.0.0.1")},
				{IP: net.ParseIP("192.168.1.1")},
				{IP: net.ParseIP("::1")},
				{IP: net.ParseIP("2001:db8::1")},
			},
			wantIfacesNames: []string{"eth2", "eth1", "eth0", "eth3"},
		},
		{
			name:       "Empty list",
			input:      []netlink.Neigh{},
			wantNeighs: []netlink.Neigh{},
		},
		{
			name: "Single element",
			input: []netlink.Neigh{
				{IP: net.ParseIP("192.168.1.1")},
			},
			ifacesNames: []string{"eth0"},
			wantNeighs: []netlink.Neigh{
				{IP: net.ParseIP("192.168.1.1")},
			},
			wantIfacesNames: []string{"eth0"},
		},
		{
			name: "Sort by link index",
			input: []netlink.Neigh{
				{IP: net.ParseIP("192.168.1.1"), LinkIndex: 3},
				{IP: net.ParseIP("192.168.1.2"), LinkIndex: 1},
				{IP: net.ParseIP("192.168.1.3"), LinkIndex: 2},
			},
			ifacesNames: []string{"eth2", "eth0", "eth1"},
			wantNeighs: []netlink.Neigh{
				{IP: net.ParseIP("192.168.1.2"), LinkIndex: 1},
				{IP: net.ParseIP("192.168.1.3"), LinkIndex: 2},
				{IP: net.ParseIP("192.168.1.1"), LinkIndex: 3},
			},
			wantIfacesNames: []string{"eth0", "eth1", "eth2"},
		},
		{
			name: "Same IP different link indexes",
			input: []netlink.Neigh{
				{IP: net.ParseIP("192.168.1.1"), LinkIndex: 3},
				{IP: net.ParseIP("192.168.1.1"), LinkIndex: 1},
				{IP: net.ParseIP("192.168.1.1"), LinkIndex: 2},
			},
			ifacesNames: []string{"eth2", "eth0", "eth1"},
			wantNeighs: []netlink.Neigh{
				{IP: net.ParseIP("192.168.1.1"), LinkIndex: 1},
				{IP: net.ParseIP("192.168.1.1"), LinkIndex: 2},
				{IP: net.ParseIP("192.168.1.1"), LinkIndex: 3},
			},
			wantIfacesNames: []string{"eth0", "eth1", "eth2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			neighs, ifaceNames := sortedNeighs(tt.input, tt.ifacesNames)

			if len(neighs) != len(tt.wantNeighs) {
				t.Errorf("Expected %d elements, got %d", len(tt.wantNeighs), len(neighs))
			}

			if diff := cmp.Diff(neighs, tt.wantNeighs); diff != "" {
				t.Errorf("Unexpected sorted neighbors result (-got +want):\n%s", diff)
			}

			if diff := cmp.Diff(ifaceNames, tt.wantIfacesNames); diff != "" {
				t.Errorf("Unexpected sorted interface names result (-got +want):\n%s", diff)
			}
		})
	}
}
