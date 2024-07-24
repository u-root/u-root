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
		fmt.Fprint(cmd.Out, neighHelp)
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
			state, err = cmd.parseInt("STATE")
			if err != nil {
				return nil, err
			}

		case "router":
			flag |= netlink.NTF_ROUTER
		case "extern_learn":
			flag |= netlink.NTF_EXT_LEARNED
		default:
			return nil, fmt.Errorf("unsupported option %q, expected: %v", c, cmd.ExpectedValues)
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
			nudStr := cmd.nextToken("STATE")

			nud, ok = neighStatesMap[strings.ToLower(nudStr)]
			if !ok {
				return nil, false, 0, fmt.Errorf("invalid state %q", nudStr)
			}

		default:
			return nil, false, 0, fmt.Errorf("unsupported option %q, expected: %v", c, cmd.ExpectedValues)
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
	flags, state, err := cmd.neighFlagState(proxy, nud)
	if err != nil {
		return err
	}

	neighs := make([]netlink.Neigh, 0)
	linkNames := make([]string, 0)

	for _, iface := range ifaces {
		linkNeighs, err := cmd.handle.NeighListExecute(netlink.Ndmsg{
			Family: uint8(cmd.Family),
			Index:  uint32(iface.Attrs().Index),
			Flags:  flags,
			State:  state,
		})
		if err != nil {
			return err
		}

		neighs = append(neighs, linkNeighs...)

		for range linkNeighs {
			linkNames = append(linkNames, iface.Attrs().Name)
		}
	}

	filteredNeighs, filteredLinkNames := filterNeighsByAddr(neighs, linkNames, address)

	return cmd.printNeighs(filteredNeighs, filteredLinkNames)
}

func filterNeighsByAddr(neighs []netlink.Neigh, linkNames []string, addr *net.IP) ([]netlink.Neigh, []string) {
	filtered := make([]netlink.Neigh, 0)
	filteredLinkNames := make([]string, 0)

	for idx, neigh := range neighs {
		if addr != nil {
			if *addr != nil && !neigh.IP.Equal(*addr) {
				continue
			}
		}
		if neigh.State != netlink.NUD_NOARP {
			filtered = append(filtered, neigh)
			filteredLinkNames = append(filteredLinkNames, linkNames[idx])
		}
	}
	return filtered, filteredLinkNames
}

func (cmd cmd) printNeighs(neighs []netlink.Neigh, ifacesNames []string) error {
	if cmd.Opts.JSON {
		pNeighs := make([]Neigh, 0, len(neighs))

		for idx, v := range neighs {
			neigh := Neigh{
				Dst:    v.IP,
				Dev:    ifacesNames[idx],
				LLAddr: v.HardwareAddr.String(),
			}

			if !cmd.Opts.Brief {
				neigh.State = getState(v.State)
			}

			pNeighs = append(pNeighs, neigh)
		}

		return printJSON(cmd, pNeighs)
	}

	neighFmt := "%s dev %s%s%s %s\n"
	neighBriefFmt := "%-39s %-13s %-9s\n"
	for idx, v := range neighs {
		if cmd.Opts.Brief {
			fmt.Fprintf(cmd.Out, neighBriefFmt, v.IP, ifacesNames[idx], v.HardwareAddr)
		} else {
			llAddr := ""
			routerStr := ""

			if v.HardwareAddr != nil {
				llAddr = fmt.Sprintf(" lladdr %s", v.HardwareAddr)
			}

			if v.Flags&netlink.NTF_ROUTER != 0 {
				routerStr = " router"
			}

			fmt.Fprintf(cmd.Out, neighFmt, v.IP, ifacesNames[idx], llAddr, routerStr, getState(v.State))
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

	flags, state, err = cmd.neighFlagState(proxy, nud)
	if err != nil {
		return err
	}

	for _, iface := range ifaces {

		msg := netlink.Ndmsg{
			Family: uint8(cmd.Family),
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

func (cmd cmd) neighFlagState(proxy bool, nud int) (uint8, uint16, error) {
	var flags uint8
	var state uint16

	if cmd.Family < 0 || cmd.Family > 255 {
		return 0, 0, fmt.Errorf("invalid family %d", cmd.Family)
	}

	if proxy {
		flags |= netlink.NTF_PROXY
	}

	if nud != -1 && nud <= math.MaxUint16 {
		state = uint16(nud)
	}

	return flags, state, nil
}
