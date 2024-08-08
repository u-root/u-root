// Copyright 2012-20124 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trafficctl

import (
	"errors"
	"fmt"
	"os"

	"github.com/florianl/go-tc"
)

var (
	ErrInvalidActionControl = errors.New("invalid action control parameter")
)

const (
	GActUnspec  = -1
	GActOk      = 0
	GActReclass = 1
	GActShot    = 2
	GActPipe    = 3
	GActTrap    = 8
	GActJump    = 1 << 28
	GactGoTo    = 2 << 28
)

func parseActionGAT(args []string) (*[]*tc.Action, error) {
	if len(args) < 1 {
		return nil, ErrNotEnoughArgs
	}

	GActMap := map[string]int{
		"continue":   GActUnspec,
		"drop":       GActShot,
		"shot":       GActShot,
		"pass":       GActOk,
		"ok":         GActOk,
		"reclassify": GActReclass,
		"pipe":       GActPipe,
		"goto":       GactGoTo,
		"jump":       GActJump,
		"trap":       GActTrap,
		"help":       0xFFFF,
	}

	act, ok := GActMap[args[0]]
	if !ok {
		return nil, ErrInvalidActionControl
	}

	if act == 0xFFFF {
		printGActHelp()
		os.Exit(0)
	}

	gact := &tc.Gact{
		Tm: &tc.Tcft{},
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

func printGActHelp() {
	fmt.Printf("%s\n", gactHelp)
}
