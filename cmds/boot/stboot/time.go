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
	"github.com/u-root/u-root/pkg/rtc"
)

const (
	timestampPath string = "/etc/timestamp"
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
	rtc, err := rtc.OpenRTC()
	if err != nil {
		return time.Time{}, err
	}
	return rtc.Read()
}

func writeRTCTime(t time.Time) error {
	rtc, err := rtc.OpenRTC()
	if err != nil {
		return err
	}
	return rtc.Set(t)
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
