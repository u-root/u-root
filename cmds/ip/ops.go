// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"math"
	"net"
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

var neighStates = map[int]string{
	netlink.NUD_NONE:       "NONE",
	netlink.NUD_INCOMPLETE: "INCOMPLETE",
	netlink.NUD_REACHABLE:  "REACHABLE",
	netlink.NUD_STALE:      "STALE",
	netlink.NUD_DELAY:      "DELAY",
	netlink.NUD_PROBE:      "PROBE",
	netlink.NUD_FAILED:     "FAILED",
	netlink.NUD_NOARP:      "NOARP",
	netlink.NUD_PERMANENT:  "PERMANENT",
}

func getState(state int) string {
	ret := make([]string, 0)
	for st, name := range neighStates {
		if state&st != 0 {
			ret = append(ret, name)
		}
	}
	if len(ret) == 0 {
		return "UNKNOWN"
	}
	return strings.Join(ret, ",")
}

func showNeighbours(w io.Writer, withAddresses bool) error {
	ifaces, err := net.Interfaces()
	if err != nil {
		return err
	}
	for _, iface := range ifaces {
		neighs, err := netlink.NeighList(iface.Index, 0)
		if err != nil {
			return fmt.Errorf("Can't list neighbours? %v", err)
		}

		for _, v := range neighs {
			if v.State&netlink.NUD_NOARP != 0 {
				continue
			}
			entry := fmt.Sprintf("%s dev %s", v.IP.String(), iface.Name)
			if v.HardwareAddr != nil {
				entry += fmt.Sprintf(" lladdr %s", v.HardwareAddr)
			}
			if v.Flags&netlink.NTF_ROUTER != 0 {
				entry += " router"
			}
			entry += " " + getState(v.State)
			fmt.Println(entry)
		}
	}
	return nil
}
