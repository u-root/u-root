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
	xfrmStateHelp = `Usage: ip xfrm state { add | update } ID [ ALGO-LIST ] [ mode MODE ]
        [ mark MARK [ mask MASK ] ] [ reqid REQID ] [ replay-window SIZE ] 
        [ flag FLAG-LIST ] [ LIMIT-LIST ] [ encap ENCAP ]
        [ output-mark OUTPUT-MARK [ mask MASK ] [ if_id IF_ID ] 
Usage: ip xfrm state allocspi ID [ mode MODE ] [ mark MARK [ mask MASK ] ]
        [ reqid REQID ] 
Usage: ip xfrm state { delete | get } ID [ mark MARK [ mask MASK ] ]
Usage: ip xfrm state deleteall [ ID ] [ mode MODE ] [ reqid REQID ]
Usage: ip xfrm state list [ nokeys ] [ ID ] [ mode MODE ] [ reqid REQID ]
Usage: ip xfrm state flush [ proto XFRM-PROTO ]
Usage: ip xfrm state count
ID := [ src ADDR ] [ dst ADDR ] [ proto XFRM-PROTO ] [ spi SPI ]
XFRM-PROTO := esp | ah | comp | route2 | hao
ALGO-LIST := [ ALGO-LIST ] ALGO
ALGO := { enc | auth } ALGO-NAME ALGO-KEYMAT |
        auth-trunc ALGO-NAME ALGO-KEYMAT ALGO-TRUNC-LEN |
        aead ALGO-NAME ALGO-KEYMAT ALGO-ICV-LEN |
MODE := transport | tunnel | beet | ro | in_trigger
LIMIT-LIST := [ LIMIT-LIST ] limit LIMIT
LIMIT := { time-soft | time-hard | time-use-soft | time-use-hard } SECONDS |
         { byte-soft | byte-hard } SIZE | { packet-soft | packet-hard } COUNT
ENCAP := { espinudp | espinudp-nonike | espintcp } SPORT DPORT OADDR`
)

func (cmd *cmd) xfrmState() error {
	switch cmd.findPrefix("add", "update", "allocspi", "delete", "deleteall", "show", "list", "flush", "count", "help") {
	case "add":
		xfrmState, err := cmd.parseXfrmStateAddUpdate()
		if err != nil {
			return err
		}

		return cmd.handle.XfrmStateAdd(xfrmState)
	case "update":
		xfrmState, err := cmd.parseXfrmStateAddUpdate()
		if err != nil {
			return err
		}

		return cmd.handle.XfrmStateUpdate(xfrmState)
	case "allocspi":
		xfrmState, err := cmd.parseXfrmStateAllocSPI()
		if err != nil {
			return err
		}

		if _, err := netlink.XfrmStateAllocSpi(xfrmState); err != nil {
			return err
		}

	case "delete":
		xfrmState, err := cmd.parseXfrmStateDeleteGet()
		if err != nil {
			return err
		}

		return cmd.handle.XfrmStateDel(xfrmState)
	case "get":
		xfrmState, err := cmd.parseXfrmStateDeleteGet()
		if err != nil {
			return err
		}

		xfrmState, err = cmd.handle.XfrmStateGet(xfrmState)
		if err != nil {
			return err
		}

		printXfrmState(cmd.Out, *xfrmState, true)
	case "list", "show":

		xfrmState, noKeys, err := cmd.parseXfrmStateListDeleteAll()
		if err != nil {
			return err
		}

		states, err := netlink.XfrmStateList(cmd.Family)
		if err != nil {
			return err
		}

		cmd.printFilteredXfrmStates(states, xfrmState, noKeys)
	case "count":
		states, err := cmd.handle.XfrmStateList(cmd.Family)
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.Out, "XFRM states: %d\n", len(states))
	case "flush":
		return cmd.xfrmStateFlush()
	case "deleteall":
		return cmd.xfrmStateDeleteAll()
	case "help":
		fmt.Fprint(cmd.Out, xfrmStateHelp)

		return nil
	default:
		return cmd.usage()
	}

	return nil
}

func (cmd *cmd) xfrmStateFlush() error {
	if !cmd.tokenRemains() {
		return cmd.handle.XfrmStateFlush(0)
	}

	if cmd.nextToken("proto") != "proto" {
		return cmd.usage()
	}

	proto, err := cmd.parseXfrmProto()
	if err != nil {
		return err
	}

	return cmd.handle.XfrmStateFlush(proto)
}

