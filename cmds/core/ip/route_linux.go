// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"net"

	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

const routeHelp = `Usage: ip route { list | flush } SELECTOR

       ip route get ADDRESS
                [ from ADDRESS] [ iif STRING ]
                [ oif STRING ] [ vrf NAME ]
     
       ip route { add | del | append | replace } ROUTE

	   ip route help
SELECTOR := [ root PREFIX ] [ match PREFIX ] [ exact PREFIX ]
            [ table TABLE_ID ] [ proto RTPROTO ]
            [ type TYPE ] [ scope SCOPE ]
ROUTE := NODE_SPEC [ INFO_SPEC ]
NODE_SPEC := [ TYPE ] PREFIX [ tos TOS ]
             [ table TABLE_ID ] [ proto RTPROTO ]
             [ scope SCOPE ] [ metric METRIC ] OPTIONS
INFO_SPEC := [ nexthop NH ]...
NH := [ via ADDRESS ]
FAMILY := [ inet | inet6 | mpls | bridge | link ]
OPTIONS := FLAGS [ mtu NUMBER ] [ advmss NUMBER ]
           [ rtt TIME ] [ rttvar TIME ] [ reordering NUMBER ]
           [ window NUMBER ] [ cwnd NUMBER ] [ initcwnd NUMBER ]
           [ ssthresh NUMBER ] [ realms REALM ] [ src ADDRESS ]
           [ rto_min TIME ] [ hoplimit NUMBER ] [ initrwnd NUMBER ]
           [ features FEATURES ] [ quickack BOOL ] [ congctl NAME ]
		   [ fastopen_no_cookie BOOL ]
TYPE := { unicast | local | broadcast | multicast | throw |
          unreachable | prohibit | blackhole | nat }
TABLE_ID := [ local | main | default | all | NUMBER ]
SCOPE := [ host | link | global | NUMBER ]
BOOL := [1|0]
OPTIONS := OPTION [ OPTIONS ]
`

var routeTypes = map[string]int{
	"unicast":     unix.RTN_UNICAST,
	"local":       unix.RTN_LOCAL,
	"broadcast":   unix.RTN_BROADCAST,
	"multicast":   unix.RTN_MULTICAST,
	"throw":       unix.RTN_THROW,
	"unreachable": unix.RTN_UNREACHABLE,
	"prohibit":    unix.RTN_PROHIBIT,
	"blackhole":   unix.RTN_BLACKHOLE,
	"nat":         unix.RTN_NAT,
}

func routeAdddefault(w io.Writer) error {
	nh, nhval, err := parseNextHop()
	if err != nil {
		return err
	}
	// TODO: NHFLAGS.
	l, err := parseDeviceName(true)
	if err != nil {
		return err
	}
	switch nh {
	case "via":
		fmt.Fprintf(w, "Add default route %v via %v", nhval, l.Attrs().Name)
		r := &netlink.Route{LinkIndex: l.Attrs().Index, Gw: nhval}
		if err := netlink.RouteAdd(r); err != nil {
			return fmt.Errorf("error adding default route to %v: %v", l.Attrs().Name, err)
		}
		return nil
	}
	return usage()
}

func routeAdd(w io.Writer) error {
	ns := parseNodeSpec()
	switch ns {
	case "default":
		return routeAdddefault(w)
	default:
		route, d, err := parseRouteAddAppendReplaceDel(ns)
		if err != nil {
			return err
		}

		if err := netlink.RouteAdd(route); err != nil {
			return fmt.Errorf("error adding route %s -> %s: %v", route.Dst.IP, d.Attrs().Name, err)
		}
		return nil
	}
}

func routeAppend() error {
	ns := parseNodeSpec()
	route, d, err := parseRouteAddAppendReplaceDel(ns)
	if err != nil {
		return err
	}

	if err := netlink.RouteAppend(route); err != nil {
		return fmt.Errorf("error appending route %s -> %s: %v", route.Dst.IP, d.Attrs().Name, err)
	}
	return nil
}

func routeReplace() error {
	ns := parseNodeSpec()
	route, d, err := parseRouteAddAppendReplaceDel(ns)
	if err != nil {
		return err
	}

	if err := netlink.RouteReplace(route); err != nil {
		return fmt.Errorf("error appending route %s -> %s: %v", route.Dst.IP, d.Attrs().Name, err)
	}
	return nil
}

func routeDel() error {
	ns := parseNodeSpec()
	route, d, err := parseRouteAddAppendReplaceDel(ns)
	if err != nil {
		return err
	}

	if err := netlink.RouteDel(route); err != nil {
		return fmt.Errorf("error deleting route %s -> %s: %v", route.Dst.IP, d.Attrs().Name, err)
	}
	return nil
}

