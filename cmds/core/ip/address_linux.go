// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"

	"github.com/vishvananda/netlink"
)

func address(w io.Writer) error {
	var err error
	var addr *netlink.Addr
	if len(arg) == 1 {
		return showLinks(w, true)
	}
	cursor++
	whatIWant = []string{"add", "del"}
	cmd := arg[cursor]

	c := findPrefix(cmd, whatIWant)
	switch c {
	case "add", "del":
		cursor++
		whatIWant = []string{"CIDR format address"}
		addr, err = netlink.ParseAddr(arg[cursor])
		if err != nil {
			return err
		}
	default:
		return usage()
	}
	iface, err := parseDeviceName()
	if err != nil {
		return err
	}
	switch c {
	case "add":
		if err := netlink.AddrAdd(iface, addr); err != nil {
			return fmt.Errorf("adding %v to %v failed: %v", arg[1], arg[2], err)
		}
	case "del":
		if err := netlink.AddrDel(iface, addr); err != nil {
			return fmt.Errorf("deleting %v from %v failed: %v", arg[1], arg[2], err)
		}
	default:
		return fmt.Errorf("devip: arg[0] changed: can't happen")
	}
	return nil
}