func (cmd *cmd) xfrmStateDeleteAll() error {
	filter, noKeys, err := cmd.parseXfrmStateListDeleteAll()
	if err != nil {
		return err
	}

	if noKeys {
		return fmt.Errorf("deleteall does not support nokeys")
	}

	states, err := cmd.handle.XfrmStateList(cmd.Family)
	if err != nil {
		return err
	}

	for _, state := range states {
		if filter != nil {
			if filter.Src != nil && !filter.Src.Equal(state.Src) {
				continue
			}
			if filter.Dst != nil && !filter.Dst.Equal(state.Dst) {
				continue
			}
			if filter.Proto != 0 && filter.Proto != state.Proto {
				continue
			}
			if filter.Spi != 0 && filter.Spi != state.Spi {
				continue
			}
			if filter.Mode != 0 && filter.Mode != state.Mode {
				continue
			}
			if filter.Reqid != 0 && filter.Reqid != state.Reqid {
				continue
			}
		}

		if err := cmd.handle.XfrmStateDel(&state); err != nil {
			return err
		}
	}

	return nil
}

func (cmd *cmd) printFilteredXfrmStates(states []netlink.XfrmState, filter *netlink.XfrmState, noKeys bool) {
	for _, state := range states {
		if filter != nil {
			if filter.Src != nil && !filter.Src.Equal(state.Src) {
				continue
			}
			if filter.Dst != nil && !filter.Dst.Equal(state.Dst) {
				continue
			}
			if filter.Proto != 0 && filter.Proto != state.Proto {
				continue
			}
			if filter.Spi != 0 && filter.Spi != state.Spi {
				continue
			}
			if filter.Mode != 0 && filter.Mode != state.Mode {
				continue
			}
			if filter.Reqid != 0 && filter.Reqid != state.Reqid {
				continue
			}
		}
		printXfrmState(cmd.Out, state, noKeys)
		fmt.Fprintln(cmd.Out)
	}
}

func printXfrmState(w io.Writer, state netlink.XfrmState, noKeys bool) {
	fmt.Fprintf(w, "src %s dst %s\n", state.Src, state.Dst)
	fmt.Fprintf(w, "\tproto %s spi 0x%x mode %s\n", state.Proto, state.Spi, state.Mode)

	options := "\t"

	if state.Reqid != 0 {
		options += fmt.Sprintf("reqid %d", state.Reqid)
	}

	if state.ReplayWindow != 0 {
		options += fmt.Sprintf(" replay-window %d", state.ReplayWindow)
	}

	if options != "\t" {
		fmt.Fprintln(w, options)
	}

	if state.Auth != nil {
		if noKeys {
			fmt.Fprintf(w, "\tauth %s %dbits\n", state.Auth.Name, len(state.Auth.Key)*8)
		} else {
			fmt.Fprintf(w, "\tauth %s 0x%x %dbits\n", state.Auth.Name, state.Auth.Key, len(state.Auth.Key)*8)
		}
	}
	if state.Crypt != nil {
		if noKeys {
			fmt.Fprintf(w, "\tenc %s %dbits\n", state.Crypt.Name, len(state.Crypt.Key)*8)
		} else {
			fmt.Fprintf(w, "\tenc %s 0x%x %dbits\n", state.Crypt.Name, state.Crypt.Key, len(state.Crypt.Key)*8)
		}
	}
	if state.Aead != nil {
		if noKeys {
			fmt.Fprintf(w, "\taead %s %dbits\n", state.Aead.Name, len(state.Aead.Key)*8)
		} else {
			fmt.Fprintf(w, "\taead %s 0x%x %dbits\n", state.Aead.Name, state.Aead.Key, len(state.Aead.Key)*8)
		}
	}
	if state.Encap != nil {
		fmt.Fprintf(w, "\tencap type %s sport %d dport %d addr %s\n", state.Encap.Type, state.Encap.SrcPort, state.Encap.DstPort, state.Encap.OriginalAddress)
	}

	if state.Mark != nil {
		fmt.Fprintf(w, "\tmark %d", state.Mark.Value)
		if state.Mark.Mask != 0 {
			fmt.Fprintf(w, "/%x", state.Mark.Mask)
		}
		fmt.Fprintln(w)
	}

	if state.OutputMark != nil {
		fmt.Fprintf(w, "\toutput-mark %d", state.OutputMark.Value)
		if state.OutputMark.Mask != 0 {
			fmt.Fprintf(w, "/%x", state.OutputMark.Mask)
		}
		fmt.Fprintln(w)
	}

	if state.Limits.ByteSoft != 0 || state.Limits.ByteHard != 0 {
		fmt.Fprintf(w, "\tsoft-byte-limit %d hard-byte-limit %d\n", state.Limits.ByteSoft, state.Limits.ByteHard)
	}
	if state.Limits.PacketSoft != 0 || state.Limits.PacketHard != 0 {
		fmt.Fprintf(w, "\tsoft-packet-limit %d hard-packet-limit %d\n", state.Limits.PacketSoft, state.Limits.PacketHard)
	}
	if state.Limits.TimeSoft != 0 || state.Limits.TimeHard != 0 {
		fmt.Fprintf(w, "\tsoft-add-expires-seconds %d hard-add-expires-seconds %d\n", state.Limits.TimeSoft, state.Limits.TimeHard)
	}
	if state.Limits.TimeUseSoft != 0 || state.Limits.TimeUseHard != 0 {
		fmt.Fprintf(w, "\tsoft-use-expires-seconds %d hard-use-expires-seconds %d\n", state.Limits.TimeUseSoft, state.Limits.TimeUseHard)
	}

	fmt.Fprintf(w, "statistics: replay-window %d replay %d failed %d bytes %d packets %d\n", state.Statistics.ReplayWindow, state.Statistics.Replay, state.Statistics.Failed, state.Statistics.Bytes, state.Statistics.Packets)
}

