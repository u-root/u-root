// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/beevik/ntp"
	"github.com/u-root/u-root/pkg/rtc"
)

const (
	timestampPath string = "/etc/timestamp"
	ntpTimePool   string = "0.beevik-ntp.pool.ntp.org"
)

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
		reboot("Set system time. Need reboot.")
	}
	return nil
}
