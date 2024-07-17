// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"

	"github.com/vishvananda/netlink"
)

const (
	xfrmPolicyHelp = `Usage: ip xfrm policy { add | update } SELECTOR dir DIR
	[ mark MARK [ mask MASK ] ] [ index INDEX ] [ ptype PTYPE ]
	[ action ACTION ] [ priority PRIORITY ] [ if_id IF_ID ] [ TMPL-LIST ]
Usage: ip xfrm policy { delete | get } { SELECTOR | index INDEX } dir DIR
	[ mark MARK [ mask MASK ] ] [ if_id IF_ID ]
Usage: ip xfrm policy { deleteall | list }[ SELECTOR ] [ dir DIR ]
	[ index INDEX ][ action ACTION ] [ priority PRIORITY ]
Usage: ip xfrm policy flush 
Usage: ip xfrm policy count
SELECTOR := [ src ADDR[/PLEN] ] [ dst ADDR[/PLEN] ] [ dev DEV ] [ UPSPEC ]
UPSPEC := proto { { tcp | udp | sctp | dccp } [ sport PORT ] [ dport PORT ]
DIR := in | out | fwd
ACTION := allow | block
TMPL-LIST := [ TMPL-LIST ] tmpl TMPL
TMPL := ID [ mode MODE ] [ reqid REQID ] [ level LEVEL ]
ID := [ src ADDR ] [ dst ADDR ] [ proto XFRM-PROTO ] [ spi SPI ]
XFRM-PROTO := esp | ah | comp | route2 | hao
MODE := transport | tunnel | beet | ro | in_trigger
LEVEL := required | use

`
)

func parseXfrmPolicyTmpl() (*netlink.XfrmPolicyTmpl, error) {
	var err error

	tmpl := &netlink.XfrmPolicyTmpl{}

	cursor++

	for {
		cursor++

		if cursor == len(arg) {
			break
		}
		expectedValues = []string{"src", "dst", "proto", "spi", "mode", "reqid", "level"}
		switch arg[cursor] {
		case "src":
			tmpl.Src, err = parseAddress()
			if err != nil {
				return nil, err
			}
		case "dst":
			tmpl.Dst, err = parseAddress()
			if err != nil {
				return nil, err
			}
		case "proto":
			tmpl.Proto, err = parseXfrmProto()
			if err != nil {
				return nil, err
			}
		case "spi":
			tmpl.Spi, err = parseInt("SPI")
			if err != nil {
				return nil, err
			}
		case "mode":
			tmpl.Mode, err = parseXfrmMode()
			if err != nil {
				return nil, err
			}
		case "reqid":
			tmpl.Reqid, err = parseInt("REQID")
			if err != nil {
				return nil, err
			}
		case "level":
			cursor++
			expectedValues = []string{"required", "use"}
			if arg[cursor] == "use" {
				tmpl.Optional = 1
			}
		default:
			return nil, usage()
		}
	}

	return tmpl, nil
}

func xfrmPolicy(w io.Writer) error {
	cursor++
	expectedValues = []string{"add", "update", "delete", "get", "deleteall", "show", "list", "flush", "count", "set", "help"}
	switch findPrefix(arg[cursor], expectedValues) {
	case "add":
		policy, err := parseXfrmPolicyAddUpdate()
		if err != nil {
			return err
		}

		return netlink.XfrmPolicyAdd(policy)
	case "update":
		policy, err := parseXfrmPolicyAddUpdate()
		if err != nil {
			return err
		}

		return netlink.XfrmPolicyUpdate(policy)
	case "delete":
		policy, err := parseXfrmPolicyDeleteGet()
		if err != nil {
			return err
		}

		return netlink.XfrmPolicyDel(policy)
	case "get":
		policy, err := parseXfrmPolicyDeleteGet()
		if err != nil {
			return err
		}

		policy, err = netlink.XfrmPolicyGet(policy)
		if err != nil {
			return err
		}

		printXfrmPolicy(w, policy)
	case "deleteall":
		policy, err := parseXfrmPolicyListDeleteAll()
		if err != nil {
			return err
		}

		return netlink.XfrmPolicyDel(policy)
	case "list", "show":
		policy, err := parseXfrmPolicyListDeleteAll()
		if err != nil {
			return err
		}

		return printFilteredXfrmPolicies(w, policy, family)
	case "flush":
		return netlink.XfrmPolicyFlush()
	case "count":
		policies, err := netlink.XfrmPolicyList(family)
		if err != nil {
			return err
		}

		fmt.Fprintf(w, "XFRM policies: %d\n", len(policies))
	case "help":
		fmt.Fprint(w, xfrmPolicyHelp)

		return nil
	default:
		return usage()
	}

	return nil
}

