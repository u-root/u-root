// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trafficctl

import (
	"fmt"
	"io"
	"net"

	"github.com/florianl/go-tc"
	"github.com/florianl/go-tc/core"
	"golang.org/x/sys/unix"
)

// Args holds all possible args for qdisc subcommand
// tc qdisc [ add | del | replace | change | show ] dev STRING
// [ handle QHANDLE ] [ root | ingress | clsact | parent CLASSID ]
// [ estimator INTERVAL TIME_CONSTANT ]
// [ stab [ help | STAB_OPTIONS] ]
// [ ingress_block BLOCK_INDEX ] [ egress_block BLOCK_INDEX ]
// [ [ QDISC_KIND ] [ help | OPTIONS ] ]
type Args struct {
	dev    string
	kind   string
	parent *uint32
	handle *uint32
	obj    *tc.Object
}

// ParseQDiscArgs takes an io.Writer and []string slice with arguments to parse.
// It returns a structure of type Args for qdisc operation.
func ParseQdiscArgs(stdout io.Writer, args []string) (*Args, error) {
	ret := &Args{
		obj: &tc.Object{},
	}
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
			ret.dev = val
		case "handle":
			indirect, err := ParseHandle(val)
			if err != nil {
				return nil, err
			}
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
			ret.obj.Kind = "ingress"
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
			ret.obj.Kind = "clsact"

			if ret.parent != nil {
				return nil, ErrInvalidArg
			}

			indirectPar := tc.HandleIngress // is the same as clsact handle
			ret.parent = &indirectPar

			indirectHan := core.BuildHandle(tc.HandleIngress, 0)
			ret.handle = &indirectHan
			i--
		case "parent":
			parent, err := ParseClassID(val)
			if err != nil {
				return nil, err
			}
			indirect := uint32(parent)
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
			fmt.Fprint(stdout, QdiscHelp)
		default:
			var qdiscParseFn func(io.Writer, []string) (*tc.Object, error)
			if qdiscParseFn = supportetQdisc(args[i]); qdiscParseFn == nil {
				return nil, fmt.Errorf("%w: invalid qdisc: %s", ErrInvalidArg, args[i])
			}
			var err error
			ret.obj, err = qdiscParseFn(stdout, args[i+1:])
			if err != nil {
				return nil, err
			}
			ret.kind = ret.obj.Kind

			return ret, nil
		}
	}

	return ret, nil
}

// Originally from tc:
// Usage: tc qdisc [ add | del | replace | change | show ] dev STRING
//        [ handle QHANDLE ] [ root | ingress | clsact | parent CLASSID ]
//        [ estimator INTERVAL TIME_CONSTANT ]
//        [ stab [ help | STAB_OPTIONS] ]
//        [ ingress_block BLOCK_INDEX ] [ egress_block BLOCK_INDEX ]
//        [ [ QDISC_KIND ] [ help | OPTIONS ] ]

//        tc qdisc { show | list } [ dev STRING ] [ QDISC_ID ] [ invisible ]
// Where:
// QDISC_KIND := { [p|b]fifo | tbf | prio | cbq | red | etc. }
// OPTIONS := ... try tc qdisc add <desired QDISC_KIND> help
// STAB_OPTIONS := ... try tc qdisc add stab help
// QDISC_ID := { root | ingress | handle QHANDLE | parent CLASSID }

const QdiscHelp = `Usage: tc qdisc [ add | del | replace | change | show ] dev STRING
   		[ handle QHANDLE ] [ root | ingress | clsact | parent CLASSID ]
  		[ [ QDISC_KIND ] [ help | OPTIONS ] ]

	tc qdisc show [ dev STRING ] [ QDISC_ID ]

Where:
	QDISC_KIND := { codel | qfq | htb | hfsc }
	OPTIONS := ... try tc qdisc add <desired QDISC_KIND> help
	QDISC_ID := { root | ingress | handle QHANDLE | parent CLASSID }
`

