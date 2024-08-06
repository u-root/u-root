// Copyright 2012-20124 the u-root Authors. All rights reserved
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

// possible args:
// "dev DEV parent qdisc-id [ classid class-id ] qdisc [ qdisc specific parameters ]"
type CArgs struct {
	dev    *string
	parent *uint32
	handle *uint32
	obj    *tc.Object
}

func ParseClassArgs(stdout io.Writer, args []string) (*CArgs, error) {
	ret := &CArgs{}
	if len(args) < 1 {
		return nil, fmt.Errorf("ParseClassArgs() = %v", ErrNotEnoughArgs)
	}

	for i := 0; i <= len(args[1:]); i = i + 2 {
		var val string
		if len(args[1:]) > i {
			val = args[i+1]
		}
		fmt.Printf("args[%d]: %s\n", i, args[i])
		switch args[i] {
		case "dev":
			ret.dev = &val
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
		case "estimator":
			return nil, ErrNotImplemented
		case "help":
			PrintClassHelp(stdout)
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

// possible args:
// "tc class [ add | del | change | replace | show ] dev STRING [ classid CLASSID ] [ root | parent CLASSID ] [ [ QDISC_KIND ] [ help | OPTIONS ] ]"
func (t *Trafficctl) ShowClass(cArgs *CArgs, stdout io.Writer) error {
	ifs, err := net.Interfaces()
	if err != nil {
		return err
	}

	for _, iface := range ifs {
		if cArgs.dev != nil && iface.Name != *cArgs.dev {
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
				fmt.Fprintf(stdout, "%20s\t%s\n", iface.Name, class.Kind)
			}
		}
	}

	return nil
}

func (t *Trafficctl) AddClass(args *CArgs, stdout io.Writer) error {
	iface, err := net.InterfaceByName(*args.dev)
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

func (t *Trafficctl) DeleteClass(args *CArgs, stdout io.Writer) error {
	iface, err := net.InterfaceByName(*args.dev)
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

func (t *Trafficctl) ChangeClass(args *CArgs, stdout io.Writer) error {
	return nil
}

func (t *Trafficctl) ReplaceClass(args *CArgs, stdout io.Writer) error {
	return nil
}

const (
	ClassHelp = `Usage:
tc class [ add | del | change | replace | show ] dev STRING
	[ classid CLASSID ] [ root | parent CLASSID ]
	[ [ QDISC_KIND ] [ help | OPTIONS ] ]

	tc class show [ dev STRING ] [ root | parent CLASSID ]
	"Where:
	QDISC_KIND := { prio | etc. }"
	OPTIONS := ... try tc class add <desired QDISC_KIND> help`
)

func PrintClassHelp(stdout io.Writer) {
	fmt.Fprintf(stdout, "%s", ClassHelp)
}
