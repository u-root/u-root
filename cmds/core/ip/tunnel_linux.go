// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"net"
	"reflect"

	"github.com/vishvananda/netlink"
)

const (
	tunnelHelp = `Usage: ip tunnel { add | change | del | show | prl | 6rd } [ NAME ]
        [ mode { gre | ipip | sit | vti } ]
        [ remote ADDR ] [ local ADDR ]
        [ [i|o]seq ] [ [i|o]key KEY ] [ [i|o]csum ]
        [ prl-default ADDR ] [ prl-nodefault ADDR ] [ prl-delete ADDR ]
        [ 6rd-prefix ADDR ] [ 6rd-relay_prefix ADDR ] [ 6rd-reset ]
        [ ttl TTL ] [ tos TOS ] [ [no]pmtudisc ] [ dev PHYS_DEV ]

Where: NAME := STRING
       ADDR := { IP_ADDRESS | any }
       TOS  := { 1..255 }
       TTL  := { 1..255 | inherit }
       KEY  := { NUMBER }
`
)

var (
	filterMap      = map[string][]string{"gre": {"gre", "ip6gre"}, "ip6gre": {"ip6gre"}, "ipip": {"ipip", "ip6tln"}, "ip6tln": {"ip6tln"}, "vti": {"vti", "vti6"}, "vti6": {"vti6"}, "sit": {"sit"}}
	allTunnelTypes = []string{"gre", "ipip", "ip6tln", "ip6gre", "vti", "vti6", "sit"}
)

func tunnel(w io.Writer) error {
	cursor++
	if len(arg[cursor:]) == 0 {
		return routeShow(w)
	}

	whatIWant = []string{"add", "del", "show"}
	var c string

	switch c = findPrefix(arg[cursor], whatIWant); c {
	case "show", "add", "del":
		options, err := parseTunnel()
		if err != nil {
			return err
		}

		switch c {
		case "add":
			return tunnelAdd(options)
		case "del":
			return tunnelDelete(options)
		case "show":
			return showTunnels(w, options)
		}
	}
	return usage()
}

type options struct {
	// populated when adding a tunnel
	mode string
	// populated when showing tunnelss
	modes  []string
	dev    string
	name   string
	remote string
	local  string
	iKey   int
	oKey   int
	ttl    int
	tos    int
}

func defaultOptions() options {
	return options{
		modes: []string{},
		iKey:  -1,
		oKey:  -1,
		ttl:   -1,
		tos:   -1,
	}
}

func parseTunnel() (*options, error) {
	options := defaultOptions()
	cursor++

	whatIWant = []string{"name", "mode", "remote", "local", "ttl", "tos", "ikey", "okey", "dev"}
	for cursor < len(arg) {
		switch arg[cursor] {
		case "mode":
			cursor++
			whatIWant = []string{"gre", "ip6gre", "ipip", "ip6tln", "vti", "vti6", "sit"}

			switch arg[cursor] {
			case "gre", "ip6gre", "ipip", "ip6tln", "vti", "vti6", "sit":
				options.mode = arg[cursor]
				options.modes = append(options.modes, filterMap[arg[cursor]]...)
			default:
				return nil, fmt.Errorf("invalid mode %s", arg[cursor])
			}
		case "remote":
			cursor++
			whatIWant = []string{"IP_ADDRESS, any"}

			options.remote = arg[cursor]
		case "local":
			cursor++
			whatIWant = []string{"IP_ADDRESS, any"}

			options.local = arg[cursor]
		case "ttl":
			whatIWant = []string{"TTL (0...255) | inherit"}
			if arg[cursor+1] == "inherit" {
				cursor++
				options.ttl = 0
				continue
			}

			ttl, err := parseUint8("TTL (0...255) | inherit")
			if err != nil {
				return nil, err
			}

			options.ttl = int(ttl)
		case "tos":
			tos, err := parseUint8("TOS (0...255)")
			if err != nil {
				return nil, err
			}

			options.tos = int(tos)
		case "ikey":
			cursor++
			whatIWant = []string{"key"}

			iKey, err := parseUint16("KEY")
			if err != nil {
				return nil, err
			}

			options.iKey = int(iKey)
		case "okey":
			cursor++
			whatIWant = []string{"key"}

			oKey, err := parseUint16("KEY")
			if err != nil {
				return nil, err
			}

			options.oKey = int(oKey)
		case "dev":
			cursor++
			whatIWant = []string{"PHYS_DEV"}

			options.dev = arg[cursor]
		default:
			options.name = arg[cursor]
		}

		cursor++
	}

	if reflect.DeepEqual(options.modes, []string{}) {
		options.modes = allTunnelTypes
	}

	return &options, nil
}

