// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"fmt"
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

func addrScopeStr(scope netlink.Scope) string {
	switch scope {
	case netlink.SCOPE_UNIVERSE:
		return "global"
	default:
		return scope.String()
	}
}

func routeTypeToString(routeType int) string {
	for key, value := range routeTypes {
		if value == routeType {
			return key
		}
	}
	return "unknown"
}

func (cmd *cmd) routeAdddefault() error {
	nh, nhval, err := cmd.parseNextHop()
	if err != nil {
		return err
	}
	// TODO: NHFLAGS.
	l, err := cmd.parseDeviceName(true)
	if err != nil {
		return err
	}
	switch nh {
	case "via":
		fmt.Fprintf(cmd.Out, "Add default route %v via %v", nhval, l.Attrs().Name)
		r := &netlink.Route{LinkIndex: l.Attrs().Index, Gw: nhval}
		if err := cmd.handle.RouteAdd(r); err != nil {
			return fmt.Errorf("error adding default route to %v: %w", l.Attrs().Name, err)
		}
		return nil
	}
	return cmd.usage()
}

func (cmd *cmd) routeAdd() error {
	ns := cmd.nextToken("default", "CIDR")
	switch ns {
	case "default":
		return cmd.routeAdddefault()
	default:
		route, d, err := cmd.parseRouteAddAppendReplaceDel(ns)
		if err != nil {
			return err
		}

		link, err := netlink.LinkByName(d)
		if err != nil {
			return fmt.Errorf("error getting link %s: %w", d, err)
		}

		route.LinkIndex = link.Attrs().Index

		if err := cmd.handle.RouteAdd(route); err != nil {
			return fmt.Errorf("error adding route %s -> %s: %w", route.Dst.IP, d, err)
		}
		return nil
	}
}

func (cmd *cmd) routeAppend() error {
	ns := cmd.nextToken("default", "CIDR")
	route, d, err := cmd.parseRouteAddAppendReplaceDel(ns)
	if err != nil {
		return err
	}

	link, err := netlink.LinkByName(d)
	if err != nil {
		return fmt.Errorf("error getting link %s: %w", d, err)
	}

	route.LinkIndex = link.Attrs().Index

	if err := cmd.handle.RouteAppend(route); err != nil {
		return fmt.Errorf("error appending route %s -> %s: %w", route.Dst.IP, d, err)
	}
	return nil
}

func (cmd *cmd) routeReplace() error {
	ns := cmd.nextToken("default", "CIDR")
	route, d, err := cmd.parseRouteAddAppendReplaceDel(ns)
	if err != nil {
		return err
	}

	link, err := netlink.LinkByName(d)
	if err != nil {
		return fmt.Errorf("error getting link %s: %w", d, err)
	}

	route.LinkIndex = link.Attrs().Index

	if err := cmd.handle.RouteReplace(route); err != nil {
		return fmt.Errorf("error appending route %s -> %s: %w", route.Dst.IP, d, err)
	}
	return nil
}

func (cmd *cmd) routeDel() error {
	ns := cmd.nextToken("default", "CIDR")
	route, d, err := cmd.parseRouteAddAppendReplaceDel(ns)
	if err != nil {
		return err
	}

	link, err := netlink.LinkByName(d)
	if err != nil {
		return fmt.Errorf("error getting link %s: %w", d, err)
	}

	route.LinkIndex = link.Attrs().Index

	if err := cmd.handle.RouteDel(route); err != nil {
		return fmt.Errorf("error deleting route %s -> %s: %w", route.Dst.IP, d, err)
	}
	return nil
}

