// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os/exec"
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

// NewTimeService builds a TimeService with the current system date and time
func NewTimeService() (*TimeService, error) {
	return &TimeService{
		Date: getCurrentDate(),
		Time: getCurrentTime(),
	}, nil
}