// ShowQdisc implements the functionality of `tc qdisc show ... `
func (t *Trafficctl) ShowQdisc(stdout io.Writer, args *Args) error {
	qdiscs, err := t.Tc.Qdisc().Get()
	if err != nil {
		return err
	}

	for _, qdisc := range qdiscs {
		iface, err := net.InterfaceByIndex(int(qdisc.Ifindex))
		if err != nil {
			return err
		}

		if args.dev == "" || args.dev == iface.Name {
			fmt.Fprintf(stdout, "%20s\tqdisc %s %s %s",
				iface.Name, qdisc.Kind,
				RenderClassID(qdisc.Handle, false),
				RenderClassID(qdisc.Parent, true),
			)
			if qdisc.Kind == "htb" {
				htb := qdisc.Htb
				if htb.Init != nil {
					fmt.Fprintf(stdout, " r2q %d default 0x%x", htb.Init.Rate2Quantum, htb.Init.Defcls)
				}
				if htb.DirectQlen != nil {
					fmt.Fprintf(stdout, " direct_qlen %d", *htb.DirectQlen)
				}
			}
			fmt.Fprintf(stdout, "\n")
		}
	}
	return nil
}

// AddQdisc implements the functionality of `tc qdisc add ... `
func (t *Trafficctl) AddQdisc(stdout io.Writer, args *Args) error {
	iface, err := getDevice(args.dev)
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
		Msg:       msg,
		Attribute: args.obj.Attribute,
	}

	if err := t.Tc.Qdisc().Add(obj); err != nil {
		return fmt.Errorf("Qdisc.Add() = %w", err)
	}
	return nil
}

// DeleteQdisc implements the functionality of `tc qdisc del ... `
func (t *Trafficctl) DeleteQdisc(stdout io.Writer, args *Args) error {
	iface, err := getDevice(args.dev)
	if err != nil {
		return err
	}

	qdiscs, err := t.Tc.Qdisc().Get()
	if err != nil {
		return err
	}

	var q tc.Object
	var found bool
	for _, qdisc := range qdiscs {
		if qdisc.Ifindex == uint32(iface.Index) {
			q = qdisc
			found = true
		}
	}

	if !found {
		return fmt.Errorf("on device '%s' no qdisc '%s' was found: %w", args.dev, args.kind, ErrInvalidArg)
	}

	if err := t.Tc.Qdisc().Delete(&q); err != nil {
		return fmt.Errorf("Qdisc.Delete() = %w", err)
	}

	return nil
}

// ReplaceQdisc implements the functionality of `tc qdisc replace ... `
func (t *Trafficctl) ReplaceQdisc(stdout io.Writer, args *Args) error {
	iface, err := getDevice(args.dev)
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
			Kind: args.obj.Kind,
		},
	}

	if err := t.Tc.Qdisc().Replace(obj); err != nil {
		return fmt.Errorf("Qdisc.Replace() = %w", err)
	}
	return nil
}

// ChangeQdisc implements the functionality of `tc qdisc change ... `
func (t *Trafficctl) ChangeQdisc(stdout io.Writer, args *Args) error {
	iface, err := getDevice(args.dev)
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
			Kind: args.obj.Kind,
		},
	}

	if err := t.Tc.Qdisc().Change(obj); err != nil {
		return fmt.Errorf("Qdisc.Change() = %w", err)
	}

	return nil
}

func supportetQdisc(qd string) func(io.Writer, []string) (*tc.Object, error) {
	supported := map[string]func(io.Writer, []string) (*tc.Object, error){
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
		// QFQ is listed as Classfull QDisk in man page of tc, but tc implementation
		// complains that it is a classless QDisc, so im treating it as classless
		"qfq": ParseQFQArgs,
		"red": nil,
		"sfb": nil,
		"sfq": nil,
		"tbf": nil,
	}

	ret := supported[qd]

	if ret != nil {
		return ret
	}

	return supportetQdiscClassfull(qd)
}

func supportetQdiscClassfull(qd string) func(io.Writer, []string) (*tc.Object, error) {
	supported := map[string]func(io.Writer, []string) (*tc.Object, error){
		// Classful qdiscs
		"cbs":      nil, // (not supported for adding byt go-tc library)
		"htb":      ParseHTBQDiscArgs,
		"hfsc":     ParseHFSCQDiscArgs,
		"hfscqopt": nil, // (not supported for adding byt go-tc library)
		"dsmark":   nil, // (not supported for adding byt go-tc library)
		"drr":      nil, // (not supported for adding byt go-tc library)
		"cbq":      nil,
		"atm":      nil, // (not supported for adding byt go-tc library)
		"taprio":   nil, // (not supported for adding byt go-tc library)
	}

	ret, ok := supported[qd]
	if !ok {
		return nil
	}

	return ret
}