func parseRouteAddAppendReplaceDel(ns string) (*netlink.Route, netlink.Link, error) {
	var err error

	route := &netlink.Route{}

	_, route.Dst, err = net.ParseCIDR(ns)
	if err != nil {
		return nil, nil, err
	}

	d, err := parseDeviceName(true)
	if err != nil {
		return nil, nil, err
	}

	route.LinkIndex = d.Attrs().Index

	for {
		cursor++

		if cursor == len(arg) {
			break
		}
		whatIWant = []string{"type", "tos", "table", "proto", "scope", "metric", "mtu", "advmss", "rtt", "rttvar", "reordering", "window", "cwnd", "initcwnd", "ssthresh", "realms", "src", "rto_min", "hoplimit", "initrwnd", "congctl", "features", "quickack", "fastopen_no_cookie"}
		switch arg[cursor] {
		case "tos":
			cursor++
			route.Tos, err = parseInt()
			if err != nil {
				return nil, nil, err
			}

		case "table":
			cursor++
			route.Table, err = parseInt()
			if err != nil {
				return nil, nil, err
			}

		case "proto":
			cursor++
			proto, err := parseInt()
			if err != nil {
				return nil, nil, err
			}
			route.Protocol = netlink.RouteProtocol(proto)

		case "scope":
			cursor++
			scope, err := parseUint8()
			if err != nil {
				return nil, nil, err
			}
			route.Scope = netlink.Scope(scope)
		case "metric":
			cursor++
			route.Priority, err = parseInt()
			if err != nil {
				return nil, nil, err
			}
		case "mtu":
			cursor++
			route.MTU, err = parseInt()
			if err != nil {
				return nil, nil, err
			}
		case "advmss":
			cursor++
			route.AdvMSS, err = parseInt()
			if err != nil {
				return nil, nil, err
			}
		case "rtt":
			cursor++
			route.Rtt, err = parseInt()
			if err != nil {
				return nil, nil, err
			}
		case "rttvar":
			cursor++
			route.RttVar, err = parseInt()
			if err != nil {
				return nil, nil, err
			}
		case "reordering":
			cursor++
			route.Reordering, err = parseInt()
			if err != nil {
				return nil, nil, err
			}
		case "window":
			cursor++
			route.Window, err = parseInt()
			if err != nil {
				return nil, nil, err
			}
		case "cwnd":
			cursor++
			route.Cwnd, err = parseInt()
			if err != nil {
				return nil, nil, err
			}
		case "initcwnd":
			cursor++
			route.InitCwnd, err = parseInt()
			if err != nil {
				return nil, nil, err
			}
		case "ssthresh":
			cursor++
			route.Ssthresh, err = parseInt()
			if err != nil {
				return nil, nil, err
			}
		case "realms":
			cursor++
			route.Realm, err = parseInt()
			if err != nil {
				return nil, nil, err
			}
		case "src":
			cursor++
			route.Src = net.ParseIP(arg[cursor])
			if route.Src == nil {
				return nil, nil, fmt.Errorf("invalid source address: %v", arg[cursor])
			}
		case "rto_min":
			cursor++
			route.RtoMin, err = parseInt()
			if err != nil {
				return nil, nil, err
			}
		case "hoplimit":
			cursor++
			route.Hoplimit, err = parseInt()
			if err != nil {
				return nil, nil, err
			}
		case "initrwnd":
			cursor++
			route.InitRwnd, err = parseInt()
			if err != nil {
				return nil, nil, err
			}
		case "congctl":
			cursor++
			route.Congctl = parseString()
		case "features":
			cursor++
			route.Features, err = parseInt()
			if err != nil {
				return nil, nil, err
			}
		case "quickack":
			cursor++
			switch arg[cursor] {
			case "1":
				route.QuickACK = 1
			case "0":
				route.QuickACK = 0
			default:
				return nil, nil, usage()
			}
		case "fastopen_no_cookie":
			cursor++
			switch arg[cursor] {
			case "1":
				route.FastOpenNoCookie = 1
			case "0":
				route.FastOpenNoCookie = 0
			default:
				return nil, nil, usage()
			}
		default:
			return nil, nil, usage()
		}
	}

	return route, d, nil
}

func routeShow(w io.Writer) error {
	filter, filterMask, root, match, exact, err := parseRouteShowListFlush()
	if err != nil {
		return err
	}

	return showRoutes(w, filter, filterMask, root, match, exact, inet6)
}

