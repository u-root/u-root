// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"fmt"
	"net"
	"strconv"

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
             [ scope SCOPE ] [ metric NUMBER ]

INFO_SPEC := NH OPTIONS [ nexthop NH ]...

NH := [ via [ FAMILY ] ADDRESS ] [ dev STRING ] NHFLAGS

FAMILY := [ inet | inet6 | mpls ]

OPTIONS := [ mtu NUMBER ] [ advmss NUMBER ]
           [ rtt TIME ] [ rttvar TIME ] [ reordering NUMBER ]
           [ window NUMBER ] [ cwnd NUMBER ] [ initcwnd NUMBER ]
           [ ssthresh NUMBER ] [ realms REALM ] [ src ADDRESS ]
           [ rto_min TIME ] [ hoplimit NUMBER ] [ initrwnd NUMBER ]
           [ features FEATURES ] [ quickack BOOL ] [ congctl NAME ]
		   [ fastopen_no_cookie BOOL ]

TYPE := { unicast | local | broadcast | multicast | throw |
          unreachable | prohibit | blackhole | nat }

TABLE_ID := [ local | main | default | NUMBER ]

SCOPE := [ host | link | global | NUMBER ]

NHFLAGS := [ onlink | pervasive ]

RTPROTO := [ kernel | boot | static | NUMBER ]

BOOL := [1|0]
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

func routeTypeStr(routeType int) string {
	for key, value := range routeTypes {
		if value == routeType {
			return key
		}
	}
	return "unknown"
}

// route is the entry point for 'ip route' command.
func (cmd *cmd) route() error {
	if !cmd.tokenRemains() {
		return cmd.routeShow()
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
	default:
		return cmd.usage()
	}
}

// routeAdd performs the 'ip route add' command.
func (cmd *cmd) routeAdd() error {
	route, err := cmd.parseRouteAddAppendReplaceDel(defaultLinkIdxResolver)
	if err != nil {
		return err
	}

	if err := cmd.handle.RouteAdd(route); err != nil {
		return fmt.Errorf("adding route for %s: %w", route.Dst, err)
	}
	return nil
}

// routeAppend performs the 'ip route append' command.
func (cmd *cmd) routeAppend() error {
	route, err := cmd.parseRouteAddAppendReplaceDel(defaultLinkIdxResolver)
	if err != nil {
		return err
	}

	if err := cmd.handle.RouteAppend(route); err != nil {
		return fmt.Errorf("appending route for %s: %w", route.Dst.IP, err)
	}
	return nil
}

// routeReplace performs the 'ip route replace' command.
func (cmd *cmd) routeReplace() error {
	route, err := cmd.parseRouteAddAppendReplaceDel(defaultLinkIdxResolver)
	if err != nil {
		return err
	}

	if err := cmd.handle.RouteReplace(route); err != nil {
		return fmt.Errorf("appending route for %s: %w", route.Dst.IP, err)
	}
	return nil
}

// routeDel performs the 'ip route del' command.
func (cmd *cmd) routeDel() error {
	route, err := cmd.parseRouteAddAppendReplaceDel(defaultLinkIdxResolver)
	if err != nil {
		return err
	}

	if err := cmd.handle.RouteDel(route); err != nil {
		return fmt.Errorf("deleting route for %s: %w", route.Dst.IP, err)
	}
	return nil
}

// defaultLinkIdxResolver resolves the link index by name.
// It is used by all all functions that call parseRouteAddAppendReplaceDel.
// It is a separate function to allow for testing with a mock implementation.
func defaultLinkIdxResolver(name string) (int, error) {
	iface, err := netlink.LinkByName(name)
	if err != nil {
		return 0, fmt.Errorf("link %q not found: %w", name, err)
	}
	return iface.Attrs().Index, nil
}