func (cmd *cmd) parseRouteAddAppendReplaceDel(ns string) (*netlink.Route, string, error) {
	var err error

	route := &netlink.Route{}

	_, route.Dst, err = net.ParseCIDR(ns)
	if err != nil {
		return nil, "", err
	}

	d := cmd.nextToken("dev", "device-name")
	if d == "dev" {
		d = cmd.nextToken("device-name")
	}

	for cmd.tokenRemains() {
		switch cmd.nextToken("type", "tos", "table", "proto", "scope", "metric", "mtu", "advmss", "rtt", "rttvar", "reordering", "window", "cwnd", "initcwnd", "ssthresh", "realms", "src", "rto_min", "hoplimit", "initrwnd", "congctl", "features", "quickack", "fastopen_no_cookie") {
		case "tos":
			route.Tos, err = cmd.parseInt("TOS")
			if err != nil {
				return nil, "", err
			}

		case "table":
			route.Table, err = cmd.parseInt("TABLE_ID")
			if err != nil {
				return nil, "", err
			}

		case "proto":
			proto, err := cmd.parseInt("RTPROTO")
			if err != nil {
				return nil, "", err
			}

			route.Protocol = netlink.RouteProtocol(proto)

		case "scope":
			scope, err := cmd.parseUint8("SCOPE")
			if err != nil {
				return nil, "", err
			}
			route.Scope = netlink.Scope(scope)
		case "metric":
			route.Priority, err = cmd.parseInt("METRIC")
			if err != nil {
				return nil, "", err
			}
		case "mtu":
			route.MTU, err = cmd.parseInt("NUMBER")
			if err != nil {
				return nil, "", err
			}
		case "advmss":
			route.AdvMSS, err = cmd.parseInt("NUMBER")
			if err != nil {
				return nil, "", err
			}
		case "rtt":
			route.Rtt, err = cmd.parseInt("TIME")
			if err != nil {
				return nil, "", err
			}
		case "rttvar":
			route.RttVar, err = cmd.parseInt("TIME")
			if err != nil {
				return nil, "", err
			}
		case "reordering":
			route.Reordering, err = cmd.parseInt("NUMBER")
			if err != nil {
				return nil, "", err
			}
		case "window":
			route.Window, err = cmd.parseInt("NUMBER")
			if err != nil {
				return nil, "", err
			}
		case "cwnd":
			route.Cwnd, err = cmd.parseInt("NUMBER")
			if err != nil {
				return nil, "", err
			}
		case "initcwnd":
			route.InitCwnd, err = cmd.parseInt("NUMBER")
			if err != nil {
				return nil, "", err
			}
		case "ssthresh":
			route.Ssthresh, err = cmd.parseInt("NUMBER")
			if err != nil {
				return nil, "", err
			}
		case "realms":
			route.Realm, err = cmd.parseInt("REALM")
			if err != nil {
				return nil, "", err
			}
		case "src":
			token := cmd.nextToken("ADDRESS")
			route.Src = net.ParseIP(token)
			if route.Src == nil {
				return nil, "", fmt.Errorf("invalid source address: %v", token)
			}
		case "rto_min":
			route.RtoMin, err = cmd.parseInt("TIME")
			if err != nil {
				return nil, "", err
			}
		case "hoplimit":
			route.Hoplimit, err = cmd.parseInt("NUMBER")
			if err != nil {
				return nil, "", err
			}
		case "initrwnd":
			route.InitRwnd, err = cmd.parseInt("NUMBER")
			if err != nil {
				return nil, "", err
			}
		case "congctl":
			route.Congctl = cmd.nextToken("NAME")
		case "features":
			route.Features, err = cmd.parseInt("FEATURES")
			if err != nil {
				return nil, "", err
			}
		case "quickack":
			switch cmd.nextToken("0", "1") {
			case "1":
				route.QuickACK = 1
			case "0":
				route.QuickACK = 0
			default:
				return nil, "", cmd.usage()
			}
		case "fastopen_no_cookie":
			switch cmd.nextToken("0", "1") {
			case "1":
				route.FastOpenNoCookie = 1
			case "0":
				route.FastOpenNoCookie = 0
			default:
				return nil, "", cmd.usage()
			}
		default:
			return nil, "", cmd.usage()
		}
	}

	return route, d, nil
}

func (cmd *cmd) routeShow() error {
	filter, filterMask, root, match, exact, err := cmd.parseRouteShowListFlush()
	if err != nil {
		return err
	}

	routeList, ifaceNames, err := cmd.filteredRouteList(filter, filterMask, root, match, exact)
	if err != nil {
		return err
	}

	return cmd.showRoutes(routeList, ifaceNames)
}

func (cmd *cmd) showAllRoutes() error {
	routeList, ifaceNames, err := cmd.filteredRouteList(nil, 0, nil, nil, nil)
	if err != nil {
		return err
	}

	return cmd.showRoutes(routeList, ifaceNames)
}

func (cmd *cmd) routeFlush() error {
	filter, filterMask, root, match, exact, err := cmd.parseRouteShowListFlush()
	if err != nil {
		return err
	}

	routes, _, err := cmd.filteredRouteList(filter, filterMask, root, match, exact)
	if err != nil {
		return err
	}

	for _, route := range routes {
		if err := cmd.handle.RouteDel(&route); err != nil {
			return err
		}
	}

	return nil
}

