// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
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

// address is the entry point for 'ip address' subcommand.
func (cmd *cmd) address() error {
	if !cmd.tokenRemains() {
		return cmd.addressShow()
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
			return fmt.Errorf("adding %v to %v failed: %w", addr.IP, cmd.currentToken(), err)
		}

		return nil
	case "replace":
		iface, addr, err := cmd.parseAddrAddReplace()
		if err != nil {
			return err
		}

		if err := cmd.handle.AddrReplace(iface, addr); err != nil {
			return fmt.Errorf("replacing %v on %v failed: %w", addr.IP, cmd.currentToken(), err)
		}

		return nil
	case "del":
		iface, addr, err := cmd.parseAddrAddReplace()
		if err != nil {
			return err
		}

		if err := cmd.handle.AddrDel(iface, addr); err != nil {
			return fmt.Errorf("deleting %v from %v failed: %w", addr.IP, cmd.currentToken(), err)
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

// parseAddrAddReplace returns arguments to 'ip addr add' or 'ip addr replace' from the cmdline.
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

// addressShow performs 'ip address show' command.
func (cmd *cmd) addressShow() error {
	name, types := cmd.parseAddrShow()

	links, err := cmd.getLinkDevices(true, linkNameFilter([]string{name}), linkTypeFilter(types))
	if err != nil {
		return fmt.Errorf("address show: %w", err)
	}

	err = cmd.printLinks(true, links)
	if err != nil {
		return fmt.Errorf("address show: %w", err)
	}

	return nil
}

// parseAddrShow returns arguments to 'ip addr show' from the cmdline.
func (cmd *cmd) parseAddrShow() (linkName string, types []string) {
	return cmd.parseLinkShow() // same remaining cmdline syntax
}

// parseAddrFlush returns arguments to 'ip addr flush' from the cmdline.
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

// addressFlush performs 'ip address flush' command.
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

			fmt.Printf("Deleting %v from %v\n", a, iface.Attrs().Name)

			if err := cmd.handle.AddrDel(iface, &a); err != nil {
				if idx != cmd.Opts.Loops {
					continue
				}

				return fmt.Errorf("deleting %v from %v failed: %w", a, iface, err)
			}

			break
		}
	}

	return nil
}

func skipAddr(addr netlink.Addr, filter netlink.Addr) bool {
	if filter.Scope != 0 && addr.Scope != filter.Scope {
		return true
	}

	if filter.Label != "" && addr.Label != filter.Label {
		return true
	}

	return false
}
