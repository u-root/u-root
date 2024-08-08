// Copyright 2012-20124 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trafficctl

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/florianl/go-tc"
	"github.com/florianl/go-tc/core"
)

var (
	ErrInvalidFilterType = errors.New("invalid filtertype")
)

// FArgs hold all possible args for qdisc subcommand
// tc filter [ add | del | change | replace | show ] [ dev STRING ]
// tc filter [ add | del | change | replace | show ] [ block BLOCK_INDEX ]
// tc filter get dev STRING parent CLASSID protocol PROTO handle FILTERID pref PRIO FILTER_TYPE
// tc filter get block BLOCK_INDEX protocol PROTO handle FILTERID pref PRIO FILTER_TYPE
// [ pref PRIO ] protocol PROTO [ chain CHAIN_INDEX ]
// [ estimator INTERVAL TIME_CONSTANT ]
// [ root | ingress | egress | parent CLASSID ]
// [ handle FILTERID ] [ [ FILTER_TYPE ] [ help | OPTIONS ] ]
// tc filter show [ dev STRING ] [ root | ingress | egress | parent CLASSID ]
// tc filter show [ block BLOCK_INDEX ]
type FArgs struct {
	dev       string
	kind      *string
	parent    *uint32
	handle    *uint32
	protocol  *string
	pref      *uint32
	filterObj *tc.Object
}

func ParseFilterArgs(args []string, stdout io.Writer) (*FArgs, error) {
	ret := &FArgs{}
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
		case "parent":
			parent, err := strconv.Atoi(val)
			if err != nil {
				return ret, err
			}
			if parent < 0 || parent >= 0x7FFFFFFF {
				return nil, ErrOutOfBounds
			}
			indirect := uint32(parent)
			ret.parent = &indirect
		case "protocol", "proto":
			proto := args[i+1]
			ret.protocol = &proto
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
		case "preference", "pref":
			val, err := strconv.Atoi(val)
			if err != nil {
				return nil, err
			}
			if val < 0 || val >= 0x7FFFFFFF {
				return nil, ErrOutOfBounds
			}
			indirect := uint32(val)
			ret.pref = &indirect
		case "block":
			return nil, ErrNotImplemented
		case "chain":
			return nil, ErrNotImplemented
		case "estimator":
			return nil, ErrNotImplemented
		case "root":
			if ret.parent != nil {
				return nil, ErrInvalidArg
			}
			indirect := tc.HandleRoot
			ret.parent = &indirect
			// We have a one piece argument. To get to the next arg properly
			i--
		case "ingress":
			if ret.parent != nil {
				return nil, ErrInvalidArg
			}
			indirectPar := tc.HandleIngress // is the same as clsact handle
			ret.parent = &indirectPar
			// We have a one piece argument. To get to the next arg properly
			indirectHan := core.BuildHandle(tc.HandleIngress, 0)
			ret.handle = &indirectHan

			i--
		case "egress":
			if ret.parent != nil {
				return nil, ErrInvalidArg
			}
			indirectPar := tc.HandleIngress // is the same as clsact handle
			ret.parent = &indirectPar
			// We have a one piece argument. To get to the next arg properly
			indirectHan := core.BuildHandle(tc.HandleIngress, tc.HandleMinEgress)
			ret.handle = &indirectHan

			i--
		case "help":
			PrintFilterHelp(stdout)
		default: // I hope we parsed all the stuff until here
			// args[i] is the actual filter type
			// Resolve Qdisc and parameters
			var filterParse func([]string) (*tc.Object, error)
			var err error
			if filterParse = supportedFilters(args[i]); filterParse == nil {
				return nil, fmt.Errorf("%w: invalid filter: %s", ErrInvalidArg, args[i])
			}
			k := args[i]
			ret.kind = &k

			ret.filterObj, err = filterParse(args[i+1:])
			if err != nil {
				return nil, err
			}
			return ret, nil
		}
	}
	return ret, nil
}

func (t *Trafficctl) ShowFilter(fargs *FArgs, stdout io.Writer) error {

	return nil
}

func (t *Trafficctl) AddFilter(fargs *FArgs, stdout io.Writer) error {
	if err := t.Tc.Filter().Add(fargs.filterObj); err != nil {
		return err
	}
	return nil
}

func (t *Trafficctl) DeleteFilter(fargs *FArgs, stdout io.Writer) error {
	return nil
}

func (t *Trafficctl) ReplaceFilter(fargs *FArgs, stdout io.Writer) error {
	return nil
}

func (t *Trafficctl) ChangeFilter(fargs *FArgs, stdout io.Writer) error {
	return nil
}

func (t *Trafficctl) GetFilter(fargs *FArgs, stdout io.Writer) error {
	return nil
}

const (
	Filterhelp = `Usage:
	tc filter [ add | del | change | replace | show ] [ dev STRING ]
	tc filter [ add | del | change | replace | show ] [ block BLOCK_INDEX ]
	tc filter get dev STRING parent CLASSID protocol PROTO handle FILTERID pref PRIO FILTER_TYPE
	tc filter get block BLOCK_INDEX protocol PROTO handle FILTERID pref PRIO FILTER_TYPE
		[ pref PRIO ] protocol PROTO [ chain CHAIN_INDEX ]
		[ estimator INTERVAL TIME_CONSTANT ]
		[ root | ingress | egress | parent CLASSID ]
		[ handle FILTERID ] [ [ FILTER_TYPE ] [ help | OPTIONS ] ]
	tc filter show [ dev STRING ] [ root | ingress | egress | parent CLASSID ]
	tc filter show [ block BLOCK_INDEX ]

	Where:
	FILTER_TYPE := { u32 | bpf | fw | route | etc. }
	FILTERID := ... format depends on classifier, see there
	OPTIONS := ... try tc filter add <desired FILTER_KIND> help
`
)

func PrintFilterHelp(stdout io.Writer) {
	fmt.Fprint(stdout,
		Filterhelp)
}

func supportedFilters(f string) func([]string) (*tc.Object, error) {
	supported := map[string]func([]string) (*tc.Object, error){
		"basic":    parseBasicParams,
		"bpf":      nil,
		"cgroup":   nil,
		"flow":     nil,
		"flower":   nil,
		"fw":       nil,
		"route":    nil,
		"u32":      nil,
		"matchall": nil,
	}

	ret, ok := supported[f]
	if !ok {
		return nil
	}

	return ret
}