func (cmd *cmd) parseRouteShowListFlush() (*netlink.Route, uint64, *net.IPNet, *net.IPNet, *net.IPNet, error) {
	var (
		filterMask uint64
		filter     netlink.Route
		root       *net.IPNet
		match      *net.IPNet
		exact      *net.IPNet
	)

	for cmd.tokenRemains() {
		switch cmd.nextToken("scope", "table", "proto", "root", "match", "exact", "type") {
		case "scope":
			filterMask |= netlink.RT_FILTER_SCOPE
			scope, err := cmd.parseUint8("SCOPE")
			if err != nil {
				return nil, 0, nil, nil, nil, err
			}
			filter.Scope = netlink.Scope(scope)

		case "table":
			filterMask |= netlink.RT_FILTER_TABLE
			table, err := cmd.parseInt("TABLE_ID")
			if err != nil {
				return nil, 0, nil, nil, nil, err
			}
			filter.Table = table

		case "proto":
			filterMask |= netlink.RT_FILTER_PROTOCOL
			proto, err := cmd.parseInt("RTPROTO")
			if err != nil {
				return nil, 0, nil, nil, nil, err
			}
			filter.Protocol = netlink.RouteProtocol(proto)

		case "root":
			token := cmd.nextToken("PREFIX")
			_, prefix, err := net.ParseCIDR(token)
			if err != nil {
				return nil, 0, nil, nil, nil, err
			}
			root = prefix

		case "match":
			token := cmd.nextToken("PREFIX")
			_, prefix, err := net.ParseCIDR(token)
			if err != nil {
				return nil, 0, nil, nil, nil, err
			}
			match = prefix

		case "exact":
			token := cmd.nextToken("PREFIX")
			_, prefix, err := net.ParseCIDR(token)
			if err != nil {
				return nil, 0, nil, nil, nil, err
			}
			exact = prefix
		case "type":
			if routeType, ok := routeTypes[cmd.nextToken()]; ok {
				filter.Type = routeType
				filterMask |= netlink.RT_FILTER_TYPE
			} else {
				return nil, 0, nil, nil, nil, cmd.usage()
			}
		default:
			return nil, 0, nil, nil, nil, cmd.usage()
		}
	}

	return &filter, filterMask, root, match, exact, nil
}

// RouteJSON represents a route entry for JSON output format.
type RouteJSON struct {
	Dst      string   `json:"dst"`
	Dev      string   `json:"dev"`
	Protocol string   `json:"protocol"`
	Scope    string   `json:"scope"`
	PrefSrc  string   `json:"prefsrc"`
	Flags    []string `json:"flags,omitempty"`
}

// showRoutes prints the routes in the system.
func (cmd *cmd) showRoutes(routes []netlink.Route, ifaceNames []string) error {
	if cmd.Opts.JSON {
		obj := make([]RouteJSON, 0, len(routes))

		for idx, route := range routes {

			pRoute := RouteJSON{
				Dst:   route.Dst.String(),
				Dev:   ifaceNames[idx],
				Scope: route.Scope.String(),
			}

			if !cmd.Opts.Numeric {
				pRoute.Protocol = rtProto[int(route.Protocol)]
				pRoute.Scope = route.Scope.String()
			} else {
				pRoute.Protocol = fmt.Sprintf("%d", route.Protocol)
				pRoute.Scope = fmt.Sprintf("%d", route.Scope)
			}

			if route.Src != nil {
				pRoute.PrefSrc = route.Src.String()
			}

			if len(route.ListFlags()) != 0 {
				pRoute.Flags = route.ListFlags()
			}

			obj = append(obj, pRoute)
		}

		return printJSON(*cmd, obj)
	}

	for idx, route := range routes {
		if route.Dst == nil {
			cmd.defaultRoute(route, ifaceNames[idx])
		} else {
			cmd.showRoute(route, ifaceNames[idx])
		}
	}
	return nil
}

