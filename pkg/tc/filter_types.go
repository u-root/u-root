// Copyright 2012-20124 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trafficctl

import (
	"fmt"
	"io"
	"strconv"

	"github.com/florianl/go-tc"
)

func parseBasicParams(out io.Writer, params []string) (*tc.Object, error) {
	b := &tc.Basic{}
	var err error

	for i := 0; i < len(params); i = i + 2 {
		switch params[i] {
		case "match":
			return nil, ErrNotImplemented
		case "action":
			// Only generic actions allowed here
			b.Actions, err = ParseActionGAT(params[1:], out)
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
			fmt.Fprintf(out, "%s\n", "tc filter ... basic [ match EMATCH_TREE ] [ action ACTION_SPEC ] [ classid CLASSID ]")
			return nil, nil
		default:
			//not sure yet
		}
	}

	ret := &tc.Object{}
	ret.Kind = "basic"
	ret.Basic = b

	return ret, nil
}
