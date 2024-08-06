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

const (
	QFQHelp = `Usage: ... qfq weight NUMBER maxpkt BYTES`

	maxUint32 = 0xFFFF_FFFF
)

func ParseHFSCClassArgs(out io.Writer, args []string) (*tc.Object, error) {
	ret := &tc.Object{}
	ret.Kind = "hfsc"
	return ret, nil
}

func HFSCGetSC(args []string) (*tc.ServiceCurve, error) {
	m1, d, m2, err := hfscGetSC1(args)
	if err != nil {
		return nil, err
	}

	m1, d, m2, err = hfscGetSC2(args)
	if err != nil {
		return nil, err
	}

	return &tc.ServiceCurve{
		M1: uint32(m1),
		D:  d,
		M2: uint32(m2),
	}, nil
}

func hfscGetSC1(args []string) (uint64, uint32, uint64, error) {
	if len(args) < 2 {
		return 0, 0, 0, ErrNotEnoughArgs
	}

	var d uint32
	var m1, m2 uint64
	var err error
	for i := 0; i < len(args); i = i + 2 {
		switch args[i] {
		case "m1":
			m1, err = ParseRate(args[i+1])
			if err != nil {
				return 0, 0, 0, err
			}
		case "d":
			d, err = parseTime(args[i+1])
			if err != nil {
				return 0, 0, 0, err
			}
		case "m2":
			m2, err = ParseRate(args[i+1])
			if err != nil {
				return 0, 0, 0, err
			}
		default:
			// Fallthrough if umax,dmax, rate
			if args[i] == "umax" || args[i] == "dmax" || args[i] == "rate" {
				return m1, d, m2, nil
			}
			return 0, 0, 0, ErrInvalidArg
		}
	}
	return m1, d, m2, nil
}

func hfscGetSC2(args []string) (uint64, uint32, uint64, error) {
	var m1, m2 uint64
	var umax, rate uint64
	var dmax uint32

	var d uint32
	var err error

	for i := 0; i < len(args); i = i + 2 {
		switch args[i] {
		case "umax":
			umax, err = ParseSize(args[i+1])
			if err != nil {
				return 0, 0, 0, err
			}
		case "dmax":
			dmax, err = parseTime(args[i+1])
			if err != nil {
				return 0, 0, 0, err
			}
		case "rate":
			rate, err = ParseRate(args[i+1])
			if err != nil {
				return 0, 0, 0, err
			}
		default:
			//What is this?
		}
	}
	return m1, d, m2, nil
}

func supportetClasses(cl string) func(io.Writer, []string) (*tc.Object, error) {
	supported := map[string]func(io.Writer, []string) (*tc.Object, error){
		// Classful qdiscs
		"cbs":      nil, // (not supported for adding byt go-tc library)
		"htb":      ParseHTBClassArgs,
		"hfsc":     nil,
		"hfscqopt": nil, // (not supported for adding byt go-tc library)
		"dsmark":   nil, // (not supported for adding byt go-tc library)
		"drr":      nil, // (not supported for adding byt go-tc library)
		"cbq":      nil,
		"atm":      nil, // (not supported for adding byt go-tc library)
		"taprio":   nil, // (not supported for adding byt go-tc library)
	}

	ret, ok := supported[cl]
	if !ok {
		return nil
	}

	return ret
}