func (cmd *cmd) filteredRouteList(route *netlink.Route, filterMask uint64, root, match, exact *net.IPNet) ([]netlink.Route, []string, error) {
	var matchedRoutes []netlink.Route
	var ifaceNames []string

	routes, err := netlink.RouteListFiltered(cmd.Family, route, filterMask)
	if err != nil {
		return matchedRoutes, nil, err
	}

	if root == nil && match == nil && exact == nil {
		matchedRoutes = routes
	} else {
		matchedRoutes, err = matchRoutes(routes, root, match, exact)
		if err != nil {
			return matchedRoutes, nil, err
		}
	}

	for _, route := range matchedRoutes {
		link, err := cmd.handle.LinkByIndex(route.LinkIndex)
		if err != nil {
			return matchedRoutes, nil, err
		}

		ifaceNames = append(ifaceNames, link.Attrs().Name)
	}

	return matchedRoutes, ifaceNames, nil
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

func (cmd *cmd) showRoutesForAddress(addr net.IP, options *netlink.RouteGetOptions) error {
	routes, err := cmd.handle.RouteGetWithOptions(addr, options)
	if err != nil {
		return err
	}

	for _, route := range routes {
		link, err := cmd.handle.LinkByIndex(route.LinkIndex)
		if err != nil {
			return err
		}
		if route.Dst == nil {
			cmd.defaultRoute(route, link.Attrs().Name)
		} else {
			cmd.showRoute(route, link.Attrs().Name)
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
	defaultFmt   = "%vdefault via %v dev %s proto %s metric %d\n"
	routeFmt     = "%v%v dev %s proto %s scope %s src %s metric %d\n"
	route6Fmt    = "%v%s dev %s proto %s metric %d\n"
	routeVia6Fmt = "%v%s via %s dev %s proto %s metric %d\n"
)

func (cmd *cmd) defaultRoute(r netlink.Route, name string) {
	gw := r.Gw

	var proto string

	if !cmd.Opts.Numeric {
		proto = rtProto[int(r.Protocol)]
	} else {
		proto = fmt.Sprintf("%d", r.Protocol)
	}

	metric := r.Priority

	var detail string

	if cmd.Opts.Details {
		detail = routeTypeToString(r.Type) + " "
	}

	fmt.Fprintf(cmd.Out, defaultFmt, detail, gw, name, proto, metric)
}

func (cmd *cmd) showRoute(r netlink.Route, name string) {
	switch cmd.Family {
	// print only ipv4 per default
	case netlink.FAMILY_ALL, netlink.FAMILY_V4:
		if r.Dst.IP.To4() == nil {
			return
		}

		cmd.printIPv4Route(r, name)

	case netlink.FAMILY_V6:
		if r.Dst.IP.To4() != nil {
			return
		}

		cmd.printIPv6Route(r, name)
	}
}

func (cmd *cmd) printIPv4Route(r netlink.Route, name string) {
	dest := r.Dst.String()

	var proto, scope string

	if !cmd.Opts.Numeric {
		proto = rtProto[int(r.Protocol)]
		scope = addrScopeStr(r.Scope)
	} else {
		proto = fmt.Sprintf("%d", r.Protocol)
		scope = fmt.Sprintf("%d", r.Scope)
	}

	src := r.Src
	metric := r.Priority

	var detail string

	if cmd.Opts.Details {
		detail = routeTypeToString(r.Type) + " "
	}

	fmt.Fprintf(cmd.Out, routeFmt, detail, dest, name, proto, scope, src, metric)
}

func (cmd *cmd) printIPv6Route(r netlink.Route, name string) {
	dest := r.Dst

	var proto string

	if !cmd.Opts.Numeric {
		proto = rtProto[int(r.Protocol)]
	} else {
		proto = fmt.Sprintf("%d", r.Protocol)
	}

	metric := r.Priority

	var detail string

	if cmd.Opts.Details {
		detail = routeTypeToString(r.Type) + " "
	}

	if r.Gw != nil {
		gw := r.Gw
		fmt.Fprintf(cmd.Out, routeVia6Fmt, detail, dest, gw, name, proto, metric)
	} else {
		fmt.Fprintf(cmd.Out, route6Fmt, detail, dest, name, proto, metric)
	}
}

func (cmd *cmd) routeGet() error {
	addr, err := cmd.parseAddress()
	if err != nil {
		return err
	}

	options, err := cmd.parseRouteGet()
	if err != nil {
		return err
	}

	return cmd.showRoutesForAddress(addr, options)
}

func (cmd *cmd) parseRouteGet() (*netlink.RouteGetOptions, error) {
	var opts netlink.RouteGetOptions
	for cmd.tokenRemains() {
		switch cmd.nextToken("from", "iif", "oif", "vrf") {
		case "oif":
			opts.Oif = cmd.nextToken("OIF")
		case "iif":
			opts.Iif = cmd.nextToken("IIF")

		case "vrf":
			opts.VrfName = cmd.nextToken("VRF_NAME")
		case "from":
			opts.SrcAddr = net.ParseIP(cmd.nextToken("ADDRESS"))
		default:
			return nil, cmd.usage()
		}
	}

	return &opts, nil
}

func (cmd *cmd) route() error {
	if !cmd.tokenRemains() {
		return cmd.showAllRoutes()
	}

	switch cmd.findPrefix("show", "add", "append", "replace", "del", "list", "flush", "get", "help") {
	case "add":
		return cmd.routeAdd()
	case "append":
		return cmd.routeAppend()
	case "replace":
		return cmd.routeReplace()
	case "del":
		return cmd.routeDel()
	case "show", "list":
		return cmd.routeShow()
	case "flush":
		return cmd.routeFlush()
	case "get":
		return cmd.routeGet()
	case "help":
		fmt.Fprint(cmd.Out, routeHelp)
		return nil
	}
	return cmd.usage()
}
