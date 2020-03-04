// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
	"github.com/u-root/u-root/pkg/rtc"
)

// pollNTP queries the specified NTP server.
// On error the query is repeated infinitally.
func pollNTP() (time.Time, error) {
	bytes, err := data.get(ntpServerFile)
	if err != nil {
		reboot("Bootstrap URLs: %v", err)
	}
	var servers []string
	if err = json.Unmarshal(bytes, &servers); err != nil {
		return time.Time{}, err
	}
	for _, server := range servers {
		log.Printf("Query NTP server %s", server)
		t, err := ntp.Time(server)
		if err == nil {
			return t, nil
		}
		log.Printf("NTP error: %v", err)
	}
	//time.Sleep(3 * time.Second)
	return time.Time{}, errors.New("No NTP server resposnes")
}

// validateSystemTime sets RTC and OS time according to
// realtime clock, timestamp and ntp
func validateSystemTime(builtTime time.Time) error {
	rtc, err := rtc.OpenRTC()
	if err != nil {
		return fmt.Errorf("opening RTC failed: %v", err)
	}
	rtcTime, err := rtc.Read()
	if err != nil {
		return fmt.Errorf("reading RTC failed: %v", err)
	}

	log.Printf("Systemtime: %v", rtcTime.UTC())
	if rtcTime.UTC().Before(builtTime.UTC()) {
		log.Printf("Systemtime is invalid: %v", rtcTime.UTC())
		log.Printf("Receive time via NTP")
		ntpTime, err := pollNTP()
		if err != nil {
			return err
		}
		if ntpTime.UTC().Before(builtTime.UTC()) {
			return errors.New("NTP spoof may happened")
		}
		log.Printf("Update RTC to %v", ntpTime.UTC())
		err = rtc.Set(ntpTime)
		if err != nil {
			return fmt.Errorf("writing RTC failed: %v", err)
		}
		reboot("Set system time. Need reboot.")
	}
	return nil
}
