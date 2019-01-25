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
	"log"
	"os"
	"syscall"
	"time"
	"unsafe"
)

const startYear = 1900

var write = flag.Bool("w", false, "Set hwclock from system clock in UTC")

func readRtc(rtcFile *os.File) (time.Time, error) {

	var (
		ptm     unix.RTCTime
		ret time.Time
	)

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(rtcFile.Fd()),
		uintptr(unix.RTC_RD_TIME),
		uintptr(unsafe.Pointer(&ptm)),
	)
	if errno != 0 {
		return ret, errno 
	}

	ret = time.Date(int(ptm.Year)+startYear,
		time.Month(ptm.Mon),
		int(ptm.Mday),
		int(ptm.Hour),
		int(ptm.Min),
		int(ptm.Sec),
		0,
		time.UTC)

	return ret, nil
}

func setRtcFromSysClock(rtcFile *os.File) error {
	timeUTC := time.Now().UTC()

	rtc, err := os.Open("/dev/rtc")
	if err != nil {
		return err
	}
	stm := unix.RTCTime{Sec: int32(timeUTC.Second()),
		Min:   int32(timeUTC.Minute()),
		Hour:  int32(timeUTC.Hour()),
		Mday:  int32(timeUTC.Day()),
		Mon:   int32(timeUTC.Month()),
		Year:  int32(timeUTC.Year() - startYear),
		Wday:  int32(0),
		Yday:  int32(0),
		Isdst: int32(0)}

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		rtc.Fd(),
		uintptr(unix.RTC_SET_TIME),
		uintptr(unsafe.Pointer(&stm)),
	)
	return errno
}

func main() {

	flag.Parse()

	rtcFile, err := os.Open("/dev/rtc")
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	if *write {
		if err := setRtcFromSysClock(rtcFile); err != nil {
			log.Fatalf("%v\n", err)
		}
		return
	}

	t, err := readRtc(rtcFile)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	fmt.Println(t)
}
