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
		return routeShow(w)
	}

	whatIWant = []string{"add", "del", "show", "list", "lst", "help"}
	c := findPrefix(arg[cursor], whatIWant)

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

	whatIWant = []string{"mode", "user", "group", "one_queue", "pi", "vnet_hdr", "multi_queue", "name", "dev"}
	for cursor < len(arg)-1 {
		cursor++
		switch arg[cursor] {
		case "mode":
			cursor++
			whatIWant = []string{"tun, tap"}

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

	if options.user >= 0 && options.user <= math.MaxUint32 {
		link.Owner = uint32(options.user)
	}

	if options.group >= 0 && options.group <= math.MaxUint32 {
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

func tuntapShow(w io.Writer) error {
	links, err := netlink.LinkList()
	if err != nil {
		return err
	}

	for _, link := range links {
		tunTap, ok := link.(*netlink.Tuntap)
		if !ok {
			continue
		}

		pi := ""
		if tunTap.Flags&netlink.TUNTAP_NO_PI == 0 {
			pi = "pi "
		}

		queue := ""
		if tunTap.Flags&netlink.TUNTAP_ONE_QUEUE != 0 {
			queue = "one_queue "
		} else if tunTap.Flags&netlink.TUNTAP_MULTI_QUEUE != 0 {
			queue = "multi_queue "
		}

		vnetHdr := ""
		if tunTap.Flags&netlink.TUNTAP_VNET_HDR != 0 {
			vnetHdr = "vnet_hdr "
		}

		persist := "persist "
		if tunTap.NonPersist {
			persist = "non-persist "
		}

		user := ""
		if tunTap.Owner != 0 {
			user = fmt.Sprintf("user %d ", tunTap.Owner)
		}

		group := ""
		if tunTap.Group != 0 {
			group = fmt.Sprintf("group %d", tunTap.Group)
		}

		fmt.Fprintf(w, "%s: %s %s%s%s%s%s%s\n", tunTap.Name, tunTap.Mode, pi, queue, vnetHdr, persist, user, group)
	}

	return nil
}
