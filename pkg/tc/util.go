// Copyright 2012-20124 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trafficctl

import (
	"strconv"
	"strings"
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
