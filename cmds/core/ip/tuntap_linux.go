// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
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

func tuntap(w io.Writer) error {
	cursor++
	if len(arg[cursor:]) == 0 {
		return tuntapShow(w)
	}

	expectedValues = []string{"add", "del", "show", "list", "lst", "help"}
	c := findPrefix(arg[cursor], expectedValues)

	options, err := parseTunTap()
	if err != nil {
		return err
	}

	switch c {
	case "add":
		return tuntapAdd(options)
	case "del":
		return tuntapDel(options)
	case "show", "list", "lst":
		return tuntapShow(w)
	case "help":
		fmt.Fprint(w, tuntapHelp)

		return nil
	default:
		return usage()
	}
}

type tuntapOptions struct {
	mode  netlink.TuntapMode
	user  int
	group int
	name  string
	flags netlink.TuntapFlag
}

var defaultTuntapOptions = tuntapOptions{
	mode:  netlink.TUNTAP_MODE_TUN,
	user:  -1,
	group: -1,
	name:  "",
	flags: netlink.TUNTAP_DEFAULTS,
}

func parseTunTap() (tuntapOptions, error) {
	var err error

	options := defaultTuntapOptions

	expectedValues = []string{"mode", "user", "group", "one_queue", "pi", "vnet_hdr", "multi_queue", "name", "dev"}
	for cursor < len(arg)-1 {
		cursor++
		switch arg[cursor] {
		case "mode":
			cursor++
			expectedValues = []string{"tun, tap"}

			switch arg[cursor] {
			case "tun":
				options.mode = netlink.TUNTAP_MODE_TUN
			case "tap":
				options.mode = netlink.TUNTAP_MODE_TAP
			default:
				return tuntapOptions{}, fmt.Errorf("invalid mode %s", arg[cursor])
			}
		case "user":
			options.user, err = parseInt("USER")
			if err != nil {
				return tuntapOptions{}, err
			}
		case "group":
			options.group, err = parseInt("GROUP")
			if err != nil {
				return tuntapOptions{}, err
			}
		case "dev", "name":
			options.name = parseString("NAME")
		case "one_queue":
			options.flags |= netlink.TUNTAP_ONE_QUEUE
		case "pi":
			options.flags &^= netlink.TUNTAP_NO_PI
		case "vnet_hdr":
			options.flags |= netlink.TUNTAP_VNET_HDR
		case "multi_queue":
			options.flags |= netlink.TUNTAP_MULTI_QUEUE_DEFAULTS
			options.flags &^= netlink.TUNTAP_ONE_QUEUE

		default:
			return tuntapOptions{}, usage()
		}
	}

	return options, nil
}

func tuntapAdd(options tuntapOptions) error {
	link := &netlink.Tuntap{
		LinkAttrs: netlink.LinkAttrs{
			Name: options.name,
		},
		Mode: options.mode,
	}

	if options.user >= 0 && options.user <= math.MaxUint16 {
		link.Owner = uint32(options.user)
	}

	if options.group >= 0 && options.group <= math.MaxUint16 {
		link.Group = uint32(options.group)
	}

	link.Flags = options.flags

	if err := netlink.LinkAdd(link); err != nil {
		return err
	}

	return nil
}

func tuntapDel(options tuntapOptions) error {
	links, err := netlink.LinkList()
	if err != nil {
		return err
	}

	filteredTunTaps := make([]*netlink.Tuntap, 0)

	for _, link := range links {
		tunTap, ok := link.(*netlink.Tuntap)
		if !ok {
			continue
		}

		if options.name != "" && tunTap.Name != options.name {
			continue
		}

		if options.mode != 0 && tunTap.Mode != options.mode {
			continue
		}

		filteredTunTaps = append(filteredTunTaps, tunTap)
	}

	if len(filteredTunTaps) != 1 {
		return fmt.Errorf("found %d matching tun/tap devices", len(filteredTunTaps))
	}

	if err := netlink.LinkDel(filteredTunTaps[0]); err != nil {
		return err
	}

	return nil
}

type Tuntap struct {
	IfName string   `json:"ifname"`
	Flags  []string `json:"flags"`
}

func tuntapShow(w io.Writer) error {
	links, err := netlink.LinkList()
	if err != nil {
		return err
	}

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

	if f.json {
		return printJSON(w, prints)
	}

	for _, print := range prints {
		output := fmt.Sprintf("%s:", print.IfName)

		for _, flag := range print.Flags {
			output += fmt.Sprintf(" %s", flag)
		}

		// Print the final output
		fmt.Fprintln(w, output)
	}

	return nil
}
