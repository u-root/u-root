// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/vishvananda/netlink"
)

func TestParseTuntap(t *testing.T) {
	tests := []struct {
		name     string
		cmd      cmd
		dev      string
		wantOpts tuntapOptions
		wantErr  bool
	}{
		{
			name: "default",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "tuntap", "add"},
				Out:    new(bytes.Buffer),
			},
			wantOpts: defaultTuntapOptions,
		},
		{
			name: "tun mode",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "tuntap", "add", "mode", "tun"},
				Out:    new(bytes.Buffer),
			},
			wantOpts: tuntapOptions{
				Mode:  netlink.TUNTAP_MODE_TUN,
				User:  -1,
				Group: -1,
				Flags: netlink.TUNTAP_DEFAULTS,
			},
		},
		{
			name: "tap mode",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "tuntap", "add", "mode", "tap"},
				Out:    new(bytes.Buffer),
			},
			wantOpts: tuntapOptions{
				Mode:  netlink.TUNTAP_MODE_TAP,
				User:  -1,
				Group: -1,
				Flags: netlink.TUNTAP_DEFAULTS,
			},
		},
		{
			name: "invalid mode",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "tuntap", "add", "mode", "xyz"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "invalid user",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "tuntap", "add", "user", "avc"},
				Out:    new(bytes.Buffer),
			},
			wantOpts: tuntapOptions{
				Mode:  netlink.TUNTAP_MODE_TAP,
				User:  1,
				Group: 1,
				Flags: netlink.TUNTAP_DEFAULTS,
			},
			wantErr: true,
		},
		{
			name: "invalid group",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "tuntap", "add", "group", "avc"},
				Out:    new(bytes.Buffer),
			},
			wantOpts: tuntapOptions{
				Mode:  netlink.TUNTAP_MODE_TAP,
				User:  1,
				Group: 1,
				Flags: netlink.TUNTAP_DEFAULTS,
			},
			wantErr: true,
		},
		{
			name: "user and group",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "tuntap", "add", "user", "10", "group", "10"},
				Out:    new(bytes.Buffer),
			},
			wantOpts: tuntapOptions{
				Mode:  netlink.TUNTAP_MODE_TUN,
				User:  10,
				Group: 10,
				Flags: netlink.TUNTAP_DEFAULTS,
			},
		},
		{
			name: "name",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "tuntap", "add", "name", "foo"},
				Out:    new(bytes.Buffer),
			},
			wantOpts: tuntapOptions{
				Mode:  netlink.TUNTAP_MODE_TUN,
				User:  -1,
				Group: -1,
				Name:  "foo",
				Flags: netlink.TUNTAP_DEFAULTS,
			},
		},
		{
			name: "one_queue",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "tuntap", "add", "one_queue"},
				Out:    new(bytes.Buffer),
			},
			wantOpts: tuntapOptions{
				Mode:  netlink.TUNTAP_MODE_TUN,
				User:  -1,
				Group: -1,
				Flags: netlink.TUNTAP_DEFAULTS | netlink.TUNTAP_ONE_QUEUE,
			},
		},
		{
			name: "pi",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "tuntap", "add", "pi"},
				Out:    new(bytes.Buffer),
			},
			wantOpts: tuntapOptions{
				Mode:  netlink.TUNTAP_MODE_TUN,
				User:  -1,
				Group: -1,
				Flags: netlink.TUNTAP_DEFAULTS &^ netlink.TUNTAP_NO_PI,
			},
		},
		{
			name: "vnet_hdr",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "tuntap", "add", "vnet_hdr"},
				Out:    new(bytes.Buffer),
			},
			wantOpts: tuntapOptions{
				Mode:  netlink.TUNTAP_MODE_TUN,
				User:  -1,
				Group: -1,
				Flags: netlink.TUNTAP_DEFAULTS | netlink.TUNTAP_VNET_HDR,
			},
		},
		{
			name: "multi queue",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "tuntap", "add", "multi_queue"},
				Out:    new(bytes.Buffer),
			},
			wantOpts: tuntapOptions{
				Mode:  netlink.TUNTAP_MODE_TUN,
				User:  -1,
				Group: -1,
				Flags: netlink.TUNTAP_MULTI_QUEUE_DEFAULTS,
			},
		},
		{
			name: "invalid arg",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "tuntap", "add", "yxz"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts, err := tt.cmd.parseTunTap()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAddrShow() error = %v, wantErr %t", err, tt.wantErr)
			}

			if !tt.wantErr {
				diff := cmp.Diff(tt.wantOpts, opts)

				if diff != "" {
					t.Errorf("parseAddrShow() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestFilterTunTaps(t *testing.T) {
	// Mock links
	mockTun := &netlink.Tuntap{LinkAttrs: netlink.LinkAttrs{Name: "tun0"}, Mode: netlink.TUNTAP_MODE_TUN}
	mockTap := &netlink.Tuntap{LinkAttrs: netlink.LinkAttrs{Name: "tap0"}, Mode: netlink.TUNTAP_MODE_TAP}
	mocklink := &netlink.GenericLink{LinkAttrs: netlink.LinkAttrs{Name: "eth0"}}

	tests := []struct {
		name       string
		links      []netlink.Link
		options    tuntapOptions
		want       *netlink.Tuntap
		wantErr    bool
		errMessage string
	}{
		{
			name:    "filter by name",
			links:   []netlink.Link{mockTun, mockTap, mocklink},
			options: tuntapOptions{Name: "tun0"},
			want:    mockTun,
			wantErr: false,
		},
		{
			name:    "filter by mode",
			links:   []netlink.Link{mockTun, mockTap},
			options: tuntapOptions{Mode: netlink.TUNTAP_MODE_TAP},
			want:    mockTap,
			wantErr: false,
		},
		{
			name:       "no match found",
			links:      []netlink.Link{mockTun, mockTap},
			options:    tuntapOptions{Name: "nonexistent"},
			want:       nil,
			wantErr:    true,
			errMessage: "found 0 matching tun/tap devices",
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := filterTunTaps(tt.links, tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("filterTunTaps() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMessage {
				t.Errorf("filterTunTaps() error message = %v, wantErrMessage %v", err.Error(), tt.errMessage)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterTunTaps() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrintTunTaps(t *testing.T) {
	// Mock netlink.Tuntap instances
	mockTun := &netlink.Tuntap{LinkAttrs: netlink.LinkAttrs{Name: "tun0"}, Mode: netlink.TUNTAP_MODE_TUN}
	mockTap := &netlink.Tuntap{LinkAttrs: netlink.LinkAttrs{Name: "tap0"}, Mode: netlink.TUNTAP_MODE_TAP}
	mockLink := &netlink.GenericLink{LinkAttrs: netlink.LinkAttrs{Name: "eth0"}}
	mockTapOpts := &netlink.Tuntap{LinkAttrs: netlink.LinkAttrs{Name: "tap1"}, Mode: netlink.TUNTAP_MODE_TAP, Owner: 1, Group: 1, NonPersist: true, Flags: netlink.TUNTAP_ONE_QUEUE | netlink.TUNTAP_VNET_HDR}
	mockTapOpts2 := &netlink.Tuntap{LinkAttrs: netlink.LinkAttrs{Name: "tap1"}, Mode: netlink.TUNTAP_MODE_TAP, Flags: netlink.TUNTAP_MULTI_QUEUE | netlink.TUNTAP_NO_PI}

	tests := []struct {
		name     string
		links    []netlink.Link
		cmd      cmd
		expected string
	}{
		{
			name:  "Print single tun device",
			links: []netlink.Link{mockTun, mockLink},
			cmd: cmd{
				Opts: flags{JSON: false},
			},
			expected: "tun0: tun persist\n",
		},
		{
			name:  "Print single tap device",
			links: []netlink.Link{mockTap},
			cmd: cmd{
				Opts: flags{JSON: false},
			},
			expected: "tap0: tap persist\n",
		},
		{
			name:  "Print single tap device with multiple flags",
			links: []netlink.Link{mockTapOpts},
			cmd: cmd{
				Opts: flags{JSON: false},
			},
			expected: "tap1: tap one_queue vnet_hdr non-persist user 1 group 1\n",
		},
		{
			name:  "Print single tap device with various flags",
			links: []netlink.Link{mockTapOpts2},
			cmd: cmd{
				Opts: flags{JSON: false},
			},
			expected: "tap1: tap multi_queue persist\n",
		},
		{
			name:  "Print single tap device with multiple flags",
			links: []netlink.Link{mockTapOpts},
			cmd: cmd{
				Opts: flags{JSON: true},
			},
			expected: `[{"ifname":"tap1","flags":["tap","one_queue","vnet_hdr","non-persist","user 1","group 1"]}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer

			tt.cmd.Out = &out

			err := tt.cmd.printTunTaps(tt.links)
			if err != nil {
				t.Errorf("printTunTaps() error = %v", err)
				return
			}

			if got := out.String(); got != tt.expected {
				t.Errorf("printTunTaps() got = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTunTapDevice(t *testing.T) {
	tests := []struct {
		name     string
		options  tuntapOptions
		expected *netlink.Tuntap
	}{
		{
			name: "Valid options",
			options: tuntapOptions{
				Name:  "test0",
				Mode:  netlink.TUNTAP_MODE_TUN,
				User:  1000,
				Group: 1000,
				Flags: netlink.TUNTAP_DEFAULTS,
			},
			expected: &netlink.Tuntap{
				LinkAttrs: netlink.LinkAttrs{
					Name: "test0",
				},
				Mode:  netlink.TUNTAP_MODE_TUN,
				Owner: 1000,
				Group: 1000,
				Flags: netlink.TUNTAP_DEFAULTS,
			},
		},
		{
			name: "User out of range",
			options: tuntapOptions{
				Name:  "test1",
				Mode:  netlink.TUNTAP_MODE_TUN,
				User:  -1,
				Group: 1000,
				Flags: netlink.TUNTAP_DEFAULTS,
			},
			expected: &netlink.Tuntap{
				LinkAttrs: netlink.LinkAttrs{
					Name: "test1",
				},
				Mode:  netlink.TUNTAP_MODE_TUN,
				Group: 1000,
				Flags: netlink.TUNTAP_DEFAULTS,
			},
		},
		{
			name: "Group out of range",
			options: tuntapOptions{
				Name:  "test2",
				Mode:  netlink.TUNTAP_MODE_TUN,
				User:  1000,
				Group: -1,
				Flags: netlink.TUNTAP_DEFAULTS,
			},
			expected: &netlink.Tuntap{
				LinkAttrs: netlink.LinkAttrs{
					Name: "test2",
				},
				Mode:  netlink.TUNTAP_MODE_TUN,
				Owner: 1000,
				Flags: netlink.TUNTAP_DEFAULTS,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tunTapDevice(tt.options)
			if result.LinkAttrs.Name != tt.expected.LinkAttrs.Name ||
				result.Mode != tt.expected.Mode ||
				result.Owner != tt.expected.Owner ||
				result.Group != tt.expected.Group ||
				result.Flags != tt.expected.Flags {
				t.Errorf("tunTapDevice() = %v, want %v", result, tt.expected)
			}
		})
	}
}
