// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
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

func (cmd cmd) neigh() error {
	if !cmd.tokenRemains() {
		return cmd.showAllNeighbours(-1, false)
	}

	switch c := cmd.findPrefix("show", "add", "del", "replace", "flush", "get", "help"); c {
	case "add", "del", "replace":
		neigh, err := cmd.parseNeighAddDelReplaceParams()
		if err != nil {
			return err
		}

		switch c {
		case "add":
			return cmd.handle.NeighAdd(neigh)
		case "del":
			return cmd.handle.NeighDel(neigh)
		case "replace":
			return cmd.handle.NeighSet(neigh)
		}

	case "show":
		return cmd.neighShow()
	case "flush":
		return cmd.neighFlush()
	case "get":
		ip, err := cmd.parseAddress()
		if err != nil {
			return err
		}

		iface, err := cmd.parseDeviceName(true)
		if err != nil {
			return err
		}

		return cmd.showNeighbours(-1, false, &ip, iface)
	case "help":
		fmt.Fprint(cmd.out, neighHelp)
		return nil
	}

	return cmd.usage()
}

func (cmd cmd) parseNeighAddDelReplaceParams() (*netlink.Neigh, error) {
	addr, err := cmd.parseAddress()
	if err != nil {
		return nil, err
	}

	var (
		iface       netlink.Link
		llAddr      net.HardwareAddr
		deviceFound bool
		state       int
		flag        int
	)

	for cmd.tokenRemains() {
		switch c := cmd.nextToken("dev", "lladdr", "nud", "router", "extern_learn"); c {
		case "dev":
			iface, err = cmd.parseDeviceName(true)
			if err != nil {
				return nil, err
			}
			deviceFound = true
		case "lladdr":
			llAddr, err = cmd.parseHardwareAddress()
			if err != nil {
				return nil, err
			}
		case "nud":
			state, err = parseValue[int](cmd, "STATE")
			if err != nil {
				return nil, err
			}

		case "router":
			flag |= netlink.NTF_ROUTER
		case "extern_learn":
			flag |= netlink.NTF_EXT_LEARNED
		default:
			return nil, fmt.Errorf("unsupported option %q, expected: %v", c, cmd.expectedValues)
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

func (cmd cmd) parseNeighShowFlush() (iface netlink.Link, proxy bool, nud int, err error) {
	nud = -1

	var ok bool

	for cmd.tokenRemains() {
		switch c := cmd.nextToken("dev", "proxy", "nud"); c {
		case "dev":
			dev, err := cmd.parseDeviceName(true)
			iface = dev
			if err != nil {
				return nil, false, 0, err
			}
		case "proxy":
			proxy = true
		case "nud":
			nudStr, err := parseValue[string](cmd, "STATE")
			if err != nil {
				return nil, false, 0, err
			}

			nud, ok = neighStatesMap[strings.ToLower(nudStr)]
			if !ok {
				return nil, false, 0, fmt.Errorf("invalid state %q", nudStr)
			}

		default:
			return nil, false, 0, fmt.Errorf("unsupported option %q, expected: %v", c, cmd.expectedValues)
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

var neighStatesMap = map[string]int{
	"none":       netlink.NUD_NONE,
	"incomplete": netlink.NUD_INCOMPLETE,
	"reachable":  netlink.NUD_REACHABLE,
	"stale":      netlink.NUD_STALE,
	"delay":      netlink.NUD_DELAY,
	"probe":      netlink.NUD_PROBE,
	"failed":     netlink.NUD_FAILED,
	"noarp":      netlink.NUD_NOARP,
	"permanent":  netlink.NUD_PERMANENT,
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

func (cmd cmd) showAllNeighbours(nud int, proxy bool) error {
	ifaces, err := cmd.handle.LinkList()
	if err != nil {
		return err
	}

	return cmd.showNeighbours(nud, proxy, nil, ifaces...)
}

type Neigh struct {
	Dst    net.IP `json:"dst"`
	Dev    string `json:"dev"`
	LLAddr string `json:"lladdr,omitempty"`
	State  string `json:"state,omitempty"`
}

func (cmd cmd) showNeighbours(nud int, proxy bool, address *net.IP, ifaces ...netlink.Link) error {
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
		neighs, err := cmd.handle.NeighListExecute(netlink.Ndmsg{
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

	if cmd.opts.json {
		neighs := make([]Neigh, 0, len(filteredNeighs))

		for idx, v := range filteredNeighs {
			neigh := Neigh{
				Dst:    v.IP,
				Dev:    ifaceNames[idx],
				LLAddr: v.HardwareAddr.String(),
			}

			if !cmd.opts.brief {
				neigh.State = getState(v.State)
			}

			neighs = append(neighs, neigh)
		}

		return printJSON(cmd, neighs)
	}

	neighFmt := "%s dev %s%s%s %s\n"
	neighBriefFmt := "%-39s %-13s %-9s\n"
	for idx, v := range filteredNeighs {
		if cmd.opts.brief {
			fmt.Fprintf(cmd.out, neighBriefFmt, v.IP, ifaceNames[idx], v.HardwareAddr)
		} else {
			llAddr := ""
			routerStr := ""

			if v.HardwareAddr != nil {
				llAddr = fmt.Sprintf(" lladdr %s", v.HardwareAddr)
			}

			if v.Flags&netlink.NTF_ROUTER != 0 {
				routerStr = " router"
			}

			fmt.Fprintf(cmd.out, neighFmt, v.IP, ifaceNames[idx], llAddr, routerStr, getState(v.State))
		}
	}

	return nil
}

func (cmd cmd) neighShow() error {
	iface, proxy, nud, err := cmd.parseNeighShowFlush()
	if err != nil {
		return err
	}

	if iface != nil {
		return cmd.showNeighbours(nud, proxy, nil, iface)
	}

	return cmd.showAllNeighbours(nud, proxy)
}

func (cmd cmd) neighFlush() error {
	var (
		ifaces []netlink.Link
		flags  uint8
		state  uint16
	)

	iface, proxy, nud, err := cmd.parseNeighShowFlush()
	if err != nil {
		return err
	}

	if iface == nil {
		ifaces, err = cmd.handle.LinkList()
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

		neighbors, err := cmd.handle.NeighListExecute(msg)
		if err != nil {
			return fmt.Errorf("failed to list neighbors: %w", err)
		}

		for _, neigh := range neighbors {
			if err := cmd.handle.NeighDel(&neigh); err != nil {
				return fmt.Errorf("failed to delete neighbor: %w", err)
			}
		}
	}

	return nil
}