func ParseHTBClassArgs(out io.Writer, args []string) (*tc.Object, error) {
	const linkLayerMask = 0x0F
	// rate <rate> and burst <bytes> is required
	if len(args) < 4 {
		return nil, ErrInvalidArg
	}

	buffer, cbuffer := uint32(0), uint32(0)
	mtu := uint32(1600)
	mpu := uint16(0)
	overhead := uint16(0)
	linkLayer := uint8(0)
	ceil64 := uint64(0)
	rate64 := uint64(0)

	var ceilBool bool

	opt := tc.HtbOpt{
		Rate: tc.RateSpec{},
		Ceil: tc.RateSpec{},
	}

	for i := 0; i < len(args); i = i + 2 {
		switch args[i] {
		case "prio":
			prio, err := strconv.ParseUint(args[i+1], 10, 32)
			if err != nil {
				return nil, err
			}
			opt.Prio = uint32(prio)
		case "mtu":
			m, err := strconv.ParseUint(args[i+1], 10, 32)
			if err != nil {
				return nil, err
			}
			mtu = uint32(m)
		case "mpu":
			m, err := strconv.ParseUint(args[i+1], 10, 16)
			if err != nil {
				return nil, err
			}
			mpu = uint16(m)
		case "overhead":
			o, err := strconv.ParseUint(args[i+1], 10, 16)
			if err != nil {
				return nil, err
			}
			overhead = uint16(o)
		case "linklayer":
			ll, err := ParseLinkLayer(args[i+1])
			if err != nil {
				return nil, err
			}
			linkLayer = ll
		case "quantum":
			q, err := strconv.ParseUint(args[i+1], 10, 32)
			if err != nil {
				return nil, err
			}
			opt.Quantum = uint32(q)
		case "burst", "buffer", "maxburst":
			b, err := ParseSize(args[i+1])
			if err != nil {
				return nil, err
			}
			buffer = uint32(b)
		case "cburst", "cbuffer", "cmaxburst":
			b, err := ParseSize(args[i+1])
			if err != nil {
				return nil, err
			}
			cbuffer = uint32(b)
		case "ceil":
			c, err := ParseRate(args[i+1])
			if err != nil {
				return nil, err
			}
			ceil64 = c
			ceilBool = true
		case "rate":
			r, err := ParseRate(args[i+1])
			if err != nil {
				return nil, err
			}
			rate64 = r
		case "help":
			fmt.Fprintf(out, "%s", HTBHelp)
		}
	}

	if rate64 >= maxUint32 {
		opt.Rate.Rate = maxUint32
	} else {
		opt.Rate.Rate = uint32(rate64)
	}

	if !ceilBool {
		ceil64 = rate64
	}

	if ceil64 >= maxUint32 {
		opt.Ceil.Rate = maxUint32
	} else {
		opt.Ceil.Rate = uint32(ceil64)
	}

	hz, err := GetHz()
	if err != nil {
		return nil, err
	}

	if buffer == 0 {
		buffer = uint32(uint64(rate64)/uint64(hz) + uint64(mtu))
	}

	if cbuffer == 0 {
		cbuffer = uint32(uint64(ceil64)/uint64(hz) + uint64(mtu))
	}

	opt.Rate.Overhead = overhead
	opt.Rate.Linklayer = linkLayer & linkLayerMask
	opt.Rate.Mpu = mpu
	for (mtu >> uint32(opt.Rate.CellLog)) > 255 {
		opt.Rate.CellLog++
	}

	opt.Ceil.Overhead = overhead
	opt.Ceil.Linklayer = linkLayer & linkLayerMask
	opt.Ceil.Mpu = mpu
	for (mtu >> uint32(opt.Ceil.CellLog)) > 255 {
		opt.Ceil.CellLog++
	}

	if opt.Buffer == 0 {
		opt.Buffer, err = CalcXMitTime(rate64, buffer)
		if err != nil {
			return nil, err
		}
	}

	if opt.Cbuffer == 0 {
		opt.Cbuffer, err = CalcXMitTime(ceil64, cbuffer)
		if err != nil {
			return nil, err
		}
	}

	ret := &tc.Object{}

	ret.Kind = "htb"
	ret.Htb = &tc.Htb{
		Parms: &opt,
	}

	ret.Htb.Rate64 = &rate64
	ret.Htb.Ceil64 = &ceil64

	fmt.Printf("%v\n", opt.Buffer)

	return ret, nil
}