func routeFlush() error {
	var f int

	if inet6 {
		f = netlink.FAMILY_V6
	} else {
		f = netlink.FAMILY_V4
	}

	filter, filterMask, root, match, exact, err := parseRouteShowListFlush()
	if err != nil {
		return err
	}

	routes, err := filteredRouteList(f, filter, filterMask, root, match, exact)
	if err != nil {
		return err
	}

	for _, route := range routes {
		if err := netlink.RouteDel(&route); err != nil {
			return err
		}
	}

	return nil
}

func parseRouteShowListFlush() (*netlink.Route, uint64, *net.IPNet, *net.IPNet, *net.IPNet, error) {
	var (
		filterMask uint64
		filter     netlink.Route
		root       *net.IPNet
		match      *net.IPNet
		exact      *net.IPNet
	)

	if routeType, ok := routeTypes[arg[cursor]]; ok {
		filter.Type = routeType
		filterMask |= netlink.RT_FILTER_TYPE
		cursor++
	}

	for {
		cursor++

		if cursor == len(arg) {
			break
		}

		switch arg[cursor] {
		case "scope":
			filterMask |= netlink.RT_FILTER_SCOPE
			cursor++
			scope, err := parseUint8()
			if err != nil {
				return nil, 0, nil, nil, nil, err
			}
			filter.Scope = netlink.Scope(scope)

		case "table":
			filterMask |= netlink.RT_FILTER_TABLE
			cursor++
			table, err := parseInt()
			if err != nil {
				return nil, 0, nil, nil, nil, err
			}
			filter.Table = table

		case "proto":
			filterMask |= netlink.RT_FILTER_PROTOCOL
			cursor++
			proto, err := parseInt()
			if err != nil {
				return nil, 0, nil, nil, nil, err
			}
			filter.Protocol = netlink.RouteProtocol(proto)

		case "root":
			cursor++
			_, prefix, err := net.ParseCIDR(arg[cursor])
			if err != nil {
				return nil, 0, nil, nil, nil, err
			}
			root = prefix

		case "match":
			cursor++
			_, prefix, err := net.ParseCIDR(arg[cursor])
			if err != nil {
				return nil, 0, nil, nil, nil, err
			}
			match = prefix

		case "exact":
			cursor++
			_, prefix, err := net.ParseCIDR(arg[cursor])
			if err != nil {
				return nil, 0, nil, nil, nil, err
			}
			exact = prefix
		default:
			return nil, 0, nil, nil, nil, usage()
		}

	}

	return &filter, filterMask, root, match, exact, nil
}

// showRoutes prints the routes in the system.
// If filterMask is not zero, only routes that match the filter are printed.
func showRoutes(w io.Writer, route *netlink.Route, filterMask uint64, root, match, exact *net.IPNet, inet6 bool) error {
	var f int

	if inet6 {
		f = netlink.FAMILY_V6
	} else {
		f = netlink.FAMILY_V4
	}

	routes, err := filteredRouteList(f, route, filterMask, root, match, exact)
	if err != nil {
		return err
	}

	for _, route := range routes {
		link, err := netlink.LinkByIndex(route.LinkIndex)
		if err != nil {
			return err
		}
		if route.Dst == nil {
			defaultRoute(w, route, link)
		} else {
			showRoute(w, route, link, f)
		}
	}
	return nil
}

func filteredRouteList(f int, route *netlink.Route, filterMask uint64, root, match, exact *net.IPNet) ([]netlink.Route, error) {
	var matchedRoutes []netlink.Route

	routes, err := netlink.RouteListFiltered(f, route, filterMask)
	if err != nil {
		return matchedRoutes, err
	}

	if root == nil && match == nil && exact == nil {
		matchedRoutes = routes
	} else {
		matchedRoutes, err = matchRoutes(routes, root, match, exact)
		if err != nil {
			return matchedRoutes, err
		}
	}

	return matchedRoutes, nil
}

// matchRoutes matches routes against a given prefix.
func matchRoutes(routes []netlink.Route, root, match, exact *net.IPNet) ([]netlink.Route, error) {
	matchedRoutes := []netlink.Route{}

	for _, route := range routes {
		if root != nil && !root.Contains(route.Dst.IP) {
			continue
		}

		if match != nil && !match.Contains(route.Dst.IP) {
			continue
		}

		if exact != nil && !exact.IP.Equal(route.Dst.IP) {
			continue
		}

		matchedRoutes = append(matchedRoutes, route)
	}

	return matchedRoutes, nil
}

func showRoutesForAddress(w io.Writer, addr net.IP, options *netlink.RouteGetOptions, inet6 bool) error {
	var f int

	if inet6 {
		f = netlink.FAMILY_V6
	} else {
		f = netlink.FAMILY_V4
	}

	routes, err := netlink.RouteGetWithOptions(addr, options)
	if err != nil {
		return err
	}

	for _, route := range routes {
		link, err := netlink.LinkByIndex(route.LinkIndex)
		if err != nil {
			return err
		}
		if route.Dst == nil {
			defaultRoute(w, route, link)
		} else {
			showRoute(w, route, link, f)
		}
	}
	return nil
}