// parseRouteAddAppendReplaceDel parses the arguments to 'ip route add', 'ip route append',
// 'ip route replace' or 'ip route delete' from the cmdline.
func (cmd *cmd) parseRouteAddAppendReplaceDel(resolveLinkIdxFn func(string) (int, error)) (*netlink.Route, error) {
	if resolveLinkIdxFn == nil {
		return nil, fmt.Errorf("internal parser error: resolveLinkFn is not set")
	}

	var (
		route = netlink.Route{}
		err   error
	)

	// TYPE
	if routeType, ok := probeRouteType(cmd.peekToken("TYPE")); ok {
		cmd.Cursor++
		route.Type = int(routeType)
	}

	// PREFIX
	dest := cmd.nextToken("default", "PREFIX (in CIDR notation)")
	if dest == "default" {
		route.Dst = nil
	} else {
		_, route.Dst, err = net.ParseCIDR(dest)
		if err != nil {
			return nil, err
		}
	}

	// Eventually further specification of the destination address
	// (e.g. tos, table, proto, scope, metric, etc.)
	// Otherwise, NH (via or dev)
LOOP:
	for cmd.tokenRemains() {
		token := cmd.nextToken("dev", "via", "tos", "table", "proto", "scope", "metric")
		switch token {
		case "via", "dev":
			cmd.Cursor--
			break LOOP // handled below
		case "tos":
			token = cmd.nextToken("TOS")
			route.Tos, err = parseTOS(token)
			if err != nil {
				return nil, err
			}
		case "table":
			token = cmd.nextToken("local", "main", "default", "all", "NUMBER")
			route.Table, err = parseTable(token)
			if err != nil {
				return nil, err
			}
		case "proto":
			token = cmd.nextToken("kernel", "boot", "static", "NUMBER")
			route.Protocol, err = parseProto(token)
			if err != nil {
				return nil, err
			}
		case "scope":
			token = cmd.nextToken("host", "link", "global", "NUMBER")
			route.Scope, err = parseScope(token)
			if err != nil {
				return nil, err
			}
		case "metric":
			token = cmd.nextToken("NUMBER")
			route.Priority, err = parseMetric(token)
			if err != nil {
				return nil, err
			}
		default:
			return nil, cmd.usage()
		}
	}

	// NH (next hop)

	// "dev" or "via" must be present otherwise the loop above would not have exited. So there will be at least one
	// NexthopInfo item in the route.MultiPath slice. The NexthopInfo will be populated with the information from
	// the "dev" or "via" token in the loop below.
	// Furhter NexthopInfo items will be appended if "nexthop" token is found in the loop below.

	nextHopIdx := -1

	for cmd.tokenRemains() {
		token := cmd.nextToken("dev", "via", "nexthop", "mtu", "advmss", "rtt", "rttvar", "reordering", "window", "cwnd", "initcwnd", "ssthresh", "realms", "src", "rto_min", "hoplimit", "initrwnd", "features", "quickack", "congctl", "fastopen_no_cookie")
		switch token {
		case "dev":
			token = cmd.nextToken("DEVICE")
			idx, err := resolveLinkIdxFn(token)
			if err != nil {
				return nil, err
			}
			if nextHopIdx < 0 { // first nexthop item is stored in the top-level route fields
				route.LinkIndex = idx
				// NHFLAGS
				if cmd.tokenRemains() {
					if nhFlags, ok := probeNHFlags(cmd.peekToken("NHFLAGS")); ok {
						cmd.Cursor++
						route.Flags = int(nhFlags)
					}
				}
			} else { // next hop item is a nexthop item and is stored in the MultiPath slice
				route.MultiPath[nextHopIdx].LinkIndex = idx
				// NHFLAGS
				if cmd.tokenRemains() {
					if nhFlags, ok := probeNHFlags(cmd.peekToken("NHFLAGS")); ok {
						cmd.Cursor++
						route.MultiPath[nextHopIdx].Flags = int(nhFlags)
					}
				}
			}

		case "via":
			if family, ok := probeFamily(cmd.peekToken("FAMILY")); ok {
				cmd.Cursor++
				route.Family = family
			}
			token = cmd.nextToken("ADDRESS")
			gw := net.ParseIP(token)
			if gw == nil {
				return nil, fmt.Errorf("invalid address: %s", token)
			}
			if nextHopIdx < 0 { // first nexthop item is stored in the top-level route fields
				route.Gw = gw
			} else { // next hop item is a nexthop item and is stored in the MultiPath slice
				route.MultiPath[nextHopIdx].Gw = gw
			}
		case "nexthop":
			route.MultiPath = append(route.MultiPath, &netlink.NexthopInfo{})
			nextHopIdx++
			continue
		case "mtu":
			token = cmd.nextToken("NUMBER")
			num, err := strconv.ParseInt(token, 10, 0)
			if err != nil {
				return nil, err
			}
			route.MTU = int(num)
		case "advmss":
			token = cmd.nextToken("NUMBER")
			num, err := strconv.ParseInt(token, 10, 0)
			if err != nil {
				return nil, err
			}
			route.AdvMSS = int(num)
		case "rtt":
			token = cmd.nextToken("TIME")
			route.Rtt, err = parseTime(token)
			if err != nil {
				return nil, err
			}
		case "rttvar":
			token = cmd.nextToken("TIME")
			route.RttVar, err = parseTime(token)
			if err != nil {
				return nil, err
			}
		case "reordering":
			token = cmd.nextToken("NUMBER")
			reordering, err := strconv.ParseInt(token, 10, 0)
			if err != nil {
				return nil, err
			}
			route.Reordering = int(reordering)
		case "window":
			token = cmd.nextToken("NUMBER")
			window, err := strconv.ParseInt(token, 10, 0)
			if err != nil {
				return nil, err
			}
			route.Window = int(window)
		case "cwnd":
			token = cmd.nextToken("NUMBER")
			cwnd, err := strconv.ParseInt(token, 10, 0)
			if err != nil {
				return nil, err
			}
			route.Cwnd = int(cwnd)
		case "initcwnd":
			token = cmd.nextToken("NUMBER")
			initcwnd, err := strconv.ParseInt(token, 10, 0)
			if err != nil {
				return nil, err
			}
			route.InitCwnd = int(initcwnd)
		case "ssthresh":
			token = cmd.nextToken("NUMBER")
			sshresh, err := strconv.ParseInt(token, 10, 0)
			if err != nil {
				return nil, err
			}
			route.Ssthresh = int(sshresh)
		case "realms":
			token = cmd.nextToken("NUMBER")
			realms, err := strconv.ParseInt(token, 10, 0)
			if err != nil {
				return nil, err
			}
			route.Realm = int(realms)
		case "src":
			token = cmd.nextToken("ADDRESS")
			route.Src = net.ParseIP(token)
			if route.Src == nil {
				return nil, fmt.Errorf("invalid src address: %s", token)
			}
		case "rto_min":
			token = cmd.nextToken("NUMBER")
			route.RtoMin, err = parseTime(token)
			if err != nil {
				return nil, err
			}
		case "hoplimit":
			token = cmd.nextToken("NUMBER")
			hoplimit, err := strconv.ParseInt(token, 10, 0)
			if err != nil {
				return nil, err
			}
			route.Hoplimit = int(hoplimit)
		case "initrwnd":
			token = cmd.nextToken("NUMBER")
			initrwnd, err := strconv.ParseInt(token, 10, 0)
			if err != nil {
				return nil, err
			}
			route.InitRwnd = int(initrwnd)
		case "features":
			token = cmd.nextToken("NUMBER")
			features, err := strconv.ParseInt(token, 10, 0)
			if err != nil {
				return nil, err
			}
			route.Features = int(features)
		case "quickack":
			token = cmd.nextToken("BOOL 0 | 1")
			// according to the SYNAPSYS documentation, this is a boolean value taking 0 or 1,
			// so use ParseBool to prevent other numbers from being accepted
			quickack, err := strconv.ParseBool(token)
			if err != nil {
				return nil, err
			}
			// netlink.Route.QuickACK is an int, so convert the boolean back to 0 or 1
			if quickack {
				route.QuickACK = 1
			}
		case "congctl":
			token = cmd.nextToken("NAME")
			route.Congctl = token
		case "fastopen_no_cookie":
			token = cmd.nextToken("BOOL 0 | 1")
			// according to the SYNAPSYS documentation, this is a boolean value taking 0 or 1,
			// so use ParseBool to prevent other numbers from being accepted
			fastopenNoCookie, err := strconv.ParseBool(token)
			if err != nil {
				return nil, err
			}
			// netlink.Route.FastOpenNoCookie is an int, so convert the boolean back to 0 or 1
			if fastopenNoCookie {
				route.FastOpenNoCookie = 1
			}
		default:
			return nil, cmd.usage()
		}
	}

	return &route, nil
}

