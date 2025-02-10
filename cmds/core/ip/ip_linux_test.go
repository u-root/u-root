// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/vishvananda/netlink"
)

func TestParseFlags(t *testing.T) {
	testcases := []struct {
		name    string
		args    []string
		out     io.Writer
		wantCmd cmd
		wantErr bool
	}{
		{
			name: "no args",
			args: []string{"ip"},
			out:  &bytes.Buffer{},
			wantCmd: cmd{
				Opts: flags{
					Loops: 1,
				},
			},
		},
		{
			name: "inet4",
			args: []string{"ip", "-4"},
			out:  &bytes.Buffer{},
			wantCmd: cmd{
				Opts: flags{
					Loops: 1,
					Inet4: true,
				},
				Family: netlink.FAMILY_V4,
			},
		},
		{
			name: "inet6",
			args: []string{"ip", "-6"},
			wantCmd: cmd{
				Opts: flags{
					Loops: 1,
					Inet6: true,
				},
				Family: netlink.FAMILY_V6,
			},
		},
		{
			name:    "mpls",
			args:    []string{"ip", "-M"},
			wantErr: true,
		},
		{
			name:    "bridge",
			args:    []string{"ip", "-B"},
			wantErr: true,
		},
		{
			name: "link",
			args: []string{"ip", "-0"},
			wantCmd: cmd{
				Opts: flags{
					Loops: 1,
					Link:  true,
				},
				Family: netlink.FAMILY_ALL,
			},
		},
		{
			name: "family",
			args: []string{"ip", "--family=inet"},
			wantCmd: cmd{
				Opts: flags{
					Loops:  1,
					Family: "inet",
				},
				Family: netlink.FAMILY_V4,
			},
		},
		{
			name: "family inet6",
			args: []string{"ip", "--family=inet6"},
			wantCmd: cmd{
				Opts: flags{
					Loops:  1,
					Family: "inet6",
				},
				Family: netlink.FAMILY_V6,
			},
		},
		{
			name:    "family err",
			args:    []string{"ip", "--family=abc"},
			wantErr: true,
		},
		{
			name:    "resolve",
			args:    []string{"ip", "-r"},
			wantErr: true,
		},
		{
			name:    "color",
			args:    []string{"ip", "--color=all"},
			wantErr: true,
		},
		{
			name: "rcvBuf",
			args: []string{"ip", "--rcvbuf=100"},
			wantCmd: cmd{
				Opts: flags{
					Loops:  1,
					RcvBuf: "100",
				},
				Family: netlink.FAMILY_ALL,
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := parseFlags(tt.args, tt.out)
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("unexpected error: %v", err)
				}
			} else if tt.wantErr {
				t.Fatalf("expected error, got nil")
			}

			if !tt.wantErr {
				diff := cmp.Diff(cmd, tt.wantCmd, cmpopts.IgnoreFields(cmd, "Args", "Out", "handle"))
				if diff != "" {
					t.Errorf("got diff between cmds:\n%v", diff)
				}
			}
		})
	}
}

