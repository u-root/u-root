// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"

	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

func routeAdddefault(w io.Writer) error {
	nh, nhval, err := parseNextHop()
	if err != nil {
		return err
	}
	// TODO: NHFLAGS.
	l, err := parseDeviceName()
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
		addr, err := netlink.ParseAddr(ns)
		if err != nil {
			return usage()
		}
		d, err := parseDeviceName()
		if err != nil {
			return usage()
		}
		r := &netlink.Route{LinkIndex: d.Attrs().Index, Dst: addr.IPNet}
		if err := netlink.RouteAdd(r); err != nil {
			return fmt.Errorf("error adding route %s -> %s: %v", addr, d.Attrs().Name, err)
		}
		return nil
	}
}

func routeDel() error {
	cursor++
	addr, err := netlink.ParseAddr(arg[cursor])
	if err != nil {
		return usage()
	}
	d, err := parseDeviceName()
	if err != nil {
		return usage()
	}
	r := &netlink.Route{LinkIndex: d.Attrs().Index, Dst: addr.IPNet}
	if err := netlink.RouteDel(r); err != nil {
		return fmt.Errorf("error adding route %s -> %s: %v", addr, d.Attrs().Name, err)
	}
	return nil
}

func routeShow(w io.Writer) error {
	return showRoutes(w, *inet6)
}

func showRoutes(w io.Writer, inet6 bool) error {
	var f int
	if inet6 {
		f = netlink.FAMILY_V6
	} else {
		f = netlink.FAMILY_V4
	}

	routes, err := netlink.RouteList(nil, f)
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

func route(w io.Writer) error {
	cursor++
	if len(arg[cursor:]) == 0 {
		return routeShow(w)
	}

	whatIWant = []string{"show", "add", "del"}
	switch findPrefix(arg[cursor], whatIWant) {
	case "add":
		return routeAdd(w)
	case "del":
		return routeDel()
	case "show":
		return routeShow(w)
	}
	return usage()
}