func parseXfrmPolicyAddUpdate() (*netlink.XfrmPolicy, error) {
	var err error

	policy := &netlink.XfrmPolicy{}

	for {
		cursor++

		if cursor == len(arg) {
			break
		}

		expectedValues = []string{"src", "dst", "dir", "proto", "sport", "dport", "mark", "index", "action", "priority", "if_id", "tmpl"}
		switch arg[cursor] {
		case "src":
			policy.Src, err = parseIPNet()
			if err != nil {
				return nil, err
			}
		case "dst":
			policy.Dst, err = parseIPNet()
			if err != nil {
				return nil, err
			}
		case "proto":
			policy.Proto, err = parseXfrmProto()
			if err != nil {
				return nil, err
			}
		case "sport":
			policy.SrcPort, err = parseInt("SPORT")
			if err != nil {
				return nil, err
			}
		case "dport":
			policy.DstPort, err = parseInt("DPORT")
			if err != nil {
				return nil, err
			}
		case "dir":
			policy.Dir, err = parseXfrmDir()
			if err != nil {
				return nil, err
			}
		case "mark":
			policy.Mark, err = parseXfrmMark()
			if err != nil {
				return nil, err
			}
		case "index":
			policy.Index, err = parseInt("INDEX")
			if err != nil {
				return nil, err
			}
		case "action":
			policy.Action, err = parseXfrmAction()
			if err != nil {
				return nil, err
			}
		case "priority":
			policy.Priority, err = parseInt("PRIORITY")
			if err != nil {
				return nil, err
			}
		case "if_id":
			policy.Ifid, err = parseInt("IF_ID")
			if err != nil {
				return nil, err
			}
		case "tmpl":
			tmpl, err := parseXfrmPolicyTmpl()
			if err != nil {
				return nil, err
			}
			policy.Tmpls = append(policy.Tmpls, *tmpl)

		default:
			return nil, usage()
		}
	}

	return policy, nil
}

func parseXfrmPolicyDeleteGet() (*netlink.XfrmPolicy, error) {
	var (
		indexSpecified    bool
		selectorSpecified bool
		err               error
	)

	policy := &netlink.XfrmPolicy{}

	for {
		cursor++

		if cursor == len(arg) {
			break
		}

		expectedValues = []string{"src", "dst", "dir", "proto", "sport", "dport", "mark", "index", "if_id"}
		switch arg[cursor] {
		case "src":
			policy.Src, err = parseIPNet()
			if err != nil {
				return nil, err
			}
			selectorSpecified = true
		case "dst":
			policy.Dst, err = parseIPNet()
			if err != nil {
				return nil, err
			}
			selectorSpecified = true
		case "proto":
			policy.Proto, err = parseXfrmProto()
			if err != nil {
				return nil, err
			}
			selectorSpecified = true
		case "sport":
			policy.SrcPort, err = parseInt("SPORT")
			if err != nil {
				return nil, err
			}
			selectorSpecified = true
		case "dport":
			policy.DstPort, err = parseInt("DPORT")
			if err != nil {
				return nil, err
			}
			selectorSpecified = true
		case "dir":
			policy.Dir, err = parseXfrmDir()
			if err != nil {
				return nil, err
			}
		case "mark":
			policy.Mark, err = parseXfrmMark()
			if err != nil {
				return nil, err
			}
		case "index":
			policy.Index, err = parseInt("INDEX")
			if err != nil {
				return nil, err
			}
			indexSpecified = true
		case "if_id":
			policy.Ifid, err = parseInt("IF_ID")
			if err != nil {
				return nil, err
			}
		default:
			return nil, usage()
		}
	}

	if selectorSpecified && indexSpecified {
		return nil, fmt.Errorf("cannot specify both SELECTOR and index")
	}
	return policy, nil
}

