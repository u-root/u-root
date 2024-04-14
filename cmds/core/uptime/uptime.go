// Copyright 2013-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Get the time the machine has been up
// Synopsis:
//
//	uptime
package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

var (
	errAverageFormat = errors.New("/proc/loadavg has less then 3 fields")
	errUptimeFormat  = errors.New("/proc/uptime is empty")
)

// loadavg takes in the contents of proc/loadavg,it then extracts and returns the three load averages as a string
func loadavg(contents string) (loadaverage string, err error) {
	loadavg := strings.Fields(contents)
	if len(loadavg) < 3 {
		return "", fmt.Errorf("%q: %w", contents, errAverageFormat)
	}
	return strings.Join(loadavg, ", "), nil
}

// uptime takes in the contents of proc/uptime it then extracts and returns the uptime in the format Days , Hours , Minutes ,Seconds
func uptime(contents string) (*time.Time, error) {
	uptimeArray := strings.Fields(contents)
	if len(uptimeArray) == 0 {
		return nil, errUptimeFormat
	}
	uptimeDuration, err := time.ParseDuration(string(uptimeArray[0]) + "s")
	if err != nil {
		return nil, err
	}
	uptime := time.Time{}.Add(uptimeDuration)

	return &uptime, nil
}

func run(stdout io.Writer, uptimePath, loadavgPath string) error {
	procUptimeOutput, err := os.ReadFile(uptimePath)
	if err != nil {
		return err
	}

	uptimeTime, err := uptime(string(procUptimeOutput))
	if err != nil {
		return err
	}

	procLoadAvgOutput, err := os.ReadFile(loadavgPath)
	if err != nil {
		return err
	}
	loadAverage, err := loadavg(string(procLoadAvgOutput))
	if err != nil {
		return err
	}

	// Subtracted one from time.Day() because time.Add(Duration) starts counting at 1 day instead of zero days.
	fmt.Fprintf(stdout, " %s up %d days, %d hours, %d min, loadaverage: %s\n", time.Now().Format("15:04:05"), (uptimeTime.Day() - 1), uptimeTime.Hour(), uptimeTime.Minute(), loadAverage)
	return nil
}

func main() {
	if err := run(os.Stdout, "/proc/uptime", "/proc/loadavg"); err != nil {
		log.Fatal(err)
	}
}