func probeRouteType(token string) (uint32, bool) {
	// The netlink package itself does not define these constants, but it uses e.g.
	// unix.RTN_UNICAST in the code as default in some cases. So netlink.Route.Type value is
	// set using the unix package constants. As the unix package does not document these values properly,
	// values are chosen with reference to https://manpages.debian.org/testing/manpages/rtnetlink.7.en.html.
	switch token {
	case "unicast":
		return unix.RTN_UNICAST, true
	case "local":
		return unix.RTN_LOCAL, true
	case "broadcast":
		return unix.RTN_BROADCAST, true
	case "multicast":
		return unix.RTN_MULTICAST, true
	case "throw":
		return unix.RTN_THROW, true
	case "unreachable":
		return unix.RTN_UNREACHABLE, true
	case "prohibit":
		return unix.RTN_PROHIBIT, true
	case "blackhole":
		return unix.RTN_BLACKHOLE, true
	case "nat":
		return unix.RTN_NAT, true
	default:
		return 0, false
	}
}

func probeFamily(token string) (int, bool) {
	switch token {
	case "inet":
		return netlink.FAMILY_V4, true
	case "inet6":
		return netlink.FAMILY_V6, true
	case "mpls":
		return netlink.FAMILY_MPLS, true
	// case "bridge": // netlink.FAMILY_BRIDGE is not defined in the netlink package
	// 	return ??, true
	// case "link":  // netlink.FAMILY_LINK is not defined in the netlink package
	// 	return ??, true
	default:
		return 0, false
	}
}

