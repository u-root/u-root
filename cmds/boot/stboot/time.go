// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/beevik/ntp"
	"github.com/systemboot/systemboot/pkg/recovery"
	"golang.org/x/sys/unix"
)

const (
	timestampPath string = "/etc/timestamp"
	rtcPath       string = "/dev/rtc0"
	ntpTimePool   string = "0.beevik-ntp.pool.ntp.org"
)

func getRTCMonth(month int32) (time.Month, error) {
	switch month {
	case 0:
		return time.January, nil
	case 1:
		return time.February, nil
	case 2:
		return time.March, nil
	case 3:
		return time.April, nil
	case 4:
		return time.May, nil
	case 5:
		return time.June, nil
	case 6:
		return time.July, nil
	case 7:
		return time.August, nil
	case 8:
		return time.September, nil
	case 9:
		return time.October, nil
	case 10:
		return time.November, nil
	case 11:
		return time.December, nil
	}
	return 0, fmt.Errorf("invalid rtc month representation %d", month)
}

func getRTCYear(year int32) int {
	return int(year + 1900)
}

func readRTCTime() (time.Time, error) {
	fd, err := unix.Open(rtcPath, unix.O_RDWR, 0)
	if err != nil {
		return time.Time{}, err
	}
	rtc, err := unix.IoctlGetRTCTime(fd)
	if err != nil {
		unix.Close(fd)
		return time.Time{}, err
	}
	unix.Close(fd)
	year := getRTCYear(rtc.Year)
	month, err := getRTCMonth(rtc.Mon)
	if err != nil {
		// Use mid of the year for max 6 months miscalculation
		month = time.June
	}
	localTime := time.Date(year, month, int(rtc.Mday), int(rtc.Hour), int(rtc.Min), int(rtc.Sec), 0, time.Local)
	return localTime, nil
}

func writeRTCTime(t time.Time) error {
	fd, err := unix.Open(rtcPath, unix.O_RDWR, 0)
	if err != nil {
		return err
	}
	var rtc unix.RTCTime
	rtc.Sec = int32(t.Local().Second())
	rtc.Min = int32(t.Local().Minute())
	rtc.Hour = int32(t.Local().Hour())
	rtc.Mday = int32(t.Local().Day())
	rtc.Mon = int32(t.Local().Month() - 1)
	rtc.Year = int32(t.Local().Year() - 1900)
	err = unix.IoctlSetRTCTime(fd, &rtc)
	if err != nil {
		unix.Close(fd)
		return err
	}
	unix.Close(fd)
	return nil
}

// validateSystemTime sets RTC and OS time according to
// realtime clock, timestamp and ntp
func validateSystemTime() error {
	data, err := ioutil.ReadFile(timestampPath)
	if err != nil {
		return err
	}
	unixTime, err := strconv.Atoi(strings.Trim(string(data), "\n"))
	if err != nil {
		return err
	}
	stampTime := time.Unix(int64(unixTime), 0)
	if err != nil {
		return err
	}
	rtcTime, err := readRTCTime()
	if err != nil {
		return err
	}
	log.Printf("Systemtime: %v", rtcTime.UTC())
	if rtcTime.UTC().Before(stampTime.UTC()) {
		log.Printf("Systemtime is invalid: %v", rtcTime.UTC())
		log.Printf("Receive time via NTP from %s", ntpTimePool)
		ntpTime, err := ntp.Time(ntpTimePool)
		if err != nil {
			return err
		}
		if ntpTime.UTC().Before(stampTime.UTC()) {
			return errors.New("NTP spoof may happened")
		}
		log.Printf("Update RTC to %v", ntpTime.UTC())
		err = writeRTCTime(ntpTime)
		if err != nil {
			return err
		}
		recover := recovery.SecureRecoverer{
			Reboot: true,
		}
		recover.Recover("system time update")
	}
	// tv := syscall.Timeval{
	// 	Sec:  rtcTime.Unix(),
	// 	Usec: int64(rtcTime.UnixNano() / 1000 % 1000),
	// }
	// return syscall.Settimeofday(&tv)
	return nil
}