func (cmd *cmd) parseXfrmStateAddUpdate() (*netlink.XfrmState, error) {
	var err error

	state := &netlink.XfrmState{}

	for cmd.tokenRemains() {
		switch cmd.nextToken("src", "dst", "proto", "spi", "enc", "auth", "auth-trunc", "aead", "comp", "mode", "mark", "reqid", "replay-window", "limit", "encap", "output-mark", "if_id") {
		case "src":
			state.Src, err = cmd.parseAddress()
			if err != nil {
				return nil, err
			}
		case "dst":
			state.Dst, err = cmd.parseAddress()
			if err != nil {
				return nil, err
			}
		case "proto":
			state.Proto, err = cmd.parseXfrmProto()
			if err != nil {
				return nil, err
			}
		case "spi":
			state.Spi, err = cmd.parseInt("SPI")
			if err != nil {
				return nil, err
			}
		case "enc":
			name := cmd.nextToken("ALGO-NAME")

			key, err := cmd.parseByte("ALGO-KEYMAT")
			if err != nil {
				return nil, err
			}

			state.Crypt = &netlink.XfrmStateAlgo{
				Name: name,
				Key:  key,
			}
		case "auth":
			name := cmd.nextToken("ALGO-NAME")

			key, err := cmd.parseByte("ALGO-KEYMAT")
			if err != nil {
				return nil, err
			}

			state.Auth = &netlink.XfrmStateAlgo{
				Name: name,
				Key:  key,
			}

		case "auth-trunc":
			name := cmd.nextToken("ALGO-NAME")

			key, err := cmd.parseByte("ALGO-KEYMAT")
			if err != nil {
				return nil, err
			}
			truncLen, err := cmd.parseInt("ALGO-TRUNC-LEN")
			if err != nil {
				return nil, err
			}

			state.Auth = &netlink.XfrmStateAlgo{
				Name:        name,
				Key:         key,
				TruncateLen: truncLen,
			}
		case "aead":
			name := cmd.nextToken("ALGO-NAME")

			key, err := cmd.parseByte("ALGO-KEYMAT")
			if err != nil {
				return nil, err
			}
			icvLen, err := cmd.parseInt("ALGO-ICV-LEN")
			if err != nil {
				return nil, err
			}

			state.Aead = &netlink.XfrmStateAlgo{
				Name:   name,
				Key:    key,
				ICVLen: icvLen,
			}
		case "comp":
			return nil, fmt.Errorf("comp not implemented")
		case "mode":
			state.Mode, err = cmd.parseXfrmMode()
			if err != nil {
				return nil, err
			}
		case "mark":
			state.Mark, err = cmd.parseXfrmMark()
			if err != nil {
				return nil, err
			}
		case "reqid":
			state.Reqid, err = cmd.parseInt("REQID")
			if err != nil {
				return nil, err
			}
		case "replay-window":
			state.ReplayWindow, err = cmd.parseInt("SIZE")
			if err != nil {
				return nil, err
			}
		case "limit":
			state.Limits, err = cmd.parseXfrmLimit()
			if err != nil {
				return nil, err
			}
		case "encap":
			state.Encap, err = cmd.parseXfrmEncap()
			if err != nil {
				return nil, err
			}
		case "output-mark":
			state.OutputMark, err = cmd.parseXfrmMark()
			if err != nil {
				return nil, err
			}
		case "if_id":
			state.Ifid, err = cmd.parseInt("IF_ID")
			if err != nil {
				return nil, err
			}
		default:
			return nil, cmd.usage()
		}
	}

	return state, nil
}

