// Copyright 2019-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// hwclock sets or reads the hwclock in UTC format.
//
// Synopsis:
//     hwclock [-w]
//
// Description:
//     It prints the current hwclock time in UTC if called without any flags.
//     It sets the hwclock to the system clock in UTC if called with -w.
//
// Options:
//     -w: set hwclock to system clock in UTC
//
// Author:
//     claudiozumbo@gmail.com.

package main

import (
	"flag"
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"syscall"
	"time"
	"unsafe"
)

const startYear = 1900

type tm struct {
	tm_sec   int32
	tm_min   int32
	tm_hour  int32
	tm_mday  int32
	tm_mon   int32
	tm_year  int32
	tm_wday  int32
	tm_yday  int32
	tm_isdst int32
}

func readRtc() (time.Time, error) {

	var (
		ptm     tm
		retTime time.Time
	)

	rtcFile, err := os.Open("/dev/rtc")
	if err != nil {
		return retTime, err
	}

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(rtcFile.Fd()),
		uintptr(unix.RTC_RD_TIME),
		uintptr(unsafe.Pointer(&ptm)),
	)
	if errno != 0 {
		return retTime, fmt.Errorf("ioctl RTC_RD_TIME returned with errno: %v\n", errno)
	}

	retTime = time.Date(int(ptm.tm_year)+startYear,
		time.Month(ptm.tm_mon),
		int(ptm.tm_mday),
		int(ptm.tm_hour),
		int(ptm.tm_min),
		int(ptm.tm_sec),
		0,
		time.UTC)

	return retTime, nil
}

func setRtcFromSysClock() error {
	timeUTC := time.Now().UTC()

	rtc, err := os.Open("/dev/rtc")
	if err != nil {
		return err
	}
	stm := tm{tm_sec: int32(timeUTC.Second()),
		tm_min:   int32(timeUTC.Minute()),
		tm_hour:  int32(timeUTC.Hour()),
		tm_mday:  int32(timeUTC.Day()),
		tm_mon:   int32(timeUTC.Month()),
		tm_year:  int32(timeUTC.Year() - startYear),
		tm_wday:  int32(0),
		tm_yday:  int32(0),
		tm_isdst: int32(0)}

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		rtc.Fd(),
		uintptr(unix.RTC_SET_TIME),
		uintptr(unsafe.Pointer(&stm)),
	)
	if errno != 0 {
		return fmt.Errorf("ioctl RTC_SET_TIME returned with errno: %v", errno)
	}
	return nil
}

func main() {
	writePtr := flag.Bool("w", false, "Set hwclock from system clock in UTC")
	flag.Parse()

	if *writePtr {
		if err := setRtcFromSysClock(); err != nil {
			panic(err)
		}
	} else {
		t, err := readRtc()
		if err != nil {
			panic(err)
		}
		fmt.Println(t)
	}
}