func parseXfrmPolicyListDeleteAll() (*netlink.XfrmPolicy, error) {
	var err error

	policy := &netlink.XfrmPolicy{}

	for {
		cursor++

		if cursor == len(arg) {
			break
		}

		expectedValues = []string{"src", "dst", "dir", "proto", "sport", "dport", "index", "action", "priority"}
		switch arg[cursor] {
		case "src":
			policy.Src, err = parseIPNet()
			if err != nil {
				return nil, err
			}
		case "dst":
			policy.Dst, err = parseIPNet()
			if err != nil {
				return nil, err
			}
		case "proto":
			policy.Proto, err = parseXfrmProto()
			if err != nil {
				return nil, err
			}
		case "sport":
			policy.SrcPort, err = parseInt("SPORT")
			if err != nil {
				return nil, err
			}
		case "dport":
			policy.DstPort, err = parseInt("DPORT")
			if err != nil {
				return nil, err
			}
		case "dir":
			policy.Dir, err = parseXfrmDir()
			if err != nil {
				return nil, err
			}
		case "mark":
			policy.Mark, err = parseXfrmMark()
			if err != nil {
				return nil, err
			}
		case "index":
			policy.Index, err = parseInt("INDEX")
			if err != nil {
				return nil, err
			}
		case "if_id":
			policy.Ifid, err = parseInt("IF_ID")
			if err != nil {
				return nil, err
			}
		default:
			return nil, usage()
		}
	}

	return policy, nil
}

func printXfrmPolicy(w io.Writer, policy *netlink.XfrmPolicy) {
	fmt.Fprintf(w, "src %s dst %s\n", policy.Src, policy.Dst)
	fmt.Fprintf(w, "\tdir %s priorioty %d\n", policy.Dir, policy.Priority)
	fmt.Fprintf(w, "\tproto %s sport %d dport %d\n", policy.Proto, policy.SrcPort, policy.DstPort)
	fmt.Fprintf(w, "\taction %s if_id %d\n", policy.Action, policy.Ifid)

	if policy.Mark != nil {
		fmt.Fprintf(w, "\tmark %d", policy.Mark.Value)
		if policy.Mark.Mask != 0 {
			fmt.Fprintf(w, "/%x", policy.Mark.Mask)
		}
		fmt.Fprintln(w)
	}

	for _, tmpl := range policy.Tmpls {
		fmt.Fprintf(w, "\ttmpl src %s dst %s\n\t\tproto %s reqid %d mode %s spi %d\n", tmpl.Src, tmpl.Dst, tmpl.Proto, tmpl.Reqid, tmpl.Mode, tmpl.Spi)
	}
}

func printFilteredXfrmPolicies(w io.Writer, filter *netlink.XfrmPolicy, family int) error {
	policies, err := netlink.XfrmPolicyList(family)
	if err != nil {
		return err
	}

	for _, policy := range policies {
		if filter != nil {
			if filter.Src != nil && filter.Src.String() == policy.Src.String() {
				continue
			}
			if filter.Dst != nil && filter.Dst.String() == policy.Dst.String() {
				continue
			}
			if filter.Proto != 0 && filter.Proto != policy.Proto {
				continue
			}
			if filter.SrcPort != 0 && filter.SrcPort != policy.SrcPort {
				continue
			}
			if filter.DstPort != 0 && filter.DstPort != policy.DstPort {
				continue
			}
			if filter.Dir != 0 && filter.Dir != policy.Dir {
				continue
			}
			if filter.Mark != nil {
				if policy.Mark == nil {
					continue
				}
				if filter.Mark.Value != policy.Mark.Value && filter.Mark.Mask != policy.Mark.Mask {
					continue
				}
			}
			if filter.Index != 0 && filter.Index != policy.Index {
				continue
			}
			if filter.Ifid != 0 && filter.Ifid != policy.Ifid {
				continue
			}
		}
		printXfrmPolicy(w, &policy)
		fmt.Fprintln(w)
	}

	return nil
}
