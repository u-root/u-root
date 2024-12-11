// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

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
Usage: ip xfrm policy { list }[ SELECTOR ] [ dir DIR ]
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

func (cmd *cmd) parseXfrmPolicyTmpl() (*netlink.XfrmPolicyTmpl, error) {
	var err error

	tmpl := &netlink.XfrmPolicyTmpl{}

	for cmd.tokenRemains() {
		switch cmd.nextToken("src", "dst", "proto", "spi", "mode", "reqid", "level") {
		case "src":
			tmpl.Src, err = cmd.parseAddress()
			if err != nil {
				return nil, err
			}
		case "dst":
			tmpl.Dst, err = cmd.parseAddress()
			if err != nil {
				return nil, err
			}
		case "proto":
			tmpl.Proto, err = cmd.parseXfrmProto()
			if err != nil {
				return nil, err
			}
		case "spi":
			tmpl.Spi, err = cmd.parseInt("SPI")
			if err != nil {
				return nil, err
			}
		case "mode":
			tmpl.Mode, err = cmd.parseXfrmMode()
			if err != nil {
				return nil, err
			}
		case "reqid":
			tmpl.Reqid, err = cmd.parseInt("REQID")
			if err != nil {
				return nil, err
			}
		case "level":
			c := cmd.nextToken("required", "use")
			switch c {
			case "use":
				tmpl.Optional = 1
			case "required":
				tmpl.Optional = 0
			default:
				return nil, cmd.usage()
			}
		default:
			return nil, cmd.usage()
		}
	}

	return tmpl, nil
}

func (cmd *cmd) xfrmPolicy() error {
	switch cmd.nextToken("add", "update", "delete", "get", "deleteall", "show", "list", "flush", "count", "set", "help") {
	case "add":
		policy, err := cmd.parseXfrmPolicyAddUpdate()
		if err != nil {
			return err
		}

		return cmd.handle.XfrmPolicyAdd(policy)
	case "update":
		policy, err := cmd.parseXfrmPolicyAddUpdate()
		if err != nil {
			return err
		}

		return cmd.handle.XfrmPolicyUpdate(policy)
	case "delete":
		policy, err := cmd.parseXfrmPolicyDeleteGet()
		if err != nil {
			return err
		}

		return cmd.handle.XfrmPolicyDel(policy)
	case "get":
		policy, err := cmd.parseXfrmPolicyDeleteGet()
		if err != nil {
			return err
		}

		policy, err = cmd.handle.XfrmPolicyGet(policy)
		if err != nil {
			return err
		}

		printXfrmPolicy(cmd.Out, *policy)
	case "deleteall":
		policy, err := cmd.parseXfrmPolicyListDeleteAll()
		if err != nil {
			return err
		}

		return cmd.handle.XfrmPolicyDel(policy)
	case "list", "show":
		policy, err := cmd.parseXfrmPolicyListDeleteAll()
		if err != nil {
			return err
		}

		policies, err := netlink.XfrmPolicyList(cmd.Family)
		if err != nil {
			return err
		}

		printFilteredXfrmPolicies(cmd.Out, policies, policy)

		return nil
	case "flush":
		return cmd.handle.XfrmPolicyFlush()
	case "count":
		policies, err := cmd.handle.XfrmPolicyList(cmd.Family)
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.Out, "XFRM policies: %d\n", len(policies))
	case "help":
		fmt.Fprint(cmd.Out, xfrmPolicyHelp)

		return nil
	default:
		return cmd.usage()
	}

	return nil
}

func (cmd *cmd) parseXfrmPolicyAddUpdate() (*netlink.XfrmPolicy, error) {
	var err error

	policy := &netlink.XfrmPolicy{}

	for cmd.tokenRemains() {
		switch cmd.nextToken("src", "dst", "dir", "proto", "sport", "dport", "mark", "index", "action", "priority", "if_id", "tmpl") {
		case "src":
			policy.Src, err = cmd.parseIPNet()
			if err != nil {
				return nil, err
			}
		case "dst":
			policy.Dst, err = cmd.parseIPNet()
			if err != nil {
				return nil, err
			}
		case "proto":
			policy.Proto, err = cmd.parseXfrmProto()
			if err != nil {
				return nil, err
			}
		case "sport":
			policy.SrcPort, err = cmd.parseInt("SPORT")
			if err != nil {
				return nil, err
			}
		case "dport":
			policy.DstPort, err = cmd.parseInt("DPORT")
			if err != nil {
				return nil, err
			}
		case "dir":
			policy.Dir, err = cmd.parseXfrmDir()
			if err != nil {
				return nil, err
			}
		case "mark":
			policy.Mark, err = cmd.parseXfrmMark()
			if err != nil {
				return nil, err
			}
		case "index":
			policy.Index, err = cmd.parseInt("INDEX")
			if err != nil {
				return nil, err
			}
		case "action":
			policy.Action, err = cmd.parseXfrmAction()
			if err != nil {
				return nil, err
			}
		case "priority":
			policy.Priority, err = cmd.parseInt("PRIORITY")
			if err != nil {
				return nil, err
			}
		case "if_id":
			policy.Ifid, err = cmd.parseInt("IF_ID")
			if err != nil {
				return nil, err
			}
		case "tmpl":
			tmpl, err := cmd.parseXfrmPolicyTmpl()
			if err != nil {
				return nil, err
			}
			policy.Tmpls = append(policy.Tmpls, *tmpl)

		default:
			return nil, cmd.usage()
		}
	}

	return policy, nil
}

