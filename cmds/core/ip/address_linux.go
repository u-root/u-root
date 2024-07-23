// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"

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

func (cmd cmd) address() error {
	if !cmd.tokenRemains() {
		return cmd.showAllLinks(true)
	}

	c := cmd.findPrefix("add", "replace", "del", "show", "flush", "help")
	switch c {
	case "show":
		return cmd.addressShow()
	case "add":
		tokenAddr := cmd.nextToken("CIDR format address")
		addr, err := netlink.ParseAddr(tokenAddr)
		if err != nil {
			return err
		}

		iface, err := cmd.parseDeviceName(true)
		if err != nil {
			return err
		}

		if err := cmd.handle.AddrAdd(iface, addr); err != nil {
			return fmt.Errorf("adding %v to %v failed: %v", tokenAddr, cmd.currentToken(), err)
		}

		return nil
	case "replace":
		tokenAddr := cmd.nextToken("CIDR format address")
		addr, err := netlink.ParseAddr(tokenAddr)
		if err != nil {
			return err
		}

		iface, err := cmd.parseDeviceName(true)
		if err != nil {
			return err
		}

		if err := cmd.handle.AddrReplace(iface, addr); err != nil {
			return fmt.Errorf("replacing %v on %v failed: %v", tokenAddr, cmd.currentToken(), err)
		}

		return nil
	case "del":
		tokenAddr := cmd.nextToken("CIDR format address")
		addr, err := netlink.ParseAddr(tokenAddr)
		if err != nil {
			return err
		}

		iface, err := cmd.parseDeviceName(true)
		if err != nil {
			return err
		}

		if err := cmd.handle.AddrDel(iface, addr); err != nil {
			return fmt.Errorf("deleting %v from %v failed: %v", tokenAddr, cmd.currentToken(), err)
		}

		return nil
	case "flush":
		return cmd.addressFlush()
	case "help":
		fmt.Fprint(cmd.out, addressHelp)
		return nil
	default:
		return cmd.usage()
	}
}

func (cmd cmd) addressShow() error {
	device, err := cmd.parseDeviceName(false)
	if errors.Is(err, ErrNotFound) {
		return cmd.showAllLinks(true)
	}
	typeName, err := cmd.parseType()
	if errors.Is(err, ErrNotFound) {
		return cmd.showLink(device, true)
	}

	return cmd.showLink(device, true, typeName)
}

func (cmd cmd) addressFlush() error {
	iface, err := cmd.parseDeviceName(true)
	if err != nil {
		return err
	}
	addr, err := cmd.handle.AddrList(iface, netlink.FAMILY_ALL)
	if err != nil {
		return err
	}

	for _, a := range addr {
		for idx := 1; idx <= cmd.opts.loops; idx++ {
			if err := cmd.handle.AddrDel(iface, &a); err != nil {
				if idx != cmd.opts.loops {
					continue
				}

				return fmt.Errorf("deleting %v from %v failed: %v", a, iface, err)
			}

			break
		}
	}

	return nil
}
