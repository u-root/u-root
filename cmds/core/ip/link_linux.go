// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"net"

	"github.com/vishvananda/netlink"
)

func linkSet() error {
	iface, err := parseDeviceName()
	if err != nil {
		return err
	}

	cursor++
	whatIWant = []string{"address", "up", "down", "master"}
	switch findPrefix(arg[cursor], whatIWant) {
	case "address":
		return setHardwareAddress(iface)
	case "up":
		if err := netlink.LinkSetUp(iface); err != nil {
			return fmt.Errorf("%v can't make it up: %v", iface.Attrs().Name, err)
		}
	case "down":
		if err := netlink.LinkSetDown(iface); err != nil {
			return fmt.Errorf("%v can't make it down: %v", iface.Attrs().Name, err)
		}
	case "master":
		cursor++
		whatIWant = []string{"device name"}
		master, err := netlink.LinkByName(arg[cursor])
		if err != nil {
			return err
		}
		return netlink.LinkSetMaster(iface, master)
	default:
		return usage()
	}
	return nil
}

func setHardwareAddress(iface netlink.Link) error {
	cursor++
	hwAddr, err := net.ParseMAC(arg[cursor])
	if err != nil {
		return fmt.Errorf("%v cant parse mac addr %v: %v", iface.Attrs().Name, hwAddr, err)
	}
	err = netlink.LinkSetHardwareAddr(iface, hwAddr)
	if err != nil {
		return fmt.Errorf("%v cant set mac addr %v: %v", iface.Attrs().Name, hwAddr, err)
	}
	return nil
}

func linkAdd() error {
	name, err := parseName()
	if err != nil {
		return err
	}
	attrs := netlink.LinkAttrs{Name: name}

	cursor++
	whatIWant = []string{"type"}
	if arg[cursor] != "type" {
		return usage()
	}

	cursor++
	whatIWant = []string{"bridge"}
	if arg[cursor] != "bridge" {
		return usage()
	}
	return netlink.LinkAdd(&netlink.Bridge{LinkAttrs: attrs})
}

func linkShow(w io.Writer) error {
	cursor++
	whatIWant = []string{"<nothing>", "<device name>"}
	if len(arg[cursor:]) == 0 {
		return showLinks(w, false)
	}
	return nil
}

func link(w io.Writer) error {
	if len(arg) == 1 {
		return linkShow(w)
	}

	cursor++
	whatIWant = []string{"show", "set", "add"}
	cmd := arg[cursor]

	switch findPrefix(cmd, whatIWant) {
	case "show":
		return linkShow(w)
	case "set":
		return linkSet()
	case "add":
		return linkAdd()
	}
	return usage()
}
