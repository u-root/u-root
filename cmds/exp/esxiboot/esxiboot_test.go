// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"reflect"
	"testing"
)

func TestCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantCmd *cmd
	}{
		{
			name: "Test with config flag",
			args: []string{"esxiboot", "--config", "/path/to/config"},
			wantCmd: &cmd{
				cfg: "/path/to/config",
			},
		},
		{
			name: "Test with cdrom flag",
			args: []string{"esxiboot", "--cdrom", "/dev/cdrom"},
			wantCmd: &cmd{
				cdrom: "/dev/cdrom",
			},
		},
		{
			name: "Test with device flag",
			args: []string{"esxiboot", "--device", "/dev/sda"},
			wantCmd: &cmd{
				diskDev: "/dev/sda",
			},
		},
		{
			name: "Test with append flag",
			args: []string{"esxiboot", "--device", "/dev/sda", "--append", "quiet splash"},
			wantCmd: &cmd{
				diskDev:       "/dev/sda",
				appendCmdline: []string{"quiet splash"},
			},
		},
		{
			name: "Test with dry-run flag",
			args: []string{"esxiboot", "--device", "/dev/sda", "--dry-run"},
			wantCmd: &cmd{
				diskDev: "/dev/sda",
				dryRun:  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCmd := command(tt.args)
			if !reflect.DeepEqual(gotCmd, tt.wantCmd) {
				t.Errorf("got: %+v, want: %+v", gotCmd, tt.wantCmd)
			}
		})
	}
}
