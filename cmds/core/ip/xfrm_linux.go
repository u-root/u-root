// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netlink/nl"
)

const (
	xfrmHelp = `Usage: ip xfrm XFRM-OBJECT { COMMAND | help }
where  XFRM-OBJECT := policy | monitor`

	xfrmMonitorHelp = `Usage: ip xfrm monitor [ nokeys ] [ all | OBJECTS | help ]
OBJECTS := { acquire | expire | SA | aevent | policy | report }`
)

var xfrmFilterMap = map[string][]nl.XfrmMsgType{
	"acquire": {nl.XFRM_MSG_ACQUIRE},
	"expire":  {nl.XFRM_MSG_EXPIRE, nl.XFRM_MSG_POLEXPIRE},
	"SA":      {nl.XFRM_MSG_NEWSA, nl.XFRM_MSG_DELSA, nl.XFRM_MSG_UPDSA, nl.XFRM_MSG_GETSA, nl.XFRM_MSG_FLUSHSA},
	"aevent":  {nl.XFRM_MSG_NEWAE, nl.XFRM_MSG_GETAE},
	"policy":  {nl.XFRM_MSG_NEWPOLICY, nl.XFRM_MSG_UPDPOLICY, nl.XFRM_MSG_DELPOLICY, nl.XFRM_MSG_GETPOLICY, nl.XFRM_MSG_FLUSHPOLICY},
	"report":  {nl.XFRM_MSG_REPORT},
}

func (cmd *cmd) xfrm() error {
	switch cmd.findPrefix("state", "policy", "monitor", "help") {
	case "state":
		return cmd.xfrmState()
	case "policy":
		return cmd.xfrmPolicy()
	case "monitor":
		return cmd.xfrmMonitor()
	case "help":
		fmt.Fprint(cmd.Out, xfrmHelp)
		return nil
	default:
		return cmd.usage()
	}
}

func (cmd *cmd) parseXfrmProto() (netlink.Proto, error) {
	switch c := cmd.nextToken("esp", "ah", "comp", "route2", "hao"); c {
	case "esp":
		return netlink.XFRM_PROTO_ESP, nil
	case "ah":
		return netlink.XFRM_PROTO_AH, nil
	case "comp":
		return netlink.XFRM_PROTO_COMP, nil
	case "route2":
		return netlink.XFRM_PROTO_ROUTE2, nil
	case "hao":
		return netlink.XFRM_PROTO_HAO, nil
	default:
		return netlink.XFRM_PROTO_IPSEC_ANY, cmd.usage()
	}
}

func (cmd *cmd) parseXfrmMode() (netlink.Mode, error) {
	switch c := cmd.nextToken("transport", "tunnel", "ro", "in_trigger", "beet"); c {
	case "transport":
		return netlink.XFRM_MODE_TRANSPORT, nil
	case "tunnel":
		return netlink.XFRM_MODE_TUNNEL, nil
	case "ro":
		return netlink.XFRM_MODE_ROUTEOPTIMIZATION, nil
	case "in_trigger":
		return netlink.XFRM_MODE_IN_TRIGGER, nil
	case "beet":
		return netlink.XFRM_MODE_BEET, nil
	default:
		return netlink.XFRM_MODE_MAX, cmd.usage()
	}
}

func (cmd *cmd) parseXfrmDir() (netlink.Dir, error) {
	switch c := cmd.findPrefix("in", "out", "fwd"); c {
	case "in":
		return netlink.XFRM_DIR_IN, nil
	case "out":
		return netlink.XFRM_DIR_OUT, nil
	case "fwd":
		return netlink.XFRM_DIR_FWD, nil
	default:
		return netlink.XFRM_DIR_IN, cmd.usage()
	}
}

func (cmd *cmd) parseXfrmAction() (netlink.PolicyAction, error) {
	switch c := cmd.findPrefix("allow", "block"); c {
	case "allow":
		return netlink.XFRM_POLICY_ALLOW, nil
	case "block":
		return netlink.XFRM_POLICY_BLOCK, nil
	default:
		return netlink.XFRM_POLICY_ALLOW, cmd.usage()
	}
}

func (cmd *cmd) parseXfrmMark() (*netlink.XfrmMark, error) {
	mark, err := cmd.parseUint32("MARK")
	if err != nil {
		return nil, err
	}

	if !cmd.tokenRemains() {
		return &netlink.XfrmMark{Value: mark}, nil
	}

	// mask is optional
	if cmd.nextToken() != "mask" {
		cmd.lastToken("MARK")
		return &netlink.XfrmMark{Value: mark}, nil
	}

	mask, err := cmd.parseUint32("MASK")
	if err != nil {
		return nil, err
	}

	return &netlink.XfrmMark{Value: mark, Mask: mask}, nil
}

func (cmd *cmd) parseXfrmEncap() (*netlink.XfrmStateEncap, error) {
	var (
		encap netlink.XfrmStateEncap
		err   error
	)

	switch cmd.nextToken("espinudp", "espinudp-nonike", "espintcp") {
	case "espinudp":
		encap.Type = netlink.XFRM_ENCAP_ESPINUDP
	case "espinudp-nonike":
		encap.Type = netlink.XFRM_ENCAP_ESPINUDP_NONIKE
	case "espintcp":
		return nil, fmt.Errorf("espintcp not supported yet")
	}

	encap.SrcPort, err = cmd.parseInt("SPORT")
	if err != nil {
		return nil, err
	}

	encap.DstPort, err = cmd.parseInt("DPORT")
	if err != nil {
		return nil, err
	}

	encap.OriginalAddress, err = cmd.parseAddress()
	if err != nil {
		return nil, err
	}

	return &encap, nil
}

