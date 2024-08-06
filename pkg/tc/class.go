// Copyright 2012-20124 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trafficctl

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/florianl/go-tc"
	"golang.org/x/sys/unix"
)

// possible args:
// "dev DEV parent qdisc-id [ classid class-id ] qdisc [ qdisc specific parameters ]"
type CArgs struct {
	dev    *string
	parent *uint32
	handle *uint32
}

func ParseClassArgs(args []string, stdout io.Writer) (*CArgs, error) {
	ret := &CArgs{}
	if len(args) < 1 {
		fmt.Println(args)
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
			qdiscID, err := strconv.Atoi(val)
			if err != nil {
				return ret, err
			}
			if qdiscID < 0 || qdiscID >= 0x7FFFFFFF {
				return nil, ErrOutOfBounds
			}
			indirect := uint32(qdiscID)
			ret.parent = &indirect
		case "root":
			indirect := tc.HandleRoot
			ret.parent = &indirect
			// We have a one piece argument. To get to the next arg properly
			i--
		case "classid":
			// Dealing with major:minor of ID
			if val == "root" {
				indirect := tc.HandleRoot
				ret.handle = &indirect
				continue
			}
			if val == "none" {
				indirect := uint32(0x0)
				ret.handle = &indirect
				continue
			}
			// Split the string
			cid := strings.Split(val, ":")
			//MinorID
			minClassID, err := strconv.Atoi(cid[1])
			if err != nil {
				return ret, err
			}

			if minClassID >= (1 << 16) {
				return ret, ErrOutOfBounds
			}

			majClassID, err := strconv.Atoi(cid[0])
			if err != nil {
				return ret, nil
			}

			if majClassID >= (1 << 16) {
				return ret, ErrOutOfBounds
			}

			majClassID <<= 16

			indirect := uint32(majClassID) + uint32(minClassID)
			ret.handle = &indirect
		case "estimator":
			return nil, ErrNotImplemented
		case "help":
			PrintClassHelp(stdout)
			return nil, nil
		default:
			// Resolve Qdisc and parameters
			if qobj := supportetQdisc(args[i]); qobj == nil {
				return nil, fmt.Errorf("%w: invalid class: %s", ErrInvalidArg, args[i])
			}
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
	return nil
}

func (t *Trafficctl) DeleteClass(args *CArgs, stdout io.Writer) error {
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
