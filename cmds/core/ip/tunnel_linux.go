// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"fmt"
	"net"
	"strconv"

	"github.com/vishvananda/netlink"
)

const (
	tunnelHelp = `Usage: ip tunnel show [ NAME ] 
	[ mode { gre | ipip | sit | vti } ]
	[ remote ADDR ] [ local ADDR ] [ [i|o]key KEY ]
	[ ttl TTL ] [ tos TOS ] [ dev PHYS_DEV ]

       ip tunnel { add | del } [ NAME ]
	[ mode { gre | ipip | sit | vti } ]
	[ remote ADDR ] [ local ADDR ] [ [i|o]key KEY ]
	[ ttl TTL ] [ tos TOS ] 

Where: NAME := STRING
	   ADDR := { IP_ADDRESS | any }
	   TOS  := { 1..255 }
	   TTL  := { 1..255 | inherit }
	   KEY  := { NUMBER }
	   PHYS_DEV := STRING
`
)

var (
	filterMap      = map[string][]string{"gre": {"gre", "ip6gre"}, "ip6gre": {"ip6gre"}, "ipip": {"ipip", "ip6tln"}, "ip6tln": {"ip6tln"}, "vti": {"vti", "vti6"}, "vti6": {"vti6"}, "sit": {"sit"}}
	allTunnelTypes = []string{"gre", "ipip", "ip6tln", "ip6gre", "vti", "vti6", "sit"}
)

