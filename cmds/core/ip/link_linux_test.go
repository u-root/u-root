// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"net"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/vishvananda/netlink"
)

func TestParseLinkShow(t *testing.T) {
	tests := []struct {
		name      string
		cmd       cmd
		wantDev   netlink.Link
		wantTypes []string
		wantErr   bool
	}{
		{
			name: "Successful parsing",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "link", "show", "type", "dummy", "abc", "dev", "lo"},
				Out:    new(bytes.Buffer),
			},
			wantDev:   &netlink.Dummy{LinkAttrs: netlink.LinkAttrs{Name: "lo"}},
			wantTypes: []string{"dummy", "abc"},
		},
		{
			name: "Successful parsing",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "link", "show", "dev", "xyz"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.cmd
			gotDev, gotType, err := cmd.parseLinkShow()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseLinkShow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if tt.wantDev != nil {
					if gotDev.Attrs().Name != tt.wantDev.Attrs().Name {
						t.Errorf("parseLinkShow() gotDev = %v, want %v", gotDev, tt.wantDev)
					}
				}
				if c := cmp.Diff(gotType, tt.wantTypes); c != "" {
					t.Errorf("parseLinkShow() diff:\n%v", c)
				}
			}
		})
	}
}

func TestParseLinkAttrs(t *testing.T) {
	tests := []struct {
		name      string
		cmd       cmd
		wantType  string
		wantAttrs netlink.LinkAttrs
		wantErr   bool
	}{
		{
			name: "Successful parsing",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "link", "add", "name", "eth0", "type", "dummy", "txqueuelen", "1000", "address", "00:0c:29:3e:5c:7f", "mtu", "1500", "index", "1", "numtxqueues", "2", "numrxqueues", "2"},
				Out:    new(bytes.Buffer),
			},
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
			name: "invalid txqlen",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "link", "add", "name", "eth0", "type", "dummy", "txqueuelen", "abc", "address", "00:0c:29:3e:5c:7f", "mtu", "1500", "index", "1", "numtxqueues", "2", "numrxqueues", "2"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "invalid address",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "link", "add", "name", "eth0", "type", "dummy", "txqueuelen", "1000", "address", "00:0c:29:3e:5c:7f:00", "mtu", "1500", "index", "1", "numtxqueues", "2", "numrxqueues", "2"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "invalid mtu",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "link", "add", "name", "eth0", "type", "dummy", "txqueuelen", "1000", "address", "00:0c:29:3e:5c:7f", "mtu", "abc", "index", "1", "numtxqueues", "2", "numrxqueues", "2"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "invalid index",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "link", "add", "name", "eth0", "type", "dummy", "txqueuelen", "1000", "address", "00:0c:29:3e:5c:7f", "mtu", "1500", "index", "abc", "numtxqueues", "2", "numrxqueues", "2"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "invalid numtxqueues",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "link", "add", "name", "eth0", "type", "dummy", "txqueuelen", "1000", "address", "00:0c:29:3e:5c:7f", "mtu", "1500", "index", "1", "numtxqueues", "abc", "numrxqueues", "2"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "invalid numrxqueues",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "link", "add", "name", "eth0", "type", "dummy", "txqueuelen", "1000", "address", "00:0c:29:3e:5c:7f", "mtu", "1500", "index", "1", "numtxqueues", "2", "numrxqueues", "abc"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "invalid arg",
			cmd: cmd{
				Cursor: 2,
				Args:   []string{"ip", "link", "add", "name", "eth0", "abc"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.cmd
			gotType, gotAttrs, err := cmd.parseLinkAttrs()
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
