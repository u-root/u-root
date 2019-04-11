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

	"golang.org/x/sys/unix"
)

func showLinks(w io.Writer, withAddresses bool) error {
	ifaces, err := linkList()
	if err != nil {
		return fmt.Errorf("Can't enumerate interfaces? %v", err)
	}

	for _, l := range ifaces {
		fmt.Fprintf(w, "%d: %s: <%s> mtu %d state %s\n", l.Index, l.Attributes.Name,
			strings.Replace(strings.ToUpper(fmt.Sprintf("%x", l.Flags)), "|", ",", -1),
			l.Attributes.MTU, strings.ToUpper(string(l.Attributes.OperationalState)))

		fmt.Fprintf(w, "    link/%x %s\n", l.Type, l.Attributes.Address)

		if withAddresses {
			showLinkAddresses(w, &net.Interface{Index: int(l.Index)})
		}
	}
	return nil
}

func showLinkAddresses(w io.Writer, link *net.Interface) error {
	addrs, err := addrList(link, unix.AF_UNSPEC)
	if err != nil {
		return fmt.Errorf("Can't enumerate addresses: %v", err)
	}

	for _, addr := range addrs {
		var inet string
		switch len(addr.Attributes.Address) {
		case 4:
			inet = "inet"
		case 16:
			inet = "inet6"
		default:
			return fmt.Errorf("Can't figure out IP protocol version")
		}

		fmt.Fprintf(w, "    %s %s", inet, addr.Attributes.Address)
		if addr.Attributes.Broadcast != nil {
			fmt.Fprintf(w, " brd %s", addr.Attributes.Broadcast)
		}
		fmt.Fprintf(w, " scope %s %s\n", addrScopes[addr.Scope], addr.Attributes.Label)

		var validLft, preferredLft string
		if addr.Attributes.CacheInfo.Prefered == math.MaxUint32 {
			preferredLft = "forever"
		} else {
			preferredLft = fmt.Sprintf("%dsec", addr.Attributes.CacheInfo.Prefered)
		}
		if addr.Attributes.CacheInfo.Valid == math.MaxUint32 {
			validLft = "forever"
		} else {
			validLft = fmt.Sprintf("%dsec", addr.Attributes.CacheInfo.Valid)
		}
		fmt.Fprintf(w, "       valid_lft %s preferred_lft %s\n", validLft, preferredLft)
	}
	return nil
}

func getState(state uint16) string {
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
		neighs, err := neighList(&iface, 0)
		if err != nil {
			return fmt.Errorf("Can't list neighbours? %v", err)
		}

		for _, v := range neighs {
			if v.State&NUD_NOARP != 0 {
				continue
			}
			entry := fmt.Sprintf("%s dev %s", v.Attributes.Address.String(), iface.Name)
			if v.Attributes.LLAddress != nil {
				entry += fmt.Sprintf(" lladdr %s", v.Attributes.LLAddress)
			}
			if v.Flags&unix.NTF_ROUTER != 0 {
				entry += " router"
			}
			entry += " " + getState(v.State)
			fmt.Println(entry)
		}
	}
	return nil
}
