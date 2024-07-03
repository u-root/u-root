// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/vishvananda/netlink"
)

const addressHelp = `Usage: ip address {add|replace} ADDR dev IFNAME 

       ip address del IFADDR dev IFNAME 

	   ip address flush

       ip address [ show [ dev IFNAME ] [ type TYPE ]

	   ip address help

TYPE := { bareudp | bond | bond_slave | bridge | bridge_slave |
          dummy | erspan | geneve | gre | gretap | ifb |
          ip6erspan | ip6gre | ip6gretap | ip6tnl |
          ipip | ipoib | ipvlan | ipvtap |
          macsec | macvlan | macvtap |
          netdevsim | nlmon | rmnet | sit | team | team_slave |
          vcan | veth | vlan | vrf | vti | vxcan | vxlan | wwan |
          xfrm }
`

func address(w io.Writer) error {
	if len(arg) == 1 {
		return showAllLinks(w, true)
	}
	cursor++
	whatIWant = []string{"add", "replace", "del", "show", "flush", "help"}
	cmd := arg[cursor]

	c := findPrefix(cmd, whatIWant)
	switch c {
	case "show":
		return addressShow(w)
	case "add", "change", "replace", "del":
		return addressAddReplaceDel(c)
	case "flush":
		return addressFlush()
	case "help":
		fmt.Fprint(w, addressHelp)
		return nil
	default:
		return usage()
	}
}

func addressShow(w io.Writer) error {
	device, err := parseDeviceName(false)
	if errors.Is(err, ErrNotFound) {
		return showAllLinks(w, true)
	}
	typeName, err := parseType()
	if errors.Is(err, ErrNotFound) {
		return showLink(w, device, true)
	}

	return showLink(w, device, true, typeName)
}

func addressAddReplaceDel(cmd string) error {
	cursor++
	whatIWant = []string{"CIDR format address"}
	addr, err := netlink.ParseAddr(arg[cursor])
	if err != nil {
		return err
	}

	iface, err := parseDeviceName(true)
	if err != nil {
		return err
	}

	c := findPrefix(cmd, whatIWant)
	switch c {
	case "add":
		if err := netlink.AddrAdd(iface, addr); err != nil {
			return fmt.Errorf("adding %v to %v failed: %v", arg[1], arg[2], err)
		}
	case "replace":
		if err := netlink.AddrReplace(iface, addr); err != nil {
			return fmt.Errorf("replacing %v on %v failed: %v", arg[1], arg[2], err)
		}
	case "del":
		if err := netlink.AddrDel(iface, addr); err != nil {
			return fmt.Errorf("deleting %v from %v failed: %v", arg[1], arg[2], err)
		}
	default:
		return fmt.Errorf("subcommand %s not yet implemented, expected: %v", c, whatIWant)
	}
	return nil
}

func addressFlush() error {
	iface, err := parseDeviceName(true)
	if err != nil {
		return err
	}
	addr, err := netlink.AddrList(iface, netlink.FAMILY_ALL)
	if err != nil {
		return err
	}

	for _, a := range addr {
		if err := netlink.AddrDel(iface, &a); err != nil {
			return err
		}
	}

	return nil
}
