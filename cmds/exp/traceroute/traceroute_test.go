// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"errors"
	"testing"

	"github.com/u-root/u-root/pkg/traceroute"
)

func TestParseFlags(t *testing.T) {
	for _, tt := range []struct {
		name    string
		cmdline []string
		exp     *traceroute.Flags
		err     error
	}{
		{
			name:    "ModuleICMP4",
			cmdline: []string{"progName", "-4", "-m", "icmp", "www.google.com"},
			exp: &traceroute.Flags{
				Host:   "www.google.com",
				Module: "icmp",
				Proto:  "icmp4",
			},
		},
		{
			name:    "DirectICMP4",
			cmdline: []string{"progName", "-4", "--icmp", "www.google.com"},
			exp: &traceroute.Flags{
				Host:   "www.google.com",
				Module: "icmp",
				Proto:  "icmp4",
			},
		},
		{
			name:    "ModuleTCP4",
			cmdline: []string{"progName", "-4", "-m", "tcp", "www.google.com"},
			exp: &traceroute.Flags{
				Host:   "www.google.com",
				Module: "icmp",
				Proto:  "icmp4",
			},
		},
		{
			name:    "DirectTCP4",
			cmdline: []string{"progName", "-4", "--tcp", "www.google.com"},
			exp: &traceroute.Flags{
				Host:   "www.google.com",
				Module: "icmp",
				Proto:  "icmp4",
			},
		},
		{
			name:    "ModuleUDP4",
			cmdline: []string{"progName", "-4", "-m", "udp", "www.google.com"},
			exp: &traceroute.Flags{
				Host:   "www.google.com",
				Module: "icmp",
				Proto:  "icmp4",
			},
		},
		{
			name:    "DirectUDP4",
			cmdline: []string{"progName", "-4", "--udp", "www.google.com"},
			exp: &traceroute.Flags{
				Host:   "www.google.com",
				Module: "icmp",
				Proto:  "icmp4",
			},
		},
		{
			name:    "ModuleICMP6",
			cmdline: []string{"progName", "-6", "-m", "icmp", "www.google.com"},
			exp: &traceroute.Flags{
				Host:   "www.google.com",
				Module: "icmp",
				Proto:  "icmp6",
			},
		},
		{
			name:    "DirectICMP6",
			cmdline: []string{"progName", "-6", "--icmp", "www.google.com"},
			exp: &traceroute.Flags{
				Host:   "www.google.com",
				Module: "icmp",
				Proto:  "icmp6",
			},
		},
		{
			name:    "ModuleTCP6",
			cmdline: []string{"progName", "-6", "-m", "tcp", "www.google.com"},
			exp: &traceroute.Flags{
				Host:   "www.google.com",
				Module: "icmp",
				Proto:  "icmp6",
			},
		},
		{
			name:    "DirectTCP6",
			cmdline: []string{"progName", "-6", "--tcp", "www.google.com"},
			exp: &traceroute.Flags{
				Host:   "www.google.com",
				Module: "icmp",
				Proto:  "icmp6",
			},
		},
		{
			name:    "ModuleUDP6",
			cmdline: []string{"progName", "-6", "-m", "udp", "www.google.com"},
			exp: &traceroute.Flags{
				Host:   "www.google.com",
				Module: "icmp",
				Proto:  "icmp6",
			},
		},
		{
			name:    "DirectUDP6",
			cmdline: []string{"progName", "-6", "--udp", "www.google.com"},
			exp: &traceroute.Flags{
				Host:   "www.google.com",
				Module: "icmp",
				Proto:  "icmp6",
			},
		},
		{
			name:    "FailInvalidFlags",
			cmdline: []string{"progName", "-6", "--udp", "www.google.com", "random stuff to error out", "somemore"},
			exp: &traceroute.Flags{
				Host:   "www.google.com",
				Module: "icmp",
				Proto:  "icmp6",
			},
			err: errFlags,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			_, err := parseFlags(tt.cmdline)
			if !errors.Is(err, tt.err) {
				t.Error(err)
			}
		})
	}
}
