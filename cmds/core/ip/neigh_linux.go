// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"bytes"
	"fmt"
	"math"
	"net"
	"sort"
	"strconv"
	"strings"

	"github.com/vishvananda/netlink"
)

const neighHelp = `Usage: 	ip neigh { add | del | replace }
            ADDR [ lladdr LLADDR ] [ nud STATE ] [ dev DEV ] [ router ] [ extern_learn ] 

	ip neigh { show | flush } [ proxy ] [ to PREFIX ] [ dev DEV ] [ nud STATE ]

	ip neigh get ADDR dev DEV

STATE := { delay | failed | incomplete | noarp | none | permanent | probe | reachable | stale }
`

func (cmd *cmd) neigh() error {
	if !cmd.tokenRemains() {
		return cmd.showAllNeighbours(nil, nil, -1, false)
	}

	switch c := cmd.findPrefix("show", "add", "delete", "replace", "flush", "get", "help"); c {
	case "add", "delete", "replace":
		neigh, err := cmd.parseNeighAddDelReplaceParams()
		if err != nil {
			return err
		}

		switch c {
		case "add":
			return cmd.handle.NeighAdd(neigh)
		case "delete":
			return cmd.handle.NeighDel(neigh)
		case "replace":
			return cmd.handle.NeighSet(neigh)
		default:
			fmt.Fprint(cmd.Out, neighHelp)
			return nil
		}

	case "show":
		return cmd.neighShow()
	case "flush":
		return cmd.neighFlush()
	case "get":
		ip, iface, err := cmd.parseNeighGet()
		if err != nil {
			return err
		}

		return cmd.showNeighbours(-1, false, &ip, nil, iface)
	case "help":
		fmt.Fprint(cmd.Out, neighHelp)
		return nil
	}

	return cmd.usage()
}

func (cmd cmd) parseNeighGet() (net.IP, netlink.Link, error) {
	ip, err := cmd.parseAddress()
	if err != nil {
		return nil, nil, err
	}

	iface, err := cmd.parseDeviceName(true)
	if err != nil {
		return nil, nil, err
	}

	return ip, iface, nil
}