func showTunnels(w io.Writer, op *options) error {
	links, err := netlink.LinkList()
	if err != nil {
		return fmt.Errorf("failed to list interfaces: %v", err)
	}

	var tunnels []netlink.Link

	for _, l := range links {
		found := false
		for _, t := range op.modes {
			if l.Type() == t {
				found = true
			}
		}

		if !found {
			continue
		}

		if op.name != "" && l.Attrs().Name != op.name {
			continue
		}

		if !equalRemotes(l, op.remote) {
			continue
		}

		if !equalLocals(l, op.local) {
			continue
		}

		if op.dev != "" && l.Attrs().Name != op.dev {
			continue
		}

		if !equalTOS(l, op.tos) {
			continue
		}

		if !equalTTL(l, op.ttl) {
			continue
		}

		if !equalIKey(l, op.iKey) {
			continue
		}

		if !equalOKey(l, op.oKey) {
			continue
		}

		tunnels = append(tunnels, l)
	}

	return printTunnels(w, tunnels)
}

func printTunnels(w io.Writer, tunnels []netlink.Link) error {
	var (
		remote  string
		local   string
		ttl     string
		tlnType string
	)

	for _, t := range tunnels {
		switch v := t.(type) {
		case *netlink.Gretun:
			remote = v.Remote.String()
			local = v.Local.String()
			tlnType = "gre"
			ttl = fmt.Sprintf("ttl %d", v.Ttl)
		case *netlink.Iptun:
			remote = v.Remote.String()
			local = v.Local.String()
			tlnType = "ip"
			ttl = fmt.Sprintf("ttl %d", v.Ttl)
		case *netlink.Ip6tnl:
			remote = v.Remote.String()
			local = v.Local.String()
			tlnType = "ipv6"
			ttl = fmt.Sprintf("ttl %d", v.Ttl)
		case *netlink.Vti:
			remote = v.Remote.String()
			local = v.Local.String()
			tlnType = "ip"
		case *netlink.Sittun:
			remote = v.Remote.String()
			local = v.Local.String()
			tlnType = "ipv6"
			ttl = fmt.Sprintf("ttl %d", v.Ttl)
		default:
			return fmt.Errorf("unsupported tunnel type %s", t.Type())
		}

		if remote == "0.0.0.0" || remote == "::" {
			remote = "any"
		}

		if local == "0.0.0.0" || local == "::" {
			local = "any"
		}

		if ttl == "ttl 0" || ttl == "ttl 255" {
			ttl = "ttl inherit"
		}

		fmt.Fprintf(w, "%s %s/ip remote %s local %s %s\n", t.Attrs().Name, tlnType, remote, local, ttl)

	}

	return nil
}

// Function to check if the remote IP matches for the given link.
func equalRemotes(l netlink.Link, remote string) bool {
	if remote == "" || remote == "any" {
		return true
	}

	remoteIP := net.ParseIP(remote)
	switch v := l.(type) {
	case *netlink.Gretun:
		return remoteIP.Equal(v.Remote)
	case *netlink.Iptun:
		return remoteIP.Equal(v.Remote)
	case *netlink.Ip6tnl:
		return remoteIP.Equal(v.Remote)
	case *netlink.Vti:
		return remoteIP.Equal(v.Remote)
	case *netlink.Sittun:
		return remoteIP.Equal(v.Remote)
	default:
		return false
	}
}

// Function to check if the local IP matches for the given tunnel.
func equalLocals(l netlink.Link, local string) bool {
	if local == "" || local == "any" {
		return true
	}

	localIP := net.ParseIP(local)

	switch v := l.(type) {
	case *netlink.Gretun:
		return localIP.Equal(v.Local)
	case *netlink.Iptun:
		return localIP.Equal(v.Local)
	case *netlink.Ip6tnl:
		return localIP.Equal(v.Local)
	case *netlink.Vti:
		return localIP.Equal(v.Local)
	case *netlink.Sittun:
		return localIP.Equal(v.Local)
	default:
		return false
	}
}

func equalTTL(l netlink.Link, ttl int) bool {
	if ttl == -1 || ttl == 0 || ttl == 255 {
		return true
	}

	switch v := l.(type) {
	case *netlink.Gretun:
		return ttl == int(v.Ttl)
	case *netlink.Iptun:
		return ttl == int(v.Ttl)
	case *netlink.Ip6tnl:
		return ttl == int(v.Ttl)
	// vti does not have TTL field
	case *netlink.Vti:
		return true
	case *netlink.Sittun:
		return ttl == int(v.Ttl)
	default:
		return false
	}
}