func probeNHFlags(token string) (netlink.NextHopFlag, bool) {
	switch token {
	case "onlink":
		return netlink.FLAG_ONLINK, true
	case "pervasive":
		return netlink.FLAG_PERVASIVE, true
	default:
		return 0, false
	}
}

func parseTOS(token string) (int, error) {
	// According to the documentation, the TOS value is a 8 bit hexadecimal number.
	// https://manpages.debian.org/bookworm/iproute2/ip-route.8.en.html#tos
	n, err := strconv.ParseInt(token, 16, 8)
	if err != nil {
		// However, also allow decimal numbers.
		n, err = strconv.ParseInt(token, 10, 8)
		if err != nil {
			return 0, fmt.Errorf("invalid TOS value %q: %w", token, err)
		}
	}
	return int(n), nil
}

func parseTable(token string) (int, error) {
	// The netlink package itself does not define these constants, but it uses e.g.
	// unix.RT_TABLE_UNSPEC in the code . So netlink.Route.Table value is set using
	// the unix package constants. As the unix package does not document these values properly,
	// values are chosen with reference to https://manpages.debian.org/testing/manpages/rtnetlink.7.en.html.
	switch token {
	case "local":
		return unix.RT_TABLE_LOCAL, nil
	case "main":
		return unix.RT_TABLE_MAIN, nil
	case "default":
		return unix.RT_TABLE_DEFAULT, nil
	default:
		n, err := strconv.ParseInt(token, 10, 0)
		if err != nil {
			return 0, fmt.Errorf("invalid table id %q: %w", token, err)
		}
		return int(n), nil
	}
}

func parseProto(token string) (netlink.RouteProtocol, error) {
	// The netlink package itself does not define these constants, but it uses the respective unix constants e.g.
	// in netlink.RouteProtocol.String().
	switch token {
	case "kernel":
		return unix.RTPROT_KERNEL, nil
	case "boot":
		return unix.RTPROT_BOOT, nil
	case "static":
		return unix.RTPROT_STATIC, nil
	default:
		n, err := strconv.ParseInt(token, 10, 0)
		if err != nil {
			return 0, fmt.Errorf("invalid protocol %q: %w", token, err)
		}
		return netlink.RouteProtocol(n), nil
	}
}

