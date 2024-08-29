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

func TestVrf(t *testing.T) {
	tests := []struct {
		name    string
		cmd     cmd
		wantErr bool
	}{
		{
			name: "Help",
			cmd: cmd{
				Cursor: 1,
				Args:   []string{"ip", "vrf", "help"},
				Out:    new(bytes.Buffer),
			},
		},
		{
			name: "Wrong arguments",
			cmd: cmd{
				Cursor: 1,
				Args:   []string{"ip", "vrf", "xyz"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cmd.vrf()
			if (err != nil) != tt.wantErr {
				t.Errorf("vrf() = %v, want %t", err, tt.wantErr)
			}
		})
	}
}

func TestPrintVrf(t *testing.T) {
	links := []netlink.Link{
		&netlink.Vrf{
			LinkAttrs: netlink.LinkAttrs{
				Name: "testVrf1",
			},
			Table: 100,
		},
		&netlink.Bond{},
		&netlink.Vrf{
			LinkAttrs: netlink.LinkAttrs{
				Name: "testVrf2",
			},
			Table: 200,
		},
	}

	tests := []struct {
		name string
		opts flags
		want string
	}{
		{
			name: "Non-JSON output",
			opts: flags{JSON: false},
			want: "Name              Table\n-----------------------\ntestVrf1          100\ntestVrf2          200\n",
		},
		{
			name: "JSON output",
			opts: flags{JSON: true},
			want: `[{"name":"testVrf1","table":100},{"name":"testVrf2","table":200}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			cmd := cmd{Out: &out, Opts: tt.opts}

			err := cmd.printVrf(links)
			if err != nil {
				t.Errorf("printVrf() error = %v", err)
				return
			}

			got := out.String()
			if got != tt.want {
				t.Errorf("printVrf() got = %v, want %v", got, tt.want)
			}
		})
	}
}
