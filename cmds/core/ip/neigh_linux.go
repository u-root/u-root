// Copyright 2024 the u-root Authors. All rights reserved
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

const neighHelp = `Usage: ip neigh { add | del | replace }
                ADDR [ lladdr LLADDR ] [ nud STATE ] 
                [ dev DEV ] [ router ] [ extern_learn ] 

       ip neigh { show | flush } [ proxy ] [ dev DEV ] [ nud STATE ]

       ip neigh get ADDR dev DEV

STATE := { delay | failed | incomplete | noarp | none |
           permanent | probe | reachable | stale }
`

func neigh(w io.Writer) error {
	if len(arg) == 1 {
		return showAllNeighbours(w, -1, false)
	}

	cursor++
	expectedValues = []string{"show", "add", "del", "replace", "flush", "get", "help"}
	cmd := arg[cursor]

	switch c := findPrefix(cmd, expectedValues); c {
	case "add", "del", "replace":
		neigh, err := parseNeighAddDelReplaceParams()
		if err != nil {
			return err
		}

		switch c {
		case "add":
			return netlink.NeighAdd(neigh)
		case "del":
			return netlink.NeighDel(neigh)
		case "replace":
			return netlink.NeighSet(neigh)
		}

	case "show":
		return neighShow(w)
	case "flush":
		return neighFlush()
	case "get":
		cursor++
		expectedValues = []string{"CIDR format address"}
		ipAddr := net.ParseIP(arg[cursor])
		if ipAddr == nil {
			return fmt.Errorf("invalid IP address: %s", arg[cursor])
		}
		iface, err := parseDeviceName(true)
		if err != nil {
			return err
		}

		return showNeighbours(w, -1, false, &ipAddr, iface)
	case "help":
		fmt.Fprint(w, neighHelp)
		return nil
	}

	return usage()
}

func parseNeighAddDelReplaceParams() (*netlink.Neigh, error) {
	cursor++
	expectedValues = []string{"CIDR format address"}
	addr := net.ParseIP(arg[cursor])
	if addr == nil {
		return nil, fmt.Errorf("invalid IP address: %s", arg[cursor])
	}

	var (
		iface       netlink.Link
		llAddr      net.HardwareAddr
		deviceFound bool
		state       int
		flag        int
		err         error
	)

	for {
		if cursor == len(arg)-1 {
			break
		}

		cursor++
		expectedValues = []string{"dev", "lladdr", "nud", "router", "extern_learn"}
		switch arg[cursor] {
		case "dev":
			iface, err = parseDeviceName(true)
			if err != nil {
				return nil, err
			}
			deviceFound = true
		case "lladdr":
			llAddr, err = parseHardwareAddress()
			if err != nil {
				return nil, err
			}
		case "nud":
			state, err = parseInt("STATE")
			if err != nil {
				return nil, err
			}

		case "router":
			flag |= netlink.NTF_ROUTER
		case "extern_learn":
			flag |= netlink.NTF_EXT_LEARNED
		default:
			return nil, fmt.Errorf("unsupported option %q, expected: %v", arg[cursor], expectedValues)
		}
	}

	if !deviceFound {
		return nil, fmt.Errorf("device not specified")
	}

	family := netlink.FAMILY_V4
	if addr.To4() == nil {
		family = netlink.FAMILY_V6
	}

	return &netlink.Neigh{
		LinkIndex:    iface.Attrs().Index,
		Family:       family,
		IP:           addr,
		HardwareAddr: llAddr,
		State:        state,
		Flags:        flag,
	}, nil
}

