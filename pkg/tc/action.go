// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trafficctl

import (
	"errors"
	"fmt"
	"io"

	"github.com/florianl/go-tc"
)

var ErrInvalidActionControl = errors.New("invalid action control parameter")

const (
	GActUnspec  = -1
	GActOk      = 0
	GActReclass = 1
	GActShot    = 2
	GActPipe    = 3
	GActTrap    = 8
	GActJump    = 1 << 28
	GActGoTo    = 2 << 28
)

// ParseActionGAT parses options of the filter action category and returns
// a pointer to a slice of []*tc.Action
func ParseActionGAT(out io.Writer, args []string) (*[]*tc.Action, error) {
	if len(args) < 1 {
		return nil, ErrNotEnoughArgs
	}

	var act int
	switch args[0] {
	case "continue":
		act = GActUnspec
	case "drop":
		act = GActShot
	case "shot":
		act = GActShot
	case "pass":
		act = GActOk
	case "ok":
		act = GActOk
	case "reclassify":
		act = GActReclass
	case "pipe":
		act = GActPipe
	case "goto":
		act = GActGoTo
	case "jump":
		act = GActJump
	case "trap":
		act = GActTrap
	case "help":
		fmt.Fprintf(out, "%s\n", gactHelp)
		return nil, nil
	default:
		fmt.Fprintf(out, "%s\n", gactHelp)
		return nil, ErrInvalidActionControl
	}

	gact := &tc.Gact{
		Parms: &tc.GactParms{
			Action: uint32(act),
		},
	}

	a := &tc.Action{
		Kind: "gact",
		Gact: gact,
	}

	ret := make([]*tc.Action, 0)
	ret = append(ret, a)

	return &ret, nil
}

const (
	gactHelp = `Usage:
tc ... action CONTROL [ RAND ] [ INDEX ]
	CONTROL := { reclassify | drop | continue | pass | pipe |
		goto chain CHAIN_INDEX |
		jump JUMP_COUNT }

	RAND := random RANDTYPE CONTROL VAL
	RANDTYPE := { netrand | determ }
	VAL := number not exceeding 10000
	JUMP_COUNT := absolute jump from start of action list
	INDEX := index value used`
)