func (cmd *cmd) parseXfrmPolicyDeleteGet() (*netlink.XfrmPolicy, error) {
	var (
		indexSpecified    bool
		selectorSpecified bool
		err               error
	)

	policy := &netlink.XfrmPolicy{}

	for cmd.tokenRemains() {
		switch cmd.nextToken("src", "dst", "dir", "proto", "sport", "dport", "mark", "index", "if_id") {
		case "src":
			policy.Src, err = cmd.parseIPNet()
			if err != nil {
				return nil, err
			}
			selectorSpecified = true
		case "dst":
			policy.Dst, err = cmd.parseIPNet()
			if err != nil {
				return nil, err
			}
			selectorSpecified = true
		case "proto":
			policy.Proto, err = cmd.parseXfrmProto()
			if err != nil {
				return nil, err
			}
			selectorSpecified = true
		case "sport":
			policy.SrcPort, err = cmd.parseInt("SPORT")
			if err != nil {
				return nil, err
			}
			selectorSpecified = true
		case "dport":
			policy.DstPort, err = cmd.parseInt("DPORT")
			if err != nil {
				return nil, err
			}
			selectorSpecified = true
		case "dir":
			policy.Dir, err = cmd.parseXfrmDir()
			if err != nil {
				return nil, err
			}
		case "mark":
			policy.Mark, err = cmd.parseXfrmMark()
			if err != nil {
				return nil, err
			}
		case "index":
			policy.Index, err = cmd.parseInt("INDEX")
			if err != nil {
				return nil, err
			}
			indexSpecified = true
		case "if_id":
			policy.Ifid, err = cmd.parseInt("IF_ID")
			if err != nil {
				return nil, err
			}
		default:
			return nil, cmd.usage()
		}
	}

	if selectorSpecified && indexSpecified {
		return nil, fmt.Errorf("cannot specify both SELECTOR and index")
	}
	return policy, nil
}

func (cmd *cmd) parseXfrmPolicyListDeleteAll() (*netlink.XfrmPolicy, error) {
	var err error

	policy := &netlink.XfrmPolicy{}

	for cmd.tokenRemains() {
		switch cmd.nextToken("src", "dst", "dir", "proto", "sport", "dport", "index", "action", "priority", "mark", "if_id") {
		case "src":
			policy.Src, err = cmd.parseIPNet()
			if err != nil {
				return nil, err
			}
		case "dst":
			policy.Dst, err = cmd.parseIPNet()
			if err != nil {
				return nil, err
			}
		case "proto":
			policy.Proto, err = cmd.parseXfrmProto()
			if err != nil {
				return nil, err
			}
		case "sport":
			policy.SrcPort, err = cmd.parseInt("SPORT")
			if err != nil {
				return nil, err
			}
		case "dport":
			policy.DstPort, err = cmd.parseInt("DPORT")
			if err != nil {
				return nil, err
			}
		case "dir":
			policy.Dir, err = cmd.parseXfrmDir()
			if err != nil {
				return nil, err
			}
		case "mark":
			policy.Mark, err = cmd.parseXfrmMark()
			if err != nil {
				return nil, err
			}
		case "index":
			policy.Index, err = cmd.parseInt("INDEX")
			if err != nil {
				return nil, err
			}
		case "if_id":
			policy.Ifid, err = cmd.parseInt("IF_ID")
			if err != nil {
				return nil, err
			}
		default:
			return nil, cmd.usage()
		}
	}

	return policy, nil
}

func printXfrmPolicy(w io.Writer, policy netlink.XfrmPolicy) {
	fmt.Fprintf(w, "src %s dst %s\n", policy.Src, policy.Dst)
	fmt.Fprintf(w, "\t%s priority %d\n", policy.Dir, policy.Priority)
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

func printFilteredXfrmPolicies(w io.Writer, policies []netlink.XfrmPolicy, filter *netlink.XfrmPolicy) {
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
		printXfrmPolicy(w, policy)
		fmt.Fprintln(w)
	}
}
