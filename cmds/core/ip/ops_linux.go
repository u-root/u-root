// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/vishvananda/netlink"
)

func showAllLinks(w io.Writer, withAddresses bool, filterByType ...string) error {
	links, err := netlink.LinkList()
	if err != nil {
		return fmt.Errorf("can't enumerate interfaces: %v", err)
	}
	return showLinks(w, withAddresses, links, filterByType...)
}

func showLink(w io.Writer, link netlink.Link, withAddresses bool, filterByType ...string) error {
	return showLinks(w, withAddresses, []netlink.Link{link}, filterByType...)
}

func showLinks(w io.Writer, withAddresses bool, links []netlink.Link, filterByType ...string) error {
	for _, v := range links {
		if withAddresses {

			addrs, err := netlink.AddrList(v, family)
			if err != nil {
				return fmt.Errorf("can't enumerate addresses: %v", err)
			}

			// if there are no addresses and the link is not a vrf (only wihout -4 or -6), skip it
			if len(addrs) == 0 && (v.Type() != "vrf" || family != netlink.FAMILY_ALL) {
				continue
			}
		}

		found := true

		// check if the link type is in the filter list if the filter list is not empty
		if len(filterByType) > 0 {
			found = false
		}

		for _, t := range filterByType {
			if v.Type() == t {
				found = true
			}
		}

		if !found {
			continue
		}

		l := v.Attrs()

		master := ""
		if l.MasterIndex != 0 {
			link, err := netlink.LinkByIndex(l.MasterIndex)
			if err != nil {
				return fmt.Errorf("can't get link with index %d: %v", l.MasterIndex, err)
			}
			master = fmt.Sprintf("master %s ", link.Attrs().Name)
		}

		fmt.Fprintf(w, "%d: %s: <%s> mtu %d %sstate %s\n", l.Index, l.Name,
			strings.Replace(strings.ToUpper(l.Flags.String()), "|", ",", -1),
			l.MTU, master, strings.ToUpper(l.OperState.String()))

		fmt.Fprintf(w, "    link/%s %s\n", l.EncapType, l.HardwareAddr)

		if withAddresses {
			showLinkAddresses(w, v)
		}
	}
	return nil
}

func showLinkAddresses(w io.Writer, link netlink.Link) error {
	addrs, err := netlink.AddrList(link, netlink.FAMILY_ALL)
	if err != nil {
		return fmt.Errorf("can't enumerate addresses: %v", err)
	}

	for _, addr := range addrs {

		var inet string
		switch len(addr.IPNet.IP) {
		case 4:
			inet = "inet"
		case 16:
			inet = "inet6"
		default:
			return fmt.Errorf("can't figure out IP protocol version: IP length is %d", len(addr.IPNet.IP))
		}

		fmt.Fprintf(w, "    %s %s", inet, addr.IP)

		if addr.Broadcast != nil {
			fmt.Fprintf(w, " brd %s", addr.Broadcast)
		}

		fmt.Fprintf(w, " scope %s %s\n", addrScopes[netlink.Scope(addr.Scope)], addr.Label)

		var validLft, preferredLft string
		// TODO: fix vishnavanda/netlink. *Lft should be uint32, not int.
		if uint32(addr.PreferedLft) == math.MaxUint32 {
			preferredLft = "forever"
		} else {
			preferredLft = fmt.Sprintf("%dsec", addr.PreferedLft)
		}

		if uint32(addr.ValidLft) == math.MaxUint32 {
			validLft = "forever"
		} else {
			validLft = fmt.Sprintf("%dsec", addr.ValidLft)
		}

		fmt.Fprintf(w, "       valid_lft %s preferred_lft %s\n", validLft, preferredLft)
	}
	return nil
}
