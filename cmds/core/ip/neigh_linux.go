// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/vishvananda/netlink"
)

func neigh(w io.Writer) error {
	if len(arg) != 1 {
		return errors.New("neigh subcommands not supported yet")
	}
	return showNeighbours(w, true)
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
			return fmt.Errorf("can't list neighbours: %v", err)
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
