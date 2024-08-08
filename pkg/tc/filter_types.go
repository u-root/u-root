// Copyright 2012-20124 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trafficctl

import (
	"fmt"
	"os"
	"strconv"

	"github.com/florianl/go-tc"
)

func parseBasicParams(params []string) (*tc.Object, error) {
	b := &tc.Basic{}
	var err error

	for i := 0; i < len(params); i = i + 2 {
		switch params[i] {
		case "match":
			return nil, ErrNotImplemented
		case "action":
			// Only generic actions allowed here
			b.Actions, err = parseActionGAT(params[1:])
			if err != nil {
				return nil, err
			}
		case "classid", "flowid":
			id, err := strconv.Atoi(params[1])
			if err != nil {
				return nil, err
			}
			if id < 0x0 || id >= 0x7FFFFFFF {
				return nil, ErrOutOfBounds
			}
			indirect := uint32(id)
			b.ClassID = &indirect
		case "help":
			fmt.Printf("%s\n", "tc filter ... basic [ match EMATCH_TREE ] [ action ACTION_SPEC ] [ classid CLASSID ]")
			os.Exit(0)
		default:
			//not sure yet
		}
	}

	ret := &tc.Object{}
	ret.Basic = b

	return ret, nil
}
