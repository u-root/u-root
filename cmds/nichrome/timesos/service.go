// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

type TimeService struct {
	Date string
	Time string
}

func getCurrentDate() string {
	return time.Now().Format("2006-01-02")
}

func getCurrentTime() string {
	return time.Now().Format("15:04")
}

func parseDate(d TimeJsonMsg) time.Time {
	// split date message into integers for each field
	YYYY, _ := strconv.Atoi(d.Date[:4])
	MM, _ := strconv.Atoi(d.Date[5:7])
	DD, _ := strconv.Atoi(d.Date[8:])

	hh, _ := strconv.Atoi(d.Time[:2])
	mm, _ := strconv.Atoi(d.Time[3:])

	return time.Date(YYYY, time.Month(MM), DD, hh, mm, 0, 0, time.UTC)
}

// Update sets the TimeService fields to the current system time
func (ts *TimeService) Update() {
	ts.Date = getCurrentDate()
	ts.Time = getCurrentTime()
}

// AutoSetTime calls the ntpdate u-root command to get
// the current date from time.google.com
func (ts TimeService) AutoSetTime() error {
	return exec.Command("ntpdate").Run()
}

// ManSetTime sets the system time similarly to u-root's "date" command with
// user-entered fields
func (ts TimeService) ManSetTime(new TimeJsonMsg) error {
	userTime := parseDate(new)
	tv := syscall.NsecToTimeval(userTime.UnixNano())
	if err := syscall.Settimeofday(&tv); err != nil {
		return err
	}
	return nil
}

// NewTimeService builds a TimeService with the current system date and time
func NewTimeService() (*TimeService, error) {
	return &TimeService{
		Date: getCurrentDate(),
		Time: getCurrentTime(),
	}, nil
}