func parseNeighShowFlush() (iface netlink.Link, proxy bool, nud int, err error) {
	nud = -1

	for {
		if cursor == len(arg)-1 {
			break
		}

		cursor++
		expectedValues = []string{"dev", "proxy", "nud"}
		switch arg[cursor] {
		case "dev":
			dev, err := parseDeviceName(true)
			iface = dev
			if err != nil {
				return nil, false, 0, err
			}
		case "proxy":
			proxy = true
		case "nud":
			nud, err = parseInt("STATE")
			if err != nil {
				return nil, false, 0, err
			}
		default:
			return nil, false, 0, fmt.Errorf("unsupported option %q, expected: %v", arg[cursor], expectedValues)
		}
	}

	return iface, proxy, nud, nil
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

func showAllNeighbours(w io.Writer, nud int, proxy bool) error {
	ifaces, err := netlink.LinkList()
	if err != nil {
		return err
	}

	return showNeighbours(w, nud, proxy, nil, ifaces...)
}

type Neigh struct {
	Dst    net.IP `json:"dst"`
	Dev    string `json:"dev"`
	LLAddr string `json:"lladdr"`
	State  string `json:"state"`
}

func showNeighbours(w io.Writer, nud int, proxy bool, address *net.IP, ifaces ...netlink.Link) error {
	var (
		flags uint8
		state uint16
	)

	if proxy {
		flags |= netlink.NTF_PROXY
	}

	if nud != -1 && nud <= math.MaxUint16 {
		state = uint16(nud)
	}

	filteredNeighs := make([]netlink.Neigh, 0)
	ifaceNames := make([]string, 0)

	for _, iface := range ifaces {
		neighs, err := netlink.NeighListExecute(netlink.Ndmsg{
			Family: netlink.FAMILY_ALL,
			Index:  uint32(iface.Attrs().Index),
			Flags:  flags,
			State:  state,
		})
		if err != nil {
			return err
		}

		for _, v := range neighs {
			if address != nil && !v.IP.Equal(*address) {
				continue
			}

			if v.State&netlink.NUD_NOARP != 0 {
				continue
			}

			filteredNeighs = append(filteredNeighs, v)
			ifaceNames = append(ifaceNames, iface.Attrs().Name)
		}
	}

	if f.json {
		neighs := make([]Neigh, 0, len(filteredNeighs))

		for idx, v := range filteredNeighs {
			neighs = append(neighs, Neigh{
				Dst:    v.IP,
				Dev:    ifaceNames[idx],
				LLAddr: v.HardwareAddr.String(),
				State:  getState(v.State),
			})
		}

		return printJSON(w, neighs)
	}

	for idx, v := range filteredNeighs {
		entry := fmt.Sprintf("%s dev %s", v.IP.String(), ifaceNames[idx])
		if v.HardwareAddr != nil {
			entry += fmt.Sprintf(" lladdr %s", v.HardwareAddr)
		}

		if v.Flags&netlink.NTF_ROUTER != 0 {
			entry += " router"
		}

		entry += " " + getState(v.State) + "\n"
		fmt.Fprintln(w, entry)
	}

	return nil
}

func neighShow(w io.Writer) error {
	iface, proxy, nud, err := parseNeighShowFlush()
	if err != nil {
		return err
	}

	if iface != nil {
		return showNeighbours(w, nud, proxy, nil, iface)
	}

	return showAllNeighbours(w, nud, proxy)
}

func neighFlush() error {
	var (
		ifaces []netlink.Link
		flags  uint8
		state  uint16
	)

	iface, proxy, nud, err := parseNeighShowFlush()
	if err != nil {
		return err
	}

	if iface == nil {
		ifaces, err = netlink.LinkList()
		if err != nil {
			return fmt.Errorf("failed to list interfaces: %w", err)
		}
	} else {
		ifaces = append(ifaces, iface)
	}

	if proxy {
		flags |= netlink.NTF_PROXY
	}

	if nud != -1 && nud <= math.MaxUint16 {
		state = uint16(nud)
	}

	for _, iface := range ifaces {

		msg := netlink.Ndmsg{
			Family: netlink.FAMILY_ALL,
			Index:  uint32(iface.Attrs().Index),
			Flags:  flags,
			State:  state,
		}

		neighbors, err := netlink.NeighListExecute(msg)
		if err != nil {
			return fmt.Errorf("failed to list neighbors: %w", err)
		}

		for _, neigh := range neighbors {
			if err := netlink.NeighDel(&neigh); err != nil {
				return fmt.Errorf("failed to delete neighbor: %w", err)
			}
		}
	}

	return nil
}
