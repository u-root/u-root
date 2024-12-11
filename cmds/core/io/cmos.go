// Copyright 2012-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (!tinygo || tinygo.enable) && (amd64 || 386) && !freebsd && !windows

package main

import (
	"fmt"

	"github.com/u-root/u-root/pkg/cmos"
	"github.com/u-root/u-root/pkg/memio"
)

func init() {
	usageMsg += `io (cr index)... # read from CMOS register index [14-127]
io (cw index value)... # write value to CMOS register index [14-127]
io (rtcr index)... # read from RTC register index [0-13]
io (rtcw index value)... # write value to RTC register index [0-13]
`
	addCmd(readCmds, "cr", &cmd{cmosRead, 7, 8})
	addCmd(readCmds, "rtcr", &cmd{rtcRead, 7, 8})
	addCmd(writeCmds, "cw", &cmd{cmosWrite, 7, 8})
	addCmd(writeCmds, "rtcw", &cmd{rtcWrite, 7, 8})
}

func cmosRead(reg int64, data memio.UintN) error {
	c, err := cmos.New()
	if err != nil {
		return err
	}
	defer c.Close()
	regVal := memio.Uint8(reg)
	if regVal < 14 {
		return fmt.Errorf("byte %d is inside the range 0-13 which is reserved for RTC", regVal)
	}
	return c.Read(regVal, data)
}

func cmosWrite(reg int64, data memio.UintN) error {
	c, err := cmos.New()
	if err != nil {
		return err
	}
	defer c.Close()
	regVal := memio.Uint8(reg)
	if regVal < 14 {
		return fmt.Errorf("byte %d is inside the range 0-13 which is reserved for RTC", regVal)
	}
	return c.Write(regVal, data)
}

func rtcRead(reg int64, data memio.UintN) error {
	c, err := cmos.New()
	if err != nil {
		return err
	}
	defer c.Close()
	regVal := memio.Uint8(reg)
	if regVal > 13 {
		return fmt.Errorf("byte %d is outside the range 0-13 reserved for RTC", regVal)
	}
	return c.Read(regVal, data)
}

func rtcWrite(reg int64, data memio.UintN) error {
	c, err := cmos.New()
	if err != nil {
		return err
	}
	defer c.Close()
	regVal := memio.Uint8(reg)
	if regVal > 13 {
		return fmt.Errorf("byte %d is outside the range 0-13 reserved for RTC", regVal)
	}
	return c.Write(regVal, data)
}
