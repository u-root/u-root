// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trafficctl

import (
	"fmt"
	"io"
	"net"

	"github.com/florianl/go-tc"
	"golang.org/x/sys/unix"
)

// ParseClassArgs takes an io.Writer for output operation and a []string with the provided
// arguments to parse. It builds a struct of type Args for further operation.
// Further more it selects the class and calls into the class related parsing function.
func ParseClassArgs(stdout io.Writer, args []string) (*Args, error) {
	ret := &Args{}
	if len(args) < 1 {
		return nil, fmt.Errorf("ParseClassArgs() = %w", ErrNotEnoughArgs)
	}

	for i := 0; i <= len(args[1:]); i = i + 2 {
		var val string
		if len(args[1:]) > i {
			val = args[i+1]
		}
		switch args[i] {
		case "dev":
			ret.dev = val
		case "parent":
			parent, err := ParseClassID(args[i+1])
			if err != nil {
				return nil, err
			}
			indirect := uint32(parent)
			ret.parent = &indirect
		case "root":
			indirect := tc.HandleRoot
			ret.parent = &indirect
			// We have a one piece argument. To get to the next arg properly
			i--
		case "classid":
			classid, err := ParseClassID(args[i+1])
			if err != nil {
				return nil, err
			}
			indirect := uint32(classid)
			ret.handle = &indirect
		case "help":
			fmt.Fprint(stdout, ClassHelp)
			return nil, nil
		default:
			// Resolve Qdisc and parameters
			var classParse func(io.Writer, []string) (*tc.Object, error)
			if classParse = supportetClasses(args[i]); classParse == nil {
				return nil, fmt.Errorf("%w: invalid class: %s", ErrInvalidArg, args[i])
			}

			var err error
			ret.obj, err = classParse(stdout, args[i+1:])
			if err != nil {
				return nil, err
			}
			return ret, nil
		}
	}

	return ret, nil
}

// ShowClass realizes the `tc class show dev <DEV>` functionality
func (t *Trafficctl) ShowClass(stdout io.Writer, args *Args) error {
	ifs, err := net.Interfaces()
	if err != nil {
		return err
	}

	for _, iface := range ifs {
		if args.dev != "" && iface.Name != args.dev {
			continue
		}
		msg := &tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(iface.Index),
			Handle:  0,
			Parent:  0,
			Info:    0,
		}

		classes, err := t.Tc.Class().Get(msg)
		if err != nil {
			return err
		}

		for _, class := range classes {
			if class.Ifindex == uint32(iface.Index) {
				fmt.Fprintf(stdout, "%20s\tclass %s %s %s",
					iface.Name, class.Kind,
					RenderClassID(class.Handle, false),
					RenderClassID(class.Parent, true),
				)
				if class.Kind == "htb" && class.Htb.Parms != nil {
					parms := class.Htb.Parms

					burst, err := CalcXMitSize(uint64(parms.Rate.Rate), parms.Buffer)
					if err != nil {
						return err
					}

					cburst, err := CalcXMitSize(uint64(parms.Ceil.Rate), parms.Cbuffer)
					if err != nil {
						return err
					}

					fmt.Fprintf(stdout,
						" prio %d rate %db ceil %db burst %db cburst %db",
						parms.Prio, parms.Rate.Rate, parms.Ceil.Rate, burst, cburst,
					)

				}
				fmt.Fprintf(stdout, "\n")
			}
		}
	}

	return nil
}

// AddClass realizes the `tc class add dev <DEV> ... ` functionality
func (t *Trafficctl) AddClass(stdout io.Writer, args *Args) error {
	iface, err := getDevice(args.dev)
	if err != nil {
		return err
	}

	// Get the qdiscs
	qdiscs, err := t.Tc.Qdisc().Get()
	if err != nil {
		return err
	}

	// Look for the same device
	var q tc.Object
	for _, qdisc := range qdiscs {
		if qdisc.Ifindex == uint32(iface.Index) {
			q = qdisc
		}
	}

	q.Attribute = args.obj.Attribute
	q.Handle = *args.handle
	q.Parent = *args.parent

	if err := t.Tc.Class().Add(&q); err != nil {
		return err
	}

	return nil
}

// DeleteClass realizes the `tc class del dev <DEV> ...` functionality
func (t *Trafficctl) DeleteClass(stdout io.Writer, args *Args) error {
	iface, err := getDevice(args.dev)
	if err != nil {
		return err
	}

	msg := tc.Msg{
		Family:  unix.AF_UNSPEC,
		Ifindex: uint32(iface.Index),
		Handle:  *args.handle,
	}

	obj := &tc.Object{
		Msg: msg,
	}

	if err := t.Tc.Class().Delete(obj); err != nil {
		return fmt.Errorf("Class.Delete() = %w", err)
	}

	return nil
}

// ChangeClass implements the changing of a classful qdisc with `tc class change ...`
func (t *Trafficctl) ChangeClass(stdout io.Writer, args *Args) error {
	return ErrNotImplemented
}

// ReplaceClass implements the replacement of a classful qdisc with `tc class replace ...`
func (t *Trafficctl) ReplaceClass(stdout io.Writer, args *Args) error {
	return ErrNotImplemented
}

// Origianlly from tc:
// Usage: tc class [ add | del | change | replace | show ] dev STRING
//        [ classid CLASSID ] [ root | parent CLASSID ]
//        [ [ QDISC_KIND ] [ help | OPTIONS ] ]

//        tc class show [ dev STRING ] [ root | parent CLASSID ]
// Where:
// QDISC_KIND := { prio | cbq | etc. }
// OPTIONS := ... try tc class add <desired QDISC_KIND> help

const ClassHelp = `Usage: tc class [ add | del | show ] dev STRING
	[ classid CLASSID ] [ root | parent CLASSID ]
	[ [ QDISC_KIND ] [ help | OPTIONS ] ]

	tc class show [ dev STRING ] [ root | parent CLASSID ]
Where:
	QDISC_KIND := { htb | hfcs }
	OPTIONS := ... try tc class add <desired QDISC_KIND> help
	`
