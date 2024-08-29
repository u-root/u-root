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
		name        string
		args        []string
		expectedCmd *cmd
	}{
		{
			name: "Default values",
			args: []string{"newsshd"},
			expectedCmd: &cmd{
				hostKeyFile: "/etc/ssh_host_rsa_key",
				pubKeyFile:  "key.pub",
				port:        "2222",
			},
		},
		{
			name: "Custom hostkeyfile and pubkeyfile",
			args: []string{"newsshd", "--hostkeyfile", "/custom/host_key", "--pubkeyfile", "custom_key.pub"},
			expectedCmd: &cmd{
				hostKeyFile: "/custom/host_key",
				pubKeyFile:  "custom_key.pub",
				port:        "2222",
			},
		},
		{
			name: "Shorthand flags",
			args: []string{"newsshd", "-h", "/shorthand/host_key", "-k", "shorthand_key.pub", "-p", "2022"},
			expectedCmd: &cmd{
				hostKeyFile: "/shorthand/host_key",
				pubKeyFile:  "shorthand_key.pub",
				port:        "2022",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCmd := command(tt.args)
			if !reflect.DeepEqual(gotCmd, tt.expectedCmd) {
				t.Errorf("%s: command() = %+v, want %+v", tt.name, gotCmd, tt.expectedCmd)
			}
		})
	}
}
