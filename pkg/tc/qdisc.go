// Copyright 2012-20124 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trafficctl

import (
	"fmt"
	"io"
	"net"
	"strconv"

	"github.com/florianl/go-tc"
	"github.com/florianl/go-tc/core"
	"golang.org/x/sys/unix"
)

// QArgs holds all possible args for qdisc subcommand
// tc qdisc [ add | del | replace | change | show ] dev STRING
// [ handle QHANDLE ] [ root | ingress | clsact | parent CLASSID ]
// [ estimator INTERVAL TIME_CONSTANT ]
// [ stab [ help | STAB_OPTIONS] ]
// [ ingress_block BLOCK_INDEX ] [ egress_block BLOCK_INDEX ]
// [ [ QDISC_KIND ] [ help | OPTIONS ] ]
type QArgs struct {
	dev    *string
	kind   *string
	parent *uint32
	handle *uint32
	obj    *tc.Object
}

func ParseQDiscArgs(args []string, stdout io.Writer) (*QArgs, error) {
	ret := &QArgs{}
	if len(args) < 1 {
		return nil, ErrNotEnoughArgs
	}

	for i := 0; i < len(args); i = i + 2 {
		var val string
		if len(args[1:]) > i {
			val = args[i+1]
		}

		switch args[i] {
		case "dev":
			ret.dev = &val
		case "handle":
			handle, err := strconv.Atoi(val)
			if err != nil {
				return nil, err
			}
			if handle < 0 || handle >= 0x7FFFFFFF {
				return nil, ErrOutOfBounds
			}
			indirect := uint32(handle)
			ret.handle = &indirect
		case "root":
			if ret.parent != nil {
				return nil, ErrInvalidArg
			}
			indirect := tc.HandleRoot
			ret.parent = &indirect
			// We have a one piece argument. To get to the next arg properly
			i--
		case "ingress":
			k := "ingress"
			ret.kind = &k
			if ret.parent != nil {
				return nil, ErrInvalidArg
			}
			indirectPar := tc.HandleIngress // is the same as clsact handle
			ret.parent = &indirectPar
			// We have a one piece argument. To get to the next arg properly
			indirectHan := core.BuildHandle(tc.HandleIngress, 0)
			ret.handle = &indirectHan

			i--
		case "clsact":
			k := "clsact"
			ret.kind = &k

			if ret.parent != nil {
				return nil, ErrInvalidArg
			}

			indirectPar := tc.HandleIngress // is the same as clsact handle
			ret.parent = &indirectPar

			indirectHan := core.BuildHandle(tc.HandleIngress, 0)
			ret.handle = &indirectHan
			i--
		case "parent":
			qdiscID, err := strconv.Atoi(val)
			if err != nil {
				return ret, err
			}
			if qdiscID < 0 || qdiscID >= 0x7FFFFFFF {
				return nil, ErrOutOfBounds
			}
			indirect := uint32(qdiscID)
			ret.parent = &indirect
		case "estimator":
			return nil, ErrNotImplemented
		case "stab":
			return nil, ErrNotImplemented
		case "ingress_block":
			return nil, ErrNotImplemented
		case "egress_block":
			return nil, ErrNotImplemented
		case "help":
			PrintQdiscHelp(stdout)
		default:
			// Resolve Qdisc and parameters
			var qdiscParseFn func([]string) (*tc.Object, error)
			if qdiscParseFn = supportetQdisc(args[i]); qdiscParseFn == nil {
				return nil, fmt.Errorf("%w: invalid qdisc: %s", ErrInvalidArg, args[i])
			}
			var err error
			ret.obj, err = qdiscParseFn(args[i+1:])
			if err != nil {
				return nil, err
			}

		}
	}

	return ret, nil
}

const (
	QdiscHelp = `Usage:
	tc qdisc [ add | del | replace | change | show ] dev STRING
   		[ handle QHANDLE ] [ root | ingress | clsact | parent CLASSID ]
   		[ estimator INTERVAL TIME_CONSTANT ]
  		[ stab [ help | STAB_OPTIONS] ]
  		[ ingress_block BLOCK_INDEX ] [ egress_block BLOCK_INDEX ]
  		[ [ QDISC_KIND ] [ help | OPTIONS ] ]

	tc qdisc { show | list } [ dev STRING ] [ QDISC_ID ] [ invisible ]

	Where:
	QDISC_KIND := { [p|b]fifo | tbf | prio | red | etc. }
	OPTIONS := ... try tc qdisc add <desired QDISC_KIND> help
	STAB_OPTIONS := ... try tc qdisc add stab help
	QDISC_ID := { root | ingress | handle QHANDLE | parent CLASSID }`
)

func PrintQdiscHelp(stdout io.Writer) {
	fmt.Fprintf(stdout, "%s", QdiscHelp)
}

func (t *Trafficctl) ShowQdisc(stdout io.Writer) error {
	qdiscs, err := t.Tc.Qdisc().Get()
	if err != nil {
		return err
	}

	for _, qdisc := range qdiscs {
		iface, err := net.InterfaceByIndex(int(qdisc.Ifindex))
		if err != nil {
			return err
		}
		fmt.Fprintf(stdout, "%20s\t%s\n", iface.Name, qdisc.Kind)
	}
	return nil
}

