// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

var wellKnownPortsMap = map[string]string{
	"20":  "ftp-data",
	"21":  "ftp",
	"22":  "ssh-scp",
	"23":  "telnet",
	"25":  "smtp",
	"53":  "domain",
	"80":  "http",
	"88":  "kerberos",
	"110": "pop3",
	"119": "nntp",
	"123": "ntp",
	"143": "imap",
	"443": "https",
	"465": "smtps",
	"563": "nntps",
	"989": "ftps-data",
	"990": "ftps",
	"993": "imaps",
	"995": "pop3s",
}

func (cmd cmd) wellKnownPorts(port string) string {
	if name, ok := wellKnownPortsMap[port]; ok && !cmd.Opts.numerical {
		return name
	}

	return port
}

func listDevices() error {
	links, err := netlink.LinkList()
	if err != nil {
		return err
	}

	for idx, link := range links {
		fmt.Printf("%d.%s [%s]\n", idx, link.Attrs().Name, link.Attrs().OperState)
	}

	return nil
}
