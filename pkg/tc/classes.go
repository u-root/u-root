// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trafficctl

import (
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/florianl/go-tc"
)

const maxUint32 = 0xFFFF_FFFF

// ParseHFSCClassArgs parses the cmdline arguments for `tc class add ... hfsc ...`
// and returns a *tc.Object.
func ParseHFSCClassArgs(out io.Writer, args []string) (*tc.Object, error) {
	if len(args) > 0 && args[0] == "help" {
		fmt.Fprint(out, HFSCHelp)
		return nil, ErrExitAfterHelp
	}
	ret := &tc.Object{}
	var fscOK, rscOK, uscOK bool
	hfsc := &tc.Hfsc{}

	if len(args) < 2 {
		return nil, ErrNotEnoughArgs
	}

	for i := 0; i < len(args); i = i + 2 {
		switch args[i] {
		case "rt":
			sc, err := HFSCGetSC(args[i+1:])
			if err != nil {
				return nil, err
			}
			hfsc.Rsc = sc
			rscOK = true
		case "ls":
			sc, err := HFSCGetSC(args[i+1:])
			if err != nil {
				return nil, err
			}
			hfsc.Fsc = sc
			fscOK = true
		case "sc":
			sc, err := HFSCGetSC(args[i+1:])
			if err != nil {
				return nil, err
			}
			hfsc.Fsc = sc
			hfsc.Rsc = sc
			rscOK = true
			fscOK = true
		case "ul":
			sc, err := HFSCGetSC(args[i+1:])
			if err != nil {
				return nil, err
			}
			hfsc.Usc = sc
			uscOK = true
		case "help":
		default:
			return nil, ErrInvalidArg
		}
	}

	if !(rscOK || fscOK || uscOK) {
		return nil, ErrInvalidArg
	}

	if uscOK && !fscOK {
		return nil, ErrInvalidArg
	}

	ret.Kind = "hfsc"
	ret.Hfsc = hfsc
	return ret, nil
}

func HFSCGetSC(args []string) (*tc.ServiceCurve, error) {
	var sc *tc.ServiceCurve
	sc1, err := hfscGetSC1(args)
	if err != nil {
		return nil, err
	}

	sc2, err := hfscGetSC2(args)
	if err != nil {
		return nil, err
	}

	if sc1 != nil && sc2 == nil {
		sc = sc1
	} else if sc1 == nil && sc2 != nil {
		sc = sc2
	} else {
		return nil, ErrInvalidArg
	}

	return sc, nil
}

func hfscGetSC1(args []string) (*tc.ServiceCurve, error) {
	if len(args) < 2 {
		return nil, ErrNotEnoughArgs
	}

	var d uint32
	var m1, m2 uint64
	var err error
	for i := 0; i < len(args); i = i + 2 {
		switch args[i] {
		case "m1":
			m1, err = ParseRate(args[i+1])
			if err != nil {
				return nil, err
			}
		case "d":
			d, err = parseTime(args[i+1])
			if err != nil {
				return nil, err
			}
		case "m2":
			m2, err = ParseRate(args[i+1])
			if err != nil {
				return nil, err
			}
		default:
			// Fallthrough if umax,dmax, rate
			if args[i] == "umax" || args[i] == "dmax" || args[i] == "rate" {
				return nil, nil
			}
			return nil, ErrInvalidArg
		}
	}

	cf, err := getClockfactor()
	if err != nil {
		return nil, err
	}
	sc := &tc.ServiceCurve{
		M1: uint32(m1),
		D:  cf * d,
		M2: uint32(m2),
	}

	return sc, nil
}

func hfscGetSC2(args []string) (*tc.ServiceCurve, error) {
	var umax, rate uint64
	var dmax uint32
	var err error

	for i := 0; i < len(args); i = i + 2 {
		switch args[i] {
		case "umax":
			umax, err = ParseSize(args[i+1])
			if err != nil {
				return nil, err
			}
		case "dmax":
			dmax, err = parseTime(args[i+1])
			if err != nil {
				return nil, err
			}
		case "rate":
			rate, err = ParseRate(args[i+1])
			if err != nil {
				return nil, err
			}
		default:
			// Just exit without error
			return nil, nil
		}
	}

	sc := &tc.ServiceCurve{}
	cf, err := getClockfactor()
	if err != nil {
		return nil, err
	}

	if dmax != 0 && math.Ceil(1.0*float64(umax)*TimeUnitsPerSecs/float64(dmax)) > float64(rate) {
		sc.M1 = uint32(math.Ceil(1.0 * float64(umax) * TimeUnitsPerSecs / float64(dmax)))
		sc.D = cf * dmax
		sc.M2 = uint32(rate)
	} else {
		sc.M1 = 0
		sc.D = uint32(float64(cf) * math.Ceil(float64(dmax)-float64(umax)*TimeUnitsPerSecs/float64(rate)))
		sc.M2 = uint32(rate)
	}

	return sc, nil
}

func supportetClasses(cl string) func(io.Writer, []string) (*tc.Object, error) {
	supported := map[string]func(io.Writer, []string) (*tc.Object, error){
		// Classful qdiscs
		"cbs":      nil, // (not supported for adding byt go-tc library)
		"htb":      ParseHTBClassArgs,
		"hfsc":     ParseHFSCClassArgs,
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

// ParseHTBClassArgs parses the cmdline arguments for `tc class add ... htb ...`
// and returns a *tc.Object.
func ParseHTBClassArgs(out io.Writer, args []string) (*tc.Object, error) {
	if len(args) > 0 && args[0] == "help" {
		fmt.Fprint(out, HTBHelp)
		return nil, ErrExitAfterHelp
	}
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
		default:
			return nil, ErrInvalidArg
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

	return ret, nil
}
