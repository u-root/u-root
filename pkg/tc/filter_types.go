// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trafficctl

import (
	"fmt"
	"io"
	"strconv"

	"github.com/florianl/go-tc"
)

// Originally from tc filter basic help
// Usage: ... basic [ match EMATCH_TREE ]
//                  [ action ACTION_SPEC ] [ classid CLASSID ]

// Where:	SELECTOR := SAMPLE SAMPLE ...
// 	FILTERID := X:Y:Z
// 	ACTION_SPEC := ... look at individual actions

// NOTE: CLASSID is parsed as hexadecimal input.

const BasicHelp = `Usage: ... basic [ action ACTION_SPEC ] [ classid CLASSID | flowid CLASSID ]
`

// ParseBasicParams parses the cmdline arguments for `tc filter ... basic ...`
// and returns a *tc.Object.
func ParseBasicParams(out io.Writer, params []string) (*tc.Object, error) {
	b := &tc.Basic{}
	var err error

	for i := 0; i < len(params); i = i + 2 {
		switch params[i] {
		case "match":
			return nil, ErrNotImplemented
		case "action":
			// Only generic actions allowed here
			b.Actions, err = ParseActionGAT(out, params[1:])
			if err != nil {
				return nil, err
			}
		case "classid", "flowid":
			id, err := ParseClassID(params[1])
			if err != nil {
				return nil, err
			}
			indirect := uint32(id)
			b.ClassID = &indirect
		case "help":
			fmt.Fprint(out, BasicHelp)
			return nil, ErrExitAfterHelp
		default:
			return nil, ErrInvalidArg
		}
	}

	ret := &tc.Object{}
	ret.Kind = "basic"
	ret.Basic = b

	return ret, nil
}

// Originally from tc filter u32 help
// Usage: ... u32 [ match SELECTOR ... ] [ link HTID ] [ classid CLASSID ]
//                [ action ACTION_SPEC ] [ offset OFFSET_SPEC ]
//                [ ht HTID ] [ hashkey HASHKEY_SPEC ]
//                [ sample SAMPLE ] [skip_hw | skip_sw]
// or         u32 divisor DIVISOR

// Where: SELECTOR := SAMPLE SAMPLE ...
//        SAMPLE := { ip | ip6 | udp | tcp | icmp | u{32|16|8} | mark }
//                  SAMPLE_ARGS [ divisor DIVISOR ]
//        FILTERID := X:Y:Z

// NOTE: CLASSID is parsed at hexadecimal input.

const U32Help = `Usage ... u32 [ match ...] [ classid CLASSID | flowid CLASSID ]
For further information see https://linux-tc-notes.sourceforge.net/tc/doc/cls_u32.txt
`

const TCU32Terminal = 1

// ParseU32Params parses the cmdline arguments for `tc filter ... u32 ...` and
// returns a *tc.Object. ParseU32Params recognizes a limited sub-language of
// the language that "tc" of iproute2 recognizes. Reference:
// <https://linux-tc-notes.sourceforge.net/tc/doc/cls_u32.txt>.
func ParseU32Params(out io.Writer, params []string) (*tc.Object, error) {
	u32 := &tc.U32{
		Sel: &tc.U32Sel{
			Flags: TCU32Terminal,
		},
	}
	i := 0
	for i < len(params) {
		switch params[i] {
		case "classid", "flowid":
			if len(params)-i == 1 {
				return nil, ErrInvalidArg
			}
			id, err := ParseClassID(params[i+1])
			if err != nil {
				return nil, err
			}
			u32.ClassID = &id
			i = i + 2
		case "match":
			if len(params)-i <= 5 || params[i+1] != "u32" || params[i+4] != "at" {
				return nil, ErrInvalidArg
			}

			val64, err := strconv.ParseUint(params[i+2], 0, 32)
			if err != nil {
				return nil, err
			}

			mask64, err := strconv.ParseUint(params[i+3], 0, 32)
			if err != nil {
				return nil, err
			}

			off64, err := strconv.ParseUint(params[i+5], 0, 32)
			if err != nil {
				return nil, err
			}

			if u32.Sel.NKeys == 255 {
				return nil, ErrOutOfBounds
			}

			u32.Sel.NKeys++
			key := tc.U32Key{
				Mask: HToNL(uint32(mask64)),
				Val:  HToNL(uint32(val64)),
				Off:  uint32(off64),
			}
			u32.Sel.Keys = append(u32.Sel.Keys, key)

			i = i + 6
		case "help":
			fmt.Fprint(out, U32Help)
			return nil, ErrExitAfterHelp
		default:
			return nil, ErrInvalidArg
		}
	}

	if u32.ClassID == nil || u32.Sel.NKeys == 0 {
		return nil, ErrInvalidArg
	}

	ret := &tc.Object{}
	ret.Kind = "u32"
	ret.U32 = u32

	return ret, nil
}