func parseScope(token string) (netlink.Scope, error) {
	// The netlink package itself does not define these constants, but it uses the respective unix constants e.g.
	// in netlink.Scope.String().
	switch token {
	case "host":
		return netlink.SCOPE_HOST, nil
	case "link":
		return netlink.SCOPE_LINK, nil
	case "global":
		return netlink.SCOPE_UNIVERSE, nil
	default:
		n, err := strconv.ParseInt(token, 10, 8)
		if err != nil {
			return 0, fmt.Errorf("invalid scope %q: %w", token, err)
		}
		return netlink.Scope(n), nil
	}
}

func parseMetric(token string) (int, error) {
	n, err := strconv.ParseInt(token, 10, 0)
	if err != nil {
		return 0, fmt.Errorf("invalid metric %q: %w", token, err)
	}
	return int(n), nil
}

func parseTime(token string) (int, error) {
	// val, err := time.ParseDuration(token)
	// var parseError *time.ParseError
	// if errors.As(err, &parseError) {
	// 	raw, err := strconv.ParseInt(token, 10, 0)
	// 	if err != nil {
	// 		return 0, err
	// 	}
	// 	return int(raw), nil
	// } else if err != nil {
	// 	return 0, err
	// }

	// The man page for iproute2 states:
	// "If no suffix is specified the units are raw values passed directly to
	// the routing code to maintain compatibility with previous releases.
	// Otherwise if a suffix of s, sec or secs is used to specify seconds and
	// ms, msec or msecs to specify milliseconds."
	// https://manpages.debian.org/bookworm/iproute2/ip-route.8.en.html#rtt
	//
	// So actually the above code would be smart. However the documentation misses
	// the importatnt part of what time unit the raw value is in. Casting time.Duration
	// to int will give the value in nanoseconds. Not sure if this is the correct.
	// Therefore only support the compatibility mode for now.

	val, err := strconv.ParseInt(token, 10, 0)
	if err != nil {
		return 0, err
	}

	return int(val), nil
}

// routeShow performs the 'ip route show' command.
func (cmd *cmd) routeShow() error {
	filter, filterMask, root, match, exact, err := cmd.parseRouteShowListFlush()
	if err != nil {
		return err
	}

	routeList, ifaceNames, err := cmd.filteredRouteList(filter, filterMask, root, match, exact)
	if err != nil {
		return err
	}

	return cmd.printRoutes(routeList, ifaceNames)
}

// routeFlush performs the 'ip route flush' command.
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

// parseRouteShowListFlush parses the arguments to 'ip route show', 'ip route list' or 'ip route flush' from the cmdline.
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

// RouteJSON represents a route entry for JSON output format.
type RouteJSON struct {
	Dst      string   `json:"dst"`
	Dev      string   `json:"dev"`
	Protocol string   `json:"protocol"`
	Scope    string   `json:"scope"`
	PrefSrc  string   `json:"prefsrc"`
	Flags    []string `json:"flags,omitempty"`
}

// printRoutes prints the routes in the system.
func (cmd *cmd) printRoutes(routes []netlink.Route, ifaceNames []string) error {
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
			cmd.printDefaultRoute(route, ifaceNames[idx])
		} else {
			cmd.printRoute(route, ifaceNames[idx])
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

func (cmd *cmd) printDefaultRoute(r netlink.Route, name string) {
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
		detail = routeTypeStr(r.Type) + " "
	}

	fmt.Fprintf(cmd.Out, defaultFmt, detail, gw, name, proto, metric)
}

func (cmd *cmd) printRoute(r netlink.Route, name string) {
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
		detail = routeTypeStr(r.Type) + " "
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
		detail = routeTypeStr(r.Type) + " "
	}

	if r.Gw != nil {
		gw := r.Gw
		fmt.Fprintf(cmd.Out, routeVia6Fmt, detail, dest, gw, name, proto, metric)
	} else {
		fmt.Fprintf(cmd.Out, route6Fmt, detail, dest, name, proto, metric)
	}
}

func addrScopeStr(scope netlink.Scope) string {
	switch scope {
	case netlink.SCOPE_UNIVERSE:
		return "global"
	default:
		return scope.String()
	}
}

// routeGet is the entry point for 'ip route get' command.
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

// parseRouteGet parses the arguments to 'ip route get' from the comdline.
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
			cmd.printDefaultRoute(route, link.Attrs().Name)
		} else {
			cmd.printRoute(route, link.Attrs().Name)
		}
	}
	return nil
}
