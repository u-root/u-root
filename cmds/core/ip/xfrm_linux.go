// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
where  XFRM-OBJECT := state | policy | monitor`

	xfrmMonitorHelp = `Usage: ip xfrm monitor [ nokeys ] [ all-nsid ] [ all | OBJECTS | help ]
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

func xfrm(w io.Writer) error {
	cursor++

	expectedValues = []string{"state", "policy", "monitor"}
	switch findPrefix(arg[cursor], expectedValues) {
	case "state":
		return xfrmState(w)
	case "monitor":
		return xfrmMonitor(w)
	case "help":
		fmt.Fprint(w, xfrmHelp)
		return nil
	default:
		return usage()
	}
}

func parseXfrmProto() (netlink.Proto, error) {
	cursor++
	expectedValues = []string{"esp", "ah", "comp", "route2", "hao"}

	switch arg[cursor] {
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
		return netlink.XFRM_PROTO_IPSEC_ANY, fmt.Errorf("invalid proto %s", arg[cursor])
	}
}

func parseXfrmMode() (netlink.Mode, error) {
	cursor++
	expectedValues = []string{"esp", "ah", "comp", "route2", "hao", "ipsec-any"}

	switch arg[cursor] {
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
		return netlink.XFRM_MODE_MAX, fmt.Errorf("invalid mode %s", arg[cursor])
	}
}

func parseXfrmDir() (netlink.Dir, error) {
	cursor++
	expectedValues = []string{"in", "out", "fwd"}

	switch arg[cursor] {
	case "in":
		return netlink.XFRM_DIR_IN, nil
	case "out":
		return netlink.XFRM_DIR_OUT, nil
	case "fwd":
		return netlink.XFRM_DIR_FWD, nil
	default:
		return netlink.XFRM_DIR_IN, fmt.Errorf("invalid mode %s", arg[cursor])
	}
}

func parseXfrmAction() (netlink.PolicyAction, error) {
	cursor++
	expectedValues = []string{"allow", "block"}

	switch arg[cursor] {
	case "allow":
		return netlink.XFRM_POLICY_ALLOW, nil
	case "block":
		return netlink.XFRM_POLICY_BLOCK, nil
	default:
		return netlink.XFRM_POLICY_ALLOW, fmt.Errorf("invalid mode %s", arg[cursor])
	}
}

func parseXfrmMark() (*netlink.XfrmMark, error) {
	cursor++
	expectedValues = []string{"MARK"}

	mark, err := parseUint32("MARK")
	if err != nil {
		return nil, err
	}

	cursor++
	if len(arg) == cursor {
		return &netlink.XfrmMark{Value: mark}, nil
	}

	// mask is optional
	if arg[cursor] != "mask" {
		cursor--
		return &netlink.XfrmMark{Value: mark}, nil
	}

	mask, err := parseUint32("MASK")
	if err != nil {
		return nil, err
	}

	return &netlink.XfrmMark{Value: mark, Mask: mask}, nil
}

func parseXfrmEncap() (*netlink.XfrmStateEncap, error) {
	var (
		encap netlink.XfrmStateEncap
		err   error
	)

	cursor++
	expectedValues = []string{"espinudp", "espinudp-nonike", "espintcp"}

	switch arg[cursor] {
	case "espinudp":
		encap.Type = netlink.XFRM_ENCAP_ESPINUDP
	case "espinudp-nonike":
		encap.Type = netlink.XFRM_ENCAP_ESPINUDP_NONIKE
	case "espintcp":
		return nil, fmt.Errorf("espintcp not supported yet")
	}

	cursor++
	expectedValues = []string{"SPORT"}
	encap.SrcPort, err = parseInt("SPORT")
	if err != nil {
		return nil, err
	}

	cursor++
	expectedValues = []string{"DPORT"}
	encap.DstPort, err = parseInt("DPORT")
	if err != nil {
		return nil, err
	}

	cursor++
	expectedValues = []string{"OADDR"}
	encap.OriginalAddress, err = parseAddress()
	if err != nil {
		return nil, err
	}

	return &encap, nil
}

func parseXfrmLimit() (netlink.XfrmStateLimits, error) {
	var (
		err    error
		limits netlink.XfrmStateLimits
	)

	cursor++
	expectedValues = []string{"time-soft", "time-hard", "time-use-soft", "time-use-hard", "byte-soft", "byte-hard", "packet-soft", "packet-hard"}

	switch arg[cursor] {
	case "time-soft":
		limits.TimeSoft, err = parseUint64("SECONDS")
		if err != nil {
			return netlink.XfrmStateLimits{}, err
		}
	case "time-hard":
		limits.TimeHard, err = parseUint64("SECONDS")
		if err != nil {
			return netlink.XfrmStateLimits{}, err
		}
	case "time-use-soft":
		limits.TimeUseSoft, err = parseUint64("SECONDS")
		if err != nil {
			return netlink.XfrmStateLimits{}, err
		}
	case "time-use-hard":
		limits.TimeUseHard, err = parseUint64("SECONDS")
		if err != nil {
			return netlink.XfrmStateLimits{}, err
		}
	case "byte-soft":
		limits.ByteSoft, err = parseUint64("SIZE")
		if err != nil {
			return netlink.XfrmStateLimits{}, err
		}
	case "byte-hard":
		limits.ByteHard, err = parseUint64("SIZE")
		if err != nil {
			return netlink.XfrmStateLimits{}, err
		}
	case "packet-soft":
		limits.PacketSoft, err = parseUint64("COUNT")
		if err != nil {
			return netlink.XfrmStateLimits{}, err
		}
	case "packet-hard":
		limits.PacketHard, err = parseUint64("COUNT")
		if err != nil {
			return netlink.XfrmStateLimits{}, err
		}
	default:
		return netlink.XfrmStateLimits{}, fmt.Errorf("unknown limit option %s", arg[cursor])
	}

	return limits, nil
}

func xfrmMonitor(w io.Writer) error {
	updates := make(chan netlink.XfrmMsg)
	errChan := make(chan error)
	done := make(chan struct{})
	defer close(done)

	var filter []nl.XfrmMsgType

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	expectedValues = []string{"all", "acquire", "expire", "SA", "aevent", "policy", "help"}
	for {
		cursor++

		if cursor == len(arg) {
			break
		}

		switch v := arg[cursor]; v {
		case "help":
			fmt.Fprint(w, xfrmMonitorHelp)
			return nil
		case "all":
			for _, v := range xfrmFilterMap {
				filter = append(filter, v...)
			}
		case "acquire", "expire", "SA", "aevent", "policy":
			filter = append(filter, xfrmFilterMap[v]...)
		default:
			return usage()

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

				printXfrmMsgExpire(w, msg.XfrmState)
			default:
				fmt.Fprintf(w, "unsupported msg type: %x", msg.Type())
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