func (cmd *cmd) parseNeighAddDelReplaceParams() (*netlink.Neigh, error) {
	addr, err := cmd.parseAddress()
	if err != nil {
		return nil, err
	}

	var (
		iface       netlink.Link
		llAddr      net.HardwareAddr
		deviceFound bool
		state       int = netlink.NUD_PERMANENT
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
			state, err = parseNUD(cmd.nextToken("STATE"))
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

func (cmd *cmd) parseNeighShowFlush() (addr net.IP, subNet *net.IPNet, iface netlink.Link, proxy bool, nud int, err error) {
	nud = -1

	for cmd.tokenRemains() {
		switch c := cmd.nextToken("dev", "proxy", "nud", "to"); c {
		case "to":
			addr, subNet, err = cmd.parseAddressorCIDR()
			if err != nil {
				return nil, nil, nil, false, 0, err
			}
		case "dev":
			dev, err := cmd.parseDeviceName(true)
			iface = dev
			if err != nil {
				return nil, nil, nil, false, 0, err
			}
		case "proxy":
			proxy = true
		case "nud":
			nud, err = parseNUD(cmd.nextToken("STATE"))
			if err != nil {
				return nil, nil, nil, false, 0, err
			}

		default:
			return nil, nil, nil, false, 0, fmt.Errorf("unsupported option %q, expected: %v", c, cmd.ExpectedValues)
		}
	}

	return addr, subNet, iface, proxy, nud, nil
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

func parseNUD(input string) (int, error) {
	var nud int

	if val, ok := neighStatesMap[strings.ToLower(input)]; ok {
		nud = val
	} else {
		var err error
		nudInt64, err := strconv.ParseInt(input, 10, 0)
		if err != nil {
			return nud, fmt.Errorf(`argument "%v" is wrong: nud state is bad`, input)
		}

		if nudInt64 < 0 {
			return 0, fmt.Errorf(`argument "%v" is wrong: nud state is bad`, input)
		}

		nud = int(nudInt64)

		if _, ok := neighStates[nud]; !ok {
			return nud, fmt.Errorf(`argument "%v" is wrong: nud state is bad`, input)
		}

	}

	return nud, nil
}

func getState(state int) string {
	ret := make([]string, 0)
	for st, name := range neighStates {
		if state == st {
			ret = append(ret, name)
		}
	}
	if len(ret) == 0 {
		return "UNKNOWN"
	}
	return strings.Join(ret, ",")
}

func (cmd *cmd) showAllNeighbours(address net.IP, subNet *net.IPNet, nud int, proxy bool) error {
	ifaces, err := cmd.handle.LinkList()
	if err != nil {
		return err
	}

	return cmd.showNeighbours(nud, proxy, &address, subNet, ifaces...)
}

// NeighJSON represents a neighbor object for JSON output format.
type NeighJSON struct {
	Dst    net.IP `json:"dst"`
	Dev    string `json:"dev"`
	LLAddr string `json:"lladdr,omitempty"`
	State  string `json:"state,omitempty"`
}

func (cmd *cmd) showNeighbours(nud int, proxy bool, address *net.IP, subNet *net.IPNet, ifaces ...netlink.Link) error {
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

	filteredNeighs, filteredLinkNames := filterNeighsByAddr(neighs, linkNames, address, subNet)

	return cmd.printNeighs(filteredNeighs, filteredLinkNames)
}

func filterNeighsByAddr(neighs []netlink.Neigh, linkNames []string, addr *net.IP, network *net.IPNet) ([]netlink.Neigh, []string) {
	filtered := make([]netlink.Neigh, 0)
	filteredLinkNames := make([]string, 0)

	for idx, neigh := range neighs {
		// If network is specified, check if IP is in the subnet
		if network != nil {
			if !network.Contains(neigh.IP) {
				continue
			}
		} else if addr != nil && *addr != nil {
			// Otherwise do exact IP matching
			if !neigh.IP.Equal(*addr) {
				continue
			}
		}

		if neigh.State != netlink.NUD_NOARP {
			filtered = append(filtered, neigh)

			if linkNames != nil {
				filteredLinkNames = append(filteredLinkNames, linkNames[idx])
			}
		}
	}
	return filtered, filteredLinkNames
}

var flagOrder = []int{
	netlink.NTF_ROUTER, // Display first
	netlink.NTF_EXT_LEARNED,
}

var flagToString = map[int]string{
	netlink.NTF_EXT_LEARNED: "extern_learn",
	netlink.NTF_ROUTER:      "router",
}

func (cmd *cmd) printNeighs(neighs []netlink.Neigh, ifacesNames []string) error {
	neighs, ifacesNames = sortedNeighs(neighs, ifacesNames)

	if cmd.Opts.JSON {
		pNeighs := make([]NeighJSON, 0, len(neighs))

		for idx, v := range neighs {
			neigh := NeighJSON{
				Dst:    v.IP,
				Dev:    ifacesNames[idx],
				LLAddr: v.HardwareAddr.String(),
			}

			if !cmd.Opts.Brief {
				neigh.State = getState(v.State)
			}

			pNeighs = append(pNeighs, neigh)
		}

		return printJSON(*cmd, pNeighs)
	}

	neighFmt := "%s dev %s%s%s %s\n"
	neighBriefFmt := "%-39s %-13s %-9s\n"
	for idx, v := range neighs {
		if cmd.Opts.Brief {
			fmt.Fprintf(cmd.Out, neighBriefFmt, v.IP, ifacesNames[idx], v.HardwareAddr)
		} else {
			llAddr := ""
			flags := ""

			if v.HardwareAddr != nil {
				llAddr = fmt.Sprintf(" lladdr %s", v.HardwareAddr)
			}

			for _, flag := range flagOrder {
				if v.Flags&flag != 0 {
					flags += fmt.Sprintf(" %s", flagToString[flag])
				}
			}

			fmt.Fprintf(cmd.Out, neighFmt, v.IP, ifacesNames[idx], llAddr, flags, getState(v.State))
		}
	}

	return nil
}

func sortedNeighs(neighs []netlink.Neigh, ifacesNames []string) ([]netlink.Neigh, []string) {
	type pair struct {
		neigh     netlink.Neigh
		ifaceName string
	}

	pairs := make([]pair, len(neighs))
	for i := range neighs {
		pairs[i] = pair{neighs[i], ifacesNames[i]}
	}

	sort.SliceStable(pairs, func(i, j int) bool {
		// First priority: IPv4 before IPv6
		isIPv4_i := pairs[i].neigh.IP.To4() != nil
		isIPv4_j := pairs[j].neigh.IP.To4() != nil
		if isIPv4_i != isIPv4_j {
			return isIPv4_i
		}

		// Second priority: By device index (hardware order)
		if pairs[i].neigh.LinkIndex != pairs[j].neigh.LinkIndex {
			return pairs[i].neigh.LinkIndex < pairs[j].neigh.LinkIndex
		}

		// Third: Use address comparison
		return bytes.Compare(pairs[i].neigh.IP, pairs[j].neigh.IP) < 0
	})

	for i, p := range pairs {
		neighs[i] = p.neigh
		ifacesNames[i] = p.ifaceName
	}

	return neighs, ifacesNames
}

func (cmd *cmd) neighShow() error {
	addr, subNet, iface, proxy, nud, err := cmd.parseNeighShowFlush()
	if err != nil {
		return err
	}

	if iface != nil {
		return cmd.showNeighbours(nud, proxy, &addr, subNet, iface)
	}

	return cmd.showAllNeighbours(addr, subNet, nud, proxy)
}

func (cmd *cmd) neighFlush() error {
	var (
		ifaces []netlink.Link
		flags  uint8
		state  uint16
	)

	address, subNet, iface, proxy, nud, err := cmd.parseNeighShowFlush()
	if err != nil {
		return err
	}

	if address == nil && subNet == nil && iface == nil && !proxy && nud == -1 {
		//nolint:revive,staticcheck // This message is analog to the one in iproute2
		return fmt.Errorf("Flush requires arguments.")
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

		// Filter neighbors based on NUD state (keep specific state or exclude PERMANENT and NOARP when no state specified)
		filteredForStates := make([]netlink.Neigh, 0, len(neighbors))
		for _, neigh := range neighbors {
			if (nud != -1 && neigh.State == nud) || (nud == -1 && neigh.State != netlink.NUD_PERMANENT && neigh.State != netlink.NUD_NOARP) {
				filteredForStates = append(filteredForStates, neigh)
			}
		}

		neighbors = filteredForStates

		filteredNeighs, _ := filterNeighsByAddr(neighbors, nil, &address, subNet)

		for _, neigh := range filteredNeighs {
			if err := cmd.handle.NeighDel(&neigh); err != nil {
				return fmt.Errorf("failed to delete neighbor: %w", err)
			}
		}
	}

	return nil
}

func (cmd *cmd) neighFlagState(proxy bool, nud int) (uint8, uint16, error) {
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
