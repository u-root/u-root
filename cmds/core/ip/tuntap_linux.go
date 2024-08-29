// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"fmt"
	"math"

	"github.com/vishvananda/netlink"
)

const (
	tuntapHelp = `Usage: ip tuntap { add | del | show | list | lst | help } [ dev | name ] NAME ]
       [ mode { tun | tap } ] [ user USER ] [ group GROUP ]
       [ one_queue ] [ pi ] [ vnet_hdr ] [ multi_queue ]

Where: USER  := { STRING | NUMBER }
       GROUP := { STRING | NUMBER }`
)

func (cmd *cmd) tuntap() error {
	if !cmd.tokenRemains() {
		return cmd.tuntapShow()
	}

	c := cmd.findPrefix("add", "del", "show", "list", "lst", "help")

	options, err := cmd.parseTunTap()
	if err != nil {
		return err
	}

	switch c {
	case "add":
		return cmd.tuntapAdd(options)
	case "del":
		return cmd.tuntapDel(options)
	case "show", "list", "lst":
		return cmd.tuntapShow()
	case "help":
		fmt.Fprint(cmd.Out, tuntapHelp)

		return nil
	default:
		return cmd.usage()
	}
}

type tuntapOptions struct {
	Mode  netlink.TuntapMode
	User  int
	Group int
	Name  string
	Flags netlink.TuntapFlag
}

var defaultTuntapOptions = tuntapOptions{
	Mode:  netlink.TUNTAP_MODE_TUN,
	User:  -1,
	Group: -1,
	Name:  "",
	Flags: netlink.TUNTAP_DEFAULTS,
}

func (cmd *cmd) parseTunTap() (tuntapOptions, error) {
	var err error
	options := defaultTuntapOptions

	for cmd.tokenRemains() {
		switch cmd.findPrefix("mode", "user", "group", "one_queue", "pi", "vnet_hdr", "multi_queue", "name", "dev") {
		case "mode":
			switch cmd.nextToken("tun, tap") {
			case "tun":
				options.Mode = netlink.TUNTAP_MODE_TUN
			case "tap":
				options.Mode = netlink.TUNTAP_MODE_TAP
			default:
				return tuntapOptions{}, fmt.Errorf("invalid mode %s", cmd.currentToken())
			}
		case "user":
			options.User, err = cmd.parseInt("USER")
			if err != nil {
				return tuntapOptions{}, err
			}
		case "group":
			options.Group, err = cmd.parseInt("GROUP")
			if err != nil {
				return tuntapOptions{}, err
			}
		case "dev", "name":
			options.Name = cmd.nextToken("NAME")
		case "one_queue":
			options.Flags |= netlink.TUNTAP_ONE_QUEUE
		case "pi":
			options.Flags &^= netlink.TUNTAP_NO_PI
		case "vnet_hdr":
			options.Flags |= netlink.TUNTAP_VNET_HDR
		case "multi_queue":
			options.Flags = netlink.TUNTAP_MULTI_QUEUE_DEFAULTS
			options.Flags &^= netlink.TUNTAP_ONE_QUEUE

		default:
			return tuntapOptions{}, cmd.usage()
		}
	}

	return options, nil
}

func (cmd *cmd) tuntapAdd(options tuntapOptions) error {
	link := tunTapDevice(options)

	if err := cmd.handle.LinkAdd(link); err != nil {
		return err
	}

	return nil
}

func tunTapDevice(options tuntapOptions) *netlink.Tuntap {
	link := &netlink.Tuntap{
		LinkAttrs: netlink.LinkAttrs{
			Name: options.Name,
		},
		Mode: options.Mode,
	}

	if options.User >= 0 && options.User <= math.MaxUint16 {
		link.Owner = uint32(options.User)
	}

	if options.Group >= 0 && options.Group <= math.MaxUint16 {
		link.Group = uint32(options.Group)
	}

	link.Flags = options.Flags

	return link
}

func (cmd *cmd) tuntapDel(options tuntapOptions) error {
	links, err := cmd.handle.LinkList()
	if err != nil {
		return err
	}

	tuntap, err := filterTunTaps(links, options)
	if err != nil {
		return err
	}

	if err := cmd.handle.LinkDel(tuntap); err != nil {
		return err
	}

	return nil
}

func filterTunTaps(links []netlink.Link, options tuntapOptions) (*netlink.Tuntap, error) {
	filteredTunTaps := make([]*netlink.Tuntap, 0)

	for _, link := range links {
		tunTap, ok := link.(*netlink.Tuntap)
		if !ok {
			continue
		}

		if options.Name != "" && tunTap.Name != options.Name {
			continue
		}

		if options.Mode != 0 && tunTap.Mode != options.Mode {
			continue
		}

		filteredTunTaps = append(filteredTunTaps, tunTap)
	}

	if len(filteredTunTaps) != 1 {
		return nil, fmt.Errorf("found %d matching tun/tap devices", len(filteredTunTaps))
	}

	return filteredTunTaps[0], nil
}

type Tuntap struct {
	IfName string   `json:"ifname"`
	Flags  []string `json:"flags"`
}

func (cmd *cmd) tuntapShow() error {
	links, err := cmd.handle.LinkList()
	if err != nil {
		return err
	}

	return cmd.printTunTaps(links)
}

func (cmd *cmd) printTunTaps(links []netlink.Link) error {
	prints := make([]Tuntap, 0)

	for _, link := range links {
		tunTap, ok := link.(*netlink.Tuntap)
		if !ok {
			continue
		}

		var obj Tuntap

		obj.Flags = append(obj.Flags, tunTap.Mode.String())

		if tunTap.Flags&netlink.TUNTAP_NO_PI == 1 {
			obj.Flags = append(obj.Flags, "pi")
		}

		if tunTap.Flags&netlink.TUNTAP_ONE_QUEUE != 0 {
			obj.Flags = append(obj.Flags, "one_queue")
		} else if tunTap.Flags&netlink.TUNTAP_MULTI_QUEUE != 0 {
			obj.Flags = append(obj.Flags, "multi_queue")
		}

		if tunTap.Flags&netlink.TUNTAP_VNET_HDR != 0 {
			obj.Flags = append(obj.Flags, "vnet_hdr")
		}

		if tunTap.NonPersist {
			obj.Flags = append(obj.Flags, "non-persist")
		} else {
			obj.Flags = append(obj.Flags, "persist")
		}

		if tunTap.Owner != 0 {
			obj.Flags = append(obj.Flags, fmt.Sprintf("user %d", tunTap.Owner))
		}

		if tunTap.Group != 0 {
			obj.Flags = append(obj.Flags, fmt.Sprintf("group %d", tunTap.Group))
		}

		obj.IfName = tunTap.Name

		prints = append(prints, obj)
	}

	if cmd.Opts.JSON {
		return printJSON(*cmd, prints)
	}

	for _, print := range prints {
		output := fmt.Sprintf("%s:", print.IfName)

		for _, flag := range print.Flags {
			output += fmt.Sprintf(" %s", flag)
		}

		fmt.Fprintln(cmd.Out, output)
	}

	return nil
}
