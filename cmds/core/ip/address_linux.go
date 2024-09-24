// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/vishvananda/netlink"
)

const addressHelp = `Usage: ip address {add|replace} ADDR dev IFNAME [ LIFETIME ]

       ip address del IFADDR dev IFNAME 

       ip address flush dev IFNAME [ scope SCOPE-ID ] [ label LABEL ]

       ip address [ show [ dev IFNAME ] [ type TYPE ]

	   ip address help

SCOPE-ID := [ host | link | global | NUMBER ]
LIFETIME := [ valid_lft LFT ] [ preferred_lft LFT ]
LFT := forever | SECONDS
TYPE := { bareudp | bond | bond_slave | bridge | bridge_slave |
          dummy | erspan | geneve | gre | gretap | ifb |
          ip6erspan | ip6gre | ip6gretap | ip6tnl |
          ipip | ipoib | ipvlan | ipvtap |
          macsec | macvlan | macvtap |
          netdevsim | nlmon | rmnet | sit | team | team_slave |
          vcan | veth | vlan | vrf | vti | vxcan | vxlan | wwan |
          xfrm }
`

var stringScope = map[string]netlink.Scope{
	"global": netlink.SCOPE_UNIVERSE,
	"host":   netlink.SCOPE_HOST,
	"link":   netlink.SCOPE_LINK,
}

func (cmd *cmd) address() error {
	if !cmd.tokenRemains() {
		return cmd.showAllLinks(true)
	}

	c := cmd.findPrefix("add", "replace", "del", "show", "flush", "help")
	switch c {
	case "show":
		return cmd.addressShow()
	case "add":
		iface, addr, err := cmd.parseAddrAddReplace()
		if err != nil {
			return err
		}

		if err := cmd.handle.AddrAdd(iface, addr); err != nil {
			return fmt.Errorf("adding %v to %v failed: %v", addr.IP, cmd.currentToken(), err)
		}

		return nil
	case "replace":
		iface, addr, err := cmd.parseAddrAddReplace()
		if err != nil {
			return err
		}

		if err := cmd.handle.AddrReplace(iface, addr); err != nil {
			return fmt.Errorf("replacing %v on %v failed: %v", addr.IP, cmd.currentToken(), err)
		}

		return nil
	case "del":
		iface, addr, err := cmd.parseAddrAddReplace()
		if err != nil {
			return err
		}

		if err := cmd.handle.AddrDel(iface, addr); err != nil {
			return fmt.Errorf("deleting %v from %v failed: %v", addr.IP, cmd.currentToken(), err)
		}

		return nil
	case "flush":
		return cmd.addressFlush()
	case "help":
		fmt.Fprint(cmd.Out, addressHelp)
		return nil
	default:
		return cmd.usage()
	}
}

func (cmd *cmd) parseAddrAddReplace() (netlink.Link, *netlink.Addr, error) {
	tokenAddr := cmd.nextToken("CIDR format address")
	addr, err := netlink.ParseAddr(tokenAddr)
	if err != nil {
		return nil, nil, err
	}

	iface, err := cmd.parseDeviceName(true)
	if err != nil {
		return nil, nil, err
	}

	for cmd.tokenRemains() {
		switch cmd.nextToken("valid_lft", "preferred_lft") {
		case "valid_lft":
			validLft := cmd.nextToken("LFT")
			if validLft != "forever" {
				validLftInt, err := strconv.ParseInt(validLft, 10, 32)
				if err != nil {
					return nil, nil, fmt.Errorf("invalid valid_lft value: %v", validLft)
				}
				addr.ValidLft = int(validLftInt)
			} else {
				addr.ValidLft = 0
			}
		case "preferred_lft":
			preferredLft := cmd.nextToken("LFT")

			if preferredLft != "forever" {
				preferredLftInt, err := strconv.ParseInt(preferredLft, 10, 32)
				if err != nil {
					return nil, nil, fmt.Errorf("invalid valid_lft value: %v", preferredLft)
				}
				addr.PreferedLft = int(preferredLftInt)
			} else {
				addr.PreferedLft = 0
			}
		}
	}
	return iface, addr, nil
}

func (cmd *cmd) addressShow() error {
	device, typeName, err := cmd.parseAddrShow()
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return cmd.showAllLinks(true)
		}

		return err
	}

	return cmd.showLink(device, true, typeName)
}

func (cmd *cmd) parseAddrShow() (netlink.Link, string, error) {
	device, err := cmd.parseDeviceName(false)
	if err != nil {
		return nil, "", err
	}
	typeName, err := cmd.parseType()
	if err != nil {
		return nil, "", err
	}

	return device, typeName, nil
}

func (cmd *cmd) parseAddrFlush() (netlink.Link, netlink.Addr, error) {
	var addr netlink.Addr

	iface, err := cmd.parseDeviceName(true)
	if err != nil {
		return nil, addr, err
	}

	for cmd.tokenRemains() {
		switch cmd.nextToken("scope", "label") {
		case "scope":
			scope := cmd.nextToken("SCOPE-ID")

			if s, ok := stringScope[scope]; ok {
				addr.Scope = int(s)
			} else {
				scopeInt, err := strconv.ParseInt(scope, 10, 32)
				if err != nil {
					return nil, addr, fmt.Errorf("invalid scope value: %v", scope)
				}
				addr.Scope = int(scopeInt)
			}
		case "label":
			addr.Label = cmd.nextToken("LABEL")
		}
	}

	return iface, addr, nil
}

func (cmd *cmd) addressFlush() error {
	iface, addr, err := cmd.parseAddrFlush()
	if err != nil {
		return err
	}

	addrs, err := cmd.handle.AddrList(iface, cmd.Family)
	if err != nil {
		return err
	}

	for _, a := range addrs {
		if skipAddr(a, addr) {
			continue
		}

		for idx := 1; idx <= cmd.Opts.Loops; idx++ {
			if err := cmd.handle.AddrDel(iface, &a); err != nil {
				if idx != cmd.Opts.Loops {
					continue
				}

				return fmt.Errorf("deleting %v from %v failed: %v", a, iface, err)
			}

			break
		}
	}

	return nil
}

func skipAddr(addr netlink.Addr, filter netlink.Addr) bool {
	if addr.Scope != 0 && addr.Scope != filter.Scope {
		return true
	}

	if addr.Label != "" && addr.Label != filter.Label {
		return true
	}

	return false
}