func (t *Trafficctl) AddQdisc(args *QArgs, stdout io.Writer) error {
	iface, err := net.InterfaceByName(*args.dev)
	if err != nil {
		return err
	}

	msg := tc.Msg{
		Family:  unix.AF_UNSPEC,
		Ifindex: uint32(iface.Index),
		Parent:  *args.parent,
		Handle:  *args.handle,
	}

	obj := &tc.Object{
		Msg: msg,
		Attribute: tc.Attribute{
			Kind: *args.kind,
		},
	}

	if err := t.Tc.Qdisc().Add(obj); err != nil {
		return fmt.Errorf("Qdisc.Add() = %v", err)
	}
	return nil
}

func (t *Trafficctl) DelQdisc(args *QArgs, stdout io.Writer) error {
	iface, err := net.InterfaceByName(*args.dev)
	if err != nil {
		return err
	}

	msg := tc.Msg{
		Family:  unix.AF_UNSPEC,
		Ifindex: uint32(iface.Index),
		Parent:  *args.parent,
		Handle:  *args.handle,
	}
	obj := &tc.Object{
		Msg: msg,
		Attribute: tc.Attribute{
			Kind: *args.kind,
		},
	}

	if err := t.Tc.Qdisc().Delete(obj); err != nil {
		return fmt.Errorf("Qdisc.Delete() = %v", err)
	}

	return nil
}

func (t *Trafficctl) ReplaceQdisc(args *QArgs, stdout io.Writer) error {
	iface, err := net.InterfaceByName(*args.dev)
	if err != nil {
		return err
	}

	msg := tc.Msg{
		Family:  unix.AF_UNSPEC,
		Ifindex: uint32(iface.Index),
		Parent:  *args.parent,
	}
	obj := &tc.Object{
		Msg: msg,
		Attribute: tc.Attribute{
			Kind: *args.kind,
		},
	}

	if err := t.Tc.Qdisc().Replace(obj); err != nil {
		return fmt.Errorf("Qdisc.Delete() = %v", err)
	}
	return nil
}

func (t *Trafficctl) ChangeQDisc(args *QArgs, stdout io.Writer) error {
	iface, err := net.InterfaceByName(*args.dev)
	if err != nil {
		return err
	}

	msg := tc.Msg{
		Family:  unix.AF_UNSPEC,
		Ifindex: uint32(iface.Index),
		Parent:  *args.parent,
	}
	obj := &tc.Object{
		Msg: msg,
		Attribute: tc.Attribute{
			Kind: *args.kind,
		},
	}

	if err := t.Tc.Qdisc().Change(obj); err != nil {
		return fmt.Errorf("Qdisc.Delete() = %v", err)
	}

	return nil
}

func (t *Trafficctl) LinkQDisc(args *QArgs, stdout io.Writer) error {
	iface, err := net.InterfaceByName(*args.dev)
	if err != nil {
		return err
	}

	msg := tc.Msg{
		Family:  unix.AF_UNSPEC,
		Ifindex: uint32(iface.Index),
		Parent:  *args.parent,
	}
	obj := &tc.Object{
		Msg: msg,
		Attribute: tc.Attribute{
			Kind: *args.kind,
		},
	}

	if err := t.Tc.Qdisc().Change(obj); err != nil {
		return fmt.Errorf("Qdisc.Delete() = %v", err)
	}
	return nil
}

func supportetQdisc(qd string) func([]string) (*tc.Object, error) {
	supported := map[string]func([]string) (*tc.Object, error){
		// Classless qdiscs
		"cake":       nil,
		"choke":      nil,
		"codel":      ParseCodelArgs,
		"pfifo":      nil,
		"qfifo":      nil,
		"fq":         nil,
		"fq_codel":   nil,
		"fq_pie":     nil,
		"gred":       nil,
		"hhf":        nil,
		"ingress":    nil,
		"mqprio":     nil,
		"multiq":     nil,
		"netem":      nil,
		"pfifo_fast": nil,
		"pie":        nil,
		"red":        nil,
		"sfb":        nil,
		"sfq":        nil,
		"tbf":        nil,
		// Classful qdiscs
		"cbs":      nil,
		"htb":      nil,
		"hfsc":     nil,
		"hfscqopt": nil,
		"dsmark":   nil,
		"drr":      nil,
		"cbq":      nil,
		"atm":      nil,
		"qfq":      nil,
		"taprio":   nil,
	}

	ret, ok := supported[qd]
	if !ok {
		return nil
	}

	return ret
}

func ParseCodelArgs(args []string) (*tc.Object, error) {
	codel := &tc.Codel{}
	for i := 0; i < len(args); i = i + 2 {
		switch args[i] {
		case "limit":
			val, err := strconv.Atoi(args[i+1])
			if err != nil {
				return nil, err
			}
			if val < 0x0 || val >= 0x7FFFFFFF {
				return nil, ErrOutOfBounds
			}
			indirect := uint32(val)
			codel.Limit = &indirect
		case "target":
			val, err := strconv.Atoi(args[i+1])
			if err != nil {
				return nil, err
			}
			if val < 0x0 || val >= 0x7FFFFFFF {
				return nil, ErrOutOfBounds
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
		}
	}
	ret := &tc.Object{}
	ret.Kind = "codel"
	ret.Codel = codel
	return ret, nil
}
