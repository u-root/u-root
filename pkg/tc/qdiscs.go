// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trafficctl

import (
	"fmt"
	"io"
	"strconv"

	"github.com/florianl/go-tc"
)

const CodelHelp = `Usage: ... codel [ limit PACKETS ] [ target TIME ]
		 [ interval TIME ] [ ecn | noecn ]
		 [ ce_threshold TIME ]
`

// ParseCodelArgs parses a []string from the commandline for the codel qdisc.
// and returns an *tc.Object accordingly.
func ParseCodelArgs(out io.Writer, args []string) (*tc.Object, error) {
	codel := &tc.Codel{}
	for i := 0; i < len(args); i = i + 2 {
		switch args[i] {
		case "limit":
			val, err := strconv.ParseUint(args[i+1], 10, 32)
			if err != nil {
				return nil, err
			}
			indirect := uint32(val)
			codel.Limit = &indirect
		case "target":
			val, err := strconv.ParseUint(args[i+1], 10, 32)
			if err != nil {
				return nil, err
			}
			indirect := uint32(val)
			codel.Target = &indirect
		case "interval":
			// This is a time value with units (ms)
			val, err := parseTime(args[i+1])
			if err != nil {
				return nil, err
			}
			codel.Interval = &val
		case "ce_threshold":
			// This is a time value with units (ms)
			val, err := parseTime(args[i+1])
			if err != nil {
				return nil, err
			}
			codel.CEThreshold = &val
		case "ecn":
			on := uint32(1)
			codel.ECN = &on
			i--
		case "noecn":
			off := uint32(0)
			codel.ECN = &off
			i--
		case "help":
			fmt.Fprintf(out, "%s", CodelHelp)
			return nil, ErrExitAfterHelp
		}
	}
	ret := &tc.Object{}
	ret.Kind = "codel"
	ret.Codel = codel
	return ret, nil
}

const QFQHelp = `Usage: ... qfq [ weight N ] [ maxpkt N ]
`

// ParseQFQArgs parses a []string from the commandline for the qfq qdisc
// via `tc qdisc ... qfq ...` and returns an *tc.Object accordingly.
func ParseQFQArgs(out io.Writer, args []string) (*tc.Object, error) {
	qfq := &tc.Qfq{}

	for i := 0; i < len(args); i = i + 2 {
		switch args[i] {
		case "weight":
			val, err := strconv.ParseUint(args[i+1], 10, 32)
			if err != nil {
				return nil, err
			}
			indirect := uint32(val)
			qfq.Weight = &indirect
		case "maxpkt":
			val, err := strconv.ParseUint(args[i+1], 10, 32)
			if err != nil {
				return nil, err
			}
			indirect := uint32(val)
			qfq.Lmax = &indirect
		case "help":
			fmt.Fprintf(out, "%s\n", QFQHelp)
			return nil, ErrExitAfterHelp
		}
	}

	ret := &tc.Object{}
	ret.Kind = "qfq"
	ret.Qfq = qfq

	return ret, nil
}

const HTBHelp = `Usage: ... qdisc add ... htb [default N] [r2q N]
	[direct_qlen P]

default  minor id of class to which unclassified packets are sent {0}
r2q      DRR quantums are computed as rate in Bps/r2q {10}
direct_qlen  Limit of the direct queue {in packets}
offload  enable hardware offload

class add ... htb rate R1 [burst B1] [mpu B] [overhead O]
	[prio P] [slot S] [pslot PS]
	[ceil R2] [cburst B2] [mtu MTU] [quantum Q]
rate     rate allocated to this class (class can still borrow)
burst    max bytes burst which can be accumulated during idle period {computed}
mpu      minimum packet size used in rate computations
overhead per-packet size overhead used in rate computations
linklay  adapting to a linklayer e.g. atm
ceil     definite upper class rate (no borrows) {rate}
cburst   burst but for ceil {computed}
mtu      max packet size we create rate map for {1600}
prio     priority of leaf; lower are served first {0}
quantum  how much bytes to serve from leaf at once {use r2q}
`

// ParseHTBQDiscArgs parses a []string from the commandline for the HTB qdisc
// via `tc qdisc ... htb ...` and returns an *tc.Object accordingly.
func ParseHTBQDiscArgs(out io.Writer, args []string) (*tc.Object, error) {
	htb := tc.Htb{
		Init: &tc.HtbGlob{
			Rate2Quantum: 10,
			Version:      3,
		},
	}
	for i := 0; i < len(args); i = i + 2 {
		switch args[i] {
		case "default":
			defcls, err := strconv.ParseUint(args[i+1], 16, 32)
			if err != nil {
				return nil, err
			}
			htb.Init.Defcls = uint32(defcls)
		case "r2q":
			r2q, err := strconv.ParseUint(args[i+1], 10, 32)
			if err != nil {
				return nil, err
			}
			htb.Init.Rate2Quantum = uint32(r2q)
		case "debug":
			return nil, ErrNotImplemented
		case "direct_qlen":
			dq, err := strconv.ParseUint(args[i+1], 10, 32)
			if err != nil {
				return nil, err
			}
			indirect := uint32(dq)
			htb.DirectQlen = &indirect
		case "offload":
			return nil, ErrNotImplemented
		case "help":
			fmt.Fprint(out, HTBHelp)
		}
	}

	ret := &tc.Object{}
	ret.Attribute.Htb = &htb
	ret.Kind = "htb"
	return ret, nil
}

const HFSCHelp = `Usage: ... hfsc [ default CLASSID ]

 default: default class for unclassified packets
 `

// ParseHFSCQDiscArgs parses a []string from the commandline for the HFSC qdisc via
// `tc qdisc ... hfsc ...` and returns an *tc.Object accordingly.
func ParseHFSCQDiscArgs(stdout io.Writer, args []string) (*tc.Object, error) {
	ret := &tc.Object{}
	hfsc := &tc.HfscQOpt{}

	for i := 0; i < len(args); i = i + 2 {
		switch args[i] {
		case "default":
			defcls, err := strconv.ParseUint(args[i+1], 10, 16)
			if err != nil {
				return nil, err
			}
			hfsc.DefCls = uint16(defcls)
		default:
			// WHAT IS THIS?
			return nil, ErrInvalidArg
		}
	}

	ret.HfscQOpt = hfsc

	return ret, nil
}