func (cmd *cmd) parseXfrmStateAllocSPI() (*netlink.XfrmState, error) {
	var err error

	state := &netlink.XfrmState{}

	for cmd.tokenRemains() {
		switch cmd.nextToken("src", "dst", "proto", "spi", "mode", "mark", "reqid") {
		case "src":
			state.Src, err = cmd.parseAddress()
			if err != nil {
				return nil, err
			}
		case "dst":
			state.Dst, err = cmd.parseAddress()
			if err != nil {
				return nil, err
			}
		case "proto":
			state.Proto, err = cmd.parseXfrmProto()
			if err != nil {
				return nil, err
			}
		case "spi":
			state.Spi, err = cmd.parseInt("SPI")
			if err != nil {
				return nil, err
			}
		case "mode":
			state.Mode, err = cmd.parseXfrmMode()
			if err != nil {
				return nil, err
			}
		case "mark":
			state.Mark, err = cmd.parseXfrmMark()
			if err != nil {
				return nil, err
			}
		case "reqid":
			state.Reqid, err = cmd.parseInt("REQID")
			if err != nil {
				return nil, err
			}
		default:
			return nil, cmd.usage()
		}
	}

	return state, nil
}

func (cmd *cmd) parseXfrmStateDeleteGet() (*netlink.XfrmState, error) {
	var err error

	state := &netlink.XfrmState{}

	for cmd.tokenRemains() {
		switch cmd.nextToken("src", "dst", "proto", "spi", "mark") {
		case "src":
			state.Src, err = cmd.parseAddress()
			if err != nil {
				return nil, err
			}
		case "dst":
			state.Dst, err = cmd.parseAddress()
			if err != nil {
				return nil, err
			}
		case "proto":
			state.Proto, err = cmd.parseXfrmProto()
			if err != nil {
				return nil, err
			}
		case "spi":
			state.Spi, err = cmd.parseInt("SPI")
			if err != nil {
				return nil, err
			}
		case "mark":
			state.Mark, err = cmd.parseXfrmMark()
			if err != nil {
				return nil, err
			}
		default:
			return nil, cmd.usage()
		}
	}

	return state, nil
}

func (cmd *cmd) parseXfrmStateListDeleteAll() (*netlink.XfrmState, bool, error) {
	var (
		noKeys bool
		err    error
	)

	state := &netlink.XfrmState{}

	for cmd.tokenRemains() {
		switch cmd.nextToken("src", "dst", "proto", "spi", "mode", "reqid", "nokeys") {
		case "src":
			state.Src, err = cmd.parseAddress()
			if err != nil {
				return nil, false, err
			}
		case "dst":
			state.Dst, err = cmd.parseAddress()
			if err != nil {
				return nil, false, err
			}
		case "proto":
			state.Proto, err = cmd.parseXfrmProto()
			if err != nil {
				return nil, false, err
			}
		case "spi":
			state.Spi, err = cmd.parseInt("SPI")
			if err != nil {
				return nil, false, err
			}
		case "mode":
			state.Mode, err = cmd.parseXfrmMode()
			if err != nil {
				return nil, false, err
			}
		case "reqid":
			state.Reqid, err = cmd.parseInt("REQID")
			if err != nil {
				return nil, false, err
			}
		case "nokeys":
			noKeys = true
		default:
			return nil, false, cmd.usage()
		}
	}

	return state, noKeys, nil
}
