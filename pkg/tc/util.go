// Copyright 2012-20124 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trafficctl

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/florianl/go-tc"
)

const (
	TimeUnitsPerSecs = 1000000
)

func parseTime(t string) (uint32, error) {
	var cutstring string
	multiplier := TimeUnitsPerSecs
	if strings.HasSuffix(t, "sec") {
		cutstring, _ = strings.CutSuffix(t, "sec")
		multiplier = TimeUnitsPerSecs
	} else if strings.HasSuffix(t, "secs") {
		cutstring, _ = strings.CutSuffix(t, "secs")
		multiplier = TimeUnitsPerSecs
	} else if strings.HasSuffix(t, "s") {
		cutstring, _ = strings.CutSuffix(t, "s")
		multiplier = TimeUnitsPerSecs
	}

	if strings.HasSuffix(t, "ms") {
		cutstring, _ = strings.CutSuffix(t, "ms")
		multiplier = TimeUnitsPerSecs / 1000
	} else if strings.HasSuffix(t, "msec") {
		cutstring, _ = strings.CutSuffix(t, "msec")
		multiplier = TimeUnitsPerSecs / 1000
	} else if strings.HasSuffix(t, "msecs") {
		cutstring, _ = strings.CutSuffix(t, "msecs")
		multiplier = TimeUnitsPerSecs / 1000
	}

	if strings.HasSuffix(t, "us") {
		cutstring, _ = strings.CutSuffix(t, "us")
		multiplier = TimeUnitsPerSecs / 1000000
	} else if strings.HasSuffix(t, "usec") {
		cutstring, _ = strings.CutSuffix(t, "usec")
		multiplier = TimeUnitsPerSecs / 1000000
	} else if strings.HasSuffix(t, "usecs") {
		cutstring, _ = strings.CutSuffix(t, "usecs")
		multiplier = TimeUnitsPerSecs / 1000000
	}

	val, err := strconv.Atoi(cutstring)
	if err != nil {
		return 0, err
	}
	if val < 0x0 || val >= 0x7FFFFFFF {
		return 0, ErrOutOfBounds
	}

	ret := uint32(val) * uint32(multiplier)

	return ret, nil
}

func ParseHandle(h string) (uint32, error) {
	// split the string at :
	maj, _, ok := strings.Cut(h, ":")
	if !ok {
		return 0, ErrInvalidArg
	}

	major, err := strconv.ParseUint(maj, 16, 16)
	if err != nil {
		return 0, err
	}

	return uint32(major) << 16, nil
}

func ParseClassID(p string) (uint32, error) {
	if p == "root" {
		return tc.HandleRoot, nil
	}

	if p == "none" {
		return 0, nil
	}

	// split the string at :
	maj, min, ok := strings.Cut(p, ":")
	if !ok {
		return 0, ErrInvalidArg
	}

	major, err := strconv.ParseUint(maj, 16, 16)
	if err != nil {
		return 0, err
	}

	minor, err := strconv.ParseUint(min, 16, 16)
	if err != nil {
		return 0, err
	}

	return uint32(major<<16) | uint32(minor), nil
}

var (
	ErrUnknownLinkLayer = errors.New("unknown linklayer value provided")
)

func ParseLinkLayer(l string) (uint8, error) {
	for _, ll := range []struct {
		name string
		val  uint8
	}{
		{name: "ethernet", val: 1},
		{name: "atm", val: 2},
		{name: "ads1", val: 2},
	} {
		if ll.name == l {
			return ll.val, nil
		}
	}
	return 0xFF, ErrUnknownLinkLayer
}

func ParseSize(s string) (uint64, error) {
	sizeStr := strings.TrimRight(s, "gkmbit")

	sz, err := strconv.ParseUint(sizeStr, 10, 32)
	if err != nil {
		return 0, err
	}

	unitMuliplier := strings.TrimLeft(s, "0123456789")

	switch unitMuliplier {
	case "k", "kb":
		sz *= 1024
	case "m", "mb":
		sz *= 1024 * 1024
	case "g", "gb":
		sz *= 1024 * 1024 * 1024
	case "kbit":
		sz *= 1024 / 8
	case "mbit":
		sz *= 1024 * 1024 / 8
	case "gbit":
		sz *= 1024 * 1024 * 1024 / 8
	}

	return sz, nil
}

func ParseRate(arg string) (uint64, error) {
	unit := strings.TrimLeft(arg, "0123456789")

	sizeStr := strings.TrimRight(arg, "bBgGKkMmTitps")
	sz, err := strconv.ParseUint(sizeStr, 10, 32)
	if err != nil {
		return 0, err
	}

	for _, entry := range []struct {
		unit  string
		value uint64
	}{
		{unit: "bit", value: 1},
		{unit: "Kibit", value: 1024},
		{unit: "mibit", value: 1024 * 1024},
		{unit: "gibit", value: 1024 * 1024 * 1024},
		{unit: "tibit", value: 1024 * 1024 * 1024 * 1024},
		{unit: "kbit", value: 1000},
		{unit: "mbit", value: 1000 * 1000},
		{unit: "gbit", value: 1000 * 1000 * 1000},
		{unit: "tit", value: 1000 * 1000 * 1000 * 1000},
		{unit: "Bps", value: 8},
		{unit: "KiBps", value: 8 * 1024},
		{unit: "Mibit", value: 8 * 1024 * 1024},
		{unit: "Gibit", value: 8 * 1024 * 1024 * 1024},
		{unit: "TiBps", value: 8 * 1024 * 1024 * 1024 * 1024},
		{unit: "KBps", value: 8 * 1000},
		{unit: "MBps", value: 8 * 1000 * 1000},
		{unit: "GBps", value: 8 * 1000 * 1000 * 1000},
		{unit: "TBps", value: 8 * 1000 * 1000 * 1000 * 1000},
	} {
		if entry.unit == unit {
			return (sz * entry.value) / 8, nil
		}
	}
	return 0, ErrInvalidArg
}

func GetHz() (int, error) {
	const HZdef = 100
	psched, err := os.Open("/proc/net/psched")
	if err != nil {
		return 0, err
	}
	defer psched.Close()

	var gb1, gb2, nom, denom int

	fmt.Fscanf(psched, "%8x %8x %8x %8x",
		&gb1,
		&gb2,
		&nom,
		&denom)

	if nom == 1000000 {
		return denom, nil
	}

	return HZdef, nil
}

func CalcXMitTime(rate uint64, size uint32) (uint32, error) {
	const TimeUnitsPerSecs = float64(1000000)
	ret := TimeUnitsPerSecs * (float64(size) / float64(rate))
	if ret >= 0xFFFF_FFFF {
		ret = maxUint32
	}

	tickInUsec, err := getTickInUsec()
	if err != nil {
		return 0, err
	}

	return uint32(ret) * tickInUsec, nil
}

func getTickInUsec() (uint32, error) {
	psched, err := os.Open("/proc/net/psched")
	if err != nil {
		return 0, err
	}
	defer psched.Close()

	var t2us, us2t, clockRes, gb int

	fmt.Fscanf(psched, "%8x %8x %8x %8x",
		&t2us,
		&us2t,
		&clockRes,
		&gb)

	if clockRes == 1000000000 {
		t2us = us2t
	}

	clockFactor := int64(clockRes / TimeUnitsPerSecs)

	return uint32(float64(t2us)/float64(us2t)) * uint32(clockFactor), nil
}