func (cmd *cmd) parseXfrmLimit() (netlink.XfrmStateLimits, error) {
	var (
		err    error
		limits netlink.XfrmStateLimits
	)

	for cmd.tokenRemains() {
		switch c := cmd.nextToken("time-soft", "time-hard", "time-use-soft", "time-use-hard", "byte-soft", "byte-hard", "packet-soft", "packet-hard"); c {
		case "time-soft":
			limits.TimeSoft, err = cmd.parseUint64("SECONDS")
			if err != nil {
				return netlink.XfrmStateLimits{}, err
			}
		case "time-hard":
			limits.TimeHard, err = cmd.parseUint64("SECONDS")
			if err != nil {
				return netlink.XfrmStateLimits{}, err
			}
		case "time-use-soft":
			limits.TimeUseSoft, err = cmd.parseUint64("SECONDS")
			if err != nil {
				return netlink.XfrmStateLimits{}, err
			}
		case "time-use-hard":
			limits.TimeUseHard, err = cmd.parseUint64("SECONDS")
			if err != nil {
				return netlink.XfrmStateLimits{}, err
			}
		case "byte-soft":
			limits.ByteSoft, err = cmd.parseUint64("SIZE")
			if err != nil {
				return netlink.XfrmStateLimits{}, err
			}
		case "byte-hard":
			limits.ByteHard, err = cmd.parseUint64("SIZE")
			if err != nil {
				return netlink.XfrmStateLimits{}, err
			}
		case "packet-soft":
			limits.PacketSoft, err = cmd.parseUint64("COUNT")
			if err != nil {
				return netlink.XfrmStateLimits{}, err
			}
		case "packet-hard":
			limits.PacketHard, err = cmd.parseUint64("COUNT")
			if err != nil {
				return netlink.XfrmStateLimits{}, err
			}
		default:
			cmd.lastToken("LIMITS")
			return limits, nil
		}
	}

	return limits, nil
}

func (cmd *cmd) xfrmMonitor() error {
	updates := make(chan netlink.XfrmMsg)
	errChan := make(chan error)
	done := make(chan struct{})
	defer close(done)

	var filter []nl.XfrmMsgType

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for cmd.tokenRemains() {
		switch v := cmd.nextToken("all", "acquire", "expire", "SA", "aevent", "policy", "help"); v {
		case "help":
			fmt.Fprint(cmd.Out, xfrmMonitorHelp)
			return nil
		case "all":
			for _, v := range xfrmFilterMap {
				filter = append(filter, v...)
			}
		case "acquire", "expire", "SA", "aevent", "policy":
			filter = append(filter, xfrmFilterMap[v]...)
		default:
			return cmd.usage()

		}
	}

	if len(filter) == 0 {
		for _, v := range xfrmFilterMap {
			filter = append(filter, v...)
		}
	}

	// TODO: implement msg types besides nl.XFRM_MSG_EXPIRE for xfrm
	if err := netlink.XfrmMonitor(updates, done, errChan, filter...); err != nil {
		return err
	}

	for {
		select {
		case msg := <-updates:
			switch msg.Type() {
			case nl.XFRM_MSG_EXPIRE:
				msg, ok := msg.(*netlink.XfrmMsgExpire)
				if !ok {
					return fmt.Errorf("invalid type %T", msg)
				}

				printXfrmMsgExpire(cmd.Out, msg.XfrmState)
			default:
				fmt.Fprintf(cmd.Out, "unsupported msg type: %x", msg.Type())
			}

		case err := <-errChan:
			return err

		case <-sig:
			return nil
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func printXfrmMsgExpire(w io.Writer, msg *netlink.XfrmState) {
	fmt.Fprintf(w, "src %s dst %s\n", msg.Src, msg.Dst)
	fmt.Fprintf(w, "    proto %s spi %d reqid %d mode %s\n", msg.Proto, msg.Spi, msg.Reqid, msg.Mode)
	fmt.Fprintf(w, "    replay-window %d\n", msg.ReplayWindow)
	fmt.Fprintf(w, "    auth-trunc %s %s %d\n", msg.Auth.Name, msg.Auth.Key, msg.Auth.TruncateLen)
	fmt.Fprintf(w, "    enc %s %s\n", msg.Crypt.Name, msg.Crypt.Key)
	fmt.Fprintf(w, "    sel src %s dst %s\n", msg.Src, msg.Dst)
	fmt.Fprintf(w, "    lifetime config:\n")
	fmt.Fprintf(w, "      limit: soft (%d)(bytes), hard (%d)(bytes)\n", msg.Limits.ByteSoft, msg.Limits.ByteHard)
	fmt.Fprintf(w, "      limit: soft (%d)(packets), hard (%d)(packets)\n", msg.Limits.PacketSoft, msg.Limits.PacketHard)
	fmt.Fprintf(w, "      expire add: soft %d(sec), hard %d(sec)\n", msg.Limits.TimeSoft, msg.Limits.TimeHard)
	fmt.Fprintf(w, "      expire use: soft %d(sec), hard %d(sec)\n", msg.Limits.TimeUseSoft, msg.Limits.TimeUseHard)
	fmt.Fprintf(w, "    lifetime current:\n")
	fmt.Fprintf(w, "      %d(bytes), %d(packets)\n", msg.Statistics.Bytes, msg.Statistics.Packets)
	fmt.Fprintf(w, "      add %d, use %d\n", msg.Statistics.AddTime, msg.Statistics.UseTime)
	fmt.Fprintf(w, "    stats:\n")
	fmt.Fprintf(w, "      replay-window %d replay %d failed %d\n", msg.Statistics.ReplayWindow, msg.Statistics.Replay, msg.Statistics.Failed)
}
