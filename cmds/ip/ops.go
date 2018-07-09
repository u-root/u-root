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

func showLinks(w io.Writer, withAddresses bool) error {
	ifaces, err := netlink.LinkList()
	if err != nil {
		return fmt.Errorf("Can't enumerate interfaces? %v", err)
	}

	for _, v := range ifaces {
		l := v.Attrs()

		fmt.Fprintf(w, "%d: %s: <%s> mtu %d state %s\n", l.Index, l.Name,
			strings.Replace(strings.ToUpper(fmt.Sprintf("%s", l.Flags)), "|", ",", -1),
			l.MTU, strings.ToUpper(l.OperState.String()))

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
		return fmt.Errorf("Can't enumerate addresses: %v", err)
	}

	for _, addr := range addrs {

		var inet string
		switch len(addr.IPNet.IP) {
		case 4:
			inet = "inet"
		case 16:
			inet = "inet6"
		default:
			return fmt.Errorf("Can't figure out IP protocol version")
		}

		fmt.Fprintf(w, "    %s %s", inet, addr.Peer)
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