// routing protocol identifier
// specified in Linux Kernel header: include/uapi/linux/rtnetlink.h
// See man IP-ROUTE(8) and RTNETLINK(7)
var rtProto = map[int]string{
	unix.RTPROT_BABEL:    "babel",
	unix.RTPROT_BGP:      "bgp",
	unix.RTPROT_BIRD:     "bird",
	unix.RTPROT_BOOT:     "boot",
	unix.RTPROT_DHCP:     "dhcp",
	unix.RTPROT_DNROUTED: "dnrouted",
	unix.RTPROT_EIGRP:    "eigrp",
	unix.RTPROT_GATED:    "gated",
	unix.RTPROT_ISIS:     "isis",
	unix.RTPROT_KERNEL:   "kernel",
	unix.RTPROT_MROUTED:  "mrouted",
	unix.RTPROT_MRT:      "mrt",
	unix.RTPROT_NTK:      "ntk",
	unix.RTPROT_OSPF:     "ospf",
	unix.RTPROT_RA:       "ra",
	unix.RTPROT_REDIRECT: "redirect",
	unix.RTPROT_RIP:      "rip",
	unix.RTPROT_STATIC:   "static",
	unix.RTPROT_UNSPEC:   "unspec",
	unix.RTPROT_XORP:     "xorp",
	unix.RTPROT_ZEBRA:    "zebra",
}

const (
	defaultFmt   = "default via %v dev %s proto %s metric %d\n"
	routeFmt     = "%v dev %s proto %s scope %s src %s metric %d\n"
	route6Fmt    = "%s dev %s proto %s metric %d\n"
	routeVia6Fmt = "%s via %s dev %s proto %s metric %d\n"
)

func defaultRoute(w io.Writer, r netlink.Route, l netlink.Link) {
	gw := r.Gw
	name := l.Attrs().Name
	proto := rtProto[int(r.Protocol)]
	metric := r.Priority
	fmt.Fprintf(w, defaultFmt, gw, name, proto, metric)
}

func showRoute(w io.Writer, r netlink.Route, l netlink.Link, f int) {
	dest := r.Dst
	name := l.Attrs().Name
	proto := rtProto[int(r.Protocol)]
	metric := r.Priority

	switch f {
	case netlink.FAMILY_V4:
		scope := addrScopes[r.Scope]
		src := r.Src
		fmt.Fprintf(w, routeFmt, dest, name, proto, scope, src, metric)
	case netlink.FAMILY_V6:
		if r.Gw != nil {
			gw := r.Gw
			fmt.Fprintf(w, routeVia6Fmt, dest, gw, name, proto, metric)
		} else {
			fmt.Fprintf(w, route6Fmt, dest, name, proto, metric)
		}
	}
}

func routeGet(w io.Writer) error {
	cursor++
	whatIWant = []string{"CIDR Address"}
	addr, _, err := net.ParseCIDR(arg[cursor])
	if err != nil {
		return err
	}

	options, err := parseRouteGet()
	if err != nil {
		return err
	}

	return showRoutesForAddress(w, addr, options, inet6)
}

func parseRouteGet() (*netlink.RouteGetOptions, error) {
	var opts netlink.RouteGetOptions
	for {
		cursor++

		if cursor == len(arg) {
			break
		}
		switch arg[cursor] {
		case "oif":
			cursor++
			opts.Oif = arg[cursor]
		case "iif":
			cursor++
			opts.Iif = arg[cursor]
		case "vrf":
			cursor++
			opts.VrfName = arg[cursor]
		case "from":
			cursor++
			opts.SrcAddr = net.ParseIP(arg[cursor])
		default:
			return nil, usage()
		}
	}

	return &opts, nil
}

func route(w io.Writer) error {
	cursor++
	if len(arg[cursor:]) == 0 {
		return routeShow(w)
	}

	whatIWant = []string{"show", "add", "append", "replace", "del", "list", "get", "help"}
	switch findPrefix(arg[cursor], whatIWant) {
	case "add":
		return routeAdd(w)
	case "append":
		return routeAppend()
	case "replace":
		return routeReplace()
	case "del":
		return routeDel()
	case "show", "list":
		return routeShow(w)
	case "flush":
		return routeFlush()
	case "get":
		return routeGet(w)
	case "help":
		fmt.Fprint(w, routeHelp)
		return nil
	}
	return usage()
}
