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
	tmSec   int32
	tmMin   int32
	tmHour  int32
	tmMday  int32
	tmMon   int32
	tmYear  int32
	tmWday  int32
	tmYday  int32
	tmIsdst int32
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
		return retTime, fmt.Errorf("ioctl RTC_RD_TIME returned with errno: %v", errno)
	}

	retTime = time.Date(int(ptm.tmYear)+startYear,
		time.Month(ptm.tmMon),
		int(ptm.tmMday),
		int(ptm.tmHour),
		int(ptm.tmMin),
		int(ptm.tmSec),
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
	stm := tm{tmSec: int32(timeUTC.Second()),
		tmMin:   int32(timeUTC.Minute()),
		tmHour:  int32(timeUTC.Hour()),
		tmMday:  int32(timeUTC.Day()),
		tmMon:   int32(timeUTC.Month()),
		tmYear:  int32(timeUTC.Year() - startYear),
		tmWday:  int32(0),
		tmYday:  int32(0),
		tmIsdst: int32(0)}

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