func equalTOS(l netlink.Link, tos int) bool {
	if tos == -1 {
		return true
	}

	switch v := l.(type) {
	case *netlink.Gretun:
		return tos == int(v.Tos)
	case *netlink.Iptun:
		return tos == int(v.Tos)
	case *netlink.Ip6tnl:
		return tos == int(v.Tos)
	// vti does not have TOS field
	case *netlink.Vti:
		return true
	case *netlink.Sittun:
		return tos == int(v.Tos)
	default:
		return false
	}
}

func equalIKey(l netlink.Link, iKey int) bool {
	if iKey == -1 {
		return true
	}

	switch v := l.(type) {
	case *netlink.Gretun:
		return iKey == int(v.IKey)
	case *netlink.Vti:
		return iKey == int(v.IKey)
	default:
		return true
	}
}

func equalOKey(l netlink.Link, oKey int) bool {
	if oKey == -1 {
		return true
	}

	switch v := l.(type) {
	case *netlink.Gretun:
		return oKey == int(v.OKey)
	case *netlink.Vti:
		return oKey == int(v.OKey)
	default:
		return true
	}
}

func tunnelAdd(op *options) error {
	if op.mode == "" {
		return fmt.Errorf("tunnel mode is required")
	}

	if op.name == "" {
		switch op.mode {
		case "gre", "ip6gre":
			op.name = "gre0"
		case "ipip":
			op.name = "tuln0"
		case "ip6tln":
			op.name = "ip6tnl0"
		case "vti", "vti6":
			op.name = "ip_vti0"
		case "sit":
			op.name = "sit0"
		}
	}

	if op.iKey < 0 {
		op.iKey = 0
	}

	if op.oKey < 0 {
		op.oKey = 0
	}

	if op.ttl < 0 {
		op.ttl = 0
	}

	if op.tos < 0 {
		op.tos = 0
	}

	var link netlink.Link

	switch op.mode {
	case "gre", "ip6gre":
		link = &netlink.Gretun{
			LinkAttrs: netlink.LinkAttrs{
				Name: op.name,
			},
			Remote: net.ParseIP(op.remote),
			Local:  net.ParseIP(op.local),
			Ttl:    uint8(op.ttl),
			Tos:    uint8(op.tos),
			IKey:   uint32(op.iKey),
			OKey:   uint32(op.oKey),
		}
	case "ipip":
		link = &netlink.Iptun{
			LinkAttrs: netlink.LinkAttrs{
				Name: op.name,
			},
			Remote: net.ParseIP(op.remote),
			Local:  net.ParseIP(op.local),
			Ttl:    uint8(op.ttl),
			Tos:    uint8(op.tos),
		}
	case "ip6tln":
		link = &netlink.Ip6tnl{
			LinkAttrs: netlink.LinkAttrs{
				Name: op.name,
			},
			Remote: net.ParseIP(op.remote),
			Local:  net.ParseIP(op.local),
			Ttl:    uint8(op.ttl),
			Tos:    uint8(op.tos),
		}
	case "vti", "vti6":
		link = &netlink.Vti{
			LinkAttrs: netlink.LinkAttrs{
				Name: op.name,
			},
			Remote: net.ParseIP(op.remote),
			Local:  net.ParseIP(op.local),
			IKey:   uint32(op.iKey),
			OKey:   uint32(op.oKey),
		}
	case "sit":
		link = &netlink.Sittun{
			LinkAttrs: netlink.LinkAttrs{
				Name: op.name,
			},
			Remote: net.ParseIP(op.remote),
			Local:  net.ParseIP(op.local),
			Ttl:    uint8(op.ttl),
			Tos:    uint8(op.tos),
		}
	default:
		return fmt.Errorf("unsupported tunnel type %s", op.mode)
	}

	if err := netlink.LinkAdd(link); err != nil {
		return fmt.Errorf("failed to add tunnel: %v", err)
	}

	return nil
}

func tunnelDelete(op *options) error {
	if op.name == "" {
		return fmt.Errorf("tunnel name is required")
	}

	link, err := netlink.LinkByName(op.name)
	if err != nil {
		return fmt.Errorf("failed to find tunnel %s: %v", op.name, err)
	}

	valid := true
	for _, t := range allTunnelTypes {
		if link.Type() == t {
			valid = true
			break
		}
	}

	if !valid {
		return fmt.Errorf("%s is not a tunnel device", op.name)
	}

	if err := netlink.LinkDel(link); err != nil {
		return fmt.Errorf("failed to delete tunnel %s: %v", op.name, err)
	}

	return nil
}
