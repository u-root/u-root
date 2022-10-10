// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv6"
)

var tests = []struct {
	iface  []string
	isIPv4 bool
	out    string
	err    string
}{
	{
		iface:  []string{"nosuchanimal"},
		isIPv4: true,
		out:    "",
		err:    "no interfaces match nosuchanimal",
	},
	{
		iface:  []string{"nosuchanimal1", "nosuchanimaltoo"},
		isIPv4: true,
		out:    "",
		err:    "more than one interface specified",
	},
}

func TestDhclient(t *testing.T) {
	var out = bytes.NewBuffer(nil)
	log.SetOutput(out)
	for _, tt := range tests {
		out.Reset()
		opts := opts{timeout: 15, retry: 5, dryRun: false, verbose: false, vverbose: false, ipv4: tt.isIPv4, ipv6: true, v6Port: dhcpv6.DefaultServerPort, v6Server: "ff02::1:2", v4Port: dhcpv4.ServerPort}
		err := run(&opts, tt.iface)

		if err != nil && tt.err == "" {
			t.Errorf("no error expected, got: \n%v", err)
		} else if err == nil && tt.err != "" {
			t.Errorf("error \n%v\nexpected, got nil error", tt.err)
		} else if err != nil && err.Error() != tt.err {
			t.Errorf("error \n%v\nexpected, got: \n%v", tt.err, err)
		}

		if tt.err == "" && !strings.HasSuffix(out.String(), tt.out) {
			t.Errorf("output expected:\n%s\ngot:\n%s", tt.out, out.String())
		}
	}
}