func (cmd *cmd) tunnel() error {
	if !cmd.tokenRemains() {
		return cmd.showAllTunnels()
	}

	var c string

	switch c = cmd.findPrefix("add", "delete", "show", "help"); c {
	case "show", "add", "delete":
		options, err := cmd.parseTunnel()
		if err != nil {
			return err
		}

		switch c {
		case "add":
			return cmd.tunnelAdd(options)
		case "delete":
			return cmd.tunnelDelete(options)
		case "show":
			return cmd.showTunnels(options)

		}
	case "help":
		fmt.Fprint(cmd.Out, tunnelHelp)
		return nil
	}
	return cmd.usage()
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

func (cmd *cmd) parseTunnel() (*options, error) {
	options := defaultOptions()

	for cmd.tokenRemains() {
		switch cmd.nextToken("name", "mode", "remote", "local", "ttl", "tos", "ikey", "okey", "dev") {
		case "mode":
			token := cmd.nextToken("gre", "ip6gre", "ipip", "ip6tln", "vti", "vti6", "sit")
			switch token {
			case "gre", "ip6gre", "ipip", "ip6tln", "vti", "vti6", "sit":
				options.mode = token
				options.modes = append(options.modes, filterMap[token]...)
			default:
				return nil, fmt.Errorf("invalid mode %s", token)
			}
		case "remote":
			options.remote = cmd.nextToken("IP_ADDRESS, any")
		case "local":
			options.local = cmd.nextToken("IP_ADDRESS, any")
		case "ttl":
			token := cmd.nextToken("0...255", "inherit")
			if token == "inherit" {
				options.ttl = 0
				continue
			}

			ttl, err := strconv.ParseUint(token, 10, 8)
			if err != nil {
				return nil, err
			}

			options.ttl = int(ttl)
		case "tos":
			tos, err := cmd.parseUint8("TOS (0...255)")
			if err != nil {
				return nil, err
			}

			options.tos = int(tos)
		case "ikey":
			iKey, err := cmd.parseUint16("KEY")
			if err != nil {
				return nil, err
			}

			options.iKey = int(iKey)
		case "okey":
			oKey, err := cmd.parseUint16("KEY")
			if err != nil {
				return nil, err
			}

			options.oKey = int(oKey)
		case "key":
			key, err := cmd.parseUint16("KEY")
			if err != nil {
				return nil, err
			}

			options.iKey = int(key)
			options.oKey = int(key)
		case "dev":
			options.dev = cmd.nextToken("PHYS_DEV")
		default:
			options.name = cmd.currentToken()
		}
	}

	if len(options.modes) == 0 {
		options.modes = allTunnelTypes
	}

	return &options, nil
}

func (cmd *cmd) showAllTunnels() error {
	return cmd.showTunnels(&options{modes: allTunnelTypes})
}

func (cmd *cmd) showTunnels(op *options) error {
	links, err := netlink.LinkList()
	if err != nil {
		return fmt.Errorf("failed to list interfaces: %w", err)
	}

	return cmd.printTunnels(filterTunnels(links, op))
}

func filterTunnels(links []netlink.Link, op *options) []netlink.Link {
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

		if op.dev != "" {
			dev, err := netlink.LinkByName(op.dev)
			if err != nil {
				return []netlink.Link{}
			}

			// Check if the tunnel's parent index matches the physical device
			if l.Attrs().ParentIndex != dev.Attrs().Index {
				continue
			}
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

	return tunnels
}

// Add key fields to the Tunnel struct
type Tunnel struct {
	IfName string `json:"ifname"`
	Mode   string `json:"mode"`
	Remote string `json:"remote"`
	Local  string `json:"local"`
	Dev    string `json:"dev,omitempty"`
	TTL    string `json:"ttl,omitempty"`
	TOS    uint8  `json:"tos,omitempty"`
	IKey   uint32 `json:"ikey,omitempty"`
	OKey   uint32 `json:"okey,omitempty"`
}

func (cmd *cmd) printTunnels(tunnels []netlink.Link) error {
	pTunnels := make([]Tunnel, 0, len(tunnels))

	for _, t := range tunnels {
		var tunnel Tunnel
		tunnel.IfName = t.Attrs().Name

		if t.Attrs().ParentIndex != 0 {
			parent, err := cmd.handle.LinkByIndex(t.Attrs().ParentIndex)
			if err != nil {
				return fmt.Errorf("failed to get parent link: %w", err)
			}
			tunnel.Dev = parent.Attrs().Name
		}

		switch v := t.(type) {
		case *netlink.Gretun:
			tunnel.Remote = v.Remote.String()
			tunnel.Local = v.Local.String()
			tunnel.Mode = "gre"
			tunnel.TTL = fmt.Sprintf("%d", v.Ttl)
			tunnel.TOS = v.Tos
			tunnel.IKey = v.IKey
			tunnel.OKey = v.OKey
		case *netlink.Iptun:
			tunnel.Remote = v.Remote.String()
			tunnel.Local = v.Local.String()
			tunnel.Mode = "any"
			tunnel.TTL = fmt.Sprintf("%d", v.Ttl)
			tunnel.TOS = v.Tos
		case *netlink.Ip6tnl:
			tunnel.Remote = v.Remote.String()
			tunnel.Local = v.Local.String()
			tunnel.Mode = "ip6tln"
			tunnel.TTL = fmt.Sprintf("%d", v.Ttl)
			tunnel.TOS = v.Tos
		case *netlink.Vti:
			tunnel.Remote = v.Remote.String()
			tunnel.Local = v.Local.String()
			tunnel.Mode = "ip"
			tunnel.TTL = "inherit"
			tunnel.IKey = v.IKey
			tunnel.OKey = v.OKey
		case *netlink.Sittun:
			tunnel.Remote = v.Remote.String()
			tunnel.Local = v.Local.String()
			tunnel.Mode = "sit"
			tunnel.TTL = fmt.Sprintf("%d", v.Ttl)
			tunnel.TOS = v.Tos

		default:
			return fmt.Errorf("unsupported tunnel type %s", t.Type())
		}

		if tunnel.Remote == "0.0.0.0" || tunnel.Remote == "::" {
			tunnel.Remote = "any"
		}

		if tunnel.Local == "0.0.0.0" || tunnel.Local == "::" {
			tunnel.Local = "any"
		}

		if tunnel.TTL == "0" || tunnel.TTL == "255" {
			tunnel.TTL = "inherit"
		}

		pTunnels = append(pTunnels, tunnel)
	}

	if cmd.Opts.JSON {
		return printJSON(*cmd, pTunnels)
	}

	for _, t := range pTunnels {
		optsString := ""
		if t.Dev != "" {
			optsString = fmt.Sprintf(" dev %s", t.Dev)
		}

		if t.TTL != "" {
			optsString += fmt.Sprintf(" ttl %s", t.TTL)
		}

		if t.TOS != 0 {
			optsString += fmt.Sprintf(" tos 0x%x", t.TOS)
		}

		if t.IKey == t.OKey && t.IKey != 0 {
			optsString += fmt.Sprintf(" key %d", t.IKey)
		} else {
			if t.IKey != 0 {
				optsString += fmt.Sprintf(" ikey %d", t.IKey)
			}
			if t.OKey != 0 {
				optsString += fmt.Sprintf(" okey %d", t.OKey)
			}
		}

		fmt.Fprintf(cmd.Out, "%s: %s/ip remote %s local %s%s\n",
			t.IfName, t.Mode, t.Remote, t.Local, optsString)
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

func normalizeOptsForAddingTunnel(op *options) error {
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

	return nil
}

func (cmd *cmd) tunnelAdd(op *options) error {
	if err := normalizeOptsForAddingTunnel(op); err != nil {
		return err
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

	if err := cmd.handle.LinkAdd(link); err != nil {
		return fmt.Errorf("failed to add tunnel: %w", err)
	}

	return nil
}

func (cmd *cmd) tunnelDelete(op *options) error {
	if op.name == "" {
		return fmt.Errorf("tunnel name is required")
	}

	link, err := cmd.handle.LinkByName(op.name)
	if err != nil {
		return fmt.Errorf("failed to find tunnel %s: %w", op.name, err)
	}

	valid := false
	for _, t := range allTunnelTypes {
		if link.Type() == t {
			valid = true
			break
		}
	}

	if !valid {
		return fmt.Errorf("%s is not a tunnel device", op.name)
	}

	if err := cmd.handle.LinkDel(link); err != nil {
		return fmt.Errorf("failed to delete tunnel %s: %w", op.name, err)
	}

	return nil
}