func TestRunSubCommand(t *testing.T) {
	tests := []struct {
		name    string
		cmd     cmd
		wantErr bool
	}{
		{
			name: "Addr",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"addr", "help"},
				Out:    new(bytes.Buffer),
			},
		},
		{
			name: "Addr invalid",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"addr", "invalid"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "link",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"link", "help"},
				Out:    new(bytes.Buffer),
			},
		},
		{
			name: "link invalid",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"link", "invalid"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "route",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"route", "help"},
				Out:    new(bytes.Buffer),
			},
		},
		{
			name: "route invalid",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"route", "invalid"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "neigh",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"neigh", "help"},
				Out:    new(bytes.Buffer),
			},
		},
		{
			name: "neigh invalid",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"neigh", "invalid"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "monitor",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"monitor", "help"},
				Out:    new(bytes.Buffer),
			},
		},
		{
			name: "monitor invalid",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"monitor", "abc"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "tunnel",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"tunnel", "help"},
				Out:    new(bytes.Buffer),
			},
		},
		{
			name: "tunnel invalid",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"tunnel", "invalid"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "tuntap",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"tuntap", "help"},
				Out:    new(bytes.Buffer),
			},
		},
		{
			name: "tuntap invalid",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"tuntap", "ac"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "tcpmetrics",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"tcpmetrics", "help"},
				Out:    new(bytes.Buffer),
			},
		},
		{
			name: "tcpmetrics invalid",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"tcpmetrics", "abv"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "VRF",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"vrf", "help"},
				Out:    new(bytes.Buffer),
			},
		},
		{
			name: "VRF invalid",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"vrf", "abc"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "xfrm",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"xfrm", "help"},
				Out:    new(bytes.Buffer),
			},
		},
		{
			name: "xfrm invalid",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"xfrm", "a"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "xfrm monitor",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"xfrm", "monitor", "help"},
				Out:    new(bytes.Buffer),
			},
		},
		{
			name: "xfrm monitor invalid",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"xfrm", "monitor", "a"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "xfrm state",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"xfrm", "state", "help"},
				Out:    new(bytes.Buffer),
			},
		},
		{
			name: "xfrm state invalid",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"xfrm", "state", "s"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "xfrm policy",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"xfrm", "policy", "help"},
				Out:    new(bytes.Buffer),
			},
		},
		{
			name: "xfrm policy invalid",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"xfrm", "policy", "aa"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "Help",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"help"},
				Out:    new(bytes.Buffer),
			},
		},
		{
			name: "Fail",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"yz"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "Addr wrong arg",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"addr", "xyz"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "link wrong arg",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"link", "xyz"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "route wrong arg",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"route", "xyz"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "neigh wrong arg",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"neigh", "xyz"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "monitor wrong arg",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"monitor", "xyz"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "tunnel wrong arg",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"tunnel", "xyz"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "tuntap wrong arg",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"tuntap", "xyz"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "tcpmetrics wrong arg",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"tcpmetrics", "xyz"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "VRF wrong arg",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"vrf", "xyz"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
		{
			name: "xfrm wrong arg",
			cmd: cmd{
				Cursor: 0,
				Args:   []string{"xfrm", "xyz"},
				Out:    new(bytes.Buffer),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cmd.runSubCommand()
			if (err != nil) != tt.wantErr {
				t.Errorf("vrf() = %v, want %t", err, tt.wantErr)
			}
		})
	}
}

func TestBatchCmds(t *testing.T) {
	dir := t.TempDir()

	tests := []struct {
		name        string
		fileContent string
		force       bool
		wantErr     bool
	}{
		{"Valid Commands", "vrf help\n\naddr help", false, false},
		{"Invalid Command", "link ax", false, true},
		{"Invalid Command with Force", "vrf xy\naddr help", true, false},
		{name: "Empty File", fileContent: "", force: false, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp(dir, "test")
			if err != nil {
				t.Fatal(err)
			}

			if _, err := tmpFile.WriteString(tt.fileContent); err != nil {
				t.Fatal(err)
			}
			if err := tmpFile.Close(); err != nil {
				t.Fatal(err)
			}

			// Setup cmd struct
			cmd := cmd{
				Out: new(bytes.Buffer),
				Opts: flags{
					Batch: tmpFile.Name(),
					Force: tt.force,
				},
			}

			err = cmd.batchCmds()

			// Assert expectations
			if (err != nil) != tt.wantErr {
				t.Errorf("batchCmds() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		name     string
		cmd      cmd
		expected string
		wantErr  bool
	}{
		{
			name: "Normal execution",
			cmd: cmd{
				Args:           []string{"arg1", "arg2"},
				Cursor:         1,
				ExpectedValues: []string{"arg1", "arg2"},
			},
			wantErr: true,
		},
		{
			name: "Batch execution",
			cmd: cmd{
				Opts: flags{
					Batch: "testdata/batch.txt",
				},
				Cursor: 0,
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var out bytes.Buffer

			test.cmd.Out = &out

			err := test.cmd.run()

			if test.wantErr && err == nil || !test.wantErr && err != nil {
				t.Errorf("expected %v, got %v", test.wantErr, err)
			}
		})
	}
}
